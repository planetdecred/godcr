package ui

import (
	"strings"

	"gioui.org/widget"
)

// HandleInputs handles all ui inputs
func (win *Window) HandleInputs() {
	if win.tabs.Changed() {
		win.selected = win.tabs.Selected
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
		win.current = win.TransactionsPage
		return
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

	for i, editor := range win.inputs.seeds {
		if editor.Text() == "" {
			win.outputs.seeds[i].HintColor = win.theme.Color.Danger
			return ""
		}

		text += editor.Text() + " "
	}

	return text
}

func (win *Window) resetSeeds() {
	for i := 0; i < len(win.inputs.seeds); i++ {
		win.inputs.seeds[i].SetText("")
	}
}

func (win *Window) editorSeedsEventsHandler() {
	for i, editor := range win.inputs.seeds {
		if editor.Focused() &&
			(win.combined.editorsEventsHandlerIndex != i || strings.Trim(editor.Text(), " ") == "") {
			win.combined.editorsEventsHandlerIndex = -1
		}

		for _, e := range editor.Events(win.gtx) {
			switch e.(type) {
			case widget.ChangeEvent:
				win.combined.editorsEventsHandlerIndex = i
			case widget.SubmitEvent:
				if i < len(win.inputs.seeds)-1 {
					win.inputs.seeds[i+1].Focus()
				}
			}
		}
	}
}

func (win *Window) onSuggestionSeedsClicked() {
	for i := 0; i < len(win.combined.autocompleteButtons); i++ {
		if win.combined.autocompleteButtons[i].Clicked(win.gtx) {
			win.inputs.seeds[win.combined.editorsEventsHandlerIndex].SetText(win.combined.suggestionsWords[i])
			win.inputs.seeds[win.combined.editorsEventsHandlerIndex].Move(len(win.combined.suggestionsWords[i]))
			win.combined.suggestionsWords = nil

			if win.combined.editorsEventsHandlerIndex < len(win.inputs.seeds)-1 {
				win.inputs.seeds[win.combined.editorsEventsHandlerIndex+1].Focus()
			}
		}
	}
}
