package ui

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/skip2/go-qrcode"
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
	body := func() {
		layout.Stack{}.Layout(p.gtx,
			layout.Expanded(func() {
				layout.Flex{}.Layout(p.gtx,
					layout.Flexed(0.9, func() {
						p.ReceivePageContents(common)
					}),
					layout.Rigid(func() {
						p.rightNav()
					}),
				)
			}),
		)
	}
	common.LayoutWithWallets(p.gtx, body)
}

func (p *receivePage) ReceivePageContents(common pageCommon) {
	pageContent := []func(){
		func() {
			p.selectedAcountColumn(common)
		},
		func() {
			p.qrCodeAddressColumn(common)
		},
		func() {
			layout.Flex{}.Layout(p.gtx,
				layout.Flexed(0.35, func() {
				}),
				layout.Flexed(1, func() {
					if p.addressCopiedLabel.Text != "" {
						p.addressCopiedLabel.Layout(p.gtx)
					}
				}),
			)
		},
		func() {
			layout.Flex{}.Layout(p.gtx,
				layout.Flexed(0.35, func() {
				}),
				layout.Flexed(1, func() {
					// win.Err()
				}),
			)
		},
	}
	p.pageContainer.Layout(p.gtx, len(pageContent), func(i int) {
		layout.Inset{Left: unit.Dp(3)}.Layout(p.gtx, pageContent[i])
	})

	// if newAddr {
	// 	win.generateNewAddress()
	// }
}

func (p *receivePage) rightNav() {
	layout.Flex{Axis: layout.Vertical}.Layout(p.gtx,
		layout.Rigid(func() {
			p.moreBtn.Layout(p.gtx, &p.moreBtnW)
		}),
		layout.Rigid(func() {
			// layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(20)}.Layout(p.gtx, func() {
			p.infoBtn.Layout(p.gtx, &p.infoBtnW)
		}),
	)
}

func (p *receivePage) selectedAcountColumn(common pageCommon) {
	current := common.info.Wallets[*common.selectedWallet]

	p.selectedWalletNameLabel.Text = current.Name
	p.selectedWalletBalLabel.Text = current.Balance

	account := common.info.Wallets[*common.selectedWallet].Accounts[0]
	p.selectedAccountNameLabel.Text = account.Name
	p.selectedAccountBalanceLabel.Text = account.SpendableBalance

	layout.Flex{}.Layout(p.gtx,
		layout.Flexed(0.22, func() {
		}),
		layout.Flexed(1, func() {
			layout.Stack{}.Layout(p.gtx,
				layout.Stacked(func() {
					selectedDetails := func() {
						layout.UniformInset(unit.Dp(10)).Layout(p.gtx, func() {
							layout.Flex{Axis: layout.Vertical}.Layout(p.gtx,
								layout.Rigid(func() {
									layout.Flex{}.Layout(p.gtx,
										layout.Rigid(func() {
											layout.Inset{Bottom: unit.Dp(5)}.Layout(p.gtx, func() {
												p.selectedAccountNameLabel.Layout(p.gtx)
											})
										}),
										layout.Rigid(func() {
											layout.Inset{Left: unit.Dp(20)}.Layout(p.gtx, func() {
												p.selectedAccountBalanceLabel.Layout(p.gtx)
											})
										}),
									)
								}),
								layout.Rigid(func() {
									layout.Inset{Left: unit.Dp(20)}.Layout(p.gtx, func() {
										layout.Flex{}.Layout(p.gtx,
											layout.Rigid(func() {
												layout.Inset{Bottom: unit.Dp(5)}.Layout(p.gtx, func() {
													p.selectedWalletNameLabel.Layout(p.gtx)
												})
											}),
											layout.Rigid(func() {
												layout.Inset{Left: unit.Dp(22)}.Layout(p.gtx, func() {
													p.selectedWalletBalLabel.Layout(p.gtx)
												})
											}),
										)
									})
								}),
							)
						})
					}
					decredmaterial.Card{}.Layout(p.gtx, selectedDetails)
				}),
			)
		}),
	)
}

func (p *receivePage) qrCodeAddressColumn(common pageCommon) {
	addrs := common.info.Wallets[*common.selectedWallet].Accounts[0].CurrentAddress
	qrCode, err := qrcode.New(addrs, qrcode.Highest)
	if err != nil {
		// win.err = err.Error()
		return
	}
	// win.err = ""
	qrCode.DisableBorder = true
	layout.Flex{Axis: layout.Vertical}.Layout(p.gtx,
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(16)}.Layout(p.gtx, func() {
				layout.Flex{}.Layout(p.gtx,
					layout.Flexed(0.2, func() {
					}),
					layout.Flexed(1, func() {
						img := common.theme.Image(paint.NewImageOp(qrCode.Image(520)))
						img.Src.Rect.Max.X = 521
						img.Src.Rect.Max.Y = 521
						img.Scale = 0.5
						img.Layout(p.gtx)
					}),
				)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(16)}.Layout(p.gtx, func() {
				layout.Flex{}.Layout(p.gtx,
					layout.Flexed(0.1, func() {
					}),
					layout.Flexed(1, func() {
						p.receiveAddressColumn(addrs)
					}),
				)
			})
		}),
	)
}

func (p *receivePage) receiveAddressColumn(addrs string) {
	layout.Flex{}.Layout(p.gtx,
		layout.Rigid(func() {
			p.receiveAddressLabel.Text = addrs
			p.receiveAddressLabel.Layout(p.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(16)}.Layout(p.gtx, func() {
				p.copyBtn.Layout(p.gtx, &p.copyBtnW)
			})
		}),
	)
}

func (p *receivePage) generateNewAddress() {
	layout.Flex{}.Layout(p.gtx,
		layout.Flexed(0.71, func() {
		}),
		layout.Flexed(1, func() {
			inset := layout.Inset{
				Top: unit.Dp(45),
			}
			inset.Layout(p.gtx, func() {
				p.gtx.Constraints.Width.Min = 40
				p.gtx.Constraints.Height.Min = 40
				p.newAddrBtn.Layout(p.gtx, &p.newAddrBtnW)
			})
		}),
	)
}
