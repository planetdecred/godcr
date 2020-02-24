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

func (align Alignment) Layout(gtx *layout.Context, widget func()) func() {
	return func() {
		align.Direction.Layout(gtx, widget)
	}
}
