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
		Height int
		Width  int
		Color  color.NRGBA
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
func (l *Line) Layout(gtx C) D {
	st := op.Save(gtx.Ops)
	line := image.Rectangle{
		Max: image.Point{
			X: l.Width,
			Y: l.Height * int(gtx.Metric.PxPerDp),
		},
	}
	clip.Rect(line).Add(gtx.Ops)
	paint.Fill(gtx.Ops, l.Color)
	st.Load()
	return layout.Dimensions{Size: line.Max}
}
