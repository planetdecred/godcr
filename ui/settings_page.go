package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageSettings = "Settings"

type settingsPage struct {
	pageContainer layout.List
	theme         *decredmaterial.Theme
	walletInfo    *wallet.MultiWalletInfo
	wal           *wallet.Wallet

	currencyConversion  decredmaterial.IconButton
	updateConnectToPeer decredmaterial.IconButton
	updateUserAgent     decredmaterial.IconButton
	changeStartupPass   decredmaterial.IconButton

	spendUnconfirmed *widget.Bool
	startupPassword  *widget.Bool
	beepNewBlocks    *widget.Bool
	connectToPeer    *widget.Bool
	userAgent        *widget.Bool

	peerLabel decredmaterial.Label

	line *decredmaterial.Line

	isStartupPassword bool
	peerAddr          string
	agentValue        string
	errChann          chan error
}

func (win *Window) SettingsPage(common pageCommon) layout.Widget {
	icon := common.icons.chevronRight
	pg := &settingsPage{
		pageContainer: layout.List{
			Axis: layout.Vertical,
		},
		theme:      common.theme,
		walletInfo: win.walletInfo,
		wal:        common.wallet,

		spendUnconfirmed: new(widget.Bool),
		startupPassword:  new(widget.Bool),
		beepNewBlocks:    new(widget.Bool),
		connectToPeer:    new(widget.Bool),
		userAgent:        new(widget.Bool),

		line:     common.theme.Line(),
		errChann: common.errorChannels[PageSettings],

		currencyConversion:  common.theme.PlainIconButton(new(widget.Clickable), icon),
		updateConnectToPeer: common.theme.PlainIconButton(new(widget.Clickable), icon),
		updateUserAgent:     common.theme.PlainIconButton(new(widget.Clickable), icon),
		changeStartupPass:   common.theme.PlainIconButton(new(widget.Clickable), icon),
	}
	pg.line.Height = 2
	pg.line.Color = common.theme.Color.Background

	color := common.theme.Color.LightGray
	zeroInset := layout.UniformInset(values.MarginPadding0)

	pg.peerLabel = common.theme.Body1("")
	pg.peerLabel.Color = common.theme.Color.Gray

	pg.currencyConversion.Color, pg.currencyConversion.Inset = color, zeroInset
	pg.updateConnectToPeer.Color, pg.updateConnectToPeer.Inset = color, zeroInset
	pg.updateUserAgent.Color, pg.updateUserAgent.Inset = color, zeroInset
	pg.changeStartupPass.Color, pg.changeStartupPass.Inset = color, zeroInset

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *settingsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.updateSettingOptions()

	body := func(gtx C) D {
		page := SubPage{
			title: "Settings",
			back: func() {
				*common.page = PageMore
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

	return common.Layout(gtx, body)
}

func (pg *settingsPage) general() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, "General", func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSectionWithSwitch(gtx, "Spending unconfirmed funds", pg.spendUnconfirmed)
				}),
				layout.Rigid(pg.lineSeparator()),
				layout.Rigid(func(gtx C) D {
					return pg.subSection(gtx, "Currency conversion", func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := pg.theme.Body1("None")
								txt.Color = pg.theme.Color.Gray
								return txt.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return pg.currencyConversion.Layout(gtx)
							}),
						)
					})
				}),
			)
		})
	}
}

func (pg *settingsPage) notification() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, "Notification", func(gtx C) D {
			return pg.subSectionWithSwitch(gtx, "Beep for new blocks", pg.beepNewBlocks)
		})
	}
}

func (pg *settingsPage) security() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, "Security", func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSectionWithSwitch(gtx, "Startup password", pg.startupPassword)
				}),
				layout.Rigid(func(gtx C) D {
					if pg.isStartupPassword {
						return pg.conditionalDisplay(gtx, func(gtx C) D {
							return pg.subSection(gtx, "Change startup password", func(gtx C) D {
								return pg.changeStartupPass.Layout(gtx)
							})
						})
					}
					return layout.Dimensions{}
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
					return pg.subSectionWithSwitch(gtx, "Connect to specific peer", pg.connectToPeer)
				}),
				layout.Rigid(func(gtx C) D {
					if pg.peerAddr != "" {
						return pg.conditionalDisplay(gtx, func(gtx C) D {
							return pg.subSection(gtx, "Change specific peer", func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return pg.peerLabel.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return pg.updateConnectToPeer.Layout(gtx)
									}),
								)
							})
						})
					}
					return layout.Dimensions{}
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
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(pg.subSectionLabel("User agent")),
							layout.Rigid(func(gtx C) D {
								txt := pg.theme.Body2("For exchange rate fetching")
								txt.Color = pg.theme.Color.Gray
								return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									return txt.Layout(gtx)
								})
							}),
						)
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
				if pg.agentValue != "" {
					return pg.conditionalDisplay(gtx, func(gtx C) D {
						return pg.subSection(gtx, "Change user agent", func(gtx C) D {
							return pg.updateUserAgent.Layout(gtx)
						})
					})
				}
				return layout.Dimensions{}
			}),
		)
	}
}

func (pg *settingsPage) mainSection(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
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
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(pg.subSectionLabel(title)),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, body)
		}),
	)
}

func (pg *settingsPage) subSectionWithSwitch(gtx layout.Context, title string, option *widget.Bool) layout.Dimensions {
	return pg.subSection(gtx, title, func(gtx C) D {
		return pg.theme.Switch(option).Layout(gtx)
	})
}

func (pg *settingsPage) conditionalDisplay(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(pg.lineSeparator()),
		layout.Rigid(body),
	)
}

func (pg *settingsPage) subSectionLabel(title string) layout.Widget {
	return func(gtx C) D {
		return pg.theme.Body1(title).Layout(gtx)
	}
}

func (pg *settingsPage) lineSeparator() layout.Widget {
	m := values.MarginPadding10
	return func(gtx C) D {
		return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
			pg.line.Width = gtx.Constraints.Max.X
			return pg.line.Layout(gtx)
		})
	}
}

func (pg *settingsPage) handle(common pageCommon) {
	if pg.spendUnconfirmed.Changed() {
		pg.wal.SpendUnconfirmed(pg.spendUnconfirmed.Value)
	}

	if pg.beepNewBlocks.Changed() {
		pg.wal.BeepNewBlocks(pg.beepNewBlocks.Value)
	}

	for pg.changeStartupPass.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ChangeStartupPasswordTemplate,
				title:    "Change startup password",
				confirm: func(oldPass, newPass string) {
					pg.wal.ChangeStartupPassphrase(oldPass, newPass, pg.errChann)
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
						pg.wal.SetStartupPassphrase(pass, pg.errChann)
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
					pg.wal.RemoveStartupPassphrase(pass, pg.errChann)
				},
				confirmText: "Confirm",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	if pg.connectToPeer.Changed() {
		if pg.connectToPeer.Value {
			go func() {
				common.modalReceiver <- &modalLoad{
					template: ConnectToSpecificPeerTemplate,
					title:    "Connect to specific peer",
					confirm: func(ipAddress string) {
						if ipAddress != "" {
							pg.wal.ConnectToPeer(ipAddress)
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
		go func() {
			common.modalReceiver <- &modalLoad{
				template: RemoveSpecificPeerTemplate,
				title:    "Turn off connect to specific peer?",
				confirm: func() {
					pg.wal.RemoveConnectToPeer()
					common.closeModal()
				},
				confirmText: "Turn off",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	for pg.updateConnectToPeer.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ChangeSpecificPeerTemplate,
				title:    "Change specific peer",
				confirm: func(ipAddress string) {
					if ipAddress != "" {
						pg.wal.ConnectToPeer(ipAddress)
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

	for pg.updateUserAgent.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: UserAgentTemplate,
				title:    "Change user agent",
				confirm: func(agent string) {
					if agent != "" {
						pg.wal.UserAgent(agent)
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
							pg.wal.UserAgent(agent)
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
		go func() {
			common.modalReceiver <- &modalLoad{
				template: RemoveUserAgentTemplate,
				title:    "Turn off user agent?",
				confirm: func() {
					pg.wal.RemoveUserAgent()
					common.closeModal()
				},
				confirmText: "Turn off",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	select {
	case err := <-pg.errChann:
		if err.Error() == "invalid_passphrase" {
			e := "Password is incorrect"
			common.Notify(e, false)
			return
		}
		common.Notify(err.Error(), false)
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

	isSpendUnconfirmed := pg.wal.IsSpendUnconfirmed()
	pg.spendUnconfirmed.Value = false
	if isSpendUnconfirmed {
		pg.spendUnconfirmed.Value = true
	}

	beep := pg.wal.IsBeepNewBlocks()
	pg.beepNewBlocks.Value = false
	if beep {
		pg.beepNewBlocks.Value = true
	}

	pg.peerAddr = pg.wal.GetConnectToPeerValue()
	pg.connectToPeer.Value = false
	if pg.peerAddr != "" {
		pg.peerLabel.Text = pg.peerAddr
		pg.connectToPeer.Value = true
	}

	pg.agentValue = pg.wal.GetUserAgent()
	pg.userAgent.Value = false
	if pg.agentValue != "" {
		pg.userAgent.Value = true
	}
}
