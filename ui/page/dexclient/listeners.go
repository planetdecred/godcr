package dexclient

import (
	"encoding/hex"
	"fmt"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
	"decred.org/dcrdex/dex/msgjson"
	"github.com/planetdecred/godcr/dexc"
)

func (pg *Page) connectDex(h string, password []byte) {
	pg.DL.Dexc.ConnectDexes(h, password)
	go pg.listenerMessages()
	go pg.readNotifications()
	pg.updateOrderBook()
}

func (pg *Page) updateOrderBook() {
	pg.orderBook, _ = pg.DL.Dexc.Book(testDexHost, dexc.DefaultAssetID, 0)
	pg.miniTradeFormWdg.orderBook = pg.orderBook
}

func (pg *Page) listenerMessages() {
	msgs := pg.DL.Dexc.MessageSource(testDexHost)
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
	switch msg.Route {
	case msgjson.BookOrderRoute:
		fmt.Println(">>>>>>>>>>>>>>>>>>>>> BookOrderRoute Receive message source: Notification", msg.Route)
		pg.updateOrderBook()
	case msgjson.EpochOrderRoute:
		fmt.Println(">>>>>>>>>>>>>>>>>>>>> EpochOrderRoute Receive message source: Notification", msg.Route)
		pg.updateOrderBook()
	case msgjson.UnbookOrderRoute:
		fmt.Println(">>>>>>>>>>>>>>>>>>>>> UnbookOrderRoute Receive message source: Notification", msg.Route)
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

// minifyOrder creates a MiniOrder from a TradeNote. The epoch and order ID must
// be supplied.
func minifyOrder(oid dex.Bytes, trade *msgjson.TradeNote, epoch uint64) *core.MiniOrder {
	return &core.MiniOrder{
		Qty:   float64(trade.Quantity) / dexc.ConversionFactor,
		Rate:  float64(trade.Rate) / dexc.ConversionFactor,
		Sell:  trade.Side == msgjson.SellOrderNum,
		Token: token(oid),
		Epoch: epoch,
	}
}

// token is a short representation of a byte-slice-like ID, such as a match ID
// or an order ID. The token is meant for display where the 64-character
// hexadecimal IDs are untenable.
func token(id []byte) string {
	if len(id) < 4 {
		return ""
	}
	return hex.EncodeToString(id[:4])
}
