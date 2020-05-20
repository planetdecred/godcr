package ui

import (
	"image/color"
	"strings"
	"time"

	"gioui.org/io/key"
	"gioui.org/unit"
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
		}
	}

	if win.combined.sel.Changed() {
		win.selectedAccount = win.combined.sel.Selected
	}

	for _, evt := range win.inputs.spendingPassword.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.err = ""
			win.outputs.err.Color = win.theme.Color.Danger
			win.resetButton()
			win.outputs.spendingPassword.HintColor = win.theme.Color.Hint
			return
		}
		log.Debug("Pass evt", evt)
	}

	for _, evt := range win.inputs.oldSpendingPassword.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.err = ""
			win.outputs.err.Color = win.theme.Color.Danger
			win.resetButton()
			win.outputs.oldSpendingPassword.HintColor = win.theme.Color.Hint
			return
		}
		log.Debug("Pass evt", evt)
	}

	for _, evt := range win.inputs.matchSpending.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.err = ""
			win.outputs.err.Color = win.theme.Color.Danger
			win.outputs.matchSpending.Color = win.theme.Color.Text
			win.resetButton()
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

	// VERIFY MESSAGE

	if win.inputs.verifyMessDiag.Clicked(win.gtx) {
		win.current = PageVerifyMessage
		return
	}

	// CHANGE WALLET PASSWORD

	for win.inputs.changePasswordDiag.Clicked(win.gtx) {
		win.err = ""
		win.dialog = win.editPasswordDiag
		win.states.dialog = true
	}

	for win.inputs.savePassword.Clicked(win.gtx) {
		op := win.inputs.oldSpendingPassword.Text()
		if op == "" {
			win.outputs.oldSpendingPassword.HintColor = win.theme.Color.Danger
			win.err = "Old Wallet password required and cannot be empty"
			return
		}
		np := win.validatePasswords()
		if np == "" {
			return
		}

		err := win.wallet.ChangeWalletPassphrase(win.walletInfo.Wallets[win.selected].ID, op, np)
		if err != nil {
			log.Debug("Error changing wallet password " + err.Error())
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				win.err = "Passphrase is incorrect"
			} else {
				win.err = err.Error()
			}
			return
		}

		win.err = "Password changed successfully"
		win.outputs.err.Color = win.theme.Color.Success
		win.outputs.savePassword.Text = "Changed"
		win.outputs.savePassword.Background = win.theme.Color.Success
		win.resetPasswords()
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

	if win.inputs.cancelDialog.Clicked(win.gtx) {
		win.states.dialog = false
		win.err = ""
		win.resetButton()
		win.resetPasswords()
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

	if win.inputs.receiveIcons.copy.Clicked(win.gtx) {
		clipboard.WriteAll(win.walletInfo.Wallets[win.selected].Accounts[win.selectedAccount].CurrentAddress)
		win.addressCopiedLabel.Text = "Address Copied"
		time.AfterFunc(time.Second*3, func() {
			win.addressCopiedLabel.Text = ""
		})
		return
	}

	if strings.Trim(win.inputs.spendingPassword.Text(), " ") != "" {
		strength := dcrlibwallet.ShannonEntropy(win.inputs.spendingPassword.Text()) / 4.0
		switch {
		case strength > 0.6:
			win.outputs.passwordBar.ProgressColor = win.theme.Color.Success
		case strength > 0.3 && strength <= 0.6:
			win.outputs.passwordBar.ProgressColor = color.RGBA{248, 231, 27, 255}
		default:
			win.outputs.passwordBar.ProgressColor = win.theme.Color.Danger
		}

		win.outputs.passwordBar.Progress = strength * 100
	}

	if win.inputs.signMessageDiag.Clicked(win.gtx) {
		win.current = PageSignMessage
		return
	}
}

func (win *Window) validatePasswords() string {
	pass := win.validatePassword()
	if pass == "" {
		return ""
	}

	match := win.inputs.matchSpending.Text()
	if match == "" {
		win.outputs.matchSpending.HintColor = win.theme.Color.Danger
		win.err = "Enter new wallet password again and it cannot be empty."
		return ""
	}

	if match != pass {
		win.err = "New wallet passwords does no match. Try again."
		return ""
	}

	return pass
}

func (win *Window) validatePassword() string {
	pass := win.inputs.spendingPassword.Text()
	if pass == "" {
		win.outputs.spendingPassword.HintColor = win.theme.Color.Danger
		win.err = "Wallet password required and cannot be empty."
		return ""
	}

	return pass
}

func (win *Window) validateOldPassword() string {
	pass := win.inputs.oldSpendingPassword.Text()
	if pass == "" {
		win.outputs.oldSpendingPassword.HintColor = win.theme.Color.Danger
		win.err = "Old Wallet password required and cannot be empty"
		return ""
	}

	return pass
}

func (win *Window) resetPasswords() {
	win.outputs.spendingPassword.HintColor = win.theme.Color.InvText
	win.inputs.spendingPassword.SetText("")
	win.outputs.matchSpending.HintColor = win.theme.Color.InvText
	win.inputs.matchSpending.SetText("")
	win.outputs.oldSpendingPassword.HintColor = win.theme.Color.InvText
	win.inputs.oldSpendingPassword.SetText("")
	win.outputs.passwordBar.Progress = 0
}

func (win *Window) resetButton() {
	win.outputs.savePassword.Text = "Change"
	win.outputs.savePassword.Background = win.theme.Color.Primary
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
