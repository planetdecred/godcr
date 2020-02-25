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
	if win.inputs.createWallet.Clicked(win.gtx) {
		pass := win.inputs.spendingPassword.Text()
		win.wallet.CreateWallet(pass)
		log.Debug("Create Wallet clicked")
		win.states.loading = true
		return
	}

	if win.inputs.restoreWallet.Clicked(win.gtx) {
		log.Debug("Restore Wallet clicked")
		return
	}

	if win.inputs.deleteWallet.Clicked(win.gtx) {
		pass := win.inputs.spendingPassword.Text()
		win.wallet.DeleteWallet(win.selected, pass)
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
