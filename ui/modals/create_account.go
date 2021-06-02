package modals

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/values"
)

type CreateAccount struct {
	title string
	*common
}

const CreateAccountModal = "CreateAccount"

func (m *Modals) registerCreateAccountModal() {
	m.modals[CreateAccountModal] = &CreateAccount{
		title:  "Create New Account",
		common: m.common,
	}
}

func (m *CreateAccount) getTitle() string {
	return m.title
}

func (m *CreateAccount) onCancel()  {}
func (m *CreateAccount) onConfirm() {}

func (m *CreateAccount) Layout(gtx layout.Context) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, 
				layout.Rigid(func(gtx C) D {
					m.alert.Color = m.theme.Color.Gray
					return layout.Inset{Top: values.MarginPadding7, Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return m.alert.Layout(gtx, values.MarginPadding15)
					})
				}),
				layout.Rigid(func(gtx C) D {
					info := m.theme.Body1("Accounts")
					info.Color = m.theme.Color.Gray
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return info.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					info := m.theme.Body1(" cannot ")
					info.Color = m.theme.Color.DeepBlue
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return info.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					info := m.theme.Body1("be deleted when created")
					info.Color = m.theme.Color.Gray
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return info.Layout(gtx)
					})
				}),
			)
		},
		m.walletName.Layout,
		m.spendingPassword.Layout,
	}
}
