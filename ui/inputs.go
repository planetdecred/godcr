package ui

// HandleInputs handles all ui inputs
func (win *Window) HandleInputs() {
	for i, tab := range win.buttons.tabs {
		if tab.Clicked(win.gtx) {
			win.selected = i
			log.Debugf("Tab %d selected", i)
		}
	}
	if win.buttons.createWallet.Clicked(win.gtx) {
		log.Debug("Create Wallet clicked")
	}

	if win.buttons.restoreWallet.Clicked(win.gtx) {
		log.Debug("Restore Wallet clicked")
	}

	if win.buttons.deleteWallet.Clicked(win.gtx) {
		log.Debug("Delete Wallet clicked")
	}

	if win.buttons.cancelDialog.Clicked(win.gtx) {
		log.Debug("Cancel dialog clicked")
	}
	if win.buttons.confirmDialog.Clicked(win.gtx) {
		log.Debug("Confirm dialog clicked")
	}
}
