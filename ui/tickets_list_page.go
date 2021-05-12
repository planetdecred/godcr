package ui

import (
	"image/color"
	"sort"
	"time"

	"github.com/planetdecred/godcr/wallet"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageTicketsList = "TicketsList"

type ticketPageList struct {
	th           *decredmaterial.Theme
	tickets      **wallet.Tickets
	ticketsList  layout.List
	filterSorter int

	toggleViewType     *widget.Clickable
	orderDropDown      *decredmaterial.DropDown
	ticketTypeDropDown *decredmaterial.DropDown
	walletDropDown     *decredmaterial.DropDown
	isGridView         bool
}

func (win *Window) TicketPageList(c pageCommon) layout.Widget {
	pg := &ticketPageList{
		th:             c.theme,
		tickets:        &win.walletTickets,
		ticketsList:    layout.List{Axis: layout.Vertical},
		toggleViewType: new(widget.Clickable),
		isGridView:     true,
	}
	pg.orderDropDown = c.theme.DropDown([]decredmaterial.DropDownItem{{Text: "Newest"}, {Text: "Oldest"}}, 1)
	pg.ticketTypeDropDown = c.theme.DropDown([]decredmaterial.DropDownItem{
		{Text: "All"},
		{Text: "Unmined"},
		{Text: "Immature"},
		{Text: "Live"},
		{Text: "Voted"},
		{Text: "Missed"},
		{Text: "Expired"},
		{Text: "Revoked"},
	}, 1)

	return func(gtx C) D {
		pg.handler(c)
		return pg.layout(gtx, c)
	}
}

func (pg *ticketPageList) layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		page := SubPage{
			title: "All tickets",
			back: func() {
				c.changePage(PageTickets)
			},
			body: func(gtx C) D {
				walletID := c.info.Wallets[pg.walletDropDown.SelectedIndex()].ID
				tickets := (*pg.tickets).Confirmed[walletID]
				return layout.Stack{Alignment: layout.N}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, func(gtx C) D {
							return c.theme.Card().Layout(gtx, func(gtx C) D {
								gtx.Constraints.Min = gtx.Constraints.Max
								if len(tickets) == 0 {
									txt := c.theme.Body1("No tickets yet")
									txt.Color = c.theme.Color.Gray2
									txt.Alignment = text.Middle
									return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D { return txt.Layout(gtx) })
								}
								return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
									if pg.isGridView {
										return pg.ticketListGridLayout(gtx, c, tickets)
									}
									return pg.ticketListLayout(gtx, c, tickets)
								})
							})
						})
					}),
					layout.Stacked(func(gtx C) D {
						return pg.dropDowns(gtx)
					}),
				)
			},
			extraItem: pg.toggleViewType,
			extra: func(gtx C) D {
				wrap := c.theme.Card()
				wrap.Color = c.theme.Color.Gray1
				wrap.Radius.NE = 8 // top - left
				wrap.Radius.SW = 8 // bottom - left
				wrap.Radius.NW = 8 // top - right
				wrap.Radius.SE = 8 // bottom - right
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
						wrapIcon := c.theme.Card()
						wrapIcon.Color = c.theme.Color.Surface
						wrapIcon.Radius.NE = 7
						wrapIcon.Radius.SW = 7
						wrapIcon.Radius.NW = 7
						wrapIcon.Radius.SE = 7

						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								if !pg.isGridView {
									wrapIcon.Color = color.NRGBA{}
								}
								return wrapIcon.Layout(gtx, func(gtx C) D {
									return insetIcon.Layout(gtx, func(gtx C) D {
										ic := c.icons.listGridIcon
										ic.Scale = 1
										return ic.Layout(gtx)
									})
								})
							}),
							layout.Rigid(func(gtx C) D {
								if pg.isGridView {
									wrapIcon.Color = color.NRGBA{}
								} else {
									wrapIcon.Color = c.theme.Color.Surface
								}
								return wrapIcon.Layout(gtx, func(gtx C) D {
									return insetIcon.Layout(gtx, func(gtx C) D {
										ic := c.icons.list
										ic.Scale = 1
										return ic.Layout(gtx)
									})
								})
							}),
						)
					})
				})
			},
		}
		return c.SubPageLayout(gtx, page)
	}

	return c.Layout(gtx, func(gtx C) D {
		return c.UniformPadding(gtx, body)
	})
}

func (pg *ticketPageList) dropDowns(gtx layout.Context) layout.Dimensions {
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
					}.Layout(gtx, func(gtx C) D {
						return pg.ticketTypeDropDown.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left: values.MarginPadding5,
					}.Layout(gtx, func(gtx C) D {
						return pg.orderDropDown.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *ticketPageList) ticketListLayout(gtx layout.Context, c pageCommon, tickets []wallet.Ticket) layout.Dimensions {
	return pg.ticketsList.Layout(gtx, len(tickets), func(gtx C, index int) D {
		st := ticketIconStatus(&c, tickets[index].Info.Status)
		if st == nil {
			return layout.Dimensions{}
		}

		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				wrapIcon := c.theme.Card()
				wrapIcon.Color = st.background
				wrapIcon.Radius.NE = 8
				wrapIcon.Radius.SW = 8
				wrapIcon.Radius.NW = 8
				wrapIcon.Radius.SE = 8
				st.icon.Scale = 0.6
				return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
					return wrapIcon.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding10).Layout(gtx, st.icon.Layout)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if index == 0 {
							return layout.Dimensions{}
						}
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						separator := pg.th.Separator()
						separator.Width = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, func(gtx C) D {
							return separator.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Top:    values.MarginPadding6,
							Bottom: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									dtime := c.theme.Label(values.TextSize14, tickets[index].DateTime)
									dtime.Color = c.theme.Color.Gray2
									return endToEndRow(gtx, func(gtx C) D { return c.layoutBalance(gtx, tickets[index].Amount) }, dtime.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											txt := c.theme.Label(values.MarginPadding14, tickets[index].Info.Status)
											txt.Color = st.color
											return txt.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Left:  values.MarginPadding4,
												Right: values.MarginPadding4,
											}.Layout(gtx, func(gtx C) D {
												ic := c.icons.imageBrightness1
												ic.Color = c.theme.Color.Gray2
												return c.icons.imageBrightness1.Layout(gtx, values.MarginPadding5)
											})
										}),
										layout.Rigid(func(gtx C) D {
											return c.theme.Label(values.MarginPadding14, tickets[index].WalletName).Layout(gtx)
										}),
									)
								}),
							)
						})
					}),
				)
			}),
		)
	})
}

func (pg *ticketPageList) ticketListGridLayout(gtx layout.Context, c pageCommon, tickets []wallet.Ticket) layout.Dimensions {
	// TODO: GridWrap's items not able to scroll vertically, will update when it fixed
	return pg.ticketsList.Layout(gtx, 1, func(gtx C, index int) D {
		return c.theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min = gtx.Constraints.Max
			return decredmaterial.GridWrap{
				Axis:      layout.Horizontal,
				Alignment: layout.End,
			}.Layout(gtx, len(tickets), func(gtx C, index int) D {
				return layout.Inset{Right: values.MarginPadding8, Bottom: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
					return ticketLiveItemnInfo(gtx, c, &tickets[index])
				})
			})
		})
	})
}

func (pg *ticketPageList) initWalletDropDown(common pageCommon) {
	if len(common.info.Wallets) == 0 || pg.walletDropDown != nil {
		return
	}

	var walletDropDownItems []decredmaterial.DropDownItem
	for i := range common.info.Wallets {
		item := decredmaterial.DropDownItem{
			Text: common.info.Wallets[i].Name,
			Icon: common.icons.walletIcon,
		}
		walletDropDownItems = append(walletDropDownItems, item)
	}
	pg.walletDropDown = common.theme.DropDown(walletDropDownItems, 2)
}

func (pg *ticketPageList) handler(c pageCommon) {
	pg.initWalletDropDown(c)

	if pg.toggleViewType.Clicked() {
		pg.isGridView = !pg.isGridView
	}

	sortSelection := pg.orderDropDown.SelectedIndex()
	if pg.filterSorter != sortSelection {
		pg.filterSorter = sortSelection
		newestFirst := pg.filterSorter == 0
		for _, wal := range c.info.Wallets {
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
