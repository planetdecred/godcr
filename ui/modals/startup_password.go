package modals

import (
	"gioui.org/layout"
)

type SetStartupPassword struct {
	title     string
	*common
}

const SetStartupPasswordModal = "SetStartupPassword"

func (m *Modals) registerSetStartupPasswordModal() {
	m.modals[SetStartupPasswordModal] = &SetStartupPassword{
		common: m.common,
		title:     "Startup Password",
	}
}

func (m *SetStartupPassword) getTitle() string {
	return m.title
}

func (m *SetStartupPassword) onCancel()  {}
func (m *SetStartupPassword) onConfirm() {}

func (m *SetStartupPassword) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		m.spendingPassword.Layout,
		m.passwordStrength.Layout,
		m.matchSpendingPassword.Layout,
	}
}
