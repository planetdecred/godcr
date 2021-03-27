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
	modal          *Modal
	titleLabel     Label
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
		titleLabel:     t.H6("Enter password to confirm"),
		passwordEditor: t.EditorPassword(editorWidget, "Password"),
		cancelButton:   cancelButton,
		confirmButton:  confirmButton,
		modal:          t.Modal(),
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

	p.handleEvents(confirm, cancel)
	p.updateColors()

	widgets := []func(gtx C) D{
		func(gtx C) D {
			return p.titleLabel.Layout(gtx)
		},
		func(gtx C) D {
			return p.theme.Separator().Layout(gtx)
		},
		func(gtx C) D {
			return p.passwordEditor.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top: unit.Dp(20),
				}.Layout(gtx, func(gtx C) D {
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
			})
		},
	}
	return p.modal.Layout(gtx, widgets, 1350)
}

func (p *Password) WithError(e string) {
	p.passwordEditor.IsRequired = true
	p.passwordEditor.SetError(e)
}

func (p *Password) updateColors() {
	p.confirmButton.Background = p.theme.Color.Hint

	if p.passwordEditor.Editor.Len() > 0 {
		p.confirmButton.Background = p.theme.Color.Primary
	}
}

func (p *Password) handleEvents(confirm func([]byte), cancel func()) {
	for p.confirmButton.Button.Clicked() {
		if p.passwordEditor.Editor.Len() > 0 {
			confirm([]byte(p.passwordEditor.Editor.Text()))
			p.reset()
		}
	}

	for p.cancelButton.Button.Clicked() {
		p.reset()
		p.passwordEditor.IsRequired = false
		p.passwordEditor.SetError("")
		cancel()
	}
}

func (p *Password) reset() {
	p.passwordEditor.Editor.SetText("")
}
