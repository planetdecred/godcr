package ui

import (
	"gioui.org/layout"
	// "gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
	// "github.com/skip2/go-qrcode"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const PageReceive = "receive"

var pageContainer = &layout.List{Axis: layout.Vertical}

type receivePage struct {
	pageContainer layout.List
	gtx           *layout.Context

	newAddrBtn, gotItBtn                                 decredmaterial.Button
	copyBtn, infoBtn, moreBtn                            decredmaterial.IconButton
	copyBtnW, infoBtnW, moreBtnW, gotItBtnW, newAddrBtnW widget.Button

	// pageTitleLabel              decredmaterial.Label
	selectedAccountNameLabel, selectedAccountBalanceLabel decredmaterial.Label
	receiveAddressLabel, addressCopiedLabel               decredmaterial.Label
	selectedWalletBalLabel, selectedWalletNameLabel       decredmaterial.Label
}

func (win *Window) ReceivePage(common pageCommon) layout.Widget {
	moreBtn := common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationMoreVert)))
	moreBtn.Padding = unit.Dp(5)
	moreBtn.Size = unit.Dp(35)
	infoBtn := common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionInfo)))
	infoBtn.Padding = unit.Dp(5)
	infoBtn.Size = unit.Dp(35)
	copyBtn := common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy)))
	copyBtn.Padding = unit.Dp(5)
	copyBtn.Size = unit.Dp(30)
	copyBtn.Background = common.theme.Color.Background
	copyBtn.Color = common.theme.Color.Text
	receiveAddressLabel := common.theme.H6("")
	receiveAddressLabel.Color = common.theme.Color.Primary

	page := &receivePage{
		pageContainer: layout.List{
			Axis: layout.Vertical,
		},
		gtx:                         common.gtx,
		moreBtn:                     moreBtn,
		infoBtn:                     infoBtn,
		copyBtn:                     copyBtn,
		gotItBtn:                    common.theme.Button("Got It"),
		newAddrBtn:                  common.theme.Button("Generate new address"),
		receiveAddressLabel:         receiveAddressLabel,
		selectedAccountNameLabel:    common.theme.H6(""),
		selectedWalletNameLabel:     common.theme.Body2(""),
		selectedWalletBalLabel:      common.theme.Body2(""),
		selectedAccountBalanceLabel: common.theme.H6(""),
		addressCopiedLabel:          common.theme.Caption(""),
	}

	return func() {
		page.Layout(common)
		// page.Handle(common)
	}
}

func (p *receivePage) Layout(common pageCommon) {
	// 	body := func() {
	// 		layout.Stack{}.Layout(gtx,
	// 			layout.Expanded(func() {
	// 				layout.Flex{}.Layout(gtx,
	// 					layout.Rigid(func() {
	// 						win.combined.sel.Layout(gtx, func() {

	// 						})
	// 					}),
	// 					layout.Rigid(func() {
	// 						win.ReceivePageContents()
	// 					}),
	// 				)
	// 			}),
	// 		)
	// 	}
	// 	win.TabbedPage(body)
	// }

	// func (win *Window) ReceivePageContents() {
	pageContent := []func(){
		func() {
			p.pageHeaderColumn()
		},
		// func() {
		// 	win.selectedAcountColumn()
		// },
		// func() {
		// 	win.qrCodeAddressColumn()
		// },
		// func() {
		// 	layout.Flex{}.Layout(gtx,
		// 		layout.Flexed(0.35, func() {
		// 		}),
		// 		layout.Flexed(1, func() {
		// 			if win.addressCopiedLabel.Text != "" {
		// 				win.addressCopiedLabel.Layout(gtx)
		// 			}
		// 		}),
		// 	)
		// },
		// func() {
		// 	layout.Flex{}.Layout(gtx,
		// 		layout.Flexed(0.35, func() {
		// 		}),
		// 		layout.Flexed(1, func() {
		// 			win.Err()
		// 		}),
		// 	)
		// },
	}
	common.LayoutWithWallets(p.gtx, func() {
		p.pageContainer.Layout(common.gtx, len(pageContent), func(i int) {
			layout.Inset{Left: unit.Dp(3)}.Layout(p.gtx, pageContent[i])
		})
	})

	// if newAddr {
	// 	win.generateNewAddress()
	// }
}

func (p *receivePage) pageHeaderColumn() {
	layout.Flex{}.Layout(p.gtx,
		layout.Flexed(.6, func() {
			// win.outputs.pageTitle.Layout(gtx)
		}),
		layout.Flexed(.4, func() {
			layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(20)}.Layout(p.gtx, func() {
				layout.Flex{}.Layout(p.gtx,
					layout.Flexed(.5, func() {
						p.infoBtn.Layout(p.gtx, &p.infoBtnW)
					}),
					layout.Flexed(.5, func() {
						p.moreBtn.Layout(p.gtx, &p.moreBtnW)
					}),
				)
			})
		}),
	)
}

// func (win *Window) qrCodeAddressColumn() {
// 	addrs := win.walletInfo.Wallets[win.selected].Accounts[win.selectedAccount].CurrentAddress
// 	qrCode, err := qrcode.New(addrs, qrcode.Highest)
// 	if err != nil {
// 		win.err = err.Error()
// 		return
// 	}
// 	win.err = ""
// 	qrCode.DisableBorder = true
// 	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 		layout.Rigid(func() {
// 			layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func() {
// 				layout.Flex{}.Layout(gtx,
// 					layout.Flexed(0.2, func() {
// 					}),
// 					layout.Flexed(1, func() {
// 						img := win.theme.Image(paint.NewImageOp(qrCode.Image(520)))
// 						img.Src.Rect.Max.X = 521
// 						img.Src.Rect.Max.Y = 521
// 						img.Scale = 0.5
// 						img.Layout(gtx)
// 					}),
// 				)
// 			})
// 		}),
// 		layout.Rigid(func() {
// 			layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func() {
// 				layout.Flex{}.Layout(gtx,
// 					layout.Flexed(0.1, func() {
// 					}),
// 					layout.Flexed(1, func() {
// 						win.receiveAddressColumn(addrs)
// 					}),
// 				)
// 			})
// 		}),
// 	)
// }

// func (win *Window) receiveAddressColumn(addrs string) {
// 	layout.Flex{}.Layout(gtx,
// 		layout.Rigid(func() {
// 			win.outputs.receiveAddressLabel.Text = addrs
// 			win.outputs.receiveAddressLabel.Layout(gtx)
// 		}),
// 		layout.Rigid(func() {
// 			layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func() {
// 				win.outputs.copy.Layout(gtx, &win.inputs.receiveIcons.copy)
// 			})
// 		}),
// 	)
// }

// func (win *Window) selectedAcountColumn() {
// 	info := win.walletInfo.Wallets[win.selected]
// 	win.outputs.selectedWalletNameLabel.Text = info.Name
// 	win.outputs.selectedWalletBalLabel.Text = info.Balance

// 	account := win.walletInfo.Wallets[win.selected].Accounts[win.selectedAccount]
// 	win.outputs.selectedAccountNameLabel.Text = account.Name
// 	win.outputs.selectedAccountBalanceLabel.Text = account.SpendableBalance

// 	layout.Flex{}.Layout(gtx,
// 		layout.Flexed(0.22, func() {
// 		}),
// 		layout.Flexed(1, func() {
// 			layout.Stack{}.Layout(gtx,
// 				layout.Stacked(func() {
// 					selectedDetails := func() {
// 						layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
// 							layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 								layout.Rigid(func() {
// 									layout.Flex{}.Layout(gtx,
// 										layout.Rigid(func() {
// 											layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func() {
// 												win.outputs.selectedAccountNameLabel.Layout(gtx)
// 											})
// 										}),
// 										layout.Rigid(func() {
// 											layout.Inset{Left: unit.Dp(20)}.Layout(gtx, func() {
// 												win.outputs.selectedAccountBalanceLabel.Layout(gtx)
// 											})
// 										}),
// 									)
// 								}),
// 								layout.Rigid(func() {
// 									layout.Inset{Left: unit.Dp(20)}.Layout(gtx, func() {
// 										layout.Flex{}.Layout(gtx,
// 											layout.Rigid(func() {
// 												layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func() {
// 													win.outputs.selectedWalletNameLabel.Layout(gtx)
// 												})
// 											}),
// 											layout.Rigid(func() {
// 												layout.Inset{Left: unit.Dp(22)}.Layout(gtx, func() {
// 													win.outputs.selectedWalletBalLabel.Layout(gtx)
// 												})
// 											}),
// 										)
// 									})
// 								}),
// 							)
// 						})
// 					}
// 					decreddecredmaterial.Card{}.Layout(gtx, selectedDetails)
// 				}),
// 			)
// 		}),
// 	)
// }

// func (win *Window) generateNewAddress() {
// 	layout.Flex{}.Layout(gtx,
// 		layout.Flexed(0.71, func() {
// 		}),
// 		layout.Flexed(1, func() {
// 			inset := layout.Inset{
// 				Top: unit.Dp(45),
// 			}
// 			inset.Layout(gtx, func() {
// 				gtx.Constraints.Width.Min = 40
// 				gtx.Constraints.Height.Min = 40
// 				win.outputs.newAddressDiag.Layout(gtx, &win.inputs.receiveIcons.newAddressDiag)
// 			})
// 		}),
// 	)
// }
