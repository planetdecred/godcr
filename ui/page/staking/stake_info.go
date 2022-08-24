package staking

import (
	"fmt"
	"image/color"

	"gioui.org/layout"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *Page) initStakePriceWidget() *Page {
	pg.stakeSettings = pg.Theme.NewClickable(false)
	_, pg.infoButton = components.SubpageHeaderButtons(pg.Load)

	pg.stake = pg.Theme.Switch()
	return pg
}

func (pg *Page) stakePriceSection(gtx C) D {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding11,
				}.Layout(gtx, func(gtx C) D {
					col := pg.Theme.Color.GrayText2
					leftWg := func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										title := pg.Theme.Label(values.TextSize16, values.String(values.StrTicketPrice)+": ")
										title.Color = col
										return title.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Center.Layout(gtx, func(gtx C) D {
											if pg.WL.MultiWallet.IsSyncing() {
												title := pg.Theme.Label(values.TextSize16, values.String(values.StrLoadingPrice))
												title.Color = col
												return title.Layout(gtx)
											}

											return components.LayoutBalanceSize(gtx, pg.Load, pg.ticketPrice, values.TextSize16)
										})
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left:  values.MarginPadding8,
											Right: values.MarginPadding4,
										}.Layout(gtx, pg.Theme.Icons.TimerIcon.Layout12dp)
									}),
									layout.Rigid(func(gtx C) D {
										secs, _ := pg.WL.MultiWallet.NextTicketPriceRemaining()
										txt := pg.Theme.Label(values.TextSize16, nextTicketRemaining(int(secs)))
										txt.Color = col

										if pg.WL.MultiWallet.IsSyncing() {
											txt.Text = values.String(values.StrSyncingState)
										}
										return txt.Layout(gtx)
									}),
								)
							}),
							pg.dataRows(values.String(values.StrLiveTickets), pg.ticketOverview.Live),
							pg.dataRows(values.String(values.StrCanBuy), pg.CalculateTotalTicketsCanBuy()),
						)
					}

					rightWg := func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								title := pg.Theme.Label(values.TextSize16, values.String(values.StrStake))
								title.Color = col
								if !pg.WL.SelectedWallet.Wallet.IsWatchingOnlyWallet() {
									return title.Layout(gtx)
								}
								return D{}
							}),
							layout.Rigid(func(gtx C) D {
								if !pg.WL.SelectedWallet.Wallet.IsWatchingOnlyWallet() {
									return layout.Inset{
										Right: values.MarginPadding40,
										Left:  values.MarginPadding4,
									}.Layout(gtx, pg.stake.Layout)
								}
								return D{}
							}),
							layout.Rigid(func(gtx C) D {
								icon := pg.Theme.Icons.HeaderSettingsIcon
								// Todo -- darkmode icons
								// if pg.ticketBuyerWallet.IsAutoTicketsPurchaseActive() {
								// 	icon = pg.Theme.Icons.SettingsInactiveIcon
								// }
								if !pg.WL.SelectedWallet.Wallet.IsWatchingOnlyWallet() {
									return pg.stakeSettings.Layout(gtx, icon.Layout24dp)
								}
								return D{}
							}),
							layout.Rigid(func(gtx C) D {
								pg.infoButton.Inset = layout.UniformInset(values.MarginPadding0)
								pg.infoButton.Size = values.MarginPadding22
								return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, pg.infoButton.Layout)
							}),
						)
					}

					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(leftWg),
						layout.Rigid(rightWg),
					)
				})
			}),
			layout.Rigid(pg.balanceProgressBarLayout),
		)
	})
}

func (pg *Page) dataRows(title string, count int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					label := pg.Theme.Label(values.TextSize16, title+":")
					label.Color = pg.Theme.Color.GrayText2
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
						label := pg.Theme.Label(values.TextSize16, fmt.Sprintf("%d", count))
						label.Color = pg.Theme.Color.GrayText2
						return label.Layout(gtx)
					})
				}),
			)
		})
	})
}

func (pg *Page) CalculateTotalTicketsCanBuy() int {
	totalBalance, _ := components.CalculateTotalWalletsBalance(pg.Load)
	ticketPrice, err := pg.WL.MultiWallet.TicketPrice()
	if err != nil {
		log.Errorf("ticketPrice error:", err)
		return 0
	}
	canBuy := totalBalance.Spendable.ToCoin() / dcrutil.Amount(ticketPrice.TicketPrice).ToCoin()
	if canBuy < 0 {
		canBuy = 0
	}

	return int(canBuy)
}

func (pg *Page) balanceProgressBarLayout(gtx C) D {
	totalBalance, _ := components.CalculateTotalWalletsBalance(pg.Load)

	items := []decredmaterial.ProgressBarItem{
		{
			Value: int(totalBalance.LockedByTickets.ToCoin()),
			Color: pg.Theme.Color.NavyBlue,
		},
		{
			Value: int(totalBalance.Spendable.ToCoin()),
			Color: pg.Theme.Color.Turquoise300,
		},
	}

	labelWdg := func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.layoutIconAndText(gtx, "Staked"+": ", totalBalance.LockedByTickets.String(), items[0].Color)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.layoutIconAndText(gtx, values.String(values.StrLabelSpendable)+": ", totalBalance.Spendable.String(), items[1].Color)
				}),
			)
		})
	}
	pb := pg.Theme.MultiLayerProgressBar(int((totalBalance.Spendable + totalBalance.LockedByTickets).ToCoin()), items)
	pb.Height = values.MarginPadding16
	return pb.Layout(gtx, labelWdg)

}

func (pg *Page) layoutIconAndText(gtx C, title string, val string, col color.NRGBA) D {
	return layout.Inset{Right: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					ic := decredmaterial.NewIcon(pg.Theme.Icons.ImageBrightness1)
					ic.Color = col
					return ic.Layout(gtx, values.MarginPadding8)
				})
			}),
			layout.Rigid(func(gtx C) D {
				txt := pg.Theme.Label(values.TextSize14, title)
				txt.Color = pg.Theme.Color.GrayText2
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := pg.Theme.Label(values.TextSize14, val)
				txt.Color = pg.Theme.Color.GrayText2
				return txt.Layout(gtx)
			}),
		)
	})
}

func (pg *Page) stakingRecordStatistics(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(pg.stakingRecord(pg.totalRewards, fmt.Sprintf("%s %s", values.String(values.StrTotal), values.String(values.StrReward)))),
		layout.Rigid(pg.stakingRecord(fmt.Sprintf("%d", pg.ticketOverview.Voted), values.String(values.StrVoted))),
		layout.Rigid(pg.stakingRecord(fmt.Sprintf("%d", pg.ticketOverview.Revoked), values.String(values.StrRevoked))),
		layout.Rigid(pg.stakingRecord(fmt.Sprintf("%d", pg.ticketOverview.Immature), values.String(values.StrImmature))),
		layout.Rigid(pg.stakingRecord(fmt.Sprintf("%d", pg.ticketOverview.Unmined), values.String(values.StrUmined))),
		layout.Rigid(pg.stakingRecord(fmt.Sprintf("%d", pg.ticketOverview.Expired), values.String(values.StrExpired))),
	)
}

func (pg *Page) stakingRecord(count, status string) layout.Widget {
	return func(gtx C) D {
		return components.EndToEndRow(gtx,
			pg.Theme.Label(values.TextSize14, status).Layout,
			pg.Theme.Label(values.TextSize14, count).Layout,
		)
	}
}
