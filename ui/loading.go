package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
)

// Loading lays out the loading widget with a faded background
func (win *Window) Loading() {
	materialplus.Modal{
		Background: materialplus.Faded(win.theme.Background),
		Direction:  layout.Center,
	}.Layout(win.gtx, func() {
		layout.Center.Layout(win.gtx, func() {
			lbl := win.theme.Label(unit.Dp(100), "Loading")
			lbl.Color = win.theme.Primary
			lbl.Layout(win.gtx)
		})
	})

	new(widget.Button).Layout(win.gtx)
}
