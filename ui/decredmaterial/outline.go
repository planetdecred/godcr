package decredmaterial

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

type Outline struct {
	BorderColor color.RGBA
	Weight      int
}

func (t *Theme) Outline() Outline {
	return Outline{
		BorderColor: t.Color.Primary,
		Weight:      2,
	}
}

func (o Outline) Layout(gtx layout.Context, w layout.Widget) D {
	var minHeight int

	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			borderRadius := float32(gtx.Px(unit.Dp(4)))
			clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
				NE: borderRadius, NW: borderRadius, SE: borderRadius, SW: borderRadius,
			}.Add(gtx.Ops)
			fill(gtx, o.BorderColor)
			minHeight = gtx.Constraints.Min.Y

			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.Y = minHeight - o.Weight
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					clip.RRect{
						Rect: f32.Rectangle{Max: f32.Point{
							X: float32(gtx.Constraints.Min.X),
							Y: float32(gtx.Constraints.Min.Y),
						}},
						NE: borderRadius, NW: borderRadius, SE: borderRadius, SW: borderRadius,
					}.Add(gtx.Ops)
					return fill(gtx, rgb(0xffffff))
				})
			})
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return w(gtx)
				})
			})
		}),
	)

	return dims
}
