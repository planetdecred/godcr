package ui

import (
	"errors"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// Window represents the app window (and UI in general). There should only be one.
// Window uses an internal state of booleans to determine what the window is currently displaying.
type Window struct {
	window *app.Window
	theme  *decredmaterial.Theme
	gtx    *layout.Context

	wallet             *wallet.Wallet
	walletInfo         *wallet.MultiWalletInfo
	walletSyncStatus   *wallet.SyncStatus
	walletTransactions *wallet.Transactions

	current layout.Widget
	dialog  layout.Widget
	tabs    *decredmaterial.Tabs

	selectedAccount int

	selected int
	states

	inputs

	err string

	combined
	outputs
}

// CreateWindow creates and initializes a new window with start
// as the first page displayed.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func CreateWindow(wal *wallet.Wallet) (*Window, error) {
	win := new(Window)
	win.window = app.NewWindow(app.Title("GoDcr - decred wallet"))
	theme := decredmaterial.NewTheme()
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	win.theme = theme
	win.gtx = layout.NewContext(win.window.Queue())

	win.walletInfo = new(wallet.MultiWalletInfo)
	win.walletSyncStatus = new(wallet.SyncStatus)
	win.walletTransactions = new(wallet.Transactions)

	win.wallet = wal
	win.states.loading = true
	win.tabs = decredmaterial.NewTabs()
	win.tabs.Flex.Spacing = layout.SpaceBetween
	win.current = win.OverviewPage
	win.dialog = func() {}

	win.initWidgets()
	return win, nil
}

func (win *Window) unloaded() {
	lbl := win.theme.H3("Multiwallet not loaded\nIs another instance open?")
	for {
		e := <-win.window.Events()
		switch evt := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			win.gtx.Reset(evt.Config, evt.Size)
			lbl.Layout(win.gtx)
			evt.Frame(win.gtx.Ops)
		}
	}
}

// Loop runs main event handling and page rendering loop
func (win *Window) Loop(shutdown chan int) {
	for {
		select {
		case e := <-win.wallet.Send:
			if e.Err != nil {
				err := e.Err.Error()
				log.Error("Wallet Error: " + err)
				if err == dcrlibwallet.ErrWalletDatabaseInUse {
					close(shutdown)
					win.unloaded()
					return
				}
				win.err = err
				if win.states.loading {
					log.Warn("Attemping to get multiwallet info")
					win.wallet.GetMultiWalletInfo()
				}
				win.window.Invalidate()
				break
			}
			//log.Debugf("Updating with %+v", e.Resp)
			win.updateStates(e.Resp)

		case update := <-win.wallet.Sync:
			switch update.Stage {
			case wallet.SyncCompleted:
				win.outputs.syncHeader = win.outputs.icons.check
				win.updateSyncStatus(false, true)
			case wallet.SyncStarted:
				// dcrlibwallet triggers the SyncStart method several times
				// without sending a SyncComplete signal when sync is done.
				if !win.walletInfo.Synced {
					win.updateSyncStatus(true, false)
				}
			case wallet.SyncCanceled:
				win.outputs.syncHeader = win.outputs.icons.sync
				win.updateSyncStatus(false, false)
			case wallet.HeadersFetchProgress:
				win.updateSyncProgress(update.ProgressReport)
			case wallet.AddressDiscoveryProgress:
				win.updateSyncProgress(update.ProgressReport)
			case wallet.HeadersRescanProgress:
				win.updateSyncProgress(update.ProgressReport)
			case wallet.PeersConnected:
				win.updateConnectedPeers(update.ConnectedPeers)
			case wallet.BlockAttached:
				if win.walletInfo.Synced {
					win.wallet.GetMultiWalletInfo()
				}
			}

		case e := <-win.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				close(shutdown)
				return
			case system.FrameEvent:
				win.gtx.Reset(evt.Config, evt.Size)
				s := win.states

				if s.loading {
					win.Loading()
				} else if s.dialog {
					decredmaterial.Modal{Direction: layout.Center}.Layout(win.gtx, func() { win.theme.Background(win.gtx, win.current) }, win.dialog)
				} else {
					win.theme.Background(win.gtx, win.current)
				}

				win.HandleInputs()
				evt.Frame(win.gtx.Ops)
			case key.Event:
				win.KeysEventsHandler(&evt)
			case nil:
				// Ignore
			default:
				log.Tracef("Unhandled window event %+v\n", e)
			}
		}
	}
}
