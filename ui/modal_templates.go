package ui

import (
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const (
	CreateWalletTemplate           = "CreateWallet"
	RenameWalletTemplate           = "RenameWallet"
	CreateAccountTemplate          = "CreateNewAccount"
	RenameAccountTemplate          = "RenameAccount"
	PasswordTemplate               = "Password"
	ChangePasswordTemplate         = "ChangePassword"
	ConfirmRemoveTemplate          = "ConfirmRemove"
	VerifyMessageInfoTemplate      = "VerifyMessageInfo"
	SignMessageInfoTemplate        = "SignMessageInfo"
	PrivacyInfoTemplate            = "PrivacyInfo"
	RescanWalletTemplate           = "RescanWallet"
	ChangeStartupPasswordTemplate  = "ChangeStartupPassword"
	SetStartupPasswordTemplate     = "SetStartupPassword"
	RemoveStartupPasswordTemplate  = "RemoveStartupPassword"
	UnlockWalletTemplate           = "UnlockWallet"
	ConnectToSpecificPeerTemplate  = "ConnectToSpecificPeer"
	ChangeSpecificPeerTemplate     = "ChangeSpecificPeer"
	UserAgentTemplate              = "UserAgent"
	SetupMixerInfoTemplate         = "ConfirmSetupMixer"
	ConfirmMixerAcctExistTemplate  = "MixerAcctExistTemplate"
	SecurityToolsInfoTemplate      = "SecurityToolsInfo"
	ImportWatchOnlyWalletTemplate  = "ImportWatchOnlyWallet"
	UnlockWalletRestoreTemplate    = "UnlockWalletRestoreTemplate"
	SendInfoTemplate               = "SendInfo"
	ReceiveInfoTemplate            = "ReceiveInfo"
	TransactionDetailsInfoTemplate = "TransactionDetailsInfoInfo"
)

type ModalTemplate struct {
	th                    *decredmaterial.Theme
	walletName            decredmaterial.Editor
	oldSpendingPassword   decredmaterial.Editor
	spendingPassword      decredmaterial.Editor
	matchSpendingPassword decredmaterial.Editor
	extendedPublicKey     decredmaterial.Editor
	confirm               decredmaterial.Button
	cancel                decredmaterial.Button
	alert                 *widget.Icon
	passwordStrength      decredmaterial.ProgressBarStyle
	keyEvent              chan *key.Event
	refreshWindow         func()
	confirmKeyPressed     bool
}

type modalLoad struct {
	template    string
	title       string
	confirm     interface{}
	loading     bool
	confirmText string
	cancel      func()
	cancelText  string
	isReset     bool
}

func (win *Window) LoadModalTemplates() *ModalTemplate {
	cancel := win.theme.Button(new(widget.Clickable), "Cancel")
	confirm := win.theme.Button(new(widget.Clickable), "Confirm")
	cancel.TextSize, confirm.TextSize = values.TextSize16, values.TextSize16

	spendingPassword := win.theme.EditorPassword(new(widget.Editor), "Spending password")
	spendingPassword.Editor.SingleLine, spendingPassword.Editor.Submit = true, true

	matchSpendingPassword := win.theme.EditorPassword(new(widget.Editor), "Confirm spending password")
	matchSpendingPassword.Editor.Submit, matchSpendingPassword.Editor.SingleLine = true, true

	oldSpendingPassword := win.theme.EditorPassword(new(widget.Editor), "Old spending password")
	oldSpendingPassword.Editor.Submit, oldSpendingPassword.Editor.SingleLine = true, true

	walletName := win.theme.Editor(new(widget.Editor), "")
	walletName.Editor.SingleLine, walletName.Editor.Submit = true, true

	extendedPublicKey := win.theme.Editor(new(widget.Editor), "Extended public key")
	extendedPublicKey.Editor.Submit = true

	return &ModalTemplate{
		th:                    win.theme,
		confirm:               confirm,
		cancel:                cancel,
		walletName:            walletName,
		oldSpendingPassword:   oldSpendingPassword,
		spendingPassword:      spendingPassword,
		matchSpendingPassword: matchSpendingPassword,
		extendedPublicKey:     extendedPublicKey,
		alert:                 mustIcon(widget.NewIcon(icons.AlertError)),
		passwordStrength:      win.theme.ProgressBar(0),
		keyEvent:              win.keyEvents,
		refreshWindow:         win.refresh,
	}
}

func (m *ModalTemplate) importWatchOnlyWallet() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
		func(gtx C) D {
			return m.extendedPublicKey.Layout(gtx)
		},
	}
}

func (load *modalLoad) setLoading(isLoading bool) {
	load.loading = isLoading
}

func (m *ModalTemplate) createNewWallet() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
		func(gtx C) D {
			return m.spendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			return m.passwordStrength.Layout(gtx)
		},
		func(gtx C) D {
			return m.matchSpendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) renameWallet() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) createNewAccount() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					m.alert.Color = m.th.Color.Gray
					return layout.Inset{Top: values.MarginPadding7, Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return m.alert.Layout(gtx, unit.Dp(15))
					})
				}),
				layout.Rigid(func(gtx C) D {
					info := m.th.Body1("Accounts")
					info.Color = m.th.Color.Gray
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, info.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					info := m.th.Body1(" cannot ")
					info.Color = m.th.Color.DeepBlue
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, info.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					info := m.th.Body1("be deleted when created")
					info.Color = m.th.Color.Gray
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, info.Layout)
				}),
			)
		},
		func(gtx C) D {
			return m.walletName.Layout(gtx)
		},
		func(gtx C) D {
			return m.spendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) removeWallet() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			info := m.th.Body1("Make sure to have the seed phrase backed up before removing the wallet")
			info.Color = m.th.Color.Gray
			return info.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) Password() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return m.spendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) changePassword() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			m.oldSpendingPassword.Editor.SingleLine = true
			return m.oldSpendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			return m.spendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			return m.passwordStrength.Layout(gtx)
		},
		func(gtx C) D {
			return m.matchSpendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) setStartupPassword() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return m.spendingPassword.Layout(gtx)
		},
		func(gtx C) D {
			return m.passwordStrength.Layout(gtx)
		},
		func(gtx C) D {
			return m.matchSpendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) verifyMessageInfo() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := m.th.Body1("After you or your counterparty has genrated a signature, you can use this form to verify the" +
				" validity of the  signature. \n \nOnce you have entered the address, the message and the corresponding " +
				"signature, you will see VALID if the signature appropriately matches the address and message, otherwise INVALID.")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) signMessageInfo() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := m.th.Body1("Signing a message with an address' private key allows you to prove that you are the owner of a given address" +
				" to a possible counterparty.")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) rescanWallet() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := m.th.Body1("Rescanning may help resolve some balance errors. This will take some time, as it scans the entire" +
				" blockchain for transactions")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) securityToolsInfo() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := m.th.Body1("Various tools that help in different aspects of crypto currency security will be located here.")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) sendInfo() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := m.th.Body1("Input or scan the destination wallet address and input the amount to send funds.")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) receiveInfo() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := m.th.Label(values.TextSize20, "Each time you receive a payment, a new address is generated to protect your privacy.")
			text.Color = m.th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) privacyInfo() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := mustIcon(widget.NewIcon(icons.ImageLens))
					ic.Color = m.th.Color.Gray
					return ic.Layout(gtx, values.MarginPadding8)
				}),
				layout.Rigid(func(gtx C) D {
					text := m.th.Body1("When you turn on the mixer, your unmixed DCRs in this wallet (unmixed balance) will be gradually mixed.")
					text.Color = m.th.Color.Gray
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
		func(gtx C) D {
			txt := m.th.Body1("Important: keep this app opened while mixer is running.")
			txt.Font.Weight = text.Bold
			return txt.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := mustIcon(widget.NewIcon(icons.ImageLens))
					ic.Color = m.th.Color.Gray
					return ic.Layout(gtx, values.MarginPadding8)
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

func (m *ModalTemplate) setupMixerInfo() []layout.Widget {
	return []layout.Widget{
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

func (m *ModalTemplate) warnExistMixerAcct() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							m.alert.Color = m.th.Color.DeepBlue
							return m.alert.Layout(gtx, values.MarginPadding50)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					label := m.th.H6("Account name is taken")
					label.Font.Weight = text.Bold
					return label.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			txt := m.th.Body1("There are existing accounts named mixed or unmixed. Please change the name to something else for now. You can change them back after the setup.")
			txt.Color = m.th.Color.Gray
			return txt.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) unlockWalletRestore() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			info := m.th.Body1("The restoration process to discover your accounts was interrupted in the last sync.")
			info.Color = m.th.Color.Gray
			return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return info.Layout(gtx)
			})
		},
		func(gtx C) D {
			info := m.th.Body1("Unlock to resume the process.")
			info.Color = m.th.Color.Gray
			return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return info.Layout(gtx)
			})
		},
		func(gtx C) D {
			return m.spendingPassword.Layout(gtx)
		},
	}
}

func (m *ModalTemplate) transactionDetailsInfo() []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					t := m.th.Body1("Tap on")
					t.Color = m.th.Color.Gray
					return t.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					t := m.th.Body1("blue text")
					t.Color = m.th.Color.Primary
					m := values.MarginPadding2
					return layout.Inset{
						Left:  m,
						Right: m,
					}.Layout(gtx, func(gtx C) D {
						return t.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					t := m.th.Body1("to copy the item.")
					t.Color = m.th.Color.Gray
					return t.Layout(gtx)
				}),
			)
		},
	}
}

func (m *ModalTemplate) Layout(th *decredmaterial.Theme, load *modalLoad) []layout.Widget {
	if !load.isReset {
		m.resetFields()
		load.isReset = true
		load.setLoading(false)
	}

	title := []layout.Widget{
		func(gtx C) D {
			txt := load.title
			if load.template == TransactionDetailsInfoTemplate {
				txt = "How to copy"
			}
			if load.template == PrivacyInfoTemplate {
				txt = "How to use the mixer?"
			}
			t := th.H6(txt)
			t.Font.Weight = text.Bold
			return t.Layout(gtx)
		},
	}

	w := m.handle(th, load)
	w = append(title, w...)
	w = append(w, m.actions(th, load)...)
	return w
}

func (m *ModalTemplate) actions(th *decredmaterial.Theme, load *modalLoad) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if load.cancelText == "" {
							return layout.Dimensions{}
						}
						m.cancel.Text = load.cancelText
						m.cancel.Font.Weight = text.Bold
						m.cancel.Background = th.Color.Surface
						m.cancel.Color = th.Color.Primary
						return m.cancel.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if load.confirmText == "" {
							return layout.Dimensions{}
						}

						m.confirm.Text = load.confirmText
						if load.template == ConfirmRemoveTemplate {
							m.confirm.Background, m.confirm.Color = th.Color.Surface, th.Color.Danger
						}
						if load.template == RescanWalletTemplate {
							m.confirm.Background, m.confirm.Color = th.Color.Surface, th.Color.Primary
						}
						if load.loading {
							th := material.NewTheme(gofont.Collection())
							return layout.Inset{Top: unit.Dp(7)}.Layout(gtx, func(gtx C) D {
								return material.Loader(th).Layout(gtx)
							})
						}
						m.confirm.Font.Weight = text.Bold
						return m.confirm.Layout(gtx)
					}),
				)
			})
		},
	}
}

func (m *ModalTemplate) handle(th *decredmaterial.Theme, load *modalLoad) (template []layout.Widget) {
	m.walletName.Editor.SingleLine = true
	switch load.template {
	case CreateWalletTemplate:
		if m.spendingPassword.Editor.Text() == m.matchSpendingPassword.Editor.Text() {
			// reset error label when password and matching password fields match
			m.matchSpendingPassword.SetError("")
		}

		if m.editorsNotEmpty(m.walletName.Editor, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) &&
			m.passwordsMatch(m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
			if m.confirm.Button.Clicked() || handleSubmitEvent(m.walletName.Editor, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
				load.setLoading(true)
				load.confirm.(func(string, string))(m.walletName.Editor.Text(), m.spendingPassword.Editor.Text())
			}
		}

		if m.cancel.Button.Clicked() {
			load.cancel()
		}

		computePasswordStrength(&m.passwordStrength, th, m.spendingPassword.Editor)

		template = m.createNewWallet()
		m.walletName.Hint = "Wallet name"
		return
	case RenameWalletTemplate, RenameAccountTemplate, ConnectToSpecificPeerTemplate, ChangeSpecificPeerTemplate, UserAgentTemplate:
		if m.editorsNotEmpty(m.walletName.Editor) {
			if m.confirm.Button.Clicked() || handleSubmitEvent(m.walletName.Editor) {
				load.confirm.(func(string))(m.walletName.Editor.Text())
			}
		}
		if m.cancel.Button.Clicked() {
			load.cancel()
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
		if m.editorsNotEmpty(m.walletName.Editor, m.spendingPassword.Editor) {
			if m.confirm.Button.Clicked() || handleSubmitEvent(m.walletName.Editor, m.spendingPassword.Editor) {
				load.setLoading(true)
				load.confirm.(func(string, string))(m.walletName.Editor.Text(), m.spendingPassword.Editor.Text())
			}
		}
		if m.cancel.Button.Clicked() {
			load.cancel()
		}

		template = m.createNewAccount()
		m.walletName.Hint = "Account name"
		return
	case PasswordTemplate, UnlockWalletTemplate, RemoveStartupPasswordTemplate:
		if m.editorsNotEmpty(m.spendingPassword.Editor) {
			if m.confirm.Button.Clicked() || handleSubmitEvent(m.spendingPassword.Editor) {
				load.setLoading(true)
				load.confirm.(func(string))(m.spendingPassword.Editor.Text())
			}
		}
		if m.cancel.Button.Clicked() {
			load.cancel()
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

		if m.editorsNotEmpty(m.oldSpendingPassword.Editor, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) &&
			m.passwordsMatch(m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
			if m.confirm.Button.Clicked() || handleSubmitEvent(m.oldSpendingPassword.Editor, m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
				load.setLoading(true)
				load.confirm.(func(string, string))(m.oldSpendingPassword.Editor.Text(), m.spendingPassword.Editor.Text())
			}
		}

		if m.cancel.Button.Clicked() {
			load.cancel()
		}

		computePasswordStrength(&m.passwordStrength, th, m.spendingPassword.Editor)

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
		if m.editorsNotEmpty(m.walletName.Editor, m.extendedPublicKey.Editor) {
			if m.confirm.Button.Clicked() || handleSubmitEvent(m.walletName.Editor, m.extendedPublicKey.Editor) {
				load.setLoading(true)
				load.confirm.(func(string, string))(m.walletName.Editor.Text(), m.extendedPublicKey.Editor.Text())
			}
		}
		if m.cancel.Button.Clicked() {
			load.cancel()
		}

		m.walletName.Hint = "Wallet name"

		template = m.importWatchOnlyWallet()
		return
	case ConfirmRemoveTemplate:
		m.handleButtonEvents(load)
		template = m.removeWallet()
		return
	case VerifyMessageInfoTemplate:
		m.handleButtonEvents(load)
		template = m.verifyMessageInfo()
		return
	case SignMessageInfoTemplate:
		m.handleButtonEvents(load)
		template = m.signMessageInfo()
		return
	case RescanWalletTemplate:
		m.handleButtonEvents(load)
		template = m.rescanWallet()
		return
	case PrivacyInfoTemplate:
		m.handleButtonEvents(load)
		template = m.privacyInfo()
		return
	case SetStartupPasswordTemplate:
		if m.spendingPassword.Editor.Text() == m.matchSpendingPassword.Editor.Text() {
			// reset error label when password and matching password fields match
			m.matchSpendingPassword.SetError("")
		}

		if m.editorsNotEmpty(m.spendingPassword.Editor, m.matchSpendingPassword.Editor) &&
			m.passwordsMatch(m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
			if m.confirm.Button.Clicked() || handleSubmitEvent(m.spendingPassword.Editor, m.matchSpendingPassword.Editor) {
				load.confirm.(func(string))(m.spendingPassword.Editor.Text())
			}
		}

		if m.cancel.Button.Clicked() {
			load.cancel()
		}

		computePasswordStrength(&m.passwordStrength, th, m.spendingPassword.Editor)
		m.spendingPassword.Hint = "Startup password"
		m.matchSpendingPassword.Hint = "Confirm startup password"

		template = m.setStartupPassword()
		return
	case SetupMixerInfoTemplate:
		m.handleButtonEvents(load)
		template = m.setupMixerInfo()
		return
	case ConfirmMixerAcctExistTemplate:
		m.handleButtonEvents(load)
		template = m.warnExistMixerAcct()
		return
	case UnlockWalletRestoreTemplate:
		if m.editorsNotEmpty(m.spendingPassword.Editor) {
			if m.confirm.Button.Clicked() || handleSubmitEvent(m.spendingPassword.Editor) {
				load.confirm.(func(string))(m.spendingPassword.Editor.Text())
			}
		}

		m.spendingPassword.Hint = "Spending password"

		template = m.unlockWalletRestore()
		return
	case SendInfoTemplate:
		m.handleButtonEvents(load)
		template = m.sendInfo()
		return
	case ReceiveInfoTemplate:
		m.handleButtonEvents(load)
		template = m.receiveInfo()
		return
	case TransactionDetailsInfoTemplate:
		m.handleButtonEvents(load)
		template = m.transactionDetailsInfo()
		return
	default:
		return
	}
}

func (m *ModalTemplate) handleButtonEvents(load *modalLoad) {
	m.handleConfirmEvent()
	if m.confirm.Button.Clicked() || m.confirmKeyPressed {
		m.confirmKeyPressed = false
		load.confirm.(func())()
	}

	if m.cancel.Button.Clicked() {
		load.cancel()
	}
	m.refreshWindow()
}

func (m *ModalTemplate) handleConfirmEvent() {
	select {
	case evt := <-m.keyEvent:
		if (evt.Name == key.NameReturn || evt.Name == key.NameEnter) && evt.State == key.Press {
			m.confirmKeyPressed = true
			return
		}
	default:
	}
}

// editorsNotEmpty checks that the editor fields are not empty. It returns false if they are empty and true if they are
// not and false if it doesn't. It sets the background of the confirm button to decredmaterial Hint color if fields
// are empty. It sets it to decredmaterial Primary color if they are not empty.
func (m *ModalTemplate) editorsNotEmpty(editors ...*widget.Editor) bool {
	m.confirm.Color = m.th.Color.Surface
	for _, e := range editors {
		if e.Text() == "" {
			m.confirm.Background = m.th.Color.Hint
			return false
		}
	}

	m.confirm.Background = m.th.Color.Primary
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
		m.confirm.Background = m.th.Color.Hint
		return false
	}

	m.matchSpendingPassword.SetError("")
	m.confirm.Background = m.th.Color.Primary
	return true
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
