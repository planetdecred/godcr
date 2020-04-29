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

func (o Outline) Layout(gtx *layout.Context, w layout.Widget) {
	var minHeight int

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			borderRadius := float32(gtx.Px(unit.Dp(4)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: borderRadius, NW: borderRadius, SE: borderRadius, SW: borderRadius,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, o.BorderColor)
			minHeight = gtx.Constraints.Height.Min

			layout.Center.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(1)).Layout(gtx, func() {
					gtx.Constraints.Height.Min = minHeight - o.Weight
					gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
					clip.Rect{
						Rect: f32.Rectangle{Max: f32.Point{
							X: float32(gtx.Constraints.Width.Min),
							Y: float32(gtx.Constraints.Height.Min),
						}},
						NE: borderRadius, NW: borderRadius, SE: borderRadius, SW: borderRadius,
					}.Op(gtx.Ops).Add(gtx.Ops)
					fill(gtx, rgb(0xffffff))
				})
			})
		}),
		layout.Stacked(func() {
			layout.Center.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
					w()
				})
			})
		}),
	)
}
