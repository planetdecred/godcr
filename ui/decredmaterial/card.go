package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
)

type Card struct {
	layout.Inset
	Color color.RGBA
}

func (c Card) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return fill(gtx, color.RGBA{A: 64})
		}),
		layout.Stacked(func(gtx C) D {
			return c.Inset.Layout(gtx, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return fill(gtx, c.Color)
					}),
					layout.Stacked(w),
				)
			})
		}),
	)
	return dims
}
