package page

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/dcrlibwallet"
)

const VerifyMessagePageID = "VerifyMessage"

type VerifyMessagePage struct {
	*load.Load

	addressEditor          decredmaterial.Editor
	messageEditor          decredmaterial.Editor
	signatureEditor        decredmaterial.Editor
	clearBtn, verifyButton decredmaterial.Button
	verifyMessage          decredmaterial.Label

	verifyMessageStatus *widget.Icon

	backButton     decredmaterial.IconButton
	infoButton     decredmaterial.IconButton
	addressIsValid bool
}

func NewVerifyMessagePage(l *load.Load) *VerifyMessagePage {
	pg := &VerifyMessagePage{
		Load:          l,
		verifyMessage: l.Theme.Body1(""),
	}

	pg.addressEditor = l.Theme.Editor(new(widget.Editor), "Address")
	pg.addressEditor.Editor.SingleLine = true
	pg.addressEditor.Editor.Submit = true

	pg.messageEditor = l.Theme.Editor(new(widget.Editor), "Message")
	pg.messageEditor.Editor.SingleLine = true
	pg.messageEditor.Editor.Submit = true

	pg.signatureEditor = l.Theme.Editor(new(widget.Editor), "Signature")
	pg.signatureEditor.Editor.Submit = true

	buttonTextSize := values.TextSize14
	pg.verifyButton = l.Theme.Button(new(widget.Clickable), "Verify message")
	pg.verifyButton.TextSize = buttonTextSize
	pg.verifyButton.Font.Weight = text.Bold
	pg.verifyButton.Background = l.Theme.Color.Hint

	pg.clearBtn = l.Theme.Button(new(widget.Clickable), "Clear all")
	pg.clearBtn.TextSize = buttonTextSize
	pg.clearBtn.Background = color.NRGBA{}
	pg.clearBtn.Color = l.Theme.Color.Hint
	pg.clearBtn.Font.Weight = text.Bold

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

func (pg *VerifyMessagePage) OnResume() {

}

func (pg *VerifyMessagePage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Verify message",
			BackButton: pg.backButton,
			InfoButton: pg.infoButton,
			Back: func() {
				pg.clearInputs()
				pg.ChangePage(WalletPageID)
				pg.ChangePage(*pg.ReturnPage)
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(pg.description()),
							layout.Rigid(pg.inputRow(pg.addressEditor)),
							layout.Rigid(pg.inputRow(pg.signatureEditor)),
							layout.Rigid(pg.inputRow(pg.messageEditor)),
							layout.Rigid(pg.verifyAndClearButtons()),
							layout.Rigid(pg.verifyMessageResponse()),
						)
					})
				})
			},
			InfoTemplate: modal.VerifyMessageInfoTemplate,
		}
		return sp.Layout(gtx)
	}
	return components.UniformPadding(gtx, body)
}

func (pg *VerifyMessagePage) inputRow(editor decredmaterial.Editor) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, editor.Layout)
	}
}

func (pg *VerifyMessagePage) description() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		desc := pg.Theme.Caption("Enter the address, signature, and message to verify:")
		desc.Color = pg.Theme.Color.Gray
		return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, desc.Layout)
	}
}

func (pg *VerifyMessagePage) verifyAndClearButtons() layout.Widget {
	return func(gtx C) D {
		dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, pg.clearBtn.Layout)
						}),
						layout.Rigid(pg.verifyButton.Layout),
					)
				})
			}),
		)
		return dims
	}
}

func (pg *VerifyMessagePage) verifyMessageResponse() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if pg.verifyMessageStatus != nil {
			return layout.Inset{Top: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
				pg.Theme.Separator().Layout(gtx)
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return pg.verifyMessageStatus.Layout(gtx, values.MarginPadding20)
							})
						}),
						layout.Rigid(pg.verifyMessage.Layout),
					)
				})
			})
		}
		return layout.Dimensions{}
	}
}

func (pg *VerifyMessagePage) Handle() {
	pg.updateColors()

	for _, evt := range pg.addressEditor.Editor.Events() {
		if pg.addressEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				pg.validateAddress()
			}
		}
	}

	if pg.verifyButton.Button.Clicked() || handleSubmitEvent(pg.addressEditor.Editor, pg.messageEditor.Editor, pg.signatureEditor.Editor) {
		if pg.validateAllInputs() {
			fmt.Println("im here")
			pg.verifyMessage.Text = ""
			pg.verifyMessageStatus = nil
			valid, err := pg.WL.MultiWallet.VerifyMessage(pg.addressEditor.Editor.Text(), pg.messageEditor.Editor.Text(), pg.signatureEditor.Editor.Text())
			if err != nil || !valid {
				pg.verifyMessage.Text = "Invalid signature or message"
				pg.verifyMessage.Color = pg.Theme.Color.Danger
				pg.verifyMessageStatus = pg.Icons.NavigationCancel
				return
			}

			pg.verifyMessageStatus = pg.Icons.ActionCheck
			pg.verifyMessageStatus.Color = pg.Theme.Color.Success
			pg.verifyMessage.Text = "Valid signature"
			pg.verifyMessage.Color = pg.Theme.Color.Success
		}
	}

	if pg.clearBtn.Button.Clicked() {
		pg.clearInputs()
	}
}
func (pg *VerifyMessagePage) validateAllInputs() bool {
	if !pg.validateAddress() || !pg.stringNotEmpty(pg.messageEditor.Editor.Text(), pg.signatureEditor.Editor.Text()) {
		return false
	}
	return true
}

func (pg *VerifyMessagePage) updateColors() {
	pg.clearBtn.Color, pg.verifyButton.Background = pg.Theme.Color.Hint, pg.Theme.Color.Hint
	if pg.stringNotEmpty(pg.addressEditor.Editor.Text()) ||
	 pg.stringNotEmpty(pg.messageEditor.Editor.Text()) ||
		pg.stringNotEmpty(pg.signatureEditor.Editor.Text()) {
		pg.clearBtn.Color = pg.Theme.Color.Primary
	}
	if pg.addressIsValid && pg.stringNotEmpty(pg.messageEditor.Editor.Text(), pg.signatureEditor.Editor.Text()) {
		pg.clearBtn.Color, pg.verifyButton.Background = pg.Theme.Color.Primary, pg.Theme.Color.Primary
	}
	return true
}

func (pg *VerifyMessagePage) updateColors() {
	if pg.stringNotEmpty(pg.addressEditor.Editor.Text()) || pg.stringNotEmpty(pg.messageEditor.Editor.Text()) ||
	 pg.stringNotEmpty(pg.signatureEditor.Editor.Text()) {
		pg.clearBtn.Color = pg.Theme.Color.Primary
	} 
}

func (pg *VerifyMessagePage) clearInputs() {
	pg.verifyMessageStatus = nil
	pg.verifyButton.Background = pg.Theme.Color.Hint
	pg.addressEditor.Editor.SetText("")
	pg.signatureEditor.Editor.SetText("")
	pg.messageEditor.Editor.SetText("")
	pg.verifyMessage.Text = ""
	pg.addressEditor.SetError("")
}

func (pg *VerifyMessagePage) validateAddress() bool {
	address := pg.addressEditor.Editor.Text()
	pg.addressEditor.SetError("")
	exist, _ := pg.WL.Wallet.HaveAddress(address)

	var valid bool
	var errorMessage string

	switch {
	case !pg.stringNotEmpty(address):
		errorMessage, valid = "Please enter a valid address", false
	case !pg.WL.MultiWallet.IsAddressValid(address):
		errorMessage, valid = "Invalid address", false
	case !exist:
		errorMessage, valid = "Address not owned by any wallet", false
	default:
		errorMessage, valid = "", true
	}
	if !valid {
		pg.addressEditor.SetError(errorMessage)
	}

	pg.addressIsValid = valid
	return valid
}

func (pg *VerifyMessagePage) stringNotEmpty(texts ...string) bool {
	for _, t := range texts {
		if strings.Trim(t, " ") == "" {
			return false
		}
	}

	return true
}

func (pg *VerifyMessagePage) OnClose() {}
