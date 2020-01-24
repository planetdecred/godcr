package widgets

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"image"
	"image/color"
)

// Fill draws the rectangle and adds color to it.
func Fill(ctx *layout.Context, col color.RGBA, rect f32.Rectangle) {
	d := image.Point{X: int(rect.Max.X), Y: int(rect.Max.Y)}
	paint.ColorOp{Color: col}.Add(ctx.Ops)
	paint.PaintOp{Rect: rect}.Add(ctx.Ops)
	ctx.Dimensions = layout.Dimensions{Size: d}
}