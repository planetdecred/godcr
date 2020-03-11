package ui

import (
	// "image"
	// "time"
	"fmt"
	// "gioui.org/io/pointer"
	"gioui.org/layout"
	// "gioui.org/op/paint"
	"gioui.org/unit"
	// "gioui.org/widget"
	// "gioui.org/widget/material"

	// "github.com/atotto/clipboard"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	// "github.com/raedahgroup/godcr-gio/ui"
	// "github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	// "github.com/raedahgroup/godcr-gio/ui/units"
	// "github.com/raedahgroup/godcr-gio/ui/values"
	"github.com/raedahgroup/godcr-gio/wallet"
	// "github.com/skip2/go-qrcode"
)

var (
	listContainer = &layout.List{Axis: layout.Vertical}
	selectedAcctNum int32
	generateNew bool
	addr string
	// pageTitle         = "Receiving DCR"
)

func (win *Window) Receive() {
	// if win.walletInfo.LoadedWallets == 0 {
		win.setDefaultPageValues()
	// }
	
	body := func() {
		layout.Stack{}.Layout(win.gtx,
			layout.Expanded(func() {
				win.ReceivePageContents()
			}),
		)
	}
	win.TabbedPage(body)
}

func (win *Window) ReceivePageContents() {
	ReceivePageContent := []func(){
		func() {
			win.pageFirstColumn()
		},

		func() {
			// layout.Align(layout.Center).Layout(win.gtx, func() {
			win.selectedAccountColumn()
			// })
		},

	}

	listContainer.Layout(win.gtx, len(ReceivePageContent), func(i int) {
		layout.UniformInset(unit.Dp(10)).Layout(win.gtx, ReceivePageContent[i])
	})

}

func (win *Window) pageFirstColumn() {
	layout.Flex{Spacing: layout.SpaceBetween}.Layout(win.gtx,
		layout.Rigid(func() {
			win.theme.H4("Receiving DCR").Layout(win.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{}.Layout(win.gtx, func() {
				layout.Flex{Spacing: layout.SpaceBetween}.Layout(win.gtx,
					layout.Rigid(func() {
						// if win.infoBtnWdg.Clicked(gtx) {
						// 	win.isInfoBtnModal = true
						// 	win.isGenerateNewAddBtnModal = false
						// 	win.isAccountModalOpen = false
						// }
						win.outputs.info.Layout(win.gtx, &win.inputs.info)
					}),
					layout.Rigid(func() {
						// if win.moreBtnWdg.Clicked(gtx) {
						// 	win.isGenerateNewAddBtnModal = true
						// 	win.isInfoBtnModal = false
						// 	win.isAccountModalOpen = false
						// }
						win.outputs.more.Layout(win.gtx, &win.inputs.more)
					}),
				)
			})
		}),
	)
}

func (win *Window) selectedAccountColumn() {
	// info := win.walletInfo.Wallets[win.selected]
	layout.Stack{Alignment: layout.Center}.Layout(win.gtx,
		layout.Stacked(func() {
			selectedDetails := func() {
				layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
					layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(win.gtx,
						layout.Rigid(func() {
									layout.Inset{}.Layout(win.gtx, func() {
										layout.Flex{}.Layout(win.gtx,
											layout.Rigid(func() {
												layout.Inset{Bottom: unit.Dp(5)}.Layout(win.gtx, func() {
													win.outputs.selectedAccountNameLabel.Layout(win.gtx)
												})
											}),
											layout.Rigid(func() {
												layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
													win.outputs.selectedAccountBalanceLabel.Layout(win.gtx)
												})
											}),
										)
									})

						}),
						layout.Rigid(func() {
							layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
								layout.Flex{}.Layout(win.gtx,
									layout.Rigid(func() {
										layout.Inset{Bottom: unit.Dp(5)}.Layout(win.gtx, func() {
											win.outputs.selectedWalletNameLabel.Layout(win.gtx)
										})
									}),
									layout.Rigid(func() {
										layout.Inset{Left: unit.Dp(22)}.Layout(win.gtx, func() {
											win.outputs.selectedWalletBalLabel.Layout(win.gtx)
										})
									}),
								)
							})
						}),
						// layout.Rigid(func() {
						// 	layout.Inset{Left: unit.Dp(15)}.Layout(win.gtx, func() {
						// if win.dropDownBtnWdg.Clicked(win.gtx) {
						// 	if win.isAccountModalOpen {
						// 		win.isAccountModalOpen = false
						// 	} else {
						// 		win.isAccountModalOpen = true
						// 		win.isInfoBtnModal = false
						// 		win.isGenerateNewAddBtnModal = false
						// 	}
						// }
						// 		win.outputs.dropdown.Layout(win.gtx, &win.inputs.dropdown)
						// 	})
						// }),
					)

				})
			}
			decredmaterial.Card{}.Layout(win.gtx, selectedDetails)
		}),
	)
}

func (win *Window) setDefaultPageValues() {
	wallets := win.walletInfo.Wallets

	for i := range wallets {
		if len(wallets[i].Accounts) == 0 {
			continue
		}

		win.setSelectedAccount(wallets[i], wallets[i].Accounts[0], false)
		break
	}
}

func (win *Window) setSelectedAccount(wallet wallet.InfoShort, account wallet.Account, generateNew bool) {

	fmt.Println(account.Name)
	win.outputs.selectedAccountNameLabel.Text = account.Name
	win.outputs.selectedWalletNameLabel.Text = wallet.Name 
	win.outputs.selectedWalletBalLabel.Text = dcrutil.Amount(account.SpendableBalance).String()
	win.outputs.selectedAccountBalanceLabel.Text = wallet.Balance
	
}

