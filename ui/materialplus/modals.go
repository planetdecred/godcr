package materialplus

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/ui/layouts"
	"github.com/raedahgroup/godcr-gio/ui/styles"
)

type ConfirmCancel struct {
	Body    layout.Widget
	Confirm material.Button
	Cancel  material.IconButton
}

func (dialog ConfirmCancel) Layout(gtx *layout.Context, confirm, cancel *widget.Button) {
	modal := func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(0.20, func() { dialog.Cancel.Layout(gtx, cancel) }),
			layout.Flexed(0.60, dialog.Body),
			layout.Flexed(0.20, func() { dialog.Confirm.Layout(gtx, confirm) }),
		)
	}
	layouts.Modal(gtx, modal, styles.RGBA(0xffffff6))
}
