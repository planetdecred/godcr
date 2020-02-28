package ui

func (win *Window) SendPage() {
	bd := func() {
		if len(win.walletInfo.Wallets) == 0 {
			win.theme.H3("No wallets").Layout(win.gtx)
			return
		}
		strs := make([]string, len(win.walletInfo.Wallets[win.selected].Accounts))

		for i, acct := range win.walletInfo.Wallets[win.selected].Accounts {
			strs[i] = acct.Name
		}
		win.combined.sel.Options = strs

		win.combined.sel.Layout(win.gtx, func() {

		})
	}
	win.TabbedPage(bd)
}
