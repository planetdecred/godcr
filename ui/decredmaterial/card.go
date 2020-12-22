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
	Color  color.RGBA
	Radius CornerRadius
}

type CornerRadius struct {
	NE int
	NW int
	SE int
	SW int
}

const (
	defaultRadius = 10
)

func (t *Theme) Card() Card {
	return Card{
		Color: t.Color.Surface,
		Radius: CornerRadius{
			NE: defaultRadius,
			SE: defaultRadius,
			NW: defaultRadius,
			SW: defaultRadius,
		},
	}
}

func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return c.Inset.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						tr := float32(gtx.Px(unit.Dp(float32(c.Radius.NE))))
						tl := float32(gtx.Px(unit.Dp(float32(c.Radius.NW))))
						br := float32(gtx.Px(unit.Dp(float32(c.Radius.SE))))
						bl := float32(gtx.Px(unit.Dp(float32(c.Radius.SW))))
						clip.RRect{
							Rect: f32.Rectangle{Max: f32.Point{
								X: float32(gtx.Constraints.Min.X),
								Y: float32(gtx.Constraints.Min.Y),
							}},
							NE: tl, NW: tr, SE: br, SW: bl,
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
