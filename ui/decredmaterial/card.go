package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
)

type Card struct {
	layout.Inset
}

func (c Card) Layout(gtx *layout.Context, w layout.Widget) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fill(gtx, color.RGBA{A: 64})
		}),
		layout.Stacked(func() {
			c.Inset.Layout(gtx, func() {
				layout.Stack{}.Layout(gtx,
					layout.Expanded(func() {
						fill(gtx, argb(0x0fffffff))
					}),
					layout.Stacked(w),
				)
			})
		}),
	)
}
