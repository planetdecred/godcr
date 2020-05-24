package ui

import (
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageVerifyMessage = "verifymessage"

type verifyMessagePage struct {
	addressInput, messageInput, signInput    decredmaterial.Editor
	addressInputW, messageInputW, signInputW *widget.Editor
	clearBtnW, verifyBtnW                    widget.Button
	clearBtn, verifyBtn                      decredmaterial.Button
	verifyMessage                            decredmaterial.Label

	backButtonW widget.Button
	backButton  decredmaterial.IconButton
}

func (win *Window) VerifyMessagePage(c pageCommon) layout.Widget {
	page := &verifyMessagePage{
		addressInput:  c.theme.Editor("Address"),
		messageInput:  c.theme.Editor("Message"),
		signInput:     c.theme.Editor("Signature"),
		verifyMessage: c.theme.H6(""),
		verifyBtn:     c.theme.Button("Verify"),
		clearBtn:      c.theme.Button("Clear All"),
		messageInputW: &widget.Editor{
			SingleLine: true,
		},
		addressInputW: &widget.Editor{
			SingleLine: true,
		},
		signInputW: &widget.Editor{
			SingleLine: true,
		},
		backButton: c.theme.PlainIconButton(c.icons.navigationArrowBack),
	}
	page.messageInput.IsRequired, page.addressInput.IsRequired, page.signInput.IsRequired = true, true, true
	page.messageInput.IsVisible, page.addressInput.IsVisible, page.signInput.IsVisible = true, true, true
	page.verifyBtn.TextSize, page.clearBtn.TextSize, page.clearBtn.TextSize = unit.Dp(13), unit.Dp(13), unit.Dp(13)
	page.verifyBtn.Background = c.theme.Color.Hint
	page.clearBtn.Background = color.RGBA{0, 0, 0, 0}
	page.clearBtn.Color = win.theme.Color.Primary
	page.backButton.Color = c.theme.Color.Hint
	page.backButton.Size = unit.Dp(32)

	return func() {
		page.Layout(c)
		page.handler(c)
	}
}

func (page *verifyMessagePage) Layout(c pageCommon) {
	body := func() {
		layout.UniformInset(unit.Dp(5)).Layout(c.gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(c.gtx,
				layout.Rigid(page.header(&c)),
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(c.gtx,
						layout.Rigid(page.inputRow(&c, page.addressInput, page.addressInputW)),
						layout.Rigid(page.inputRow(&c, page.signInput, page.signInputW)),
						layout.Rigid(page.inputRow(&c, page.messageInput, page.messageInputW)),
						layout.Rigid(page.verifyAndClearButtons(&c)),
					)
				}),
			)
		})
	}

	c.Layout(c.gtx, body)
}

func (page *verifyMessagePage) header(c *pageCommon) layout.Widget {
	return func() {
		var msg = "After you or your counterparty has genrated a signature, you can use this form to verify the signature." +
			"\nOnce you have entered the address, the message and the corresponding signature, you will see VALID if the signature" +
			"appropriately matches \nthe address and message otherwise INVALID."

		layout.Inset{Bottom: unit.Dp(30)}.Layout(c.gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(c.gtx,
				layout.Rigid(func() {
					layout.W.Layout(c.gtx, func() {
						page.backButton.Layout(c.gtx, &page.backButtonW)
					})
					layout.Inset{Left: unit.Dp(44)}.Layout(c.gtx, func() {
						c.theme.H5("Verify Wallet Message").Layout(c.gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(15)}.Layout(c.gtx, func() {})
				}),
				layout.Rigid(func() {
					txt := c.theme.Label(unit.Dp(10), msg)
					txt.Layout(c.gtx)
				}),
			)
		})
	}
}

func (page *verifyMessagePage) inputRow(c *pageCommon, out decredmaterial.Editor, in *widget.Editor) layout.Widget {
	return func() {
		layout.Inset{Bottom: unit.Dp(15)}.Layout(c.gtx, func() {
			out.Layout(c.gtx, in)
		})
	}
}

func (page *verifyMessagePage) verifyAndClearButtons(c *pageCommon) layout.Widget {
	gtx := c.gtx
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Flexed(.6, func() {
				layout.Inset{Bottom: unit.Dp(5), Top: unit.Dp(10)}.Layout(gtx, func() {
					page.verifyMessage.Layout(gtx)
				})
			}),
			layout.Flexed(.4, func() {
				layout.Flex{}.Layout(gtx,
					layout.Flexed(.5, func() {
						layout.Inset{Left: unit.Dp(0), Right: unit.Dp(10)}.Layout(gtx, func() {
							page.clearBtn.Layout(gtx, &page.clearBtnW)
						})
					}),
					layout.Flexed(.5, func() {
						page.verifyBtn.Layout(gtx, &page.verifyBtnW)
					}),
				)
			}),
		)
	}
}

func (page *verifyMessagePage) validateInputs(c *pageCommon) bool {
	page.addressInput.ErrorLabel.Text = ""
	page.verifyBtn.Background = c.theme.Color.Hint

	if strings.Trim(page.addressInputW.Text(), " ") == "" {
		page.addressInput.ErrorLabel.Text = "Please enter a valid address"
		return false
	}
	if isValid, _ := c.wallet.IsAddressValid(page.addressInputW.Text()); !isValid {
		page.addressInput.ErrorLabel.Text = "Invalid address"
		return false
	}
	if strings.Trim(page.messageInputW.Text(), " ") == "" {
		return false
	}
	if strings.Trim(page.signInputW.Text(), " ") == "" {
		return false
	}

	page.verifyBtn.Background = c.theme.Color.Primary
	return true
}

func (page *verifyMessagePage) handlerEditorEvents(c *pageCommon, w *widget.Editor) {
	for _, evt := range w.Events(c.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			page.validateInputs(c)
			return
		}
	}
}

func (page *verifyMessagePage) clearInputs(c *pageCommon) {
	page.verifyBtn.Background = c.theme.Color.Hint
	page.addressInputW.SetText("")
	page.signInputW.SetText("")
	page.messageInputW.SetText("")
	page.verifyMessage.Text = ""
}

func (page *verifyMessagePage) handler(c pageCommon) {
	page.handlerEditorEvents(&c, page.addressInputW)
	page.handlerEditorEvents(&c, page.messageInputW)
	page.handlerEditorEvents(&c, page.signInputW)

	if page.verifyBtnW.Clicked(c.gtx) && page.validateInputs(&c) {
		page.verifyMessage.Text = ""
		valid, err := c.wallet.VerifyMessage(page.addressInputW.Text(), page.messageInputW.Text(), page.signInputW.Text())

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

	if page.clearBtnW.Clicked(c.gtx) {
		page.clearInputs(&c)
	}

	if page.backButtonW.Clicked(c.gtx) {
		page.clearInputs(&c)
		*c.page = PageWallet
	}
}
