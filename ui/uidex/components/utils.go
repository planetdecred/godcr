package components

import (
	"fmt"
	"math"
	"time"

	"decred.org/dcrdex/client/core"
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
	var seconds = math.Floor(float64(time.Now().Unix()*1000 - int64(t)))

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
	s, seconds = timeMod(seconds, 1000)
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

	return ""
}
