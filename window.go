package main

import (
	"fmt"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/page"
)

const (
	landingPage = "landing"
)

type window struct {
	window  *app.Window
	theme   *material.Theme
	gtk     *layout.Context
	pages   map[string]page.Page
	current string
}

func createWindow(start string) (*window, error) {
	win := new(window)
	win.window = app.NewWindow()
	win.theme = material.NewTheme()
	win.gtk = layout.NewContext(win.window.Queue())
	pages := make(map[string]page.Page, 1)

	pages[landingPage] = new(page.Landing)

	for _, p := range pages {
		p.Init(win.theme, win.gtk)
	}

	if _, ok := pages[start]; !ok {
		return nil, fmt.Errorf("No such page")
	}
	win.current = start
	win.pages = pages
	return win, nil
}

func (win *window) loop() {
	for {
		e := <-win.window.Events()
		fmt.Println(e)
		switch e := e.(type) {
		case system.DestroyEvent:
			fmt.Println(e.Err)
			return
		case system.FrameEvent:
			win.gtk.Reset(e.Config, e.Size)
			win.pages[win.current].Draw()
			e.Frame(win.gtk.Ops)
		}
	}
}
