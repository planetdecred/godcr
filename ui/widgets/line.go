package widgets

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
		height float32
		color  color.RGBA
	}
)

// NewLine returns a line widget instance
func NewLine() *Line {
	col := ui.GrayColor
	col.A = 150

	return &Line{
		height: 1,
		color:  col,
	}
}

// SetColor sets the color of the line widget
func (l *Line) SetColor(color color.RGBA) *Line {
	l.color = color
	return l
}

// SetHeight sets the line thickness
func (l *Line) SetHeight(height float32) *Line {
	l.height = height
	return l
}

// Draw renders the line widget
func (l *Line) Draw(ctx *layout.Context) {
	rect := f32.Rectangle{
		Max: f32.Point{
			X: float32(ctx.Constraints.Width.Max),
			Y: l.height,
		},
	}
	op.TransformOp{}.Offset(f32.Point{
		X: 0,
		Y: 0,
	}).Add(ctx.Ops)
	paint.ColorOp{Color: l.color}.Add(ctx.Ops)
	paint.PaintOp{Rect: rect}.Add(ctx.Ops)
}
