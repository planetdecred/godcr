package modals

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/values"
)

type TransactionDetailsInfo struct {
	title     string
	*common
}

const TransactionDetailsInfoModal = "TransactionDetailsInfo"

func (m *Modals) registerTransactionDetailsInfoModal() {
	m.modals[TransactionDetailsInfoModal] = &TransactionDetailsInfo{
		title:     "Transaction Details",
	}
}

func (m *TransactionDetailsInfo) getTitle() string {
	return m.title
}

func (m *TransactionDetailsInfo) onCancel()  {}
func (m *TransactionDetailsInfo) onConfirm() {}

func (m *TransactionDetailsInfo) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					t := m.theme.Body1("Tap on")
					t.Color = m.theme.Color.Gray
					return t.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					t := m.theme.Body1("blue text")
					t.Color = m.theme.Color.Primary
					m := values.MarginPadding2
					return layout.Inset{
						Left:  m,
						Right: m,
					}.Layout(gtx, func(gtx C) D {
						return t.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					t := m.theme.Body1("to copy the item.")
					t.Color = m.theme.Color.Gray
					return t.Layout(gtx)
				}),
			)
		},
	}
}
