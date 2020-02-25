package materialplus

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

// PaintArea paints an area with the given color and dimensions
func PaintArea(gtx *layout.Context, col color.RGBA, x, y int) {
	dim := image.Point{
		X: x,
		Y: y,
	}

	rect := f32.Rectangle{
		Max: f32.Point{
			X: float32(dim.X),
			Y: float32(dim.Y),
		},
	}

	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: rect}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: dim}
}
