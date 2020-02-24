package layouts

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
)

// Modal lays out out a widget stacked above an area filled with Background
type Modal struct {
	Background color.RGBA
	layout.Direction
}

// Layout a widget.
func (m Modal) Layout(gtx *layout.Context, w layout.Widget) {
	layout.Stack{Alignment: m.Direction}.Layout(gtx,
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
			FillWithColor(gtx, m.Background)
			new(widget.Button).Layout(gtx)
		}),
		layout.Stacked(w),
	)
}
