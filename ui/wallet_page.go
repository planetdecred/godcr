package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
)

// WalletsPage lays out the main wallet page
func (win *Window) WalletsPage() {
	tabbed := func() {
		if win.walletInfo.LoadedWallets == 0 {
			return
		}
		win.TabbedPage(
			func() {
				info := win.walletInfo.Wallets[win.selected]
				layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.H3(info.Name).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.theme.H5(info.Balance).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						(&layout.List{Axis: layout.Vertical}).Layout(win.gtx, len(info.Accounts), func(i int) {
							acct := info.Accounts[i]
							layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
								layout.Rigid(func() {
									win.theme.Body1(acct.Name).Layout(win.gtx)
								}),
								layout.Rigid(func() {
									win.theme.Body1(acct.TotalBalance).Layout(win.gtx)
								}),
							)
						})
					}),
					layout.Flexed(.3, func() {
						layout.Flex{}.Layout(win.gtx,
							layout.Rigid(func() {
								dbtn := win.theme.DangerButton("Delete wallet")
								dbtn.Layout(win.gtx, &win.inputs.deleteWallet)
							}),
							layout.Rigid(func() {
								win.theme.Editor("Enter password").Layout(win.gtx, &win.inputs.spendingPassword)
							}),
						)
					}),
				)
			},
		)
	}
	toMax(win.gtx)
	materialplus.Modal{}.Layout(win.gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Flexed(.3, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(.3, win.Header),
				)
			}),
			layout.Rigid(tabbed),
		)
	})

}
