package styles

import (
	"gioui.org/layout"
)

type Style interface {
	Layout(gtx *layout.Context, widget func()) func()
}

func WithStyles(gtx *layout.Context, widget func(), styles ...Style) func() {
	for _, style := range styles {
		widget = style.Layout(gtx, widget)
	}

	return widget
}

func WithStyle(gtx *layout.Context, style Style, widget func()) func() {
	return style.Layout(gtx, widget)
}
