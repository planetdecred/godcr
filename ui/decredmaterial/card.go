package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
)

type Card struct {
	layout.Inset
	Color color.RGBA
}

func (c Card) Layout(gtx *layout.Context, w layout.Widget) {
	layout.Stack{}.Layout(gtx,
		layout.Stacked(func() {
			c.Inset.Layout(gtx, func() {
				layout.Stack{}.Layout(gtx,
					layout.Expanded(func() {
						fill(gtx, c.Color)
					}),
					layout.Stacked(w),
				)
			})
		}),
	)
}
