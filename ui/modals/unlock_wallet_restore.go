package modals

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/values"
)

type  UnlockWalletRestore struct {
	title     string
	*common
}

const UnlockWalletRestoreModal = "UnlockWalletRestore"

func (m *Modals) registerUnlockWalletRestoreModal() {
	m.modals[UnlockWalletRestoreModal] = &UnlockWalletRestore{
		common: m.common,
		title:     "Unlock Wallet Restore",
	}
}

func (m *UnlockWalletRestore) getTitle() string {
	return m.title
}

func (m *UnlockWalletRestore) onCancel()  {}
func (m *UnlockWalletRestore) onConfirm() {}

func (m *UnlockWalletRestore) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			info := m.theme.Body1("The restoration process to discover your accounts was interrupted in the last sync.")
			info.Color = m.theme.Color.Gray
			return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return info.Layout(gtx)
			})
		},
		func(gtx C) D {
			info := m.theme.Body1("Unlock to resume the process.")
			info.Color = m.theme.Color.Gray
			return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return info.Layout(gtx)
			})
		},
		func(gtx C) D {
			return m.spendingPassword.Layout(gtx)
		},
	}
}
