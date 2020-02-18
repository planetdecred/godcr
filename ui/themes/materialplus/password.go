package materialplus

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui"
)

type colors struct {
	cancelLabelColor            color.RGBA
	createButtonBackgroundColor color.RGBA
	nextButtonBackgroundColor   color.RGBA
}

type editor struct {
	widget   *widget.Editor
	material material.Editor
	line     *Line
}

// Password represents a form for collecting user password
type Password struct {
	colors     colors
	titleLabel material.Label

	spendingEditor *editor
	confirmEditor  *editor

	cancelButton         *widget.Button
	createButton         *widget.Button
	cancelButtonMaterial material.Button
	createButtonMaterial material.Button

	errorLabel material.Label
}

// Password initializes and returns an instance of Password
func (t *Theme) Password() *Password {
	cancelButtonMaterial := t.Button("Cancel")
	cancelButtonMaterial.Background = ui.WhiteColor
	cancelButtonMaterial.Color = ui.LightBlueColor

	errorLabel := t.Body2("")
	errorLabel.Color = ui.DangerColor

	spendingEditor := &editor{
		widget:   new(widget.Editor),
		material: t.Editor("Spending password"),
		line:     t.Line(),
	}
	confirmEditor := &editor{
		widget:   new(widget.Editor),
		material: t.Editor("Confirm spending password"),
		line:     t.Line(),
	}

	confirmEditor.line.Color = ui.LightBlueColor
	spendingEditor.line.Color = ui.LightBlueColor

	p := &Password{
		titleLabel: t.H5("Create spending password"),
		errorLabel: errorLabel,

		confirmEditor:        confirmEditor,
		spendingEditor:       spendingEditor,
		cancelButton:         new(widget.Button),
		createButton:         new(widget.Button),
		cancelButtonMaterial: cancelButtonMaterial,
		createButtonMaterial: t.Button("Create"),

		colors: colors{},
	}

	return p
}

func (p *Password) updateColors() {
	p.colors.cancelLabelColor = ui.GrayColor
	p.colors.createButtonBackgroundColor = ui.GrayColor

	if p.bothPasswordsMatch() && p.confirmEditor.widget.Len() > 0 {
		p.colors.createButtonBackgroundColor = ui.LightBlueColor
	}
}

func (p *Password) processButtonClicks(gtx *layout.Context, createFunc func(string), cancelFunc func()) {
	for p.createButton.Clicked(gtx) {
		if p.validate() {
			createFunc(p.spendingEditor.widget.Text())
		}
	}

	for p.cancelButton.Clicked(gtx) {
		p.Reset()
		cancelFunc()
	}
}

// Draw renders the widget to screen
func (p *Password) Draw(gtx *layout.Context, createFunc func(string), cancelFunc func()) {
	p.processButtonClicks(gtx, createFunc, cancelFunc)
	p.updateColors()
	p.validate()

	widgets := []func(){
		func() {
			p.titleLabel.Layout(gtx)
		},
		func() {
			inset := layout.Inset{
				Top: unit.Dp(20),
			}
			inset.Layout(gtx, func() {
				p.spendingEditor.layout(gtx)
			})
		},
		func() {
			p.confirmEditor.layout(gtx)
		},
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(2, func() {
					p.errorLabel.Layout(gtx)
				}),
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
						p.cancelButtonMaterial.Color = p.colors.cancelLabelColor
						p.cancelButtonMaterial.Layout(gtx, p.cancelButton)
					})
				}),
				layout.Rigid(func() {
					gtx.Constraints.Width.Min = 70
					layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
						p.createButtonMaterial.Background = p.colors.createButtonBackgroundColor
						p.createButtonMaterial.Layout(gtx, p.createButton)
					})
				}),
			)
		},
	}

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(widgets), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
	})
}

func (p *Password) validate() bool {
	p.errorLabel.Text = ""

	if p.spendingEditor.widget.Len() > 0 && p.confirmEditor.widget.Len() > 0 && !p.bothPasswordsMatch() {
		p.errorLabel.Text = "Both passwords do not match"
		return false
	}

	return true
}

func (p *Password) bothPasswordsMatch() bool {
	return p.confirmEditor.widget.Text() == p.spendingEditor.widget.Text()
}

// Reset empties the contents of the password form
func (p *Password) Reset() {
	p.spendingEditor.widget.SetText("")
	p.confirmEditor.widget.SetText("")
}

func (e *editor) layout(gtx *layout.Context) {
	e.material.Layout(gtx, e.widget)
	inset := layout.Inset{
		Top: unit.Dp(float32(gtx.Constraints.Height.Min)),
	}
	inset.Layout(gtx, func() {
		e.line.Draw(gtx)
	})
}
