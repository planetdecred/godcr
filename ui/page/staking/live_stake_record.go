package staking

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *Page) initLiveStakeWidget() *Page {
	pg.toTickets = pg.Theme.TextAndIconButton(values.String(values.StrSeeAll), pg.Theme.Icons.NavigationArrowForward)
	pg.toTickets.Color = pg.Theme.Color.Primary
	pg.toTickets.BackgroundColor = color.NRGBA{}

	pg.ticketsLive = pg.Theme.NewClickableList(layout.Vertical)

	return pg
}

func (pg *Page) stakeLiveSection(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
					title := pg.Theme.Label(values.TextSize14, values.String(values.StrLiveTickets))
					title.Color = pg.Theme.Color.GrayText2
					return pg.titleRow(gtx, title.Layout, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							pg.stakingCountIcon(pg.Theme.Icons.TicketUnminedIcon, pg.ticketOverview.Unmined),
							pg.stakingCountIcon(pg.Theme.Icons.TicketImmatureIcon, pg.ticketOverview.Immature),
							pg.stakingCountIcon(pg.Theme.Icons.TicketLiveIcon, pg.ticketOverview.Live),
							layout.Rigid(func(gtx C) D {
								if len(pg.liveTickets) > 0 {
									return pg.toTickets.Layout(gtx)
								}
								return D{}
							}),
						)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				if len(pg.liveTickets) == 0 {
					noLiveStake := pg.Theme.Label(values.TextSize16, values.String(values.StrNoActiveTickets))
					noLiveStake.Color = pg.Theme.Color.GrayText3
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Center.Layout(gtx, noLiveStake.Layout)
				}
				return pg.ticketsLive.Layout(gtx, len(pg.liveTickets), func(gtx C, index int) D {
					return ticketListLayout(gtx, pg.Load, pg.liveTickets[index], index, true)
				})
			}),
		)
	})
}

func (pg *Page) stakingCountIcon(icon *decredmaterial.Image, count int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		return layout.Inset{Right: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					label := pg.Theme.Label(values.TextSize16, values.String(values.StrLiveTickets)+":")
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

func (pg *Page) stakingRecordSection(gtx C) D {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					title := pg.Theme.Label(values.TextSize14, values.String(values.StrTicketRecord))
					title.Color = pg.Theme.Color.GrayText2

					if pg.ticketOverview.All == 0 {
						return pg.titleRow(gtx, title.Layout, func(gtx C) D { return D{} })
					}
					return pg.titleRow(gtx, title.Layout, pg.toTickets.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				wdgs := []layout.Widget{
					pg.ticketRecordIconCount(pg.Theme.Icons.TicketUnminedIcon, pg.ticketOverview.Unmined, values.String(values.StrUmined)),
					pg.ticketRecordIconCount(pg.Theme.Icons.TicketImmatureIcon, pg.ticketOverview.Immature, values.String(values.StrImmature)),
					pg.ticketRecordIconCount(pg.Theme.Icons.TicketLiveIcon, pg.ticketOverview.Live, values.String(values.StrLive)),
					pg.ticketRecordIconCount(pg.Theme.Icons.TicketVotedIcon, pg.ticketOverview.Voted, values.String(values.StrVoted)),
					pg.ticketRecordIconCount(pg.Theme.Icons.TicketExpiredIcon, pg.ticketOverview.Expired, values.String(values.StrExpired)),
					pg.ticketRecordIconCount(pg.Theme.Icons.TicketRevokedIcon, pg.ticketOverview.Revoked, values.String(values.StrRevoked)),
				}
				return decredmaterial.GridWrap{
					Axis:      layout.Horizontal,
					Alignment: layout.End,
				}.Layout(gtx, len(wdgs), func(gtx C, i int) D {
					return wdgs[i](gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return decredmaterial.LinearLayout{
					Width:       decredmaterial.MatchParent,
					Height:      decredmaterial.WrapContent,
					Background:  pg.Theme.Color.Success2,
					Padding:     layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding16},
					Border:      decredmaterial.Border{Radius: decredmaterial.Radius(8)},
					Direction:   layout.Center,
					Alignment:   layout.Middle,
					Orientation: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
							txt := pg.Theme.Label(values.TextSize14, values.String(values.StrRewardsEarned))
							txt.Color = pg.Theme.Color.Turquoise700
							return txt.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := pg.Theme.Icons.StakeyIcon
								return layout.Inset{Right: values.MarginPadding6}.Layout(gtx, ic.Layout24dp)
							}),
							layout.Rigid(func(gtx C) D {
								award := pg.Theme.Color.Text
								noAward := pg.Theme.Color.GrayText3
								if pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.DarkModeConfigKey, false) {
									award = pg.Theme.Color.Gray3
									noAward = pg.Theme.Color.Gray3
								}

								if pg.totalRewards == "0 DCR" {
									txt := pg.Theme.Label(values.TextSize16, values.String(values.StrNoReward))
									txt.Color = noAward
									return txt.Layout(gtx)
								}

								return components.LayoutBalanceColor(gtx, pg.Load, pg.totalRewards, award)
							}),
						)
					}),
				)
			}),
		)
	})
}

func (pg *Page) ticketRecordIconCount(icon *decredmaterial.Image, count int, status string) layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Dp(unit.Dp(110))
		gtx.Constraints.Max.X = gtx.Dp(unit.Dp(110))
		return layout.Inset{Bottom: values.MarginPadding16, Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return icon.Layout24dp(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								label := pg.Theme.Label(values.TextSize16, fmt.Sprintf("%d", count))
								return label.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								txt := pg.Theme.Label(values.TextSize12, status)
								txt.Color = pg.Theme.Color.GrayText2
								return txt.Layout(gtx)
							}),
						)
					})
				}),
			)
		})
	}
}
