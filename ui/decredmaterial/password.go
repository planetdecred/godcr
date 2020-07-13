package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Password struct {
	theme          *Theme
	passwordEditor Editor
	confirmButton  Button
	cancelButton   Button
}

// Password initializes and returns an instance of Password
func (t *Theme) Password() *Password {
	cancelButton := t.Button(new(widget.Clickable), "Cancel")
	cancelButton.Background = t.Color.Surface
	cancelButton.Color = t.Color.Primary
	confirmButton := t.Button(new(widget.Clickable), "Confirm")

	editorWidget := &widget.Editor{
		SingleLine: true,
		Mask:       '*',
	}
	p := &Password{
		theme:          t,
		passwordEditor: t.Editor(editorWidget, "Password"),
		cancelButton:   cancelButton,
		confirmButton:  confirmButton,
	}

	return p
}

// Layout renders the widget to screen. The confirm function passed by the calling page is called when the confirm button
// is clicked, and the form passes validation. The entered password is passed as an argument to the confirm func.
// The cancel func is called when the cancel button is clicked
func (p *Password) Layout(gtx layout.Context, confirm func([]byte), cancel func()) layout.Dimensions {
	if !p.passwordEditor.Editor.Focused() {
		p.passwordEditor.Editor.Focus()
	}

	p.handleEvents(gtx, confirm, cancel)
	p.updateColors()

	widgets := []func(gtx C) D{
		func(gtx C) D {
			return p.passwordEditor.Layout(gtx)
		},
		func(gtx C) D {
			inset := layout.Inset{
				Top: unit.Dp(20),
			}
			return inset.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return p.confirmButton.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						inset := layout.Inset{
							Left: unit.Dp(10),
						}
						return inset.Layout(gtx, func(gtx C) D {
							return p.cancelButton.Layout(gtx)
						})
					}),
				)
			})
		},
	}
	return p.theme.Modal(gtx, "Enter password to confirm", widgets)
}

func (p *Password) updateColors() {
	p.confirmButton.Background = p.theme.Color.Hint

	if p.passwordEditor.Editor.Len() > 0 {
		p.confirmButton.Background = p.theme.Color.Primary
	}
}

func (p *Password) handleEvents(gtx layout.Context, confirm func([]byte), cancel func()) {
	for p.confirmButton.Button.Clicked() {
		if p.passwordEditor.Editor.Len() > 0 {
			confirm([]byte(p.passwordEditor.Editor.Text()))
			p.reset()
		}
	}

	for p.cancelButton.Button.Clicked() {
		p.reset()
		cancel()
	}
}

func (p *Password) reset() {
	p.passwordEditor.Editor.SetText("")
}
