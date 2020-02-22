package ui

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

func FillWithColor(gtx *layout.Context, color color.RGBA, blocked bool) {
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{
		Max: f32.Point{
			X: float32(gtx.Constraints.Width.Max),
			Y: float32(gtx.Constraints.Height.Max),
		},
	}}.Add(gtx.Ops)
	if blocked {
		new(widget.Button).Layout(gtx)
	}
}
