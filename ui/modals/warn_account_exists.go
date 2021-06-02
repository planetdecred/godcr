package modals

import (
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/godcr/ui/values"
)

type WarnExistsMixerAccount struct {
	title string
	*common
}

const WarnExistsMixerAccountModal = "WarnExistsMixerAccount"

func (m *Modals) registerWarnExistsMixerAccountModal() {
	m.modals[WarnExistsMixerAccountModal] = &WarnExistsMixerAccount{
		title:  "Account Exists",
		common: m.common,
	}
}

func (m *WarnExistsMixerAccount) getTitle() string {
	return m.title
}

func (m *WarnExistsMixerAccount) onCancel()  {}
func (m *WarnExistsMixerAccount) onConfirm() {}

func (m *WarnExistsMixerAccount) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							m.alert.Color = m.theme.Color.DeepBlue
							return m.alert.Layout(gtx, values.MarginPadding50)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					label := m.theme.H6("Account name is taken")
					label.Font.Weight = text.Bold
					return label.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			txt := m.theme.Body1("There are existing accounts named mixed or unmixed. Please change the name to something else for now. You can change them back after the setup.")
			txt.Color = m.theme.Color.Gray
			return txt.Layout(gtx)
		},
	}
}
