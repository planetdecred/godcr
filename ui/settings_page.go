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

	currencyConversion decredmaterial.IconButton
	connectToPeer      decredmaterial.IconButton
	userAgent          decredmaterial.IconButton
	changeStartupPass  decredmaterial.IconButton

	spendUnconfirm  *widget.Bool
	startupPassword *widget.Bool
	notificationW   *widget.Bool

	line *decredmaterial.Line

	isStartupPassword bool
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

		spendUnconfirm:  new(widget.Bool),
		startupPassword: new(widget.Bool),
		notificationW:   new(widget.Bool),

		line:     common.theme.Line(),
		errChann: common.errorChannels[PageSettings],

		currencyConversion: common.theme.PlainIconButton(new(widget.Clickable), icon),
		connectToPeer:      common.theme.PlainIconButton(new(widget.Clickable), icon),
		userAgent:          common.theme.PlainIconButton(new(widget.Clickable), icon),
		changeStartupPass:  common.theme.PlainIconButton(new(widget.Clickable), icon),
	}
	pg.line.Height = 2
	pg.line.Color = common.theme.Color.Background

	color := common.theme.Color.LightGray
	zeroInset := layout.UniformInset(values.MarginPadding0)

	pg.currencyConversion.Color, pg.currencyConversion.Inset = color, zeroInset
	pg.connectToPeer.Color, pg.connectToPeer.Inset = color, zeroInset
	pg.userAgent.Color, pg.userAgent.Inset = color, zeroInset
	pg.changeStartupPass.Color, pg.changeStartupPass.Inset = color, zeroInset

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *settingsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.updateStartupPasswordSetting()
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
					return pg.subSection(gtx, "Spending unconfirmed funds", func(gtx C) D {
						return pg.theme.Switch(pg.spendUnconfirm).Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.lineSeparator(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.subSection(gtx, "Currency conversion", func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := pg.theme.Body2("None")
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
			return pg.subSection(gtx, "Beep for new blocks", func(gtx C) D {
				return pg.theme.Switch(pg.notificationW).Layout(gtx)
			})
		})
	}
}

func (pg *settingsPage) security() layout.Widget {
	return func(gtx C) D {
		return pg.mainSection(gtx, "Security", func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.subSection(gtx, "Startup password", func(gtx C) D {
						return pg.theme.Switch(pg.startupPassword).Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if pg.isStartupPassword {
						return pg.lineSeparator(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.isStartupPassword {
						return pg.subSection(gtx, "Change startup password", func(gtx C) D {
							return pg.changeStartupPass.Layout(gtx)
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
					return pg.subSection(gtx, "Connect to specific peer", func(gtx C) D {
						return pg.connectToPeer.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.lineSeparator(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(pg.subSectionLabel("User agent")),
								layout.Rigid(func(gtx C) D {
									txt := pg.theme.Body2("For exchange rate fetching")
									txt.Color = pg.theme.Color.Gray
									return txt.Layout(gtx)
								}),
							)
						}),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return pg.userAgent.Layout(gtx)
							})
						}),
					)
				}),
			)
		})
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

func (pg *settingsPage) subSectionLabel(title string) layout.Widget {
	return func(gtx C) D {
		return pg.theme.Body1(title).Layout(gtx)
	}
}

func (pg *settingsPage) lineSeparator(gtx layout.Context) layout.Dimensions {
	m := values.MarginPadding10
	return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
		pg.line.Width = gtx.Constraints.Max.X
		return pg.line.Layout(gtx)
	})
}

func (pg *settingsPage) handle(common pageCommon) {
	notImplemented := "functionality not yet implemented"

	if pg.spendUnconfirm.Changed() {
		pg.wal.SpendUnconfirmed(pg.spendUnconfirm.Value)
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

	for pg.connectToPeer.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ConnectToSpecificPeerTemplate,
				title:    "Connect to specific peer",
				confirm: func(ipAddress string) {
					common.Notify(notImplemented, true)
					common.closeModal()
				},
				confirmText: "Connect",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
	}

	for pg.userAgent.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: UserAgentTemplate,
				title:    "Set up user agent",
				confirm: func(agent string) {
					common.Notify(notImplemented, true)
					common.closeModal()
				},
				confirmText: "Set up",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
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

func (pg *settingsPage) updateStartupPasswordSetting() {
	isSet := pg.wal.IsStartupSecuritySet()
	if isSet {
		pg.startupPassword.Value = true
		pg.isStartupPassword = true
	} else {
		pg.startupPassword.Value = false
		pg.isStartupPassword = false
	}
}
