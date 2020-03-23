package ui

func (win *Window) SendPage() {
	bd := func() {
		if len(win.walletInfo.Wallets) == 0 {
			win.theme.H3("No wallets").Layout(win.gtx)
			return
		}

		win.selectedAcountDiag()
	}
	win.TabbedPage(bd)
}
