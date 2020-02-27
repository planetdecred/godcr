package ui

import (
	"errors"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// Window represents the app window (and UI in general). There should only be one.
// Window uses an internal state of booleans to determine what the window is currently displaying.
type Window struct {
	window *app.Window
	theme  *materialplus.Theme
	gtx    *layout.Context

	wallet     *wallet.Wallet
	walletInfo *wallet.MultiWalletInfo

	current layout.Widget
	dialog  layout.Widget
	tabs    *materialplus.Tabs

	selected int
	states

	inputs

	err string

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
	theme := decredTheme()
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	win.theme = theme
	win.gtx = layout.NewContext(win.window.Queue())

	win.walletInfo = new(wallet.MultiWalletInfo)

	win.wallet = wal
	win.states.loading = true
	win.inputs.tabs = make([]*widget.Button, 0)
	win.tabs = materialplus.NewTabs()
	win.tabs.Flex.Spacing = layout.SpaceBetween
	win.current = win.WalletsPage
	win.dialog = func() {}

	win.initWidgets()
	return win, nil
}

// Loop runs main event handling and page rendering loop
func (win *Window) Loop(shutdown chan int) {
	for {
		select {
		case e := <-win.wallet.Send:
			if e.Err != nil {
				err := e.Err.Error()
				log.Error("Wallet Error: " + err)
				win.err = err
				if win.states.loading {
					log.Warn("Attemping to get multiwallet info")
					win.wallet.GetMultiWalletInfo()
				}
				win.window.Invalidate()
				break
			}
			log.Debugf("Updating with %+v", e.Resp)
			win.updateStates(e.Resp)

		case update := <-win.wallet.Sync:
			switch update.Stage {
			case wallet.SyncCompleted:
				win.outputs.sync = win.outputs.icons.check
				win.updateSyncStatus(false, true)
			case wallet.SyncStarted:
				win.updateSyncStatus(true, false)
			case wallet.SyncCanceled:
				win.outputs.sync = win.outputs.icons.sync
				win.updateSyncStatus(false, false)
			}

		case e := <-win.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				close(shutdown)
				return
			case system.FrameEvent:
				win.gtx.Reset(evt.Config, evt.Size)
				if len(win.inputs.tabs) != win.walletInfo.LoadedWallets {
					win.inputs.tabs = make([]*widget.Button, win.walletInfo.LoadedWallets)
					for i := range win.inputs.tabs {
						win.inputs.tabs[i] = new(widget.Button)
					}
				}
				s := win.states
				win.theme.Background(win.gtx, win.current)

				if s.loading {
					win.Loading()
				} else if s.dialog {
					win.dialog()
				}

				win.HandleInputs()
				evt.Frame(win.gtx.Ops)

			case nil:
				// Ignore
			default:
				log.Tracef("Unhandled window event %+v\n", e)
			}
		}
	}
}
