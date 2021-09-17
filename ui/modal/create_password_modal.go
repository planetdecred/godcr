package modal

import (
	"fmt"
	"strconv"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWallet = "create_wallet_modal"

type CreatePasswordModal struct {
	*load.Load

	modal                 decredmaterial.Modal
	walletName            decredmaterial.Editor
	passwordEditor        decredmaterial.Editor
	confirmPasswordEditor decredmaterial.Editor
	passwordStrength      decredmaterial.ProgressBarStyle

	isLoading          bool
	isCancelable       bool
	walletNameEnabled  bool
	showWalletWarnInfo bool

	dialogTitle string
	randomID    string

	materialLoader material.LoaderStyle

	btnPositve  decredmaterial.Button
	btnNegative decredmaterial.Button

	callback func(walletName, password string, m *CreatePasswordModal) bool // return true to dismiss dialog
}

func NewCreatePasswordModal(l *load.Load) *CreatePasswordModal {
	cm := &CreatePasswordModal{
		Load:             l,
		randomID:         fmt.Sprintf("%s-%d", CreateWallet, generateRandomNumber()),
		modal:            *l.Theme.ModalFloatTitle(),
		passwordStrength: l.Theme.ProgressBar(0),
		btnPositve:       l.Theme.Button(new(widget.Clickable), "Confirm"),
		btnNegative:      l.Theme.Button(new(widget.Clickable), "Cancel"),
		isCancelable:     true,
	}

	cm.btnNegative.TextSize = values.TextSize16
	cm.btnNegative.Font.Weight = text.Medium

	cm.btnPositve.Background = cm.Theme.Color.InactiveGray
	cm.btnPositve.Font.Weight = text.Bold

	cm.walletName = l.Theme.Editor(new(widget.Editor), "Wallet name")
	cm.walletName.Editor.SingleLine, cm.walletName.Editor.Submit = true, true

	cm.passwordEditor = l.Theme.EditorPassword(new(widget.Editor), "Spending password")
	cm.passwordEditor.Editor.SingleLine, cm.passwordEditor.Editor.Submit = true, true

	cm.confirmPasswordEditor = l.Theme.EditorPassword(new(widget.Editor), "Spending password")
	cm.confirmPasswordEditor.Editor.SingleLine, cm.confirmPasswordEditor.Editor.Submit = true, true

	th := material.NewTheme(gofont.Collection())
	cm.materialLoader = material.Loader(th)

	return cm
}

func (cm *CreatePasswordModal) ModalID() string {
	return cm.randomID
}

func (cm *CreatePasswordModal) OnResume() {
}

func (cm *CreatePasswordModal) OnDismiss() {

}

func (cm *CreatePasswordModal) Show() {
	cm.ShowModal(cm)
}

func (cm *CreatePasswordModal) Dismiss() {
	cm.DismissModal(cm)
}

func (cm *CreatePasswordModal) Title(title string) *CreatePasswordModal {
	cm.dialogTitle = title
	return cm
}

func (cm *CreatePasswordModal) EnableName(enable bool) *CreatePasswordModal {
	cm.walletNameEnabled = enable
	return cm
}

func (cm *CreatePasswordModal) PasswordHint(hint string) *CreatePasswordModal {
	cm.passwordEditor.Hint = hint
	return cm
}

func (cm *CreatePasswordModal) ConfirmPasswordHint(hint string) *CreatePasswordModal {
	cm.confirmPasswordEditor.Hint = hint
	return cm
}

func (cm *CreatePasswordModal) ShowWalletInfoTip(show bool) *CreatePasswordModal {
	cm.showWalletWarnInfo = show
	return cm
}

func (cm *CreatePasswordModal) PasswordCreated(callback func(walletName, password string, m *CreatePasswordModal) bool) *CreatePasswordModal {
	cm.callback = callback
	return cm
}

func (cm *CreatePasswordModal) SetLoading(loading bool) {
	cm.isLoading = loading
}

func (cm *CreatePasswordModal) SetCancelable(min bool) *CreatePasswordModal {
	cm.isCancelable = min
	return cm
}

func (cm *CreatePasswordModal) SetError(err string) {

}

func (cm *CreatePasswordModal) Handle() {
	if editorsNotEmpty(cm.passwordEditor.Editor) || editorsNotEmpty(cm.walletName.Editor) ||
		editorsNotEmpty(cm.confirmPasswordEditor.Editor) {
		cm.btnPositve.Background = cm.Theme.Color.Primary
	} else {
		cm.btnPositve.Background = cm.Theme.Color.InactiveGray
	}

	if cm.passwordEditor.Editor.Text() == cm.confirmPasswordEditor.Editor.Text() {
		// reset error label when password and matching password fields match
		cm.confirmPasswordEditor.SetError("")
	}

	if cm.btnPositve.Button.Clicked() || handleSubmitEvent(cm.walletName.Editor, cm.passwordEditor.Editor, cm.confirmPasswordEditor.Editor) {

		nameValid := true
		if cm.walletNameEnabled {
			nameValid = editorsNotEmpty(cm.walletName.Editor)
		}

		if nameValid && editorsNotEmpty(cm.passwordEditor.Editor, cm.confirmPasswordEditor.Editor) &&
			cm.passwordsMatch(cm.passwordEditor.Editor, cm.confirmPasswordEditor.Editor) {

			cm.SetLoading(true)
			if cm.callback(cm.walletName.Editor.Text(), cm.passwordEditor.Editor.Text(), cm) {
				cm.Dismiss()
			}
		}

	}

	if cm.btnNegative.Button.Clicked() {
		if !cm.isLoading {
			cm.Dismiss()
		}
	}

	if cm.modal.BackdropClicked(cm.isCancelable) {
		if !cm.isLoading {
			cm.Dismiss()
		}
	}

	computePasswordStrength(&cm.passwordStrength, cm.Theme, cm.passwordEditor.Editor)

}
func (cm *CreatePasswordModal) passwordsMatch(editors ...*widget.Editor) bool {
	if len(editors) < 2 {
		return false
	}

	password := editors[0]
	matching := editors[1]

	if password.Text() != matching.Text() {
		cm.confirmPasswordEditor.SetError("passwords do not match")
		cm.btnPositve.Background = cm.Theme.Color.Hint
		return false
	}

	cm.confirmPasswordEditor.SetError("")
	cm.btnPositve.Background = cm.Theme.Color.Primary
	return true
}

func (cm *CreatePasswordModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			t := cm.Theme.H6(cm.dialogTitle)
			t.Font.Weight = text.Bold
			return t.Layout(gtx)
		},
		func(gtx C) D {
			if cm.walletNameEnabled {
				return cm.walletName.Layout(gtx)
			}
			return layout.Dimensions{}
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(cm.passwordEditor.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding20, Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								if cm.showWalletWarnInfo {
									txt := cm.Theme.Label(values.MarginPadding14, "This spending password is for the new wallet only")
									txt.Color = cm.Theme.Color.Gray4
									return txt.Layout(gtx)
								}
								return layout.Dimensions{}
							}),
							layout.Rigid(func(gtx C) D {
								txt := cm.Theme.Label(values.MarginPadding14, strconv.Itoa(cm.passwordEditor.Editor.Len()))
								txt.Color = cm.Theme.Color.Gray4
								return layout.E.Layout(gtx, txt.Layout)
								// return txt.Layout(gtx)
							}),
						)
					})
				}),
			)
		},
		func(gtx C) D {
			return cm.passwordStrength.Layout(gtx)
		},
		func(gtx C) D {
			return cm.confirmPasswordEditor.Layout(gtx)
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
