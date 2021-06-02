package modals

import (
	"gioui.org/layout"
)

type RescanWallet struct {
	title string
	*common
}

const RescanWalletModal = "RescanWallet"

func (m *Modals) registerRescanWalletModal() {
	m.modals[RescanWalletModal] = &RescanWallet{
		title:  "Rescan Wallet",
		common: m.common,
	}
}

func (m *RescanWallet) getTitle() string {
	return m.title
}

func (m *RescanWallet) onCancel()  {}
func (m *RescanWallet) onConfirm() {}

func (m *RescanWallet) Layout(gtx layout.Context) []layout.Widget {
	text := m.theme.Body1("Rescanning may help resolve some balance errors. This will take some time, as it scans the entire" +
				" blockchain for transactions")
	text.Color = m.theme.Color.Gray

	return []layout.Widget{
		text.Layout,
	}
}
