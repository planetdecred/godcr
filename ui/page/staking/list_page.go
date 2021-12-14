package staking

import (
	"context"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	tpage "github.com/planetdecred/godcr/ui/page/transaction"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const listPageID = "StakingList"

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
	ticketsList *decredmaterial.ClickableList
	scrollBar   *widget.List

	orderDropDown      *decredmaterial.DropDown
	ticketTypeDropDown *decredmaterial.DropDown
	walletDropDown     *decredmaterial.DropDown
	backButton         decredmaterial.IconButton

	wallets []*dcrlibwallet.Wallet
}

func newListPage(l *load.Load) *ListPage {
	pg := &ListPage{
		Load:        l,
		ticketsList: l.Theme.NewClickableList(layout.Vertical),
		scrollBar: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)

	pg.orderDropDown = createOrderDropDown(l.Theme)
	pg.wallets = pg.WL.SortedWalletList()
	components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets, 0) // first in the set during layout pos 0
	pg.ticketTypeDropDown = l.Theme.DropDown([]decredmaterial.DropDownItem{
		{Text: "All"},
		{Text: "Unmined"},
		{Text: "Immature"},
		{Text: "Live"},
		{Text: "Voted"},
		{Text: "Expired"},
		{Text: "Revoked"},
	}, 1, 2)

	return pg
}

func (pg *ListPage) ID() string {
	return listPageID
}

func (pg *ListPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
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

	tickets, err := stakeToTransactionItems(pg.Load, txs, newestFirst, func(filter int32) bool {
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
							return pg.Theme.List(pg.scrollBar).Layout(gtx, 1, func(gtx C, index int) D {
								return pg.Theme.Card().Layout(gtx, func(gtx C) D {
									tickets := pg.tickets

									if len(tickets) == 0 {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X

										txt := pg.Theme.Body1("No tickets yet")
										txt.Color = pg.Theme.Color.GrayText3
										txt.Alignment = text.Middle
										return layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding16}.Layout(gtx, txt.Layout)
									}

									return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
										return pg.ticketListLayout(gtx, tickets)
									})
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
						return pg.ticketTypeDropDown.Layout(gtx, pg.orderDropDown.Width-4, true)
					}),
				)
			},
		}
		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *ListPage) ticketListLayout(gtx layout.Context, tickets []*transactionItem) layout.Dimensions {
	return pg.ticketsList.Layout(gtx, len(tickets), func(gtx C, index int) D {
		var ticket = tickets[index]

		return ticketListLayout(gtx, pg.Load, ticket, index, false)
	})
}

func (pg *ListPage) Handle() {
	for pg.orderDropDown.Changed() {
		pg.fetchTickets()
	}

	for pg.walletDropDown.Changed() {
		pg.fetchTickets()
	}

	for pg.ticketTypeDropDown.Changed() {
		pg.fetchTickets()
	}

	if clicked, selectedItem := pg.ticketsList.ItemClicked(); clicked {
		pg.ChangeFragment(tpage.NewTransactionDetailsPage(pg.Load, pg.tickets[selectedItem].transaction))
	}

	decredmaterial.DisplayOneDropdown(pg.ticketTypeDropDown, pg.orderDropDown, pg.walletDropDown)
}

func (pg *ListPage) OnClose() {
	pg.ctxCancel()
}
