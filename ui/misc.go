package ui

import (
	"image/color"

	"gioui.org/layout"
)

// Center is shorthand for layout.Align(layout.Center).Layout(gtx, widget)
func Center(gtx *layout.Context, widget func()) {
	layout.Align(layout.Center).Layout(gtx, widget)
}

// LayoutWithBackGround renders widget Stacked in front of a background filled with color
func LayoutWithBackGround(gtx *layout.Context, color color.RGBA, block bool, widget func()) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			FillWithColor(gtx, color, block)
		}),
		layout.Stacked(func() {
			widget()
		}),
	)
}
