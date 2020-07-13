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
	page := &verifyMessagePage{
		addressInput:  c.theme.Editor(new(widget.Editor), "Address"),
		messageInput:  c.theme.Editor(new(widget.Editor), "Message"),
		signInput:     c.theme.Editor(new(widget.Editor), "Signature"),
		verifyMessage: c.theme.H6(""),
		verifyBtn:     c.theme.Button(new(widget.Clickable), "Verify"),
		clearBtn:      c.theme.Button(new(widget.Clickable), "Clear All"),
		backButton:    c.theme.PlainIconButton(new(widget.Clickable), c.icons.navigationArrowBack),
	}

	page.messageInput.IsRequired, page.addressInput.IsRequired, page.signInput.IsRequired = true, true, true
	page.messageInput.IsVisible, page.addressInput.IsVisible, page.signInput.IsVisible = true, true, true
	page.messageInput.Editor.SingleLine, page.addressInput.Editor.SingleLine, page.messageInput.Editor.SingleLine = true, true, true
	page.verifyBtn.TextSize, page.clearBtn.TextSize, page.clearBtn.TextSize = values.TextSize14, values.TextSize14, values.TextSize14
	page.verifyBtn.Background = c.theme.Color.Hint
	page.clearBtn.Background = color.RGBA{0, 0, 0, 0}
	page.clearBtn.Color = c.theme.Color.Primary
	page.backButton.Color = c.theme.Color.Hint
	page.backButton.Size = values.MarginPadding30

	return func(gtx C) D {
		page.handler(c)
		return page.Layout(c)
	}
}

func (page *verifyMessagePage) Layout(c pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding5).Layout(c.gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(c.gtx,
				layout.Rigid(page.header(&c)),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(c.gtx,
						layout.Rigid(page.inputRow(&c, page.addressInput)),
						layout.Rigid(page.inputRow(&c, page.signInput)),
						layout.Rigid(page.inputRow(&c, page.messageInput)),
						layout.Rigid(page.verifyAndClearButtons(&c)),
					)
				}),
			)
		})
	}
	return c.Layout(c.gtx, body)
}

func (page *verifyMessagePage) header(c *pageCommon) layout.Widget {
	return func(gtx C) D {
		var msg = "After you or your counterparty has genrated a signature, you can use this form to verify the signature." +
			"\nOnce you have entered the address, the message and the corresponding signature, you will see VALID if the signature" +
			"appropriately matches \nthe address and message otherwise INVALID."

		return layout.Inset{Bottom: values.MarginPadding30}.Layout(c.gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(c.gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.W.Layout(c.gtx, func(gtx C) D {
								return page.backButton.Layout(c.gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding45}.Layout(c.gtx, func(gtx C) D {
								return c.theme.H5("Verify Wallet Message").Layout(c.gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(c.gtx, func(gtx C) D {
						return layout.Dimensions{}
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := c.theme.Label(values.MarginPadding10, msg)
					return txt.Layout(c.gtx)
				}),
			)
		})
	}
}

func (page *verifyMessagePage) inputRow(c *pageCommon, editor decredmaterial.Editor) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding15}.Layout(c.gtx, func(gtx C) D {
			return editor.Layout(c.gtx)
		})
	}
}

func (page *verifyMessagePage) verifyAndClearButtons(c *pageCommon) layout.Widget {
	return func(gtx C) D {
		dims := layout.Flex{}.Layout(gtx,
			layout.Flexed(.6, func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding5, Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
					return page.verifyMessage.Layout(gtx)
				})
			}),
			layout.Flexed(.4, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(.5, func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return page.clearBtn.Layout(gtx)
						})
					}),
					layout.Flexed(.5, func(gtx C) D {
						return page.verifyBtn.Layout(gtx)
					}),
				)
			}),
		)
		return dims
	}
}

func (page *verifyMessagePage) validateInputs(c *pageCommon) bool {
	page.addressInput.ErrorLabel.Text = ""
	page.verifyBtn.Background = c.theme.Color.Hint

	if strings.Trim(page.addressInput.Editor.Text(), " ") == "" {
		page.addressInput.ErrorLabel.Text = "Please enter a valid address"
		return false
	}
	if isValid, _ := c.wallet.IsAddressValid(page.addressInput.Editor.Text()); !isValid {
		page.addressInput.ErrorLabel.Text = "Invalid address"
		return false
	}
	if strings.Trim(page.messageInput.Editor.Text(), " ") == "" {
		return false
	}
	if strings.Trim(page.signInput.Editor.Text(), " ") == "" {
		return false
	}

	page.verifyBtn.Background = c.theme.Color.Primary
	return true
}

func (page *verifyMessagePage) handlerEditorEvents(c *pageCommon, w *widget.Editor) {
	for _, evt := range w.Events() {
		switch evt.(type) {
		case widget.ChangeEvent:
			page.validateInputs(c)
			return
		}
	}
}

func (page *verifyMessagePage) clearInputs(c *pageCommon) {
	page.verifyBtn.Background = c.theme.Color.Hint
	page.addressInput.Editor.SetText("")
	page.signInput.Editor.SetText("")
	page.messageInput.Editor.SetText("")
	page.verifyMessage.Text = ""
}

func (page *verifyMessagePage) handler(c pageCommon) {
	page.handlerEditorEvents(&c, page.addressInput.Editor)
	page.handlerEditorEvents(&c, page.messageInput.Editor)
	page.handlerEditorEvents(&c, page.signInput.Editor)

	if page.verifyBtn.Button.Clicked() && page.validateInputs(&c) {
		page.verifyMessage.Text = ""
		valid, err := c.wallet.VerifyMessage(page.addressInput.Editor.Text(), page.messageInput.Editor.Text(), page.signInput.Editor.Text())

		if err != nil {
			page.verifyMessage.Color = c.theme.Color.Danger
			page.verifyMessage.Text = err.Error()
			return
		}

		if !valid {
			page.verifyMessage.Text = "Invalid Signature"
			page.verifyMessage.Color = c.theme.Color.Danger
			return
		}

		page.verifyMessage.Text = "Valid Signature"
		page.verifyMessage.Color = c.theme.Color.Success
	}

	if page.clearBtn.Button.Clicked() {
		page.clearInputs(&c)
	}

	if page.backButton.Button.Clicked() {
		page.clearInputs(&c)
		*c.page = PageWallet
	}
}
