package uiwallet

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
)

// Loading lays out the loading widget with a faded background
func (w *Wallet) Loading(gtx layout.Context) {
	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := w.theme.H1("Loading")
		lbl.Alignment = text.Middle
		return lbl.Layout(gtx)
	})
	new(widget.Clickable).Layout(gtx)
}
