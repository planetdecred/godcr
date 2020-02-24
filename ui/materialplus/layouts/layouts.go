package layouts

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

// FillWithColor paints a color on the contexts max contraints.
// Restores the dimensions after painting.
func FillWithColor(gtx *layout.Context, color color.RGBA) {
	cs := gtx.Constraints
	dmin := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	rect := f32.Rectangle{
		Max: f32.Point{X: float32(cs.Width.Max), Y: float32(cs.Height.Max)},
	}
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: rect}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: dmin}
	op.InvalidateOp{}.Add(gtx.Ops)
}
