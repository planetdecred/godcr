package decredmaterial

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

const (
	WrapContent = -1
	MatchParent = -2
)

type LinearLayout struct {
	Width       int
	Height      int
	Orientation layout.Axis
	Background  color.NRGBA
	Border      Border
	Margin      layout.Inset
	Padding     layout.Inset
	Direction   layout.Direction
}

func (ll LinearLayout) Layout(gtx C, children ...layout.FlexChild) D {
	// draw margin
	return ll.Margin.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				ll.applyDimension(&gtx)
				// // draw background and and clip the background to border radius
				tr := float32(gtx.Px(unit.Dp(ll.Border.Radius.TopRight)))
				tl := float32(gtx.Px(unit.Dp(ll.Border.Radius.TopLeft)))
				br := float32(gtx.Px(unit.Dp(ll.Border.Radius.BottomRight)))
				bl := float32(gtx.Px(unit.Dp(ll.Border.Radius.BottomLeft)))
				clip.RRect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(gtx.Constraints.Min.X),
						Y: float32(gtx.Constraints.Min.Y),
					}},
					NW: tl, NE: tr, SE: br, SW: bl,
				}.Add(gtx.Ops)
				return fill(gtx, ll.Background)
			}),
			layout.Stacked(func(gtx C) D {
				ll.applyDimension(&gtx)
				return ll.Border.Layout(gtx, func(gtx C) D {
					// draw padding
					return ll.Padding.Layout(gtx, func(gtx C) D {
						// draw layout direction
						return ll.Direction.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: ll.Orientation}.Layout(gtx, children...)
						})
					})
				})
			}),
		)
	})
}

func (ll LinearLayout) applyDimension(gtx *C) {
	if ll.Width == MatchParent {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
	} else if ll.Width != WrapContent {
		gtx.Constraints.Min.X = ll.Width
		gtx.Constraints.Max.X = ll.Width
	}

	if ll.Height == MatchParent {
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	} else if ll.Height != WrapContent {
		gtx.Constraints.Min.Y = ll.Height
		gtx.Constraints.Max.Y = ll.Height
	}
}
