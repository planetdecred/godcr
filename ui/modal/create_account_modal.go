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

const CreateAccount = "create_account_modal"

type CreateAccountModal struct {
	*load.Load
	randomID string
	modal    decredmaterial.Modal

	dialogTitle string

	accountName    decredmaterial.Editor
	passwordEditor decredmaterial.Editor

	isLoading      bool
	materialLoader material.LoaderStyle

	callback func(accountName, password string, m *CreateAccountModal) bool // return true to dismiss dialog

	btnPositve  decredmaterial.Button
	btnNegative decredmaterial.Button
}

func NewCreateAccountModal(l *load.Load) *CreateAccountModal {
	cm := &CreateAccountModal{
		Load:        l,
		randomID:    fmt.Sprintf("%s-%d", CreateAccount, generateRandomNumber()),
		modal:       *l.Theme.ModalFloatTitle(),
		btnPositve:  l.Theme.Button(new(widget.Clickable), "Confirm"),
		btnNegative: l.Theme.Button(new(widget.Clickable), "Cancel"),
	}

	cm.btnNegative.TextSize = values.TextSize16
	cm.btnNegative.Font.Weight = text.Bold

	cm.btnPositve.Background = cm.Theme.Color.InactiveGray
	cm.btnPositve.Font.Weight = text.Bold

	cm.accountName = l.Theme.Editor(new(widget.Editor), "Account name")
	cm.accountName.Editor.SingleLine, cm.accountName.Editor.Submit = true, true

	cm.passwordEditor = l.Theme.EditorPassword(new(widget.Editor), "Spending password")
	cm.passwordEditor.Editor.SingleLine, cm.passwordEditor.Editor.Submit = true, true

	th := material.NewTheme(gofont.Collection())
	cm.materialLoader = material.Loader(th)

	return cm
}

func (cm *CreateAccountModal) ModalID() string {
	return cm.randomID
}

func (cm *CreateAccountModal) OnResume() {
}

func (cm *CreateAccountModal) OnDismiss() {

}

func (cm *CreateAccountModal) Show() {
	cm.ShowModal(cm)
}

func (cm *CreateAccountModal) Dismiss() {
	cm.DismissModal(cm)
}

func (cm *CreateAccountModal) Title(title string) *CreateAccountModal {
	cm.dialogTitle = title
	return cm
}

func (cm *CreateAccountModal) PasswordHint(hint string) *CreateAccountModal {
	cm.passwordEditor.Hint = hint
	return cm
}

func (cm *CreateAccountModal) PasswordCreated(callback func(accountName, password string, m *CreateAccountModal) bool) *CreateAccountModal {
	cm.callback = callback
	return cm
}

func (cm *CreateAccountModal) SetLoading(loading bool) {
	cm.isLoading = loading
}

func (cm *CreateAccountModal) SetError(err string) {

}

func (cm *CreateAccountModal) Handle() {
	if editorsNotEmpty(cm.passwordEditor.Editor) || editorsNotEmpty(cm.accountName.Editor) {
		cm.btnPositve.Background = cm.Theme.Color.Primary
	} else {
		cm.btnPositve.Background = cm.Theme.Color.InactiveGray
	}

	if cm.btnPositve.Button.Clicked() || handleSubmitEvent(cm.accountName.Editor, cm.passwordEditor.Editor) {
		if editorsNotEmpty(cm.passwordEditor.Editor, cm.accountName.Editor) {
			cm.SetLoading(true)
			if cm.callback(cm.accountName.Editor.Text(), cm.passwordEditor.Editor.Text(), cm) {
				cm.Dismiss()
			}
		}

	}

	if cm.btnNegative.Button.Clicked() {
		if !cm.isLoading {
			cm.Dismiss()
		}
	}
}

func (cm *CreateAccountModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			t := cm.Theme.H6(cm.dialogTitle)
			t.Font.Weight = text.Bold
			return t.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					img := cm.Icons.ActionInfo
					img.Color = cm.Theme.Color.Gray3
					inset := layout.Inset{Right: values.MarginPadding4}
					return inset.Layout(gtx, func(gtx C) D {
						return img.Layout(gtx, values.MarginPadding20)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := cm.Theme.Label(values.MarginPadding16, "Accounts")
							txt.Color = cm.Theme.Color.Gray4
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							txt := cm.Theme.Label(values.MarginPadding16, "cannot")
							txt.Font.Weight = text.Bold
							txt.Color = cm.Theme.Color.Gray4
							inset := layout.Inset{Right: values.MarginPadding2, Left: values.MarginPadding2}
							return inset.Layout(gtx, txt.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							txt := cm.Theme.Label(values.MarginPadding16, "be deleted once created")
							txt.Color = cm.Theme.Color.Gray4
							return txt.Layout(gtx)
						}),
					)
				}),
			)
		},
		func(gtx C) D {
			return cm.accountName.Layout(gtx)
		},
		func(gtx C) D {
			return cm.passwordEditor.Layout(gtx)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {

						cm.btnNegative.Background = cm.Theme.Color.Surface
						cm.btnNegative.Color = cm.Theme.Color.Primary
						return cm.btnNegative.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if cm.isLoading {
							return cm.materialLoader.Layout(gtx)
						}

						return cm.btnPositve.Layout(gtx)
					}),
				)
			})
		},
	}

	return cm.modal.Layout(gtx, w, 850)
}
