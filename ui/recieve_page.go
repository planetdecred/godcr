package ui

import (
	// "image"
	// "time"
	// // "fmt"
	// "gioui.org/io/pointer"
	"gioui.org/layout"
	// "gioui.org/op/paint"
	"gioui.org/unit"
	// "gioui.org/widget"
	// "gioui.org/widget/material"

	// "github.com/atotto/clipboard"
	// "github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	// "github.com/raedahgroup/godcr-gio/ui"
	// "github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	// "github.com/raedahgroup/godcr-gio/ui/units"
	// "github.com/raedahgroup/godcr-gio/ui/values"
	// "github.com/raedahgroup/godcr-gio/wallet"
	// "github.com/skip2/go-qrcode"
)

var (
	listContainer = &layout.List{Axis: layout.Vertical}
	// pageTitle         = "Receiving DCR"
)

func (win *Window) Receive() {
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
		// func() {
		// 	layout.Align(layout.Center).Layout(gtx, func() {
		// 		if win.errorLabel.Text != "" {
		// 			win.errorLabel.Layout(gtx)
		// 		}
		// 	})
		// },
		func() {
			// layout.Align(layout.Center).Layout(win.gtx, func() {
			win.selectedAccountColumn()
			// })
		},
		// func() {
		// 	win.generateAddressQrCode(gtx)
		// },
		// func() {
		// 	layout.Align(layout.Center).Layout(gtx, func() {
		// 		if win.addressCopiedLabel.Text != "" {
		// 			win.addressCopiedLabel.Layout(gtx)
		// 		}
		// 	})
		// },
	}

	listContainer.Layout(win.gtx, len(ReceivePageContent), func(i int) {
		layout.UniformInset(unit.Dp(10)).Layout(win.gtx, ReceivePageContent[i])
	})
	// if win.isGenerateNewAddBtnModal {
	// 	win.drawMoreModal(gtx)
	// }
	// if win.isInfoBtnModal {
	// 	win.drawInfoModal(gtx)
	// }
	// if win.isAccountModalOpen {
	// 	win.accountSelectedModal(gtx)
	// }
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
	info := win.walletInfo.Wallets[win.selected]
	layout.Stack{Alignment: layout.Center}.Layout(win.gtx,
		layout.Stacked(func() {
			selectedDetails := func() {
				layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
					layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(win.gtx,
						layout.Rigid(func() {
							list := layout.List{Axis: layout.Vertical}
							list.Layout(win.gtx, len(info.Accounts), func(i int) {
								acct := info.Accounts[i]
								if acct.Name != "imported" {
									layout.Inset{}.Layout(win.gtx, func() {
										layout.Flex{}.Layout(win.gtx,
											layout.Rigid(func() {
												layout.Inset{Bottom: unit.Dp(5)}.Layout(win.gtx, func() {
													win.theme.H6(acct.Name).Layout(win.gtx)
												})
											}),
											layout.Rigid(func() {
												layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
													win.theme.H6(acct.TotalBalance).Layout(win.gtx)
												})
											}),
										)
									})
								}
							})
						}),
						layout.Rigid(func() {
							layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
								layout.Flex{}.Layout(win.gtx,
									layout.Rigid(func() {
										layout.Inset{Bottom: unit.Dp(5)}.Layout(win.gtx, func() {
											win.theme.Body2(info.Name).Layout(win.gtx)
										})
									}),
									layout.Rigid(func() {
										layout.Inset{Left: unit.Dp(22)}.Layout(win.gtx, func() {
											win.theme.Body2(info.Balance).Layout(win.gtx)
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
