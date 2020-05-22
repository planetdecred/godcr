package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	//"gioui.org/op"
	"gioui.org/op/paint"
)

type (
	// Line represents a rectangle widget with an initial thickness of 1
	Line struct {
		Height int
		Width  int
		Color  color.RGBA
	}
)

// Line returns a line widget instance
func (t *Theme) Line() *Line {
	col := t.Color.Primary
	col.A = 150

	return &Line{
		Height: 1,
		Color:  col,
	}
}

// Layout renders the line widget
func (l *Line) Layout(gtx *layout.Context) {
	paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{
		Max: f32.Point{
			X: float32(l.Width),
			Y: float32(l.Height),
		},
	}}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{
		Size: image.Point{X: l.Width, Y: l.Height},
	}
}
