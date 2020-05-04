package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/atotto/clipboard"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

const PageSignMessage = "sign_message"

type SignMessagePage struct {
	theme    *decredmaterial.Theme
	wallet   *wallet.Wallet
	walletID int

	isPasswordModalOpen bool
	isSigningMessage    bool

	titleLabel         decredmaterial.Label
	subtitleLabel      decredmaterial.Label
	addressErrorLabel  decredmaterial.Label
	messageErrorLabel  decredmaterial.Label
	errorLabel         decredmaterial.Label
	signedMessageLabel decredmaterial.Label

	addressEditorMaterial decredmaterial.Editor
	messageEditorMaterial decredmaterial.Editor

	addressEditorWidget *widget.Editor
	messageEditorWidget *widget.Editor

	clearButtonMaterial          decredmaterial.Button
	signButtonMaterial           decredmaterial.Button
	copyButtonMaterial           decredmaterial.Button
	pasteInAddressButtonMaterial decredmaterial.Button
	pasteInMessageButtonMaterial decredmaterial.Button

	passwordModal *decredmaterial.Password

	clearButtonWidget          *widget.Button
	signButtonWidget           *widget.Button
	copyButtonWidget           *widget.Button
	pasteInAddressButtonWidget *widget.Button
	pasteInMessageButtonWidget *widget.Button
}

var signMessagePage *SignMessagePage

const (
	editorWidthRatio = 0.99
)

func (pg *SignMessagePage) Draw(gtx *layout.Context) {
	pg.handleEvents(gtx)
	pg.updateColors()
	pg.validate(true)

	w := []func(){
		func() {
			pg.titleLabel.Layout(gtx)
		},
		func() {
			inset := layout.Inset{
				Top:    unit.Dp(5),
				Bottom: unit.Dp(15),
			}
			inset.Layout(gtx, func() {
				pg.subtitleLabel.Layout(gtx)
			})
		},
		func() {
			pg.errorLabel.Layout(gtx)
		},
		func() {
			pg.drawAddressEditor(gtx)
		},
		func() {
			pg.drawMessageEditor(gtx)
		},
		func() {
			pg.drawButtonsRow(gtx)
		},
		func() {
			pg.drawResult(gtx)
		},
	}

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(w), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, w[i])
	})

	if pg.isPasswordModalOpen {
		pg.passwordModal.Layout(gtx, pg.confirm, pg.cancel)
	}
}

func (pg *SignMessagePage) drawAddressEditor(gtx *layout.Context) {
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(editorWidthRatio, func() {
			pg.addressEditorMaterial.Layout(gtx, pg.addressEditorWidget)
		}),
		layout.Rigid(func() {
			pg.pasteInAddressButtonMaterial.Layout(gtx, pg.pasteInAddressButtonWidget)
		}),
	)

	if pg.addressErrorLabel.Text != "" {
		inset := layout.Inset{
			Top: unit.Dp(25),
		}
		inset.Layout(gtx, func() {
			pg.addressErrorLabel.Layout(gtx)
		})
	}
}

func (pg *SignMessagePage) drawMessageEditor(gtx *layout.Context) {
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(editorWidthRatio, func() {
			pg.messageEditorMaterial.Layout(gtx, pg.messageEditorWidget)
		}),
		layout.Rigid(func() {
			pg.pasteInMessageButtonMaterial.Layout(gtx, pg.pasteInMessageButtonWidget)
		}),
	)
	if pg.messageErrorLabel.Text != "" {
		inset := layout.Inset{
			Top: unit.Dp(25),
		}
		inset.Layout(gtx, func() {
			pg.messageErrorLabel.Layout(gtx)
		})
	}
}

func (pg *SignMessagePage) drawButtonsRow(gtx *layout.Context) {
	layout.E.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func() {
				inset := layout.Inset{
					Right: unit.Dp(5),
				}
				inset.Layout(gtx, func() {
					pg.clearButtonMaterial.Layout(gtx, pg.clearButtonWidget)
				})
			}),
			layout.Rigid(func() {
				pg.signButtonMaterial.Layout(gtx, pg.signButtonWidget)
			}),
		)
	})
}

func (pg *SignMessagePage) drawResult(gtx *layout.Context) {
	if pg.signedMessageLabel.Text == "" {
		return
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			pg.signedMessageLabel.Layout(gtx)
		}),
		layout.Rigid(func() {
			pg.copyButtonMaterial.Layout(gtx, pg.copyButtonWidget)
		}),
	)
}

func (pg *SignMessagePage) updateColors() {
	if pg.isSigningMessage || pg.addressEditorWidget.Text() == "" || pg.messageEditorWidget.Text() == "" {
		pg.signButtonMaterial.Background = pg.theme.Color.Hint
	} else {
		pg.signButtonMaterial.Background = pg.theme.Color.Primary
	}
}

func (pg *SignMessagePage) handleEvents(gtx *layout.Context) {
	for pg.clearButtonWidget.Clicked(gtx) {
		pg.clearForm()
	}

	for pg.signButtonWidget.Clicked(gtx) {
		if !pg.isSigningMessage && pg.validate(false) {
			pg.isPasswordModalOpen = true
		}
	}

	for pg.pasteInAddressButtonWidget.Clicked(gtx) {
		pg.addressEditorWidget.Insert(GetClipboardContent())
	}

	for pg.pasteInMessageButtonWidget.Clicked(gtx) {
		pg.messageEditorWidget.Insert(GetClipboardContent())
	}

	for pg.copyButtonWidget.Clicked(gtx) {
		clipboard.WriteAll(pg.signedMessageLabel.Text)
	}
}

func (pg *SignMessagePage) confirm(password []byte) {
	pg.isPasswordModalOpen = false
	pg.isSigningMessage = true

	pg.signButtonMaterial.Text = "Signing..."
	pg.wallet.SignMessage(pg.walletID, password, pg.addressEditorWidget.Text(), pg.messageEditorWidget.Text())
}

func (pg *SignMessagePage) cancel() {
	pg.isPasswordModalOpen = false
}

func (pg *SignMessagePage) validate(ignoreEmpty bool) bool {
	isAddressValid := pg.validateAddress(ignoreEmpty)
	isMessageValid := pg.validateMessage(ignoreEmpty)

	if !isAddressValid || !isMessageValid {
		return false
	}
	return true
}

func (pg *SignMessagePage) validateAddress(ignoreEmpty bool) bool {
	pg.addressErrorLabel.Text = ""
	address := pg.addressEditorWidget.Text()

	if address == "" && !ignoreEmpty {
		pg.addressErrorLabel.Text = "please enter a valid address"
		return false
	}

	if address != "" {
		isValid, _ := pg.wallet.IsAddressValid(address)
		if !isValid {
			pg.addressErrorLabel.Text = "invalid address"
			return false
		}
	}
	return true
}

func (pg *SignMessagePage) validateMessage(ignoreEmpty bool) bool {
	message := pg.messageEditorWidget.Text()
	if message == "" && !ignoreEmpty {
		pg.messageErrorLabel.Text = "please enter a message to sign"
		return false
	}
	return true
}

func (pg *SignMessagePage) clearForm() {
	pg.addressEditorWidget.SetText("")
	pg.messageEditorWidget.SetText("")
	pg.errorLabel.Text = ""
}

func (win *Window) newSignMessagePage() *SignMessagePage {
	pg := &SignMessagePage{}
	pg.theme = win.theme
	pg.wallet = win.wallet

	pg.passwordModal = pg.theme.Password()

	pg.titleLabel = pg.theme.H5("Sign Message")
	pg.subtitleLabel = pg.theme.Body2("Enter the address and message you want to sign")
	pg.errorLabel = pg.theme.Caption("")
	pg.addressErrorLabel = pg.theme.Caption("")
	pg.signedMessageLabel = pg.theme.H5("")
	pg.messageErrorLabel = pg.theme.Caption("")

	pg.addressEditorMaterial = pg.theme.Editor("Address")
	pg.addressEditorWidget = &widget.Editor{
		SingleLine: true,
	}

	pg.messageEditorMaterial = pg.theme.Editor("Message")
	pg.messageEditorWidget = &widget.Editor{
		SingleLine: true,
	}

	pg.clearButtonMaterial = pg.theme.Button("Clear all")
	pg.clearButtonWidget = new(widget.Button)

	pg.signButtonMaterial = pg.theme.Button("Sign")
	pg.signButtonWidget = new(widget.Button)

	pg.pasteInAddressButtonMaterial = pg.theme.Button("Paste")
	pg.pasteInAddressButtonWidget = new(widget.Button)

	pg.copyButtonMaterial = pg.theme.Button("Copy")
	pg.copyButtonWidget = new(widget.Button)

	pg.pasteInMessageButtonMaterial = pg.theme.Button("Paste")
	pg.pasteInMessageButtonWidget = new(widget.Button)

	pg.pasteInMessageButtonMaterial.Background = pg.theme.Color.Surface
	pg.pasteInMessageButtonMaterial.Color = pg.theme.Color.Primary
	pg.pasteInAddressButtonMaterial.Background = pg.theme.Color.Surface
	pg.pasteInAddressButtonMaterial.Color = pg.theme.Color.Primary
	pg.clearButtonMaterial.Background = pg.theme.Color.Background
	pg.clearButtonMaterial.Color = pg.theme.Color.Gray
	pg.addressErrorLabel.Color = pg.theme.Color.Danger
	pg.errorLabel.Color = pg.theme.Color.Danger
	pg.messageErrorLabel.Color = pg.theme.Color.Danger

	return pg
}

func (win *Window) SignMessagePage() {
	if signMessagePage == nil {
		signMessagePage = win.newSignMessagePage()
	}
	signMessagePage.walletID = win.walletInfo.Wallets[win.selected].ID

	if win.signatureResult != nil {
		if win.signatureResult.Err != nil {
			signMessagePage.errorLabel.Text = win.signatureResult.Err.Error()
		} else {
			signMessagePage.signedMessageLabel.Text = win.signatureResult.Signature
		}
		win.signatureResult = nil
		signMessagePage.isSigningMessage = false
		signMessagePage.signButtonMaterial.Text = "Sign"
	}

	body := func() {
		layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
			signMessagePage.Draw(win.gtx)
		})
	}
	win.Page(body)
}
