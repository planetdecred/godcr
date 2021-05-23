// util contains functions that don't contain layout code. They could be considered helpers that aren't particularly
// bounded to a page.

package uidex

import (
	"strings"

	"gioui.org/gesture"
	"gioui.org/widget"
)

// createClickGestures returns a slice of click gestures
func createClickGestures(count int) []*gesture.Click {
	var gestures = make([]*gesture.Click, count)
	for i := 0; i < count; i++ {
		gestures[i] = &gesture.Click{}
	}
	return gestures
}

// coinImageBySymbol returns image widget for supported asset coins
func coinImageBySymbol(icons *pageIcons, coinName string) *widget.Image {
	m := map[string]*widget.Image{
		"btc": icons.btc,
		"dcr": icons.dcr,
		"bch": icons.bch,
		"ltc": icons.ltc,
	}
	coin, ok := m[strings.ToLower(coinName)]

	if !ok {
		return nil
	}
	return coin
}
