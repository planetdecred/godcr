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
// Implementations can store the event Duplex for
// communication with the window
//
// Draw draws the implementation's widgets to the given
// layout context with regards to the given states.
// The given states must have a *wallet.MultiWalletInfo as the first
// element.
// Draw returns any window event not handled by page itself.
// Draw is only called once per frame for the active page.
type Page interface {
	Init(*materialplus.Theme, *wallet.Wallet)
	Draw(gtx *layout.Context, states ...interface{}) interface{}
}
