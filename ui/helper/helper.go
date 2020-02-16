package helper

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

const (
	// WindowHeight is the height of the main window
	WindowHeight = 500
	// WindowWidth is the width of the main window
	WindowWidth = 500
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

// Fill paints an area completely
func Fill(gtx *layout.Context, col color.RGBA) {
	PaintArea(gtx, col, gtx.Constraints.Width.Min, gtx.Constraints.Height.Min)
}
