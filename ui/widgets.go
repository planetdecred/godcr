package ui

import (
	"fmt"
	"image/color"

	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var msg = `After you or your counterparty has genrated a signature, you can use this 
form to verify the signature. 

Once you have entered the address, the message and the corresponding 
signature, you will see VALID if the signature appropraitely matches 
the address and message otherwise INVALID.`

type inputs struct {
	createDiag, deleteDiag, cancelDialog                           widget.Button
	restoreDiag, addAcctDiag                                       widget.Button
	createWallet, restoreWallet, deleteWallet, renameWallet        widget.Button
	addAccount, toggleWalletRename                                 widget.Button
	toOverview, toWallets, toTransactions, toSend, toSettings      widget.Button
	toRestoreWallet                                                widget.Button
	toReceive                                                      widget.Button
	toTransactionsFilters                                     widget.Button
	applyFiltersTransactions                                  widget.Button
	sync, syncHeader, vMsgInfoBtn                                  widget.Button
	pasteAddr, pasteMsg, pasteSign, clearAddr, clearMsg, clearSign widget.Button
	spendingPassword, matchSpending, rename, dialog                widget.Editor
	addressInput, messageInput, signInput                          widget.Editor
	clearBtn, verifyBtn, verifyMessDiag, verifyInfo                widget.Button

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

	transactionFilterDirection *widget.Enum
	transactionFilterSort      *widget.Enum
}

type combined struct {
	sel *decredmaterial.Select

	transaction struct {
		status, direction *decredmaterial.Icon
		amount, time      decredmaterial.Label
	}
}

type outputs struct {
	ic struct {
		create, clear, done *decredmaterial.Icon
	}
	icons struct {
		add, check, cancel, sync, info, more, copy, verifyInfo         decredmaterial.IconButton
		pasteAddr, pasteMsg, pasteSign, clearAddr, clearMsg, clearSign decredmaterial.IconButton
	}
	spendingPassword, matchSpending, dialog, rename                            decredmaterial.Editor
	addressInput, messageInput, signInput                                      decredmaterial.Editor
	toOverview, toWallets, toTransactions, toRestoreWallet, toSend, toSettings decredmaterial.IconButton
	toReceive                                                                  decredmaterial.IconButton
	createDiag, cancelDiag, addAcctDiag                                        decredmaterial.IconButton
	clearBtn, verifyBtn, verifyMessDiag                                        decredmaterial.Button

	createWallet, restoreDiag, restoreWallet, deleteWallet, deleteDiag, gotItDiag decredmaterial.Button
	toggleWalletRename, renameWallet, syncHeader                                  decredmaterial.IconButton
	applyFiltersTransactions                                                      decredmaterial.Button
	sync, moreDiag, vMsgInfoBtn                                                   decredmaterial.Button

	addAccount, newAddressDiag                                     decredmaterial.Button
	info, more, copy, verifyInfo                                   decredmaterial.IconButton
	pasteAddr, pasteMsg, pasteSign, clearAddr, clearMsg, clearSign decredmaterial.IconButton

	tabs                                                                        []decredmaterial.TabItem
	notImplemented, noWallet, pageTitle, pageInfo, vMsgInfo, verifyMessage, err decredmaterial.Label

	seedEditors      []decredmaterial.Editor
	seedsSuggestions []decredmaterial.Button

	//receive page labels
	selectedAccountNameLabel, selectedAccountBalanceLabel           decredmaterial.Label
	receiveAddressLabel, accountModalTitleLabel, addressCopiedLabel decredmaterial.Label
	selectedWalletBalLabel, selectedWalletNameLabel                 decredmaterial.Label

	toTransactionsFilters struct {
		sortNewest, sortOldest decredmaterial.IconButton
	}
	transactionFilterDirection []decredmaterial.RadioButton
	transactionFilterSort      []decredmaterial.RadioButton
}

func (win *Window) initWidgets() {
	theme := win.theme

	win.combined.sel = theme.Select()

	win.outputs.ic.clear = mustIcon(decredmaterial.NewIcon(icons.ContentClear))
	win.outputs.ic.done = mustIcon(decredmaterial.NewIcon(icons.ActionDone))
	win.outputs.ic.create = mustIcon(decredmaterial.NewIcon(icons.ContentCreate))

	win.outputs.icons.add = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentAdd)))
	win.outputs.icons.sync = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NotificationSync)))
	win.outputs.icons.cancel = decredmaterial.IconButton{
		Icon:       win.outputs.ic.clear,
		Size:       unit.Dp(40),
		Background: color.RGBA{},
		Color:      win.theme.Color.Hint,
		Padding:    unit.Dp(0),
	}

	win.outputs.icons.check = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationCheck)))
	win.outputs.icons.check.Background = theme.Color.Success
	win.outputs.icons.more = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationMoreVert)))
	win.outputs.icons.more.Padding = unit.Dp(5)
	win.outputs.icons.more.Size = unit.Dp(35)
	win.outputs.icons.info = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionInfo)))
	win.outputs.icons.info.Padding = unit.Dp(5)
	win.outputs.icons.info.Size = unit.Dp(35)
	win.outputs.icons.verifyInfo = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionInfo)))
	win.outputs.icons.verifyInfo.Padding = unit.Dp(5)
	win.outputs.icons.verifyInfo.Size = unit.Dp(35)
	win.outputs.icons.copy = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy)))
	win.outputs.icons.copy.Padding = unit.Dp(5)
	win.outputs.icons.copy.Size = unit.Dp(30)
	win.outputs.icons.copy.Background = theme.Color.Background
	win.outputs.icons.copy.Color = theme.Color.Text
	win.outputs.icons.pasteAddr = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentPaste)))
	win.outputs.icons.pasteAddr.Padding = unit.Dp(5)
	win.outputs.icons.pasteAddr.Size = unit.Dp(30)
	win.outputs.icons.pasteAddr.Background = theme.Color.Background
	win.outputs.icons.pasteAddr.Color = theme.Color.Text
	win.outputs.icons.pasteMsg = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentPaste)))
	win.outputs.icons.pasteMsg.Padding = unit.Dp(5)
	win.outputs.icons.pasteMsg.Size = unit.Dp(30)
	win.outputs.icons.pasteMsg.Background = theme.Color.Background
	win.outputs.icons.pasteMsg.Color = theme.Color.Text
	win.outputs.icons.pasteSign = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentPaste)))
	win.outputs.icons.pasteSign.Padding = unit.Dp(5)
	win.outputs.icons.pasteSign.Size = unit.Dp(30)
	win.outputs.icons.pasteSign.Background = theme.Color.Background
	win.outputs.icons.pasteSign.Color = theme.Color.Text
	win.outputs.icons.clearSign = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentClear)))
	win.outputs.icons.clearSign.Padding = unit.Dp(5)
	win.outputs.icons.clearSign.Size = unit.Dp(30)
	win.outputs.icons.clearSign.Background = theme.Color.Background
	win.outputs.icons.clearSign.Color = theme.Color.Text
	win.outputs.icons.clearAddr = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentClear)))
	win.outputs.icons.clearAddr.Padding = unit.Dp(5)
	win.outputs.icons.clearAddr.Size = unit.Dp(30)
	win.outputs.icons.clearAddr.Background = theme.Color.Background
	win.outputs.icons.clearAddr.Color = theme.Color.Text
	win.outputs.icons.clearMsg = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentClear)))
	win.outputs.icons.clearMsg.Padding = unit.Dp(5)
	win.outputs.icons.clearMsg.Size = unit.Dp(30)
	win.outputs.icons.clearMsg.Background = theme.Color.Background
	win.outputs.icons.clearMsg.Color = theme.Color.Text

	win.outputs.spendingPassword = theme.Editor("Enter password")
	win.inputs.spendingPassword.SingleLine = true

	win.outputs.matchSpending = theme.Editor("Enter password again")
	win.inputs.matchSpending.SingleLine = true

	// verify message widgets
	win.outputs.addressInput = theme.Editor("Address")
	win.outputs.signInput = theme.Editor("Signature")
	win.outputs.messageInput = theme.Editor("Message")
	win.outputs.verifyBtn = theme.Button("Verify")
	win.outputs.verifyBtn.TextSize = unit.Dp(13)
	win.outputs.clearBtn = theme.Button("Clear All")
	win.outputs.clearBtn.Background = win.theme.Color.Transparent
	win.outputs.clearBtn.Color = win.theme.Color.Primary
	win.outputs.clearBtn.TextSize = unit.Dp(13)

	win.outputs.verifyMessDiag = theme.Button("Verify Message")

	//
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
	win.outputs.vMsgInfo = theme.Body1(msg)

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
	win.outputs.toTransactionsFilters.sortNewest = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentFilterList)))
	win.outputs.toTransactionsFilters.sortOldest = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentSort)))

	win.outputs.noWallet = theme.H3("No wallet loaded")

	win.outputs.err = theme.Caption("")
	win.outputs.err.Color = theme.Color.Danger
	win.outputs.verifyMessage = win.theme.H6("")
	win.outputs.sync = theme.Button("Reconnect")
	win.outputs.syncHeader = win.outputs.icons.sync
	win.outputs.moreDiag = theme.Button("more")
	win.outputs.vMsgInfoBtn = theme.Button("Got it")

	win.outputs.more = win.outputs.icons.more
	win.outputs.info = win.outputs.icons.info
	win.outputs.copy = win.outputs.icons.copy
	win.outputs.pasteAddr = win.outputs.icons.pasteAddr
	win.outputs.pasteMsg = win.outputs.icons.pasteMsg
	win.outputs.pasteSign = win.outputs.icons.pasteSign
	win.outputs.clearAddr = win.outputs.icons.clearAddr
	win.outputs.clearMsg = win.outputs.icons.clearMsg
	win.outputs.clearSign = win.outputs.icons.clearSign
	win.outputs.verifyInfo = win.outputs.icons.verifyInfo

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

	win.inputs.transactionFilterDirection = new(widget.Enum)
	win.inputs.transactionFilterDirection.SetValue("0")
	win.inputs.transactionFilterSort = new(widget.Enum)
	win.inputs.transactionFilterSort.SetValue("0")

	txFiltersDirection := []string{"All", "Sent", "Received", "Transfer"}
	txSort := []string{"Newest", "Oldest"}

	for i := 0; i < len(txFiltersDirection); i++ {
		win.outputs.transactionFilterDirection = append(
			win.outputs.transactionFilterDirection,
			theme.RadioButton(fmt.Sprint(i), txFiltersDirection[i]))
	}

	for i := 0; i < len(txSort); i++ {
		win.outputs.transactionFilterSort = append(
			win.outputs.transactionFilterSort,
			theme.RadioButton(fmt.Sprint(i), txSort[i]))
	}

	win.outputs.applyFiltersTransactions = theme.Button("Ok")
}

func mustIcon(ic *decredmaterial.Icon, err error) *decredmaterial.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
