package decredmaterial

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Shadow struct {
	surface color.NRGBA
}

const (
	shadowElevation = 5
	shadowRadius    = 8
)

func (t *Theme) Shadow() *Shadow {
	return &Shadow{
		surface: t.Color.Surface,
	}
}

func (s *Shadow) Layout(gtx C, w layout.Widget) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			s.layoutShadow(gtx)
			surface := clip.UniformRRect(f32.Rectangle{Max: layout.FPt(gtx.Constraints.Min)}, shadowRadius)
			paint.FillShape(gtx.Ops, s.surface, surface.Op(gtx.Ops))
			return D{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}

func (s *Shadow) layoutShadow(gtx C) D {
	sz := gtx.Constraints.Min
	offset := pxf(gtx.Metric, unit.Dp(shadowElevation))
	shadowSize := float32(gtx.Px(unit.Dp(shadowElevation)))
	rect := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))

	// shadow layers arranged from the biggest to the smallest
	gradientBox(gtx.Ops, rect, shadowSize, offset/0.8, color.NRGBA{A: 0x5})
	gradientBox(gtx.Ops, rect, shadowSize, offset, color.NRGBA{A: 0x10})
	gradientBox(gtx.Ops, rect, shadowSize, offset/1.5, color.NRGBA{A: 0x15})
	gradientBox(gtx.Ops, rect, shadowSize, offset/2.5, color.NRGBA{A: 0x20})

	return layout.Dimensions{Size: sz}
}

func gradientBox(ops *op.Ops, r f32.Rectangle, rr, spread float32, col color.NRGBA) {
	paint.FillShape(ops, col, clip.RRect{
		Rect: outset(r, spread),
		SE:   rr + spread, SW: rr + spread, NW: rr + spread, NE: rr + spread,
	}.Op(ops))
}

func outset(r f32.Rectangle, rr float32) f32.Rectangle {
	r.Min.X -= rr
	r.Min.Y -= rr
	r.Max.X += rr
	r.Max.Y += rr
	return r
}

func pxf(c unit.Metric, v unit.Value) float32 {
	switch v.U {
	case unit.UnitPx:
		return v.V
	case unit.UnitDp:
		s := c.PxPerDp
		if s == 0 {
			s = 1
		}
		return s * v.V
	case unit.UnitSp:
		s := c.PxPerSp
		if s == 0 {
			s = 1
		}
		return s * v.V
	default:
		panic("unknown unit")
	}
}
