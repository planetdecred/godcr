package staking

import (
	"fmt"
	"image/color"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *Page) initLiveStakeWidget() *Page {
	pg.toTickets = pg.Theme.TextAndIconButton("See All", pg.Icons.NavigationArrowForward)
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
					title := pg.Theme.Label(values.TextSize14, "Live Tickets")
					title.Color = pg.Theme.Color.GrayText2
					return pg.titleRow(gtx, title.Layout, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							pg.stakingCountIcon(pg.Icons.TicketUnminedIcon, pg.ticketOverview.Unmined),
							pg.stakingCountIcon(pg.Icons.TicketImmatureIcon, pg.ticketOverview.Immature),
							pg.stakingCountIcon(pg.Icons.TicketLiveIcon, pg.ticketOverview.Live),
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
					noLiveStake := pg.Theme.Label(values.TextSize16, "No active tickets.")
					noLiveStake.Color = pg.Theme.Color.GrayText3
					return noLiveStake.Layout(gtx)
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
		if count == 0 {
			return D{}
		}
		return layout.Inset{Right: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return icon.Layout16dp(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
						label := pg.Theme.Label(values.TextSize14, fmt.Sprintf("%d", count))
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
					title := pg.Theme.Label(values.TextSize14, "Ticket Record")
					title.Color = pg.Theme.Color.GrayText2

					if pg.ticketOverview.All == 0 {
						return pg.titleRow(gtx, title.Layout, func(gtx C) D { return D{} })
					}
					return pg.titleRow(gtx, title.Layout, pg.toTickets.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				wdgs := []layout.Widget{
					pg.ticketRecordIconCount(pg.Icons.TicketUnminedIcon, pg.ticketOverview.Unmined, "Unmined"),
					pg.ticketRecordIconCount(pg.Icons.TicketImmatureIcon, pg.ticketOverview.Immature, "Immature"),
					pg.ticketRecordIconCount(pg.Icons.TicketLiveIcon, pg.ticketOverview.Live, "Live"),
					pg.ticketRecordIconCount(pg.Icons.TicketVotedIcon, pg.ticketOverview.Voted, "Voted"),
					pg.ticketRecordIconCount(pg.Icons.TicketExpiredIcon, pg.ticketOverview.Expired, "Expired"),
					pg.ticketRecordIconCount(pg.Icons.TicketRevokedIcon, pg.ticketOverview.Revoked, "Revoked"),
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
							txt := pg.Theme.Label(values.TextSize14, "Rewards Earned")
							txt.Color = pg.Theme.Color.Turquoise700
							return txt.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := pg.Icons.StakeyIcon
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
									txt := pg.Theme.Label(values.TextSize16, "Stakey sees no rewards")
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
