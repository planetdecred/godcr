package tickets

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const OverviewPageID = "Tickets"

type Page struct {
	*load.Load

	ticketPageContainer *layout.List
	ticketsLive         *layout.List
	ticketsActivity     *layout.List

	purchaseTicket decredmaterial.Button

	tickets     **wallet.Tickets
	ticketPrice string

	autoPurchaseEnabled *decredmaterial.Switch
	toTickets           decredmaterial.TextAndIconButton
	toTicketsActivity   decredmaterial.TextAndIconButton
	ticketTooltips      []tooltips
}

func NewTicketPage(l *load.Load) *Page {
	pg := &Page{
		Load:    l,
		tickets: l.WL.Tickets,

		ticketsLive:         &layout.List{Axis: layout.Horizontal},
		ticketsActivity:     &layout.List{Axis: layout.Vertical},
		ticketPageContainer: &layout.List{Axis: layout.Vertical},
		purchaseTicket:      l.Theme.Button(new(widget.Clickable), "Purchase"),

		autoPurchaseEnabled: l.Theme.Switch(),
		toTickets:           l.Theme.TextAndIconButton(new(widget.Clickable), "See All", l.Icons.NavigationArrowForward),
		toTicketsActivity:   l.Theme.TextAndIconButton(new(widget.Clickable), "See All", l.Icons.NavigationArrowForward),
	}

	pg.purchaseTicket.TextSize = values.TextSize12
	pg.purchaseTicket.Background = l.Theme.Color.Primary

	pg.toTickets.Color = l.Theme.Color.Primary
	pg.toTickets.BackgroundColor = l.Theme.Color.Surface

	pg.toTicketsActivity.Color = l.Theme.Color.Primary
	pg.toTicketsActivity.BackgroundColor = l.Theme.Color.Surface

	return pg
}

func (pg *Page) ID() string {
	return OverviewPageID
}

func (pg *Page) OnResume() {
	pg.ticketPrice = dcrutil.Amount(pg.WL.TicketPrice()).String()
	go pg.WL.GetVSPList()
	// TODO: automatic ticket purchase functionality
	pg.autoPurchaseEnabled.Disabled()
}

func (pg *Page) Layout(gtx layout.Context) layout.Dimensions {
	return components.UniformPadding(gtx, func(gtx layout.Context) layout.Dimensions {
		sections := []func(gtx C) D{
			func(ctx layout.Context) layout.Dimensions {
				return pg.ticketPriceSection(gtx)
			},
			func(ctx layout.Context) layout.Dimensions {
				return pg.ticketsLiveSection(gtx)
			},
			func(ctx layout.Context) layout.Dimensions {
				return pg.ticketsActivitySection(gtx)
			},
			func(ctx layout.Context) layout.Dimensions {
				return pg.stackingRecordSection(gtx)
			},
		}

		return pg.ticketPageContainer.Layout(gtx, len(sections), func(gtx C, i int) D {
			return sections[i](gtx)
		})
	})
}

func (pg *Page) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.Inset{
		Bottom: values.MarginPadding8,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, body)
		})
	})
}

func (pg *Page) titleRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(leftWidget),
		layout.Rigid(rightWidget),
	)
}

func (pg *Page) ticketPriceSection(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding11,
				}.Layout(gtx, func(gtx C) D {
					leftWg := func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								title := pg.Theme.Label(values.TextSize14, "Ticket Price")
								title.Color = pg.Theme.Color.Gray2
								return title.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Left:  values.MarginPadding8,
									Right: values.MarginPadding4,
								}.Layout(gtx, func(gtx C) D {
									ic := pg.Icons.TimerIcon
									ic.Scale = 1
									return ic.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								secs, _ := pg.WL.MultiWallet.NextTicketPriceRemaining()
								txt := pg.Theme.Label(values.TextSize14, nextTicketRemaining(int(secs)))
								txt.Color = pg.Theme.Color.Gray2
								return txt.Layout(gtx)
							}),
						)
					}
					return pg.titleRow(gtx, leftWg, pg.autoPurchaseEnabled.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					ic := pg.Icons.TicketPurchasedIcon
					return layout.Center.Layout(gtx, ic.Layout48dp)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return components.LayoutBalanceSize(gtx, pg.Load, pg.ticketPrice, values.TextSize28)
					})
				})
			}),
			layout.Rigid(pg.purchaseTicket.Layout),
		)
	})
}

func (pg *Page) ticketsLiveSection(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
					tit := pg.Theme.Label(values.TextSize14, "Live Tickets")
					tit.Color = pg.Theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, func(gtx C) D {
						ticketLiveCounter := (*pg.tickets).LiveCounter
						var elements []layout.FlexChild
						for i := 0; i < len(ticketLiveCounter); i++ {
							item := ticketLiveCounter[i]
							elements = append(elements, layout.Rigid(func(gtx C) D {
								return layout.Inset{Right: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											st := ticketStatusIcon(pg.Load, item.Status)
											if st == nil {
												return layout.Dimensions{}
											}
											return st.icon.Layout16dp(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
												label := pg.Theme.Label(values.TextSize14, fmt.Sprintf("%d", item.Count))
												label.Color = pg.Theme.Color.DeepBlue
												return label.Layout(gtx)
											})
										}),
									)
								})
							}))
						}
						elements = append(elements, layout.Rigid(pg.toTickets.Layout))
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx, elements...)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				tickets := (*pg.tickets).LiveRecent
				for range tickets {
					pg.ticketTooltips = append(pg.ticketTooltips, tooltips{
						statusTooltip:     pg.Load.Theme.Tooltip(),
						walletNameTooltip: pg.Load.Theme.Tooltip(),
						dateTooltip:       pg.Load.Theme.Tooltip(),
						daysBehindTooltip: pg.Load.Theme.Tooltip(),
						durationTooltip:   pg.Load.Theme.Tooltip(),
					})
				}

				return pg.ticketsLive.Layout(gtx, len(tickets), func(gtx C, index int) D {
					return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						return ticketCard(gtx, pg.Load, &tickets[index], pg.ticketTooltips[index])
					})
				})
			}),
		)
	})
}

func (pg *Page) ticketsActivitySection(gtx layout.Context) layout.Dimensions {
	tickets := (*pg.tickets).RecentActivity
	if len(tickets) == 0 {
		return layout.Dimensions{}
	}

	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					tit := pg.Theme.Label(values.TextSize14, "Recent Activity")
					tit.Color = pg.Theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, pg.toTicketsActivity.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.ticketsActivity.Layout(gtx, len(tickets), func(gtx C, index int) D {
					return ticketActivityRow(gtx, pg.Load, tickets[index], index)
				})
			}),
		)
	})
}

func (pg *Page) stackingRecordSection(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					tit := pg.Theme.Label(values.TextSize14, "Staking Record")
					tit.Color = pg.Theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, func(gtx C) D { return layout.Dimensions{} })
				})
			}),
			layout.Rigid(func(gtx C) D {
				stackingRecords := (*pg.tickets).StackingRecordCounter
				return decredmaterial.GridWrap{
					Axis:      layout.Horizontal,
					Alignment: layout.End,
				}.Layout(gtx, len(stackingRecords), func(gtx layout.Context, i int) layout.Dimensions {
					item := stackingRecords[i]
					width := unit.Value{U: unit.UnitDp, V: 118}
					gtx.Constraints.Min.X = gtx.Px(width)

					return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								st := ticketStatusIcon(pg.Load, item.Status)
								if st == nil {
									return layout.Dimensions{}
								}
								return st.icon.Layout24dp(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											label := pg.Theme.Label(values.TextSize16, fmt.Sprintf("%d", item.Count))
											label.Color = pg.Theme.Color.DeepBlue
											return label.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
												txt := pg.Theme.Label(values.TextSize12, strings.Title(strings.ToLower(item.Status)))
												txt.Color = pg.Theme.Color.Gray2
												return txt.Layout(gtx)
											})
										}),
									)
								})
							}),
						)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				wrapper := pg.Theme.Card()
				wrapper.Color = pg.Theme.Color.Success2
				return wrapper.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Inset{
							Top:    values.MarginPadding16,
							Bottom: values.MarginPadding16,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Bottom: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
										txt := pg.Theme.Label(values.TextSize14, "Rewards Earned")
										txt.Color = pg.Theme.Color.Success
										return txt.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											ic := pg.Icons.StakeyIcon
											ic.Scale = 1.0
											return ic.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return components.LayoutBalance(gtx, pg.Load, "16.5112316")
										}),
									)
								}),
							)
						})
					})
				})
			}),
		)
	})
}

func (pg *Page) Handle() {
	if pg.purchaseTicket.Button.Clicked() {
		newTicketPurchaseModal(pg.Load).
			Show()
	}

	if pg.toTickets.Button.Clicked() {
		pg.ChangeFragment(newListPage(pg.Load))
	}

	if pg.toTicketsActivity.Button.Clicked() {
		pg.ChangeFragment(newTicketActivityPage(pg.Load))
	}
}

func (pg *Page) OnClose() {}
