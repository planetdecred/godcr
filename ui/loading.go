package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

// Loading lays out the loading widget with a faded background
func (win *Window) Loading() {
	win.theme.Background(win.gtx, func() {
		layout.Center.Layout(win.gtx, func() {
			lbl := win.theme.Label(unit.Dp(100), "Loading")
			lbl.Color = win.theme.Primary
			lbl.Layout(win.gtx)
		})
	})

	new(widget.Button).Layout(win.gtx)
}
