package ui

import (
	"gioui.org/layout"
)

// Landing page lays out the create wallet and restore wallet buttons
func (win *Window) Landing() {
	toMax(win.gtx)
	layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(win.gtx,
		layout.Flexed(headerHeight, func() {
			win.Header()
		}),
		layout.Rigid(func() {
			layout.Center.Layout(win.gtx, func() {
				layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceBetween,
				}.Layout(win.gtx,
					layout.Flexed(.3, func() {
						win.theme.Button("Create Wallet").Layout(win.gtx, &win.inputs.createWallet)
					}),
					layout.Flexed(.3, func() {
						win.theme.Button("Restore Wallet").Layout(win.gtx, &win.inputs.restoreWallet)
					}),
					layout.Flexed(.1, func() {
						win.widgets.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
					}),
				)
			})
		}),
	)
}
