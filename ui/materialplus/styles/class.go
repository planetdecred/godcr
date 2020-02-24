package styles

import "gioui.org/layout"

type Class struct {
	Alignment
	Background
}

func WithClass(gtx *layout.Context, class Class, widget layout.Widget) layout.Widget {
	return WithStyles(gtx, widget, class.Alignment, class.Background)
}
