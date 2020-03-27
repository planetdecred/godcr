package ui

import (
	"sort"
	"strings"
	"time"

	"gioui.org/io/key"
	// "gioui.org/unit"
	"gioui.org/widget"
	"github.com/atotto/clipboard"
	"github.com/raedahgroup/dcrlibwallet"
)

var (
	old     int
	newAddr bool
)

// HandleInputs handles all ui inputs
func (win *Window) HandleInputs() {
	if win.tabs.Changed() {
		if win.tabs.Selected != win.selected {
			win.combined.sel.Selected = 0
			win.selectedAccount = 0
			win.selected = win.tabs.Selected
			win.sortTransactions()
		}
	}
	if win.combined.sel.Changed() {
		win.selectedAccount = win.combined.sel.Selected
	}

	for _, evt := range win.inputs.spendingPassword.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.outputs.spendingPassword.HintColor = win.theme.Color.Hint
			return
		}
		log.Debug("Pass evt", evt)
	}

	for _, evt := range win.inputs.matchSpending.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.outputs.matchSpending.Color = win.theme.Color.Text
			win.outputs.matchSpending.HintColor = win.theme.Color.Hint
			return
		}
		log.Debug("Match evt", evt)
	}

	for _, evt := range win.inputs.addressInput.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.err = ""
			win.outputs.verifyMessage.Text = ""
			win.outputs.addressInput.HintColor = win.theme.Color.Hint
			return
		}
		log.Debug("address evt", evt)
	}

	for _, evt := range win.inputs.signInput.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.err = ""
			win.outputs.verifyMessage.Text = ""
			win.outputs.signInput.HintColor = win.theme.Color.Hint
			return
		}
		log.Debug("sign evt", evt)
	}

	for _, evt := range win.inputs.messageInput.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.err = ""
			win.outputs.verifyMessage.Text = ""
			win.outputs.messageInput.HintColor = win.theme.Color.Hint
			return
		}
		log.Debug("Match evt", evt)
	}

	// CREATE WALLET
	if win.inputs.createDiag.Clicked(win.gtx) {
		win.dialog = win.CreateDiag
		win.states.dialog = true
	}

	if win.inputs.createWallet.Clicked(win.gtx) {
		pass := win.validatePasswords()
		if pass == "" {
			return
		}
		win.wallet.CreateWallet(pass)
		win.resetPasswords()
		log.Debug("Create Wallet clicked")
		win.states.loading = true
		return
	}

	// RESTORE WALLET

	if win.inputs.toRestoreWallet.Clicked(win.gtx) {
		win.current = PageRestore
		return
	}

	win.editorSeedsEventsHandler()
	win.onSuggestionSeedsClicked()

	if win.inputs.restoreDiag.Clicked(win.gtx) && win.validateSeeds() != "" {
		win.dialog = win.RestoreDiag
		win.states.dialog = true
	}

	if win.inputs.restoreWallet.Clicked(win.gtx) {
		pass := win.validatePasswords()
		if pass == "" {
			return
		}

		win.wallet.RestoreWallet(win.validateSeeds(), pass)
		win.states.loading = true
		log.Debug("Restore Wallet clicked")
		return
	}

	// EDIT WALLET

	if win.inputs.toggleWalletRename.Clicked(win.gtx) {
		win.dialog = win.editWalletDiag
		if win.states.renamingWallet {
			win.outputs.toggleWalletRename.Icon = win.outputs.ic.create
			win.outputs.toggleWalletRename.Color = win.theme.Color.Primary
		} else {
			win.inputs.rename.SetText(win.walletInfo.Wallets[win.selected].Name)
			win.outputs.toggleWalletRename.Icon = win.outputs.ic.clear
			win.outputs.toggleWalletRename.Color = win.theme.Color.Danger
		}
		win.states.dialog = true
	}

	// RENAME WALLET

	// if win.inputs.toggleWalletRename.Clicked(win.gtx) {
	// if win.states.renamingWallet {
	// 	win.outputs.toggleWalletRename.Icon = win.outputs.ic.create
	// 	win.outputs.toggleWalletRename.Color = win.theme.Color.Primary
	// } else {
	// 	win.inputs.rename.SetText(win.walletInfo.Wallets[win.selected].Name)
	// 	win.outputs.rename.TextSize = unit.Dp(48)
	// 	win.outputs.toggleWalletRename.Icon = win.outputs.ic.clear
	// 	win.outputs.toggleWalletRename.Color = win.theme.Color.Danger
	// }

	win.states.renamingWallet = !win.states.renamingWallet
	// }

	if win.inputs.renameWallet.Clicked(win.gtx) {
		name := win.inputs.rename.Text()
		if name == "" {
			return
		}
		err := win.wallet.RenameWallet(win.walletInfo.Wallets[win.selected].ID, name)
		if err != nil {
			log.Debug("Error renaming wallet")
		} else {
			win.walletInfo.Wallets[win.selected].Name = name
			win.states.renamingWallet = false
			win.outputs.toggleWalletRename.Icon = win.outputs.ic.create
			win.outputs.toggleWalletRename.Color = win.theme.Color.Primary
			win.reloadTabs()
		}
	}

	// VERIFY MESSAGE

	if win.inputs.verifyMessDiag.Clicked(win.gtx) {
		win.states.dialog = true
		win.dialog = win.verifyMessageDiag
		return
	}

	if strings.Trim(win.inputs.addressInput.Text(), " ") == "" || strings.Trim(win.inputs.signInput.Text(), " ") == "" || strings.Trim(win.inputs.messageInput.Text(), " ") == "" {
		win.outputs.verifyBtn.Background = win.theme.Color.Hint
		win.outputs.verifyMessage.Text = ""
		win.err = ""
		if win.inputs.verifyBtn.Clicked(win.gtx) {
			return
		}
	} else {
		win.outputs.verifyBtn.Background = win.theme.Color.Primary
		if win.inputs.verifyBtn.Clicked(win.gtx) {
			addr := win.inputs.addressInput.Text()
			if addr == "" {
				return
			}
			sign := win.inputs.signInput.Text()
			if sign == "" {
				return
			}
			msg := win.inputs.messageInput.Text()
			if msg == "" {
				return
			}

			valid, err := win.wallet.VerifyMessage(addr, msg, sign)
			if err != nil {
				win.err = err.Error()
				return
			}
			if !valid {
				win.outputs.verifyMessage.Text = "Invalid Signature"
				win.outputs.verifyMessage.Color = win.theme.Color.Danger
				return
			}

			win.outputs.verifyMessage.Text = "Valid Signature"
			win.outputs.verifyMessage.Color = win.theme.Color.Success
		}
	}

	data, err := clipboard.ReadAll()
	if err != nil {
		win.err = err.Error()
	}

	//signature control
	if win.inputs.clearSign.Clicked(win.gtx) {
		win.inputs.signInput.SetText("")
		return
	}
	if win.inputs.pasteSign.Clicked(win.gtx) {
		win.inputs.signInput.SetText(data)
		return
	}
	//address control
	if win.inputs.clearAddr.Clicked(win.gtx) {
		win.inputs.addressInput.SetText("")
		return
	}
	if win.inputs.pasteAddr.Clicked(win.gtx) {
		win.inputs.addressInput.SetText(data)
		return
	}
	//mesage control
	if win.inputs.clearMsg.Clicked(win.gtx) {
		win.inputs.messageInput.SetText("")
		return
	}
	if win.inputs.pasteMsg.Clicked(win.gtx) {
		win.inputs.messageInput.SetText(data)
		return
	}

	if win.inputs.clearBtn.Clicked(win.gtx) {
		win.resetVerifyFields()
		return
	}

	if win.inputs.verifyInfo.Clicked(win.gtx) {
		win.dialog = win.msgInfoDiag
	}

	if win.inputs.hideMsgInfo.Clicked(win.gtx) {
		win.dialog = win.verifyMessageDiag
	}

	// DELETE WALLET

	if win.inputs.deleteDiag.Clicked(win.gtx) {
		win.states.dialog = true
		win.dialog = win.DeleteDiag
		return
	}

	if win.inputs.deleteWallet.Clicked(win.gtx) {
		pass := win.validatePassword()
		if pass == "" {
			return
		}
		win.wallet.DeleteWallet(win.walletInfo.Wallets[win.selected].ID, pass)
		win.resetPasswords()
		win.states.loading = true
		log.Debug("Delete Wallet clicked")
		return
	}

	// ADD ACCOUNT

	if win.inputs.addAcctDiag.Clicked(win.gtx) {
		win.outputs.dialog.Hint = "Enter account name"
		win.dialog = win.AddAccountDiag
		win.states.dialog = true
	}

	if win.inputs.addAccount.Clicked(win.gtx) {
		pass := win.validatePassword()
		if pass == "" {
			return
		}
		win.wallet.AddAccount(win.walletInfo.Wallets[win.selected].ID, win.inputs.dialog.Text(), pass)
		win.states.loading = true
	}

	// NAVIGATION

	if win.inputs.toWallets.Clicked(win.gtx) {
		//win.current = win.WalletsPage
		win.current = PageWallet
		return
	}

	if win.inputs.toOverview.Clicked(win.gtx) {
		win.current = PageOverview
		return
	}

	if win.inputs.toReceive.Clicked(win.gtx) {
		win.current = PageReceive
		return
	}

	if win.inputs.toSend.Clicked(win.gtx) {
		win.current = PageSend
		return
	}

	if win.inputs.toTransactions.Clicked(win.gtx) {
		win.current = PageTransactions
		return
	}

	if win.inputs.toTransactionsFilters.Clicked(win.gtx) {
		win.states.dialog = true
		win.dialog = win.transactionsFilters
	}

	if win.inputs.applyFiltersTransactions.Clicked(win.gtx) {
		win.states.dialog = false
		win.sortTransactions()
	}

	if win.inputs.sync.Clicked(win.gtx) || win.inputs.syncHeader.Clicked(win.gtx) {
		//log.Info("Sync clicked :", win.walletInfo.Synced, win.walletInfo.Syncing)
		if win.walletInfo.Synced || win.walletInfo.Syncing {
			win.wallet.CancelSync()
			win.outputs.sync = win.theme.Button("Reconnect")
		} else {
			win.wallet.StartSync()
			win.outputs.sync = win.theme.DangerButton("Cancel")
			win.outputs.syncHeader = win.outputs.icons.cancel
		}
	}

	if win.inputs.cancelDialog.Clicked(win.gtx) {
		win.states.dialog = false
		win.err = ""
		win.resetVerifyFields()
		log.Debug("Cancel dialog clicked")
		return
	}

	// RECEIVE PAGE
	if win.inputs.receiveIcons.info.Clicked(win.gtx) {
		win.states.dialog = true
		win.dialog = win.infoDiag
	}

	if win.inputs.receiveIcons.more.Clicked(win.gtx) {
		newAddr = !newAddr
	}

	if win.inputs.receiveIcons.gotItDiag.Clicked(win.gtx) {
		win.states.dialog = false
	}

	if win.inputs.receiveIcons.newAddressDiag.Clicked(win.gtx) {
		wallet := win.walletInfo.Wallets[win.selected]
		account := wallet.Accounts[win.selectedAccount]
		addr, err := win.wallet.NextAddress(wallet.ID, account.Number)
		if err != nil {
			log.Debug("Error generating new address" + err.Error())
			win.err = err.Error()
		} else {
			win.walletInfo.Wallets[win.selected].Accounts[win.selectedAccount].CurrentAddress = addr
			newAddr = false
		}
	}

	for win.inputs.receiveIcons.copy.Clicked(win.gtx) {
		clipboard.WriteAll(win.walletInfo.Wallets[win.selected].Accounts[win.selectedAccount].CurrentAddress)
		win.addressCopiedLabel.Text = "Address Copied"
		time.AfterFunc(time.Second*3, func() {
			win.addressCopiedLabel.Text = ""
		})
		return
	}
}

func (win *Window) validatePasswords() string {
	pass := win.inputs.spendingPassword.Text()
	if pass == "" {
		win.outputs.spendingPassword.HintColor = win.theme.Color.Danger
		return pass
	}

	match := win.inputs.matchSpending.Text()
	if match == "" {
		win.outputs.matchSpending.HintColor = win.theme.Color.Danger
		return ""
	}

	if match != pass {
		win.outputs.matchSpending.Color = win.theme.Color.Danger
		return ""
	}

	return pass
}

func (win *Window) validatePassword() string {
	pass := win.inputs.spendingPassword.Text()
	if pass == "" {
		win.outputs.spendingPassword.HintColor = win.theme.Color.Danger
	}
	return pass
}

func (win *Window) resetPasswords() {
	win.outputs.spendingPassword.HintColor = win.theme.Color.InvText
	win.inputs.spendingPassword.SetText("")
	win.outputs.matchSpending.HintColor = win.theme.Color.InvText
	win.inputs.matchSpending.SetText("")
}

func (win *Window) resetVerifyFields() {
	win.inputs.addressInput.SetText("")
	win.inputs.signInput.SetText("")
	win.inputs.messageInput.SetText("")
	win.outputs.verifyMessage.Text = ""
}

func (win *Window) validateSeeds() string {
	text := ""
	win.err = ""

	for i, editor := range win.inputs.seedEditors.editors {
		if editor.Text() == "" {
			win.outputs.seedEditors[i].HintColor = win.theme.Color.Danger
			return ""
		}

		text += editor.Text() + " "
	}

	if !dcrlibwallet.VerifySeed(text) {
		win.err = "Invalid seed phrase"
		return ""
	}

	return text
}

func (win *Window) resetSeeds() {
	for i := 0; i < len(win.inputs.seedEditors.editors); i++ {
		win.inputs.seedEditors.editors[i].SetText("")
	}
}

func (win *Window) editorSeedsEventsHandler() {
	for i := 0; i < len(win.inputs.seedEditors.editors); i++ {
		editor := &win.inputs.seedEditors.editors[i]

		if editor.Focused() && win.inputs.seedEditors.focusIndex != i {
			win.inputs.seedsSuggestions = nil
			win.outputs.seedsSuggestions = nil
			win.inputs.seedEditors.focusIndex = i

			return
		}

		for _, e := range editor.Events(win.gtx) {
			switch e.(type) {
			case widget.ChangeEvent:
				win.inputs.seedsSuggestions = nil
				win.outputs.seedsSuggestions = nil

				if strings.Trim(editor.Text(), " ") == "" {
					return
				}

				for _, word := range dcrlibwallet.PGPWordList() {
					if strings.HasPrefix(strings.ToLower(word), strings.ToLower(editor.Text())) {
						if len(win.inputs.seedsSuggestions) < 2 {
							var btn struct {
								text   string
								button widget.Button
							}

							btn.text = word
							win.inputs.seedsSuggestions = append(win.inputs.seedsSuggestions, btn)
							win.outputs.seedsSuggestions = append(win.outputs.seedsSuggestions, win.theme.Button(word))
						}
					}
				}

			case widget.SubmitEvent:
				if i < len(win.inputs.seedEditors.editors)-1 {
					win.inputs.seedEditors.editors[i+1].Focus()
				}
			}
		}
	}
}

func (win *Window) onSuggestionSeedsClicked() {
	for i := 0; i < len(win.inputs.seedsSuggestions); i++ {
		btn := win.inputs.seedsSuggestions[i]
		if btn.button.Clicked(win.gtx) {
			for i := 0; i < len(win.inputs.seedEditors.editors); i++ {
				editor := &win.inputs.seedEditors.editors[i]
				if editor.Focused() {
					editor.SetText(btn.text)
					editor.Move(len(btn.text))

					if i < len(win.inputs.seedEditors.editors)-1 {
						win.inputs.seedEditors.editors[i+1].Focus()
					} else {
						win.inputs.seedEditors.focusIndex = -1
					}
				}
			}
		}
	}
}

// KeysEventsHandler handlers all pressed keys events
func (win *Window) KeysEventsHandler(evt *key.Event) {
	if evt.Name == key.NameTab {
		for i := 0; i < len(win.inputs.seedEditors.editors); i++ {
			editor := &win.inputs.seedEditors.editors[i]
			if editor.Focused() && win.inputs.seedsSuggestions != nil {
				editor.SetText(win.inputs.seedsSuggestions[0].text)
				editor.Move(len(win.inputs.seedsSuggestions[0].text))
			}
		}
	}
}

func (win *Window) sortTransactions() {
	newestFirst := true
	if win.inputs.transactionFilterSort.Value(win.gtx) == "1" {
		newestFirst = false
	}

	walletSelected := win.walletInfo.Wallets[win.selected].ID
	transactions := win.walletTransactions.Txs[walletSelected]
	sort.SliceStable(transactions, func(i, j int) bool {
		backTime := time.Unix(transactions[j].Txn.Timestamp, 0)
		frontTime := time.Unix(transactions[i].Txn.Timestamp, 0)
		if newestFirst {
			return backTime.Before(frontTime)
		}
		return frontTime.Before(backTime)
	})
}
