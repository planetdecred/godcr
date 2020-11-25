package decredmaterial

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

type Card struct {
	layout.Inset
	Color   color.RGBA
	Rounded bool
}

const (
	cardRadius = 10
)

func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	radius := 0
	if c.Rounded {
		radius = cardRadius
	}

	dims := layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return c.Inset.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						rr := float32(gtx.Px(unit.Dp(float32(radius))))
						clip.RRect{
							Rect: f32.Rectangle{Max: f32.Point{
								X: float32(gtx.Constraints.Min.X),
								Y: float32(gtx.Constraints.Min.Y),
							}},
							NE: rr, NW: rr, SE: rr, SW: rr,
						}.Add(gtx.Ops)
						return fill(gtx, c.Color)
					}),
					layout.Stacked(w),
				)
			})
		}),
	)
	return dims
}
