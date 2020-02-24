package layouts

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
)

// Clicker lays out a button over widget
type Clicker layout.Widget

// Layout lays out btn over the Clicker widget
func (c Clicker) Layout(gtx *layout.Context, btn *widget.Button) {
	c()
	pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
	btn.Layout(gtx)
}
