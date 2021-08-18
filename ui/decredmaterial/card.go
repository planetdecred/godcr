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
	Color  color.NRGBA
	Radius CornerRadius
}

type CornerRadius struct {
	TopLeft     float32
	TopRight    float32
	BottomRight float32
	BottomLeft  float32
}

func Radius(radius float32) CornerRadius {
	return CornerRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomRight: radius,
		BottomLeft:  radius,
	}
}

const (
	defaultRadius = 14
)

func (t *Theme) Card() Card {
	return Card{
		Color:  t.Color.Surface,
		Radius: Radius(defaultRadius),
	}
}

func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := c.Inset.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				tr := float32(gtx.Px(unit.Dp(c.Radius.TopRight)))
				tl := float32(gtx.Px(unit.Dp(c.Radius.TopLeft)))
				br := float32(gtx.Px(unit.Dp(c.Radius.BottomRight)))
				bl := float32(gtx.Px(unit.Dp(c.Radius.BottomLeft)))
				clip.RRect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(gtx.Constraints.Min.X),
						Y: float32(gtx.Constraints.Min.Y),
					}},
					NW: tl, NE: tr, SE: br, SW: bl,
				}.Add(gtx.Ops)
				return fill(gtx, c.Color)
			}),
			layout.Stacked(w),
		)
	})
	return dims
}
