package layouts

import (
	"image/color"

	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/styles"
)

func Modal(gtx *layout.Context, w layout.Widget, shadow color.RGBA) {
	layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Stacked(func() {
			styles.FillWithColor(gtx, shadow)
		}),
		layout.Stacked(w),
	)
}
