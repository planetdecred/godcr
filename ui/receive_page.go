package ui

import (
	"bytes"
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	qrcode "github.com/yeqown/go-qrcode"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/draw"
)

const PageReceive = "Receive"

type receivePage struct {
	pageContainer     layout.List
	theme             *decredmaterial.Theme
	isNewAddr, isInfo bool
	addrs             string
	newAddr, copy     decredmaterial.Button
	info, more        decredmaterial.IconButton
	card              decredmaterial.Card
	receiveAddress    decredmaterial.Label

	backdrop *widget.Clickable
}

func (win *Window) ReceivePage(common pageCommon) layout.Widget {
	page := &receivePage{
		pageContainer: layout.List{
			Axis: layout.Vertical,
		},
		theme:          common.theme,
		info:           common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.ActionInfo))),
		copy:           common.theme.Button(new(widget.Clickable), "Copy"),
		more:           common.theme.PlainIconButton(new(widget.Clickable), common.icons.navMoreIcon),
		newAddr:        common.theme.Button(new(widget.Clickable), "Generate new address"),
		receiveAddress: common.theme.Label(values.TextSize20, ""),
		card:           common.theme.Card(),
		backdrop:       new(widget.Clickable),
	}

	page.info.Inset, page.info.Size = layout.UniformInset(values.MarginPadding5), values.MarginPadding20
	page.copy.Background = color.NRGBA{}
	page.copy.Color = common.theme.Color.Primary
	page.copy.Inset = layout.Inset{
		Top:    values.MarginPadding19p5,
		Bottom: values.MarginPadding19p5,
		Left:   values.MarginPadding16,
		Right:  values.MarginPadding16,
	}
	page.more.Color = common.theme.Color.Gray3
	page.more.Inset = layout.UniformInset(values.MarginPadding0)
	page.newAddr.Inset = layout.Inset{
		Top:    values.MarginPadding20,
		Bottom: values.MarginPadding20,
		Left:   values.MarginPadding16,
		Right:  values.MarginPadding16,
	}
	page.newAddr.Color = common.theme.Color.Text
	page.newAddr.Background = common.theme.Color.Surface
	page.newAddr.TextSize = values.TextSize16

	return func(gtx C) D {
		page.Handle(common)
		return page.Layout(gtx, common)
	}
}

func (pg *receivePage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.pageBackdropLayout(gtx)

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.pageSections(gtx, func(gtx C) D {
				return common.accountSelectorLayout(gtx, "Receiving account")
			})
		},
		func(gtx C) D {
			return pg.theme.Separator().Layout(gtx)
		},
		func(gtx C) D {
			return pg.pageSections(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.titleLayout(gtx, common)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							return layout.Flex{
								Axis:      layout.Vertical,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									if pg.addrs != "" {
										return pg.addressLayout(gtx, common)
									}
									return layout.Dimensions{}
								}),
								layout.Rigid(func(gtx C) D {
									return pg.addressQRCodeLayout(gtx, common)
								}),
							)
						})
					}),
				)
			})
		},
	}

	dims := common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return pg.topNav(gtx, common)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return common.theme.Card().Layout(gtx, func(gtx C) D {
						return pg.pageContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
							return pageContent[i](gtx)
						})
					})
				}),
			)
		})
	})

	return dims
}

func (pg *receivePage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return pg.theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding16).Layout(gtx, body)
	})
}

// pageBackdropLayout layout of background overlay when the popup button generate new address is show,
// click outside of the generate new address button to hide the button
func (pg *receivePage) pageBackdropLayout(gtx layout.Context) {
	if pg.isNewAddr {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
		m := op.Record(gtx.Ops)
		pg.backdrop.Layout(gtx)
		op.Defer(gtx.Ops, m.Stop())
	}
}

func (pg *receivePage) topNav(gtx layout.Context, common pageCommon) layout.Dimensions {
	m := values.MarginPadding20
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					common.subPageBackButton.Icon = common.icons.contentClear
					return common.subPageBackButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: m}.Layout(gtx, func(gtx C) D {
						return pg.theme.H6("Receive DCR").Layout(gtx)
					})
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return common.subPageInfoButton.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *receivePage) titleLayout(gtx layout.Context, common pageCommon) layout.Dimensions {
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return common.theme.Body1("Your Address").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if pg.isNewAddr {
						m := op.Record(gtx.Ops)
						layout.Inset{Top: values.MarginPadding30, Left: unit.Dp(-152)}.Layout(gtx, func(gtx C) D {
							return pg.newAddr.Layout(gtx)
						})
						op.Defer(gtx.Ops, m.Stop())
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return pg.more.Layout(gtx)
				}),
			)
		}),
	)
}

func (pg *receivePage) addressLayout(gtx layout.Context, c pageCommon) layout.Dimensions {
	card := decredmaterial.Card{
		Inset: layout.Inset{
			Top:    values.MarginPadding14,
			Bottom: values.MarginPadding16,
		},
		Color: c.theme.Color.LightGray,
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			pg.receiveAddress.Text = pg.addrs
			pg.receiveAddress.Alignment = text.Middle
			pg.receiveAddress.MaxLines = 1
			card.Radius.NE = 8
			card.Radius.SW = 8
			card.Radius.NW = 0
			card.Radius.SE = 0
			return card.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
					return pg.receiveAddress.Layout(gtx)
				})
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding1}.Layout(gtx, func(gtx C) D { return layout.Dimensions{} })
		}),
		layout.Rigid(func(gtx C) D {
			card.Radius.NE = 0
			card.Radius.SW = 0
			card.Radius.NW = 8
			card.Radius.SE = 8
			return card.Layout(gtx, func(gtx C) D {
				return pg.copy.Layout(gtx)
			})
		}),
	)
}

func (pg *receivePage) addressQRCodeLayout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.addrs = common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress
	opt := qrcode.WithLogoImageFilePNG("ui/assets/decredicons/qrcodeSymbol.png")
	qrCode, err := qrcode.New(pg.addrs, opt)
	if err != nil {
		log.Error("Error generating address qrCode: " + err.Error())
		return layout.Dimensions{}
	}

	var buff bytes.Buffer
	err = qrCode.SaveTo(&buff)
	if err != nil {
		log.Error(err.Error())
		return layout.Dimensions{}
	}
	imgdec, _, err := image.Decode(bytes.NewReader(buff.Bytes()))
	if err != nil {
		log.Error(err.Error())
		return layout.Dimensions{}
	}

	imgs := image.NewRGBA(image.Rectangle{Max: image.Point{X: 180, Y: 180}})
	draw.ApproxBiLinear.Scale(imgs, imgs.Bounds(), imgdec, imgdec.Bounds(), draw.Src, nil)

	src := paint.NewImageOp(imgs)
	img := widget.Image{
		Src:   src,
		Scale: 1,
	}

	return img.Layout(gtx)
}

func (pg *receivePage) Handle(common pageCommon) {
	if pg.backdrop.Clicked() {
		pg.isNewAddr = false
	}

	if pg.more.Button.Clicked() {
		pg.isNewAddr = !pg.isNewAddr
		if pg.isInfo {
			pg.isInfo = false
		}
	}

	if pg.newAddr.Button.Clicked() {
		wall := common.info.Wallets[*common.selectedWallet]

		var generateNewAddress func(wall wallet.InfoShort)
		generateNewAddress = func(wall wallet.InfoShort) {
			oldAddr := wall.Accounts[*common.selectedAccount].CurrentAddress
			newAddr, err := common.wallet.NextAddress(wall.ID, wall.Accounts[*common.selectedAccount].Number)
			if err != nil {
				log.Debug("Error generating new address" + err.Error())
				return
			}
			if newAddr == oldAddr {
				log.Info("Call again to generate new address")
				generateNewAddress(wall)
				return
			}
			common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress = newAddr
			pg.isNewAddr = false
		}
		generateNewAddress(wall)
	}

	if common.subPageInfoButton.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template:   ReceiveInfoTemplate,
				title:      "Receive DCR",
				cancel:     common.closeModal,
				cancelText: "Got it",
			}
		}()
	}

	if common.subPageBackButton.Button.Clicked() {
		*common.page = PageOverview
	}

	if pg.copy.Button.Clicked() {
		go func() {
			common.clipboard <- WriteClipboard{Text: common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].CurrentAddress}
		}()
		pg.copy.Text = "Copied!"
		pg.copy.Color = common.theme.Color.Success
		time.AfterFunc(time.Second*3, func() {
			pg.copy.Text = "Copy"
			pg.copy.Color = common.theme.Color.Primary
		})
		return
	}
}
