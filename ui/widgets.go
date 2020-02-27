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
	sync                                      widget.Button
	tabs                                      []*widget.Button

	spendingPassword, matchSpending, renameWallet widget.Editor
}

type outputs struct {
	icons struct {
		add, check, cancel, sync material.IconButton
	}
	spendingPassword, matchSpending     material.Editor
	toOverview, toWallets               material.IconButton
	createDiag, restoreDiag, cancelDiag material.IconButton

	createWallet, restoreWallet, deleteWallet, deleteDiag material.Button
	sync                                                  material.IconButton

	notImplemented, noWallet, err material.Label
}

func (win *Window) initWidgets() {
	theme := win.theme

	win.outputs.icons.add = theme.IconButton(mustIcon(material.NewIcon(icons.ContentAdd)))
	win.outputs.icons.sync = theme.IconButton(mustIcon(material.NewIcon(icons.NotificationSync)))
	win.outputs.icons.cancel = theme.IconButton(mustIcon(material.NewIcon(icons.NavigationCancel)))
	win.outputs.icons.cancel.Background = theme.Danger
	win.outputs.icons.check = theme.IconButton(mustIcon(material.NewIcon(icons.NavigationCheck)))
	win.outputs.icons.check.Background = theme.Success

	win.outputs.spendingPassword = theme.Editor("Enter password")
	win.inputs.spendingPassword.SingleLine = true

	win.outputs.matchSpending = theme.Editor("Enter password again")
	win.inputs.matchSpending.SingleLine = true

	win.outputs.createDiag = theme.IconButton(mustIcon(material.NewIcon(icons.ContentCreate)))
	win.outputs.createWallet = theme.Button("create")

	win.outputs.deleteDiag = theme.DangerButton("Delete Wallet")
	win.outputs.deleteWallet = theme.DangerButton("delete")

	win.outputs.cancelDiag = win.outputs.icons.cancel

	win.outputs.notImplemented = theme.H3("Not Implemented")

	win.outputs.toWallets = theme.IconButton(mustIcon(material.NewIcon(icons.ActionAccountBalanceWallet)))
	win.outputs.toOverview = theme.IconButton(mustIcon(material.NewIcon(icons.ActionHome)))

	win.outputs.noWallet = theme.H3("No wallet loaded")

	win.outputs.err = theme.Caption("")
	win.outputs.err.Color = theme.Danger

	win.outputs.sync = win.outputs.icons.sync
}
