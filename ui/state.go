package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// states represents a combination of booleans that determine what the wallet is displaying.
type states struct {
	deleting bool // true if a wallet is being deleted
	loading  bool // true if the window is in the middle of an operation that cannot be stopped
	creating bool
	dialog   bool
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
		//win.states[page.StateWalletCreated] = t
	case wallet.LoadedWallets:
		win.wallet.GetMultiWalletInfo()
		win.states.loading = true
		win.resetInputs()
	case wallet.DeletedWallet:
		win.resetInputs()
		//win.states[page.StateDeletedWallet] = t
	}

	log.Debugf("Updated state %+v", win.states)
}

// reload combines the window's state to determine what widget to layout
// then invalidates the gioui window.
func (win *Window) reload() {
	log.Debugf("Reloaded with info %+v", win.walletInfo)
	current := win.WalletsPage()
	if win.states.loading {
		current = func() {
			layout.Stack{}.Layout(win.gtx,
				layout.Stacked(win.current),
				layout.Expanded(win.Loading),
			)
		}
	}
	win.current = current
	win.window.Invalidate()
}

// updateSyncStatus updates the sync status in the walletInfo state.
func (win Window) updateSyncStatus(syncing, synced bool) {
	win.walletInfo.Syncing = syncing
	win.walletInfo.Synced = synced
}
