package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// Loading lays out the loading widget with a faded background
func (win *Window) Loading() {

	win.theme.Surface(win.gtx, func() {
		layout.Center.Layout(win.gtx, func() {
			lbl := win.theme.H1("Loading")
			lbl.Layout(win.gtx)
		})

	})

	new(widget.Button).Layout(win.gtx)
}
