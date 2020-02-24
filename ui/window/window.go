package window

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr-gio/ui"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/page"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// Window represents the app window (and UI in general). There should only be one.
type Window struct {
	window     *app.Window
	theme      *materialplus.Theme
	gtx        *layout.Context
	pages      map[string]page.Page
	current    string
	wallet     *wallet.Wallet
	states     map[string]interface{}
	uiEvents   chan interface{}
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
func CreateWindow(start string, wal *wallet.Wallet) (*Window, error) {
	win := new(Window)
	win.window = app.NewWindow(app.Title("GoDcr - decred wallet"))
	win.theme = materialplus.NewTheme(ui.DecredPalette)
	win.gtx = layout.NewContext(win.window.Queue())

	pages := make(map[string]page.Page)

	win.uiEvents = make(chan interface{}, 2) // Buffered so Loop can send and receive in the goroutine

	win.states = make(map[string]interface{})
	pages[page.LandingID] = new(page.Landing)
	pages[page.LoadingID] = new(page.Loading)
	pages[page.WalletsID] = new(page.Wallets)
	//pages[page.UITestID] = new(page.UITest)

	win.walletInfo = new(wallet.MultiWalletInfo)
	win.states[page.StateWalletInfo] = win.walletInfo
	for _, p := range pages {
		p.Init(win.theme, wal, win.states)
	}

	if _, ok := pages[start]; !ok {
		return nil, fmt.Errorf("no such page")
	}

	win.current = start
	win.pages = pages
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
		win.states[page.StateWalletCreated] = t
	case wallet.DeletedWallet:
		win.states[page.StateDeletedWallet] = t
		win.wallet.GetMultiWalletInfo()
	}
}

// updateSyncStatus updates the sync status in the walletInfo state.
func (win Window) updateSyncStatus(syncing, synced bool) {
	win.walletInfo.Syncing = syncing
	win.walletInfo.Synced = synced
}
