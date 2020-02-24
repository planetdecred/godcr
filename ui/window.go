package ui

import (
	"errors"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// Window represents the app window (and UI in general). There should only be one.
type Window struct {
	window      *app.Window
	theme       *materialplus.Theme
	gtx         *layout.Context
	current     func(*layout.Context, *materialplus.Theme, *wallet.MultiWalletInfo)
	wallet      *wallet.Wallet
	walletInfo  *wallet.MultiWalletInfo
	infoLoading bool
	buttons     struct {
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
	theme := decredTheme()
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	win.theme = theme
	win.gtx = layout.NewContext(win.window.Queue())

	win.walletInfo = new(wallet.MultiWalletInfo)

	win.current = blank
	win.wallet = wal
	return win, nil
}

// Loop runs main event handling and page rendering loop
func (win *Window) Loop(shutdown chan int) {
	for {
		select {
		case e := <-win.wallet.Send:
			log.Debugf("Recieved event %+v", e)
			if e.Err != nil {
				win.window.Invalidate()
				break
			}
			switch evt := e.Resp.(type) {
			case *wallet.LoadedWallets:
				win.wallet.GetMultiWalletInfo()
				if evt.Count == 0 {
					win.current = blank
				} else {
					win.current = blank
				}
			case *wallet.MultiWalletInfo:
				*win.walletInfo = *evt
			default:
				win.updateState(e.Resp)
			}
			// set error if it exists
			if e.Err != nil {
				//win.states[page.StateError] = e.Err
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

				win.current(win.gtx, win.theme, win.walletInfo)
				if win.infoLoading {
					loading(win.gtx, win.theme, win.walletInfo)
				}
				log.Tracef("Page {%s} rendered in %v", win.current, time.Since(start))
				evt.Frame(win.gtx.Ops)
				win.HandleInputs()
			case nil:
				// Ignore
			default:
				log.Tracef("Unhandled window event %+v\n", e)
			}
		}
	}
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
		win.reloadInfo()
		//win.states[page.StateWalletCreated] = t
	case wallet.DeletedWallet:
		//win.states[page.StateDeletedWallet] = t
		win.reloadInfo()
	}
}

// updateSyncStatus updates the sync status in the walletInfo state.
func (win Window) updateSyncStatus(syncing, synced bool) {
	win.walletInfo.Syncing = syncing
	win.walletInfo.Synced = synced
}

func (win *Window) reloadInfo() {
	win.infoLoading = true
	win.wallet.GetMultiWalletInfo()
}
