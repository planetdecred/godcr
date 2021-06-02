package modals

import (
	"gioui.org/layout"
)

type SecurityToolsInfo struct {
	title string
	*common
}

const SecurityToolsInfoModal = "SecurityToolsInfo"

func (m *Modals) registerSecurityToolsInfoModal() {
	m.modals[SecurityToolsInfoModal] = &SecurityToolsInfo{
		title:  "Security Tools Info",
		common: m.common,
	}
}

func (m *SecurityToolsInfo) getTitle() string {
	return m.title
}

func (m *SecurityToolsInfo) onCancel()  {}
func (m *SecurityToolsInfo) onConfirm() {}

func (m *SecurityToolsInfo) Layout(gtx layout.Context) []layout.Widget {
	text := m.theme.Body1("Various tools that help in different aspects of crypto currency security will be located here.")
	text.Color = m.theme.Color.Gray

	return []layout.Widget{
		text.Layout,
	}
}
