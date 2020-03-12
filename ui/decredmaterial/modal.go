package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
)

// Modal lays out a widget Stacked (with Directrion) after a Stacked area filled with Background.
// The Stacked background is laid out with max Contraints.
type Modal struct {
	layout.Direction
	Overlay bool
}

// Layout the modal
func (m Modal) Layout(gtx *layout.Context, w layout.Widget) {
	if m.Overlay {
		fillMax(gtx, color.RGBA{A: 64})
	}
	layout.Stack{Alignment: m.Direction}.Layout(gtx,
		layout.Expanded(func() {
			fill(gtx, argb(0x0fffffff))
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
			new(widget.Button).Layout(gtx)
		}),
		layout.Stacked(w),
	)
}
