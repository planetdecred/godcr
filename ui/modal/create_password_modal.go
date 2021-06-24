package modal

import (
	"fmt"
	"github.com/planetdecred/godcr/ui"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWallet = "create_wallet_modal"

type CreatePasswordModal struct {
	*ui.Common
	randomID string
	modal    decredmaterial.Modal

	dialogTitle string

	walletNameEnabled     bool
	walletName            decredmaterial.Editor
	passwordEditor        decredmaterial.Editor
	confirmPasswordEditor decredmaterial.Editor
	passwordStrength      decredmaterial.ProgressBarStyle

	isLoading      bool
	materialLoader material.LoaderStyle

	callback   func(walletName, password string, m *CreatePasswordModal) bool // return true to dismiss dialog
	btnPositve decredmaterial.Button

	btnNegative decredmaterial.Button
}

func NewCreatePasswordModal(common *ui.Common) *CreatePasswordModal {
	cm := &CreatePasswordModal{
		Common:       common,
		randomID:         fmt.Sprintf("%s-%d", CreateWallet, ui.GenerateRandomNumber()),
		modal:            *common.Theme.ModalFloatTitle(),
		passwordStrength: common.Theme.ProgressBar(0),
		btnPositve:       common.Theme.Button(new(widget.Clickable), "Confirm"),
		btnNegative:      common.Theme.Button(new(widget.Clickable), "Cancel"),
	}

	cm.btnPositve.TextSize, cm.btnNegative.TextSize = values.TextSize16, values.TextSize16
	cm.btnPositve.Font.Weight, cm.btnNegative.Font.Weight = text.Bold, text.Bold

	cm.walletName = common.Theme.Editor(new(widget.Editor), "Wallet name")
	cm.walletName.Editor.SingleLine, cm.walletName.Editor.Submit = true, true

	cm.passwordEditor = common.Theme.EditorPassword(new(widget.Editor), "Spending password")
	cm.passwordEditor.Editor.SingleLine, cm.passwordEditor.Editor.Submit = true, true

	cm.confirmPasswordEditor = common.Theme.EditorPassword(new(widget.Editor), "Spending password")
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

func (cm *CreatePasswordModal) PasswordCreated(callback func(walletName, password string, m *CreatePasswordModal) bool) *CreatePasswordModal {
	cm.callback = callback
	return cm
}

func (cm *CreatePasswordModal) SetLoading(loading bool) {
	cm.isLoading = loading
}

func (cm *CreatePasswordModal) SetError(err string) {

}

func (cm *CreatePasswordModal) handle() {
	if cm.passwordEditor.Editor.Text() == cm.confirmPasswordEditor.Editor.Text() {
		// reset error label when password and matching password fields match
		cm.confirmPasswordEditor.SetError("")
	}

	if cm.btnPositve.Button.Clicked() || ui.HandleSubmitEvent(cm.walletName.Editor, cm.passwordEditor.Editor, cm.confirmPasswordEditor.Editor) {

		nameValid := true
		if cm.walletNameEnabled {
			nameValid = ui.EditorsNotEmpty(cm.walletName.Editor)
		}

		if nameValid && ui.EditorsNotEmpty(cm.passwordEditor.Editor, cm.confirmPasswordEditor.Editor) &&
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

	ui.ComputePasswordStrength(&cm.passwordStrength, cm.Theme, cm.passwordEditor.Editor)

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

func (cm *CreatePasswordModal) Layout(gtx layout.Context) ui.D {
	w := []layout.Widget{
		func(gtx ui.C) ui.D {
			t := cm.Theme.H6(cm.dialogTitle)
			t.Font.Weight = text.Bold
			return t.Layout(gtx)
		},
		func(gtx ui.C) ui.D {
			if cm.walletNameEnabled {
				return cm.walletName.Layout(gtx)
			}
			return layout.Dimensions{}
		},
		func(gtx ui.C) ui.D {
			return cm.passwordEditor.Layout(gtx)
		},
		func(gtx ui.C) ui.D {
			return cm.passwordStrength.Layout(gtx)
		},
		func(gtx ui.C) ui.D {
			return cm.confirmPasswordEditor.Layout(gtx)
		},
		func(gtx ui.C) ui.D {
			return layout.E.Layout(gtx, func(gtx ui.C) ui.D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx ui.C) ui.D {

						cm.btnNegative.Background = cm.Theme.Color.Surface
						cm.btnNegative.Color = cm.Theme.Color.Primary
						return cm.btnNegative.Layout(gtx)
					}),
					layout.Rigid(func(gtx ui.C) ui.D {
						if cm.isLoading {
							return cm.materialLoader.Layout(gtx)
						}

						cm.btnPositve.Background, cm.btnPositve.Color = cm.Theme.Color.Surface, cm.Theme.Color.Primary
						return cm.btnPositve.Layout(gtx)
					}),
				)
			})
		},
	}

	return cm.modal.Layout(gtx, w, 850)
}
