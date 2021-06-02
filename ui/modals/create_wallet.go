package modals

import (
	"gioui.org/layout"
)

type CreateWallet struct {
	title string
	*common
}

const CreateWalletModal = "CreateWallet"

func (m *Modals) registerCreateWalletModal() {
	m.modals[CreateWalletModal] = &CreateWallet{
		title:  "Create New Wallet",
		common: m.common,
	}
}

func (m *CreateWallet) getTitle() string {
	return m.title
}

func (m *CreateWallet) onCancel()  {}
func (m *CreateWallet) onConfirm() {}

func (m *CreateWallet) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		m.walletName.Layout,
		m.spendingPassword.Layout,
		m.passwordStrength.Layout,
		m.matchSpendingPassword.Layout,
	}
}
