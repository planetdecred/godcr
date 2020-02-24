package layouts

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
)

func Modal(gtx *layout.Context, w layout.Widget, shadow color.RGBA) {
	layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
			FillWithColor(gtx, shadow)
			new(widget.Button).Layout(gtx)
		}),
		layout.Stacked(w),
	)
}
