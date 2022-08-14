package privacy

import (
	"context"

	"gioui.org/layout"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const AccountMixerPageID = "AccountMixer"

type AccountMixerPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	*listeners.AccountMixerNotificationListener

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	wallet *dcrlibwallet.Wallet

	toggleMixer *decredmaterial.Switch

	mixerCompleted bool

	totalBalance  dcrutil.Amount
	mixerProgress decredmaterial.ProgressBarStyle
}

func NewAccountMixerPage(l *load.Load) *AccountMixerPage {
	pg := &AccountMixerPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(AccountMixerPageID),
		wallet:           l.WL.SelectedWallet.Wallet,
		toggleMixer:      l.Theme.Switch(),
		mixerProgress:    l.Theme.ProgressBar(0),
	}
	pg.mixerProgress.Height = values.MarginPadding18
	pg.mixerProgress.Radius = decredmaterial.Radius(2)
	totalBalance, _ := components.CalculateTotalWalletsBalance(pg.Load) // TODO - handle error
	pg.totalBalance = totalBalance.Total

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *AccountMixerPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.listenForMixerNotifications()
	pg.toggleMixer.SetChecked(pg.WL.SelectedWallet.Wallet.IsAccountMixerActive())
}

func (pg *AccountMixerPage) toggleMixerAndProgres(gtx C, l *load.Load, button layout.Widget) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.UniformInset(values.MarginPadding40).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(l.Theme.H6(values.String(values.StrBalance)).Layout),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									return components.LayoutBalanceWithUnit(gtx, pg.Load, pg.totalBalance.String())
								})
							}),
							layout.Flexed(1, func(gtx C) D {
								return layout.E.Layout(gtx, func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, l.Theme.H6(values.String(values.StrMix)).Layout)
										}),
										layout.Rigid(button),
									)
								})
							}),
						)
					})
				}),
				layout.Rigid(l.Theme.Separator().Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := l.Theme.H6(values.String(values.StrMixer))
								txt.Color = l.Theme.Color.GrayText3
								return txt.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding20, Right: values.MarginPadding40}.Layout(gtx, pg.mixerProgress.Layout)
							}),
						)
					})
				}),
			)
		})
	})
}

func (pg *AccountMixerPage) mixedBalanceInfo(gtx C, l *load.Load, mixedBalance string) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		insetValue := values.MarginPadding15
		return layout.Inset{Right: insetValue,
			Left: insetValue,
			Top:  insetValue,
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(l.Theme.Icons.MixedTxIcon.Layout12dp),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding11}.Layout(gtx, l.Theme.H6(values.String(values.StrMixed)).Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding25}.Layout(gtx, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return components.LayoutBalanceWithUnit(gtx, pg.Load, mixedBalance)
						})
					})
				}),
			)
		})
	})
}

func (pg *AccountMixerPage) mixerImage(gtx C, l *load.Load) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		insetValue := values.MarginPadding40
		return layout.Inset{Right: insetValue, Left: insetValue}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Max.X = gtx.Constraints.Max.X/2 - 40
					return layout.W.Layout(gtx, l.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding20, Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, l.Theme.Icons.MixerIcon.Layout36dp)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.E.Layout(gtx, l.Theme.Separator().Layout)
				}),
			)
		})
	})
}

func (pg *AccountMixerPage) unmixedBalanceInfo(gtx C, l *load.Load, unmixedBalance string) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.Inset{Right: values.MarginPadding15, Left: values.MarginPadding15, Bottom: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(l.Theme.Icons.UnmixedTxIcon.Layout12dp),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding11}.Layout(gtx, l.Theme.H6(values.String(values.StrUnmixed)).Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding25}.Layout(gtx, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return components.LayoutBalanceWithUnit(gtx, pg.Load, unmixedBalance)
						})
					})
				}),
			)
		})
	})
}

func (pg *AccountMixerPage) LayoutMixerPage(gtx C, l *load.Load, mixerActive bool, button layout.Widget) D {
	mixedBalance := "0 DCR"
	unmixedBalance := "0 DCR"
	accounts, _ := pg.wallet.GetAccountsRaw() // TODO - handle errors

	for _, acct := range accounts.Acc {
		if acct.Number == pg.wallet.MixedAccountNumber() {
			mixedBalance = dcrutil.Amount(acct.TotalBalance).String()
		} else if acct.Number == pg.wallet.UnmixedAccountNumber() {
			unmixedBalance = dcrutil.Amount(acct.TotalBalance).String()
		}
	}

	return l.Theme.Card().Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			pg.toggleMixerAndProgres(gtx, l, button),
			pg.mixedBalanceInfo(gtx, l, mixedBalance),
			pg.mixerImage(gtx, l),
			pg.unmixedBalanceInfo(gtx, l, unmixedBalance),
		)
	})
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *AccountMixerPage) Layout(gtx layout.Context) layout.Dimensions {
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *AccountMixerPage) layoutDesktop(gtx layout.Context) layout.Dimensions {
	return components.UniformPadding(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding50).Layout(gtx, func(gtx C) D {
			return pg.LayoutMixerPage(gtx, pg.Load, pg.wallet.IsAccountMixerActive(), pg.toggleMixer.Layout)
		})
	})
}

func (pg *AccountMixerPage) layoutMobile(gtx layout.Context) layout.Dimensions {
	return D{}
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *AccountMixerPage) HandleUserInteractions() {
	if pg.toggleMixer.Changed() {
		if pg.toggleMixer.IsChecked() {
			go pg.showModalPasswordStartAccountMixer()
		} else {
			pg.toggleMixer.SetChecked(true)
			info := modal.NewInfoModal(pg.Load).
				Title(values.String(values.StrCancelMixer)).
				Body(values.String(values.StrSureToCancelMixer)).
				NegativeButton(values.String(values.StrNo), func() {}).
				PositiveButton(values.String(values.StrYes), func(isChecked bool) bool {
					pg.toggleMixer.SetChecked(false)
					go pg.WL.MultiWallet.StopAccountMixer(pg.WL.SelectedWallet.Wallet.ID)
					return true
				})
			pg.ParentWindow().ShowModal(info)
		}
	}

	if pg.mixerCompleted {
		pg.toggleMixer.SetChecked(false)
		pg.mixerCompleted = false
		pg.ParentWindow().Reload()
	}
}

func (pg *AccountMixerPage) showModalPasswordStartAccountMixer() {
	passwordModal := modal.NewPasswordModal(pg.Load).
		Title(values.String(values.StrConfirmToMixAccount)).
		NegativeButton(values.String(values.StrCancel), func() {
			pg.toggleMixer.SetChecked(false)
		}).
		PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := pg.WL.MultiWallet.StartAccountMixer(pg.WL.SelectedWallet.Wallet.ID, password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
			}()

			return false
		})
	pg.ParentWindow().ShowModal(passwordModal)
}

func (pg *AccountMixerPage) listenForMixerNotifications() {
	if pg.AccountMixerNotificationListener != nil {
		return
	}

	pg.AccountMixerNotificationListener = listeners.NewAccountMixerNotificationListener()
	err := pg.WL.MultiWallet.AddAccountMixerNotificationListener(pg, AccountMixerPageID)
	if err != nil {
		log.Errorf("Error adding account mixer notification listener: %+v", err)
		return
	}

	go func() {
		for {
			select {
			case n := <-pg.MixerChan:
				if n.RunStatus == wallet.MixerStarted {
					pg.Toast.Notify(values.String(values.StrMixerStart))
					pg.ParentWindow().Reload()
				}

				if n.RunStatus == wallet.MixerEnded {
					pg.mixerCompleted = true
					pg.ParentWindow().Reload()
				}

			case <-pg.ctx.Done():
				pg.WL.MultiWallet.RemoveAccountMixerNotificationListener(AccountMixerPageID)
				close(pg.MixerChan)
				pg.AccountMixerNotificationListener = nil
				return
			}
		}
	}()
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *AccountMixerPage) OnNavigatedFrom() {
	pg.ctxCancel()
}
