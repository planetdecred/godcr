package materialplus

import (
	"image/color"

	"gioui.org/layout"
)

type Card struct {
	layout.Inset
}

func (c Card) Layout(gtx *layout.Context, w layout.Widget) {

	layout.Stack{Alignment: layout.NW}.Layout(gtx,
		layout.Expanded(func() {
			gtx.Constraints.Height.Min -= 10
			gtx.Constraints.Height.Max -= 10
			fillWithColor(gtx, color.RGBA{A: 128})
		}),
		layout.Stacked(func() {
			fillWithColor(gtx, ARGB(0x0fffffff))
			c.Inset.Layout(gtx, w)
		}),
	)
}
