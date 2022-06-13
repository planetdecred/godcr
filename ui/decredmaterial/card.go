package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Card struct {
	layout.Inset
	Color      color.NRGBA
	HoverColor color.NRGBA
	Radius     CornerRadius
}

type CornerRadius struct {
	TopLeft     int
	TopRight    int
	BottomRight int
	BottomLeft  int
}

func Radius(radius int) CornerRadius {
	return CornerRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomRight: radius,
		BottomLeft:  radius,
	}
}

func TopRadius(radius int) CornerRadius {
	return CornerRadius{
		TopLeft:  radius,
		TopRight: radius,
	}
}

func BottomRadius(radius int) CornerRadius {
	return CornerRadius{
		BottomRight: radius,
		BottomLeft:  radius,
	}
}

const (
	defaultRadius = 14
)

func (t *Theme) Card() Card {
	return Card{
		Color:      t.Color.Surface,
		HoverColor: t.Color.Gray4,
		Radius:     Radius(defaultRadius),
	}
}

func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := c.Inset.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				tr := gtx.Dp(unit.Dp(c.Radius.TopRight))
				tl := gtx.Dp(unit.Dp(c.Radius.TopLeft))
				br := gtx.Dp(unit.Dp(c.Radius.BottomRight))
				bl := gtx.Dp(unit.Dp(c.Radius.BottomLeft))
				defer clip.RRect{
					Rect: image.Rectangle{Max: image.Point{
						X: gtx.Constraints.Min.X,
						Y: gtx.Constraints.Min.Y,
					}},
					NW: tl, NE: tr, SE: br, SW: bl,
				}.Push(gtx.Ops).Pop()
				return fill(gtx, c.Color)
			}),
			layout.Stacked(w),
		)
	})

	return dims
}

func (c Card) HoverableLayout(gtx layout.Context, btn *Clickable, w layout.Widget) layout.Dimensions {
	background := c.Color
	dims := c.Inset.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				tr := gtx.Dp(unit.Dp(c.Radius.TopRight))
				tl := gtx.Dp(unit.Dp(c.Radius.TopLeft))
				br := gtx.Dp(unit.Dp(c.Radius.BottomRight))
				bl := gtx.Dp(unit.Dp(c.Radius.BottomLeft))
				defer clip.RRect{
					Rect: image.Rectangle{Max: image.Point{
						X: gtx.Constraints.Min.X,
						Y: gtx.Constraints.Min.Y,
					}},
					NW: tl, NE: tr, SE: br, SW: bl,
				}.Push(gtx.Ops).Pop()

				if btn.Hoverable && btn.button.Hovered() {
					background = btn.style.HoverColor
				}

				return fill(gtx, background)
			}),
			layout.Stacked(w),
		)
	})

	return dims
}

func (c Card) GradientLayout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := c.Inset.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				tr := gtx.Dp(unit.Dp(c.Radius.TopRight))
				tl := gtx.Dp(unit.Dp(c.Radius.TopLeft))
				br := gtx.Dp(unit.Dp(c.Radius.BottomRight))
				bl := gtx.Dp(unit.Dp(c.Radius.BottomLeft))

				dr := image.Rectangle{Max: gtx.Constraints.Min}

				paint.LinearGradientOp{
					Stop1:  layout.FPt(dr.Min),
					Stop2:  layout.FPt(dr.Max),
					Color1: color.NRGBA{R: 0x10, G: 0xff, B: 0x10, A: 0xFF},
					Color2: color.NRGBA{R: 0x10, G: 0x10, B: 0xff, A: 0xFF},
				}.Add(gtx.Ops)
				defer clip.RRect{
					Rect: dr,
					NW:   tl, NE: tr, SE: br, SW: bl,
				}.Push(gtx.Ops).Pop()
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{
					Size: gtx.Constraints.Max,
				}
			}),
			layout.Stacked(w),
		)
	})

	return dims
}
