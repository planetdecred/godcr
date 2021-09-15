package modal

import (
	"fmt"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const Password = "password_modal"

type PasswordModal struct {
	*load.Load
	randomID string
	modal    decredmaterial.Modal
	password decredmaterial.Editor

	dialogTitle string

	isLoading    bool
	isCancelable bool

	materialLoader material.LoaderStyle

	positiveButtonText    string
	positiveButtonClicked func(password string, m *PasswordModal) bool // return true to dismiss dialog
	btnPositve            decredmaterial.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           decredmaterial.Button
}

func NewPasswordModal(l *load.Load) *PasswordModal {
	pm := &PasswordModal{
		Load:         l,
		randomID:     fmt.Sprintf("%s-%d", Password, generateRandomNumber()),
		modal:        *l.Theme.ModalFloatTitle(),
		btnPositve:   l.Theme.Button(new(widget.Clickable), "Confirm"),
		btnNegative:  l.Theme.OutlineButton(new(widget.Clickable), "Cancel"),
		isCancelable: true,
	}

	pm.btnPositve.Font.Weight = text.Medium

	pm.btnNegative.Font.Weight = text.Medium
	pm.btnNegative.Margin.Right = values.MarginPadding8

	pm.password = l.Theme.EditorPassword(new(widget.Editor), "Spending password")
	pm.password.Editor.SingleLine, pm.password.Editor.Submit = true, true

	th := material.NewTheme(gofont.Collection())
	pm.materialLoader = material.Loader(th)

	return pm
}

func (pm *PasswordModal) ModalID() string {
	return pm.randomID
}

func (pm *PasswordModal) OnResume() {
}

func (pm *PasswordModal) OnDismiss() {

}

func (pm *PasswordModal) Show() {
	pm.ShowModal(pm)
}

func (pm *PasswordModal) Dismiss() {
	pm.DismissModal(pm)
}

func (pm *PasswordModal) Title(title string) *PasswordModal {
	pm.dialogTitle = title
	return pm
}

func (pm *PasswordModal) Hint(hint string) *PasswordModal {
	pm.password.Hint = hint
	return pm
}

func (pm *PasswordModal) PositiveButton(text string, clicked func(password string, m *PasswordModal) bool) *PasswordModal {
	pm.positiveButtonText = text
	pm.positiveButtonClicked = clicked
	return pm
}

func (pm *PasswordModal) NegativeButton(text string, clicked func()) *PasswordModal {
	pm.negativeButtonText = text
	pm.negativeButtonClicked = clicked
	return pm
}

func (pm *PasswordModal) SetLoading(loading bool) {
	pm.isLoading = loading
}

func (pm *PasswordModal) SetCancelable(min bool) *PasswordModal {
	pm.isCancelable = min
	return pm
}

func (pm *PasswordModal) SetError(err string) {
	if err == "" {
		pm.password.ClearError()
	} else {
		pm.password.SetError(err)
	}
}

func (pm *PasswordModal) Handle() {
	pm.btnPositve.SetEnabled(editorsNotEmpty(pm.password.Editor))
	for pm.btnPositve.Clicked() || handleSubmitEvent(pm.password.Editor) {
		if pm.isLoading || !editorsNotEmpty(pm.password.Editor) {
			continue
		}

		pm.SetLoading(true)
		pm.SetError("")
		if pm.positiveButtonClicked(pm.password.Editor.Text(), pm) {
			pm.DismissModal(pm)
		}
	}

	pm.btnNegative.SetEnabled(!pm.isLoading)
	for pm.btnNegative.Clicked() {
		if !pm.isLoading {
			pm.DismissModal(pm)
			pm.negativeButtonClicked()
		}
	}

	if pm.modal.BackdropClicked(pm.isCancelable) {
		if !pm.isLoading {
			pm.Dismiss()
		}
	}
}

func (pm *PasswordModal) Layout(gtx layout.Context) D {
	title := func(gtx C) D {
		t := pm.Theme.H6(pm.dialogTitle)
		t.Font.Weight = text.Bold
		return t.Layout(gtx)
	}

	editor := func(gtx C) D {
		return pm.password.Layout(gtx)
	}

	actionButtons := func(gtx C) D {
		return layout.E.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if pm.negativeButtonText == "" || pm.isLoading {
						return layout.Dimensions{}
					}

					pm.btnNegative.Text = pm.negativeButtonText
					return pm.btnNegative.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
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
