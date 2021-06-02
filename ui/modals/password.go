package modals

import (
	"gioui.org/layout"
)

type Password struct {
	title string
	*common
}

const PasswordModal = "Password"

func (m *Modals) registerPasswordModal() {
	m.modals[PasswordModal] = &Password{
		title:  "Password",
		common: m.common,
	}
}

func (m *Password) getTitle() string {
	return m.title
}

func (m *Password) onCancel()  {}
func (m *Password) onConfirm() {}

func (m *Password) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return m.spendingPassword.Layout(gtx)
		},
	}
}
