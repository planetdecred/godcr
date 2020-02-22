package ui

import (
	"image/color"

	"gioui.org/layout"
)

// Center is shorthand for layout.Align(layout.Center).Layout(gtx, widget)
func Center(gtx *layout.Context, widget func()) {
	layout.Align(layout.Center).Layout(gtx, widget)
}

// LayoutWithBackground renders widget Stacked in front of a background filled with color
func LayoutWithBackground(gtx *layout.Context, color color.RGBA, block bool, widget func()) {
	wmin := gtx.Constraints.Width.Min
	hmin := gtx.Constraints.Height.Min
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			FillWithColor(gtx, color, block)
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = wmin
			gtx.Constraints.Height.Min = hmin
			widget()
		}),
	)
}
