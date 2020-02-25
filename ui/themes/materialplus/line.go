package materialplus

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/raedahgroup/godcr-gio/ui"
)

// Line represents a rectangle widget with an initial thickness of 1
type Line struct {
	Height float32
	Width  float32
	Color  color.RGBA
}

// Line returns a line widget instance
func (t *Theme) Line() *Line {
	col := ui.GrayColor
	col.A = 150

	return &Line{
		Height: 1,
		Width:  10,
		Color:  col,
	}
}

// Draw renders the line widget
func (l *Line) Layout(ctx *layout.Context) {
	rect := f32.Rectangle{
		Max: f32.Point{
			X: l.Width,
			Y: l.Height,
		},
	}
	op.TransformOp{}.Offset(f32.Point{
		X: 0,
		Y: 0,
	}).Add(ctx.Ops)
	paint.ColorOp{Color: l.Color}.Add(ctx.Ops)
	paint.PaintOp{Rect: rect}.Add(ctx.Ops)
}
