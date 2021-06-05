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
	CreateAccountTemplate          = "CreateNewAccount"
	ImportWatchOnlyWalletTemplate  = "ImportWatchOnlyWallet"

	VerifyMessageInfoTemplate      = "VerifyMessageInfo"
	SignMessageInfoTemplate        = "SignMessageInfo"
	PrivacyInfoTemplate            = "PrivacyInfo"
	SetupMixerInfoTemplate         = "ConfirmSetupMixer"
	SecurityToolsInfoTemplate      = "SecurityToolsInfo"
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

func (common *pageCommon) LoadModalTemplates() *ModalTemplate {
	cancel := common.theme.Button(new(widget.Clickable), "Cancel")
	confirm := common.theme.Button(new(widget.Clickable), "Confirm")
	cancel.TextSize, confirm.TextSize = values.TextSize16, values.TextSize16

	spendingPassword := common.theme.EditorPassword(new(widget.Editor), "Spending password")
	spendingPassword.Editor.SingleLine, spendingPassword.Editor.Submit = true, true

	matchSpendingPassword := common.theme.EditorPassword(new(widget.Editor), "Confirm spending password")
	matchSpendingPassword.Editor.Submit, matchSpendingPassword.Editor.SingleLine = true, true

	oldSpendingPassword := common.theme.EditorPassword(new(widget.Editor), "Old spending password")
	oldSpendingPassword.Editor.Submit, oldSpendingPassword.Editor.SingleLine = true, true

	walletName := common.theme.Editor(new(widget.Editor), "")
	walletName.Editor.SingleLine, walletName.Editor.Submit = true, true

	extendedPublicKey := common.theme.Editor(new(widget.Editor), "Extended public key")
	extendedPublicKey.Editor.Submit = true

	return &ModalTemplate{
		th:                    common.theme,
		confirm:               confirm,
		cancel:                cancel,
		walletName:            walletName,
		oldSpendingPassword:   oldSpendingPassword,
		spendingPassword:      spendingPassword,
		matchSpendingPassword: matchSpendingPassword,
		extendedPublicKey:     extendedPublicKey,
		alert:                 mustIcon(widget.NewIcon(icons.AlertError)),
		passwordStrength:      common.theme.ProgressBar(0),
		keyEvent:              common.keyEvents,
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

func verifyMessageInfo(th *decredmaterial.Theme) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := th.Body1("After you or your counterparty has genrated a signature, you can use this form to verify the" +
				" validity of the  signature. \n \nOnce you have entered the address, the message and the corresponding " +
				"signature, you will see VALID if the signature appropriately matches the address and message, otherwise INVALID.")
			text.Color = th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func signMessageInfo(th *decredmaterial.Theme) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			text := th.Body1("Signing a message with an address' private key allows you to prove that you are the owner of a given address" +
				" to a possible counterparty.")
			text.Color = th.Color.Gray
			return text.Layout(gtx)
		},
	}
}

func privacyInfo(th *decredmaterial.Theme) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := mustIcon(widget.NewIcon(icons.ImageLens))
					ic.Color = th.Color.Gray
					return ic.Layout(gtx, values.MarginPadding8)
				}),
				layout.Rigid(func(gtx C) D {
					text := th.Body1("When you turn on the mixer, your unmixed DCRs in this wallet (unmixed balance) will be gradually mixed.")
					text.Color = th.Color.Gray
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
		func(gtx C) D {
			txt := th.Body1("Important: keep this app opened while mixer is running.")
			txt.Font.Weight = text.Bold
			return txt.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					ic := mustIcon(widget.NewIcon(icons.ImageLens))
					ic.Color = th.Color.Gray
					return ic.Layout(gtx, values.MarginPadding8)
				}),
				layout.Rigid(func(gtx C) D {
					text := th.Body1("Mixer will automatically stop when unmixed balance are fully mixed.")
					text.Color = th.Color.Gray
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, text.Layout)
				}),
			)
		},
	}
}

func setupMixerInfo(th *decredmaterial.Theme) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			txt := th.Body1("Two dedicated accounts (“mixed” & “unmixed”) will be created in order to use the mixer.")
			txt.Color = th.Color.Gray
			return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, txt.Layout)
		},
		func(gtx C) D {
			txt := th.Label(values.TextSize18, "This action cannot be undone.")
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

func transactionDetailsInfo(th *decredmaterial.Theme) []layout.Widget {
	return []layout.Widget{
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					t := th.Body1("Tap on")
					t.Color = th.Color.Gray
					return t.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					t := th.Body1("blue text")
					t.Color = th.Color.Primary
					m := values.MarginPadding2
					return layout.Inset{
						Left:  m,
						Right: m,
					}.Layout(gtx, func(gtx C) D {
						return t.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					t := th.Body1("to copy the item.")
					t.Color = th.Color.Gray
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
						// if load.template == ConfirmRemoveTemplate { //TODO
						// 	m.confirm.Background, m.confirm.Color = th.Color.Surface, th.Color.Danger
						// }
						// if load.template == RescanWalletTemplate {
						// 	m.confirm.Background, m.confirm.Color = th.Color.Surface, th.Color.Primary
						// }
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
	case VerifyMessageInfoTemplate:
		m.handleButtonEvents(load)
		template = verifyMessageInfo(m.th)
		return
	case SignMessageInfoTemplate:
		m.handleButtonEvents(load)
		template = signMessageInfo(m.th)
		return
	case PrivacyInfoTemplate:
		m.handleButtonEvents(load)
		template = privacyInfo(m.th)
		return
	case TransactionDetailsInfoTemplate:
		m.handleButtonEvents(load)
		template = transactionDetailsInfo(m.th)
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
	// m.refreshWindow()
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
