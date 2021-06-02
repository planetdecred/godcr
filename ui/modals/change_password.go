package modals

import (
	"gioui.org/layout"
)

type ChangePassword struct {
	title     string
	*common
}

const ChangePasswordModal = "ChangePassword"

func (m *Modals) registerChangePasswordModal() {
	m.modals[ChangePasswordModal] = &ChangePassword{
		common: m.common,
		title:     "Change Password",
	}
}

func (m *ChangePassword) getTitle() string {
	return m.title
}

func (m *ChangePassword) onCancel()  {}
func (m *ChangePassword) onConfirm() {}

func (m *ChangePassword) Layout(gtx layout.Context) []layout.Widget {
	m.oldSpendingPassword.Editor.SingleLine = true

	return []layout.Widget{
		m.oldSpendingPassword.Layout,
		m.spendingPassword.Layout,
		m.passwordStrength.Layout,
		m.matchSpendingPassword.Layout,
	}
}
