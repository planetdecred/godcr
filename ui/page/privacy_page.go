package page

import (
	"fmt"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PagePrivacy = "Privacy"

type privacyPage struct {
	*load.Load
	wallet                               *dcrlibwallet.Wallet
	pageContainer                        layout.List
	toggleMixer, allowUnspendUnmixedAcct *widget.Bool
	toPrivacySetup                       decredmaterial.Button
	dangerZoneCollapsible                *decredmaterial.Collapsible
	acctMixerStatus                      *chan *wallet.AccountMixer

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton
}

func PrivacyPage(l *load.Load, wallet *dcrlibwallet.Wallet) *privacyPage {
	pg := &privacyPage{
		Load:                    l,
		wallet:                  wallet,
		pageContainer:           layout.List{Axis: layout.Vertical},
		toggleMixer:             new(widget.Bool),
		allowUnspendUnmixedAcct: new(widget.Bool),
		toPrivacySetup:          l.Theme.Button(new(widget.Clickable), "Set up mixer for this wallet"),
		dangerZoneCollapsible:   l.Theme.Collapsible(),
		acctMixerStatus:         &l.Receiver.AcctMixerStatus,
	}
	pg.toPrivacySetup.Background = l.Theme.Color.Primary
	pg.backButton, pg.infoButton = subpageHeaderButtons(l)

	return pg
}

func (pg *privacyPage) OnResume() {

}

func (pg *privacyPage) Layout(gtx layout.Context) layout.Dimensions {
	d := func(gtx C) D {
		sp := SubPage{
			Load:       pg.Load,
			title:      "StakeShuffle",
			walletName: pg.wallet.Name,
			backButton: pg.backButton,
			infoButton: pg.infoButton,
			back: func() {
				pg.ChangePage(WalletPageID)
			},
			infoTemplate: modal.PrivacyInfoTemplate,
			body: func(gtx layout.Context) layout.Dimensions {
				if pg.wallet.AccountMixerConfigIsSet() {
					widgets := []func(gtx C) D{
						func(gtx C) D {
							return pg.mixerInfoLayout(gtx)
						},
						pg.gutter,
						func(gtx C) D {
							return pg.mixerSettingsLayout(gtx)
						},
						pg.gutter,
						func(gtx C) D {
							return pg.dangerZoneLayout(gtx)
						},
					}
					return pg.pageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
						return widgets[i](gtx)
					})
				}
				return pg.privacyIntroLayout(gtx)
			},
		}
		return sp.Layout(gtx)
	}
	return uniformPadding(gtx, d)
}

func (pg *privacyPage) privacyIntroLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
										pg.Icons.TransactionFingerPrintIcon.Scale = 1.0
										return pg.Icons.TransactionFingerPrintIcon.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									pg.Icons.ArrowForwardIcon.Scale = 0.5
									return pg.Icons.ArrowForwardIcon.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									pg.Icons.MixerSmall.Scale = 1.0
									return layout.Inset{
										Left:  values.MarginPadding5,
										Right: values.MarginPadding5,
									}.Layout(gtx, pg.Icons.MixerSmall.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									pg.Icons.ArrowForwardIcon.Scale = 0.5
									return pg.Icons.ArrowForwardIcon.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									pg.Icons.TransactionIcon.Scale = 1.5
									return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, pg.Icons.TransactionIcon.Layout)
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							txt := pg.Theme.H6("How does CoinShuffle++ mixer enhance your privacy?")
							txt2 := pg.Theme.Body1("CoinShuffle++ mixer can mix your DCRs through CoinJoin transactions.")
							txt3 := pg.Theme.Body1("Using mixed DCRs protects you from exposing your financial activities to the public (e.g. how much you own, who pays you).")
							txt.Alignment, txt2.Alignment, txt3.Alignment = text.Middle, text.Middle, text.Middle

							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(txt.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, txt2.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, txt3.Layout)
								}),
							)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, pg.toPrivacySetup.Layout)
			}),
		)
	})
}

func (pg *privacyPage) mixerInfoStatusTextLayout(gtx layout.Context) layout.Dimensions {
	txt := pg.Theme.H6("Mixer")
	subtxt := pg.Theme.Body2("Ready to mix")
	subtxt.Color = pg.Theme.Color.Gray
	iconVisibility := false

	if pg.wallet.IsAccountMixerActive() {
		txt.Text = "Mixer is running..."
		subtxt.Text = "Keep this app opened"
		iconVisibility = true
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(txt.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if !iconVisibility {
						return layout.Dimensions{}
					}
					pg.Icons.AlertGray.Scale = 1.0
					return pg.Icons.AlertGray.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, subtxt.Layout)
				}),
			)
		}),
	)
}

func (pg *privacyPage) mixersubInfolayout(gtx layout.Context) layout.Dimensions {
	txt := pg.Theme.Body2("")

	if pg.wallet.IsAccountMixerActive() {
		txt = pg.Theme.Body2("The mixer will automatically stop when unmixed balance are fully mixed.")
		txt.Color = pg.Theme.Color.Gray
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(txt.Layout),
	)
}

func (pg *privacyPage) mixerInfoLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							ic := pg.Icons.MixerSmall
							ic.Scale = 1.0
							return ic.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return pg.mixerInfoStatusTextLayout(gtx)
							})
						}),
						layout.Rigid(material.Switch(pg.Theme.Base, pg.toggleMixer).Layout),
					)
				}),
				layout.Rigid(pg.gutter),
				layout.Rigid(func(gtx C) D {
					content := pg.Theme.Card()
					content.Color = pg.Theme.Color.LightGray
					return content.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
							var mixedBalance = "0.00"
							var unmixedBalance = "0.00"
							accounts, _ := pg.wallet.GetAccountsRaw()
							for _, acct := range accounts.Acc {
								if acct.Number == pg.wallet.MixedAccountNumber() {
									mixedBalance = dcrutil.Amount(acct.TotalBalance).String()
								} else if acct.Number == pg.wallet.UnmixedAccountNumber() {
									unmixedBalance = dcrutil.Amount(acct.TotalBalance).String()
								}
							}

							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											txt := pg.Theme.Label(values.TextSize14, "Unmixed balance")
											txt.Color = pg.Theme.Color.Gray
											return txt.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layoutBalance(gtx, pg.Load, unmixedBalance, true)
										}),
									)
								}),
								layout.Rigid(func(gtx C) D {
									pg.Icons.ArrowDownIcon.Scale = 1.0
									return layout.Center.Layout(gtx, pg.Icons.ArrowDownIcon.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											t := pg.Theme.Label(values.TextSize14, "Mixed balance")
											t.Color = pg.Theme.Color.Gray
											return t.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layoutBalance(gtx, pg.Load, mixedBalance, true)
										}),
									)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return pg.mixersubInfolayout(gtx)
						}),
					)
				}),
			)
		})
	})
}

func (pg *privacyPage) mixerSettingsLayout(gtx layout.Context) layout.Dimensions {
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

func (pg *privacyPage) shufflePortForCurrentNet() string {
	if pg.WL.Wallet.Net == "testnet3" {
		return dcrlibwallet.TestnetShufflePort
	}

	return dcrlibwallet.MainnetShufflePort
}

func (pg *privacyPage) dangerZoneLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return pg.dangerZoneCollapsible.Layout(gtx,
				func(gtx C) D {
					txt := pg.Theme.Label(values.MarginPadding15, "Danger Zone")
					txt.Color = pg.Theme.Color.Gray
					return txt.Layout(gtx)
				},
				func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(1, pg.Theme.Label(values.TextSize16, "Allow spending from unmixed accounts").Layout),
							layout.Rigid(material.Switch(pg.Theme.Base, pg.allowUnspendUnmixedAcct).Layout),
						)
					})
				},
			)
		})
	})
}

func (pg *privacyPage) gutter(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Dimensions{}
	})
}

func (pg *privacyPage) Handle() {
	if pg.toPrivacySetup.Button.Clicked() {
		go pg.showModalSetupMixerInfo()
	}

	if pg.toggleMixer.Changed() {
		if pg.toggleMixer.Value {
			go pg.showModalPasswordStartAccountMixer()
		} else {
			go pg.WL.MultiWallet.StopAccountMixer(pg.wallet.ID)
		}
	}

	select {
	case stt := <-*pg.acctMixerStatus:
		if stt.RunStatus == wallet.MixerStarted {
			log.Info("Start Successfully")
			// TODO: will use when notify implemented
			// common.notify("Start Successfully", true)
		} else {
			log.Info("Stop Successfully")
			// common.notify("Stop Successfully", true)
		}
	default:
	}
}

func (pg *privacyPage) showModalSetupMixerInfo() {
	info := modal.NewInfoModal(pg.Load).
		Title("Set up mixer by creating two needed accounts").
		Body("Each time you receive a payment, a new address is generated to protect your privacy.").
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Begin setup", func() {
			pg.showModalSetupMixerAcct()
		})
	pg.ShowModal(info)
}

func (pg *privacyPage) showModalSetupMixerAcct() {
	accounts, _ := pg.wallet.GetAccountsRaw()
	for _, acct := range accounts.Acc {
		if acct.Name == "mixed" || acct.Name == "unmixed" {
			alert := mustIcon(widget.NewIcon(icons.AlertError))
			alert.Color = pg.Theme.Color.DeepBlue

			info := modal.NewInfoModal(pg.Load).
				Icon(alert).
				Title("Account name is taken").
				Body("There are existing accounts named mixed or unmixed. Please change the name to something else for now. You can change them back after the setup.").
				PositiveButton("Go back & rename", func() {
					*pg.Page = WalletPageID
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

func (pg *privacyPage) showModalPasswordStartAccountMixer() {
	modal.NewPasswordModal(pg.Load).
		Title("Confirm to mix account").
		NegativeButton("Cancel", func() {}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {

				err := pg.WL.MultiWallet.StartAccountMixer(pg.wallet.ID, password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
				// common.notify("Start Successfully", true)
			}()

			return false
		}).Show()
}

func (pg *privacyPage) OnClose() {}
