package ui

import (
	"bytes"
	"image"
	"image/color"
	"path/filepath"
	"time"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	qrcode "github.com/yeqown/go-qrcode"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const PageReceive = "Receive"

type receivePage struct {
	*pageCommon
	pageContainer     layout.List
	isNewAddr, isInfo bool
	currentAddress    string
	qrImage           *image.Image
	newAddr, copy     decredmaterial.Button
	info, more        decredmaterial.IconButton
	card              decredmaterial.Card
	receiveAddress    decredmaterial.Label
	gtx               *layout.Context

	selector *accountSelector

	backdrop   *widget.Clickable
	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton
}

func ReceivePage(common *pageCommon) Page {
	page := &receivePage{
		pageCommon: common,
		pageContainer: layout.List{
			Axis: layout.Vertical,
		},
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

	page.backButton, page.infoButton = common.SubPageHeaderButtons()
	page.backButton.Icon = page.icons.contentClear

	page.selector = newAccountSelector(common).
		title("Receiving account").
		accountSelected(func(selectedAccount *dcrlibwallet.Account) {
			selectedWallet := page.multiWallet.WalletWithID(selectedAccount.WalletID)
			currentAddress, err := selectedWallet.CurrentAddress(selectedAccount.Number)
			if err != nil {
				log.Errorf("Error getting current address: %v", err)
			} else {
				page.currentAddress = currentAddress
			}

			page.generateQRForAddress()
		}).
		accountValidator(func(account *dcrlibwallet.Account) bool {

			// Filter out imported account and mixed.
			wal := page.multiWallet.WalletWithID(account.WalletID)
			if account.Number == MaxInt32 ||
				account.Number == wal.MixedAccountNumber() {
				return false
			}
			return true
		})

	return page
}

func (pg *receivePage) OnResume() {
	pg.selector.selectFirstWalletValidAccount()
}

func (pg *receivePage) generateQRForAddress() {
	absoluteWdPath, err := GetAbsolutePath()
	if err != nil {
		log.Error(err.Error())
	}

	opt := qrcode.WithLogoImageFilePNG(filepath.Join(absoluteWdPath, "ui/assets/decredicons/qrcodeSymbol.png"))
	qrCode, err := qrcode.New(pg.currentAddress, opt)
	if err != nil {
		log.Error("Error generating address qrCode: " + err.Error())
		return
	}

	var buff bytes.Buffer
	err = qrCode.SaveTo(&buff)
	if err != nil {
		log.Error(err.Error())
		return
	}

	imgdec, _, err := image.Decode(bytes.NewReader(buff.Bytes()))
	if err != nil {
		log.Error(err.Error())
		return
	}

	pg.qrImage = &imgdec
}

func (pg *receivePage) Layout(gtx layout.Context) layout.Dimensions {
	if pg.gtx == nil {
		pg.gtx = &gtx
	}
	pg.pageBackdropLayout(gtx)

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.pageSections(gtx, func(gtx C) D {
				return pg.selector.Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.theme.Separator().Layout(gtx)
		},
		func(gtx C) D {
			return pg.pageSections(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.titleLayout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							return layout.Flex{
								Axis:      layout.Vertical,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									if pg.currentAddress != "" {
										return pg.addressLayout(gtx)
									}
									return layout.Dimensions{}
								}),
								layout.Rigid(func(gtx C) D {
									if pg.qrImage == nil {
										return layout.Dimensions{}
									}

									return pg.theme.ImageIcon(gtx, *pg.qrImage, 360)
								}),
							)
						})
					}),
				)
			})
		},
	}

	dims := pg.UniformPadding(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
					return pg.topNav(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.theme.Card().Layout(gtx, func(gtx C) D {
					return pg.pageContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
						return pageContent[i](gtx)
					})
				})
			}),
		)
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

func (pg *receivePage) topNav(gtx layout.Context) layout.Dimensions {
	m := values.MarginPadding20
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.backButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: m}.Layout(gtx, pg.theme.H6("Receive DCR").Layout)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, pg.infoButton.Layout)
		}),
	)
}

func (pg *receivePage) titleLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := pg.theme.Body2("Your Address")
			txt.Color = pg.theme.Color.Gray
			return txt.Layout(gtx)
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
				layout.Rigid(pg.more.Layout),
			)
		}),
	)
}

func (pg *receivePage) addressLayout(gtx layout.Context) layout.Dimensions {
	card := decredmaterial.Card{
		Inset: layout.Inset{
			Top:    values.MarginPadding14,
			Bottom: values.MarginPadding16,
		},
		Color: pg.theme.Color.LightGray,
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			card.Radius = decredmaterial.CornerRadius{NE: 8, NW: 0, SE: 0, SW: 8}
			return card.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top:    values.MarginPadding30,
					Bottom: values.MarginPadding30,
					Left:   values.MarginPadding30,
					Right:  values.MarginPadding30,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Dimensions{}
				})
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			pg.receiveAddress.Text = pg.currentAddress
			pg.receiveAddress.Color = pg.theme.Color.DeepBlue
			pg.receiveAddress.Alignment = text.Middle
			pg.receiveAddress.MaxLines = 1
			return card.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.UniformInset(values.MarginPadding16).Layout(gtx, pg.receiveAddress.Layout)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding1}.Layout(gtx, func(gtx C) D { return layout.Dimensions{} })
		}),
		layout.Rigid(func(gtx C) D {
			card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 8, SE: 8, SW: 0}
			return card.Layout(gtx, pg.copy.Layout)
		}),
	)
}

func (pg *receivePage) handle() {
	common := pg.pageCommon
	gtx := pg.gtx
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
		newAddr, err := pg.generateNewAddress()
		if err != nil {
			log.Debug("Error generating new address" + err.Error())
			return
		}

		pg.currentAddress = newAddr
		pg.generateQRForAddress()
		pg.isNewAddr = false
	}

	if pg.infoButton.Button.Clicked() {
		info := newInfoModal(common).
			title("Receive DCR").
			body("Each time you receive a payment, a new address is generated to protect your privacy.").
			positiveButton("Got it", func() {})
		common.showModal(info)
	}

	if pg.backButton.Button.Clicked() {
		common.changePage(*common.returnPage)
	}

	if pg.copy.Button.Clicked() {

		clipboard.WriteOp{Text: pg.currentAddress}.Add(gtx.Ops)

		pg.copy.Text = "Copied!"
		pg.copy.Color = common.theme.Color.Success
		time.AfterFunc(time.Second*3, func() {
			pg.copy.Text = "Copy"
			pg.copy.Color = common.theme.Color.Primary
		})
		return
	}
}
func (pg *receivePage) generateNewAddress() (string, error) {
	selectedWallet := pg.multiWallet.WalletWithID(pg.selector.selectedAccount.WalletID)

generateAddress:
	newAddr, err := selectedWallet.NextAddress(pg.selector.selectedAccount.Number)
	if err != nil {
		return "", err
	}

	if newAddr == pg.currentAddress {
		goto generateAddress
	}

	return newAddr, nil
}

func (pg *receivePage) onClose() {
}
