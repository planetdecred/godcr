package ui

import (
	"gioui.org/layout"
)

// Landing lays out the windows landing page
func (win *Window) Landing() layout.Widget {
	log.Debug("On Landing")
	return func() {
		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEnd,
		}.Layout(win.gtx,
			layout.Flexed(0.3, func() {
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
					)
				})
			}),
		)
	}
}
