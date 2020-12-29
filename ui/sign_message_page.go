package ui

import (
	"fmt"

	"github.com/planetdecred/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/wallet"
)

const PageSignMessage = "SignMessage"

type signMessagePage struct {
	container  layout.List
	wallet     *wallet.Wallet
	walletID   int
	errChannel chan error

	isPasswordModalOpen, isSigningMessage                     bool
	titleLabel, subtitleLabel, errorLabel, signedMessageLabel decredmaterial.Label
	addressEditor, messageEditor                              decredmaterial.Editor
	clearButton, signButton, copyButton                       decredmaterial.Button
	passwordModal                                             *decredmaterial.Password
	result                                                    **wallet.Signature

	backButton decredmaterial.IconButton
}

func (win *Window) SignMessagePage(common pageCommon) layout.Widget {
	addressEditor := common.theme.Editor(new(widget.Editor), "Address")
	addressEditor.IsVisible = true
	addressEditor.IsRequired = true
	addressEditor.Editor.SingleLine = true
	messageEditor := common.theme.Editor(new(widget.Editor), "Message")
	messageEditor.IsVisible = true
	messageEditor.IsRequired = true
	messageEditor.Editor.SingleLine = true
	clearButton := common.theme.Button(new(widget.Clickable), "Clear all")
	clearButton.Background = common.theme.Color.Background
	clearButton.Color = common.theme.Color.Gray
	errorLabel := common.theme.Caption("")
	errorLabel.Color = common.theme.Color.Danger

	pg := &signMessagePage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		wallet:        common.wallet,
		passwordModal: common.theme.Password(),

		titleLabel:         common.theme.H5("Sign Message"),
		subtitleLabel:      common.theme.Body2("Enter the address and message you want to sign"),
		signedMessageLabel: common.theme.H5(""),
		errorLabel:         errorLabel,
		addressEditor:      addressEditor,
		messageEditor:      messageEditor,

		clearButton: clearButton,
		signButton:  common.theme.Button(new(widget.Clickable), "Sign message"),
		copyButton:  common.theme.Button(new(widget.Clickable), "Copy"),
		result:      &win.signatureResult,

		backButton: common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
	}
	pg.backButton.Color = common.theme.Color.Text
	pg.backButton.Inset = layout.UniformInset(values.MarginPadding0)

	return func(gtx C) D {
		pg.handle(common)
		pg.updateColors(common)
		pg.validate(true)
		return pg.Layout(gtx, common)
	}
}

func (pg *signMessagePage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.walletID = common.info.Wallets[*common.selectedWallet].ID
	pg.errChannel = common.errorChannels[PageSignMessage]

	w := []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.NW.Layout(gtx, func(gtx C) D {
						return pg.backButton.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return pg.titleLabel.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			inset := layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding10,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return pg.subtitleLabel.Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.errorLabel.Layout(gtx)
		},
		func(gtx C) D {
			return pg.addressEditor.Layout(gtx)
		},
		func(gtx C) D {
			return pg.messageEditor.Layout(gtx)
		},
		func(gtx C) D {
			return pg.drawButtonsRow(gtx)
		},
		func(gtx C) D {
			return pg.drawResult(gtx)
		},
	}

	body := common.Layout(gtx, func(gtx C) D {
		return common.theme.Card().Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
				return pg.container.Layout(gtx, len(w), func(gtx C, i int) D {
					return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
						return w[i](gtx)
					})
				})
			})
		})
	})

	if pg.isPasswordModalOpen {
		pg.walletID = common.info.Wallets[*common.selectedWallet].ID
		return common.Modal(gtx, body, pg.passwordModal.Layout(gtx, pg.confirm, pg.cancel))
	}

	return body
}

func (pg *signMessagePage) drawButtonsRow(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						inset := layout.Inset{
							Right: values.MarginPadding5,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return pg.clearButton.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return pg.signButton.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *signMessagePage) drawResult(gtx layout.Context) layout.Dimensions {
	if pg.signedMessageLabel.Text == "" {
		return layout.Dimensions{}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.signedMessageLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.copyButton.Layout(gtx)
		}),
	)
}

func (pg *signMessagePage) updateColors(common pageCommon) {
	if pg.isSigningMessage || pg.addressEditor.Editor.Text() == "" || pg.messageEditor.Editor.Text() == "" {
		pg.signButton.Background = common.theme.Color.Hint
	} else {
		pg.signButton.Background = common.theme.Color.Primary
	}
}

func (pg *signMessagePage) handle(common pageCommon) {
	for pg.clearButton.Button.Clicked() {
		pg.clearForm()
	}

	for pg.signButton.Button.Clicked() {
		if !pg.isSigningMessage && pg.validate(false) {
			pg.isPasswordModalOpen = true
		}
	}

	if pg.copyButton.Button.Clicked() {
		go func() {
			common.clipboard <- WriteClipboard{Text: pg.signedMessageLabel.Text}
		}()
	}

	select {
	case err := <-pg.errChannel:
		fmt.Printf("SIGNMESSAGE PAGE ERROR! %v", err)
	default:
	}

	if *pg.result != nil {
		if (*pg.result).Err != nil {
			pg.errorLabel.Text = (*pg.result).Err.Error()
		} else {
			pg.signedMessageLabel.Text = (*pg.result).Signature
		}
		*pg.result = nil
		pg.isSigningMessage = false
		pg.signButton.Text = "Sign message"
	}

	if pg.backButton.Button.Clicked() {
		pg.clearForm()
		*common.page = PageWallet
	}
}

func (pg *signMessagePage) confirm(password []byte) {
	pg.isPasswordModalOpen = false
	pg.isSigningMessage = true

	pg.signButton.Text = "Signing..."
	pg.wallet.SignMessage(pg.walletID, password, pg.addressEditor.Editor.Text(), pg.messageEditor.Editor.Text(), pg.errChannel)
}

func (pg *signMessagePage) cancel() {
	pg.isPasswordModalOpen = false
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
	}
	return true
}

func (pg *signMessagePage) validateMessage(ignoreEmpty bool) bool {
	message := pg.messageEditor.Editor.Text()
	if message == "" && !ignoreEmpty {
		pg.messageEditor.SetError("Please enter a message to sign")
		return false
	}
	return true
}

func (pg *signMessagePage) clearForm() {
	pg.addressEditor.Editor.SetText("")
	pg.messageEditor.Editor.SetText("")
	pg.errorLabel.Text = ""
}
