package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageWalletSettings = "WalletSettings"

type walletSettingsPage struct {
	theme      *decredmaterial.Theme
	walletInfo *wallet.MultiWalletInfo
	wal        *wallet.Wallet

	changePass, rescan, deleteWallet *widget.Clickable

	notificationW *widget.Bool
	errorReceiver chan error

	chevronRightIcon *widget.Icon
}

func (win *Window) WalletSettingsPage(common pageCommon) layout.Widget {
	pg := &walletSettingsPage{
		theme:         common.theme,
		walletInfo:    win.walletInfo,
		wal:           common.wallet,
		notificationW: new(widget.Bool),
		errorReceiver: make(chan error),

		changePass:   new(widget.Clickable),
		rescan:       new(widget.Clickable),
		deleteWallet: new(widget.Clickable),

		chevronRightIcon: common.icons.chevronRight,
	}

	pg.chevronRightIcon.Color = pg.theme.Color.LightGray

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *walletSettingsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	beep := pg.wal.ReadBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey)
	pg.notificationW.Value = false
	if beep {
		pg.notificationW.Value = true
	}

	body := func(gtx C) D {
		page := SubPage{
			title:      values.String(values.StrSettings),
			walletName: common.info.Wallets[*common.selectedWallet].Name,
			back: func() {
				common.changePage(PageWallet)
			},
			body: func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !common.info.Wallets[*common.selectedWallet].IsWatchingOnly {
							return pg.changePassphrase()(gtx)
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(pg.notification()),
					layout.Rigid(pg.debug()),
					layout.Rigid(pg.dangerZone()),
				)
			},
			infoTemplate: "",
		}
		return common.SubPageLayout(gtx, page)
	}
	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, body)
	})
}

func (pg *walletSettingsPage) changePassphrase() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrSpendingPassword), pg.changePass, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel(values.String(values.StrChangeSpendingPass))),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
					})
				}),
			)
		})
	}
}

func (pg *walletSettingsPage) notification() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrNotifications), nil, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel(values.String(values.StrBeepForNewBlocks))),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.theme.Switch(pg.notificationW).Layout(gtx)
					})
				}),
			)
		})
	}
}

func (pg *walletSettingsPage) debug() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrDebug), pg.rescan, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel(values.String(values.StrRescanBlockchain))),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
					})
				}),
			)
		})
	}
}

func (pg *walletSettingsPage) dangerZone() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrDangerZone), pg.deleteWallet, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel(values.String(values.StrRemoveWallet))),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
					})
				}),
			)
		})
	}
}

func (pg *walletSettingsPage) pageSections(gtx layout.Context, title string, clickable *widget.Clickable, body layout.Widget) layout.Dimensions {
	dims := func(gtx layout.Context, title string, body layout.Widget) D {
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
	}

	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			if clickable == nil {
				return dims(gtx, title, body)
			}
			return decredmaterial.Clickable(gtx, clickable, func(gtx C) D {
				return dims(gtx, title, body)
			})
		})
	})
}

func (pg *walletSettingsPage) bottomSectionLabel(title string) layout.Widget {
	return func(gtx C) D {
		return pg.theme.Body1(title).Layout(gtx)
	}
}

func (pg *walletSettingsPage) resetSelectedWallet(common pageCommon) {
	common.wallAcctSelector.selectedSendWallet = 0
	common.wallAcctSelector.selectedSendAccount = 0

	common.wallAcctSelector.selectedReceiveWallet = 0
	common.wallAcctSelector.selectedReceiveAccount = 0
}

func (pg *walletSettingsPage) handle(common pageCommon) {
	for pg.changePass.Clicked() {
		walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ChangePasswordTemplate,
				title:    values.String(values.StrChangeSpendingPass),
				confirm: func(oldPass, newPass string) {
					pg.wal.ChangeWalletPassphrase(walletID, oldPass, newPass, pg.errorReceiver)
				},
				confirmText: values.String(values.StrChange),
				cancel:      common.closeModal,
				cancelText:  values.String(values.StrCancel),
			}
		}()
		break
	}

	for pg.rescan.Clicked() {
		walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
		go func() {
			common.modalReceiver <- &modalLoad{
				template: RescanWalletTemplate,
				title:    values.String(values.StrRescanBlockchain),
				confirm: func() {
					err := pg.wal.RescanBlocks(walletID)
					if err != nil {
						if err.Error() == dcrlibwallet.ErrNotConnected {
							common.notify(values.String(values.StrNotConnected), false)
							return
						}
						common.notify(err.Error(), false)
						return
					}
					msg := values.String(values.StrRescanProgressNotification)
					common.notify(msg, true)
					go func() {
						common.modalReceiver <- &modalLoad{}
					}()
				},
				confirmText: values.String(values.StrRescan),
				cancel:      common.closeModal,
				cancelText:  values.String(values.StrCancel),
			}
		}()
		break
	}

	if pg.notificationW.Changed() {
		pg.wal.SaveConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey, pg.notificationW.Value)
	}

	for pg.deleteWallet.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: ConfirmRemoveTemplate,
				title:    values.String(values.StrRemoveWallet),
				confirm: func() {
					walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
					go func() {
						common.modalReceiver <- &modalLoad{
							template: PasswordTemplate,
							title:    values.String(values.StrConfirmToRemove),
							confirm: func(pass string) {
								pg.wal.DeleteWallet(walletID, []byte(pass), pg.errorReceiver)
								pg.resetSelectedWallet(common)
							},
							confirmText: values.String(values.StrConfirm),
							cancel:      common.closeModal,
							cancelText:  values.String(values.StrCancel),
						}
					}()
				},
				confirmText: values.String(values.StrRemove),
				cancel:      common.closeModal,
				cancelText:  values.String(values.StrCancel),
			}
		}()
		break
	}

	select {
	case err := <-pg.errorReceiver:
		if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
			e := values.String(values.StrInvalidPassphrase)
			common.notify(e, false)
		}
	default:
	}
}
