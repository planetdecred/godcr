package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const CreateWalletTemplate = "CreateWallet"
const RenameWalletTemplate = "RenameWallet"
const CreateAccountTemplate = "CreateNewAccount"
const RenameAccountTemplate = "RenameAccount"
const PasswordTemplate = "Password"
const ChangePasswordTemplate = "ChangePassword"
const ConfirmRemoveTemplate = "ConfirmRemove"
const VerifyMessageInfoTemplate = "VerifyMessageInfo"
const SignMessageInfoTemplate = "SignMessageInfo"
const PrivacyInfoTemplate = "PrivacyInfo"
const RescanWalletTemplate = "RescanWallet"
const ChangeStartupPasswordTemplate = "ChangeStartupPassword"
const SetStartupPasswordTemplate = "SetStartupPassword"
const RemoveStartupPasswordTemplate = "RemoveStartupPassword"
const UnlockWalletTemplate = "UnlockWallet"
const ConnectToSpecificPeerTemplate = "ConnectToSpecificPeer"
const ChangeSpecificPeerTemplate = "ChangeSpecificPeer"
const UserAgentTemplate = "UserAgent"
const ConfirmSetupMixerTemplate = "ConfirmSetupMixer"
const ConfirmSetupMixerAcctTemplate = "SetupMixerAcctTemplate"
const ConfirmMixerAcctExistTemplate = "MixerAcctExistTemplate"
const SecurityToolsInfoTemplate = "SecurityToolsInfo"
const ImportWatchOnlyWalletTemplate = "ImportWatchOnlyWallet"

type ModalTemplate struct {
	th                    *decredmaterial.Theme
	walletName            decredmaterial.Editor
	oldSpendingPassword   decredmaterial.Editor
	spendingPassword      decredmaterial.Editor
	matchSpendingPassword decredmaterial.Editor
	extendedPublicKey     decredmaterial.Editor
	confirm               decredmaterial.Button
	cancel                decredmaterial.Button
	alert                 decredmaterial.IconButton
	alertError            *widget.Image
	passwordStgth         decredmaterial.ProgressBarStyle
}

type modalLoad struct {
	template       string
	customTemplate string
	title          string
	confirm        interface{}
	confirmText    string
	cancel         interface{}
	cancelText     string
	isReset        bool
}

func (win *Window) LoadModalTemplates() *ModalTemplate {
	icon := win.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.AlertError)))
	icon.Size = values.MarginPadding20
	icon.Inset = layout.UniformInset(values.MarginPadding5)
	icon.Color = win.theme.Color.Gray
	icon.Background = win.theme.Color.Surface

	return &ModalTemplate{
		th:                    win.theme,
		confirm:               win.theme.Button(new(widget.Clickable), "Confirm"),
		cancel:                win.theme.Button(new(widget.Clickable), "Cancel"),
		walletName:            win.theme.Editor(new(widget.Editor), ""),
		oldSpendingPassword:   win.theme.Editor(new(widget.Editor), "Old spending password"),
		spendingPassword:      win.theme.Editor(new(widget.Editor), "Spending password"),
		matchSpendingPassword: win.theme.Editor(new(widget.Editor), "Confirm spending password"),
		extendedPublicKey:     win.theme.Editor(new(widget.Editor), "Extended public key"),
		alert:                 icon,
		passwordStgth:         win.theme.ProgressBar(0),
	}
}

func (m *ModalTemplate) importWatchOnlyWallet() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
		func(gtx C) D {
			return m.extendedPublicKey.Layout(gtx)
		},
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
			return m.passwordStgth.Layout(gtx)
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

func (m *ModalTemplate) createNewAccount(th *decredmaterial.Theme) []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return m.alert.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					info := th.Body1("Accounts")
					info.Color = th.Color.Gray
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return info.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					info := th.Body1(" cannot ")
					info.Color = th.Color.Black
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return info.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					info := th.Body1("be deleted when created")
					info.Color = th.Color.Gray
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return info.Layout(gtx)
					})
				}),
			)
		},
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
			info := th.Body1("Make sure to have the seed phrase backed up before removing the wallet")
			info.Color = th.Color.Gray
			return info.Layout(gtx)
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

func (m *ModalTemplate) changePassword() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			m.oldSpendingPassword.Editor.Mask, m.oldSpendingPassword.Editor.SingleLine = '*', true
			return m.oldSpendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			m.spendingPassword.Editor.Mask, m.spendingPassword.Editor.SingleLine = '*', true
			return m.spendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			return m.passwordStgth.Layout(gtx)
		},
		func(gtx C) D {
			m.matchSpendingPassword.Editor.Mask, m.matchSpendingPassword.Editor.SingleLine = '*', true
			return m.matchSpendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) setStartupPassword() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			m.spendingPassword.Editor.Mask, m.spendingPassword.Editor.SingleLine = '*', true
			return m.spendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			return m.passwordStgth.Layout(gtx)
		},
		func(gtx C) D {
			m.matchSpendingPassword.Editor.Mask, m.matchSpendingPassword.Editor.SingleLine = '*', true
			return m.matchSpendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) verifyMessageInfo() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			text := m.th.Body1("After you or your counterparty has genrated a signature, you can use this form to verify the" +
				" validity of the  signature. \n \nOnce you have entered the address, the message and the corresponding " +
				"signature, you will see VALID if the signature appropriately matches the address and message, otherwise INVALID.")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) signMessageInfo() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			text := m.th.Body1("Signing a message with an address' private key allows you to prove that you are the owner of a given address" +
				" to a possible counterparty.")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) rescanWallet() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			text := m.th.Body1("Rescanning may help resolve some balance errors. This will take some time, as it scans the entire" +
				" blockchain for transactions")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) securityToolsInfo() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			text := m.th.Body1("Various tools that help in different aspects of crypto currency security will be located here.")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) privacyInfo() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := mustIcon(widget.NewIcon(icons.ImageLens))
					ic.Color = m.th.Color.Gray
					return ic.Layout(gtx, values.TextSize8)
				}),
				layout.Rigid(func(gtx C) D {
					text := m.th.Body1("When you turn on the mixer, your unmixed DCRs in this wallet (unmixed balance) will be gradually mixed.")
					text.Color = m.th.Color.Gray
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
		func(gtx C) D {
			text := m.th.Label(values.TextSize18, "Important: keep this app opened while mixer is running.")
			return text.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := mustIcon(widget.NewIcon(icons.ImageLens))
					ic.Color = m.th.Color.Gray
					return ic.Layout(gtx, values.TextSize8)
				}),
				layout.Rigid(func(gtx C) D {
					text := m.th.Body1("Mixer will automatically stop when unmixed balance are fully mixed.")
					text.Color = m.th.Color.Gray
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
	}
}

func (m *ModalTemplate) setupMixerInfo() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			txt := m.th.Body1("Two dedicated accounts (“mixed” & “unmixed”) will be created in order to use the mixer.")
			txt.Color = m.th.Color.Gray
			return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, txt.Layout)
		},
		func(gtx C) D {
			txt := m.th.Label(values.TextSize18, "This action cannot be undone.")
			return txt.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) setupMixerAcct() []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			m.spendingPassword.Editor.Mask, m.spendingPassword.Editor.SingleLine = '*', true
			return m.spendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) warnExistMixerAcct(load *modalLoad) []func(gtx C) D {
	return append([]func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							m.alertError.Scale = 0.07
							return m.alertError.Layout(gtx)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					return m.th.H5("Account name is taken").Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			txt := m.th.Body1("There are existing accounts named mixed or unmixed. Please change the name to something else for now. You can change them back after the setup.")
			txt.Color = m.th.Color.Gray
			return txt.Layout(gtx)
		},
	}, m.actions(m.th, load)...)
}

func (m *ModalTemplate) Layout(th *decredmaterial.Theme, load *modalLoad) []func(gtx C) D {
	if load.customTemplate != "" {
		return m.handleCustomTemplate(load)
	}

	if !load.isReset {
		m.resetFields()
		load.isReset = true
	}

	title := []func(gtx C) D{
		func(gtx C) D {
			return th.H5(load.title).Layout(gtx)
		},
	}

	w := m.handle(th, load)
	w = append(title, w...)
	w = append(w, m.actions(th, load)...)
	return w
}

func (m *ModalTemplate) actions(th *decredmaterial.Theme, load *modalLoad) []func(gtx C) D {
	return []func(gtx C) D{
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if load.cancelText == "" {
							return layout.Dimensions{}
						}
						return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
							m.cancel.Text = load.cancelText
							m.cancel.Background = th.Color.Surface
							m.cancel.Color = th.Color.Primary
							return m.cancel.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						if load.confirmText == "" {
							return layout.Dimensions{}
						}
						return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
							m.confirm.Text = load.confirmText
							if load.template == ConfirmRemoveTemplate {
								m.confirm.Background, m.confirm.Color = th.Color.Surface, th.Color.Danger
							}
							if load.template == RescanWalletTemplate {
								m.confirm.Background, m.confirm.Color = th.Color.Surface, th.Color.Primary
							}
							return m.confirm.Layout(gtx)
						})
					}),
				)
			})
		},
	}
}

func (m *ModalTemplate) handle(th *decredmaterial.Theme, load *modalLoad) (template []func(gtx C) D) {
	switch load.template {
	case CreateWalletTemplate:
		if m.spendingPassword.Editor.Text() == m.matchSpendingPassword.Editor.Text() {
			// reset error label when password and matching password fields match
			m.matchSpendingPassword.SetError("")
		}

		if m.editorsNotEmpty(th, m.walletName.Editor, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) &&
			m.confirm.Button.Clicked() {
			if m.passwordsMatch(m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
				load.confirm.(func(string, string))(m.walletName.Editor.Text(), m.spendingPassword.Editor.Text())
			}
		}

		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		m.passwordStrength(th, m.spendingPassword.Editor)

		template = m.createNewWallet()
		m.walletName.Hint = "Wallet name"
		return
	case RenameWalletTemplate, RenameAccountTemplate, ConnectToSpecificPeerTemplate, ChangeSpecificPeerTemplate, UserAgentTemplate:
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
		if load.template == ConnectToSpecificPeerTemplate || load.template == ChangeSpecificPeerTemplate {
			m.walletName.Hint = "IP address"
		}
		if load.template == UserAgentTemplate {
			m.walletName.Hint = "User agent"
		}
		return
	case CreateAccountTemplate:
		if m.editorsNotEmpty(th, m.walletName.Editor, m.spendingPassword.Editor) && m.confirm.Button.Clicked() {
			load.confirm.(func(string, string))(m.walletName.Editor.Text(), m.spendingPassword.Editor.Text())
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		template = m.createNewAccount(th)
		m.walletName.Hint = "Account name"
		return
	case PasswordTemplate, UnlockWalletTemplate, RemoveStartupPasswordTemplate:
		if m.editorsNotEmpty(th, m.spendingPassword.Editor) && m.confirm.Button.Clicked() {
			load.confirm.(func(string))(m.spendingPassword.Editor.Text())
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		m.spendingPassword.Hint = "Spending password"
		if load.template == RemoveStartupPasswordTemplate || load.template == UnlockWalletTemplate {
			m.spendingPassword.Hint = "Startup password"
		}

		template = m.Password()
		return
	case ChangePasswordTemplate, ChangeStartupPasswordTemplate:
		if m.spendingPassword.Editor.Text() == m.matchSpendingPassword.Editor.Text() {
			// reset error label when password and matching password fields match
			m.matchSpendingPassword.SetError("")
		}

		if m.editorsNotEmpty(th, m.oldSpendingPassword.Editor, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) &&
			m.confirm.Button.Clicked() {
			if m.passwordsMatch(m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
				load.confirm.(func(string, string))(m.oldSpendingPassword.Editor.Text(), m.spendingPassword.Editor.Text())
			}
		}

		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		m.passwordStrength(th, m.spendingPassword.Editor)

		m.spendingPassword.Hint = "New spending password"
		m.matchSpendingPassword.Hint = "Confirm new spending password"
		if load.template == ChangeStartupPasswordTemplate {
			m.oldSpendingPassword.Hint = "Old startup password"
			m.spendingPassword.Hint = "New startup password"
			m.matchSpendingPassword.Hint = "Confirm new startup password"
		}

		template = m.changePassword()
		return
	case ImportWatchOnlyWalletTemplate:
		if m.confirm.Button.Clicked() {
			load.confirm.(func(string, string))(m.walletName.Editor.Text(), m.extendedPublicKey.Editor.Text())
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		m.walletName.Hint = "Wallet name"

		template = m.importWatchOnlyWallet()
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
	case VerifyMessageInfoTemplate:
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.verifyMessageInfo()
		return
	case SignMessageInfoTemplate:
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.signMessageInfo()
		return
	case RescanWalletTemplate:
		if m.confirm.Button.Clicked() {
			load.confirm.(func())()
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.rescanWallet()
		return
	case PrivacyInfoTemplate:
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.privacyInfo()
		return
	case SetStartupPasswordTemplate:
		if m.spendingPassword.Editor.Text() == m.matchSpendingPassword.Editor.Text() {
			// reset error label when password and matching password fields match
			m.matchSpendingPassword.SetError("")
		}

		if m.editorsNotEmpty(th, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) &&
			m.confirm.Button.Clicked() {
			if m.passwordsMatch(m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
				load.confirm.(func(string))(m.spendingPassword.Editor.Text())
			}
		}

		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		m.passwordStrength(th, m.spendingPassword.Editor)
		m.spendingPassword.Hint = "Startup password"
		m.matchSpendingPassword.Hint = "Confirm startup password"

		template = m.setStartupPassword()
		return
	case ConfirmSetupMixerTemplate:
		if m.confirm.Button.Clicked() {
			load.confirm.(func())()
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}

		template = m.setupMixerInfo()
		return
	case ConfirmSetupMixerAcctTemplate:
		if m.editorsNotEmpty(th, m.spendingPassword.Editor) && m.confirm.Button.Clicked() {
			load.confirm.(func(string))(m.spendingPassword.Editor.Text())
		}
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.setupMixerAcct()
		return
	case SecurityToolsInfoTemplate:
		if m.cancel.Button.Clicked() {
			load.cancel.(func())()
		}
		template = m.securityToolsInfo()
		return
	default:
		return
	}
}

func (m *ModalTemplate) handleCustomTemplate(load *modalLoad) (template []func(gtx C) D) {
	switch load.customTemplate {
	case ConfirmMixerAcctExistTemplate:
		if m.confirm.Button.Clicked() {
			load.confirm.(func())()
		}
		template = m.warnExistMixerAcct(load)
		return
	default:
		return
	}
}

// editorsNotEmpty checks that the editor fields are not empty. It returns false if they are empty and true if they are
// not and false if it doesn't. It sets the background of the confirm button to decredmaterial Hint color if fields
// are empty. It sets it to decredmaterial Primary color if they are not empty.
func (m *ModalTemplate) editorsNotEmpty(th *decredmaterial.Theme, editors ...*widget.Editor) bool {
	m.confirm.Color = th.Color.Surface
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
func (m *ModalTemplate) passwordsMatch(editors ...*widget.Editor) bool {
	if len(editors) < 2 {
		return false
	}

	password := editors[0]
	matching := editors[1]

	if password.Text() != matching.Text() {
		m.matchSpendingPassword.SetError("passwords do not match")
		return false
	}

	m.matchSpendingPassword.SetError("")
	return true
}

func (m *ModalTemplate) passwordStrength(th *decredmaterial.Theme, editors ...*widget.Editor) {
	password := editors[0]
	strength := (dcrlibwallet.ShannonEntropy(password.Text()) / 4.0)
	m.passwordStgth.Progress = float32(strength * 100)
	m.passwordStgth.Color = th.Color.Success
}

// resetFields clears all modal fields when the modal is closed
func (m *ModalTemplate) resetFields() {
	m.matchSpendingPassword.Editor.SetText("")
	m.spendingPassword.Editor.SetText("")
	m.walletName.Editor.SetText("")
	m.matchSpendingPassword.SetError("")
	m.oldSpendingPassword.Editor.SetText("")
	m.extendedPublicKey.Editor.SetText("")
}
