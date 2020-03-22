package ui

import (
	"strings"

	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// HandleInputs handles all ui inputs
func (win *Window) HandleInputs() {
	if win.tabs.Changed() {
		win.selected = win.tabs.Selected
		win.wallet.GetTransactionsByWallet(win.walletInfo.Wallets[win.tabs.Selected].ID, 0, 100, 0, 0)
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
		win.current = win.RestorePage
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

	// RENAME WALLET

	if win.inputs.toggleWalletRename.Clicked(win.gtx) {
		if win.states.renamingWallet {
			win.outputs.toggleWalletRename.Icon = win.outputs.ic.create
			win.outputs.toggleWalletRename.Color = win.theme.Color.Primary
		} else {
			win.inputs.rename.SetText(win.walletInfo.Wallets[win.selected].Name)
			win.outputs.rename.TextSize = unit.Dp(48)
			win.outputs.toggleWalletRename.Icon = win.outputs.ic.clear
			win.outputs.toggleWalletRename.Color = win.theme.Color.Danger
		}

		win.states.renamingWallet = !win.states.renamingWallet
	}

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
		win.current = win.WalletsPage
		return
	}

	if win.inputs.toOverview.Clicked(win.gtx) {
		win.current = win.Overview
		return
	}

	if win.inputs.toSend.Clicked(win.gtx) {
		win.current = win.SendPage
		return
	}

	if win.inputs.toTransactions.Clicked(win.gtx) {
		win.wallet.GetTransactionsByWallet(win.walletInfo.Wallets[win.tabs.Selected].ID, 0, 100, 0, 0)
		win.current = win.TransactionsPage
		return
	}

	if win.combined.transactionSort.Changed() {
		win.wallet.GetTransactionsByWallet(
			win.walletInfo.Wallets[win.tabs.Selected].ID, 0, 100,
			int32(win.combined.transactionStatus.Selected()),
			win.combined.transactionSort.Selected(),
		)
		log.Info(win.tabs.Selected)
	}

	if win.combined.transactionStatus.Changed() {
		log.Info(win.tabs.Selected)
		win.wallet.GetTransactionsByWallet(
			win.walletInfo.Wallets[win.tabs.Selected].ID, 0, 100,
			int32(win.combined.transactionStatus.Selected()),
			win.combined.transactionSort.Selected(),
		)
	}

	for i := 0; i < len(win.combined.transactions); i++ {
		for _, e := range win.combined.transactions[i].gesture.Events(win.gtx) {
			if e.Type == gesture.TypeClick {
				transaction := win.combined.transactions[i].data.(*wallet.TransactionInfo)
				log.Infof("To transaction details %+v", transaction)
			}
		}
	}

	// SYNC
	if win.inputs.sync.Clicked(win.gtx) {
		//log.Info("Sync clicked :", win.walletInfo.Synced, win.walletInfo.Syncing)
		if win.walletInfo.Syncing {
			win.wallet.CancelSync()
		} else if !win.walletInfo.Synced {
			win.wallet.StartSync()
			cancel := win.outputs.icons.cancel
			win.outputs.sync = cancel
		}
	}

	if win.inputs.cancelDialog.Clicked(win.gtx) {
		win.states.dialog = false
		win.err = ""
		log.Debug("Cancel dialog clicked")
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
