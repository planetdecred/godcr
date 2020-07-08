package ui

import (
	"time"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/atotto/clipboard"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/skip2/go-qrcode"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const PageReceive = "receive"

type receivePage struct {
	pageContainer layout.List
	gtx           *layout.Context

	isNewAddr, isInfo bool
	addrs             string

	newAddrBtn, minInfo                                 decredmaterial.Button
	copyBtn, infoBtn, moreBtn                           decredmaterial.IconButton
	copyBtnW, infoBtnW, moreBtnW, minInfoW, newAddrBtnW widget.Button

	selectedAccountNameLabel, selectedAccountBalanceLabel decredmaterial.Label
	receiveAddressLabel, addressCopiedLabel, pageInfo     decredmaterial.Label
	selectedWalletBalLabel, selectedWalletNameLabel       decredmaterial.Label
}

func (win *Window) ReceivePage(common pageCommon) layout.Widget {
	moreBtn := common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.NavigationMoreVert)))
	moreBtn.Padding, moreBtn.Size = values.MarginPadding5, values.MarginPadding35
	infoBtn := common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ActionInfo)))
	infoBtn.Padding, infoBtn.Size = values.MarginPadding5, values.MarginPadding35
	copyBtn := common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy)))
	copyBtn.Padding, copyBtn.Size = values.MarginPadding5, values.MarginPadding35
	copyBtn.Background = common.theme.Color.Background
	copyBtn.Color = common.theme.Color.Text
	receiveAddressLabel := common.theme.H6("")
	receiveAddressLabel.Color = common.theme.Color.Primary
	pageInfo := common.theme.Body1("Each time you request a payment, a \nnew address is created to protect \nyour privacy.")
	page := &receivePage{
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
		gtx:                         common.gtx,
		moreBtn:                     moreBtn,
		infoBtn:                     infoBtn,
		copyBtn:                     copyBtn,
		minInfo:                     common.theme.Button("Got It"),
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
		layout.Stack{Alignment: layout.NE}.Layout(p.gtx,
			layout.Expanded(func() {
				layout.Inset{Top: values.MarginPadding15}.Layout(p.gtx, func() {
					layout.Flex{}.Layout(p.gtx,
						layout.Flexed(1, func() {
							p.ReceivePageContents(common)
						}),
					)
				})
			}),
			layout.Stacked(func() {
				layout.Inset{Right: values.MarginPadding30}.Layout(p.gtx, func() {
					p.rightNav()
				})
			}),
		)
	}
	common.LayoutWithAccounts(p.gtx, body)
}

func (p *receivePage) ReceivePageContents(common pageCommon) {
	layout.Center.Layout(p.gtx, func() {
		pageContent := []func(){
			func() {
				p.selectedAccountColumn(common)
			},
			func() {
				p.qrCodeAddressColumn(common)
			},
			func() {
				if p.addrs != "" {
					p.receiveAddressColumn()
				}
			},
			func() {
				layout.Flex{}.Layout(p.gtx,
					layout.Rigid(func() {
						if p.addressCopiedLabel.Text != "" {
							p.addressCopiedLabel.Layout(p.gtx)
						}
					}),
				)
			},
		}
		p.pageContainer.Layout(p.gtx, len(pageContent), func(i int) {
			layout.Inset{}.Layout(p.gtx, pageContent[i])
		})
	})
}

func (p *receivePage) rightNav() {
	layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(p.gtx,
		layout.Rigid(func() {
			p.moreBtn.Layout(p.gtx, &p.moreBtnW)
		}),
		layout.Rigid(func() {
			if p.isNewAddr {
				p.generateNewAddress()
			}
		}),
		layout.Rigid(func() {
			// layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(p.gtx, func() {
			// 	p.infoBtn.Layout(p.gtx, &p.infoBtnW)
			// })
		}),
		layout.Rigid(func() {
			if p.isInfo {
				p.infoDiag()
			}
		}),
	)
}

func (p *receivePage) selectedAccountColumn(common pageCommon) {
	current := common.info.Wallets[*common.selectedWallet]

	p.selectedWalletNameLabel.Text = current.Name
	p.selectedWalletBalLabel.Text = current.Balance

	account := common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount]
	p.selectedAccountNameLabel.Text = account.Name
	p.selectedAccountBalanceLabel.Text = dcrutil.Amount(account.SpendableBalance).String()

	selectedDetails := func() {
		layout.UniformInset(values.MarginPadding10).Layout(p.gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(p.gtx,
				layout.Rigid(func() {
					layout.Flex{}.Layout(p.gtx,
						layout.Rigid(func() {
							layout.Inset{Bottom: values.MarginPadding5}.Layout(p.gtx, func() {
								p.selectedAccountNameLabel.Layout(p.gtx)
							})
						}),
						layout.Rigid(func() {
							layout.Inset{Left: values.MarginPadding20}.Layout(p.gtx, func() {
								p.selectedAccountBalanceLabel.Layout(p.gtx)
							})
						}),
					)
				}),
				layout.Rigid(func() {
					layout.Flex{}.Layout(p.gtx,
						layout.Rigid(func() {
							layout.Inset{Bottom: values.MarginPadding5}.Layout(p.gtx, func() {
								p.selectedWalletNameLabel.Layout(p.gtx)
							})
						}),
						layout.Rigid(func() {
							layout.Inset{Left: values.MarginPadding20}.Layout(p.gtx, func() {
								p.selectedWalletBalLabel.Layout(p.gtx)
							})
						}),
					)
				}),
			)
		})
	}
	decredmaterial.Card{Color: common.theme.Color.Surface}.Layout(p.gtx, selectedDetails)
}

func (p *receivePage) qrCodeAddressColumn(common pageCommon) {
	p.addrs = common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress
	qrCode, err := qrcode.New(p.addrs, qrcode.Highest)
	if err != nil {
		log.Error("Error generating address qrCode: " + err.Error())
		return
	}

	qrCode.DisableBorder = true
	layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding10}.Layout(p.gtx, func() {
		img := common.theme.Image(paint.NewImageOp(qrCode.Image(520)))
		img.Src.Rect.Max.X = 521
		img.Src.Rect.Max.Y = 521
		img.Scale = 0.5
		img.Layout(p.gtx)
	})
}

func (p *receivePage) receiveAddressColumn() {
	layout.Flex{}.Layout(p.gtx,
		layout.Rigid(func() {
			p.receiveAddressLabel.Text = p.addrs
			p.receiveAddressLabel.Layout(p.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Left: values.MarginPadding10}.Layout(p.gtx, func() {
				p.copyBtn.Layout(p.gtx, &p.copyBtnW)
			})
		}),
	)
}

func (p *receivePage) generateNewAddress() {
	layout.Flex{}.Layout(p.gtx,
		layout.Rigid(func() {
			inset := layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
			}
			inset.Layout(p.gtx, func() {
				p.newAddrBtn.TextSize = values.TextSize10
				p.newAddrBtn.Layout(p.gtx, &p.newAddrBtnW)
			})
		}),
	)
}

func (p *receivePage) infoDiag() {
	infoDetails := func() {
		layout.UniformInset(values.MarginPadding10).Layout(p.gtx, func() {
			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(p.gtx,
				layout.Rigid(func() {
					layout.Inset{Bottom: values.MarginPadding5}.Layout(p.gtx, func() {
						p.pageInfo.Layout(p.gtx)
					})
				}),
				layout.Rigid(func() {
					p.minInfo.TextSize = values.TextSize10
					p.minInfo.Layout(p.gtx, &p.minInfoW)
				}),
			)
		})
	}
	decredmaterial.Card{}.Layout(p.gtx, infoDetails)
}

func (p *receivePage) Handle(common pageCommon) {
	// if p.infoBtnW.Clicked(p.gtx) {
	// 	p.isInfo = !p.isInfo
	// 	if p.isNewAddr {
	// 		p.isNewAddr = false
	// 	}
	// }

	if p.moreBtnW.Clicked(p.gtx) {
		p.isNewAddr = !p.isNewAddr
		if p.isInfo {
			p.isInfo = false
		}
	}

	if p.minInfoW.Clicked(p.gtx) {
		p.isInfo = false
	}

	if p.newAddrBtnW.Clicked(p.gtx) {
		wallet := common.info.Wallets[*common.selectedWallet]
		account := common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount]

		addr, err := common.wallet.NextAddress(wallet.ID, account.Number)
		if err != nil {
			log.Debug("Error generating new address" + err.Error())
			// win.err = err.Error()
		} else {
			common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress = addr
			p.isNewAddr = false
		}
	}

	if p.copyBtnW.Clicked(p.gtx) {
		clipboard.WriteAll(common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress)
		p.addressCopiedLabel.Text = "Address Copied"
		time.AfterFunc(time.Second*3, func() {
			p.addressCopiedLabel.Text = ""
		})
		return
	}
}
