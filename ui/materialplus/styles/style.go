package styles

import (
	"gioui.org/layout"
)

type Style interface {
	Styled(*layout.Context, layout.Widget) layout.Widget
}

func WithStyles(gtx *layout.Context, widget layout.Widget, styles ...Style) layout.Widget {
	for _, style := range styles {
		widget = style.Styled(gtx, widget)
	}

	return widget
}

func WithStyle(gtx *layout.Context, style Style, widget layout.Widget) layout.Widget {
	return style.Styled(gtx, widget)
}
