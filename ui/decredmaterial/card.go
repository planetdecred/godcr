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
	Color       color.RGBA
	Rounded     bool
	CornerStyle CornerStyle
}

type CornerStyle uint8

const (
	SquareEdge CornerStyle = iota
	HalfRoundedEdge
	RoundedEdge
)

const (
	cardRadius = 10
)

func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	topRadius := 0
	bottomRadius := 0

	switch c.CornerStyle {
	case RoundedEdge:
		topRadius = cardRadius
		bottomRadius = cardRadius

	case HalfRoundedEdge:
		topRadius = 0
		bottomRadius = cardRadius

	case SquareEdge:
		topRadius = 0
		bottomRadius = 0
	}

	dims := layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return c.Inset.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						tr := float32(gtx.Px(unit.Dp(float32(topRadius))))
						br := float32(gtx.Px(unit.Dp(float32(bottomRadius))))
						clip.RRect{
							Rect: f32.Rectangle{Max: f32.Point{
								X: float32(gtx.Constraints.Min.X),
								Y: float32(gtx.Constraints.Min.Y),
							}},
							NE: tr, NW: tr, SE: br, SW: br,
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
