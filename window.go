package main

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/page"
)

// window represents the app window. There should only be one.
type window struct {
	window     *app.Window
	theme      *material.Theme
	gtx        *layout.Context
	pages      map[string]page.Page
	current    string
	walletSync chan int
}

// createWindow creates and initializes a new window.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func createWindow(start string, walletSync chan int) (*window, error) {
	win := new(window)
	win.window = app.NewWindow()
	win.theme = material.NewTheme()
	win.gtx = layout.NewContext(win.window.Queue())

	pages := make(map[string]page.Page, 1)

	pages[page.LandingID] = new(page.Landing)
	pages[page.LoadingID] = new(page.Loading)

	for _, p := range pages {
		p.Init(win.theme, win.gtx)
	}

	if _, ok := pages[start]; !ok {
		return nil, fmt.Errorf("No such page")
	}
	win.current = start
	win.pages = pages
	win.walletSync = walletSync
	return win, nil
}

// loop is the main event loop
func (win *window) loop() {
	for {
		select {
		case <-win.walletSync:
			win.current = page.LandingID
		case e := <-win.window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return
			case system.FrameEvent:
				win.gtx.Reset(e.Config, e.Size)
				win.pages[win.current].Draw()
				e.Frame(win.gtx.Ops)
			}
		}
	}
}
