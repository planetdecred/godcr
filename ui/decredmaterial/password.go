package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Password struct {
	theme *Theme

	titleLabel Label

	passwordEditorMaterial Editor
	passwordEditorWidget   *widget.Editor

	confirmButtonMaterial Button
	confirmButtonWidget   *widget.Button

	cancelButtonMaterial Button
	cancelButtonWidget   *widget.Button
}

// Password initializes and returns an instance of Password
func (t *Theme) Password() *Password {
	cancelButtonMaterial := t.Button("Cancel")
	cancelButtonMaterial.Background = t.Color.Surface
	cancelButtonMaterial.Color = t.Color.Primary

	p := &Password{
		theme:      t,
		titleLabel: t.H5("Enter password to confirm"),

		passwordEditorMaterial: t.Editor("Password"),
		passwordEditorWidget:   new(widget.Editor),

		cancelButtonMaterial:  cancelButtonMaterial,
		confirmButtonMaterial: t.Button("Confirm"),

		cancelButtonWidget:  new(widget.Button),
		confirmButtonWidget: new(widget.Button),
	}

	return p
}

// Layout renders the widget to screen. The confirm function passed by the calling page is called when the confirm button
// is clicked, and the form passes validation. The entered password is passed as an argument to the confirm func.
// The cancel func is called when the cancel button is clicked
func (p *Password) Layout(gtx *layout.Context, confirm func([]byte), cancel func()) {
	p.handleEvents(gtx, confirm, cancel)
	p.updateColors()

	widgets := []func(){
		func() {
			p.titleLabel.Layout(gtx)
		},
		func() {
			p.passwordEditorMaterial.Layout(gtx, p.passwordEditorWidget)
		},
		func() {
			inset := layout.Inset{
				Top: unit.Dp(20),
			}
			inset.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						p.confirmButtonMaterial.Layout(gtx, p.confirmButtonWidget)
					}),
					layout.Rigid(func() {
						inset := layout.Inset{
							Left: unit.Dp(10),
						}
						inset.Layout(gtx, func() {
							p.cancelButtonMaterial.Layout(gtx, p.cancelButtonWidget)
						})
					}),
				)
			})
		},
	}

	p.theme.Modal(gtx, func() {
		list := layout.List{Axis: layout.Vertical}
		list.Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
		})
	})
}

func (p *Password) updateColors() {
	p.confirmButtonMaterial.Background = p.theme.Color.Hint

	if p.passwordEditorWidget.Len() > 0 {
		p.confirmButtonMaterial.Background = p.theme.Color.Primary
	}
}

func (p *Password) handleEvents(gtx *layout.Context, confirm func([]byte), cancel func()) {
	for p.confirmButtonWidget.Clicked(gtx) {
		if p.passwordEditorWidget.Len() > 0 {
			confirm([]byte(p.passwordEditorWidget.Text()))
			p.reset()
		}
	}

	for p.cancelButtonWidget.Clicked(gtx) {
		p.reset()
		cancel()
	}
}

func (p *Password) reset() {
	p.passwordEditorWidget.SetText("")
}
