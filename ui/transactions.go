package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func (win *Window) TransactionsPage() {
	bd := func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Flexed(.2, func() {
				layout.Flex{Spacing: layout.SpaceBetween}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.H4("Transactions").Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.outputs.toSend.Layout(win.gtx, &win.inputs.toSend)
					}),
					layout.Rigid(func() {
						layout.Inset{Right: unit.Dp(20)}.Layout(win.gtx, func() {
							win.outputs.toReceive.Layout(win.gtx, &win.inputs.toReceive)
						})
					}),
				)
			}),
		)
	}
	win.TabbedPage(bd)
}
