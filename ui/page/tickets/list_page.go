package tickets

import (
	"context"
	"fmt"
	"image/color"

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

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	tickets      []Ticket
	ticketsList  layout.List
	filterSorter int
	isGridView   bool

	toggleViewType *widget.Clickable

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
		{Text: "Revoked"},
	}, 1)

	return pg
}

func (pg *ListPage) ID() string {
	return listPageID
}

func (pg *ListPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.wallets = pg.WL.SortedWalletList()
	components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets)
	pg.listenForTxNotifications()
	pg.fetchTickets()
}

func (pg *ListPage) listenForTxNotifications() {
	go func() {
		for {
			var notification interface{}

			select {
			case notification = <-pg.Receiver.NotificationsUpdate:
			case <-pg.ctx.Done():
				return
			}

			switch n := notification.(type) {
			case wallet.NewTransaction:
				selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
				if selectedWallet.ID == n.Transaction.WalletID {
					pg.fetchTickets()
					pg.RefreshWindow()
				}
			}
		}
	}()
}

func (pg *ListPage) fetchTickets() {
	var txFilter int32
	switch pg.ticketTypeDropDown.SelectedIndex() {
	case 0:
		txFilter = dcrlibwallet.TxFilterStaking
	case 1:
		txFilter = dcrlibwallet.TxFilterUnmined
	case 2:
		txFilter = dcrlibwallet.TxFilterImmature
	case 3:
		txFilter = dcrlibwallet.TxFilterLive
	case 4:
		txFilter = dcrlibwallet.TxFilterVoted
	case 5:
		txFilter = dcrlibwallet.TxFilterRevoked
	default:
		return
	}

	newestFirst := pg.orderDropDown.SelectedIndex() == 0
	selectedWalletID := pg.wallets[pg.walletDropDown.SelectedIndex()].ID
	tickets, err := getTickets(pg.WL.MultiWallet, selectedWalletID, txFilter, newestFirst)
	if err != nil {
		pg.Toast.NotifyError(err.Error())
	} else {
		fmt.Printf("ticket length %v  filter %v\n", len(tickets), txFilter)
		pg.tickets = tickets
	}
}

func (pg *ListPage) Layout(gtx C) D {
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
								if len(pg.tickets) == 0 {
									txt := pg.Theme.Body1("No tickets yet")
									txt.Color = pg.Theme.Color.Gray2
									txt.Alignment = text.Middle
									return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, txt.Layout)
								}
								return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
									if pg.isGridView {
										return pg.ticketListGridLayout(gtx, pg.tickets)
									}
									return pg.ticketListLayout(gtx, pg.tickets)
								})
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return pg.walletDropDown.Layout(gtx, 0, false)
					}),
					layout.Expanded(func(gtx C) D {
						return pg.orderDropDown.Layout(gtx, 0, true)
					}),
					layout.Expanded(func(gtx C) D {
						return pg.ticketTypeDropDown.Layout(gtx, pg.orderDropDown.Width+10, true)
					}),
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
									return insetIcon.Layout(gtx, ic.Layout16dp)
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
									return insetIcon.Layout(gtx, ic.Layout16dp)
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
		layout.Rigid(pg.walletDropDown.Layout),
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

func (pg *ListPage) ticketListLayout(gtx C, tickets []load.Ticket) D {
	return pg.ticketsList.Layout(gtx, len(tickets), func(gtx C, index int) D {
		st := ticketStatusProfile(pg.Load, tickets[index].Status)
		if st == nil {
			return D{}
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
							dims := wrapIcon.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, st.icon.Layout24dp)
							})
							progressBarWidth = dims.Size.X
							return dims
						}),
						layout.Stacked(func(gtx C) D {
							gtx.Constraints.Max.X = progressBarWidth - 4
							p := pg.Theme.ProgressBar(20)
							p.Height, p.Radius = values.MarginPadding4, values.MarginPadding2
							p.Color = st.color
							return p.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if index == 0 {
							return D{}
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
									l := func(gtx C) D {
										return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												txt := pg.Theme.Label(values.MarginPadding14, tickets[index].Status)
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
									r := func(gtx C) D {
										txt := pg.Theme.Label(values.TextSize14, tickets[index].DaysBehind)
										txt.Color = pg.Theme.Color.Gray2
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

func (pg *ListPage) ticketListGridLayout(gtx C, tickets []load.Ticket) D {
	// TODO: GridWrap's items not able to scroll vertically, will update when it fixed
	return layout.Center.Layout(gtx, func(gtx C) D {
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

	for pg.orderDropDown.Changed() {
		pg.fetchTickets()
	}

	for pg.walletDropDown.Changed() {
		pg.fetchTickets()
	}

	for pg.ticketTypeDropDown.Changed() {
		pg.fetchTickets()
	}
}

func (pg *ListPage) OnClose() {
	pg.ctxCancel()
}
