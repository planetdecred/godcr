package ui

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageTickets = "Tickets"

type ticketPage struct {
	*pageCommon

	ticketPageContainer layout.List
	ticketsLive         layout.List
	ticketsActivity     layout.List

	purchaseTicket        decredmaterial.Button

	tickets               **wallet.Tickets
	ticketPrice           string
	remainingBalance      string

	autoPurchaseEnabled *widget.Bool
	toTickets           decredmaterial.TextAndIconButton
	toTicketsActivity   decredmaterial.TextAndIconButton
}

func TicketPage(c *pageCommon) Page {
	pg := &ticketPage{
		pageCommon: c,
		tickets:    c.walletTickets,

		ticketsLive:           layout.List{Axis: layout.Horizontal},
		ticketsActivity:       layout.List{Axis: layout.Vertical},
		ticketPageContainer:   layout.List{Axis: layout.Vertical},
		purchaseTicket:        c.theme.Button(new(widget.Clickable), "Purchase"),

		autoPurchaseEnabled: new(widget.Bool),
		toTickets:           c.theme.TextAndIconButton(new(widget.Clickable), "See All", c.icons.navigationArrowForward),
		toTicketsActivity:   c.theme.TextAndIconButton(new(widget.Clickable), "See All", c.icons.navigationArrowForward),
	}

	pg.purchaseTicket.TextSize = values.TextSize12
	pg.purchaseTicket.Background = c.theme.Color.Primary


	pg.toTickets.Color = c.theme.Color.Primary
	pg.toTickets.BackgroundColor = c.theme.Color.Surface

	pg.toTicketsActivity.Color = c.theme.Color.Primary
	pg.toTicketsActivity.BackgroundColor = c.theme.Color.Surface

	return pg
}

func (pg *ticketPage) OnResume() {
	pg.ticketPrice = dcrutil.Amount(pg.wallet.TicketPrice()).String()
	go pg.GetVSPList()
}

func (pg *ticketPage) Layout(gtx layout.Context) layout.Dimensions {
	return pg.UniformPadding(gtx, func(gtx layout.Context) layout.Dimensions {
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

func (pg *ticketPage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.Inset{
		Bottom: values.MarginPadding8,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, body)
		})
	})
}

func (pg *ticketPage) titleRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(leftWidget),
		layout.Rigid(rightWidget),
	)
}

func (pg *ticketPage) ticketPriceSection(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding11,
				}.Layout(gtx, func(gtx C) D {
					tit := pg.theme.Label(values.TextSize14, "Ticket Price")
					tit.Color = pg.theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, material.Switch(pg.theme.Base, pg.autoPurchaseEnabled).Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					ic := pg.icons.ticketPurchasedIcon
					ic.Scale = 1.2
					return layout.Center.Layout(gtx, ic.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						mainText, subText := breakBalance(pg.printer, pg.ticketPrice)
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								label := pg.theme.Label(values.TextSize28, mainText)
								label.Color = pg.theme.Color.DeepBlue
								return label.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								label := pg.theme.Label(values.TextSize16, subText)
								label.Color = pg.theme.Color.DeepBlue
								return label.Layout(gtx)
							}),
						)
					})
				})
			}),
			layout.Rigid(pg.purchaseTicket.Layout),
		)
	})
}

func (pg *ticketPage) ticketsLiveSection(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
					tit := pg.theme.Label(values.TextSize14, "Live Tickets")
					tit.Color = pg.theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, func(gtx C) D {
						ticketLiveCounter := (*pg.tickets).LiveCounter
						var elements []layout.FlexChild
						for i := 0; i < len(ticketLiveCounter); i++ {
							item := ticketLiveCounter[i]
							elements = append(elements, layout.Rigid(func(gtx C) D {
								return layout.Inset{Right: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											st := ticketStatusIcon(pg.pageCommon, item.Status)
											if st == nil {
												return layout.Dimensions{}
											}
											st.icon.Scale = .5
											return st.icon.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
												label := pg.theme.Label(values.TextSize14, fmt.Sprintf("%d", item.Count))
												label.Color = pg.theme.Color.DeepBlue
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
				return pg.ticketsLive.Layout(gtx, len(tickets), func(gtx C, index int) D {
					return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						return ticketCard(gtx, pg.pageCommon, &tickets[index], pg.theme.Tooltip())
					})
				})
			}),
		)
	})
}

func (pg *ticketPage) ticketsActivitySection(gtx layout.Context) layout.Dimensions {
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
					tit := pg.theme.Label(values.TextSize14, "Recent Activity")
					tit.Color = pg.theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, pg.toTicketsActivity.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.ticketsActivity.Layout(gtx, len(tickets), func(gtx C, index int) D {
					return ticketActivityRow(gtx, pg.pageCommon, tickets[index], index)
				})
			}),
		)
	})
}

func (pg *ticketPage) stackingRecordSection(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					tit := pg.theme.Label(values.TextSize14, "Staking Record")
					tit.Color = pg.theme.Color.Gray2
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
								st := ticketStatusIcon(pg.pageCommon, item.Status)
								if st == nil {
									return layout.Dimensions{}
								}
								st.icon.Scale = 0.6
								return st.icon.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											label := pg.theme.Label(values.TextSize16, fmt.Sprintf("%d", item.Count))
											label.Color = pg.theme.Color.DeepBlue
											return label.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
												txt := pg.theme.Label(values.TextSize12, strings.Title(strings.ToLower(item.Status)))
												txt.Color = pg.theme.Color.Gray2
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
				wrapper := pg.theme.Card()
				wrapper.Color = pg.theme.Color.Success2
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
										txt := pg.theme.Label(values.TextSize14, "Rewards Earned")
										txt.Color = pg.theme.Color.Success
										return txt.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											ic := pg.icons.stakeyIcon
											ic.Scale = 1.0
											return ic.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return pg.layoutBalance(gtx, "16.5112316", false)
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

func (pg *ticketPage) handle() {
	if pg.purchaseTicket.Button.Clicked() {
		newTicketPurchaseModal(pg.pageCommon).
			Show()
	}

	if pg.toTickets.Button.Clicked() {
		pg.changePage(PageTicketsList)
	}

	if pg.toTicketsActivity.Button.Clicked() {
		pg.changePage(PageTicketsActivity)
	}
}

func (pg *ticketPage) onClose() {}
