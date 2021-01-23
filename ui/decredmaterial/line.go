package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/op"
	"gioui.org/op/clip"
	//"gioui.org/f32"
	"gioui.org/layout"
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
	/**paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{
		Max: f32.Point{
			X: float32(l.Width),
			Y: float32(l.Height),
		},
	}}.Add(gtx.Ops)
	dims := image.Point{X: l.Width, Y: l.Height}**/

	st := op.Save(gtx.Ops)
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
