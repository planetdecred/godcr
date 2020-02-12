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
}

type pinTabWidgets struct {
	spendingEditor         *widget.Editor
	spendingEditorMaterial material.Editor

	confirmEditor         *widget.Editor
	confirmEditorMaterial material.Editor

	isShowingConfirmEditor bool
}

type colors struct {
	cancelLabelColor            color.RGBA
	createButtonBackgroundColor color.RGBA
	nextButtonBackgroundColor   color.RGBA
}

// PasswordAndPin represents the spending password and pin widget
type PasswordAndPin struct {
	colors             colors
	tabContainer       *TabContainer
	passwordTabWidgets *passwordTabWidgets
	pinTabWidgets      *pinTabWidgets

	cancelButton         *widget.Button
	cancelButtonMaterial material.Button

	createButton         *widget.Button
	createButtonMaterial material.Button

	nextButton         *widget.Button
	nextButtonMaterial material.Button

	errorLabel material.Label
}

const (
	passwordTabID int32 = iota
	pinTabID
)

// PasswordAndPin returns an instance of the PasswordAndPin widget
func (t *Theme) PasswordAndPin() *PasswordAndPin {
	cancelButtonMaterial := t.Button("Cancel")
	cancelButtonMaterial.Background = ui.WhiteColor
	cancelButtonMaterial.Color = ui.LightBlueColor

	errorLabel := t.Body2("")
	errorLabel.Color = ui.DangerColor

	passwordTabWidgets := &passwordTabWidgets{
		spendingEditor:         new(widget.Editor),
		confirmEditor:          new(widget.Editor),
		spendingEditorMaterial: t.Editor("Spending password"),
		confirmEditorMaterial:  t.Editor("Confirm spending password"),
	}

	pinTabWidgets := &pinTabWidgets{
		spendingEditor:         new(widget.Editor),
		confirmEditor:          new(widget.Editor),
		spendingEditorMaterial: t.Editor("Enter spending PIN"),
		confirmEditorMaterial:  t.Editor("Enter spending PIN again"),
	}

	p := &PasswordAndPin{
		cancelButton:         new(widget.Button),
		cancelButtonMaterial: cancelButtonMaterial,

		createButton:         new(widget.Button),
		nextButton:           new(widget.Button),
		nextButtonMaterial:   t.Button("Next"),
		createButtonMaterial: t.Button("Create"),

		passwordTabWidgets: passwordTabWidgets,
		pinTabWidgets:      pinTabWidgets,
		colors:             colors{},

		errorLabel: errorLabel,
	}

	tabItems := []Tab{
		{
			ID:         passwordTabID,
			Label:      "Password",
			RenderFunc: p.passwordTab,
		},
		{
			ID:         pinTabID,
			Label:      "Pin",
			RenderFunc: p.pinTab,
		},
	}
	p.tabContainer = t.TabContainer(tabItems)

	return p
}

func (p *PasswordAndPin) pinTab(gtx *layout.Context) {
	w := []func(){
		func() {
			gtx.Constraints.Height.Min = 50
			if !p.pinTabWidgets.isShowingConfirmEditor {
				p.pinTabWidgets.spendingEditorMaterial.Layout(gtx, p.pinTabWidgets.spendingEditor)
			} else {
				p.pinTabWidgets.confirmEditorMaterial.Layout(gtx, p.pinTabWidgets.confirmEditor)
			}
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
			w[i]()
		})
	})
}

func (p *PasswordAndPin) passwordTab(gtx *layout.Context) {
	w := []func(){
		func() {
			gtx.Constraints.Height.Min = 50
			p.passwordTabWidgets.spendingEditorMaterial.Layout(gtx, p.passwordTabWidgets.spendingEditor)
		},
		func() {
			gtx.Constraints.Height.Min = 50
			p.passwordTabWidgets.confirmEditorMaterial.Layout(gtx, p.passwordTabWidgets.confirmEditor)
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
			w[i]()
		})
	})
}

func (p *PasswordAndPin) updateColors() {
	p.colors.cancelLabelColor = ui.GrayColor
	p.colors.createButtonBackgroundColor = ui.GrayColor
	p.colors.nextButtonBackgroundColor = ui.GrayColor

	currentTabID := p.tabContainer.GetCurrentTabID()
	p.colors.cancelLabelColor = ui.LightBlueColor

	if currentTabID == pinTabID {
		if p.pinTabWidgets.spendingEditor.Len() > 0 {
			p.colors.nextButtonBackgroundColor = ui.LightBlueColor
		}

		if p.pinTabWidgets.isShowingConfirmEditor && p.bothPinsMatch() {
			p.colors.createButtonBackgroundColor = ui.LightBlueColor
		}
	} else {
		if p.bothPasswordsMatch() && p.passwordTabWidgets.confirmEditor.Len() > 0 {
			p.colors.createButtonBackgroundColor = ui.LightBlueColor
		}
	}
}

// Draw renders this widget
func (p *PasswordAndPin) Draw(gtx *layout.Context, createFunc func(string, int32), cancelFunc func()) {
	go p.updateColors()
	go p.validate(p.tabContainer.GetCurrentTabID())

	currentTabID := p.tabContainer.GetCurrentTabID()

	for p.nextButton.Clicked(gtx) {
		if p.pinTabWidgets.spendingEditor.Len() > 0 {
			p.pinTabWidgets.isShowingConfirmEditor = true
		}
	}

	for p.createButton.Clicked(gtx) {
		pwd := p.passwordTabWidgets.spendingEditor.Text()

		if currentTabID == pinTabID {
			pwd = p.pinTabWidgets.spendingEditor.Text()
		}

		if p.validate(currentTabID) {
			createFunc(pwd, currentTabID)
		}
	}

	for p.cancelButton.Clicked(gtx) {
		p.reset()
		cancelFunc()
	}

	w := []func(){
		func() {
			p.tabContainer.Draw(gtx, "Create spending password")
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
						if currentTabID == passwordTabID || p.pinTabWidgets.isShowingConfirmEditor {
							p.createButtonMaterial.Background = p.colors.createButtonBackgroundColor
							p.createButtonMaterial.Layout(gtx, p.createButton)
						} else {
							p.nextButtonMaterial.Background = p.colors.nextButtonBackgroundColor
							p.nextButtonMaterial.Layout(gtx, p.nextButton)
						}
					})
				}),
			)
		},
	}

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(w), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, w[i])
	})

}

func (p *PasswordAndPin) validate(currentTabID int32) bool {
	p.errorLabel.Text = ""

	if currentTabID == passwordTabID {
		if p.passwordTabWidgets.spendingEditor.Text() != "" && p.passwordTabWidgets.confirmEditor.Text() != "" {
			if !p.bothPasswordsMatch() {
				p.errorLabel.Text = "Both passwords do not match"
				return false
			}
		}
	} else {
		if p.pinTabWidgets.spendingEditor.Len() > 0 && p.pinTabWidgets.confirmEditor.Len() > 0 {
			if !p.bothPinsMatch() {
				p.errorLabel.Text = "Both pins do not match"
				return false
			}
		}
	}

	return true
}

func (p *PasswordAndPin) bothPasswordsMatch() bool {
	return p.passwordTabWidgets.confirmEditor.Text() == p.passwordTabWidgets.spendingEditor.Text()
}

func (p *PasswordAndPin) bothPinsMatch() bool {
	return p.pinTabWidgets.confirmEditor.Text() == p.pinTabWidgets.spendingEditor.Text()
}

func (p *PasswordAndPin) reset() {
	p.passwordTabWidgets.spendingEditor.SetText("")
	p.passwordTabWidgets.confirmEditor.SetText("")
	p.pinTabWidgets.spendingEditor.SetText("")
	p.pinTabWidgets.confirmEditor.SetText("")
}
