package modals

import (
	"gioui.org/layout"
)

type RenameWallet struct {
	title string
	*common
}

const RenameWalletModal = "RenameWallet"

func (m *Modals) registerRenameWalletModal() {
	m.modals[RenameWalletModal] = &RenameWallet{
		title:  "Rename Wallet",
		common: m.common,
	}
}

func (m *RenameWallet) getTitle() string {
	return m.title
}

func (m *RenameWallet) onCancel()  {}
func (m *RenameWallet) onConfirm() {}

func (m *RenameWallet) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		m.walletName.Layout,
	}
}
