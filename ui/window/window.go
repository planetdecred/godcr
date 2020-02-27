package window

import (
	"fmt"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"

	"github.com/raedahgroup/godcr-gio/ui/page"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
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
}

// CreateWindow creates and initializes a new window with start
// as the first page displayed.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func CreateWindow(start string, wal *wallet.Wallet) (*Window, error) {
	win := new(Window)
	win.window = app.NewWindow(app.Title("GoDcr - decred wallet"))
	win.theme = materialplus.NewTheme()
	win.gtx = layout.NewContext(win.window.Queue())

	pages := make(map[string]page.Page)

	win.uiEvents = make(chan interface{}, 2) // Buffered so Loop can send and receive in the goroutine

	win.states = make(map[string]interface{})
	pages[page.LandingID] = new(page.Landing)
	pages[page.LoadingID] = new(page.Loading)
	pages[page.OverviewID] = new(page.Overview)
	pages[page.WalletsID] = new(page.Wallets)
	pages[page.UITestID] = new(page.UITest)
	pages[page.ReceivingID] = new(page.Receive)

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

// Loop runs main event handling and page rendering loop
func (win *Window) Loop(shutdown chan int) {
	for {
		select {
		case e := <-win.uiEvents:
			switch evt := e.(type) {
			case page.EventNav:
				win.current = evt.Next
			case error:
				// TODO: display error
			}
			win.window.Invalidate()
		case e := <-win.wallet.Send:
			switch evt := e.Resp.(type) {
			case *wallet.LoadedWallets:
				win.wallet.GetMultiWalletInfo()
				win.wallet.GetAllTransactions(0, 10, 0)
				if evt.Count == 0 {
					win.current = page.LandingID
				} else {
					win.current = page.WalletsID
				}
			case *wallet.MultiWalletInfo:
				*win.walletInfo = *evt
			default:
				win.updateState(e.Resp)
			}
			// set error if it exists
			if e.Err != nil {
				win.states[page.StateError] = e.Err
			}
			win.window.Invalidate()
		case e := <-win.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				close(shutdown)
				return
			case system.FrameEvent:
				win.gtx.Reset(evt.Config, evt.Size)
				start := time.Now()
				pageEvt := win.pages[win.current].Draw(win.gtx)
				log.Tracef("Page {%s} rendered in %v", win.current, time.Since(start))
				if pageEvt != nil {
					win.uiEvents <- pageEvt
				}
				evt.Frame(win.gtx.Ops)
			case nil:
				// Ignore
			default:
				//fmt.Printf("Unhandled window event %+v\n", e)
			}
		}
	}
}

// updateState checks for the event type that is passed as an argument and updates its
// respective state.
func (win *Window) updateState(t interface{}) {
	switch t := t.(type) {
	case wallet.SyncStarted:
		win.updateSyncStatus(true, false)
	case wallet.SyncCanceled:
		win.updateSyncStatus(false, false)
	case wallet.SyncCompleted:
		win.updateSyncStatus(false, true)
	case wallet.SyncHeadersFetchProgress:
		win.updateHeaderFetchProgress(t)
	case wallet.SyncPeersChanged:
		win.updateConnectedPeers(t)
	case wallet.SyncHeadersRescanProgress:
		win.updateRescanHeaderProgress(t)
	case wallet.SyncAddressDiscoveryProgress:
		win.updateAddressDiscoveryProgress(t)
	case *wallet.Transactions:
		win.updateTransactions(t)
	case *wallet.CreatedSeed:
		win.wallet.GetMultiWalletInfo()
		win.states[page.StateWalletCreated] = t
	}
}

// stateObject fetches and returns the state if it already exists in the state map.
// Otherwise, it creates a new state object and returns a pointer of that state.
func (win Window) stateObject(key string) interface{} {
	if state, ok := win.states[key]; ok {
		return state
	}
	switch key {
	case page.StateSyncStatus:
		win.states[key] = new(wallet.SyncStatus)
		return win.states[key]
	case page.StateTransactions:
		win.states[key] = new(wallet.Transactions)
		return win.states[key]
	}
	return nil
}

// updateSyncStatus updates the sync status in the walletInfo state.
func (win Window) updateSyncStatus(syncing, synced bool) {
	win.walletInfo.Syncing = syncing
	win.walletInfo.Synced = synced
}

// updateSyncProgress updates the headers fetched in the SyncStatus state
func (win Window) updateHeaderFetchProgress(resp wallet.SyncHeadersFetchProgress) {
	state := win.stateObject(page.StateSyncStatus)
	syncState := state.(*wallet.SyncStatus)
	syncState.HeadersFetchProgress = resp.Progress.HeadersFetchProgress
	syncState.HeadersToFetch = resp.Progress.TotalHeadersToFetch
	syncState.Progress = resp.Progress.TotalSyncProgress
	syncState.RemainingTime = resp.Progress.TotalTimeRemainingSeconds
	syncState.TotalSteps = wallet.TotalSyncSteps
	syncState.Steps = wallet.FetchHeadersStep
	syncState.CurrentBlockHeight = resp.Progress.CurrentHeaderHeight
	// update wallet state when new headers are fetched
	win.wallet.GetMultiWalletInfo()
}

// updateSyncProgress updates rescan Header Progress in the SyncStatus state
func (win Window) updateRescanHeaderProgress(resp wallet.SyncHeadersRescanProgress) {
	state := win.stateObject(page.StateSyncStatus)
	syncState := state.(*wallet.SyncStatus)
	syncState.RescanHeadersProgress = resp.Progress.RescanProgress
	syncState.Progress = resp.Progress.TotalSyncProgress
	syncState.RemainingTime = resp.Progress.TotalTimeRemainingSeconds
	syncState.TotalSteps = wallet.TotalSyncSteps
	syncState.Steps = wallet.RescanHeadersStep
}

// updateSyncProgress updates Address Discovery Progress in the SyncStatus state
func (win Window) updateAddressDiscoveryProgress(resp wallet.SyncAddressDiscoveryProgress) {
	state := win.stateObject(page.StateSyncStatus)
	syncState := state.(*wallet.SyncStatus)
	syncState.RescanHeadersProgress = resp.Progress.AddressDiscoveryProgress
	syncState.Progress = resp.Progress.TotalSyncProgress
	syncState.RemainingTime = resp.Progress.TotalTimeRemainingSeconds
	syncState.TotalSteps = wallet.TotalSyncSteps
	syncState.Steps = wallet.AddressDiscoveryStep
}

// updateConnectedPeers updates connected peers in the SyncStatus state
func (win Window) updateConnectedPeers(resp wallet.SyncPeersChanged) {
	state := win.stateObject(page.StateSyncStatus)
	syncState := state.(*wallet.SyncStatus)
	syncState.ConnectedPeers = resp.ConnectedPeers
}

// updateTransactionState updates the transactions state.
func (win Window) updateTransactions(resp *wallet.Transactions) {
	state := win.stateObject(page.StateTransactions)
	txState := state.(*wallet.Transactions)
	txState.Recent = resp.Recent
	txState.Txs = resp.Txs
}
