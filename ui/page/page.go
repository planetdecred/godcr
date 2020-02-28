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

// Handler represents a page handler
type Handler struct {
	ID        string
	IsNavPage bool
	Page      Page
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
func GetHandlers() []Handler {
	return []Handler{
		{
			ID:        LandingID,
			IsNavPage: false,
			Page:      new(Landing),
		},
		{
			ID:        LoadingID,
			IsNavPage: false,
			Page:      new(Loading),
		},
		{
			ID:        WalletsID,
			IsNavPage: true,
			Page:      new(Wallets),
		},
		{
			ID:        UITestID,
			IsNavPage: true,
			Page:      new(UITest),
		},
	}
}
