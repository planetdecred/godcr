package widgets

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

func fill(ctx *layout.Context, col color.RGBA, x, y int) {
	d := image.Point{X: x, Y: y}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(ctx.Ops)
	paint.PaintOp{Rect: dr}.Add(ctx.Ops)
	ctx.Dimensions = layout.Dimensions{Size: d}
}
