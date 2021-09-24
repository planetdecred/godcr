package tickets

import (
	"context"
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const listPageID = "TicketsList"

type txType int

const (
	All txType = iota
	Unmined
	Immature
	Live
	Voted
	Expired
	Revoked
)

type ListPage struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	tickets     []*transactionItem
	ticketsList layout.List
	isGridView  bool

	toggleViewType *widget.Clickable

	orderDropDown      *decredmaterial.DropDown
	ticketTypeDropDown *decredmaterial.DropDown
	walletDropDown     *decredmaterial.DropDown
	backButton         decredmaterial.IconButton

	wallets []*dcrlibwallet.Wallet
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
		{Text: "Expired"},
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
			case wallet.NewBlock:
				selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
				if selectedWallet.ID == n.WalletID {
					pg.fetchTickets()
					pg.RefreshWindow()
				}
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
	var ticketTypeDropdown = txType(pg.ticketTypeDropDown.SelectedIndex())
	switch ticketTypeDropdown {
	case Unmined:
		txFilter = dcrlibwallet.TxFilterUnmined
	case Immature:
		txFilter = dcrlibwallet.TxFilterImmature
	case Live:
		txFilter = dcrlibwallet.TxFilterLive
	case Voted:
		txFilter = dcrlibwallet.TxFilterTickets
	case Expired:
		txFilter = dcrlibwallet.TxFilterExpired
	case Revoked:
		txFilter = dcrlibwallet.TxFilterTickets
	default:
		txFilter = dcrlibwallet.TxFilterTickets
	}

	newestFirst := pg.orderDropDown.SelectedIndex() == 0
	selectedWalletID := pg.wallets[pg.walletDropDown.SelectedIndex()].ID
	multiWallet := pg.WL.MultiWallet
	w := multiWallet.WalletWithID(selectedWalletID)
	txs, err := w.GetTransactionsRaw(0, 0, txFilter, newestFirst)
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
	}

	tickets, err := ticketsToTransactionItems(pg.Load, txs, newestFirst, func(filter int32) bool {
		switch filter {
		case dcrlibwallet.TxFilterVoted:
			return ticketTypeDropdown == Voted
		case dcrlibwallet.TxFilterRevoked:
			return ticketTypeDropdown == Revoked
		}

		return filter == txFilter
	})
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
	}

	pg.tickets = tickets
}

func (pg *ListPage) Layout(gtx C) D {
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "All tickets",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return layout.Stack{Alignment: layout.N}.Layout(gtx,
					layout.Expanded(func(gtx C) D {

						return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Card().Layout(gtx, func(gtx C) D {
								tickets := pg.tickets

								if len(tickets) == 0 {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X

									txt := pg.Theme.Body1("No tickets yet")
									txt.Color = pg.Theme.Color.Gray2
									txt.Alignment = text.Middle
									return layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding16}.Layout(gtx, txt.Layout)
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
					layout.Expanded(func(gtx C) D {
						return pg.walletDropDown.Layout(gtx, 0, false)
					}),
					layout.Expanded(func(gtx C) D {
						return pg.orderDropDown.Layout(gtx, 0, true)
					}),
					layout.Expanded(func(gtx C) D {
						return pg.ticketTypeDropDown.Layout(gtx, pg.orderDropDown.Width, true)
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

func (pg *ListPage) ticketListLayout(gtx layout.Context, tickets []*transactionItem) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max
	return pg.ticketsList.Layout(gtx, len(tickets), func(gtx C, index int) D {
		var ticket = tickets[index]

		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
					return layout.Stack{Alignment: layout.S}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							wrapIcon := pg.Theme.Card()
							wrapIcon.Color = ticket.status.Background
							wrapIcon.Radius = decredmaterial.Radius(8)
							dims := wrapIcon.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, ticket.status.Icon.Layout24dp)
							})
							return dims
						}),
						layout.Expanded(func(gtx C) D {
							if !ticket.showProgress {
								return D{}
							}
							p := pg.Theme.ProgressBar(int(ticket.progress))
							p.Width = values.MarginPadding44
							p.Height = values.MarginPadding4
							p.Direction = layout.SW
							p.Radius = decredmaterial.BottomRadius(8)
							p.Color = ticket.status.ProgressBarColor
							p.TrackColor = ticket.status.ProgressTrackColor
							return p.Layout2(gtx)
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

									dtime := pg.Theme.Label(values.TextSize14, ticket.purchaseTime)
									dtime.Color = pg.Theme.Color.Gray2
									return components.EndToEndRow(gtx, func(gtx C) D {
										return components.LayoutBalance(gtx, pg.Load, dcrutil.Amount(ticket.transaction.Amount).String())
									}, dtime.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									l := func(gtx C) layout.Dimensions {
										txt := pg.Theme.Label(values.MarginPadding14, ticket.status.Title)
										txt.Color = ticket.status.Color
										return txt.Layout(gtx)
									}
									r := func(gtx C) layout.Dimensions {

										if ticket.ticketAge == "" {
											return D{}
										}

										txt := pg.Theme.Label(values.TextSize14, ticket.ticketAge)
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

func (pg *ListPage) ticketListGridLayout(gtx layout.Context, tickets []*transactionItem) layout.Dimensions {
	// TODO: GridWrap's items not able to scroll vertically, will update when it fixed
	return layout.Center.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min = gtx.Constraints.Max
			return decredmaterial.GridLayout{
				List:              &pg.ticketsList,
				HorizontalSpacing: layout.SpaceBetween,
				RowCount:          3,
			}.Layout(gtx, len(tickets), func(gtx C, index int) D {
				return layout.Inset{
					Left:   values.MarginPadding4,
					Right:  values.MarginPadding4,
					Bottom: values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					return ticketCard(gtx, pg.Load, tickets[index], false)
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
