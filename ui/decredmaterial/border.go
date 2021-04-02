package decredmaterial

import (
	"image"
	"image/color"

	// "gioui.org/f32"
	// "gioui.org/layout"
	// "gioui.org/op"
	// "gioui.org/op/clip"
	// "gioui.org/op/paint"
	// "gioui.org/unit"

	// "gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Border struct {
	Color                    color.NRGBA
	CornerRadius             unit.Value
	Width                    unit.Value
	Top, Bottom, Left, Right bool
}

func (t *Theme) Border(top, bottom, left, right bool, cornerRadius, width unit.Value) Border {
	if width == unit.Dp(0) {
		width = unit.Dp(1)
	}

	col := t.Color.Background
	col.A = 150
	return Border{
		Color:        col,
		CornerRadius: cornerRadius,
		Width:        width,
		Top:          top,
		Bottom:       bottom,
		Left:         left,
		Right:        right,
	}
}

func (b Border) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := w(gtx)
	sz := dims.Size

	// rr := float32(gtx.Px(b.CornerRadius))
	// st := op.Save(gtx.Ops)

	width := gtx.Px(b.Width)
	sz.X -= width
	sz.Y -= width

	// clip.Border{
	// 	Rect: f32.Rectangle{
	// 		Max: layout.FPt(sz),
	// 	},
	// 	NE: rr, NW: rr, SE: rr, SW: rr,
	// 	Width: float32(width),
	// }.Add(gtx.Ops)

	// paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
	// paint.PaintOp{}.Add(gtx.Ops)

	// st.Load()

	st := op.Save(gtx.Ops)
	if b.Width == unit.Dp(0) {
		b.Width = unit.Dp(1)
	}

	gtx.Constraints.Max.Y = int(gtx.Metric.PxPerDp) * 150
	line1 := image.Rectangle{
		Max: image.Point{
			X: int(gtx.Metric.PxPerDp) * 150,
			Y: width,
		},
	}
	line2 := image.Rectangle{
		Max: image.Point{
			X: width,
			Y: int(gtx.Metric.PxPerDp) * 300,
		},
	}
	clip.Rect(line1).Add(gtx.Ops)
	clip.Rect(line2).Add(gtx.Ops)
	paint.Fill(gtx.Ops, b.Color)
	st.Load()

	return dims
}
