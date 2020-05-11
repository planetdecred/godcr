package ui

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/atotto/clipboard"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/skip2/go-qrcode"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const PageReceive = "receive"

var pageContainer = &layout.List{Axis: layout.Vertical}

type receivePage struct {
	pageContainer layout.List
	gtx           *layout.Context

	isNewAddr, isInfo bool

	newAddrBtn, gotItBtn                                 decredmaterial.Button
	copyBtn, infoBtn, moreBtn                            decredmaterial.IconButton
	copyBtnW, infoBtnW, moreBtnW, gotItBtnW, newAddrBtnW widget.Button

	selectedAccountNameLabel, selectedAccountBalanceLabel decredmaterial.Label
	receiveAddressLabel, addressCopiedLabel, pageInfo     decredmaterial.Label
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
	pageInfo := common.theme.Body1("Each time you request a payment, a \nnew address is created to protect \nyour privacy.")
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
		pageInfo:                    pageInfo,
		selectedAccountNameLabel:    common.theme.H6(""),
		selectedWalletNameLabel:     common.theme.Body2(""),
		selectedWalletBalLabel:      common.theme.Body2(""),
		selectedAccountBalanceLabel: common.theme.H6(""),
		addressCopiedLabel:          common.theme.Caption(""),
	}

	return func() {
		page.Layout(common)
		page.Handle(common)
	}
}

func (p *receivePage) Layout(common pageCommon) {
	body := func() {
		layout.Stack{}.Layout(p.gtx,
			layout.Expanded(func() {
				layout.Inset{Top: unit.Dp(15)}.Layout(p.gtx, func() {
					layout.Flex{}.Layout(p.gtx,
						layout.Flexed(0.9, func() {
							p.ReceivePageContents(common)
						}),
						layout.Rigid(func() {
							p.rightNav()
						}),
					)
				})
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
			if p.isNewAddr {
				p.generateNewAddress()
			}
		}),
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(p.gtx, func() {
				p.infoBtn.Layout(p.gtx, &p.infoBtnW)
			})
		}),
		layout.Rigid(func() {
			if p.isInfo {
				p.infoDiag()
			}
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
		layout.Rigid(func() {
			inset := layout.Inset{
				Top:    unit.Dp(5),
				Bottom: unit.Dp(5),
			}
			inset.Layout(p.gtx, func() {
				p.newAddrBtn.TextSize = syncButtonTextSize
				p.newAddrBtn.Layout(p.gtx, &p.newAddrBtnW)
			})
		}),
	)
}

func (p *receivePage) infoDiag() {
	infoDetails := func() {
		layout.UniformInset(unit.Dp(5)).Layout(p.gtx, func() {
			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(p.gtx,
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(10)).Layout(p.gtx, func() {
						p.pageInfo.Layout(p.gtx)
					})
				}),
				layout.Rigid(func() {
					p.gotItBtn.TextSize = syncButtonTextSize
					p.gotItBtn.Layout(p.gtx, &p.gotItBtnW)
				}),
			)
		})
	}
	decredmaterial.Card{}.Layout(p.gtx, infoDetails)
}

func (p *receivePage) Handle(common pageCommon) {
	if p.infoBtnW.Clicked(p.gtx) {
		p.isInfo = !p.isInfo
		if p.isNewAddr {
			p.isNewAddr = false
		}
	}

	if p.moreBtnW.Clicked(p.gtx) {
		p.isNewAddr = !p.isNewAddr
		if p.isInfo {
			p.isInfo = false
		}
	}

	if p.gotItBtnW.Clicked(p.gtx) {
		p.isInfo = false
	}

	if p.newAddrBtnW.Clicked(p.gtx) {
		wallet := common.info.Wallets[*common.selectedWallet]
		account := common.info.Wallets[*common.selectedWallet].Accounts[0]

		addr, err := common.wallet.NextAddress(wallet.ID, account.Number)
		if err != nil {
			log.Debug("Error generating new address" + err.Error())
			// win.err = err.Error()
		} else {
			common.info.Wallets[*common.selectedWallet].Accounts[0].CurrentAddress = addr
			p.isNewAddr = false
		}
	}

	if p.copyBtnW.Clicked(p.gtx) {
		clipboard.WriteAll(common.info.Wallets[*common.selectedWallet].Accounts[0].CurrentAddress)
		p.addressCopiedLabel.Text = "Address Copied"
		time.AfterFunc(time.Second*3, func() {
			p.addressCopiedLabel.Text = ""
		})
		return
	}
}
