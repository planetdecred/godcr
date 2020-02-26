package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// Loading lays out the loading widget with a faded background
func (win *Window) Loading() {
	rigids := []layout.FlexChild{layout.Rigid(func() {
		layout.Center.Layout(win.gtx, func() {
			lbl := win.theme.H1("Loading")
			lbl.Color = win.theme.Primary
			lbl.Layout(win.gtx)
		})
	},
	)}

	if win.err != "" {
		rigids = append(rigids, layout.Rigid(func() {
			lbl := win.theme.Caption(win.err)
			lbl.Color = win.theme.Danger
			lbl.Layout(win.gtx)
		}))
	}

	win.theme.Background(win.gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx, rigids...)
	})

	new(widget.Button).Layout(win.gtx)
}
