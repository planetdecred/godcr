package ui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type inputs struct {
	createDiag, deleteDiag, cancelDialog                                 widget.Button
	createWallet, restoreWallet, deleteWallet, renameWallet              widget.Button
	addAccount, toggleWalletRename                                       widget.Button
	toOverview, toWallets, toTransactions, toSend, toSettings            widget.Button
	toReceive                                                            widget.Button
	toRestoreWallet                                                      widget.Button
	sync, syncHeader, hideMsgInfo, changePasswordDiag, signMessageDiag   widget.Button
	pasteAddr, pasteMsg, pasteSign, clearAddr, clearMsg, clearSign       widget.Button
	spendingPassword, matchSpending, oldSpendingPassword, rename, dialog widget.Editor
	addressInput, messageInput, signInput                                widget.Editor
	clearBtn, verifyBtn, verifyMessDiag, verifyInfo                      widget.Button
	restoreDiag, addAcctDiag, savePassword                               widget.Button

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
		add, check, cancel, sync, info, more, copy, verifyInfo         decredmaterial.IconButton
		pasteAddr, pasteMsg, pasteSign, clearAddr, clearMsg, clearSign decredmaterial.IconButton
	}
	spendingPassword, matchSpending, oldSpendingPassword, dialog, rename       decredmaterial.Editor
	messageInput, signInput                                                    decredmaterial.Editor
	addressInput                                                               decredmaterial.EditorCustom
	addressInput, signInput                                                               decredmaterial.EditorCustom
	toOverview, toWallets, toTransactions, toRestoreWallet, toSend, toSettings decredmaterial.IconButton
	toReceive                                                                  decredmaterial.IconButton
	createDiag, cancelDiag, addAcctDiag                                        decredmaterial.IconButton
	clearBtn, verifyBtn, verifyMessDiag                                        decredmaterial.Button

	createWallet, restoreDiag, restoreWallet, deleteWallet, deleteDiag, gotItDiag  decredmaterial.Button
	toggleWalletRename, renameWallet, syncHeader                                   decredmaterial.IconButton
	sync, moreDiag, hideMsgInfo, savePassword, changePasswordDiag, signMessageDiag decredmaterial.Button
	addAccount, newAddressDiag                                                     decredmaterial.Button
	info, more, copy, verifyInfo                                                   decredmaterial.IconButton
	pasteAddr, pasteMsg, pasteSign, clearAddr, clearMsg, clearSign                 decredmaterial.IconButton
	passwordBar                                                                    *decredmaterial.ProgressBar

	tabs                                                              []decredmaterial.TabItem
	notImplemented, noWallet, pageTitle, pageInfo, verifyMessage, err decredmaterial.Label

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
	win.outputs.passwordBar = theme.ProgressBar(0)
	win.outputs.oldSpendingPassword = theme.Editor("Enter old password")
	win.inputs.oldSpendingPassword.SingleLine = true

	win.outputs.matchSpending = theme.Editor("Enter password again")
	win.inputs.matchSpending.SingleLine = true

	// verify message widgets
	win.outputs.addressInput = theme.EditorCustom("Address", "Enter Address", &win.inputs.addressInput)
	win.outputs.signInput = theme.EditorCustom("Signature", "Enter Signature", &win.inputs.signInput)
	win.outputs.messageInput = theme.Editor("Message")
	win.outputs.verifyBtn = theme.Button("Verify")
	win.outputs.verifyBtn.TextSize = unit.Dp(13)
	win.outputs.clearBtn = theme.Button("Clear All")
	win.outputs.clearBtn.Background = color.RGBA{0, 0, 0, 0}
	win.outputs.clearBtn.Color = win.theme.Color.Primary
	win.outputs.clearBtn.TextSize = unit.Dp(13)

	win.outputs.verifyMessDiag = theme.Button("Verify Message")
	win.outputs.signMessageDiag = theme.Button("Sign Message")

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

	win.outputs.changePasswordDiag = theme.Button("Change Password")
	win.outputs.savePassword = theme.Button("Change")
	win.outputs.savePassword.TextSize = unit.Dp(11)

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
	win.outputs.verifyMessage = win.theme.H6("")
	win.outputs.sync = theme.Button("Reconnect")
	win.outputs.syncHeader = win.outputs.icons.sync
	win.outputs.moreDiag = theme.Button("more")
	win.outputs.hideMsgInfo = theme.Button("Got it")
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
		Size:       unit.Dp(40),
		Background: color.RGBA{},
		Color:      win.theme.Color.Primary,
		Padding:    unit.Dp(3),
	}

	win.outputs.rename = theme.Editor("")
	win.inputs.rename.SingleLine = true
	win.outputs.rename.TextSize = unit.Dp(20)

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

func mustDcrIcon(icon image.Image, err error) image.Image {
	if err != nil {
		panic(err)
	}
	return icon
}
