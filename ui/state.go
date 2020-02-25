package ui

import (
	"github.com/raedahgroup/godcr-gio/wallet"
)

// states represents a combination of booleans that determine what the wallet is displaying.
type states struct {
	loading  bool // true if the window is in the middle of an operation that cannot be stopped
	dialog   bool // true if the window dialog modal is open
	deleted  bool // true if a wallet has been deleted
	restored bool // true if a wallet has been restored
	created  bool // true if a wallet has been created
	synced   bool // true if the mutiwallet is synced
	syncing  bool // true if the multiwallet is syncing

}

// updateStates changes the wallet state based on the received update
func (win *Window) updateStates(update interface{}) {
	win.states.loading = false
	log.Debugf("Received update %+v", update)
	switch e := update.(type) {
	case *wallet.MultiWalletInfo:
		*win.walletInfo = *e
	case wallet.SyncStarted:
		win.updateSyncStatus(true, false)
	case wallet.SyncCanceled:
		win.updateSyncStatus(false, false)
	case wallet.SyncCompleted:
		win.updateSyncStatus(false, true)
	case wallet.CreatedSeed:
		win.states.created = true
	case wallet.Restored:
		win.states.restored = true
	case wallet.LoadedWallets:
		win.wallet.GetMultiWalletInfo()
		win.states.loading = true
		win.resetInputs()
	case wallet.DeletedWallet:
		win.states.deleted = true
		win.resetInputs()
	}

	log.Debugf("Updated state %+v", win.states)
}

// reload combines the window's state to determine what widget to layout
// then invalidates the gioui window.
func (win *Window) reload() {
	log.Debugf("Reloaded with info %+v", win.walletInfo)
	current := win.WalletsPage()

	if win.states.dialog {
		current = func() {
			win.Overlay(win.dialog, current)
		}
	}

	if win.states.loading {
		current = func() {
			win.Overlay(win.Loading, current)
		}
	}
	win.current = current
	win.window.Invalidate()
}

// updateSyncStatus updates the sync status in the walletInfo state.
func (win Window) updateSyncStatus(syncing, synced bool) {
	win.walletInfo.Syncing = syncing
	win.walletInfo.Synced = synced
	win.states.synced = synced
	win.states.syncing = syncing
}
