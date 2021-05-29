package modals

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type ConfirmRemove struct {
	title string
	infoLabel decredmaterial.Label
}

const ConfirmRemoveModal = "ConfirmRemove"

func (m *Modals) registerConfirmRemoveWalletModal() {
	lbl := m.theme.Body1("Make sure to have the seed phrase backed up before removing the wallet")
	lbl.Color = m.theme.Color.Gray

	m.modals[ConfirmRemoveModal] = &ConfirmRemove{
		infoLabel: lbl,
		title: "Remove Wallet",
	}
}

func (m *ConfirmRemove) getTitle() string {
	return m.title
}

func (m *ConfirmRemove) onCancel() {}
func (m *ConfirmRemove) onConfirm() {}

func (m *ConfirmRemove) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget {
		m.infoLabel.Layout,
	}
}