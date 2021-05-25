package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageSettings = "Settings"
const USDExchangeValue = "USD (Bittrex)"

type row struct {
	title     string
	clickable *widget.Clickable
	icon      *widget.Icon
	label     decredmaterial.Label
}

type settingsPage struct {
	pageContainer layout.List
	theme         *decredmaterial.Theme
	walletInfo    *wallet.MultiWalletInfo
	wal           *wallet.Wallet

	currencyConversion  *widget.Clickable
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

	currencyValue       string
	initialValue        string
	isCurrencyModalOpen bool

	currencyModal     *decredmaterial.Modal
	radioButtonsGroup *widget.Enum
}

func (win *Window) SettingsPage(common pageCommon) layout.Widget {
	chevronRightIcon := common.icons.chevronRight

	pg := &settingsPage{
		pageContainer: layout.List{
			Axis: layout.Vertical,
		},
		theme:      common.theme,
		walletInfo: win.walletInfo,
		wal:        common.wallet,

		isDarkModeOn:     new(widget.Bool),
		spendUnconfirmed: new(widget.Bool),
		startupPassword:  new(widget.Bool),
		beepNewBlocks:    new(widget.Bool),
		connectToPeer:    new(widget.Bool),
		userAgent:        new(widget.Bool),
		chevronRightIcon: chevronRightIcon,

		errorReceiver: make(chan error),

		currencyConversion:  new(widget.Clickable),
		updateConnectToPeer: new(widget.Clickable),
		updateUserAgent:     new(widget.Clickable),
		changeStartupPass:   new(widget.Clickable),

		confirm: win.theme.Button(new(widget.Clickable), "Ok"),
		cancel:  win.theme.Button(new(widget.Clickable), "Cancel"),

		currencyModal:     common.theme.Modal(),
		radioButtonsGroup: new(widget.Enum),
	}

	color := common.theme.Color.LightGray

	pg.peerLabel = common.theme.Body1("")
	pg.peerLabel.Color = common.theme.Color.Gray

	pg.agentLabel = common.theme.Body1("")
	pg.agentLabel.Color = common.theme.Color.Gray

	pg.chevronRightIcon.Color = color

	return func(gtx C) D {
		pg.handle(common, win)
		return pg.Layout(gtx, common)
	}
}

func (pg *settingsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.updateSettingOptions()

	body := func(gtx C) D {
		page := SubPage{
			title: "Settings",
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

	if pg.isCurrencyModalOpen {
		return common.Modal(gtx, common.Layout(gtx, body), pg.currencyConversionSection(gtx))
	}

	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, body)
	})
}

func (pg *settingsPage) general() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, "General", func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, "Dark mode", pg.isDarkModeOn)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, "Spending unconfirmed funds", pg.spendUnconfirmed)
				}),
				layout.Rigid(pg.lineSeparator()),
				layout.Rigid(func(gtx C) D {
					currencyConversionRow := row{
						title:     "Currency conversion",
						clickable: pg.currencyConversion,
						icon:      pg.chevronRightIcon,
						label:     pg.theme.Body2(pg.currencyValue),
					}
					return pg.clickableRow(gtx, currencyConversionRow)
				}),
			)
		})
	}
}

func (pg *settingsPage) notification() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, "Notification", func(gtx C) D {
			return pg.subSectionSwitch(gtx, "Beep for new blocks", pg.beepNewBlocks)
		})
	}
}

func (pg *settingsPage) security() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, "Security", func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, "Startup password", pg.startupPassword)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.conditionalDisplay(gtx, pg.isStartupPassword, func(gtx C) D {
						changeStartupPassRow := row{
							title:     "Change startup password",
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
		return pg.mainSection(gtx, "Connection", func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSectionSwitch(gtx, "Connect to specific peer", pg.connectToPeer)
				}),
				layout.Rigid(func(gtx C) D {
					peerAddrRow := row{
						title:     "Change specific peer",
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
								layout.Rigid(pg.subSectionLabel("Custom user agent")),
								layout.Rigid(func(gtx C) D {
									txt := pg.theme.Body2("For exchange rate fetching")
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
							return layout.E.Layout(gtx, func(gtx C) D {
								return pg.theme.Switch(pg.userAgent).Layout(gtx)
							})
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.conditionalDisplay(gtx, pg.agentValue != "", func(gtx C) D {
					userAgentRow := row{
						title:     "Change user agent",
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

func (pg *settingsPage) currencyConversionSection(gtx layout.Context) layout.Dimensions {
	w := []layout.Widget{
		func(gtx C) D {
			txt := pg.theme.H6("Currency conversion")
			txt.Color = pg.theme.Color.Text
			return txt.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.theme.RadioButton(pg.radioButtonsGroup, "None", "None").Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return pg.theme.RadioButton(pg.radioButtonsGroup, USDExchangeValue, USDExchangeValue).Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
							pg.cancel.Background, pg.cancel.Color = color.NRGBA{}, pg.theme.Color.Primary
							return pg.cancel.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
							pg.confirm.Background, pg.confirm.Color = color.NRGBA{}, pg.theme.Color.Primary
							return pg.confirm.Layout(gtx)
						})
					}),
				)
			})
		},
	}

	return pg.currencyModal.Layout(gtx, w, 1050)
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
						return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return txt.Layout(gtx)
						})
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
	return pg.subSection(gtx, title, func(gtx C) D {
		return pg.theme.Switch(option).Layout(gtx)
	})
}

func (pg *settingsPage) clickableRow(gtx layout.Context, row row) layout.Dimensions {
	return decredmaterial.Clickable(gtx, row.clickable, func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPaddingMinus5}.Layout(gtx, func(gtx C) D {
			return pg.subSection(gtx, row.title, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return row.label.Layout(gtx)
					}),
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
		return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
			return pg.theme.Separator().Layout(gtx)
		})
	}
}

func (pg *settingsPage) handle(common pageCommon, win *Window) {

	if pg.isDarkModeOn.Changed() {
		win.theme.SwitchDarkMode(pg.isDarkModeOn.Value)
		pg.wal.SaveConfigValueForKey("isDarkModeOn", pg.isDarkModeOn.Value)
		win.loadPage(common.icons)
	}

	if pg.spendUnconfirmed.Changed() {
		pg.wal.SaveConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey, pg.spendUnconfirmed.Value)
	}

	if pg.beepNewBlocks.Changed() {
		pg.wal.SaveConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey, pg.beepNewBlocks.Value)
	}

	for pg.changeStartupPass.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ChangeStartupPasswordTemplate,
				title:    "Change startup password",
				confirm: func(oldPass, newPass string) {
					pg.wal.ChangeStartupPassphrase(oldPass, newPass, pg.errorReceiver)
				},
				confirmText: "Change",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
	}

	if pg.startupPassword.Changed() {
		if pg.startupPassword.Value {
			go func() {
				common.modalReceiver <- &modalLoad{
					template: SetStartupPasswordTemplate,
					title:    "Create a startup password",
					confirm: func(pass string) {
						pg.wal.SetStartupPassphrase(pass, pg.errorReceiver)
					},
					confirmText: "Create",
					cancel:      common.closeModal,
					cancelText:  "Cancel",
				}
			}()
			return
		}
		go func() {
			common.modalReceiver <- &modalLoad{
				template: RemoveStartupPasswordTemplate,
				title:    "Confirm to turn off startup password",
				confirm: func(pass string) {
					pg.wal.RemoveStartupPassphrase(pass, pg.errorReceiver)
				},
				confirmText: "Confirm",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	specificPeerKey := dcrlibwallet.SpvPersistentPeerAddressesConfigKey
	if pg.connectToPeer.Changed() {
		if pg.connectToPeer.Value {
			go func() {
				common.modalReceiver <- &modalLoad{
					template: ConnectToSpecificPeerTemplate,
					title:    "Connect to specific peer",
					confirm: func(ipAddress string) {
						if ipAddress != "" {
							pg.wal.SaveConfigValueForKey(specificPeerKey, ipAddress)
							common.closeModal()
						}
					},
					confirmText: "Connect",
					cancel:      common.closeModal,
					cancelText:  "Cancel",
				}
			}()
			return
		}
		pg.wal.RemoveUserConfigValueForKey(specificPeerKey)
	}
	for pg.updateConnectToPeer.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ChangeSpecificPeerTemplate,
				title:    "Change specific peer",
				confirm: func(ipAddress string) {
					if ipAddress != "" {
						pg.wal.SaveConfigValueForKey(specificPeerKey, ipAddress)
						common.closeModal()
					}
				},
				confirmText: "Done",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
	}

	userAgentKey := dcrlibwallet.UserAgentConfigKey
	for pg.updateUserAgent.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: UserAgentTemplate,
				title:    "Change user agent",
				confirm: func(agent string) {
					if agent != "" {
						pg.wal.SaveConfigValueForKey(userAgentKey, agent)
						common.closeModal()
					}
				},
				confirmText: "Done",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
	}

	if pg.userAgent.Changed() {
		if pg.userAgent.Value {
			go func() {
				common.modalReceiver <- &modalLoad{
					template: UserAgentTemplate,
					title:    "Set up user agent",
					confirm: func(agent string) {
						if agent != "" {
							pg.wal.SaveConfigValueForKey(userAgentKey, agent)
							common.closeModal()
						}
					},
					confirmText: "Set up",
					cancel:      common.closeModal,
					cancelText:  "Cancel",
				}
			}()
			return
		}
		pg.wal.RemoveUserConfigValueForKey(userAgentKey)
	}

	for pg.currencyConversion.Clicked() {
		pg.isCurrencyModalOpen = true
	}

	if pg.radioButtonsGroup.Changed() {
		if pg.radioButtonsGroup.Value == "None" {
			pg.initialValue = USDExchangeValue
		} else {
			pg.initialValue = "None"
		}
		pg.wal.SaveConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey, pg.radioButtonsGroup.Value)
	}

	for pg.cancel.Button.Clicked() {
		pg.wal.SaveConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey, pg.initialValue)
		pg.isCurrencyModalOpen = false
	}

	if pg.confirm.Button.Clicked() {
		pg.initialValue = pg.radioButtonsGroup.Value
		pg.isCurrencyModalOpen = false
	}

	select {
	case err := <-pg.errorReceiver:
		if err.Error() == "invalid_passphrase" {
			e := "Password is incorrect"
			common.notify(e, false)
			return
		}
		common.notify(err.Error(), false)
	default:
	}
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

	pg.currencyValue = pg.wal.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if pg.currencyValue == "" {
		pg.currencyValue = "None"
	}
	pg.radioButtonsGroup.Value = pg.currencyValue
}
