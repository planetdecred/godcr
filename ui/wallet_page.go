package ui

import (
	"gioui.org/layout"
)

// WalletsPage layouts the main wallet page
func (win *Window) WalletsPage() layout.Widget {
	log.Debug("On Wallets")
	tabbed := func() {
		if len(win.walletInfo.Wallets) == 0 {
			return
		}
		win.TabbedWallets(
			func() {
				win.Background()
				info := win.walletInfo.Wallets[win.selected]
				win.theme.H5(info.Balance).Layout(win.gtx)
			},
			func() {
				win.Background()
				info := win.walletInfo.Wallets[win.selected]
				win.theme.H5(info.Name).Layout(win.gtx)
			},
			func(i int) {
				info := win.walletInfo.Wallets[i]
				layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.H5(info.Name).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.theme.H5(info.Balance).Layout(win.gtx)
					}),
				)

			},
		)
	}
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Flexed(.3, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(.3, win.Header),
				)
			}),
			layout.Rigid(tabbed),
		)
	}
}
