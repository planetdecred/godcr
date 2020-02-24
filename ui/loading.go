package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

func (win *Window) Loading() {
	log.Debug("With Loading")
	layout.Center.Layout(win.gtx, func() {
		win.theme.Label(unit.Dp(100), "Loading").Layout(win.gtx)
	})
	new(widget.Button).Layout(win.gtx)
}
