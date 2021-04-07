package ui

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/skip2/go-qrcode"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const PageReceive = "Receive"

type receivePage struct {
	pageContainer layout.List

	isNewAddr, isInfo bool
	addrs             string

	newAddrBtn, minInfo       decredmaterial.Button
	copyBtn, infoBtn, moreBtn decredmaterial.IconButton
	card                      decredmaterial.Card

	receiveAddressLabel, addressCopiedLabel, pageInfo decredmaterial.Label
}

func (win *Window) ReceivePage(common pageCommon) layout.Widget {
	moreBtn := common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.NavigationMoreVert)))
	moreBtn.Inset, moreBtn.Size = layout.UniformInset(values.MarginPadding5), values.MarginPadding20
	infoBtn := common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.ActionInfo)))
	infoBtn.Inset, infoBtn.Size = layout.UniformInset(values.MarginPadding5), values.MarginPadding20
	copyBtn := common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.ContentContentCopy)))
	copyBtn.Inset, copyBtn.Size = layout.UniformInset(values.MarginPadding5), values.MarginPadding20
	copyBtn.Background = common.theme.Color.LightGray
	copyBtn.Color = common.theme.Color.Text
	receiveAddressLabel := common.theme.H6("")
	receiveAddressLabel.Color = common.theme.Color.Primary
	pageInfo := common.theme.Body1("Each time you request a payment, a \nnew address is created to protect \nyour privacy.")
	page := &receivePage{
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
		moreBtn:             moreBtn,
		infoBtn:             infoBtn,
		copyBtn:             copyBtn,
		minInfo:             common.theme.Button(new(widget.Clickable), "Got It"),
		newAddrBtn:          common.theme.Button(new(widget.Clickable), "Generate new address"),
		receiveAddressLabel: receiveAddressLabel,
		pageInfo:            pageInfo,
		addressCopiedLabel:  common.theme.Caption(""),
		card:                common.theme.Card(),
	}

	return func(gtx C) D {
		page.Handle(common)
		return page.Layout(gtx, common)
	}
}

func (pg *receivePage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		return layout.Stack{Alignment: layout.NE}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return pg.ReceivePageContents(gtx, common)
						}),
					)
				})
			}),
			layout.Stacked(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
					return pg.rightNav(gtx)
				})
			}),
		)
	}
	return common.LayoutWithAccounts(gtx, body)
}

func (pg *receivePage) ReceivePageContents(gtx layout.Context, common pageCommon) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx C) D {
		pageContent := []func(gtx C) D{
			func(gtx C) D {
				return common.SelectedAccountLayout(gtx)
			},
			func(gtx C) D {
				return pg.qrCodeAddressColumn(gtx, common)
			},
			func(gtx C) D {
				if pg.addrs != "" {
					return pg.receiveAddressColumn(gtx)
				}
				return layout.Dimensions{}
			},
			func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if pg.addressCopiedLabel.Text != "" {
							return pg.addressCopiedLabel.Layout(gtx)
						}
						return layout.Dimensions{}
					}),
				)
			},
		}
		return pg.pageContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
			return layout.Inset{}.Layout(gtx, pageContent[i])
		})
	})
}

func (pg *receivePage) rightNav(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.moreBtn.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			if pg.isNewAddr {
				return pg.generateNewAddress(gtx)
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			if pg.isInfo {
				return pg.infoDiag(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

func (pg *receivePage) qrCodeAddressColumn(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.addrs = common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress
	qrCode, err := qrcode.New(pg.addrs, qrcode.Highest)
	if err != nil {
		log.Error("Error generating address qrCode: " + err.Error())
		return layout.Dimensions{}
	}

	qrCode.DisableBorder = true
	return layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		img := widget.Image{
			Src: paint.NewImageOp(qrCode.Image(100)),
		}
		return img.Layout(gtx)
	})
}

func (pg *receivePage) receiveAddressColumn(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			pg.receiveAddressLabel.Text = pg.addrs
			return pg.receiveAddressLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding10, Top: values.MarginPaddingMinus5}.Layout(gtx, func(gtx C) D {
				return pg.copyBtn.Layout(gtx)
			})
		}),
	)
}

func (pg *receivePage) generateNewAddress(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
			}
			return inset.Layout(gtx, func(gtx C) D {
				pg.newAddrBtn.TextSize = values.TextSize10
				return pg.newAddrBtn.Layout(gtx)
			})
		}),
	)
}

func (pg *receivePage) infoDiag(gtx layout.Context) layout.Dimensions {
	infoDetails := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return pg.pageInfo.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					pg.minInfo.TextSize = values.TextSize10
					return pg.minInfo.Layout(gtx)
				}),
			)
		})
	}

	return pg.card.Layout(gtx, infoDetails)
}

func (pg *receivePage) Handle(common pageCommon) {
	if pg.moreBtn.Button.Clicked() {
		pg.isNewAddr = !pg.isNewAddr
		if pg.isInfo {
			pg.isInfo = false
		}
	}

	if pg.minInfo.Button.Clicked() {
		pg.isInfo = false
	}

	if pg.newAddrBtn.Button.Clicked() {
		wallet := common.info.Wallets[*common.selectedWallet]
		account := common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount]

		addr, err := common.wallet.NextAddress(wallet.ID, account.Number)
		if err != nil {
			log.Debug("Error generating new address" + err.Error())
		} else {
			common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress = addr
			pg.isNewAddr = false
		}
	}

	if pg.copyBtn.Button.Clicked() {
		go func() {
			common.clipboard <- WriteClipboard{Text: common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress}
		}()
		pg.addressCopiedLabel.Text = "Address Copied"
		time.AfterFunc(time.Second*3, func() {
			pg.addressCopiedLabel.Text = ""
		})
		return
	}
}
