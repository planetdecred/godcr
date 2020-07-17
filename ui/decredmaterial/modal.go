package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// Modal lays out a widget Stacked (with Directrion) after a Stacked area filled with Background.
// The Stacked background is laid out with max Contraints.
type Modal struct {
	layout.Direction
}

// Layout the modal
func (m Modal) Layout(gtx layout.Context, background, dialog layout.Widget) layout.Dimensions {
	dims := layout.Stack{Alignment: m.Direction}.Layout(gtx,
		layout.Stacked(background),
		layout.Expanded(func(gtx C) D {
			return fill(gtx, argb(0x7F444444))
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return new(widget.Clickable).Layout(gtx)
		}),
		layout.Stacked(dialog),
	)
	return dims
}
