package wallets

import (
	"context"
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"github.com/planetdecred/godcr/listeners"
)

const PrivacyPageID = "Privacy"

type PrivacyPage struct {
	*load.Load

	*listeners.AccountMixerNotif

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	wallet                *dcrlibwallet.Wallet
	pageContainer         layout.List
	toPrivacySetup        decredmaterial.Button
	dangerZoneCollapsible *decredmaterial.Collapsible

	backButton              decredmaterial.IconButton
	infoButton              decredmaterial.IconButton
	toggleMixer             *decredmaterial.Switch
	allowUnspendUnmixedAcct *decredmaterial.Switch

	mixerCompleted bool
}

func NewPrivacyPage(l *load.Load, wallet *dcrlibwallet.Wallet) *PrivacyPage {
	pg := &PrivacyPage{
		Load:                    l,
		wallet:                  wallet,
		pageContainer:           layout.List{Axis: layout.Vertical},
		toggleMixer:             l.Theme.Switch(),
		allowUnspendUnmixedAcct: l.Theme.Switch(),
		toPrivacySetup:          l.Theme.Button("Set up mixer for this wallet"),
		dangerZoneCollapsible:   l.Theme.Collapsible(),
	}
	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *PrivacyPage) ID() string {
	return PrivacyPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *PrivacyPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.listenForMixerNotifications()
	pg.toggleMixer.SetChecked(pg.wallet.IsAccountMixerActive())

	if pg.wallet.AccountMixerConfigIsSet() {
		pg.allowUnspendUnmixedAcct.SetChecked(false)
	} else {
		pg.allowUnspendUnmixedAcct.SetChecked(true)
	}

	if pg.AccountMixerNotif == nil {
		pg.AccountMixerNotif = listeners.NewAccountMixerNotif(make(chan wallet.AccountMixer, 4))
	} else {
		pg.MixerCh = make(chan wallet.AccountMixer, 4)
	}
	pg.WL.MultiWallet.AddAccountMixerNotificationListener(pg, PrivacyPageID)
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *PrivacyPage) Layout(gtx layout.Context) layout.Dimensions {
	d := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "StakeShuffle",
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			InfoButton: pg.infoButton,
			Back: func() {
				pg.PopFragment()
			},
			InfoTemplate: modal.PrivacyInfoTemplate,
			Body: func(gtx layout.Context) layout.Dimensions {
				if pg.wallet.MixedAccountNumber() > 0 && pg.wallet.UnmixedAccountNumber() > 0 {
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
												return layout.Center.Layout(gtx, pg.Icons.ArrowDownIcon.Layout24dp)
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
				}
				return pg.privacyIntroLayout(gtx)
			},
		}
		return sp.Layout(gtx)
	}
	return components.UniformPadding(gtx, d)
}

func (pg *PrivacyPage) privacyIntroLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Bottom: values.MarginPadding24,
							}.Layout(gtx, func(gtx C) D {
								return pg.Icons.PrivacySetup.LayoutSize(gtx, values.MarginPadding280)
							})
						}),
						layout.Rigid(func(gtx C) D {
							txt := pg.Theme.H6("How does StakeShuffle++ mixer enhance your privacy?")
							txt2 := pg.Theme.Body1("Shuffle++ mixer can mix your DCRs through CoinJoin transactions.")
							txt3 := pg.Theme.Body1("Using mixed DCRs protects you from exposing your financial activities to")
							txt4 := pg.Theme.Body1("the public (e.g. how much you own, who pays you).")
							txt.Alignment, txt2.Alignment, txt3.Alignment, txt4.Alignment = text.Middle, text.Middle, text.Middle, text.Middle

							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(txt.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, txt2.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, txt3.Layout)
								}),
								layout.Rigid(txt4.Layout),
							)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, pg.toPrivacySetup.Layout)
			}),
		)
	})
}

func (pg *PrivacyPage) mixerSettingsLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

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
			layout.Rigid(func(gtx C) D { return row("Mixed account", "mixed") }),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Rigid(func(gtx C) D { return row("Change account", "unmixed") }),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Rigid(func(gtx C) D { return row("Account branch", fmt.Sprintf("%d", dcrlibwallet.MixedAccountBranch)) }),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Rigid(func(gtx C) D { return row("Shuffle server", dcrlibwallet.ShuffleServer) }),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Rigid(func(gtx C) D { return row("Shuffle port", pg.shufflePortForCurrentNet()) }),
		)
	})
}

func (pg *PrivacyPage) shufflePortForCurrentNet() string {
	if pg.WL.Wallet.Net == dcrlibwallet.Testnet3 {
		return dcrlibwallet.TestnetShufflePort
	}

	return dcrlibwallet.MainnetShufflePort
}

func (pg *PrivacyPage) dangerZoneLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return pg.dangerZoneCollapsible.Layout(gtx,
				func(gtx C) D {
					txt := pg.Theme.Label(values.MarginPadding15, "Danger Zone")
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
func (pg *PrivacyPage) HandleUserInteractions() {
	if pg.toPrivacySetup.Clicked() {
		go pg.showModalSetupMixerInfo()
	}

	if pg.toggleMixer.Changed() {
		if pg.toggleMixer.IsChecked() {
			go pg.showModalPasswordStartAccountMixer()
		} else {
			pg.toggleMixer.SetChecked(true)
			info := modal.NewInfoModal(pg.Load).
				Title("Cancel mixer?").
				Body("Are you sure you want to cancel mixer action?").
				NegativeButton("No", func() {}).
				PositiveButton("Yes", func() {
					pg.toggleMixer.SetChecked(false)
					go pg.WL.MultiWallet.StopAccountMixer(pg.wallet.ID)
				})
			pg.ShowModal(info)
		}
	}

	if pg.mixerCompleted {
		pg.toggleMixer.SetChecked(false)
		pg.mixerCompleted = false
		pg.RefreshWindow()
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
						pg.wallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false)
						tim.Dismiss()
					}
					return false
				})

			textModal.Title("Confirm to allow spending from unmixed accounts").
				NegativeButton("Cancel", func() {
					pg.allowUnspendUnmixedAcct.SetChecked(false)
				})
			textModal.Show()

		} else {
			pg.wallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, true)
		}

		if pg.dangerZoneCollapsible.IsExpanded() {
			pg.RefreshWindow()
		}
	}

}

func (pg *PrivacyPage) showModalSetupMixerInfo() {
	info := modal.NewInfoModal(pg.Load).
		Title("Set up mixer by creating two needed accounts").
		Body("Each time you receive a payment, a new address is generated to protect your privacy.").
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Begin setup", func() {
			pg.showModalSetupMixerAcct()
		})
	pg.ShowModal(info)
}

func (pg *PrivacyPage) showModalSetupMixerAcct() {
	accounts, _ := pg.wallet.GetAccountsRaw()
	for _, acct := range accounts.Acc {
		if acct.Name == "mixed" || acct.Name == "unmixed" {
			alert := decredmaterial.NewIcon(decredmaterial.MustIcon(widget.NewIcon(icons.AlertError)))
			alert.Color = pg.Theme.Color.DeepBlue
			info := modal.NewInfoModal(pg.Load).
				Icon(alert).
				Title("Account name is taken").
				Body("There are existing accounts named mixed or unmixed. Please change the name to something else for now. You can change them back after the setup.").
				PositiveButton("Go back & rename", func() {
					pg.PopFragment()
				})
			pg.ShowModal(info)
			return
		}
	}

	modal.NewPasswordModal(pg.Load).
		Title("Confirm to create needed accounts").
		NegativeButton("Cancel", func() {}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := pg.wallet.CreateMixerAccounts("mixed", "unmixed", password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
			}()

			return false
		}).Show()
}

func (pg *PrivacyPage) showModalPasswordStartAccountMixer() {
	modal.NewPasswordModal(pg.Load).
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
		}).Show()
}

func (pg *PrivacyPage) listenForMixerNotifications() {
	go func() {
		for {
			select {
			case n := <-pg.MixerCh:
				if n.RunStatus == wallet.MixerStarted {
					pg.Toast.Notify("Mixer start Successfully")
					pg.RefreshWindow()
				}

				if n.RunStatus == wallet.MixerEnded {
					pg.mixerCompleted = true
					pg.RefreshWindow()
				}

			case <-pg.ctx.Done():
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
func (pg *PrivacyPage) OnNavigatedFrom() {
	if pg.MixerCh != nil {
		close(pg.MixerCh)
	}
	pg.WL.MultiWallet.RemoveAccountMixerNotificationListener(PrivacyPageID)
}
