package modals

import (
	"gioui.org/layout"
)

type importWatchOnlyWallet struct {
	title string
	*common
}

const ImportWatchOnlyWalletModal = "ImportWatchOnlyWallet"

func (m *Modals) registerImportWatchOnlyWalletModal() {
	m.modals[ImportWatchOnlyWalletModal] = &importWatchOnlyWallet{
		common: m.common,
		title:  "Import Watch Only Wallet",
	}
}

func (m *importWatchOnlyWallet) getTitle() string {
	return m.title
}

func (m *importWatchOnlyWallet) onCancel()  {}
func (m *importWatchOnlyWallet) onConfirm() {}

func (m *importWatchOnlyWallet) Layout(gtx C) []layout.Widget {
	if m.walletName.Hint == "" {
		m.walletName.Hint = "Wallet name"
	}

	return []layout.Widget{
		m.walletName.Layout,
		m.extendedPublicKey.Layout,
	}
}
