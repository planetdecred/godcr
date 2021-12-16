package dexclient

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
	"decred.org/dcrdex/dex/calc"
	"decred.org/dcrdex/dex/msgjson"
	"decred.org/dcrdex/dex/order"
)

const (
	aYear   = 31536000000
	aMonth  = 2592000000
	aDay    = 86400000
	anHour  = 3600000
	aMinute = 60000
)

func sellString(ord *core.Order) string {
	if ord.Sell {
		return "sell"
	}

	return "buy"
}

func typeString(ord *core.Order) string {
	if ord.Type != order.LimitOrderType {
		return "market"
	}

	if ord.TimeInForce == order.ImmediateTiF {
		return "limit (i)"
	}

	return "limit"
}

func rateString(ord *core.Order) string {
	if ord.Type == order.MarketOrderType {
		return "market"
	}
	// return ord.Type == Market ? 'market' : Doc.formatCoinValue(ord.rate / 1e8)
	return formatCoinValue(ord.Rate)
}

// formatCoinValue formats the asset value to a string.
func formatCoinValue(val uint64) string {
	return fmt.Sprintf("%.6f", float64(val)/1e8)
}

// timeSince returns a string representation of the duration since the specified
// unix timestamp.
func timeSince(t uint64) string {
	var seconds = float64(time.Now().Unix()*1000 - int64(t))

	var result = ""
	var count = 0

	add := func(n float64, s string) bool {
		if n > 0 || count > 0 {
			count++
		}
		if n > 0 {
			result += fmt.Sprintf("%d %s ", int(n), s)
		}
		return count >= 2
	}

	var y, mo, d, h, m, s float64

	y, seconds = timeMod(seconds, aYear)
	if add(y, "y") {
		return result
	}
	mo, seconds = timeMod(seconds, aMonth)
	if add(mo, "mo") {
		return result
	}
	d, seconds = timeMod(seconds, aDay)
	if add(d, "d") {
		return result
	}
	h, seconds = timeMod(seconds, anHour)
	if add(h, "h") {
		return result
	}
	m, seconds = timeMod(seconds, aMinute)
	if add(m, "m") {
		return result
	}
	s, _ = timeMod(seconds, 1000)
	add(s, "s")
	if result == "" {
		return "0 s"
	}
	return result
}

// timeMod returns the quotient and remainder of t / dur.
func timeMod(t float64, dur float64) (float64, float64) {
	n := math.Floor(t / dur)
	return n, t - n*dur
}

// isMarketBuy will return true if the order is a market buy order.
func isMarketBuy(ord *core.Order) bool {
	return ord.Type == order.MarketOrderType && !ord.Sell
}

// settled sums the quantities of the matches that have completed.
func settled(ord *core.Order) float64 {
	if ord.Matches == nil {
		return 0
	}
	var qty func(m *core.Match) float64

	if isMarketBuy(ord) {
		qty = func(m *core.Match) float64 {
			return (float64(m.Qty*m.Rate) * 1e-8)
		}
	} else {
		qty = func(m *core.Match) float64 {
			return float64(m.Qty)
		}
	}

	var settl float64 = 0
	for _, match := range ord.Matches {
		if match.IsCancel {
			continue
		}
		redeemed := (match.Side == order.Maker && match.Status >= order.MakerRedeemed) ||
			(match.Side == order.Taker && match.Status >= order.MatchComplete)

		if redeemed {
			settl += settl + qty(match)
		} else {
			settl += settl
		}
	}

	return settl
}

// hasLiveMatches returns true if the order has matches that have not completed
// settlement yet.
func hasLiveMatches(ord *core.Order) bool {
	if ord.Matches == nil {
		return false
	}

	for _, m := range ord.Matches {
		if !m.Revoked && m.Status < order.MakerRedeemed {
			return true
		}
	}
	return false
}

// statusString converts the order status to a string
func statusString(ord *core.Order) string {
	isLive := hasLiveMatches(ord)

	switch ord.Status {
	case order.OrderStatusUnknown:
		return "unknown"
	case order.OrderStatusEpoch:
		return "epoch"
	case order.OrderStatusBooked:
		if ord.Cancelling {
			return "cancelling"
		}
		return "booked"
	case order.OrderStatusExecuted:
		if isLive {
			return "settling"
		}
		return "executed"
	case order.OrderStatusCanceled:
		if isLive {
			return "canceled/settling"
		}
		return "canceled"
	case order.OrderStatusRevoked:
		if isLive {
			return "revoked/settling"
		}
		return "revoked"
	}

	return "unknown"
}

// minifyOrder creates a MiniOrder from a TradeNote. The epoch and order ID must
// be supplied.
func minifyOrder(oid dex.Bytes, trade *msgjson.TradeNote, epoch uint64, marketID string) (*core.MiniOrder, error) {
	b, q, err := marketIDToAsset(marketID)
	if err != nil {
		return nil, err
	}
	return &core.MiniOrder{
		Qty:       float64(trade.Quantity) / float64(b.UnitInfo.Conventional.ConversionFactor),
		QtyAtomic: trade.Quantity,
		Rate:      calc.ConventionalRate(trade.Rate, b.UnitInfo, q.UnitInfo),
		MsgRate:   trade.Rate,
		Sell:      trade.Side == msgjson.SellOrderNum,
		Token:     token(oid),
		Epoch:     epoch,
	}, nil
}

func marketIDToAsset(marketID string) (baseInfo *asset.WalletInfo, quoteInfo *asset.WalletInfo, err error) {
	mktIDs := strings.Split(marketID, "_")
	baseID, ok := dex.BipSymbolID(mktIDs[0])
	if !ok {
		return nil, nil, errors.New("Invalid market")
	}
	b, err := asset.Info(baseID)
	if err != nil {
		return nil, nil, err
	}

	quoteID, ok := dex.BipSymbolID(mktIDs[1])
	if !ok {
		return nil, nil, errors.New("Invalid market")
	}
	q, err := asset.Info(quoteID)
	if err != nil {
		return nil, nil, err
	}

	return b, q, nil
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

// removeOrder removes an order from the order book.
func removeOrder(orID dex.Bytes, sells, buys []*core.MiniOrder) ([]*core.MiniOrder, []*core.MiniOrder) {
	token := token(orID)
	if s, ok := removeFromSide(sells, token); ok {
		return s, buys
	}
	b, _ := removeFromSide(buys, token)
	return sells, b
}

// removeFromSide removes an order from the list of orders
func removeFromSide(side []*core.MiniOrder, token string) ([]*core.MiniOrder, bool) {
	ord, index := findOrder(side, token)
	if ord != nil {
		return append(side[:index], side[index+1:]...), true
	}
	return side, false
}

// findOrder finds an order in a specified side
func findOrder(side []*core.MiniOrder, token string) (*core.MiniOrder, int) {
	for i, s := range side {
		if s.Token == token {
			return s, i
		}
	}
	return nil, -1
}

func minMaxRateOrderBook(orders []*core.MiniOrder) (float64, float64) {
	if len(orders) == 0 {
		return 0, 0
	}
	max := orders[0].Rate
	min := orders[0].Rate
	for _, value := range orders {
		if max < value.Rate {
			max = value.Rate
		}
		if min > value.Rate {
			min = value.Rate
		}
	}
	return min, max
}
