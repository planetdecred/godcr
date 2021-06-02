package modals

import (
	"gioui.org/layout"
)

type SendInfo struct {
	title string
	*common
}

const SendInfoModal = "SendInfo"

func (m *Modals) registerSendInfoModal() {
	m.modals[SendInfoModal] = &SendInfo{
		title:  "Send Info",
		common: m.common,
	}
}

func (m *SendInfo) getTitle() string {
	return m.title
}

func (m *SendInfo) onCancel()  {}
func (m *SendInfo) onConfirm() {}

func (m *SendInfo) Layout(gtx layout.Context) []layout.Widget {
	text := m.theme.Body1("Input or scan the destination wallet address and input the amount to send funds.")
	text.Color = m.theme.Color.Gray

	return []layout.Widget{
		text.Layout,
	}
}
