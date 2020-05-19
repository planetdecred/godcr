package ui

import (
	"image"
	"image/color"

	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type inputs struct {
	deleteDiag, cancelDialog                                             widget.Button
	deleteWallet, renameWallet                                           widget.Button
	addAccount, toggleWalletRename                                       widget.Button
	toOverview, toWallets, toTransactions, toSend, toSettings            widget.Button
	toReceive                                                            widget.Button
	changePasswordDiag, signMessageDiag                                  widget.Button
	spendingPassword, matchSpending, oldSpendingPassword, rename, dialog widget.Editor
	verifyMessDiag                                                       widget.Button
	restoreDiag, addAcctDiag, savePassword                               widget.Button

	receiveIcons struct {
		info, more, copy, gotItDiag, newAddressDiag widget.Button
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
	spendingPassword, matchSpending, oldSpendingPassword, dialog, rename decredmaterial.Editor
	toOverview, toWallets, toTransactions, toSend, toSettings            decredmaterial.IconButton
	toReceive                                                            decredmaterial.IconButton
	cancelDiag, addAcctDiag                                              decredmaterial.IconButton
	verifyMessDiag                                  decredmaterial.Button

	deleteWallet, deleteDiag, gotItDiag 								 decredmaterial.Button
	toggleWalletRename, renameWallet                                              decredmaterial.IconButton
	hideMsgInfo, savePassword, changePasswordDiag, signMessageDiag                decredmaterial.Button
	addAccount, newAddressDiag                                                    decredmaterial.Button
	info, more, copy                                                              decredmaterial.IconButton
	passwordBar                                                                   *decredmaterial.ProgressBar

	tabs                                     []decredmaterial.TabItem
	notImplemented, pageTitle, pageInfo, err decredmaterial.Label

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
	win.outputs.icons.copy = theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy)))
	win.outputs.icons.copy.Padding = unit.Dp(5)
	win.outputs.icons.copy.Size = unit.Dp(30)
	win.outputs.icons.copy.Background = theme.Color.Background
	win.outputs.icons.copy.Color = theme.Color.Text

	win.outputs.spendingPassword = theme.Editor("Enter password")
	win.outputs.spendingPassword.IsRequired = true
	win.inputs.spendingPassword.SingleLine = true

	win.outputs.passwordBar = theme.ProgressBar(0)

	win.outputs.oldSpendingPassword = theme.Editor("Enter old password")
	win.outputs.oldSpendingPassword.IsRequired = true
	win.inputs.oldSpendingPassword.SingleLine = true

	win.outputs.matchSpending = theme.Editor("Enter password again")
	win.outputs.matchSpending.IsRequired = true
	win.inputs.matchSpending.SingleLine = true

	// verify message
	win.outputs.verifyMessDiag = theme.Button("Verify Message")
	win.outputs.signMessageDiag = theme.Button("Sign Message")

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

	win.outputs.err = theme.Caption("")
	win.outputs.err.Color = theme.Color.Danger
	win.outputs.hideMsgInfo = theme.Button("Got it")
	win.outputs.more = win.outputs.icons.more
	win.outputs.info = win.outputs.icons.info
	win.outputs.copy = win.outputs.icons.copy

	//for i := 0; i <= 32; i++ {
	//	e := theme.Editor(fmt.Sprintf("Input word %d...", i+1))
	//	e.IsTitleLabel = false
	//	win.outputs.seedEditors = append(win.outputs.seedEditors, e)
	//	win.inputs.seedEditors.focusIndex = -1
	//	win.inputs.seedEditors.editors = append(win.inputs.seedEditors.editors, widget.Editor{SingleLine: true, Submit: true})
	//}
	// win.outputs.sync = theme.Button("Reconnect")

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
