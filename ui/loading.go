package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// Loading lays out the loading widget with a faded background
func (win *Window) Loading() {

	win.theme.Surface(win.gtx, func(gtx C) D {
		return layout.Center.Layout(win.gtx, func(gtx C) D {
			lbl := win.theme.H1("Loading")
			return lbl.Layout(win.gtx)
		})
	})

	new(widget.Clickable).Layout(win.gtx)
}
