package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWalletTemplate = "CreateWallet"
const RenameWalletTemplate = "RenameWallet"

type modalTemplate struct {
	walletName    decredmaterial.Editor
	password      decredmaterial.Editor
	matchPassword decredmaterial.Editor
	confirm       decredmaterial.Button
	cancel        decredmaterial.Button
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

func (win *Window) LoadTemplates(th *decredmaterial.Theme) *modalTemplate {
	return &modalTemplate{
		confirm:       th.Button(new(widget.Clickable), "Confirm"),
		cancel:        th.Button(new(widget.Clickable), "Cancel"),
		walletName:    th.Editor(new(widget.Editor), "Wallet name"),
		password:      th.Editor(new(widget.Editor), "Password"),
		matchPassword: th.Editor(new(widget.Editor), "Matching password"),
	}
}

func (m *modalTemplate) createNewWallet() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
		func(gtx C) D {
			m.password.Editor.Mask, m.password.Editor.SingleLine = '*', true
			return m.password.Layout(gtx)
		},
		func(gtx C) D {
			m.matchPassword.Editor.Mask, m.matchPassword.Editor.SingleLine = '*', true
			return m.matchPassword.Layout(gtx)
		},
	}
}

func (m *modalTemplate) renameWallet() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
	}
}

func (m *modalTemplate) Layout(th *decredmaterial.Theme, load *modalLoad) []func(gtx C) D {
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

func (m *modalTemplate) handle(th *decredmaterial.Theme, load *modalLoad) (template []func(gtx C) D) {
	switch load.template {
	case CreateWalletTemplate:
		t := load.confirm.(func(string, string))
		if m.editorsNotEmpty(th, m.walletName.Editor, m.password.Editor, m.matchPassword.Editor) &&
			m.passwordsMatch(th, m.password.Editor, m.matchPassword.Editor) &&
			m.confirm.Button.Clicked() {
			t(m.walletName.Editor.Text(), m.password.Editor.Text())
		}

		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.createNewWallet()
		return template
	case RenameWalletTemplate:
		t := load.confirm.(func(string))
		if m.editorsNotEmpty(th, m.walletName.Editor) && m.confirm.Button.Clicked() {
			t(m.walletName.Editor.Text())
		}

		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.renameWallet()
		return
	default:
		return
	}
}

// editorsNotEmpty checks that the editor fields are not emtpy. It returns false if they are empty and true if they are
// not and false if it doesn't. It sets the background of the confirm button to decredmaterial Hint color if fields
// are empty. It sets it to decredmaterial Primary color if they are not empty.
func (m *modalTemplate) editorsNotEmpty(th *decredmaterial.Theme, editors ...*widget.Editor) bool {
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
func (m *modalTemplate) passwordsMatch(th *decredmaterial.Theme, editors ...*widget.Editor) bool {
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
func (m *modalTemplate) resetFields() {
	m.matchPassword.Editor.SetText("")
	m.password.Editor.SetText("")
	m.walletName.Editor.SetText("")
}
