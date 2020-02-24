package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
)

func (win *Window) Loading() {
	layouts.FillWithColor(win.gtx, layouts.Faded(win.theme.Background))
	log.Debugf("With Loading")
	layout.Center.Layout(win.gtx, func() {
		lbl := win.theme.Label(unit.Dp(100), "Loading")
		lbl.Color = win.theme.Primary
		lbl.Layout(win.gtx)
	})
	new(widget.Button).Layout(win.gtx)
}
