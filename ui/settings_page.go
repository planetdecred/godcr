package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageSettings = "Settings"

type settingsPage struct {
	theme      *decredmaterial.Theme
	walletInfo *wallet.MultiWalletInfo
	wal        *wallet.Wallet

	changePass   decredmaterial.IconButton
	setupCoin    decredmaterial.IconButton
	rescan       decredmaterial.IconButton
	deleteWallet decredmaterial.IconButton

	notificationW *widget.Bool
	chevronIcon   *widget.Icon
	line          *decredmaterial.Line
	errChann      chan error
}

func (win *Window) SettingsPage(common pageCommon) layout.Widget {
	pg := &settingsPage{
		theme:         common.theme,
		walletInfo:    win.walletInfo,
		wal:           common.wallet,
		notificationW: new(widget.Bool),
		line:          common.theme.Line(),
		errChann:      common.errorChannels[PageSettings],

		changePass: decredmaterial.IconButton{
			material.IconButtonStyle{
				Icon:       common.icons.chevronRight,
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      common.theme.Color.LightGray,
				Inset:      layout.UniformInset(values.MarginPadding0),
				Button:     new(widget.Clickable),
			},
		},
		setupCoin: decredmaterial.IconButton{
			material.IconButtonStyle{
				Icon:       common.icons.chevronRight,
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      common.theme.Color.LightGray,
				Inset:      layout.UniformInset(values.MarginPadding0),
				Button:     new(widget.Clickable),
			},
		},
		rescan: decredmaterial.IconButton{
			material.IconButtonStyle{
				Icon:       common.icons.chevronRight,
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      common.theme.Color.LightGray,
				Inset:      layout.UniformInset(values.MarginPadding0),
				Button:     new(widget.Clickable),
			},
		},
		deleteWallet: decredmaterial.IconButton{
			material.IconButtonStyle{
				Icon:       common.icons.chevronRight,
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      common.theme.Color.LightGray,
				Inset:      layout.UniformInset(values.MarginPadding0),
				Button:     new(widget.Clickable),
			},
		},
	}
	pg.line.Height = 2
	pg.line.Color = common.theme.Color.Background

	pg.chevronIcon = common.icons.chevronRight
	pg.chevronIcon.Color = common.theme.Color.Background

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

// main settings layout
func (pg *settingsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		page := SubPage{
			title:      "Settings",
			walletName: common.info.Wallets[*common.selectedWallet].Name,
			back: func() {
				*common.page = PageWallet
			},
			body: func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(pg.changePassphrase()),
						layout.Rigid(pg.notification()),
						layout.Rigid(pg.debug()),
						layout.Rigid(pg.dangerZone()),
					)
				})
			},
			infoTemplate: "",
		}
		return common.SubPageLayoutWithoutInfo(gtx, page)
	}

	return common.Layout(gtx, body)
}

func (pg *settingsPage) changePassphrase() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, "Spending password", func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel("Change spending password")),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.changePass.Layout(gtx)
					})
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

func (pg *settingsPage) debug() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, "Debug", func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel("Rescan blockchain")),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.rescan.Layout(gtx)
					})
				}),
			)
		})
	}
}

func (pg *settingsPage) dangerZone() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, "Danger zone", func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel("Remove wallet from device")),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.deleteWallet.Layout(gtx)
					})
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
						return txt.Layout(gtx)
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
	for pg.changePass.Button.Clicked() {
		walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ChangePasswordTemplate,
				title:    "Change spending password",
				confirm: func(oldPass, newPass string) {
					pg.wal.ChangeWalletPassphrase(walletID, oldPass, newPass, pg.errChann)
				},
				confirmText: "Change",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
	}

	for pg.rescan.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: RescanWalletTemplate,
				title:    "Rescan blockchain",
				confirm: func() {
					pg.wal.RestartSpvSync()
					e := "Rescan initiated (check in overview)"
					common.Notify(e, true)
					go func() {
						common.modalReceiver <- &modalLoad{}
					}()
				},
				confirmText: "Rescan",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
	}

	for pg.deleteWallet.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ConfirmRemoveTemplate,
				title:    "Remove wallet from device",
				confirm: func() {
					walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
					go func() {
						common.modalReceiver <- &modalLoad{
							template: PasswordTemplate,
							title:    "Confirm to remove",
							confirm: func(pass string) {
								pg.wal.DeleteWallet(walletID, []byte(pass), pg.errChann)
							},
							confirmText: "Confirm",
							cancel:      common.closeModal,
							cancelText:  "Cancel",
						}
					}()
				},
				confirmText: "Remove",
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
		}
	default:
	}
}
