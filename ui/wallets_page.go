package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
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
		win.vFlex(
			// Page heading: Wallet name and total balance
			rigid(func() {
				win.vFlex(
					rigid(func() {
						tRename := rigid(func() {
							layout.Center.Layout(win.gtx, func() {
								win.outputs.toggleWalletRename.Layout(win.gtx, &win.inputs.toggleWalletRename)
							})
						})
						if win.states.renamingWallet {
							win.hFlex(
								rigid(func() {
									win.outputs.rename.Layout(win.gtx, &win.inputs.rename)
								}),
								rigid(func() {
									win.outputs.renameWallet.Layout(win.gtx, &win.inputs.renameWallet)
								}),
								tRename,
							)
						} else {
							win.hFlex(
								rigid(func() {
									win.theme.H3(info.Name).Layout(win.gtx)
								}),
								tRename,
							)
						}
					}),
					rigid(func() {
						win.theme.H5("Total Balance: " + info.Balance).Layout(win.gtx)
					}),
				)
			}),
			// Accounts list
			rigid(func() {
				// List header
				h := func() {
					win.hFlex(
						rigid(func() {
							win.theme.H5("Accounts").Layout(win.gtx)
						}),
						layout.Rigid(func() {
							layout.S.Layout(win.gtx, func() {
								win.outputs.addAcctDiag.Layout(win.gtx, &win.inputs.addAcctDiag)
							})
						}),
					)
				}
				b := func() {
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
								layout.Rigid(func() {
									win.theme.Body1("Keys: " + acct.Keys.External + " external, " + acct.Keys.Internal + " internal, " + acct.Keys.Imported + " imported").Layout(win.gtx)
								}),
								layout.Rigid(func() {
									win.theme.Body1("HD Path: " + acct.HDPath).Layout(win.gtx)
								}),
							)
						}
						layout.Inset{Top: unit.Dp(3)}.Layout(win.gtx, a)
					})
				}
				layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
					layout.Rigid(h),
					layout.Rigid(b),
				)
			}),
			// Action Buttons
			layout.Rigid(func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Rigid(func() {
						win.outputs.deleteDiag.Layout(win.gtx, &win.inputs.deleteDiag)
					}),
					layout.Rigid(func() {
						inset := layout.Inset{
							Left: unit.Dp(10),
						}
						inset.Layout(win.gtx, func() {
							win.outputs.verifyMessDiag.Layout(win.gtx, &win.inputs.verifyMessDiag)
						})
					}),
				)
			}),
		)
	}
	win.TabbedPage(body)
}
