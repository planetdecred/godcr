package ui

import "gioui.org/widget"

type inputs struct {
	deleteWallet, cancelDialog, confirmDialog widget.Button
	createWallet, restoreWallet               widget.Button
	tabs                                      []*widget.Button

	spendingPassword, renameWallet widget.Editor
}

// HandleInputs handles all ui inputs
func (win *Window) HandleInputs() {
	for i, tab := range win.inputs.tabs {
		for tab.Clicked(win.gtx) {
			win.selected = i
			log.Debugf("Tab %d selected", i)
			return
		}
	}
	pass := win.inputs.spendingPassword.Text()
	for _, evt := range win.inputs.spendingPassword.Events(win.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			win.widgets.spendingPassword.HintColor = win.theme.Color.Text
			return
		}
		log.Debug(evt)
	}

	if win.inputs.createWallet.Clicked(win.gtx) {
		if pass == "" {
			win.widgets.spendingPassword.HintColor = win.theme.Danger
			return
		}
		win.wallet.CreateWallet(pass)
		win.inputs.spendingPassword.SetText("")
		log.Debug("Create Wallet clicked")
		win.states.loading = true
		return
	}

	if win.inputs.restoreWallet.Clicked(win.gtx) {
		if pass == "" {
			win.widgets.spendingPassword.HintColor = win.theme.Danger
			return
		}
		log.Debug("Restore Wallet clicked")
		return
	}

	if win.inputs.deleteWallet.Clicked(win.gtx) {
		if pass == "" {
			win.widgets.spendingPassword.HintColor = win.theme.Danger
			return
		}
		pass := win.inputs.spendingPassword.Text()
		win.wallet.DeleteWallet(win.walletInfo.Wallets[win.selected].ID, pass)
		win.inputs.spendingPassword.SetText("")
		log.Debug("Delete Wallet clicked")
		return
	}

	if win.inputs.cancelDialog.Clicked(win.gtx) {
		log.Debug("Cancel dialog clicked")
		return
	}
	if win.inputs.confirmDialog.Clicked(win.gtx) {
		log.Debug("Confirm dialog clicked")
		return
	}
}
