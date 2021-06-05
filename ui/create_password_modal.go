package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const ModalCreateWallet = "create_wallet_modal"

type createPasswordModal struct {
	*pageCommon

	modal decredmaterial.Modal

	dialogTitle string

	walletNameEnabled     bool
	walletName            decredmaterial.Editor
	passwordEditor        decredmaterial.Editor
	confirmPasswordEditor decredmaterial.Editor
	passwordStrength      decredmaterial.ProgressBarStyle

	isLoading      bool
	materialLoader material.LoaderStyle

	callback   func(walletName, password string, m *createPasswordModal) bool // return true to dismiss dialog
	btnPositve decredmaterial.Button

	btnNegative decredmaterial.Button
}

func newCreatePasswordModal(common *pageCommon) *createPasswordModal {
	cm := &createPasswordModal{
		pageCommon:       common,
		modal:            *common.theme.ModalFloatTitle(),
		passwordStrength: common.theme.ProgressBar(0),
		btnPositve:       common.theme.Button(new(widget.Clickable), "Confirm"),
		btnNegative:      common.theme.Button(new(widget.Clickable), "Cancel"),
	}

	cm.btnPositve.TextSize, cm.btnNegative.TextSize = values.TextSize16, values.TextSize16
	cm.btnPositve.Font.Weight, cm.btnNegative.Font.Weight = text.Bold, text.Bold

	cm.walletName = common.theme.Editor(new(widget.Editor), "Wallet name")
	cm.walletName.Editor.SingleLine, cm.walletName.Editor.Submit = true, true

	cm.passwordEditor = common.theme.EditorPassword(new(widget.Editor), "Spending password")
	cm.passwordEditor.Editor.SingleLine, cm.passwordEditor.Editor.Submit = true, true

	cm.confirmPasswordEditor = common.theme.EditorPassword(new(widget.Editor), "Spending password")
	cm.confirmPasswordEditor.Editor.SingleLine, cm.confirmPasswordEditor.Editor.Submit = true, true

	th := material.NewTheme(gofont.Collection())
	cm.materialLoader = material.Loader(th)

	return cm
}

func (cm *createPasswordModal) modalID() string {
	return ModalCreateWallet + cm.dialogTitle // TODO
}

func (cm *createPasswordModal) OnResume() {
}

func (cm *createPasswordModal) OnDismiss() {

}

func (cm *createPasswordModal) show() {
	cm.showModal(cm)
}

func (cm *createPasswordModal) dismiss() {
	cm.dismissModal(cm)
}

func (cm *createPasswordModal) title(title string) *createPasswordModal {
	cm.dialogTitle = title
	return cm
}

func (cm *createPasswordModal) enableName(enable bool) *createPasswordModal {
	cm.walletNameEnabled = enable
	return cm
}

func (cm *createPasswordModal) passwordHint(hint string) *createPasswordModal {
	cm.passwordEditor.Hint = hint
	return cm
}

func (cm *createPasswordModal) confirmPasswordHint(hint string) *createPasswordModal {
	cm.confirmPasswordEditor.Hint = hint
	return cm
}

func (cm *createPasswordModal) passwordCreated(callback func(walletName, password string, m *createPasswordModal) bool) *createPasswordModal {
	cm.callback = callback
	return cm
}

func (cm *createPasswordModal) setLoading(loading bool) {
	cm.isLoading = loading
}

func (cm *createPasswordModal) setError(err string) {

}

func (cm *createPasswordModal) handle() {
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

			cm.setLoading(true)
			if cm.callback(cm.walletName.Editor.Text(), cm.passwordEditor.Editor.Text(), cm) {
				cm.dismiss()
			}
		}

	}

	if cm.btnNegative.Button.Clicked() {
		if !cm.isLoading {
			cm.dismiss()
		}
	}

	computePasswordStrength(&cm.passwordStrength, cm.theme, cm.passwordEditor.Editor)

}
func (cm *createPasswordModal) passwordsMatch(editors ...*widget.Editor) bool {
	if len(editors) < 2 {
		return false
	}

	password := editors[0]
	matching := editors[1]

	if password.Text() != matching.Text() {
		cm.confirmPasswordEditor.SetError("passwords do not match")
		cm.btnPositve.Background = cm.theme.Color.Hint
		return false
	}

	cm.confirmPasswordEditor.SetError("")
	cm.btnPositve.Background = cm.theme.Color.Primary
	return true
}

func (cm *createPasswordModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			t := cm.theme.H6(cm.dialogTitle)
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
			return cm.passwordEditor.Layout(gtx)
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

						cm.btnNegative.Background = cm.theme.Color.Surface
						cm.btnNegative.Color = cm.theme.Color.Primary
						return cm.btnNegative.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if cm.isLoading {
							return cm.materialLoader.Layout(gtx)
						}

						cm.btnPositve.Background, cm.btnPositve.Color = cm.theme.Color.Surface, cm.theme.Color.Primary
						return cm.btnPositve.Layout(gtx)
					}),
				)
			})
		},
	}

	return cm.modal.Layout(gtx, w, 850)
}
