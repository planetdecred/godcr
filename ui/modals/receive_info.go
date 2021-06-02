package modals

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/values"
)

type ReceiveInfo struct {
	title string
	*common
}

const ReceiveInfoModal = "ReceiveInfo"

func (m *Modals) registerReceiveInfoModal() {
	m.modals[ReceiveInfoModal] = &ReceiveInfo{
		title:  "Receive Info",
		common: m.common,
	}
}

func (m *ReceiveInfo) getTitle() string {
	return m.title
}

func (m *ReceiveInfo) onCancel()  {}
func (m *ReceiveInfo) onConfirm() {}

func (m *ReceiveInfo) Layout(gtx layout.Context) []layout.Widget {
	text := m.theme.Label(values.TextSize20, "Each time you receive a payment, a new address is generated to protect your privacy.")
	text.Color = m.theme.Color.Gray
	
	return []layout.Widget{
		text.Layout,
	}
}
