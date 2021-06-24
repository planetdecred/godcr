package modal

import (
	"fmt"
	"github.com/planetdecred/godcr/ui"
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const Password = "password_modal"

type passwordModal struct {
	*ui.Common
	randomID string
	modal    decredmaterial.Modal
	password decredmaterial.Editor

	dialogTitle string

	isLoading      bool
	materialLoader material.LoaderStyle

	positiveButtonText    string
	positiveButtonClicked func(password string, m *passwordModal) bool // return true to dismiss dialog
	btnPositve            decredmaterial.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           decredmaterial.Button
}

func NewPasswordModal(c *ui.Common) *passwordModal {
	pm := &passwordModal{
		Common:  c,
		randomID:    fmt.Sprintf("%s-%d", Password, ui.GenerateRandomNumber()),
		modal:       *c.Theme.ModalFloatTitle(),
		btnPositve:  c.Theme.Button(new(widget.Clickable), "Confirm"),
		btnNegative: c.Theme.Button(new(widget.Clickable), "Cancel"),
	}

	pm.btnPositve.TextSize, pm.btnNegative.TextSize = values.TextSize16, values.TextSize16
	pm.btnNegative.Background, pm.btnNegative.Color = pm.Theme.Color.Surface, pm.Theme.Color.Primary
	pm.btnPositve.Font.Weight, pm.btnNegative.Font.Weight = text.Bold, text.Bold
	pm.btnPositve.Background, pm.btnPositve.Color = pm.Theme.Color.Surface, pm.Theme.Color.Primary

	pm.password = c.Theme.EditorPassword(new(widget.Editor), "Spending password")
	pm.password.Editor.SingleLine, pm.password.Editor.Submit = true, true

	th := material.NewTheme(gofont.Collection())
	pm.materialLoader = material.Loader(th)

	return pm
}

func (pm *passwordModal) ModalID() string {
	return pm.randomID
}

func (pm *passwordModal) OnResume() {
}

func (pm *passwordModal) OnDismiss() {

}

func (pm *passwordModal) Show() {
	pm.ShowModal(pm)
}

func (pm *passwordModal) Dismiss() {
	pm.DismissModal(pm)
}

func (pm *passwordModal) title(title string) *passwordModal {
	pm.dialogTitle = title
	return pm
}

func (pm *passwordModal) hint(hint string) *passwordModal {
	pm.password.Hint = hint
	return pm
}

func (pm *passwordModal) positiveButtonStyle(background, text color.NRGBA) *passwordModal {
	pm.btnPositve.Background, pm.btnPositve.Color = background, text
	return pm
}

func (pm *passwordModal) positiveButton(text string, clicked func(password string, m *passwordModal) bool) *passwordModal {
	pm.positiveButtonText = text
	pm.positiveButtonClicked = clicked
	return pm
}

func (pm *passwordModal) negativeButton(text string, clicked func()) *passwordModal {
	pm.negativeButtonText = text
	pm.negativeButtonClicked = clicked
	return pm
}

func (pm *passwordModal) setLoading(loading bool) {
	pm.isLoading = loading
}

func (pm *passwordModal) setError(err string) {
	if err == "" {
		pm.password.ClearError()
	} else {
		pm.password.SetError(err)
	}
}

func (pm *passwordModal) handle() {

	for pm.btnPositve.Button.Clicked() {

		if pm.isLoading || !ui.EditorsNotEmpty(pm.password.Editor) {
			continue
		}

		pm.setLoading(true)
		pm.setError("")
		if pm.positiveButtonClicked(pm.password.Editor.Text(), pm) {
			pm.DismissModal(pm)
		}
	}

	for pm.btnNegative.Button.Clicked() {
		if !pm.isLoading {
			pm.DismissModal(pm)
			pm.negativeButtonClicked()
		}
	}
}

func (pm *passwordModal) Layout(gtx layout.Context) ui.D {
	title := func(gtx ui.C) ui.D {
		t := pm.Theme.H6(pm.dialogTitle)
		t.Font.Weight = text.Bold
		return t.Layout(gtx)
	}

	editor := func(gtx ui.C) ui.D {
		return pm.password.Layout(gtx)
	}

	actionButtons := func(gtx ui.C) ui.D {
		return layout.E.Layout(gtx, func(gtx ui.C) ui.D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx ui.C) ui.D {
					if pm.negativeButtonText == "" {
						return layout.Dimensions{}
					}

					pm.btnNegative.Text = pm.negativeButtonText
					return pm.btnNegative.Layout(gtx)
				}),
				layout.Rigid(func(gtx ui.C) ui.D {
					if pm.isLoading {
						return pm.materialLoader.Layout(gtx)
					}

					if pm.positiveButtonText == "" {
						return layout.Dimensions{}
					}

					pm.btnPositve.Text = pm.positiveButtonText
					return pm.btnPositve.Layout(gtx)
				}),
			)
		})
	}
	var w []layout.Widget

	w = append(w, title)
	w = append(w, editor)
	w = append(w, actionButtons)

	return pm.modal.Layout(gtx, w, 850)
}
