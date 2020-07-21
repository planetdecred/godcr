package ui

import (
	"image/color"
	"strings"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageVerifyMessage = "verifymessage"

type verifyMessagePage struct {
	addressInput, messageInput, signInput decredmaterial.Editor
	clearBtn, verifyBtn                   decredmaterial.Button
	verifyMessage                         decredmaterial.Label

	backButton decredmaterial.IconButton
}

func (win *Window) VerifyMessagePage(c pageCommon) layout.Widget {
	pg := &verifyMessagePage{
		addressInput:  c.theme.Editor(new(widget.Editor), "Address"),
		messageInput:  c.theme.Editor(new(widget.Editor), "Message"),
		signInput:     c.theme.Editor(new(widget.Editor), "Signature"),
		verifyMessage: c.theme.H6(""),
		verifyBtn:     c.theme.Button(new(widget.Clickable), "Verify"),
		clearBtn:      c.theme.Button(new(widget.Clickable), "Clear All"),
		backButton:    c.theme.PlainIconButton(new(widget.Clickable), c.icons.navigationArrowBack),
	}

	pg.messageInput.IsRequired, pg.addressInput.IsRequired, pg.signInput.IsRequired = true, true, true
	pg.messageInput.IsVisible, pg.addressInput.IsVisible, pg.signInput.IsVisible = true, true, true
	pg.messageInput.Editor.SingleLine, pg.addressInput.Editor.SingleLine, pg.messageInput.Editor.SingleLine = true, true, true
	pg.verifyBtn.TextSize, pg.clearBtn.TextSize, pg.clearBtn.TextSize = values.TextSize14, values.TextSize14, values.TextSize14
	pg.verifyBtn.Background = c.theme.Color.Hint
	pg.clearBtn.Background = color.RGBA{0, 0, 0, 0}
	pg.clearBtn.Color = c.theme.Color.Primary
	pg.backButton.Color = c.theme.Color.Hint
	pg.backButton.Size = values.MarginPadding30
	pg.backButton.Inset = layout.UniformInset(values.MarginPadding0)

	return func(gtx C) D {
		pg.handler(c)
		return pg.Layout(gtx, c)
	}
}

func (pg *verifyMessagePage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.header(&c)),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(pg.inputRow(pg.addressInput)),
						layout.Rigid(pg.inputRow(pg.signInput)),
						layout.Rigid(pg.inputRow(pg.messageInput)),
						layout.Rigid(pg.verifyAndClearButtons()),
					)
				}),
			)
		})
	}
	return c.Layout(gtx, body)
}

func (pg *verifyMessagePage) header(c *pageCommon) layout.Widget {
	return func(gtx C) D {
		var msg = "After you or your counterparty has genrated a signature, you can use this form to verify the signature." +
			"\nOnce you have entered the address, the message and the corresponding signature, you will see VALID if the signature" +
			"appropriately matches \nthe address and message otherwise INVALID."

		return layout.Inset{Bottom: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.W.Layout(gtx, func(gtx C) D {
								return pg.backButton.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding45}.Layout(gtx, func(gtx C) D {
								return c.theme.H5("Verify Wallet Message").Layout(gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return layout.Dimensions{}
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := c.theme.Label(values.MarginPadding10, msg)
					return txt.Layout(gtx)
				}),
			)
		})
	}
}

func (pg *verifyMessagePage) inputRow(editor decredmaterial.Editor) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
			return editor.Layout(gtx)
		})
	}
}

func (pg *verifyMessagePage) verifyAndClearButtons() layout.Widget {
	return func(gtx C) D {
		dims := layout.Flex{}.Layout(gtx,
			layout.Flexed(.6, func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding5, Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
					return pg.verifyMessage.Layout(gtx)
				})
			}),
			layout.Flexed(.4, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(.5, func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return pg.clearBtn.Layout(gtx)
						})
					}),
					layout.Flexed(.5, func(gtx C) D {
						return pg.verifyBtn.Layout(gtx)
					}),
				)
			}),
		)
		return dims
	}
}

func (pg *verifyMessagePage) validateInputs(c *pageCommon) bool {
	pg.addressInput.ErrorLabel.Text = ""
	pg.verifyBtn.Background = c.theme.Color.Hint

	if strings.Trim(pg.addressInput.Editor.Text(), " ") == "" {
		pg.addressInput.ErrorLabel.Text = "Please enter a valid address"
		return false
	}
	if isValid, _ := c.wallet.IsAddressValid(pg.addressInput.Editor.Text()); !isValid {
		pg.addressInput.ErrorLabel.Text = "Invalid address"
		return false
	}
	if strings.Trim(pg.messageInput.Editor.Text(), " ") == "" {
		return false
	}
	if strings.Trim(pg.signInput.Editor.Text(), " ") == "" {
		return false
	}

	pg.verifyBtn.Background = c.theme.Color.Primary
	return true
}

func (pg *verifyMessagePage) handlerEditorEvents(c *pageCommon, w *widget.Editor) {
	for _, evt := range w.Events() {
		switch evt.(type) {
		case widget.ChangeEvent:
			pg.validateInputs(c)
			return
		}
	}
}

func (pg *verifyMessagePage) clearInputs(c *pageCommon) {
	pg.verifyBtn.Background = c.theme.Color.Hint
	pg.addressInput.Editor.SetText("")
	pg.signInput.Editor.SetText("")
	pg.messageInput.Editor.SetText("")
	pg.verifyMessage.Text = ""
}

func (pg *verifyMessagePage) handler(c pageCommon) {
	pg.handlerEditorEvents(&c, pg.addressInput.Editor)
	pg.handlerEditorEvents(&c, pg.messageInput.Editor)
	pg.handlerEditorEvents(&c, pg.signInput.Editor)

	if pg.verifyBtn.Button.Clicked() && pg.validateInputs(&c) {
		pg.verifyMessage.Text = ""
		valid, err := c.wallet.VerifyMessage(pg.addressInput.Editor.Text(), pg.messageInput.Editor.Text(), pg.signInput.Editor.Text())

		if err != nil {
			pg.verifyMessage.Color = c.theme.Color.Danger
			pg.verifyMessage.Text = err.Error()
			return
		}

		if !valid {
			pg.verifyMessage.Text = "Invalid Signature"
			pg.verifyMessage.Color = c.theme.Color.Danger
			return
		}

		pg.verifyMessage.Text = "Valid Signature"
		pg.verifyMessage.Color = c.theme.Color.Success
	}

	if pg.clearBtn.Button.Clicked() {
		pg.clearInputs(&c)
	}

	if pg.backButton.Button.Clicked() {
		pg.clearInputs(&c)
		*c.page = PageWallet
	}
}
