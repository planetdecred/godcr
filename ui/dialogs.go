package ui

import (
	"gioui.org/layout"
)

func (win *Window) CreateDiag() {
	win.theme.Background(win.gtx, func() {
		toMax(win.gtx)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.E.Layout(win.gtx, func() {
					win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
				})
			}),
			layout.Rigid(func() {
				win.theme.H3("Create new wallet").Layout(win.gtx)
			}),
			layout.Rigid(func() {
				win.outputs.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
			}),
			layout.Rigid(func() {
				win.outputs.matchSpending.Layout(win.gtx, &win.inputs.matchSpending)
			}),
			layout.Rigid(func() {
				win.outputs.createWallet.Layout(win.gtx, &win.inputs.createWallet)
			}),
		)
	})
}

func (win *Window) DeleteDiag() {
	win.theme.Background(win.gtx, func() {
		toMax(win.gtx)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.E.Layout(win.gtx, func() {
					win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
				})
			}),
			layout.Rigid(func() {
				win.theme.H3("Delete wallet").Layout(win.gtx)
			}),
			layout.Rigid(func() {
				win.outputs.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
			}),
			layout.Rigid(func() {
				win.outputs.deleteWallet.Layout(win.gtx, &win.inputs.deleteWallet)
			}),
		)
	})
}
