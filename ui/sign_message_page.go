package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/atotto/clipboard"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

const PageSignMessage = "sign_message"

type signMessagePage struct {
	container  layout.List
	wallet     *wallet.Wallet
	walletID   int
	errChannel chan error

	isPasswordModalOpen, isSigningMessage                     bool
	titleLabel, subtitleLabel, errorLabel, signedMessageLabel decredmaterial.Label
	addressEditor, messageEditor                              decredmaterial.Editor
	addressEditorW, messageEditorW                            *widget.Editor
	clearButton, signButton, copyButton                       decredmaterial.Button
	passwordModal                                             *decredmaterial.Password
	clearButtonW, signButtonW, copyButtonW                    *widget.Button
	result                                                    **wallet.Signature

	backButtonW widget.Button
	backButton  decredmaterial.IconButton
}

func (win *Window) SignMessagePage(common pageCommon) layout.Widget {
	addressEditor := common.theme.Editor("Address")
	addressEditor.IsVisible = true
	addressEditor.IsRequired = true
	messageEditor := common.theme.Editor("Message")
	messageEditor.IsVisible = true
	messageEditor.IsRequired = true
	clearButton := common.theme.Button("Clear all")
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
		addressEditorW: &widget.Editor{
			SingleLine: true,
		},

		messageEditor: messageEditor,
		messageEditorW: &widget.Editor{
			SingleLine: true,
		},

		clearButton:  clearButton,
		clearButtonW: new(widget.Button),

		signButton:  common.theme.Button("Sign"),
		signButtonW: new(widget.Button),

		copyButton:  common.theme.Button("Copy"),
		copyButtonW: new(widget.Button),
		result:      &win.signatureResult,

		backButton: common.theme.PlainIconButton(common.icons.navigationArrowBack),
	}
	pg.backButton.Color = common.theme.Color.Hint
	pg.backButton.Size = unit.Dp(32)

	return func() {
		pg.Layout(common)
		pg.handle(common)
		pg.updateColors(common)
		pg.validate(true)
	}
}

func (pg *signMessagePage) Layout(common pageCommon) {
	gtx := common.gtx
	pg.walletID = common.info.Wallets[*common.selectedWallet].ID
	pg.errChannel = common.errorChannels[PageSignMessage]

	w := []func(){
		func() {
			layout.W.Layout(gtx, func() {
				pg.backButton.Layout(gtx, &pg.backButtonW)
			})
			layout.Inset{Left: unit.Dp(44)}.Layout(gtx, func() {
				pg.titleLabel.Layout(gtx)
			})
		},
		func() {
			inset := layout.Inset{
				Top:    unit.Dp(5),
				Bottom: unit.Dp(10),
			}
			inset.Layout(gtx, func() {
				pg.subtitleLabel.Layout(gtx)
			})
		},
		func() {
			pg.errorLabel.Layout(gtx)
		},
		func() {
			pg.addressEditor.Layout(gtx, pg.addressEditorW)
		},
		func() {
			pg.messageEditor.Layout(gtx, pg.messageEditorW)
		},
		func() {
			pg.drawButtonsRow(gtx)
		},
		func() {
			pg.drawResult(gtx)
		},
	}

	common.Layout(gtx, func() {
		pg.container.Layout(gtx, len(w), func(i int) {
			layout.UniformInset(unit.Dp(5)).Layout(gtx, w[i])
		})
	})

	if pg.isPasswordModalOpen {
		pg.walletID = common.info.Wallets[*common.selectedWallet].ID
		pg.passwordModal.Layout(gtx, pg.confirm, pg.cancel)
	}
}

func (pg *signMessagePage) drawButtonsRow(gtx *layout.Context) {
	layout.E.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func() {
				inset := layout.Inset{
					Right: unit.Dp(5),
				}
				inset.Layout(gtx, func() {
					pg.clearButton.Layout(gtx, pg.clearButtonW)
				})
			}),
			layout.Rigid(func() {
				pg.signButton.Layout(gtx, pg.signButtonW)
			}),
		)
	})
}

func (pg *signMessagePage) drawResult(gtx *layout.Context) {
	if pg.signedMessageLabel.Text == "" {
		return
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			pg.signedMessageLabel.Layout(gtx)
		}),
		layout.Rigid(func() {
			pg.copyButton.Layout(gtx, pg.copyButtonW)
		}),
	)
}

func (pg *signMessagePage) updateColors(common pageCommon) {
	if pg.isSigningMessage || pg.addressEditorW.Text() == "" || pg.messageEditorW.Text() == "" {
		pg.signButton.Background = common.theme.Color.Hint
	} else {
		pg.signButton.Background = common.theme.Color.Primary
	}
}

func (pg *signMessagePage) handle(common pageCommon) {
	gtx := common.gtx
	for pg.clearButtonW.Clicked(gtx) {
		pg.clearForm()
	}

	for pg.signButtonW.Clicked(gtx) {
		if !pg.isSigningMessage && pg.validate(false) {
			pg.isPasswordModalOpen = true
		}
	}

	for pg.copyButtonW.Clicked(gtx) {
		clipboard.WriteAll(pg.signedMessageLabel.Text)
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
		pg.signButton.Text = "Sign"
	}

	if pg.backButtonW.Clicked(gtx) {
		pg.clearForm()
		*common.page = PageWallet
	}
}

func (pg *signMessagePage) confirm(password []byte) {
	pg.isPasswordModalOpen = false
	pg.isSigningMessage = true

	pg.signButton.Text = "Signing..."
	pg.wallet.SignMessage(pg.walletID, password, pg.addressEditorW.Text(), pg.messageEditorW.Text(), pg.errChannel)
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
	address := pg.addressEditorW.Text()
	pg.addressEditor.ErrorLabel.Text = ""

	if address == "" && !ignoreEmpty {
		pg.addressEditor.ErrorLabel.Text = "Please enter a valid address"
		return false
	}

	if address != "" {
		isValid, _ := pg.wallet.IsAddressValid(address)
		if !isValid {
			pg.addressEditor.ErrorLabel.Text = "Invalid address"
			return false
		}
	}
	return true
}

func (pg *signMessagePage) validateMessage(ignoreEmpty bool) bool {
	message := pg.messageEditorW.Text()
	if message == "" && !ignoreEmpty {
		pg.messageEditor.ErrorLabel.Text = "Please enter a message to sign"
		return false
	}
	return true
}

func (pg *signMessagePage) clearForm() {
	pg.addressEditorW.SetText("")
	pg.messageEditorW.SetText("")
	pg.errorLabel.Text = ""
}
