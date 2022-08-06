package page

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/preference"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const SettingsPageID = "Settings"

type row struct {
	title     string
	clickable *decredmaterial.Clickable
	icon      *decredmaterial.Icon
	label     decredmaterial.Label
}

type SettingsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	pageContainer *widget.List
	wal           *wallet.Wallet

	updateConnectToPeer *decredmaterial.Clickable
	updateUserAgent     *decredmaterial.Clickable
	changeStartupPass   *decredmaterial.Clickable
	language            *decredmaterial.Clickable
	currency            *decredmaterial.Clickable
	help                *decredmaterial.Clickable
	about               *decredmaterial.Clickable

	chevronRightIcon *decredmaterial.Icon
	backButton       decredmaterial.IconButton
	infoButton       decredmaterial.IconButton

	isDarkModeOn            *decredmaterial.Switch
	spendUnconfirmed        *decredmaterial.Switch
	startupPassword         *decredmaterial.Switch
	beepNewBlocks           *decredmaterial.Switch
	connectToPeer           *decredmaterial.Switch
	userAgent               *decredmaterial.Switch
	governance              *decredmaterial.Switch
	proposalNotification    *decredmaterial.Switch
	transactionNotification *decredmaterial.Switch

	peerLabel, agentLabel decredmaterial.Label

	isStartupPassword bool
	peerAddr          string
	agentValue        string
	errorReceiver     chan error
}

func NewSettingsPage(l *load.Load) *SettingsPage {
	chevronRightIcon := l.Theme.Icons.ChevronRight

	pg := &SettingsPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(SettingsPageID),
		pageContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		wal: l.WL.Wallet,

		isDarkModeOn:            l.Theme.Switch(),
		spendUnconfirmed:        l.Theme.Switch(),
		startupPassword:         l.Theme.Switch(),
		beepNewBlocks:           l.Theme.Switch(),
		connectToPeer:           l.Theme.Switch(),
		userAgent:               l.Theme.Switch(),
		governance:              l.Theme.Switch(),
		proposalNotification:    l.Theme.Switch(),
		transactionNotification: l.Theme.Switch(),

		chevronRightIcon: decredmaterial.NewIcon(chevronRightIcon),

		errorReceiver: make(chan error),

		updateConnectToPeer: l.Theme.NewClickable(false),
		updateUserAgent:     l.Theme.NewClickable(false),
		changeStartupPass:   l.Theme.NewClickable(false),
		language:            l.Theme.NewClickable(false),
		currency:            l.Theme.NewClickable(false),
		help:                l.Theme.NewClickable(false),
		about:               l.Theme.NewClickable(false),
	}

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SettingsPage) OnNavigatedTo() {

}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *SettingsPage) Layout(gtx C) D {
	pg.updateSettingOptions()
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *SettingsPage) layoutDesktop(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.pageHeaderLayout),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, pg.pageContentLayout)
			}),
		)
	})
}

func (pg *SettingsPage) pageHeaderLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.W.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding16,
							Top:   values.MarginPaddingMinus2,
						}.Layout(gtx, pg.backButton.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						return pg.Theme.Label(values.TextSize20, values.String(values.StrSettings)).Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *SettingsPage) pageContentLayout(gtx layout.Context) layout.Dimensions {
	pageContent := []func(gtx C) D{
		pg.general(),
		pg.security(),
		pg.info(),
		// pg.notification(),
		// pg.connection(),
	}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Dp(values.MarginPadding550)
		gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding550)
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		return pg.Theme.List(pg.pageContainer).Layout(gtx, len(pageContent), func(gtx C, i int) D {
			return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, pageContent[i])
		})
	})
}

func (pg *SettingsPage) layoutMobile(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}

func (pg *SettingsPage) settingLine(gtx C) D {
	line := pg.Theme.Line(1, 0)
	line.Color = pg.Theme.Color.Gray3
	return line.Layout(gtx)
}

func (pg *SettingsPage) wrapSection(gtx C, title string, body layout.Widget) D {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := pg.Theme.Body2(title)
							txt.Color = pg.Theme.Color.GrayText2
							return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, txt.Layout)
						}),
						layout.Flexed(1, func(gtx C) D {
							if title == values.String(values.StrSecurity) {
								pg.infoButton.Inset = layout.UniformInset(values.MarginPadding0)
								pg.infoButton.Size = values.MarginPadding20
								return layout.E.Layout(gtx, pg.infoButton.Layout)
							}
							return D{}
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, pg.settingLine)
				}),
				layout.Rigid(body),
			)
		})
	})
}

func (pg *SettingsPage) general() layout.Widget {
	return func(gtx C) D {
		return pg.wrapSection(gtx, values.String(values.StrGeneral), func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					exchangeRate := row{
						title:     values.String(values.StrExchangeRate),
						clickable: pg.currency,
						icon:      pg.chevronRightIcon,
						label:     pg.Theme.Body2(pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)),
					}
					return pg.clickableRow(gtx, exchangeRate)
				}),
				layout.Rigid(func(gtx C) D {
					languageRow := row{
						title:     values.String(values.StrLanguage),
						clickable: pg.language,
						icon:      pg.chevronRightIcon,
						label:     pg.Theme.Body2(pg.WL.MultiWallet.ReadStringConfigValueForKey(load.LanguagePreferenceKey)),
					}
					return pg.clickableRow(gtx, languageRow)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, values.StringF(values.StrTxNotification, ""), pg.transactionNotification)
				}),
			)
		})
	}
}

func (pg *SettingsPage) security() layout.Widget {
	return func(gtx C) D {
		return pg.wrapSection(gtx, values.String(values.StrSecurity), func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, values.String(values.StrStartupPassword), pg.startupPassword)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.conditionalDisplay(gtx, pg.isStartupPassword, func(gtx C) D {
						changeStartupPassRow := row{
							title:     values.String(values.StrChangeStartupPassword),
							clickable: pg.changeStartupPass,
							icon:      pg.chevronRightIcon,
							label:     pg.Theme.Body1(""),
						}
						return pg.clickableRow(gtx, changeStartupPassRow)
					})
				}),
			)
		})
	}
}

func (pg *SettingsPage) info() layout.Widget {
	return func(gtx C) D {
		return pg.wrapSection(gtx, values.String(values.StrInfo), func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					helpRow := row{
						title:     values.String(values.StrHelp),
						clickable: pg.help,
						icon:      pg.chevronRightIcon,
						label:     pg.Theme.Body2(""),
					}
					return pg.clickableRow(gtx, helpRow)
				}),
				layout.Rigid(func(gtx C) D {
					aboutRow := row{
						title:     values.String(values.StrAbout),
						clickable: pg.about,
						icon:      pg.chevronRightIcon,
						label:     pg.Theme.Body2(""),
					}
					return pg.clickableRow(gtx, aboutRow)
				}),
			)
		})
	}
}

func (pg *SettingsPage) agent() layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						m10 := values.MarginPadding10
						return layout.Inset{Top: m10, Bottom: m10}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(pg.subSectionLabel(values.String(values.StrCustomUserAgent))),
								layout.Rigid(func(gtx C) D {
									txt := pg.Theme.Body2(values.String(values.StrHTTPRequest))
									txt.Color = pg.Theme.Color.GrayText2
									return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
										return txt.Layout(gtx)
									})
								}),
							)
						})
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
							return layout.E.Layout(gtx, pg.userAgent.Layout)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				agentLabel := pg.Theme.Body1(pg.agentValue)
				agentLabel.Color = pg.Theme.Color.GrayText2
				return pg.conditionalDisplay(gtx, pg.agentValue != "", func(gtx C) D {
					userAgentRow := row{
						title:     values.String(values.StrUserAgentDialogTitle),
						clickable: pg.updateUserAgent,
						icon:      pg.chevronRightIcon,
						label:     agentLabel,
					}
					return pg.clickableRow(gtx, userAgentRow)
				})
			}),
		)
	}
}

func (pg *SettingsPage) mainSection(gtx C, title string, body layout.Widget) D {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := pg.Theme.Body2(title)
								txt.Color = pg.Theme.Color.GrayText2
								return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, txt.Layout)
							}),
							layout.Flexed(1, func(gtx C) D {
								if title == values.String(values.StrSecurity) {
									pg.infoButton.Inset = layout.UniformInset(values.MarginPadding0)
									pg.infoButton.Size = values.MarginPadding20
									return layout.E.Layout(gtx, pg.infoButton.Layout)
								}
								return D{}
							}),
						)
					}),
					layout.Rigid(body),
				)
			})
		})
	})
}

func (pg *SettingsPage) subSection(gtx C, title string, body layout.Widget) D {
	return layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(pg.subSectionLabel(title)),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, body)
			}),
		)
	})
}

func (pg *SettingsPage) subSectionSwitch(gtx C, title string, option *decredmaterial.Switch) D {
	return pg.subSection(gtx, title, option.Layout)
}

func (pg *SettingsPage) clickableRow(gtx C, row row) D {
	return row.clickable.Layout(gtx, func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPaddingMinus5}.Layout(gtx, func(gtx C) D {
			return pg.subSection(gtx, row.title, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(row.label.Layout),
					layout.Rigid(func(gtx C) D {
						ic := row.icon
						ic.Color = pg.Theme.Color.Gray1
						return ic.Layout(gtx, values.MarginPadding22)
					}),
				)
			})
		})
	})
}

func (pg *SettingsPage) conditionalDisplay(gtx C, display bool, body layout.Widget) D {
	if display {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.lineSeparator()),
			layout.Rigid(body),
		)
	}
	return D{}
}

func (pg *SettingsPage) subSectionLabel(title string) layout.Widget {
	return func(gtx C) D {
		return pg.Theme.Body1(title).Layout(gtx)
	}
}

func (pg *SettingsPage) lineSeparator() layout.Widget {
	m := values.MarginPadding1
	return func(gtx C) D {
		return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.Theme.Separator().Layout)
	}
}

func (pg *SettingsPage) showWarningModalDialog(title, msg, key string) {
	info := modal.NewInfoModal(pg.Load).
		Title(title).
		Body(msg).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButtonStyle(pg.Theme.Color.Surface, pg.Theme.Color.Danger).
		PositiveButton(values.String(values.StrRemove), func(isChecked bool) bool {
			pg.WL.MultiWallet.DeleteUserConfigValueForKey(key)
			return true
		})
	pg.ParentWindow().ShowModal(info)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SettingsPage) HandleUserInteractions() {

	for pg.language.Clicked() {
		langSelectorModal := preference.NewListPreference(pg.Load,
			load.LanguagePreferenceKey, values.DefaultLangauge, values.ArrLanguages).
			Title(values.StrLanguage).
			UpdateValues(func() {
				values.SetUserLanguage(pg.WL.MultiWallet.ReadStringConfigValueForKey(load.LanguagePreferenceKey))
			})
		pg.ParentWindow().ShowModal(langSelectorModal)
		break
	}

	for pg.backButton.Button.Clicked() {
		pg.ParentNavigator().CloseCurrentPage()
	}

	for pg.currency.Clicked() {
		currencySelectorModal := preference.NewListPreference(pg.Load,
			dcrlibwallet.CurrencyConversionConfigKey, values.DefaultExchangeValue,
			values.ArrExchangeCurrencies).
			Title(values.StrExchangeRate).
			UpdateValues(func() {})
		pg.ParentWindow().ShowModal(currencySelectorModal)
		break
	}

	if pg.isDarkModeOn.Changed() {
		pg.WL.MultiWallet.SaveUserConfigValue(load.DarkModeConfigKey, pg.isDarkModeOn.IsChecked())
		pg.RefreshTheme(pg.ParentWindow())
	}

	if pg.spendUnconfirmed.Changed() {
		pg.WL.MultiWallet.SaveUserConfigValue(dcrlibwallet.SpendUnconfirmedConfigKey, pg.spendUnconfirmed.IsChecked())
	}

	if pg.governance.Changed() {
		if pg.governance.IsChecked() {
			go pg.WL.MultiWallet.Politeia.Sync()
			pg.WL.MultiWallet.SaveUserConfigValue(load.FetchProposalConfigKey, pg.governance.IsChecked())
			pg.Toast.Notify(values.StringF(values.StrPropFetching, values.String(values.StrEnabled), values.String(values.StrCheckGovernace)))
		} else {
			info := modal.NewInfoModal(pg.Load).
				Title(values.String(values.StrGovernance)).
				Body(values.String(values.StrGovernanceSettingsInfo)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButtonStyle(pg.Theme.Color.Surface, pg.Theme.Color.Danger).
				PositiveButton(values.String(values.StrDisable), func(isChecked bool) bool {
					if pg.WL.MultiWallet.Politeia.IsSyncing() {
						go pg.WL.MultiWallet.Politeia.StopSync()
					}
					pg.WL.MultiWallet.SaveUserConfigValue(load.FetchProposalConfigKey, !pg.governance.IsChecked())
					pg.WL.MultiWallet.Politeia.ClearSavedProposals()
					pg.Toast.Notify(values.StringF(values.StrPropFetching, values.String(values.StrDisabled)))
					return true
				})
			pg.ParentWindow().ShowModal(info)
		}
	}

	if pg.beepNewBlocks.Changed() {
		pg.WL.MultiWallet.SaveUserConfigValue(dcrlibwallet.BeepNewBlocksConfigKey, pg.beepNewBlocks.IsChecked())
	}

	if pg.proposalNotification.Changed() {
		pg.WL.MultiWallet.SaveUserConfigValue(load.ProposalNotificationConfigKey, pg.proposalNotification.IsChecked())
		if pg.proposalNotification.IsChecked() {
			pg.Toast.Notify(values.StringF(values.StrPropNotification, values.String(values.StrEnabled)))
		} else {
			pg.Toast.Notify(values.StringF(values.StrPropNotification, values.String(values.StrDisabled)))
		}
	}

	if pg.transactionNotification.Changed() {
		pg.WL.MultiWallet.SaveUserConfigValue(load.TransactionNotificationConfigKey, pg.transactionNotification.IsChecked())
		if pg.transactionNotification.IsChecked() {
			pg.Toast.Notify(values.StringF(values.StrTxNotification, values.String(values.StrEnabled)))
		} else {
			pg.Toast.Notify(values.StringF(values.StrTxNotification, values.String(values.StrDisabled)))
		}
	}

	if pg.infoButton.Button.Clicked() {
		info := modal.NewInfoModal(pg.Load).
			Title(values.String(values.StrSetupStartupPassword)).
			Body(values.String(values.StrStartupPasswordInfo)).
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) bool {
				return true
			})
		pg.ParentWindow().ShowModal(info)
	}

	for pg.changeStartupPass.Clicked() {
		currentPasswordModal := modal.NewPasswordModal(pg.Load).
			Title(values.String(values.StrConfirmStartupPass)).
			Hint(values.String(values.StrCurrentStartupPass)).
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
				go func() {
					var error string
					err := pg.wal.GetMultiWallet().VerifyStartupPassphrase([]byte(password))
					if err != nil {
						if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
							error = values.String(values.StrInvalidPassphrase)
						} else {
							error = err.Error()
						}
						pm.SetError(error)
						pm.SetLoading(false)
						return
					}
					pm.Dismiss()

					// change password
					newPasswordModal := modal.NewCreatePasswordModal(pg.Load).
						Title(values.String(values.StrCreateStartupPassword)).
						EnableName(false).
						PasswordHint(values.String(values.StrNewStartupPass)).
						ConfirmPasswordHint(values.String(values.StrConfirmNewStartupPass)).
						PasswordCreated(func(walletName, newPassword string, m *modal.CreatePasswordModal) bool {
							go func() {
								err := pg.wal.GetMultiWallet().ChangeStartupPassphrase([]byte(password), []byte(newPassword), dcrlibwallet.PassphraseTypePass)
								if err != nil {
									m.SetError(err.Error())
									m.SetLoading(false)
									return
								}
								pg.Toast.Notify(values.String(values.StrStartupPassConfirm))
								m.Dismiss()
							}()
							return false
						})
					pg.ParentWindow().ShowModal(newPasswordModal)
				}()

				return false
			})
		pg.ParentWindow().ShowModal(currentPasswordModal)
		break
	}

	if pg.startupPassword.Changed() {
		if pg.startupPassword.IsChecked() {
			createPasswordModal := modal.NewCreatePasswordModal(pg.Load).
				Title(values.String(values.StrCreateStartupPassword)).
				EnableName(false).
				PasswordHint(values.String(values.StrStartupPassword)).
				ConfirmPasswordHint(values.String(values.StrConfirmStartupPass)).
				NegativeButton(func() {}).
				PasswordCreated(func(walletName, password string, m *modal.CreatePasswordModal) bool {
					go func() {
						err := pg.wal.GetMultiWallet().SetStartupPassphrase([]byte(password), dcrlibwallet.PassphraseTypePass)
						if err != nil {
							m.SetError(err.Error())
							m.SetLoading(false)
							return
						}
						pg.Toast.Notify(values.StringF(values.StrStartupPasswordEnabled, values.String(values.StrEnabled)))
						m.Dismiss()
					}()
					return false
				})
			pg.ParentWindow().ShowModal(createPasswordModal)
		} else {
			currentPasswordModal := modal.NewPasswordModal(pg.Load).
				Title(values.String(values.StrConfirmRemoveStartupPass)).
				Hint(values.String(values.StrStartupPassword)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
					go func() {
						var error string
						err := pg.wal.GetMultiWallet().RemoveStartupPassphrase([]byte(password))
						if err != nil {
							if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
								error = values.String(values.StrInvalidPassphrase)
							} else {
								error = err.Error()
							}
							pm.SetError(error)
							pm.SetLoading(false)
							return
						}
						pg.Toast.Notify(values.StringF(values.StrStartupPasswordEnabled, values.String(values.StrDisabled)))
						pm.Dismiss()
					}()

					return false
				})
			pg.ParentWindow().ShowModal(currentPasswordModal)
		}
	}

	specificPeerKey := dcrlibwallet.SpvPersistentPeerAddressesConfigKey
	if pg.connectToPeer.Changed() {
		if pg.connectToPeer.IsChecked() {
			pg.showSPVPeerDialog()
			return
		}

		title := values.String(values.StrRemovePeer)
		msg := values.String(values.StrRemovePeerWarn)
		pg.showWarningModalDialog(title, msg, specificPeerKey)
	}

	for pg.updateConnectToPeer.Clicked() {
		pg.showSPVPeerDialog()
		break
	}

	userAgentKey := dcrlibwallet.UserAgentConfigKey
	for pg.updateUserAgent.Clicked() {
		pg.showUserAgentDialog()
		break
	}

	if pg.userAgent.Changed() {
		if pg.userAgent.IsChecked() {
			pg.showUserAgentDialog()
			return
		}

		title := values.String(values.StrRemoveUserAgent)
		msg := values.String(values.StrRemoveUserAgentWarn)
		pg.showWarningModalDialog(title, msg, userAgentKey)
	}

	select {
	case err := <-pg.errorReceiver:
		if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
			pg.Toast.NotifyError(values.String(values.StrInvalidPassphrase))
			return
		}
		pg.Toast.NotifyError(err.Error())
	default:
	}
}

func (pg *SettingsPage) showSPVPeerDialog() {
	textModal := modal.NewTextInputModal(pg.Load).
		Hint(values.String(values.StrIPAddress)).
		PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
		PositiveButton(values.String(values.StrConfirm), func(ipAddress string, tim *modal.TextInputModal) bool {
			if ipAddress != "" {
				pg.WL.MultiWallet.SaveUserConfigValue(dcrlibwallet.SpvPersistentPeerAddressesConfigKey, ipAddress)
			}
			return true
		})

	textModal.Title(values.String(values.StrConnectToSpecificPeer)).
		NegativeButton(values.String(values.StrCancel), func() {})
	pg.ParentWindow().ShowModal(textModal)
}

func (pg *SettingsPage) showUserAgentDialog() {
	textModal := modal.NewTextInputModal(pg.Load).
		Hint(values.String(values.StrUserAgent)).
		PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
		PositiveButton(values.String(values.StrConfirm), func(userAgent string, tim *modal.TextInputModal) bool {
			if userAgent != "" {
				pg.WL.MultiWallet.SaveUserConfigValue(dcrlibwallet.UserAgentConfigKey, userAgent)
			}
			return true
		})

	textModal.Title(values.String(values.StrChangeUserAgent)).
		NegativeButton(values.String(values.StrCancel), func() {})
	pg.ParentWindow().ShowModal(textModal)
}

func (pg *SettingsPage) updateSettingOptions() {
	isPassword := pg.WL.MultiWallet.IsStartupSecuritySet()
	pg.startupPassword.SetChecked(false)
	pg.isStartupPassword = false
	if isPassword {
		pg.startupPassword.SetChecked(isPassword)
		pg.isStartupPassword = true
	}

	isDarkModeOn := pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.DarkModeConfigKey, false)
	pg.isDarkModeOn.SetChecked(false)
	if isDarkModeOn {
		pg.isDarkModeOn.SetChecked(isDarkModeOn)
	}

	isSpendUnconfirmed := pg.WL.MultiWallet.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey, false)
	pg.spendUnconfirmed.SetChecked(false)
	if isSpendUnconfirmed {
		pg.spendUnconfirmed.SetChecked(isSpendUnconfirmed)
	}

	beep := pg.WL.MultiWallet.ReadBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey, false)
	pg.beepNewBlocks.SetChecked(false)
	if beep {
		pg.beepNewBlocks.SetChecked(beep)
	}

	pg.peerAddr = pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.SpvPersistentPeerAddressesConfigKey)
	pg.connectToPeer.SetChecked(false)
	if pg.peerAddr != "" {
		pg.peerLabel.Text = pg.peerAddr
		pg.connectToPeer.SetChecked(true)
	}

	pg.agentValue = pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.UserAgentConfigKey)
	pg.userAgent.SetChecked(false)
	if pg.agentValue != "" {
		pg.agentLabel.Text = pg.agentValue
		pg.userAgent.SetChecked(true)
	}

	governanceSet := pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey, false)
	pg.governance.SetChecked(false)
	if governanceSet {
		pg.governance.SetChecked(governanceSet)
	}

	proposalNotification := pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.ProposalNotificationConfigKey, false)
	pg.proposalNotification.SetChecked(false)
	if proposalNotification {
		pg.proposalNotification.SetChecked(proposalNotification)
	}

	transactionNotification := pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.TransactionNotificationConfigKey, false)
	pg.transactionNotification.SetChecked(false)
	if transactionNotification {
		pg.transactionNotification.SetChecked(transactionNotification)
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *SettingsPage) OnNavigatedFrom() {}
