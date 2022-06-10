package privacy

import (
	"context"
	"fmt"

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

	wallet                *dcrlibwallet.Wallet
	pageContainer         layout.List
	dangerZoneCollapsible *decredmaterial.Collapsible

	backButton              decredmaterial.IconButton
	infoButton              decredmaterial.IconButton
	toggleMixer             *decredmaterial.Switch
	allowUnspendUnmixedAcct *decredmaterial.Switch

	mixerCompleted bool
}

func NewAccountMixerPage(l *load.Load, wallet *dcrlibwallet.Wallet) *AccountMixerPage {
	pg := &AccountMixerPage{
		Load:                    l,
		GenericPageModal:        app.NewGenericPageModal(AccountMixerPageID),
		wallet:                  wallet,
		pageContainer:           layout.List{Axis: layout.Vertical},
		toggleMixer:             l.Theme.Switch(),
		allowUnspendUnmixedAcct: l.Theme.Switch(),
		dangerZoneCollapsible:   l.Theme.Collapsible(),
	}
	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *AccountMixerPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.listenForMixerNotifications()
	pg.toggleMixer.SetChecked(pg.wallet.IsAccountMixerActive())

	isSpendUnmixedFunds := pg.wallet.ReadBoolConfigValueForKey(load.SpendUnmixedFundsKey, false)
	pg.allowUnspendUnmixedAcct.SetChecked(isSpendUnmixedFunds)
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *AccountMixerPage) Layout(gtx layout.Context) layout.Dimensions {
	d := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "StakeShuffle",
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			InfoButton: pg.infoButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			InfoTemplate: modal.PrivacyInfoTemplate,
			Body: func(gtx layout.Context) layout.Dimensions {
				widgets := []func(gtx C) D{
					func(gtx C) D {
						return components.MixerInfoLayout(gtx, pg.Load, pg.wallet.IsAccountMixerActive(),
							pg.toggleMixer.Layout, func(gtx C) D {
								mixedBalance := "0.00"
								unmixedBalance := "0.00"
								accounts, _ := pg.wallet.GetAccountsRaw()
								for _, acct := range accounts.Acc {
									if acct.Number == pg.wallet.MixedAccountNumber() {
										mixedBalance = dcrutil.Amount(acct.TotalBalance).String()
									} else if acct.Number == pg.wallet.UnmixedAccountNumber() {
										unmixedBalance = dcrutil.Amount(acct.TotalBalance).String()
									}
								}

								return components.MixerInfoContentWrapper(gtx, pg.Load, func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													txt := pg.Theme.Label(values.TextSize14, "Unmixed balance")
													txt.Color = pg.Theme.Color.GrayText2
													return txt.Layout(gtx)
												}),
												layout.Rigid(func(gtx C) D {
													return components.LayoutBalance(gtx, pg.Load, unmixedBalance)
												}),
											)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Center.Layout(gtx, pg.Theme.Icons.ArrowDownIcon.Layout24dp)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													t := pg.Theme.Label(values.TextSize14, "Mixed balance")
													t.Color = pg.Theme.Color.GrayText2
													return t.Layout(gtx)
												}),
												layout.Rigid(func(gtx C) D {
													return components.LayoutBalance(gtx, pg.Load, mixedBalance)
												}),
											)
										}),
									)
								})
							})
					},
					func(gtx C) D {
						return pg.mixerSettingsLayout(gtx)
					},
					func(gtx C) D {
						return pg.dangerZoneLayout(gtx)
					},
				}
				return pg.pageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
					m := values.MarginPadding10
					if i == len(widgets) {
						m = values.MarginPadding0
					}
					return layout.Inset{Bottom: m}.Layout(gtx, widgets[i])
				})

			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	return components.UniformPadding(gtx, d)
}

func (pg *AccountMixerPage) mixerSettingsLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		mixedAccountName, _ := pg.wallet.AccountName(pg.wallet.MixedAccountNumber())
		unmixedAccountName, _ := pg.wallet.AccountName(pg.wallet.UnmixedAccountNumber())

		row := func(txt1, txt2 string) D {
			return layout.Inset{
				Left:   values.MarginPadding15,
				Right:  values.MarginPadding15,
				Top:    values.MarginPadding10,
				Bottom: values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(pg.Theme.Label(values.TextSize16, txt1).Layout),
					layout.Rigid(pg.Theme.Body2(txt2).Layout),
				)
			})
		}

		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, pg.Theme.Body2("Mixer Settings").Layout)
			}),
			layout.Rigid(func(gtx C) D { return row("Mixed account", mixedAccountName) }),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Rigid(func(gtx C) D { return row("Change account", unmixedAccountName) }),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Rigid(func(gtx C) D { return row("Account branch", fmt.Sprintf("%d", dcrlibwallet.MixedAccountBranch)) }),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Rigid(func(gtx C) D { return row("Shuffle server", dcrlibwallet.ShuffleServer) }),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Rigid(func(gtx C) D { return row("Shuffle port", pg.shufflePortForCurrentNet()) }),
		)
	})
}

func (pg *AccountMixerPage) shufflePortForCurrentNet() string {
	if pg.WL.Wallet.Net == dcrlibwallet.Testnet3 {
		return dcrlibwallet.TestnetShufflePort
	}

	return dcrlibwallet.MainnetShufflePort
}

func (pg *AccountMixerPage) dangerZoneLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return pg.dangerZoneCollapsible.Layout(gtx,
				func(gtx C) D {
					txt := pg.Theme.Label(values.TextSize16, "Danger Zone")
					txt.Color = pg.Theme.Color.Danger
					return txt.Layout(gtx)
				},
				func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(1, pg.Theme.Label(values.TextSize16, "Allow spending from unmixed accounts").Layout),
							layout.Rigid(pg.allowUnspendUnmixedAcct.Layout),
						)
					})
				},
			)
		})
	})
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
				Title("Cancel mixer?").
				Body("Are you sure you want to cancel mixer action?").
				NegativeButton("No", func() {}).
				PositiveButton("Yes", func(isChecked bool) bool {
					pg.toggleMixer.SetChecked(false)
					go pg.WL.MultiWallet.StopAccountMixer(pg.wallet.ID)
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

	if pg.allowUnspendUnmixedAcct.Changed() {
		if pg.allowUnspendUnmixedAcct.IsChecked() {
			textModal := modal.NewTextInputModal(pg.Load).
				SetTextWithTemplate(modal.AllowUnmixedSpendingTemplate).
				Hint("").
				PositiveButtonStyle(pg.Load.Theme.Color.Danger, pg.Load.Theme.Color.InvText).
				PositiveButton("Confirm", func(textInput string, tim *modal.TextInputModal) bool {
					if textInput != "I understand the risks" {
						tim.SetError("confirmation text is incorrect")
						tim.SetLoading(false)
					} else {
						pg.wallet.SetBoolConfigValueForKey(load.SpendUnmixedFundsKey, true)
						tim.Dismiss()
					}
					return false
				})

			textModal.Title("Confirm to allow spending from unmixed accounts").
				NegativeButton("Cancel", func() {
					pg.allowUnspendUnmixedAcct.SetChecked(false)
				})
			pg.ParentWindow().ShowModal(textModal)

		} else {
			pg.wallet.SetBoolConfigValueForKey(load.SpendUnmixedFundsKey, false)
		}

		if pg.dangerZoneCollapsible.IsExpanded() {
			pg.ParentWindow().Reload()
		}
	}

	if pg.backButton.Button.Clicked() {
		pg.ParentNavigator().ClosePagesAfter(components.WalletsPageID)
	}
}

func (pg *AccountMixerPage) showModalPasswordStartAccountMixer() {
	passwordModal := modal.NewPasswordModal(pg.Load).
		Title("Confirm to mix account").
		NegativeButton("Cancel", func() {
			pg.toggleMixer.SetChecked(false)
		}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := pg.WL.MultiWallet.StartAccountMixer(pg.wallet.ID, password)
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
					pg.Toast.Notify("Mixer start Successfully")
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
