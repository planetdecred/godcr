package window

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/page"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
)

// Window represents the app window (and UI in general). There should only be one.
type Window struct {
	window     *app.Window
	theme      *materialplus.Theme
	gtx        *layout.Context
	pages      map[string]page.Page
	current    string
	walletSync event.Duplex
	pageStates map[string]event.Event
	uiEvents   chan event.Event
}

// CreateWindow creates and initializes a new window with start
// as the first page displayed.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func CreateWindow(start string, walletSync event.Duplex) (*Window, error) {
	win := new(Window)
	win.window = app.NewWindow(app.Title("GoDcr - decred wallet"))
	win.theme = materialplus.NewTheme()
	win.gtx = layout.NewContext(win.window.Queue())

	pages := make(map[string]page.Page, 1)

	win.uiEvents = make(chan event.Event, 1) // Buffered so Loop can send and receive in the goroutine
	win.pageStates = make(map[string]event.Event)

	pages[page.LandingID] = new(page.Landing)
	pages[page.LoadingID] = new(page.Loading)

	for key, p := range pages {
		p.Init(win.theme)
		win.pageStates[key] = nil
	}

	if _, ok := pages[start]; !ok {
		return nil, fmt.Errorf("no such page")
	}
	win.current = start
	win.pages = pages
	win.walletSync = walletSync
	return win, nil
}

// Loop runs main event handling and page rendering loop
func (win *Window) Loop() {
	for {
	eventSelect:
		select {
		case e := <-win.uiEvents:
			switch evt := e.(type) {
			case event.Nav:
				win.current = evt.Next
			case error:
				// TODO: display error
			}

		case e := <-win.walletSync.Receive:
			switch evt := e.(type) {
			case event.WalletResponse:
				switch evt.Resp {
				case event.LoadedWalletsResp:
					loaded, err := evt.Results.PopInt()
					if err != nil {
						// Log or display error
						break eventSelect
					}
					if loaded == 0 {
						win.uiEvents <- event.Nav{
							Next: page.LandingID,
						}
					} // else overview

				default:
					// Unhandled Response
				}
			case error:
				// TODO: display error
			default:
				// How tho?
			}

		case e := <-win.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				win.walletSync.Send <- event.WalletCmd{
					Cmd: event.ShutdownCmd,
				}
				return
			case system.FrameEvent:
				fmt.Println("Frame")
				win.gtx.Reset(evt.Config, evt.Size)
				win.uiEvents <- win.pages[win.current].Draw(win.gtx, win.pageStates[win.current])
				evt.Frame(win.gtx.Ops)
			case nil:
				// Ignore
			default:
				fmt.Printf("Unhandled window event %+v\n", e)
			}
		}
	}
}
