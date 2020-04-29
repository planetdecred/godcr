package ui

import (
	"gioui.org/layout"
)

const PageSend = "send"

func (win *Window) SendPage() {
	bd := func() {
		if len(win.walletInfo.Wallets) == 0 {
			win.theme.H3("No wallets").Layout(win.gtx)
			return
		}
		win.combined.sel.Layout(win.gtx, func() {

		})

		// win.selectedAcountColumn()
	}
	win.TabbedPage(bd)
}
