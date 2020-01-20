package window

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/page"
)

// Window represents the app window (and UI in general). There should only be one.
type Window struct {
	window     *app.Window
	theme      *material.Theme
	gtx        *layout.Context
	pages      map[string]page.Page
	current    string
	walletSync event.Duplex
	pageEvt    chan event.Event
}

// CreateWindow creates and initializes a new window with start
// as the first page displayed.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func CreateWindow(start string, walletSync event.Duplex) (*Window, error) {
	win := new(Window)
	win.window = app.NewWindow(app.Title("GoDcr - decred wallet"))
	win.theme = material.NewTheme()
	win.gtx = layout.NewContext(win.window.Queue())

	pages := make(map[string]page.Page, 1)

	win.pageEvt = make(chan event.Event, 2) // Buffered so Loop can send and receive in the goroutine

	pages[page.LandingID] = new(page.Landing)
	pages[page.LoadingID] = new(page.Loading)

	for _, p := range pages {
		p.Init(win.theme)
	}

	if _, ok := pages[start]; !ok {
		return nil, fmt.Errorf("No such page")
	}
	win.current = start
	win.pages = pages
	win.walletSync = walletSync
	return win, nil
}

// Loop runs main event handling and page rendering loop
func (win *Window) Loop() {
	for {
		select {
		case e := <-win.pageEvt:
			switch evt := e.(type) {
			case event.Nav:
				win.current = evt.Next
			case error:
				// TODO: display error
			}
		case e := <-win.walletSync.Receive:
			switch evt := e.(type) {
			case event.Loaded:
				if evt.WalletsLoadedCount == 0 {
					win.current = page.LandingID
				}
			}
		case e := <-win.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
				return
			case system.FrameEvent:
				win.gtx.Reset(evt.Config, evt.Size)
				// TODO: send events
				if pageEvt := win.pages[win.current].Draw(win.gtx, nil); pageEvt != nil {
					win.pageEvt <- pageEvt
				}

				evt.Frame(win.gtx.Ops)
			}
		}
	}
}
