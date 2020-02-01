package materialplus

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui"
)

type passwordTabWidgets struct {
	spendingEditor         *widget.Editor
	spendingEditorMaterial material.Editor

	confirmEditor         *widget.Editor
	confirmEditorMaterial material.Editor

	errorLabel material.Label
}

type pinTabWidgets struct {
	spendingEditor         *widget.Editor
	spendingEditorMaterial material.Editor

	confirmEditor         *widget.Editor
	confirmEditorMaterial material.Editor

	errorLabel material.Label
}

type colors struct {
	cancelLabelColor            color.RGBA
	createButtonBackgroundColor color.RGBA
}

// PasswordAndPin represents the spending password and pin widget
type PasswordAndPin struct {
	tabContainer       *TabContainer
	passwordTabWidgets *passwordTabWidgets

	cancelButton         *widget.Button
	cancelButtonMaterial material.Button

	createButton         *widget.Button
	createButtonMaterial material.Button

	colors colors

	currentTab string

	isCreating bool
}

const (
	passwordTabLabel = "Password"
	pinTabLabel      = "Pin"
)

// PasswordAndPin returns an instance of the PasswordAndPin widget
func (t *Theme) PasswordAndPin() *PasswordAndPin {
	cancelButtonMaterial := t.Button("Cancel")
	cancelButtonMaterial.Background = ui.WhiteColor
	cancelButtonMaterial.Color = ui.LightBlueColor

	errorLabel := t.Body2("")
	errorLabel.Color = ui.DangerColor

	p := &PasswordAndPin{
		cancelButton:         new(widget.Button),
		cancelButtonMaterial: cancelButtonMaterial,

		createButton:         new(widget.Button),
		createButtonMaterial: t.Button("Create"),

		passwordTabWidgets: &passwordTabWidgets{
			spendingEditor:         new(widget.Editor),
			confirmEditor:          new(widget.Editor),
			spendingEditorMaterial: t.Editor("Spending password"),
			confirmEditorMaterial:  t.Editor("Confirm spending password"),

			errorLabel: errorLabel,
		},
		colors: colors{},
	}

	tabItems := []Tab{
		Tab{
			Label:      passwordTabLabel,
			RenderFunc: p.passwordTab,
		},
		Tab{
			Label:      pinTabLabel,
			RenderFunc: p.pinTab,
		},
	}

	p.tabContainer = t.TabContainer(tabItems)
	p.currentTab = tabItems[0].Label

	return p
}

func (p *PasswordAndPin) passwordTab(gtx *layout.Context) {
	go func() {
		p.updateColors()
		p.validatePasswordTab()
	}()

	w := []func(){
		func() {
			p.passwordTabWidgets.spendingEditorMaterial.Layout(gtx, p.passwordTabWidgets.spendingEditor)
		},
		func() {
			p.passwordTabWidgets.confirmEditorMaterial.Layout(gtx, p.passwordTabWidgets.confirmEditor)
		},
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(2, func() {
					p.passwordTabWidgets.errorLabel.Layout(gtx)
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
	list.Layout(gtx, len(w), func(i int) {
		inset := layout.Inset{
			Top:   unit.Dp(7),
			Left:  unit.Dp(25),
			Right: unit.Dp(25),
		}
		inset.Layout(gtx, func() {
			gtx.Constraints.Height.Min = 60
			w[i]()
		})
	})
}

func (p *PasswordAndPin) updateColors() {
	if p.isCreating {
		p.colors.cancelLabelColor = ui.GrayColor
	} else {
		p.colors.cancelLabelColor = ui.LightBlueColor
	}

	// create button
	if p.isCreating {
		p.colors.createButtonBackgroundColor = ui.GrayColor
	} else {
		if p.bothPasswordsMatch() && p.passwordTabWidgets.confirmEditor.Len() > 0 {
			p.colors.createButtonBackgroundColor = ui.LightBlueColor
		} else {
			p.colors.createButtonBackgroundColor = ui.GrayColor
		}
	}
}

func (p *PasswordAndPin) reset() {
	p.passwordTabWidgets.spendingEditor.SetText("")
	p.passwordTabWidgets.confirmEditor.SetText("")
}

func (p *PasswordAndPin) pinTab(gtx *layout.Context) {

}

// Draw renders this widget
func (p *PasswordAndPin) Draw(gtx *layout.Context, createFunc func(string), cancelFunc func()) {
	for p.cancelButton.Clicked(gtx) {
		if !p.isCreating {
			p.reset()
			cancelFunc()
		}
	}

	for p.createButton.Clicked(gtx) {
		if p.tabContainer.GetCurrentTabLabel() == passwordTabLabel {
			p.validatePasswordTabAndSubmit(createFunc)
		} else {
			p.validatePinTabAndSubmit(createFunc)
		}
	}
	p.tabContainer.Draw(gtx)
}

func (p *PasswordAndPin) validatePasswordTab() bool {
	if p.passwordTabWidgets.spendingEditor.Text() == "" {
		p.passwordTabWidgets.errorLabel.Text = ""
		return false
	}

	if !p.bothPasswordsMatch() {
		p.passwordTabWidgets.errorLabel.Text = "Both passwords do not match"
		return false
	}

	p.passwordTabWidgets.errorLabel.Text = ""

	return true
}

func (p *PasswordAndPin) bothPasswordsMatch() bool {
	if p.passwordTabWidgets.confirmEditor.Text() == p.passwordTabWidgets.spendingEditor.Text() {
		return true
	}

	return false
}

func (p *PasswordAndPin) validatePinTab() bool {
	return true
}

func (p *PasswordAndPin) validatePasswordTabAndSubmit(createFunc func(string)) bool {
	if p.validatePasswordTab() {
		createFunc(p.passwordTabWidgets.spendingEditor.Text())
		return true
	}
	return false
}

func (p *PasswordAndPin) validatePinTabAndSubmit(createFunc func(string)) bool {
	if p.validatePinTab() {
		//createFunc(p.pinTa)
	}
	return false
}
