package layouts

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
)

type Clicker struct {
	Widget layout.Widget
}

func (c Clicker) Layout(gtx *layout.Context, btn *widget.Button) {
	c.Widget()
	pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
	btn.Layout(gtx)
}

type Clickable interface {
	Layout(*layout.Context, widget.Button)
}
