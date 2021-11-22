package components

import (
	"strings"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
)

// CoinImageBySymbol returns image widget for supported asset coins.
func CoinImageBySymbol(icons *load.Icons, coinName string) *decredmaterial.Image {
	m := map[string]*decredmaterial.Image{
		"btc": icons.BTC,
		"dcr": icons.DCR,
		"bch": icons.BCH,
		"ltc": icons.LTC,
	}
	coin, ok := m[strings.ToLower(coinName)]

	if !ok {
		return nil
	}

	return coin
}
