package ui

import (
	"fmt"

	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type inputs struct {
	createDiag, deleteDiag, cancelDialog, restoreDiag         widget.Button
	createWallet, restoreWallet, deleteWallet                 widget.Button
	toOverview, toWallets, toTransactions, toSend, toSettings widget.Button
	toRestoreWallet                                           widget.Button
	//toReceive                                                 widget.Button
	sync widget.Button

	spendingPassword, matchSpending, renameWallet widget.Editor

	seedEditors         []widget.Editor
	seedsSuggestionsBtn []struct {
		text   string
		button widget.Button
	}
}

type combined struct {
	sel *decredmaterial.Select
}

type outputs struct {
	icons struct {
		add, check, cancel, sync decredmaterial.IconButton
	}
	spendingPassword, matchSpending                                            decredmaterial.Editor
	toOverview, toWallets, toRestoreWallet, toTransactions, toSend, toSettings decredmaterial.IconButton
	//toReceive                                                 decredmaterial.IconButton
	createDiag, cancelDiag decredmaterial.IconButton

	createWallet, restoreDiag, restoreWallet, deleteWallet, deleteDiag decredmaterial.Button
	sync                                                               decredmaterial.IconButton

	tabs                          []decredmaterial.TabItem
	notImplemented, noWallet, err decredmaterial.Label

	seedEditors         []decredmaterial.Editor
	seedsSuggestionsBtn []decredmaterial.Button
}

func (win *Window) initWidgets() {
	theme := win.theme

	win.combined.sel = theme.Select()

	win.outputs.icons.add = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentAdd)))
	win.outputs.icons.sync = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NotificationSync)))
	win.outputs.icons.cancel = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationCancel)))
	win.outputs.icons.cancel.Background = theme.Color.Danger
	win.outputs.icons.check = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationCheck)))
	win.outputs.icons.check.Background = theme.Color.Success

	win.outputs.spendingPassword = theme.Editor("Enter password")
	win.inputs.spendingPassword.SingleLine = true

	win.outputs.matchSpending = theme.Editor("Enter password again")
	win.inputs.matchSpending.SingleLine = true

	win.outputs.createDiag = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentCreate)))
	win.outputs.createWallet = theme.Button("create")

	win.outputs.toRestoreWallet = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionRestorePage)))
	win.outputs.restoreDiag = theme.Button("Restore wallet")
	win.outputs.restoreWallet = theme.Button("Restore")

	win.outputs.deleteDiag = theme.DangerButton("Delete Wallet")
	win.outputs.deleteWallet = theme.DangerButton("delete")

	win.outputs.cancelDiag = win.outputs.icons.cancel

	win.outputs.notImplemented = theme.H3("Not Implemented")

	win.outputs.toWallets = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionAccountBalanceWallet)))
	win.outputs.toOverview = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionHome)))
	win.outputs.toTransactions = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.SocialPoll)))
	win.outputs.toSettings = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionSettings)))
	win.outputs.toSend = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentSend)))

	win.outputs.noWallet = theme.H3("No wallet loaded")

	win.outputs.err = theme.Caption("")
	win.outputs.err.Color = theme.Color.Danger

	win.outputs.sync = win.outputs.icons.sync

	for i := 0; i <= 32; i++ {
		win.outputs.seedEditors = append(win.outputs.seedEditors, theme.Editor(fmt.Sprintf("Input word %d...", i+1)))
		win.inputs.seedEditors = append(win.inputs.seedEditors, widget.Editor{SingleLine: true, Submit: true})
	}
}

func mustIcon(ic *decredmaterial.Icon, err error) *decredmaterial.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
