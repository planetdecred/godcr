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
	gtx           layout.Context

	isNewAddr, isInfo bool
	addrs             string

	newAddrBtn, minInfo       decredmaterial.Button
	copyBtn, infoBtn, moreBtn decredmaterial.IconButton
	// copyBtnW, infoBtnW, moreBtnW, minInfoW, newAddrBtnW widget.Clickable

	selectedAccountNameLabel, selectedAccountBalanceLabel decredmaterial.Label
	receiveAddressLabel, addressCopiedLabel, pageInfo     decredmaterial.Label
	selectedWalletBalLabel, selectedWalletNameLabel       decredmaterial.Label
}

func (win *Window) ReceivePage(common pageCommon) layout.Widget {
	moreBtn := common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.NavigationMoreVert)))
	moreBtn.Inset, moreBtn.Size = layout.UniformInset(values.MarginPadding5), values.MarginPadding35
	infoBtn := common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.ActionInfo)))
	infoBtn.Inset, infoBtn.Size = layout.UniformInset(values.MarginPadding5), values.MarginPadding35
	copyBtn := common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.ContentContentCopy)))
	copyBtn.Inset, copyBtn.Size = layout.UniformInset(values.MarginPadding5), values.MarginPadding35
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
		minInfo:                     common.theme.Button(new(widget.Clickable), "Got It"),
		newAddrBtn:                  common.theme.Button(new(widget.Clickable), "Generate new address"),
		receiveAddressLabel:         receiveAddressLabel,
		pageInfo:                    pageInfo,
		selectedAccountNameLabel:    common.theme.H6(""),
		selectedWalletNameLabel:     common.theme.Body2(""),
		selectedWalletBalLabel:      common.theme.Body2(""),
		selectedAccountBalanceLabel: common.theme.H6(""),
		addressCopiedLabel:          common.theme.Caption(""),
	}

	return func(gtx C) D {
		page.Handle(common)
		return page.Layout(common)
	}
}

func (p *receivePage) Layout(common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		return layout.Stack{Alignment: layout.NE}.Layout(p.gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding15}.Layout(p.gtx, func(gtx C) D {
					return layout.Flex{}.Layout(p.gtx,
						layout.Flexed(1, func(gtx C) D {
							return p.ReceivePageContents(common)
						}),
					)
				})
			}),
			layout.Stacked(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding30}.Layout(p.gtx, func(gtx C) D {
					return p.rightNav()
				})
			}),
		)
	}
	return common.LayoutWithAccounts(p.gtx, body)
}

func (p *receivePage) ReceivePageContents(common pageCommon) layout.Dimensions {
	dims := layout.Center.Layout(p.gtx, func(gtx C) D {
		pageContent := []func(gtx C) D{
			func(gtx C) D {
				return p.selectedAccountColumn(common)
			},
			func(gtx C) D {
				return p.qrCodeAddressColumn(common)
			},
			func(gtx C) D {
				if p.addrs != "" {
					return p.receiveAddressColumn()
				}
				return layout.Dimensions{}
			},
			func(gtx C) D {
				return layout.Flex{}.Layout(p.gtx,
					layout.Rigid(func(gtx C) D {
						if p.addressCopiedLabel.Text != "" {
							return p.addressCopiedLabel.Layout(p.gtx)
						}
						return layout.Dimensions{}
					}),
				)
			},
		}
		return p.pageContainer.Layout(p.gtx, len(pageContent), func(gtx C, i int) D {
			return layout.Inset{}.Layout(p.gtx, pageContent[i])
		})
	})
	return dims
}

func (p *receivePage) rightNav() layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(p.gtx,
		layout.Rigid(func(gtx C) D {
			return p.moreBtn.Layout(p.gtx)
		}),
		layout.Rigid(func(gtx C) D {
			if p.isNewAddr {
				return p.generateNewAddress()
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			// layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(p.gtx, func() {
			// 	p.infoBtn.Layout(p.gtx, &p.infoBtnW)
			// })
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			if p.isInfo {
				return p.infoDiag()
			}
			return layout.Dimensions{}
		}),
	)
}

func (p *receivePage) selectedAccountColumn(common pageCommon) layout.Dimensions {
	current := common.info.Wallets[*common.selectedWallet]

	p.selectedWalletNameLabel.Text = current.Name
	p.selectedWalletBalLabel.Text = current.Balance

	account := common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount]
	p.selectedAccountNameLabel.Text = account.Name
	p.selectedAccountBalanceLabel.Text = dcrutil.Amount(account.SpendableBalance).String()

	selectedDetails := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(p.gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(p.gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(p.gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(p.gtx, func(gtx C) D {
								return p.selectedAccountNameLabel.Layout(p.gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(p.gtx, func(gtx C) D {
								return p.selectedAccountBalanceLabel.Layout(p.gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(p.gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(p.gtx, func(gtx C) D {
								return p.selectedWalletNameLabel.Layout(p.gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(p.gtx, func(gtx C) D {
								return p.selectedWalletBalLabel.Layout(p.gtx)
							})
						}),
					)
				}),
			)
		})
	}
	return decredmaterial.Card{}.Layout(p.gtx, selectedDetails)
	decredmaterial.Card{Color: common.theme.Color.Surface}.Layout(gtx, selectedDetails)
}

func (p *receivePage) qrCodeAddressColumn(common pageCommon) layout.Dimensions {
	p.addrs = common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress
	qrCode, err := qrcode.New(p.addrs, qrcode.Highest)
	if err != nil {
		log.Error("Error generating address qrCode: " + err.Error())
		return layout.Dimensions{}
	}

	qrCode.DisableBorder = true
	return layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding10}.Layout(p.gtx, func(gtx C) D {
		img := common.theme.Image(paint.NewImageOp(qrCode.Image(520)))
		img.Src.Rect.Max.X = 521
		img.Src.Rect.Max.Y = 521
		img.Scale = 0.5
		return img.Layout(p.gtx)
	})
}

func (p *receivePage) receiveAddressColumn() layout.Dimensions {
	return layout.Flex{}.Layout(p.gtx,
		layout.Rigid(func(gtx C) D {
			p.receiveAddressLabel.Text = p.addrs
			return p.receiveAddressLabel.Layout(p.gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding10}.Layout(p.gtx, func(gtx C) D {
				return p.copyBtn.Layout(p.gtx)
			})
		}),
	)
}

func (p *receivePage) generateNewAddress() layout.Dimensions {
	return layout.Flex{}.Layout(p.gtx,
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
			}
			return inset.Layout(p.gtx, func(gtx C) D {
				p.newAddrBtn.TextSize = values.TextSize10
				return p.newAddrBtn.Layout(p.gtx)
			})
		}),
	)
}

func (p *receivePage) infoDiag() layout.Dimensions {
	infoDetails := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(p.gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(p.gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding5}.Layout(p.gtx, func(gtx C) D {
						return p.pageInfo.Layout(p.gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					p.minInfo.TextSize = values.TextSize10
					return p.minInfo.Layout(p.gtx)
				}),
			)
		})
	}
	return decredmaterial.Card{}.Layout(p.gtx, infoDetails)
}

func (p *receivePage) Handle(common pageCommon) {
	// if p.infoBtnW.Clicked(p.gtx) {
	// 	p.isInfo = !p.isInfo
	// 	if p.isNewAddr {
	// 		p.isNewAddr = false
	// 	}
	// }

	if p.moreBtn.Button.Clicked() {
		p.isNewAddr = !p.isNewAddr
		if p.isInfo {
			p.isInfo = false
		}
	}

	if p.minInfo.Button.Clicked() {
		p.isInfo = false
	}

	if p.newAddrBtn.Button.Clicked() {
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

	if p.copyBtn.Button.Clicked() {
		clipboard.WriteAll(common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress)
		p.addressCopiedLabel.Text = "Address Copied"
		time.AfterFunc(time.Second*3, func() {
			p.addressCopiedLabel.Text = ""
		})
		return
	}
}
