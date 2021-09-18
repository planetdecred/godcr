package dexclient

import (
	"fmt"

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
				fmt.Errorf("listen(wc): Connection terminated for %wsc.", "test")
				return
			}
			switch msg.Type {
			case msgjson.Request:
				fmt.Sprintf(">>>>>>>>>>>>>>>>>>>>> Message source Request: %s ", msg)
			case msgjson.Notification:
				pg.noteHandlers(msg)
			case msgjson.Response:
				// client/comms.wsConn handles responses to requests we sent.
				fmt.Sprintf("A response was received in the message queue: %s", msg)
				continue
			default:
				fmt.Sprintf("Invalid message type %d from MessageSource", msg.Type)
				continue
			}

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
