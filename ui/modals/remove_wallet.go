package modals

import (
	"gioui.org/layout"
)

type RemoveWallet struct {
	title string
	*common
}

const RemoveWalletModal = "RemoveWallet"

func (m *Modals) registerRemoveWalletModal() {
	m.modals[RemoveWalletModal] = &RemoveWallet{
		title:  "Remove Wallet",
		common: m.common,
	}
}

func (m *RemoveWallet) getTitle() string {
	return m.title
}

func (m *RemoveWallet) onCancel()  {}
func (m *RemoveWallet) onConfirm() {}

func (m *RemoveWallet) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			info := m.theme.Body1("Make sure to have the seed phrase backed up before removing the wallet")
			info.Color = m.theme.Color.Gray
			return info.Layout(gtx)
		},
	}
}
