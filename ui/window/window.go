package window

import (
	"fmt"

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
	pageStates map[string][]interface{}
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
	win.pageStates = make(map[string][]interface{})

	pages[page.LandingID] = new(page.Landing)
	pages[page.LoadingID] = new(page.Loading)

	win.walletInfo = new(wallet.MultiWalletInfo)
	for key, p := range pages {
		p.Init(win.theme, wal)
		win.pageStates[key] = []interface{}{*win.walletInfo}
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
func (win *Window) Loop() {
	for {
		select {
		case e := <-win.uiEvents:
			switch evt := e.(type) {
			case page.EventNav:
				win.current = evt.Next
			case error:
				// TODO: display error
			}

		case e := <-win.wallet.Send:
			switch evt := e.(type) {
			case *wallet.LoadedWallets:
				if evt.Count == 0 {
					win.current = page.LandingID
				}
			case *wallet.MultiWalletInfo:
				win.walletInfo = evt
			case error:
				// TODO: display error
			default:
				// How tho?
			}

		case e := <-win.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				win.wallet.Shutdown()
				return
			case system.FrameEvent:
				fmt.Println("Frame")
				win.gtx.Reset(evt.Config, evt.Size)
				if pageEvt := win.pages[win.current].Draw(win.gtx, win.pageStates[win.current]...); pageEvt != nil {
					win.uiEvents <- pageEvt
				}
				evt.Frame(win.gtx.Ops)
			case nil:
				// Ignore
			default:
				fmt.Printf("Unhandled window event %+v\n", e)
			}
		}
	}
}
