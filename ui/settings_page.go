package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/preference"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageSettings = "Settings"

const (
	USDExchangeValue     = "usd_bittrex"
	DefaultExchangeValue = "none"

	languagePreferenceKey = "app_language"
)

type row struct {
	title     string
	clickable *widget.Clickable
	icon      *widget.Icon
	label     decredmaterial.Label
}

type settingsPage struct {
	common        *pageCommon
	pageContainer layout.List
	theme         *decredmaterial.Theme
	walletInfo    *wallet.MultiWalletInfo
	wal           *wallet.Wallet

	updateConnectToPeer *widget.Clickable
	updateUserAgent     *widget.Clickable
	changeStartupPass   *widget.Clickable
	chevronRightIcon    *widget.Icon
	confirm             decredmaterial.Button
	cancel              decredmaterial.Button

	isDarkModeOn     *widget.Bool
	spendUnconfirmed *widget.Bool
	startupPassword  *widget.Bool
	beepNewBlocks    *widget.Bool
	connectToPeer    *widget.Bool
	userAgent        *widget.Bool

	peerLabel, agentLabel decredmaterial.Label

	isStartupPassword bool
	peerAddr          string
	agentValue        string
	errorReceiver     chan error

	currencyPreference *preference.ListPreference
	languagePreference *preference.ListPreference
}

func SettingsPage(common *pageCommon) Page {
	chevronRightIcon := common.icons.chevronRight

	pg := &settingsPage{
		pageContainer: layout.List{
			Axis: layout.Vertical,
		},
		theme:      common.theme,
		walletInfo: common.info,
		wal:        common.wallet,
		common:     common,

		isDarkModeOn:     new(widget.Bool),
		spendUnconfirmed: new(widget.Bool),
		startupPassword:  new(widget.Bool),
		beepNewBlocks:    new(widget.Bool),
		connectToPeer:    new(widget.Bool),
		userAgent:        new(widget.Bool),
		chevronRightIcon: chevronRightIcon,

		errorReceiver: make(chan error),

		updateConnectToPeer: new(widget.Clickable),
		updateUserAgent:     new(widget.Clickable),
		changeStartupPass:   new(widget.Clickable),

		confirm: common.theme.Button(new(widget.Clickable), "Ok"),
		cancel:  common.theme.Button(new(widget.Clickable), values.String(values.StrCancel)),
	}

	languagePreference := preference.NewListPreference(common.wallet, common.theme, languagePreferenceKey,
		values.DefaultLangauge, values.ArrLanguages).
		Title(values.StrLanguage).
		PostiveButton(values.StrConfirm, func() {
			values.SetUserLanguage(pg.wal.ReadStringConfigValueForKey(languagePreferenceKey))
		}).
		NegativeButton(values.StrCancel, func() {})
	pg.languagePreference = languagePreference

	currencyMap := make(map[string]string)
	currencyMap[DefaultExchangeValue] = values.StrNone
	currencyMap[USDExchangeValue] = values.StrUsdBittrex

	currencyPreference := preference.NewListPreference(common.wallet, common.theme,
		dcrlibwallet.CurrencyConversionConfigKey, DefaultExchangeValue, currencyMap).
		Title(values.StrCurrencyConversion).
		PostiveButton(values.StrConfirm, func() {}).
		NegativeButton(values.StrCancel, func() {})
	pg.currencyPreference = currencyPreference

	color := common.theme.Color.LightGray

	pg.peerLabel = common.theme.Body1("")
	pg.peerLabel.Color = common.theme.Color.Gray

	pg.agentLabel = common.theme.Body1("")
	pg.agentLabel.Color = common.theme.Color.Gray

	pg.chevronRightIcon.Color = color

	return pg
}

func (pg *settingsPage) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.common
	pg.updateSettingOptions()

	body := func(gtx C) D {
		page := SubPage{
			title: values.String(values.StrSettings),
			back: func() {
				common.changePage(PageMore)
			},
			body: func(gtx layout.Context) layout.Dimensions {
				pageContent := []func(gtx C) D{
					pg.general(),
					pg.security(),
					pg.notification(),
					pg.connection(),
				}

				return pg.pageContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
					return layout.Inset{}.Layout(gtx, pageContent[i])
				})
			},
		}
		return common.SubPageLayout(gtx, page)
	}

	if pg.currencyPreference.IsShowing {
		return pg.currencyPreference.Layout(gtx, common.UniformPadding(gtx, body))
	}

	if pg.languagePreference.IsShowing {
		return pg.languagePreference.Layout(gtx, common.UniformPadding(gtx, body))
	}

	return common.UniformPadding(gtx, body)
}

func (pg *settingsPage) general() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, values.String(values.StrGeneral), func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, "Dark mode", pg.isDarkModeOn)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, values.String(values.StrUnconfirmedFunds), pg.spendUnconfirmed)
				}),
				layout.Rigid(pg.lineSeparator()),
				layout.Rigid(func(gtx C) D {
					currencyConversionRow := row{
						title:     values.String(values.StrCurrencyConversion),
						clickable: pg.currencyPreference.Clickable(),
						icon:      pg.chevronRightIcon,
						label:     pg.theme.Body2(pg.wal.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)),
					}
					return pg.clickableRow(gtx, currencyConversionRow)
				}),
				layout.Rigid(pg.lineSeparator()),
				layout.Rigid(func(gtx C) D {
					languageRow := row{
						title:     values.String(values.StrLanguage),
						clickable: pg.languagePreference.Clickable(),
						icon:      pg.chevronRightIcon,
						label:     pg.theme.Body2(pg.wal.ReadStringConfigValueForKey(languagePreferenceKey)),
					}
					return pg.clickableRow(gtx, languageRow)
				}),
			)
		})
	}
}

func (pg *settingsPage) notification() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, values.String(values.StrNotifications), func(gtx C) D {
			return pg.subSectionSwitch(gtx, values.String(values.StrBeepForNewBlocks), pg.beepNewBlocks)
		})
	}
}

func (pg *settingsPage) security() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, values.String(values.StrSecurity), func(gtx C) D {
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
							label:     pg.theme.Body1(""),
						}
						return pg.clickableRow(gtx, changeStartupPassRow)
					})
				}),
			)
		})
	}
}

func (pg *settingsPage) connection() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, values.String(values.StrConnection), func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, values.String(values.StrConnectToSpecificPeer), pg.connectToPeer)
				}),
				layout.Rigid(func(gtx C) D {
					peerAddrRow := row{
						title:     values.String(values.StrChangeSpecificPeer),
						clickable: pg.updateConnectToPeer,
						icon:      pg.chevronRightIcon,
						label:     pg.peerLabel,
					}
					return pg.conditionalDisplay(gtx, pg.peerAddr != "", func(gtx C) D {
						return pg.clickableRow(gtx, peerAddrRow)
					})
				}),
				layout.Rigid(pg.lineSeparator()),
				layout.Rigid(pg.agent()),
			)
		})
	}
}

func (pg *settingsPage) agent() layout.Widget {
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
									txt := pg.theme.Body2(values.String(values.StrUserAgentSummary))
									txt.Color = pg.theme.Color.Gray
									return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
										return txt.Layout(gtx)
									})
								}),
							)
						})
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
							return layout.E.Layout(gtx, pg.theme.Switch(pg.userAgent).Layout)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.conditionalDisplay(gtx, pg.agentValue != "", func(gtx C) D {
					userAgentRow := row{
						title:     values.String(values.StrUserAgentDialogTitle),
						clickable: pg.updateUserAgent,
						icon:      pg.chevronRightIcon,
						label:     pg.agentLabel,
					}
					return pg.clickableRow(gtx, userAgentRow)
				})
			}),
		)
	}
}

func (pg *settingsPage) mainSection(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			m15 := values.MarginPadding15
			return layout.Inset{Top: m15, Left: m15, Right: m15}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						txt := pg.theme.Body2(title)
						txt.Color = pg.theme.Color.Gray
						return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, txt.Layout)
					}),
					layout.Rigid(body),
				)
			})
		})
	})
}

func (pg *settingsPage) subSection(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(pg.subSectionLabel(title)),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, body)
			}),
		)
	})
}

func (pg *settingsPage) subSectionSwitch(gtx layout.Context, title string, option *widget.Bool) layout.Dimensions {
	return pg.subSection(gtx, title, pg.theme.Switch(option).Layout)
}

func (pg *settingsPage) clickableRow(gtx layout.Context, row row) layout.Dimensions {
	return decredmaterial.Clickable(gtx, row.clickable, func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPaddingMinus5}.Layout(gtx, func(gtx C) D {
			return pg.subSection(gtx, row.title, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(row.label.Layout),
					layout.Rigid(func(gtx C) D {
						return row.icon.Layout(gtx, values.MarginPadding22)
					}),
				)
			})
		})
	})
}

func (pg *settingsPage) conditionalDisplay(gtx layout.Context, display bool, body layout.Widget) layout.Dimensions {
	if display {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.lineSeparator()),
			layout.Rigid(body),
		)
	}
	return layout.Dimensions{}
}

func (pg *settingsPage) subSectionLabel(title string) layout.Widget {
	return func(gtx C) D {
		return pg.theme.Body1(title).Layout(gtx)
	}
}

func (pg *settingsPage) lineSeparator() layout.Widget {
	m := values.MarginPadding1
	return func(gtx C) D {
		return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.theme.Separator().Layout)
	}
}

func (pg *settingsPage) handle() {
	common := pg.common
	pg.languagePreference.Handle()
	pg.currencyPreference.Handle()

	if pg.isDarkModeOn.Changed() {
		pg.theme.SwitchDarkMode(pg.isDarkModeOn.Value)
		pg.wal.SaveConfigValueForKey("isDarkModeOn", pg.isDarkModeOn.Value)
	}

	if pg.spendUnconfirmed.Changed() {
		pg.wal.SaveConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey, pg.spendUnconfirmed.Value)
	}

	if pg.beepNewBlocks.Changed() {
		pg.wal.SaveConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey, pg.beepNewBlocks.Value)
	}

	for pg.changeStartupPass.Clicked() {

		newPasswordModal(common).
			title(values.String(values.StrConfirmRemoveStartupPass)).
			hint("Current startup password").
			negativeButton(values.String(values.StrCancel), func() {}).
			positiveButton(values.String(values.StrConfirm), func(password string, pm *passwordModal) bool {
				go func() {
					err := pg.wal.GetMultiWallet().VerifyStartupPassphrase([]byte(password))
					if err != nil {
						pm.setError(err.Error())
						pm.setLoading(false)
						return
					}
					pm.dismiss()

					// change password
					newCreatePasswordModal(common).
						title(values.String(values.StrCreateStartupPassword)).
						enableName(false).
						passwordHint("New startup password").
						confirmPasswordHint("Confirm new startup password").
						passwordCreated(func(walletName, newPassword string, m *createPasswordModal) bool {
							go func() {
								err := pg.wal.GetMultiWallet().ChangeStartupPassphrase([]byte(password), []byte(newPassword), dcrlibwallet.PassphraseTypePass)
								if err != nil {
									m.setError(err.Error())
									m.setLoading(false)
									return
								}
								m.dismiss()
							}()
							return false
						}).show()

				}()

				return false
			}).show()
		break
	}

	if pg.startupPassword.Changed() {
		if pg.startupPassword.Value {
			newCreatePasswordModal(common).
				title(values.String(values.StrCreateStartupPassword)).
				enableName(false).
				passwordHint("Startup password").
				confirmPasswordHint("Confirm startup password").
				passwordCreated(func(walletName, password string, m *createPasswordModal) bool {
					go func() {
						err := pg.wal.GetMultiWallet().SetStartupPassphrase([]byte(password), dcrlibwallet.PassphraseTypePass)
						if err != nil {
							m.setError(err.Error())
							m.setLoading(false)
							return
						}
						m.dismiss()
					}()
					return false
				}).show()
		} else {

			newPasswordModal(common).
				title(values.String(values.StrConfirmRemoveStartupPass)).
				hint("Startup password").
				negativeButton(values.String(values.StrCancel), func() {}).
				positiveButton(values.String(values.StrConfirm), func(password string, pm *passwordModal) bool {
					go func() {
						err := pg.wal.GetMultiWallet().RemoveStartupPassphrase([]byte(password))
						if err != nil {
							pm.setError(err.Error())
							pm.setLoading(false)
							return
						}
						pm.dismiss()
					}()

					return false
				}).show()
		}
	}

	specificPeerKey := dcrlibwallet.SpvPersistentPeerAddressesConfigKey
	if pg.connectToPeer.Changed() {
		if pg.connectToPeer.Value {
			pg.showSPVPeerDialog()
			return
		}
		pg.wal.RemoveUserConfigValueForKey(specificPeerKey)
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
		if pg.userAgent.Value {
			pg.showUserAgentDialog()
			return
		}
		pg.wal.RemoveUserConfigValueForKey(userAgentKey)
	}

	select {
	case err := <-pg.errorReceiver:
		common.modalLoad.setLoading(false)
		if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
			e := "Password is incorrect"
			common.notify(e, false)
			return
		}
		common.notify(err.Error(), false)
	default:
	}
}

func (pg *settingsPage) showSPVPeerDialog() {
	textModal := newTextInputModal(pg.common).
		hint("IP address").
		positiveButton(values.String(values.StrConfirm), func(ipAddress string, tim *textInputModal) bool {
			if ipAddress != "" {
				pg.wal.SaveConfigValueForKey(dcrlibwallet.SpvPersistentPeerAddressesConfigKey, ipAddress)
			}
			return true
		})

	textModal.title(values.String(values.StrConnectToSpecificPeer)).
		negativeButton(values.String(values.StrCancel), func() {})
	textModal.show()
}

func (pg *settingsPage) showUserAgentDialog() {
	textModal := newTextInputModal(pg.common).
		hint("User agent").
		positiveButton(values.String(values.StrConfirm), func(userAgent string, tim *textInputModal) bool {
			if userAgent != "" {
				pg.wal.SaveConfigValueForKey(dcrlibwallet.UserAgentConfigKey, userAgent)
			}
			return true
		})

	textModal.title(values.String(values.StrChangeUserAgent)).
		negativeButton(values.String(values.StrCancel), func() {})
	textModal.show()
}

func (pg *settingsPage) updateSettingOptions() {
	isPassword := pg.wal.IsStartupSecuritySet()
	pg.startupPassword.Value = false
	pg.isStartupPassword = false
	if isPassword {
		pg.startupPassword.Value = true
		pg.isStartupPassword = true
	}

	isDarkModeOn := pg.wal.ReadBoolConfigValueForKey("isDarkModeOn")
	pg.isDarkModeOn.Value = false
	if isDarkModeOn {
		pg.isDarkModeOn.Value = true
	}

	isSpendUnconfirmed := pg.wal.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey)
	pg.spendUnconfirmed.Value = false
	if isSpendUnconfirmed {
		pg.spendUnconfirmed.Value = true
	}

	beep := pg.wal.ReadBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey)
	pg.beepNewBlocks.Value = false
	if beep {
		pg.beepNewBlocks.Value = true
	}

	pg.peerAddr = pg.wal.ReadStringConfigValueForKey(dcrlibwallet.SpvPersistentPeerAddressesConfigKey)
	pg.connectToPeer.Value = false
	if pg.peerAddr != "" {
		pg.peerLabel.Text = pg.peerAddr
		pg.connectToPeer.Value = true
	}

	pg.agentValue = pg.wal.ReadStringConfigValueForKey(dcrlibwallet.UserAgentConfigKey)
	pg.userAgent.Value = false
	if pg.agentValue != "" {
		pg.agentLabel.Text = pg.agentValue
		pg.userAgent.Value = true
	}
}

func (pg *settingsPage) onClose() {}
