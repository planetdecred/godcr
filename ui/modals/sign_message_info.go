package modals

import (
	"gioui.org/layout"
)

type SignMessageInfo struct {
	title string
	*common
}

const SignMessageInfoModal = "SignMessageInfo"

func (m *Modals) registerSignMessageInfoModal() {
	m.modals[SignMessageInfoModal] = &SignMessageInfo{
		title:  "Sign Message Info",
		common: m.common,
	}
}

func (m *SignMessageInfo) getTitle() string {
	return m.title
}

func (m *SignMessageInfo) onCancel()  {}
func (m *SignMessageInfo) onConfirm() {}

func (m *SignMessageInfo) Layout(gtx layout.Context) []layout.Widget {
	text := m.theme.Body1("Signing a message with an address' private key allows you to prove that you are the owner of a given address" +
				" to a possible counterparty.")
			text.Color = m.theme.Color.Gray

	return []layout.Widget{
		text.Layout,
	}
}
