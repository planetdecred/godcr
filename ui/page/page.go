// Package page provides an interface and implementations
// for creating and using pages.
package page

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"

	"github.com/raedahgroup/godcr-gio/wallet"
)

// Handler represents a single page of the app.
//
// Init creates widgets with the given theme.
//
// Draw draws the implementation's widgets to the given
// layout context with regards to the given states.
// Draw returns any window event not handled by page itself.
// Draw is only called once per frame for the active page.
type Handler interface {
	Init(*materialplus.Theme, *wallet.Wallet, map[string]interface{})
	Draw(gtx *layout.Context) interface{}
}

// Page represents information about a page
type Page struct {
	ID        string
	IsNavPage bool
	Handler   Handler
	Button    *widget.Button
}

const (
	// StateWalletInfo is the map key for the WalletInfo state
	StateWalletInfo = "walletinfo"

	// StateSyncStatus is the map key for SyncStatus state

	// SyncStatus is the map key for SyncStatus state
	StateSyncStatus = "syncstatus"

	// StateTransactions is the map key for Transactions state
	StateTransactions = "transactions"

	// StateWalletCreated is the map key for the WalletCreated state
	StateWalletCreated = "walletCreated"

	// StateError is the map key for error
	StateError = "error"
)

// GetPages returns all pages
func GetPages() []Page {
	return []Page{
		{
			ID:        LandingID,
			IsNavPage: false,
			Handler:   new(Landing),
			Button:    new(widget.Button),
		},
		{
			ID:        OverviewID,
			IsNavPage: true,
			Handler:   new(Overview),
			Button:    new(widget.Button),
		},
		{
			ID:        LoadingID,
			IsNavPage: false,
			Handler:   new(Loading),
			Button:    new(widget.Button),
		},
		{
			ID:        WalletsID,
			IsNavPage: true,
			Handler:   new(Wallets),
			Button:    new(widget.Button),
		},
		{
			ID:        UITestID,
			IsNavPage: true,
			Handler:   new(UITest),
			Button:    new(widget.Button),
		},
	}
}
