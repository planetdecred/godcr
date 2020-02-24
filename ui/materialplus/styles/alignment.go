package styles

import (
	"gioui.org/layout"
)

type Alignment struct {
	layout.Direction
}

var (
	Centered = Alignment{layout.Center}
)

func (align Alignment) Styled(gtx *layout.Context, w layout.Widget) layout.Widget {
	return func() {
		align.Direction.Layout(gtx, w)
	}
}
