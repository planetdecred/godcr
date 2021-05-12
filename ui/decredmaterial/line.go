package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type (
	// Line represents a rectangle widget with an initial thickness of 1
	Line struct {
		Height     int
		Width      int
		Color      color.NRGBA
		isVertical bool
	}
)

// SeparatorVertical returns a vertical line widget instance
func (t *Theme) SeparatorVertical(height, width int) Line {
	vLine := t.Line(height, width)
	vLine.isVertical = true
	return vLine
}

// Line returns a line widget instance
func (t *Theme) Line(height, width int) Line {
	if height == 0 {
		height = 1
	}

	col := t.Color.Primary
	col.A = 150
	return Line{
		Height: height,
		Width:  width,
		Color:  col,
	}
}

func (t *Theme) Separator() Line {
	l := t.Line(1, 0)
	l.Color = t.Color.Gray1
	return l
}

// Layout renders the line widget
func (l Line) Layout(gtx C) D {
	st := op.Save(gtx.Ops)
	if l.Width == 0 {
		l.Width = gtx.Constraints.Max.X
	}

	if l.isVertical && l.Height == 0 {
		l.Height = gtx.Constraints.Max.Y
	}

	line := image.Rectangle{
		Max: image.Point{
			X: l.Width,
			Y: l.Height,
		},
	}
	clip.Rect(line).Add(gtx.Ops)
	paint.Fill(gtx.Ops, l.Color)
	st.Load()
	return layout.Dimensions{Size: line.Max}
}
