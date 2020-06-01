package ui

import (
	"gioui.org/layout"
)

func (win *Window) AddAccountDiag() {
	win.theme.Surface(win.gtx, func() {
		toMax(win.gtx)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.E.Layout(win.gtx, func() {
					// win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
				})
			}),
			layout.Rigid(func() {
				d := win.theme.H3("Add account ")
				d.Layout(win.gtx)
			}),
			layout.Rigid(func() {
				win.outputs.dialog.Layout(win.gtx, &win.inputs.dialog)
			}),
			layout.Rigid(func() {
				win.outputs.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
			}),
			layout.Rigid(func() {
				win.Err()
			}),
			layout.Rigid(func() {
				win.outputs.addAccount.Layout(win.gtx, &win.inputs.addAccount)
			}),
		)
	})
}
