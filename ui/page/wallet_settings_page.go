package page

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletSettingsPageID = "WalletSettings"

type WalletSettingsPage struct {
	*load.Load
	wallet *dcrlibwallet.Wallet

	changePass, rescan, deleteWallet *widget.Clickable

	notificationW *widget.Bool

	chevronRightIcon *widget.Icon
	backButton       decredmaterial.IconButton
}

func NewWalletSettingsPage(l *load.Load, wal *dcrlibwallet.Wallet) *WalletSettingsPage {
	pg := &WalletSettingsPage{
		Load:          l,
		wallet:        wal,
		notificationW: new(widget.Bool),
		changePass:    new(widget.Clickable),
		rescan:        new(widget.Clickable),
		deleteWallet:  new(widget.Clickable),

		chevronRightIcon: l.Icons.ChevronRight,
	}

	pg.chevronRightIcon.Color = l.Theme.Color.LightGray
	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
}

func (pg *WalletSettingsPage) OnResume() {

}

func (pg *WalletSettingsPage) Layout(gtx layout.Context) layout.Dimensions {

	beep := pg.wallet.ReadBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey, false)
	pg.notificationW.Value = beep
	if beep {
		pg.notificationW.Value = true
	}

	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrSettings),
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			Back: func() {
				pg.ChangePage(WalletPageID)
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !pg.wallet.IsWatchingOnlyWallet() {
							return pg.changePassphrase()(gtx)
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(pg.notification()),
					layout.Rigid(pg.debug()),
					layout.Rigid(pg.dangerZone()),
				)
			},
		}
		return sp.Layout(gtx)
	}
	return components.UniformPadding(gtx, body)
}

func (pg *WalletSettingsPage) changePassphrase() layout.Widget {
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

func (pg *WalletSettingsPage) notification() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrNotifications), nil, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(pg.bottomSectionLabel(values.String(values.StrBeepForNewBlocks))),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.Theme.Switch(pg.notificationW).Layout(gtx)
					})
				}),
			)
		})
	}
}

func (pg *WalletSettingsPage) debug() layout.Widget {
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

func (pg *WalletSettingsPage) dangerZone() layout.Widget {
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

func (pg *WalletSettingsPage) pageSections(gtx layout.Context, title string, clickable *widget.Clickable, body layout.Widget) layout.Dimensions {
	dims := func(gtx layout.Context, title string, body layout.Widget) D {
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := pg.Theme.Body2(title)
					txt.Color = pg.Theme.Color.Gray
					return txt.Layout(gtx)
				}),
				layout.Rigid(body),
			)
		})
	}

	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			if clickable == nil {
				return dims(gtx, title, body)
			}
			return decredmaterial.Clickable(gtx, clickable, func(gtx C) D {
				return dims(gtx, title, body)
			})
		})
	})
}

func (pg *WalletSettingsPage) bottomSectionLabel(title string) layout.Widget {
	return func(gtx C) D {
		return pg.Theme.Body1(title).Layout(gtx)
	}
}

func (pg *WalletSettingsPage) Handle() {
	for pg.changePass.Clicked() {
		modal.NewPasswordModal(pg.Load).
			Title(values.String(values.StrChangeSpendingPass)).
			Hint("Current spending password").
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
				go func() {
					err := pg.wallet.UnlockWallet([]byte(password))
					if err != nil {
						pm.SetError(err.Error())
						pm.SetLoading(false)
						return
					}
					pg.wallet.LockWallet()
					pm.Dismiss()

					// change password
					modal.NewCreatePasswordModal(pg.Load).
						Title(values.String(values.StrChangeSpendingPass)).
						EnableName(false).
						PasswordHint("New spending password").
						ConfirmPasswordHint("Confirm new spending password").
						PasswordCreated(func(walletName, newPassword string, m *modal.CreatePasswordModal) bool {
							go func() {
								err := pg.WL.MultiWallet.ChangePrivatePassphraseForWallet(pg.wallet.ID, []byte(password),
									[]byte(newPassword), dcrlibwallet.PassphraseTypePass)
								if err != nil {
									m.SetError(err.Error())
									m.SetLoading(false)
									return
								}
								m.Dismiss()
							}()
							return false
						}).Show()

				}()

				return false
			}).Show()
		break
	}

	for pg.rescan.Clicked() {
		go func() {
			info := modal.NewInfoModal(pg.Load).
				Title(values.String(values.StrRescanBlockchain)).
				Body("Rescanning may help resolve some balance errors. This will take some time, as it scans the entire"+
					" blockchain for transactions").
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(values.String(values.StrRescan), func() {
					err := pg.WL.MultiWallet.RescanBlocks(pg.wallet.ID)
					if err != nil {
						if err.Error() == dcrlibwallet.ErrNotConnected {
							pg.CreateToast(values.String(values.StrNotConnected), false)
							return
						}
						pg.CreateToast(err.Error(), false)
						return
					}
					msg := values.String(values.StrRescanProgressNotification)
					pg.CreateToast(msg, true)
				})

			pg.ShowModal(info)
		}()
		break
	}

	if pg.notificationW.Changed() {
		pg.wallet.SetBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey, pg.notificationW.Value)
	}

	for pg.deleteWallet.Clicked() {
		modal.NewInfoModal(pg.Load).
			Title(values.String(values.StrRemoveWallet)).
			Body("Make sure to have the seed phrase backed up before removing the wallet").
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton(values.String(values.StrRemove), func() {

				modal.NewPasswordModal(pg.Load).
					Title(values.String(values.StrConfirmToRemove)).
					NegativeButton(values.String(values.StrCancel), func() {}).
					PositiveButtonStyle(pg.Load.Theme.Color.Surface, pg.Load.Theme.Color.Danger).
					PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
						go func() {
							err := pg.WL.MultiWallet.DeleteWallet(pg.wallet.ID, []byte(password))
							if err != nil {
								pm.SetError(err.Error())
								pm.SetLoading(false)
								return
							}
							pm.Dismiss()
							pm.ChangePage(WalletPageID)
						}()
						return false
					}).Show()

			}).Show()
		break
	}
}

func (pg *WalletSettingsPage) OnClose() {}
