// Package page provides an interface and implementations
// for creating and using pages.
package page

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"

	"github.com/raedahgroup/godcr-gio/wallet"
)

// Page represents a single page of the app.
//
// Init creates widgets with the given theme.
//
// Draw draws the implementation's widgets to the given
// layout context with regards to the given states.
// Draw returns any window event not handled by page itself.
// Draw is only called once per frame for the active page.
type Page interface {
	Init(*materialplus.Theme, *wallet.Wallet, map[string]interface{})
	Draw(gtx *layout.Context) interface{}
}

const (
	// StateWalletInfo is the map key for the WalletInfo state
	StateWalletInfo = "walletinfo"
)
