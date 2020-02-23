package styles

import (
	"gioui.org/layout"
)

type Alignment layout.Alignment

const (
	Centered = Alignment(layout.Center)
)

func (align Alignment) Layout(gtx *layout.Context, widget func()) func() {
	return func() {
		layout.Align(align).Layout(gtx, widget)
	}
}
