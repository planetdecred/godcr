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

	changePass, rescan, deleteWallet *decredmaterial.Clickable

	chevronRightIcon *widget.Icon
	backButton       decredmaterial.IconButton
}

func NewWalletSettingsPage(l *load.Load, wal *dcrlibwallet.Wallet) *WalletSettingsPage {
	pg := &WalletSettingsPage{
		Load:         l,
		wallet:       wal,
		changePass:   l.Theme.NewClickable(false),
		rescan:       l.Theme.NewClickable(false),
		deleteWallet: l.Theme.NewClickable(false),

		chevronRightIcon: l.Icons.ChevronRight,
	}

	pg.changePass.Radius = decredmaterial.Radius(14)
	pg.rescan.Radius = decredmaterial.Radius(14)
	pg.deleteWallet.Radius = decredmaterial.Radius(14)

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
}

func (pg *WalletSettingsPage) ID() string {
	return WalletSettingsPageID
}

func (pg *WalletSettingsPage) OnResume() {

}

func (pg *WalletSettingsPage) Layout(gtx layout.Context) layout.Dimensions {

	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrSettings),
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if !pg.wallet.IsWatchingOnlyWallet() {
							return pg.changePassphrase()(gtx)
						}
						return layout.Dimensions{}
					}),
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
						gtx.Constraints.Min.X = gtx.Px(values.MarginPadding20)
						return pg.chevronRightIcon.Layout(gtx, pg.Theme.Color.LightGray)
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
						gtx.Constraints.Min.X = gtx.Px(values.MarginPadding20)
						return pg.chevronRightIcon.Layout(gtx, pg.Theme.Color.LightGray)
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
						gtx.Constraints.Min.X = gtx.Px(values.MarginPadding20)
						return pg.chevronRightIcon.Layout(gtx, pg.Theme.Color.LightGray)
					})
				}),
			)
		})
	}
}

func (pg *WalletSettingsPage) pageSections(gtx layout.Context, title string, clickable *decredmaterial.Clickable, body layout.Widget) layout.Dimensions {
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
			return clickable.Layout(gtx, func(gtx C) D {
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
							pg.Toast.NotifyError(values.String(values.StrNotConnected))
							return
						}
						pg.Toast.NotifyError(err.Error())
						return
					}
					msg := values.String(values.StrRescanProgressNotification)
					pg.Toast.Notify(msg)
				})

			pg.ShowModal(info)
		}()
		break
	}

	for pg.deleteWallet.Clicked() {
		modal.NewInfoModal(pg.Load).
			Title(values.String(values.StrRemoveWallet)).
			Body("Make sure to have the seed word backed up before removing the wallet").
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButtonStyle(pg.Load.Theme.Color.Surface, pg.Load.Theme.Color.Danger).
			PositiveButton(values.String(values.StrRemove), func() {
				modal.NewPasswordModal(pg.Load).
					Title(values.String(values.StrConfirmToRemove)).
					NegativeButton(values.String(values.StrCancel), func() {}).
					PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
						go func() {
							err := pg.WL.MultiWallet.DeleteWallet(pg.wallet.ID, []byte(password))
							if err != nil {
								pm.SetError(err.Error())
								pm.SetLoading(false)
								return
							}
							pm.Dismiss()
							pg.Toast.Notify("Wallet removed")
							pg.PopFragment()
						}()
						return false
					}).Show()

			}).Show()
		break
	}
}

func (pg *WalletSettingsPage) OnClose() {}
