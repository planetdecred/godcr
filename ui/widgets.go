package ui

import (
	"fmt"

	"image/color"

	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"github.com/raedahgroup/godcr-gio/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type inputs struct {
	createDiag, deleteDiag, cancelDialog                      widget.Button
	restoreDiag, addAcctDiag                                  widget.Button
	createWallet, restoreWallet, deleteWallet, renameWallet   widget.Button
	addAccount, toggleWalletRename                            widget.Button
	toOverview, toWallets, toTransactions, toSend, toSettings widget.Button
	toRestoreWallet                                           widget.Button
	//toReceive                                                 widget.Button
	sync, syncHeader widget.Button

	spendingPassword, matchSpending, rename, dialog widget.Editor

	seedEditors struct {
		focusIndex int
		editors    []widget.Editor
	}
	seedsSuggestions []struct {
		text   string
		button widget.Button
	}
}

type combined struct {
	sel *decredmaterial.Select

	transactionStatus *decredmaterial.Select
	transactionSort   *decredmaterial.Select
	transactions      []struct {
		iconStatus    *decredmaterial.Icon
		iconDirection *decredmaterial.Icon
		info          *wallet.TransactionInfo
	}
}

type outputs struct {
	ic struct {
		create, clear, done *decredmaterial.Icon
	}
	icons struct {
		add, check, cancel, sync decredmaterial.IconButton
	}
	spendingPassword, matchSpending, dialog, rename                            decredmaterial.Editor
	toOverview, toWallets, toTransactions, toRestoreWallet, toSend, toSettings decredmaterial.IconButton
	//toReceive                                                 decredmaterial.IconButton
	createDiag, cancelDiag, addAcctDiag decredmaterial.IconButton

	createWallet, restoreDiag, restoreWallet, deleteWallet, deleteDiag decredmaterial.Button
	addAccount                                                         decredmaterial.Button
	toggleWalletRename, renameWallet, syncHeader                       decredmaterial.IconButton
	sync, more                                                         decredmaterial.Button

	tabs                          []decredmaterial.TabItem
	notImplemented, noWallet, err decredmaterial.Label

	seedEditors      []decredmaterial.Editor
	seedsSuggestions []decredmaterial.Button
}

func (win *Window) initWidgets() {
	theme := win.theme

	win.combined.sel = theme.Select()

	win.outputs.ic.clear = mustIcon(decredmaterial.NewIcon(icons.ContentClear))
	win.outputs.ic.done = mustIcon(decredmaterial.NewIcon(icons.ActionDone))
	win.outputs.ic.create = mustIcon(decredmaterial.NewIcon(icons.ContentCreate))

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
	win.outputs.sync = theme.Button("Reconnect")
	win.outputs.syncHeader = win.outputs.icons.sync
	win.outputs.more = theme.Button("more")

	for i := 0; i <= 32; i++ {
		win.outputs.seedEditors = append(win.outputs.seedEditors, theme.Editor(fmt.Sprintf("Input word %d...", i+1)))
		win.inputs.seedEditors.focusIndex = -1
		win.inputs.seedEditors.editors = append(win.inputs.seedEditors.editors, widget.Editor{SingleLine: true, Submit: true})
	}
	win.outputs.sync = theme.Button("Reconnect")

	win.outputs.addAcctDiag = win.outputs.icons.add
	win.outputs.addAccount = theme.Button("add")
	win.outputs.addAcctDiag.Size = unit.Dp(24)
	win.outputs.addAcctDiag.Padding = unit.Dp(0)
	win.outputs.dialog = theme.Editor("")

	win.outputs.toggleWalletRename = decredmaterial.IconButton{
		Icon:       win.outputs.ic.create,
		Size:       unit.Dp(48),
		Background: color.RGBA{},
		Color:      win.theme.Color.Primary,
		Padding:    unit.Dp(0),
	}

	win.outputs.rename = theme.Editor("")
	win.outputs.renameWallet = decredmaterial.IconButton{
		Icon:       win.outputs.ic.done,
		Size:       unit.Dp(48),
		Background: color.RGBA{},
		Color:      win.theme.Color.Success,
		Padding:    unit.Dp(0),
	}

	win.combined.transactionStatus = theme.Select()
	win.combined.transactionStatus.Options = []string{"All", "Sent", "Received", "Transfer"}
	win.combined.transactionSort = theme.Select()
	win.combined.transactionSort.Options = []string{"Newest", "Oldest"}
}

func mustIcon(ic *decredmaterial.Icon, err error) *decredmaterial.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
