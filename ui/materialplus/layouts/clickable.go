package layouts

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
)

// Clickable lays out a button over widget
type Clickable layout.Widget

// Layout lays out btn over the Clickable
func (c Clickable) Layout(gtx *layout.Context, btn *widget.Button) {
	c()
	pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
	btn.Layout(gtx)
}
