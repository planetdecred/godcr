package ui

import (
	"image/color"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageSignMessage = "SignMessage"

type signMessagePage struct {
	common        *pageCommon
	theme         *decredmaterial.Theme
	container     layout.List
	wallet        *wallet.Wallet
	walletID      int
	errorReceiver chan error

	isSigningMessage                           bool
	titleLabel, errorLabel, signedMessageLabel decredmaterial.Label
	addressEditor, messageEditor               decredmaterial.Editor
	clearButton, signButton, copyButton        decredmaterial.Button
	result                                     **wallet.Signature
	copySignature                              *widget.Clickable
	copyIcon                                   *widget.Image
	gtx                                        *layout.Context
}

func (win *Window) SignMessagePage(common *pageCommon) Page {
	addressEditor := common.theme.Editor(new(widget.Editor), "Address")
	addressEditor.Editor.SingleLine, addressEditor.Editor.Submit = true, true
	messageEditor := common.theme.Editor(new(widget.Editor), "Message")
	messageEditor.Editor.SingleLine, messageEditor.Editor.Submit = true, true
	clearButton := common.theme.Button(new(widget.Clickable), "Clear all")
	clearButton.Background = color.NRGBA{}
	clearButton.Color = common.theme.Color.Gray
	errorLabel := common.theme.Caption("")
	errorLabel.Color = common.theme.Color.Danger
	copyIcon := common.icons.copyIcon

	pg := &signMessagePage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		common: common,
		theme:  common.theme,
		wallet: common.wallet,

		titleLabel:         common.theme.H5("Sign Message"),
		signedMessageLabel: common.theme.Body1(""),
		errorLabel:         errorLabel,
		addressEditor:      addressEditor,
		messageEditor:      messageEditor,

		clearButton:   clearButton,
		signButton:    common.theme.Button(new(widget.Clickable), "Sign message"),
		copyButton:    common.theme.Button(new(widget.Clickable), "Copy"),
		result:        &win.signatureResult,
		copySignature: new(widget.Clickable),
		copyIcon:      copyIcon,
		errorReceiver: make(chan error),
	}

	pg.signedMessageLabel.Color = common.theme.Color.Gray

	return pg
}

func (pg *signMessagePage) Layout(gtx layout.Context) layout.Dimensions {
	if pg.gtx == nil {
		pg.gtx = &gtx
	}
	common := pg.common
	pg.walletID = common.info.Wallets[*common.selectedWallet].ID

	body := func(gtx C) D {
		page := SubPage{
			title:      "Sign message",
			walletName: common.info.Wallets[*common.selectedWallet].Name,
			back: func() {
				pg.clearForm()
				common.changePage(PageWallet)
			},
			body: func(gtx layout.Context) layout.Dimensions {
				return common.theme.Card().Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(pg.description()),
							layout.Rigid(pg.editors(pg.addressEditor)),
							layout.Rigid(pg.editors(pg.messageEditor)),
							layout.Rigid(pg.drawButtonsRow()),
							layout.Rigid(pg.drawResult()),
						)
					})
				})
			},
			infoTemplate: SignMessageInfoTemplate,
		}
		return common.SubPageLayout(gtx, page)
	}

	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, body)
	})
}

func (pg *signMessagePage) description() layout.Widget {
	return func(gtx C) D {
		desc := pg.theme.Caption("Enter an address and message to sign:")
		desc.Color = pg.theme.Color.Gray
		return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, desc.Layout)
	}
}

func (pg *signMessagePage) editors(editor decredmaterial.Editor) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, editor.Layout)
	}
}

func (pg *signMessagePage) drawButtonsRow() layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Right: values.MarginPadding5,
							}
							return inset.Layout(gtx, pg.clearButton.Layout)
						}),
						layout.Rigid(pg.signButton.Layout),
					)
				})
			}),
		)
	}
}

func (pg *signMessagePage) drawResult() layout.Widget {
	return func(gtx C) D {
		if pg.signedMessageLabel.Text == "" {
			return layout.Dimensions{}
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				m := values.MarginPadding30
				return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.theme.Separator().Layout)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						border := widget.Border{Color: pg.theme.Color.LightGray, CornerRadius: values.MarginPadding10, Width: values.MarginPadding2}
						return border.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Flexed(0.9, pg.signedMessageLabel.Layout),
									layout.Flexed(0.1, func(gtx C) D {
										return layout.E.Layout(gtx, func(gtx C) D {
											return layout.Inset{Top: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
												pg.copyIcon.Scale = 1.0
												return decredmaterial.Clickable(gtx, pg.copySignature, pg.copyIcon.Layout)
											})
										})
									}),
								)
							})
						})
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:  values.MarginPaddingMinus10,
							Left: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							return pg.theme.Card().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								label := pg.theme.Body1("Signature")
								label.Color = pg.theme.Color.Gray
								return label.Layout(gtx)
							})
						})
					}),
				)
			}),
		)
	}
}

func (pg *signMessagePage) updateColors(common *pageCommon) {
	if pg.isSigningMessage || pg.addressEditor.Editor.Text() == "" || pg.messageEditor.Editor.Text() == "" {
		pg.signButton.Background = common.theme.Color.Hint
	} else {
		pg.signButton.Background = common.theme.Color.Primary
	}
}

func (pg *signMessagePage) handle() {
	gtx := pg.gtx
	common := pg.common
	pg.updateColors(common)
	pg.validate(true)

	for pg.clearButton.Button.Clicked() {
		pg.clearForm()
	}

	for pg.signButton.Button.Clicked() || handleSubmitEvent(pg.addressEditor.Editor, pg.messageEditor.Editor) {
		if !pg.isSigningMessage && pg.validate(false) {
			address := pg.addressEditor.Editor.Text()
			message := pg.messageEditor.Editor.Text()

			go func() {
				common.modalReceiver <- &modalLoad{
					template: PasswordTemplate,
					title:    "Confirm to sign",
					confirm: func(pass string) {
						pg.wallet.SignMessage(pg.walletID, []byte(pass), address, message, pg.errorReceiver)
					},
					confirmText: "Confirm",
					cancel:      common.closeModal,
					cancelText:  "Cancel",
				}
			}()
		}
	}

	if pg.copySignature.Clicked() {
		clipboard.WriteOp{Text: pg.signedMessageLabel.Text}.Add(gtx.Ops)
	}

	if *pg.result != nil {
		pg.signedMessageLabel.Text = (*pg.result).Signature
		*pg.result = nil
		pg.isSigningMessage = false
	}

	select {
	case err := <-pg.errorReceiver:
		common.modalLoad.setLoading(false)
		common.notify(err.Error(), false)
		if err.Error() != dcrlibwallet.ErrInvalidPassphrase {
			common.closeModal()
		}
	default:
	}
}

func (pg *signMessagePage) validate(ignoreEmpty bool) bool {
	isAddressValid := pg.validateAddress(ignoreEmpty)
	isMessageValid := pg.validateMessage(ignoreEmpty)
	if !isAddressValid || !isMessageValid {
		return false
	}
	return true
}

func (pg *signMessagePage) validateAddress(ignoreEmpty bool) bool {
	address := pg.addressEditor.Editor.Text()
	pg.addressEditor.SetError("")

	if address == "" && !ignoreEmpty {
		pg.addressEditor.SetError("Please enter a valid address")
		return false
	}

	if address != "" {
		isValid, _ := pg.wallet.IsAddressValid(address)
		if !isValid {
			pg.addressEditor.SetError("Invalid address")
			return false
		}

		exist, err := pg.wallet.HaveAddress(pg.walletID, address)
		if err != nil {
			return false
		}
		if !exist {
			pg.addressEditor.SetError("Address not owned by this wallet")
			return false
		}
	}
	return true
}

func (pg *signMessagePage) validateMessage(ignoreEmpty bool) bool {
	message := pg.messageEditor.Editor.Text()
	if message == "" && !ignoreEmpty {
		return false
	}
	return true
}

func (pg *signMessagePage) clearForm() {
	pg.addressEditor.Editor.SetText("")
	pg.messageEditor.Editor.SetText("")
	pg.signedMessageLabel.Text = ""
	pg.errorLabel.Text = ""
}

func (pg *signMessagePage) onClose() {}
