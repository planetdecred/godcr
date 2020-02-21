package materialplus

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
)

// Modal lays out a widget Stacked (with Directrion) after a Stacked area filled with Background.
// The Stacked background is laid out with max Contraints.
type Modal struct {
	Background color.RGBA
	layout.Direction
}

// Layout the modal
func (m Modal) Layout(gtx *layout.Context, w layout.Widget) {
	layout.Stack{Alignment: m.Direction}.Layout(gtx,
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
			fillWithColor(gtx, m.Background)
			new(widget.Button).Layout(gtx)
		}),
		layout.Stacked(w),
	)
}
