package window

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr-gio/ui"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// Window represents the app window (and UI in general). There should only be one.
type Window struct {
	window     *app.Window
	theme      *materialplus.Theme
	gtx        *layout.Context
	current    func(theme *materialplus.Theme, gtx *layout.Context)
	wallet     *wallet.Wallet
	walletInfo *wallet.MultiWalletInfo
	buttons    struct {
		deleteWallet, cancelDialog, confirmDialog *widget.Button
	}
}

// CreateWindow creates and initializes a new window with start
// as the first page displayed.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func CreateWindow(wal *wallet.Wallet) (*Window, error) {
	win := new(Window)
	win.window = app.NewWindow(app.Title("GoDcr - decred wallet"))
	win.theme = materialplus.NewTheme(ui.DecredPalette)
	win.gtx = layout.NewContext(win.window.Queue())

	win.walletInfo = new(wallet.MultiWalletInfo)

	win.current = Loading
	win.wallet = wal
	return win, nil
}

// updateState checks for the event type that is passed as an argument and updates its
// respective state.
func (win *Window) updateState(t interface{}) {
	switch t.(type) {
	case wallet.SyncStarted:
		win.updateSyncStatus(true, false)
	case wallet.SyncCanceled:
		win.updateSyncStatus(false, false)
	case wallet.SyncCompleted:
		win.updateSyncStatus(false, true)
	case *wallet.CreatedSeed:
		win.wallet.GetMultiWalletInfo()
		//win.states[page.StateWalletCreated] = t
	case wallet.DeletedWallet:
		//win.states[page.StateDeletedWallet] = t
		win.wallet.GetMultiWalletInfo()
	}
}

// updateSyncStatus updates the sync status in the walletInfo state.
func (win Window) updateSyncStatus(syncing, synced bool) {
	win.walletInfo.Syncing = syncing
	win.walletInfo.Synced = synced
}
