package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// Loading lays out the loading widget with a faded background
func (win *Window) Loading() {

	win.theme.Surface(win.gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.Center.Layout(win.gtx, func() {
					lbl := win.theme.H1("Loading")
					lbl.Color = win.theme.Color.Primary
					lbl.Layout(win.gtx)
				})
			}),
			layout.Rigid(win.Err),
		)
	})

	new(widget.Button).Layout(win.gtx)
}
