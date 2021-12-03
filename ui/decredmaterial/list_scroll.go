package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ScrollbarStyle configures the presentation of a scrollbar.
type ScrollbarStyle struct {
	material.ScrollbarStyle
}

// ListStyle configures the presentation of a layout.List with a scrollbar.
type ListStyle struct {
	material.ListStyle
}

func (t *Theme) Scrollbar(state *widget.Scrollbar) ScrollbarStyle {
	return ScrollbarStyle{material.Scrollbar(t.Base, state)}
}

func (t *Theme) List(state *widget.List) ListStyle {
	list := ListStyle{material.List(t.Base, state)}
	list.Indicator.Color = t.Color.Gray3
	return list
}

// layout the scroll track and indicator.
func (s ScrollbarStyle) layout(gtx layout.Context, axis layout.Axis, viewportStart, viewportEnd float32) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return s.ScrollbarStyle.Layout(gtx, axis, viewportStart, viewportEnd)
}

// Layout the list and its scrollbar.
func (l ListStyle) Layout(gtx layout.Context, length int, w layout.ListElement) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return l.ListStyle.Layout(gtx, length, w)
}
