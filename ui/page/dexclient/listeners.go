package dexclient

import (
	"fmt"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex/msgjson"
)

func (pg *Page) connectDex(h string, password []byte) {
	pg.DL.Dexc.ConnectDexes(h, password)
	go pg.listenerMessages()
	go pg.readNotifications()
	pg.updateOrderBook()
}

func (pg *Page) updateOrderBook() {
	orderBoook, err := pg.DL.Dexc.Book(pg.selectedMaket.host, pg.selectedMaket.marketBaseID, pg.selectedMaket.marketQuoteID)
	if err != nil {
		return
	}
	pg.orderBook = orderBoook
	pg.miniTradeFormWdg.orderBook = pg.orderBook
}

func (pg *Page) listenerMessages() {
	msgs := pg.DL.Dexc.MessageSource(pg.selectedMaket.host)
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				fmt.Println("[ERROR] Listen(wc): Connection terminated for.", "test")
				return
			}
			switch msg.Type {
			case msgjson.Request:
			case msgjson.Notification:
				pg.noteHandlers(msg)
			case msgjson.Response:
				// client/comms.wsConn handles responses to requests we sent.
				continue
			default:
				continue
			}
		case <-pg.ctx.Done():
			return
		}
	}
}

// readNotifications reads from the Core notification channel.
func (pg *Page) readNotifications() {
	ch := pg.DL.NotificationFeed()
	for {
		select {
		case n := <-ch:
			fmt.Println("Recv notification", n)
			fmt.Println("<INFO>", n.ID())
			fmt.Println("<INFO>", n.Severity())
			fmt.Println("<INFO>", n.Type())
			fmt.Println("<INFO>", n.DBNote())
			fmt.Println("<INFO>", n.String())
			fmt.Println("<INFO>", n.Subject())

			if n.Type() == core.NoteTypeFeePayment {
				pg.RefreshWindow()
			}
			pg.refreshUser()
		case <-pg.ctx.Done():
			return
		}
	}
}

func (pg *Page) noteHandlers(msg *msgjson.Message) {
	fmt.Println(">>> Receive message source: noteHandlers", msg.Route)
	switch msg.Route {
	case msgjson.BookOrderRoute:
		pg.updateOrderBook()
	case msgjson.EpochOrderRoute:
		pg.updateOrderBook()
	case msgjson.UnbookOrderRoute:
		pg.updateOrderBook()
	case msgjson.MatchProofRoute:
	case msgjson.UpdateRemainingRoute:
	case msgjson.EpochReportRoute:
	case msgjson.SuspensionRoute:
	case msgjson.ResumptionRoute:
	case msgjson.NotifyRoute:
	case msgjson.PenaltyRoute:
	case msgjson.NoMatchRoute:
	case msgjson.RevokeOrderRoute:
	case msgjson.RevokeMatchRoute:
	}
}
