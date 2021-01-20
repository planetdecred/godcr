package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageSettings = "SettingsPage"

type settingsPage struct {
	theme      *decredmaterial.Theme
	walletInfo *wallet.MultiWalletInfo
	wal        *wallet.Wallet

	currencyConversion decredmaterial.IconButton
	connectToPeer      decredmaterial.IconButton
	userAgent          decredmaterial.IconButton

	spendUnconfirm  *widget.Bool
	startupPassword *widget.Bool
	notificationW   *widget.Bool

	line     *decredmaterial.Line
	errChann chan error
}

func (win *Window) SettingsPage(common pageCommon) layout.Widget {
	icon := common.icons.chevronRight
	pg := &settingsPage{
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
	}
	pg.line.Height = 2
	pg.line.Color = common.theme.Color.Background

	color := common.theme.Color.LightGray
	zeroInset := layout.UniformInset(values.MarginPadding0)

	pg.currencyConversion.Color, pg.currencyConversion.Inset = color, zeroInset
	pg.connectToPeer.Color, pg.connectToPeer.Inset = color, zeroInset
	pg.userAgent.Color, pg.userAgent.Inset = color, zeroInset

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *settingsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		page := SubPage{
			title:      "Settings",
			walletName: "",
			back: func() {
				*common.page = PageWallet
			},
			body: func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(pg.general()),
					layout.Rigid(pg.security()),
					layout.Rigid(pg.notification()),
					layout.Rigid(pg.connection()),
				)
			},
			infoTemplate: "",
		}
		return common.SubPageLayout(gtx, page)
	}

	return common.Layout(gtx, body)
}

func (pg *settingsPage) general() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, "General", func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(pg.bottomSectionLabel("Spending unconfirmed funds")),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return pg.theme.Switch(pg.spendUnconfirm).Layout(gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					m := values.MarginPadding10
					return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
						pg.line.Width = gtx.Constraints.Max.X
						return pg.line.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(pg.bottomSectionLabel("Currency conversion")),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
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
				}),
			)
		})
	}
}

func (pg *settingsPage) notification() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, "Notification", func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel("Beep for new blocks")),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.theme.Switch(pg.notificationW).Layout(gtx)
					})
				}),
			)
		})
	}
}

func (pg *settingsPage) security() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, "Security", func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel("Startup password")),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.theme.Switch(pg.startupPassword).Layout(gtx)
					})
				}),
			)
		})
	}
}

func (pg *settingsPage) connection() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, "Connection", func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(pg.bottomSectionLabel("Connect to specific peer")),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return pg.connectToPeer.Layout(gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					m := values.MarginPadding10
					return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
						pg.line.Width = gtx.Constraints.Max.X
						return pg.line.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(pg.bottomSectionLabel("User agent")),
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

func (pg *settingsPage) pageSections(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
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

func (pg *settingsPage) bottomSectionLabel(title string) layout.Widget {
	return func(gtx C) D {
		return pg.theme.Body1(title).Layout(gtx)
	}
}

func (pg *settingsPage) handle(common pageCommon) {
	// for pg.changePass.Button.Clicked() {
	// 	walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
	// 	go func() {
	// 		common.modalReceiver <- &modalLoad{
	// 			template: ChangePasswordTemplate,
	// 			title:    "Change spending password",
	// 			confirm: func(oldPass, newPass string) {
	// 				// pg.wal.ChangeWalletPassphrase(walletID, oldPass, newPass, pg.errChann)
	// 			},
	// 			confirmText: "Change",
	// 			cancel:      common.closeModal,
	// 			cancelText:  "Cancel",
	// 		}
	// 	}()
	// 	break
	// }

	// for pg.rescan.Button.Clicked() {
	// 	walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
	// 	go func() {
	// 		common.modalReceiver <- &modalLoad{
	// 			template: RescanWalletTemplate,
	// 			title:    "Rescan blockchain",
	// 			confirm: func() {
	// 				err := pg.wal.RescanBlocks(walletID)
	// 				if err != nil {
	// 					if err.Error() == "not_connected" {
	// 						common.Notify("Not connected to decred network", false)
	// 						return
	// 					}
	// 					common.Notify(err.Error(), false)
	// 					return
	// 				}
	// 				msg := "Rescan initiated (check in overview)"
	// 				common.Notify(msg, true)
	// 				go func() {
	// 					common.modalReceiver <- &modalLoad{}
	// 				}()
	// 			},
	// 			confirmText: "Rescan",
	// 			cancel:      common.closeModal,
	// 			cancelText:  "Cancel",
	// 		}
	// 	}()
	// 	break
	// }

	// for pg.deleteWallet.Button.Clicked() {
	// 	go func() {
	// 		common.modalReceiver <- &modalLoad{
	// 			template: ConfirmRemoveTemplate,
	// 			title:    "Remove wallet from device",
	// 			confirm: func() {
	// 				walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
	// 				go func() {
	// 					common.modalReceiver <- &modalLoad{
	// 						template: PasswordTemplate,
	// 						title:    "Confirm to remove",
	// 						confirm: func(pass string) {
	// 							pg.wal.DeleteWallet(walletID, []byte(pass), pg.errChann)
	// 						},
	// 						confirmText: "Confirm",
	// 						cancel:      common.closeModal,
	// 						cancelText:  "Cancel",
	// 					}
	// 				}()
	// 			},
	// 			confirmText: "Remove",
	// 			cancel:      common.closeModal,
	// 			cancelText:  "Cancel",
	// 		}
	// 	}()
	// 	break
	// }

	// select {
	// case err := <-pg.errChann:
	// 	if err.Error() == "invalid_passphrase" {
	// 		e := "Password is incorrect"
	// 		common.Notify(e, false)
	// 	}
	// default:
	// }
}
