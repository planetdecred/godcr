package uidex

import (
	"strings"

	"github.com/planetdecred/godcr/dexc"
)

// states represents a combination of booleans that determine what the wallet is displaying.
type states struct {
	loading  bool // true if the window is in the middle of an operation that cannot be stopped
	creating bool // true if a wallet is being created or restored
}

// updateStates changes the dexc state based on the received update
func (d *DexUI) updateStates(update interface{}) {
	switch e := update.(type) {
	case dexc.User:
		d.userInfo = &e
		if d.userInfo.Info.Exchanges == nil {
			return
		}

		// Set default selected host and market
		if d.market.host == "" ||
			d.market.marketBase == "" ||
			d.market.marketQuote == "" {
			for _, exchange := range d.userInfo.Info.Exchanges {
				d.market.host = exchange.Host
				for _, market := range exchange.Markets {
					d.market.marketBase = market.BaseSymbol
					d.market.marketQuote = market.QuoteSymbol
					d.market.marketBaseID = market.BaseID
					d.market.marketQuoteID = market.QuoteID
					d.market.name = strings.ToUpper(market.BaseSymbol + "-" + market.QuoteSymbol)
					break
				}
				break
			}
		}
		d.refresh()
		return

	case dexc.MaxOrderEstimate:
		d.maxOrderEstimate = &e
	}
}
