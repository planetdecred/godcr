package widgets

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/raedahgroup/godcr-gio/helper"
)

type (
	Line struct {
		height float32
		color  color.RGBA
	}
)

func NewLine() *Line {
	return &Line{
		height: 1,
		color:  helper.DecredDarkBlueColor,
	}
}

func (l *Line) SetColor(color color.RGBA) *Line {
	l.color = color
	return l
}

func (l *Line) SetHeight(height float32) *Line {
	l.height = height
	return l
}

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
