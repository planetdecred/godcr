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
	surface         color.NRGBA
	shadowElevation float32
	shadowRadius    float32
}

func (t *Theme) Shadow() *Shadow {
	return &Shadow{
		surface:         t.Color.Surface,
		shadowRadius:    8, // raduis of the shadow
		shadowElevation: 7, // height/spread of the shadow
	}
}

func (s *Shadow) Layout(gtx C, w layout.Widget) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			sz := gtx.Constraints.Min
			offset := pxf(gtx.Metric, unit.Dp(s.shadowElevation))
			rr := float32(gtx.Px(unit.Dp(s.shadowRadius)))
			rect := f32.Rect(0, 0, float32(sz.X), float32(sz.Y))

			// shadow layers arranged from the biggest to the smallest
			gradientBox(gtx.Ops, rect, rr, offset/0.8, color.NRGBA{A: 0x3})
			gradientBox(gtx.Ops, rect, rr, offset, color.NRGBA{A: 0x4})
			gradientBox(gtx.Ops, rect, rr, offset/1.5, color.NRGBA{A: 0x7})
			gradientBox(gtx.Ops, rect, rr, offset/2.5, color.NRGBA{A: 0x10})

			return layout.Dimensions{Size: sz}

		}),
		layout.Stacked(w),
	)
}

func (s *Shadow) SetShadowRadius(radius float32) {
	s.shadowRadius = radius
}

func (s *Shadow) SetShadowElevation(elev float32) {
	s.shadowElevation = elev
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
