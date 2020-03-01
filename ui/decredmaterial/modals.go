package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
)

// ConfirmCancel lays out Body with layout
// Flex{Flexed(Cancel), Flexed(Body), Flexed(Confirm)}
type ConfirmCancel struct {
	layout.Flex
	layout.Direction
	Body struct {
		Size float32
		layout.Widget
	}
	Confirm struct {
		layout.Direction
		layout.Widget
	}
	Cancel struct {
		layout.Direction
		layout.Widget
	}
	Background color.RGBA
}

// Layout lays out the widget
func (cc ConfirmCancel) Layout(gtx *layout.Context, confirm, cancel *widget.Button) {
	modal := func() {
		s := (1 - (cc.Body.Size)) / 2
		cc.Flex.Layout(gtx,
			layout.Flexed(s, func() {
				cc.Cancel.Direction.Layout(gtx, cc.Cancel.Widget)
			}),
			layout.Flexed(cc.Body.Size, cc.Body.Widget),
			layout.Flexed(s, func() {
				cc.Confirm.Direction.Layout(gtx, cc.Confirm.Widget)
			}),
		)
	}
	Modal{
		Direction: cc.Direction,
	}.Layout(gtx, modal)
}
