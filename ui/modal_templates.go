package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWalletTemplate = "CreateWallet"
const RenameWalletTemplate = "RenameWallet"
const CreateAccountTemplate = "CreateNewAccount"
const RenameAccountTemplate = "RenameAccount"
const PasswordTemplate = "Password"
const ConfirmRemoveTemplate = "ConfirmRemove"

type ModalTemplate struct {
	walletName            decredmaterial.Editor
	spendingPassword      decredmaterial.Editor
	matchSpendingPassword decredmaterial.Editor
	confirm               decredmaterial.Button
	cancel                decredmaterial.Button
}

type modalLoad struct {
	template    string
	title       string
	confirm     interface{}
	confirmText string
	cancel      interface{}
	cancelText  string
	isReset     bool
}

func (win *Window) LoadModalTemplates() *ModalTemplate {
	return &ModalTemplate{
		confirm:               win.theme.Button(new(widget.Clickable), "Confirm"),
		cancel:                win.theme.Button(new(widget.Clickable), "Cancel"),
		walletName:            win.theme.Editor(new(widget.Editor), ""),
		spendingPassword:      win.theme.Editor(new(widget.Editor), "Spending password"),
		matchSpendingPassword: win.theme.Editor(new(widget.Editor), "Confirm spending password"),
	}
}

func (m *ModalTemplate) createNewWallet() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
		func(gtx C) D {
			m.spendingPassword.Editor.Mask, m.spendingPassword.Editor.SingleLine = '*', true
			return m.spendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			m.matchSpendingPassword.Editor.Mask, m.matchSpendingPassword.Editor.SingleLine = '*', true
			return m.matchSpendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) renameWallet() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) createNewAccount() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
		func(gtx C) D {
			m.spendingPassword.Editor.Mask, m.spendingPassword.Editor.SingleLine = '*', true
			return m.spendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) removeWallet(th *decredmaterial.Theme) []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return th.Body2("Make sure to have the seed phrase backed up before removing the wallet").Layout(gtx)
		},
	}
}

func (m *ModalTemplate) Password() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			m.spendingPassword.Editor.Mask, m.spendingPassword.Editor.SingleLine = '*', true
			return m.spendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) Layout(th *decredmaterial.Theme, load *modalLoad) []func(gtx C) D {
	if !load.isReset {
		m.resetFields()
		load.isReset = true
	}

	title := []func(gtx C) D{
		func(gtx C) D {
			return th.H6(load.title).Layout(gtx)
		},
	}

	action := []func(gtx C) D{
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
							m.cancel.Text = load.cancelText
							m.cancel.Background = th.Color.Surface
							m.cancel.Color = th.Color.Primary
							return m.cancel.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
							m.confirm.Text = load.confirmText
							m.confirm.Background, m.confirm.Color = th.Color.Primary, th.Color.Surface
							if load.template == ConfirmRemoveTemplate {
								m.confirm.Background, m.confirm.Color = th.Color.Surface, th.Color.Danger
							}
							return m.confirm.Layout(gtx)
						})
					}),
				)
			})
		},
	}

	w := m.handle(th, load)
	w = append(title, w...)
	w = append(w, action...)
	return w
}

func (m *ModalTemplate) handle(th *decredmaterial.Theme, load *modalLoad) (template []func(gtx C) D) {
	switch load.template {
	case CreateWalletTemplate:
		if m.editorsNotEmpty(th, m.walletName.Editor, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) &&
			m.passwordsMatch(th, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) &&
			m.confirm.Button.Clicked() {
			load.confirm.(func(string, string))(m.walletName.Editor.Text(), m.spendingPassword.Editor.Text())
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		m.walletName.Hint = "Wallet name"
		template = m.createNewWallet()
		return template
	case RenameWalletTemplate, RenameAccountTemplate:
		if m.editorsNotEmpty(th, m.walletName.Editor) && m.confirm.Button.Clicked() {
			load.confirm.(func(string))(m.walletName.Editor.Text())
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		template = m.renameWallet()
		m.walletName.Hint = "Wallet name"
		if load.template == RenameAccountTemplate {
			m.walletName.Hint = "Account name"
		}
		return
	case CreateAccountTemplate:
		if m.editorsNotEmpty(th, m.walletName.Editor, m.spendingPassword.Editor) && m.confirm.Button.Clicked() {
			load.confirm.(func(string, string))(m.walletName.Editor.Text(), m.spendingPassword.Editor.Text())
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		template = m.createNewAccount()
		m.walletName.Hint = "Account name"
		return
	case PasswordTemplate:
		if m.editorsNotEmpty(th, m.spendingPassword.Editor) && m.confirm.Button.Clicked() {
			load.confirm.(func(string))(m.spendingPassword.Editor.Text())
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		template = m.Password()
		return
	case ConfirmRemoveTemplate:
		if m.confirm.Button.Clicked() {
			load.confirm.(func())()
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.removeWallet(th)
		return
	default:
		return
	}
}

// editorsNotEmpty checks that the editor fields are not empty. It returns false if they are empty and true if they are
// not and false if it doesn't. It sets the background of the confirm button to decredmaterial Hint color if fields
// are empty. It sets it to decredmaterial Primary color if they are not empty.
func (m *ModalTemplate) editorsNotEmpty(th *decredmaterial.Theme, editors ...*widget.Editor) bool {
	for _, e := range editors {
		if e.Text() == "" {
			m.confirm.Background = th.Color.Hint
			return false
		}
	}

	m.confirm.Background = th.Color.Primary
	return true
}

// passwordMatches checks that the password and matching password field matches. It returns true if it matches
// and false if it doesn't. It sets the background of the confirm button to decredmaterial Hint color if the passwords
// don't match. It sets it to decredmaterial Primary color if the passwords match.
func (m *ModalTemplate) passwordsMatch(th *decredmaterial.Theme, editors ...*widget.Editor) bool {
	if len(editors) < 2 {
		return false
	}

	passWord := editors[0]
	matching := editors[1]
	if passWord.Text() != matching.Text() {
		m.confirm.Background = th.Color.Hint
		return false
	}
	m.confirm.Background = th.Color.Primary
	return true
}

// resetFields clears all modal fields when the modal is closed
func (m *ModalTemplate) resetFields() {
	m.matchSpendingPassword.Editor.SetText("")
	m.spendingPassword.Editor.SetText("")
	m.walletName.Editor.SetText("")
}
