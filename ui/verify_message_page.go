package ui

import (
	"image/color"
	"strings"

	"github.com/planetdecred/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageVerifyMessage = "VerifyMessage"

type verifyMessagePage struct {
	theme                                 *decredmaterial.Theme
	addressInput, messageInput, signInput decredmaterial.Editor
	clearBtn, verifyBtn                   decredmaterial.Button
	verifyMessage                         decredmaterial.Label

	line                *decredmaterial.Line
	verifyMessageStatus *widget.Icon
}

func (win *Window) VerifyMessagePage(c pageCommon) layout.Widget {
	pg := &verifyMessagePage{
		theme:         c.theme,
		addressInput:  c.theme.Editor(new(widget.Editor), "Address"),
		messageInput:  c.theme.Editor(new(widget.Editor), "Message"),
		signInput:     c.theme.Editor(new(widget.Editor), "Signature"),
		verifyMessage: c.theme.Body1(""),
		verifyBtn:     c.theme.Button(new(widget.Clickable), "Verify message"),
		clearBtn:      c.theme.Button(new(widget.Clickable), "Clear all"),
		line:          c.theme.Line(),
	}

	pg.messageInput.IsVisible, pg.addressInput.IsVisible, pg.signInput.IsVisible = true, true, true
	pg.messageInput.Editor.SingleLine, pg.addressInput.Editor.SingleLine, pg.messageInput.Editor.SingleLine = true, true, true
	pg.verifyBtn.TextSize, pg.clearBtn.TextSize, pg.clearBtn.TextSize = values.TextSize14, values.TextSize14, values.TextSize14
	pg.verifyBtn.Background = c.theme.Color.Hint
	pg.clearBtn.Background = color.RGBA{0, 0, 0, 0}
	pg.clearBtn.Color = c.theme.Color.Primary

	return func(gtx C) D {
		pg.handle(c)
		return pg.Layout(gtx, c)
	}
}

func (pg *verifyMessagePage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		load := SubPage{
			"Verify message",
			c.info.Wallets[*c.selectedWallet].Name,
			func() {
				pg.clearInputs(&c)
				*c.page = PageWallet
			},
			func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(pg.description()),
						layout.Rigid(pg.inputRow(pg.addressInput)),
						layout.Rigid(pg.inputRow(pg.signInput)),
						layout.Rigid(pg.inputRow(pg.messageInput)),
						layout.Rigid(pg.verifyAndClearButtons()),
						layout.Rigid(pg.verifyMessageResponse()),
					)
				})
			},
		}
		return c.SubPageLayout(gtx, load)
	}
	return c.Layout(gtx, body)
}

func (pg *verifyMessagePage) inputRow(editor decredmaterial.Editor) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
			return editor.Layout(gtx)
		})
	}
}

func (pg *verifyMessagePage) description() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		desc := pg.theme.Caption("Enter the address, signature, and message to verify:")
		desc.Color = pg.theme.Color.Gray
		return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
			return desc.Layout(gtx)
		})
	}
}

func (pg *verifyMessagePage) verifyAndClearButtons() layout.Widget {
	return func(gtx C) D {
		dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return pg.clearBtn.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.verifyBtn.Layout(gtx)
						}),
					)
				})
			}),
		)
		return dims
	}
}

func (pg *verifyMessagePage) verifyMessageResponse() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if pg.verifyMessageStatus != nil {
			return layout.Inset{Top: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
				pg.line.Width = gtx.Constraints.Max.X
				pg.line.Color = pg.theme.Color.Hint
				pg.line.Layout(gtx)

				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return pg.verifyMessageStatus.Layout(gtx, values.MarginPadding20)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return pg.verifyMessage.Layout(gtx)
						}),
					)
				})
			})
		}
		return layout.Dimensions{}
	}
}

func (pg *verifyMessagePage) handle(c pageCommon) {
	pg.verifyBtn.Background = c.theme.Color.Hint
	if pg.inputsNotEmpty(c) {
		pg.verifyBtn.Background = c.theme.Color.Primary
		if pg.verifyBtn.Button.Clicked() {
			pg.verifyMessage.Text = ""
			pg.verifyMessageStatus = nil
			valid, err := c.wallet.VerifyMessage(pg.addressInput.Editor.Text(), pg.messageInput.Editor.Text(), pg.signInput.Editor.Text())
			if err != nil {
				pg.signInput.SetError("Invalid signature")
				return
			}

			if !valid {
				pg.verifyMessageStatus = c.icons.navigationCancel
				pg.verifyMessage.Text = "Invalid signature"
				pg.verifyMessage.Color = c.theme.Color.Danger
				return
			}

			pg.verifyMessageStatus = c.icons.actionCheck
			pg.verifyMessageStatus.Color = c.theme.Color.Success
			pg.verifyMessage.Text = "Valid signature"
			pg.verifyMessage.Color = c.theme.Color.Success
		}
	}

	pg.handlerEditorEvents(&c, pg.addressInput.Editor)
	if pg.clearBtn.Button.Clicked() {
		pg.clearInputs(&c)
	}
}

func (pg *verifyMessagePage) handlerEditorEvents(c *pageCommon, w *widget.Editor) {
	for _, evt := range w.Events() {
		switch evt.(type) {
		case widget.ChangeEvent:
			pg.validateAddress(*c)
			return
		}
	}
}

func (pg *verifyMessagePage) clearInputs(c *pageCommon) {
	pg.verifyMessageStatus = nil
	pg.verifyBtn.Background = c.theme.Color.Hint
	pg.addressInput.Editor.SetText("")
	pg.signInput.Editor.SetText("")
	pg.messageInput.Editor.SetText("")
	pg.verifyMessage.Text = ""
	pg.addressInput.SetError("")
	pg.signInput.SetError("")
}

func (pg *verifyMessagePage) validateAddress(c pageCommon) bool {
	if isValid, _ := c.wallet.IsAddressValid(pg.addressInput.Editor.Text()); !isValid {
		pg.addressInput.SetError("Invalid address")
		return false
	}

	return true
}

func (pg *verifyMessagePage) inputsNotEmpty(c pageCommon) bool {
	if strings.Trim(pg.addressInput.Editor.Text(), " ") == "" {
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
