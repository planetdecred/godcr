package styles

import (
	"image/color"

	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
)

var (
	White = color.RGBA{}
	Black = RGB(0xfffff)
	Blue  = RGB(0x00ff00)
)

type Background color.RGBA

func (c Background) Styled(gtx *layout.Context, widget layout.Widget) layout.Widget {
	return func() {
		wmin := gtx.Constraints.Width.Min
		hmin := gtx.Constraints.Height.Min
		layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func() {
				layouts.FillWithColor(gtx, color.RGBA(c))
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = wmin
				gtx.Constraints.Height.Min = hmin
				widget()
			}),
		)
	}
}
func RGB(hex int32) color.RGBA {
	return RGBA((hex << 4) | 0xff)
}

func RGBA(hex int32) color.RGBA {
	return color.RGBA{R: uint8(hex >> 24), B: uint8(hex >> 16), G: uint8(hex >> 8), A: uint8(hex)}
}
