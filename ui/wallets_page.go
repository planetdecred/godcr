package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
)

// WalletsPage lays out the main wallet page
func (win *Window) WalletsPage() {
	if win.walletInfo.LoadedWallets == 0 {
		win.Page(func() {
			win.outputs.noWallet.Layout(win.gtx)
		})
		return
	}
	body := func() {
		info := win.walletInfo.Wallets[win.selected]
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Flexed(.2, func() {
				win.theme.H3(info.Name).Layout(win.gtx)
			}),
			layout.Flexed(.2, func() {
				win.theme.H5(info.Balance).Layout(win.gtx)
			}),
			layout.Flexed(.5, func() {
				(&layout.List{Axis: layout.Vertical}).Layout(win.gtx, len(info.Accounts), func(i int) {
					acct := info.Accounts[i]
					a := func() {
						layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
							layout.Rigid(func() {
								win.theme.Body1(acct.Name).Layout(win.gtx)
							}),
							layout.Rigid(func() {
								win.theme.Body1(acct.TotalBalance).Layout(win.gtx)
							}),
						)
					}
					materialplus.Card{}.Layout(win.gtx, a)
					// a()
				})
			}),
			layout.Flexed(.1, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Rigid(func() {
						win.outputs.deleteDiag.Layout(win.gtx, &win.inputs.deleteDiag)
					}),
				)
			}),
		)
	}
	win.TabbedPage(body)
}
