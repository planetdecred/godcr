package ui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
)

// Loading lays out the loading widget with a faded background
func (win *Window) Loading() {
	layout.Center.Layout(win.gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := win.theme.H1("Loading")
		lbl.Alignment = text.Middle
		return lbl.Layout(gtx)
	})
	new(widget.Clickable).Layout(win.gtx)
}
