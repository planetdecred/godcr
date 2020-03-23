package ui

import (
	"fmt"

	"image/color"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type inputs struct {
	createDiag, deleteDiag, cancelDialog                      widget.Button
	restoreDiag, addAcctDiag                                  widget.Button
	createWallet, restoreWallet, deleteWallet, renameWallet   widget.Button
	addAccount, toggleWalletRename                            widget.Button
	toOverview, toWallets, toTransactions, toSend, toSettings widget.Button
	toRestoreWallet                                           widget.Button
	toReceive                                                 widget.Button
	sync, syncHeader widget.Button
	info, more, dropdown, copy, gotIt, newAddress                          widget.Button
	spendingPassword, matchSpending, rename, dialog widget.Editor

	receiveIcons struct {
		info, more, copy, gotItDiag, newAddressDiag widget.Button
	}
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
}

type outputs struct {
	ic struct {
		create, clear, done *decredmaterial.Icon
	}
	icons struct {
		add, check, cancel, sync, info, more, copy decredmaterial.IconButton
	}
	spendingPassword, matchSpending, dialog, rename                            decredmaterial.Editor
	toOverview, toWallets, toTransactions, toRestoreWallet, toSend, toSettings decredmaterial.IconButton
	toReceive                                                                  decredmaterial.IconButton
	createDiag, cancelDiag, addAcctDiag                                        decredmaterial.IconButton

	createWallet, restoreDiag, restoreWallet, deleteWallet, deleteDiag, gotItDiag decredmaterial.Button
	toggleWalletRename, renameWallet, syncHeader                       decredmaterial.IconButton
	sync, moreDiag                                                         decredmaterial.Button

	addAccount, newAddressDiag                                                    decredmaterial.Button
	info, more, copy            decredmaterial.IconButton

	tabs                                               []decredmaterial.TabItem
	notImplemented, noWallet, pageTitle, pageInfo, err decredmaterial.Label

	seedEditors      []decredmaterial.Editor
	seedsSuggestions []decredmaterial.Button

	//receive page labels
	selectedAccountNameLabel, selectedAccountBalanceLabel           decredmaterial.Label
	receiveAddressLabel, accountModalTitleLabel, addressCopiedLabel decredmaterial.Label
	selectedWalletBalLabel, selectedWalletNameLabel                 decredmaterial.Label
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
	win.outputs.icons.more = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationMoreVert)))
	win.outputs.icons.more.Padding = unit.Dp(5)
	win.outputs.icons.more.Size = unit.Dp(35)
	win.outputs.icons.info = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionInfo)))
	win.outputs.icons.info.Padding = unit.Dp(5)
	win.outputs.icons.info.Size = unit.Dp(35)
	win.outputs.icons.copy = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy)))
	win.outputs.icons.copy.Padding = unit.Dp(5)
	win.outputs.icons.copy.Size = unit.Dp(30)
	win.outputs.icons.copy.Background = theme.Color.Background
	win.outputs.icons.copy.Color = theme.Color.Text
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

	//receive widgets
	win.outputs.gotItDiag = theme.Button("Got It")
	win.outputs.newAddressDiag = theme.Button("Generate new address")

	win.outputs.pageTitle = theme.H4("Receiving DCR")
	win.outputs.pageInfo = theme.Body1("Each time you request a payment, a \nnew address is created to protect \nyour privacy.")

	win.outputs.selectedAccountNameLabel = win.theme.H6("")
	win.outputs.selectedWalletNameLabel = win.theme.Body2("")
	win.outputs.selectedWalletBalLabel = win.theme.Body2("")
	win.outputs.selectedAccountBalanceLabel = win.theme.H6("")
	win.outputs.receiveAddressLabel = win.theme.H6("")
	win.outputs.receiveAddressLabel.Color = theme.Color.Primary
	win.outputs.addressCopiedLabel = win.theme.Caption("")

	win.outputs.toWallets = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionAccountBalanceWallet)))
	win.outputs.toOverview = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionHome)))
	win.outputs.toTransactions = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.SocialPoll)))
	win.outputs.toSettings = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionSettings)))
	win.outputs.toSend = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentSend)))
	win.outputs.toReceive = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentAddBox)))

	win.outputs.noWallet = theme.H3("No wallet loaded")

	win.outputs.err = theme.Caption("")
	win.outputs.err.Color = theme.Color.Danger
	win.outputs.sync = theme.Button("Reconnect")
	win.outputs.syncHeader = win.outputs.icons.sync
	win.outputs.moreDiag = theme.Button("more")

	win.outputs.more = win.outputs.icons.more
	win.outputs.info = win.outputs.icons.info
	win.outputs.copy = win.outputs.icons.copy

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
}

func mustIcon(ic *decredmaterial.Icon, err error) *decredmaterial.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
