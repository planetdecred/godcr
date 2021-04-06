package ui

import (
	"bytes"
	"image"
	"image/color"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	qrcode "github.com/yeqown/go-qrcode"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/draw"
)

const PageReceive = "Receive"

type walletAccount struct {
	evt          *gesture.Click
	walletIndex  int
	accountIndex int
	accountName  string
	totalBalance string
	spendable    string
}

type walletAccountWidget struct {
	title                    decredmaterial.Label
	walletAccount            decredmaterial.Modal
	wallets, accounts        layout.List
	isWalletAccountModalOpen bool
	walletAccounts           map[int][]walletAccount
	fromAccount              *widget.Clickable
}
type receivePage struct {
	pageContainer     layout.List
	theme             *decredmaterial.Theme
	isNewAddr, isInfo bool
	addrs             string
	newAddr, copy     decredmaterial.Button
	info, more        decredmaterial.IconButton
	card              decredmaterial.Card
	receiveAddress    decredmaterial.Label

	line           *decredmaterial.Line
	backdrop       *widget.Clickable
	wallAcctWidget walletAccountWidget
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
		line:           common.theme.Line(),
		backdrop:       new(widget.Clickable),

		wallAcctWidget: walletAccountWidget{
			title:                    common.theme.Label(values.TextSize24, "Receiving account"),
			fromAccount:              new(widget.Clickable),
			walletAccount:            *common.theme.Modal(),
			wallets:                  layout.List{Axis: layout.Vertical},
			accounts:                 layout.List{Axis: layout.Vertical},
			walletAccounts:           make(map[int][]walletAccount),
			isWalletAccountModalOpen: false,
		},
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
	page.more.Color = common.theme.Color.IconColor
	page.more.Inset = layout.UniformInset(values.MarginPadding0)
	page.line.Color = common.theme.Color.Background
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
				return pg.wallAcctWidget.accountSelectLayout(gtx, common)
			})
		},
		func(gtx C) D {
			pg.line.Width, pg.line.Height = gtx.Constraints.Max.X, 1
			pg.line.Color = common.theme.Color.Background
			return pg.line.Layout(gtx)
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

	if pg.wallAcctWidget.isWalletAccountModalOpen {
		return common.Modal(gtx, dims, pg.wallAcctWidget.walletAccountModalLayout(gtx, common))
	}

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
					txt := common.theme.H6("Receive DCR")
					txt.Color = common.theme.Color.DeepBlue
					return layout.Inset{Left: m}.Layout(gtx, txt.Layout)
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
		Color: c.theme.Color.Background,
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
	pg.wallAcctWidget.Handler(common)

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

func (wg *walletAccountWidget) accountSelectLayout(gtx layout.Context, common pageCommon) layout.Dimensions {
	border := widget.Border{
		Color:        common.theme.Color.BorderColor,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}
	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return decredmaterial.Clickable(gtx, wg.fromAccount, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						accountIcon := common.icons.accountIcon
						accountIcon.Scale = 1
						inset := layout.Inset{
							Right: values.MarginPadding8,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return accountIcon.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return common.theme.Body1(
							common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].Name).Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						inset := layout.Inset{
							Left: values.MarginPadding4,
							Top:  values.MarginPadding2,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return decredmaterial.Card{
								Color: common.theme.Color.LightGray,
							}.Layout(gtx, func(gtx C) D {
								m2 := values.MarginPadding2
								m4 := values.MarginPadding4
								inset := layout.Inset{
									Left:   m4,
									Top:    m2,
									Bottom: m2,
									Right:  m4,
								}
								return inset.Layout(gtx, func(gtx C) D {
									text := common.theme.Caption(common.info.Wallets[*common.selectedWallet].Name)
									text.Color = common.theme.Color.Gray
									return text.Layout(gtx)
								})
							})
						})
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									txt := common.theme.Body1(
										common.info.Wallets[*common.selectedWallet].Accounts[*common.selectedAccount].TotalBalance)
									txt.Color = common.theme.Color.DeepBlue
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									inset := layout.Inset{
										Left: values.MarginPadding15,
									}
									return inset.Layout(gtx, func(gtx C) D {
										return common.icons.dropDownIcon.Layout(gtx, values.MarginPadding20)
									})
								}),
							)
						})
					}),
				)
			})
		})
	})
}

func (wg *walletAccountWidget) walletAccountModalLayout(gtx layout.Context, c pageCommon) layout.Dimensions {
	wallAcctGroup := func(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
		return layout.Inset{
			Bottom: values.MarginPadding10,
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := c.theme.Body2(title)
					txt.Color = c.theme.Color.Text
					inset := layout.Inset{
						Bottom: values.MarginPadding15,
					}
					return inset.Layout(gtx, txt.Layout)
				}),
				layout.Rigid(body),
			)
		})
	}

	w := []func(gtx C) D{
		func(gtx C) D {
			wg.title.Color = c.theme.Color.Text
			return wg.title.Layout(gtx)
		},
		func(gtx C) D {
			return wg.wallets.Layout(gtx, len(c.info.Wallets), func(gtx C, windex int) D {
				return wallAcctGroup(gtx, c.info.Wallets[windex].Name, func(gtx C) D {
					return wg.accounts.Layout(gtx, len(c.info.Wallets[windex].Accounts), func(gtx C, aindex int) D {
						click := wg.walletAccounts[windex][aindex].evt
						pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
						click.Add(gtx.Ops)
						wg.walletAccountsHandler(gtx, c, wg.walletAccounts[windex][aindex])
						return wg.walletAccountLayout(gtx, c, wg.walletAccounts[windex][aindex])
					})
				})
			})
		},
	}

	return wg.walletAccount.Layout(gtx, w, 850)
}

func (wg *walletAccountWidget) Handler(c pageCommon) {
	for windex := 0; windex < c.info.LoadedWallets; windex++ {
		if _, ok := wg.walletAccounts[windex]; !ok {
			accounts := c.info.Wallets[windex].Accounts
			if len(accounts) != len(wg.walletAccounts[windex]) {
				wg.walletAccounts[windex] = make([]walletAccount, len(accounts))
				for aindex := range accounts {
					wg.walletAccounts[windex][aindex] = walletAccount{
						walletIndex:  windex,
						accountIndex: aindex,
						evt:          &gesture.Click{},
						accountName:  accounts[aindex].Name,
						totalBalance: accounts[aindex].TotalBalance,
						spendable:    dcrutil.Amount(accounts[aindex].SpendableBalance).String(),
					}
				}
			}
		}
	}

	if wg.fromAccount.Clicked() {
		wg.isWalletAccountModalOpen = true
	}
}

func (wg *walletAccountWidget) walletAccountsHandler(gtx layout.Context, common pageCommon, wallAcct walletAccount) {
	for _, e := range wallAcct.evt.Events(gtx) {
		if e.Type == gesture.TypeClick {
			*common.selectedWallet = wallAcct.walletIndex
			*common.selectedAccount = wallAcct.accountIndex
			wg.isWalletAccountModalOpen = false
		}
	}
}

func (wg *walletAccountWidget) walletAccountLayout(gtx layout.Context, common pageCommon, wallAcct walletAccount) layout.Dimensions {
	accountIcon := common.icons.accountIcon
	accountIcon.Scale = 1

	inset := layout.Inset{
		Bottom: values.MarginPadding10,
	}
	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(0.1, func(gtx C) D {
						inset := layout.Inset{
							Right: values.MarginPadding10,
							Top:   values.MarginPadding15,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return accountIcon.Layout(gtx)
						})
					}),
					layout.Flexed(0.8, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								accountLabel := common.theme.Body2(wallAcct.accountName)
								accountLabel.Color = common.theme.Color.Text
								accountBalLabel := common.theme.Body2(wallAcct.totalBalance)
								accountBalLabel.Color = common.theme.Color.Text
								return wg.accountTableLayout(gtx, accountLabel, accountBalLabel)
							}),
							layout.Rigid(func(gtx C) D {
								spendibleLabel := common.theme.Body2("Spendable")
								spendibleLabel.Color = common.theme.Color.Gray
								spendibleBalLabel := common.theme.Body2(wallAcct.spendable)
								spendibleBalLabel.Color = common.theme.Color.Gray
								return wg.accountTableLayout(gtx, spendibleLabel, spendibleBalLabel)
							}),
						)
					}),
					layout.Flexed(0.1, func(gtx C) D {
						inset := layout.Inset{
							Right: values.MarginPadding10,
							Top:   values.MarginPadding10,
						}

						if *common.selectedWallet == wallAcct.walletIndex && *common.selectedAccount == wallAcct.accountIndex {
							return layout.E.Layout(gtx, func(gtx C) D {
								return inset.Layout(gtx, func(gtx C) D {
									return common.icons.navigationCheck.Layout(gtx, values.MarginPadding20)
								})
							})
						}
						return layout.Dimensions{}
					}),
				)
			}),
		)
	})
}

func (wg *walletAccountWidget) accountTableLayout(gtx layout.Context, leftLabel, rightLabel decredmaterial.Label) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Top: values.MarginPadding2,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return leftLabel.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return rightLabel.Layout(gtx)
			})
		}),
	)
}
