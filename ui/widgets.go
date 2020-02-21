package ui

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type inputs struct {
	createDiag, deleteDiag, cancelDialog      widget.Button
	createWallet, restoreWallet, deleteWallet widget.Button
	toOverview, toWallets                     widget.Button
	tabs                                      []*widget.Button

	spendingPassword, matchSpending, renameWallet widget.Editor
}

type outputs struct {
	spendingPassword, matchSpending     material.Editor
	toOverview, toWallets               material.IconButton
	createDiag, restoreDiag, cancelDiag material.IconButton

	createWallet, restoreWallet, deleteWallet, deleteDiag material.Button

	notImplemented, noWallet, err material.Label
}

func (win *Window) initWidgets() {
	theme := win.theme
	win.outputs.spendingPassword = theme.Editor("Enter password")
	win.inputs.spendingPassword.SingleLine = true

	win.outputs.matchSpending = theme.Editor("Enter password again")
	win.inputs.matchSpending.SingleLine = true

	win.outputs.createDiag = theme.IconButton(mustIcon(material.NewIcon(icons.ContentAdd)))
	win.outputs.createWallet = theme.Button("create")

	win.outputs.deleteDiag = theme.DangerButton("Delete Wallet")
	win.outputs.deleteWallet = theme.DangerButton("delete")

	win.outputs.cancelDiag = theme.IconButton(mustIcon(material.NewIcon(icons.NavigationCancel)))
	win.outputs.cancelDiag.Background = theme.Danger

	win.outputs.notImplemented = theme.H3("Not Implemented")

	win.outputs.toWallets = theme.IconButton(mustIcon(material.NewIcon(icons.ActionAccountBalanceWallet)))
	win.outputs.toOverview = theme.IconButton(mustIcon(material.NewIcon(icons.ActionAccountBox)))

	win.outputs.noWallet = theme.H3("No wallet loaded")

	win.outputs.err = theme.Caption("")
	win.outputs.err.Color = theme.Danger
}
