package materialplus

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/raedahgroup/godcr-gio/ui"
)

type (
	// Line represents a rectangle widget with an initial thickness of 1
	Line struct {
		Height float32
		Color  color.RGBA
	}
)

// Line returns a line widget instance
func (t *Theme) Line() *Line {
	col := ui.GrayColor
	col.A = 150

	return &Line{
		Height: 1,
		Color:  col,
	}
}

// Draw renders the line widget
func (l *Line) Draw(ctx *layout.Context) {
	rect := f32.Rectangle{
		Max: f32.Point{
			X: float32(ctx.Constraints.Width.Max),
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
