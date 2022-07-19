package staking

import (
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *Page) initTicketList() {
	pg.ticketsList = pg.Theme.NewClickableList(layout.Vertical)
}

func (pg *Page) listenForTxNotifications() {
	if pg.TxAndBlockNotificationListener != nil {
		return
	}
	pg.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err := pg.WL.MultiWallet.AddTxAndBlockNotificationListener(pg.TxAndBlockNotificationListener, true, OverviewPageID)
	if err != nil {
		log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	go func() {
		for {
			select {
			case n := <-pg.TxAndBlockNotifChan:
				if n.Type == listeners.BlockAttached || n.Type == listeners.NewTransaction {
					pg.fetchTickets()
					pg.ParentWindow().Reload()
				}
			case <-pg.ctx.Done():
				pg.WL.MultiWallet.RemoveTxAndBlockNotificationListener(OverviewPageID)
				close(pg.TxAndBlockNotifChan)
				pg.TxAndBlockNotificationListener = nil

				return
			}
		}
	}()
}

func (pg *Page) fetchTickets() {
	txs, err := pg.WL.SelectedWallet.Wallet.GetTransactionsRaw(0, 0, dcrlibwallet.TxFilterTickets, true)
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
	}

	tickets, err := stakeToTransactionItems(pg.Load, txs, true, func(filter int32) bool {
		return filter == dcrlibwallet.TxFilterTickets
	})
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
	}

	pg.tickets = tickets
}

func (pg *Page) ticketListLayout(gtx C) D {
	return pg.pageSections(gtx, func(gtx C) D {
		tickets := pg.tickets

		if len(tickets) == 0 {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X

			txt := pg.Theme.Body1(values.String(values.StrNoTickets))
			txt.Color = pg.Theme.Color.GrayText3
			txt.Alignment = text.Middle
			return layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding16}.Layout(gtx, txt.Layout)
		}

		return pg.ticketsList.Layout(gtx, len(tickets), func(gtx C, index int) D {
			var ticket = tickets[index]

			return ticketListLayout(gtx, pg.Load, ticket, index, false)
		})
	})
}
