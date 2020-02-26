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
	log.Debugf("Received update %+v", update)

	switch e := update.(type) {
	case wallet.SyncCompleted, wallet.SyncCanceled, wallet.SyncStarted:
		switch e.(type) {
		case wallet.SyncStarted:
			win.updateSyncStatus(true, false)
		case wallet.SyncCanceled:
			win.updateSyncStatus(false, false)
		case wallet.SyncCompleted:
			win.updateSyncStatus(false, true)
		}
		return
	case wallet.MultiWalletInfo:
		*win.walletInfo = e
		win.states.loading = false
		return
	case wallet.CreatedSeed:
		win.current = win.WalletsPage
		win.states.dialog = false
	case wallet.Restored:
		win.current = win.WalletsPage
		win.states.dialog = false
	case wallet.DeletedWallet:
		win.selected = 0
		win.states.dialog = false
	}
	win.states.loading = true
	win.wallet.GetMultiWalletInfo()

	log.Debugf("Updated with multiwallet info: %+v\n and window state %+v", win.walletInfo, win.states)
}

// updateSyncStatus updates the sync status in the walletInfo state.
func (win Window) updateSyncStatus(syncing, synced bool) {
	win.walletInfo.Syncing = syncing
	win.walletInfo.Synced = synced
	win.states.synced = synced
	win.states.syncing = syncing
}
