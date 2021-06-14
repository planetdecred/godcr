package tickets

import (
	"fmt"
	"image/color"
	"sort"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const listPageID = "TicketsList"

type ListPage struct {
	*load.Load

	tickets      **wallet.Tickets
	ticketsList  layout.List
	filterSorter int
	isGridView   bool

	toggleViewType     *widget.Clickable
	orderDropDown      *decredmaterial.DropDown
	ticketTypeDropDown *decredmaterial.DropDown
	walletDropDown     *decredmaterial.DropDown
	backButton         decredmaterial.IconButton

	ticketTooltips []tooltips
	wallets        []*dcrlibwallet.Wallet
}

func newListPage(l *load.Load) *ListPage {
	pg := &ListPage{
		Load:           l,
		tickets:        l.WL.Tickets,
		ticketsList:    layout.List{Axis: layout.Vertical},
		toggleViewType: new(widget.Clickable),
		isGridView:     true,
	}
	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)

	pg.orderDropDown = createOrderDropDown(l.Theme)
	pg.ticketTypeDropDown = l.Theme.DropDown([]decredmaterial.DropDownItem{
		{Text: "All"},
		{Text: "Unmined"},
		{Text: "Immature"},
		{Text: "Live"},
		{Text: "Voted"},
		{Text: "Missed"},
		{Text: "Expired"},
		{Text: "Revoked"},
	}, 1)

	return pg
}

func (pg *ListPage) ID() string {
	return listPageID
}

func (pg *ListPage) OnResume() {
	pg.wallets = pg.WL.SortedWalletList()
	components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets)
}

func (pg *ListPage) Layout(gtx layout.Context) layout.Dimensions {
	walletID := pg.wallets[pg.walletDropDown.SelectedIndex()].ID
	tickets := (*pg.tickets).Confirmed[walletID]

	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "All tickets",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				for range tickets {
					pg.ticketTooltips = append(pg.ticketTooltips, tooltips{
						statusTooltip:     pg.Load.Theme.Tooltip(),
						walletNameTooltip: pg.Load.Theme.Tooltip(),
						dateTooltip:       pg.Load.Theme.Tooltip(),
						daysBehindTooltip: pg.Load.Theme.Tooltip(),
						durationTooltip:   pg.Load.Theme.Tooltip(),
					})
				}
				return layout.Stack{Alignment: layout.N}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Card().Layout(gtx, func(gtx C) D {
								gtx.Constraints.Min = gtx.Constraints.Max

								if pg.ticketTypeDropDown.SelectedIndex()-1 != -1 {
									tickets = filterTickets(tickets, func(ticketStatus string) bool {
										return ticketStatus == strings.ToUpper(pg.ticketTypeDropDown.Selected())
									})
								}

								if len(tickets) == 0 {
									txt := pg.Theme.Body1("No tickets yet")
									txt.Color = pg.Theme.Color.Gray2
									txt.Alignment = text.Middle
									return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, txt.Layout)
								}
								return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
									if pg.isGridView {
										return pg.ticketListGridLayout(gtx, tickets)
									}
									return pg.ticketListLayout(gtx, tickets)
								})
							})
						})
					}),
					layout.Stacked(pg.dropDowns),
				)
			},
			ExtraItem: pg.toggleViewType,
			Extra: func(gtx C) D {
				wrap := pg.Theme.Card()
				wrap.Color = pg.Theme.Color.Gray1
				wrap.Radius = decredmaterial.Radius(8)
				return wrap.Layout(gtx, func(gtx C) D {
					insetIcon := layout.Inset{
						Top:    values.MarginPadding4,
						Bottom: values.MarginPadding4,
						Left:   values.MarginPadding8,
						Right:  values.MarginPadding8,
					}
					return layout.Inset{
						Left:   values.MarginPadding2,
						Right:  values.MarginPadding2,
						Top:    values.MarginPadding3,
						Bottom: values.MarginPadding3,
					}.Layout(gtx, func(gtx C) D {
						wrapIcon := pg.Theme.Card()
						wrapIcon.Color = pg.Theme.Color.Surface
						wrapIcon.Radius = decredmaterial.Radius(7)

						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								if !pg.isGridView {
									wrapIcon.Color = color.NRGBA{}
								}
								return wrapIcon.Layout(gtx, func(gtx C) D {
									ic := pg.Icons.ListGridIcon
									ic.Scale = 1
									return insetIcon.Layout(gtx, ic.Layout)
								})
							}),
							layout.Rigid(func(gtx C) D {
								if pg.isGridView {
									wrapIcon.Color = color.NRGBA{}
								} else {
									wrapIcon.Color = pg.Theme.Color.Surface
								}
								return wrapIcon.Layout(gtx, func(gtx C) D {
									ic := pg.Icons.List
									ic.Scale = 1
									return insetIcon.Layout(gtx, ic.Layout)
								})
							}),
						)
					})
				})
			},
		}
		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *ListPage) dropDowns(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.walletDropDown.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left: values.MarginPadding5,
					}.Layout(gtx, pg.ticketTypeDropDown.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left: values.MarginPadding5,
					}.Layout(gtx, pg.orderDropDown.Layout)
				}),
			)
		}),
	)
}

func (pg *ListPage) ticketListLayout(gtx layout.Context, tickets []wallet.Ticket) layout.Dimensions {
	return pg.ticketsList.Layout(gtx, len(tickets), func(gtx C, index int) D {
		st := ticketStatusIcon(pg.Load, tickets[index].Info.Status)
		if st == nil {
			return layout.Dimensions{}
		}

		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
					var progressBarWidth int
					return layout.Stack{Alignment: layout.S}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							wrapIcon := pg.Theme.Card()
							wrapIcon.Color = st.background
							wrapIcon.Radius = decredmaterial.Radius(8)
							st.icon.Scale = 0.6
							dims := wrapIcon.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, st.icon.Layout)
							})
							progressBarWidth = dims.Size.X
							return dims
						}),
						layout.Stacked(func(gtx C) D {
							gtx.Constraints.Max.X = progressBarWidth - 4
							blockHeight := tickets[index].Info.BlockHeight
							percent := getPercentConfirmation(pg.Load, blockHeight)
							if percent >= 100 {
								return layout.Dimensions{}
							}

							p := pg.Load.Theme.ProgressBar(percent)
							p.Height, p.Radius = values.MarginPadding4, values.MarginPadding2
							p.Color = st.color
							return p.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if index == 0 {
							return layout.Dimensions{}
						}
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						separator := pg.Theme.Separator()
						separator.Width = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, separator.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Top:    values.MarginPadding6,
							Bottom: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									dtime := pg.Theme.Label(values.TextSize14, tickets[index].DateTime)
									dtime.Color = pg.Theme.Color.Gray2
									return components.EndToEndRow(gtx, func(gtx C) D {
										return components.LayoutBalance(gtx, pg.Load, tickets[index].Amount)
									}, dtime.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									l := func(gtx C) layout.Dimensions {
										return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												txt := pg.Theme.Label(values.MarginPadding14, tickets[index].Info.Status)
												txt.Color = st.color
												return txt.Layout(gtx)
											}),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{
													Left:  values.MarginPadding4,
													Right: values.MarginPadding4,
												}.Layout(gtx, func(gtx C) D {
													ic := pg.Icons.ImageBrightness1
													ic.Color = pg.Theme.Color.Gray2
													return pg.Icons.ImageBrightness1.Layout(gtx, values.MarginPadding5)
												})
											}),
											layout.Rigid(pg.Theme.Label(values.MarginPadding14, tickets[index].WalletName).Layout),
										)
									}
									r := func(gtx C) layout.Dimensions {

										timeBehind, unit := getTimeBehind(tickets[index].DateTime)
										if timeBehind == 0 && unit == "h" {
											return layout.Dimensions{}
										}

										txt := pg.Load.Theme.Label(values.TextSize14, fmt.Sprintf("%d%s", timeBehind, unit))
										txt.Color = pg.Load.Theme.Color.Gray2

										return txt.Layout(gtx)
									}
									return components.EndToEndRow(gtx, l, r)
								}),
							)
						})
					}),
				)
			}),
		)
	})
}

func (pg *ListPage) ticketListGridLayout(gtx layout.Context, tickets []wallet.Ticket) layout.Dimensions {
	// TODO: GridWrap's items not able to scroll vertically, will update when it fixed
	return layout.NW.Layout(gtx, func(gtx C) D {
		return pg.ticketsList.Layout(gtx, 1, func(gtx C, index int) D {
			return pg.Theme.Card().Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = gtx.Constraints.Max
				return decredmaterial.GridWrap{
					Axis:      layout.Horizontal,
					Alignment: layout.End,
				}.Layout(gtx, len(tickets), func(gtx C, index int) D {
					return layout.Inset{
						Left:   values.MarginPadding4,
						Right:  values.MarginPadding4,
						Bottom: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return ticketCard(gtx, pg.Load, &tickets[index], pg.ticketTooltips[index])
					})
				})
			})
		})
	})
}

func (pg *ListPage) Handle() {

	if pg.toggleViewType.Clicked() {
		pg.isGridView = !pg.isGridView
	}

	sortSelection := pg.orderDropDown.SelectedIndex()
	if pg.filterSorter != sortSelection {
		pg.filterSorter = sortSelection
		newestFirst := pg.filterSorter == 0
		for _, wal := range pg.wallets {
			tickets := (*pg.tickets).Confirmed[wal.ID]
			sort.SliceStable(tickets, func(i, j int) bool {
				backTime := time.Unix(tickets[j].Info.Ticket.Timestamp, 0)
				frontTime := time.Unix(tickets[i].Info.Ticket.Timestamp, 0)
				if newestFirst {
					return backTime.Before(frontTime)
				}
				return frontTime.Before(backTime)
			})
		}
	}
}

func (pg *ListPage) OnClose() {}
