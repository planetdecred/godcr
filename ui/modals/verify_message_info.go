package modals

import (
	"gioui.org/layout"
)

type VerifyMessageInfo struct {
	title string
	*common
}

const VerifyMessageInfoModal = "VerifyMessageInfo"

func (m *Modals) registerVerifyMessageInfoModal() {
	m.modals[VerifyMessageInfoModal] = &VerifyMessageInfo{
		title:  "Verify Message Info",
		common: m.common,
	}
}

func (m *VerifyMessageInfo) getTitle() string {
	return m.title
}

func (m *VerifyMessageInfo) onCancel()  {}
func (m *VerifyMessageInfo) onConfirm() {}

func (m *VerifyMessageInfo) Layout(gtx layout.Context) []layout.Widget {
	text := m.theme.Body1("After you or your counterparty has genrated a signature, you can use this form to verify the" +
				" validity of the  signature. \n \nOnce you have entered the address, the message and the corresponding " +
				"signature, you will see VALID if the signature appropriately matches the address and message, otherwise INVALID.")
	text.Color = m.theme.Color.Gray

	return []layout.Widget{
		text.Layout,
	}
}
