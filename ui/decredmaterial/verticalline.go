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
	VLine struct {
		Height int
		Width  int
		Color  color.NRGBA
	}
)

// Line returns a line widget instance
func (t *Theme) VLine(height, width int) VLine {
	if width == 0 {
		width = 1
	}

	col := t.Color.Primary
	col.A = 150
	return VLine{
		Height: height,
		Width:  width,
		Color:  col,
	}
}

func (t *Theme) VSeparator() VLine {
	l := t.VLine(0, 1)
	l.Color = t.Color.Gray1
	return l
}

// Layout renders the line widget
func (l VLine) Layout(gtx C) D {
	st := op.Save(gtx.Ops)
	if l.Height == 0 {
		l.Height = 30
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
