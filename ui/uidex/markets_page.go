package uidex

import (
	"fmt"
	"image"
	"io/ioutil"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/dex"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/sqweek/dialog"
)

const PageMarkets = "MarketsPage"

type navItem struct {
	evt        *gesture.Click
	imageLeft  *widget.Image
	imageRight *widget.Image
	asset      *selectedMaket
}

type drawerNav struct {
	host     string
	elements []navItem
}

// walletActionInfo the data will be show up in the unlock or create wallet modal
type walletActionInfo struct {
	image    *widget.Image
	coin     string
	coinName string
	coinID   uint32
}

type walletActionWidget struct {
	evt  *gesture.Click
	info *walletActionInfo
}

type marketsPage struct {
	theme     *decredmaterial.Theme
	pageModal *decredmaterial.Modal
	exchange  layout.List

	supportedAsset []*core.SupportedAsset
	user           **dex.User
	selectedMaket  **selectedMaket
	cert           string
	certName       string

	drawerNavItems []*drawerNav

	appPassword      decredmaterial.Editor
	appPasswordAgain decredmaterial.Editor
	accountName      decredmaterial.Editor
	walletPassword   decredmaterial.Editor
	dexServerAddress decredmaterial.Editor

	createPassword   decredmaterial.Button
	createNewWallet  decredmaterial.Button
	unlockWallet     decredmaterial.Button
	login            decredmaterial.Button
	addCertFile      decredmaterial.Button
	addDexServer     decredmaterial.Button
	register         decredmaterial.Button
	choseAssetUnlock decredmaterial.IconButton
	toWallet         decredmaterial.IconButton

	isLoggedIn          bool
	showAddWallet       bool
	showUnlockWallet    bool
	showConfirmRegister bool
	errCreateWalletChan chan error
	errLoginChan        chan error
	errInitappChan      chan error
	errRegisterChan     chan error
	errUnlockWallChan   chan error
	responseGetDex      *core.Exchange
	responseGetDexChan  chan *core.Exchange

	walletActionWidgets map[string]*walletActionWidget
	walletActionInfo    *walletActionInfo
}

func (d *Dex) MarketsPage(common pageCommon) layout.Widget {
	pg := &marketsPage{
		theme:         common.theme,
		pageModal:     common.theme.Modal(),
		exchange:      layout.List{Axis: layout.Vertical},
		user:          &d.userInfo,
		selectedMaket: &d.market,

		drawerNavItems:      make([]*drawerNav, 0),
		errCreateWalletChan: make(chan error),
		errInitappChan:      make(chan error),
		errLoginChan:        make(chan error),
		errRegisterChan:     make(chan error),
		errUnlockWallChan:   make(chan error),
		responseGetDexChan:  make(chan *core.Exchange),

		createPassword:   d.theme.Button(new(widget.Clickable), "Create password"),
		login:            d.theme.Button(new(widget.Clickable), "Login"),
		createNewWallet:  d.theme.Button(new(widget.Clickable), "Add"),
		unlockWallet:     d.theme.Button(new(widget.Clickable), "Unlock"),
		addCertFile:      d.theme.Button(new(widget.Clickable), "Add a file"),
		addDexServer:     d.theme.Button(new(widget.Clickable), "Submit"),
		register:         d.theme.Button(new(widget.Clickable), "Register"),
		toWallet:         d.theme.PlainIconButton(new(widget.Clickable), common.icons.cached),
		choseAssetUnlock: d.theme.PlainIconButton(new(widget.Clickable), common.icons.lock),

		appPassword:      d.theme.EditorPassword(new(widget.Editor), "Password"),
		appPasswordAgain: d.theme.EditorPassword(new(widget.Editor), "Password Again"),
		accountName:      d.theme.Editor(new(widget.Editor), "Account Name"),
		walletPassword:   d.theme.EditorPassword(new(widget.Editor), "Wallet Password"),
		dexServerAddress: d.theme.Editor(new(widget.Editor), "DEX Address"),

		walletActionWidgets: make(map[string]*walletActionWidget),
		walletActionInfo: &walletActionInfo{
			image:    coinImageBySymbol(&common.icons, dex.DefaultAsset),
			coin:     dex.DefaultAsset,
			coinName: "Decred",
			coinID:   dex.DefaultAssetID,
		},
	}

	iconColor := common.theme.Color.Gray3
	pg.toWallet.Color = iconColor
	pg.choseAssetUnlock.Color = iconColor

	pg.dexServerAddress.Editor.SetText("http://127.0.0.1:7232")

	pg.appPassword.Editor.SetText("")
	pg.appPassword.Editor.SingleLine = true
	pg.appPasswordAgain.Editor.SetText("")
	pg.appPasswordAgain.Editor.SingleLine = true

	// Get initial values
	for _, asset := range common.dexc.SupportedAsset() {
		pg.supportedAsset = append(pg.supportedAsset, asset)
		pg.walletActionWidgets[asset.Symbol] = &walletActionWidget{
			evt: &gesture.Click{},
			info: &walletActionInfo{
				image:    coinImageBySymbol(&common.icons, asset.Symbol),
				coin:     asset.Symbol,
				coinName: asset.Info.Name,
				coinID:   asset.ID,
			},
		}
	}

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *marketsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	dims := common.Layout(gtx, func(gtx C) D {
		card := pg.theme.Card()
		card.Radius = decredmaterial.CornerRadius{}
		return card.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min = gtx.Constraints.Max
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.navDrawerLayout(gtx, common)
				}),
				layout.Rigid(func(gtx C) D {
					l := common.theme.Line(gtx.Constraints.Max.X, 1)
					l.Color = common.theme.Color.Gray1
					return l.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.marketsLayout(gtx, common)
				}),
			)
		})
	})

	u := ((*pg.user).Info)

	if !u.Initialized {
		return pg.initAppPasswordModal(gtx, common)
	}

	if !pg.isLoggedIn && u.Initialized {
		return pg.loginModal(gtx, common)
	}

	// Show add wallet from initialize
	if len(u.Exchanges) == 0 && u.Initialized && u.Assets[dex.DefaultAssetID].Wallet == nil {
		return pg.createNewWalletModal(gtx, common)
	}

	// Show unlock wallet from initialize
	if u.Assets[dex.DefaultAssetID] != nil &&
		u.Assets[dex.DefaultAssetID].Wallet != nil &&
		!u.Assets[dex.DefaultAssetID].Wallet.Open {
		return pg.unlockWalletModal(gtx, common)
	}

	if len(u.Exchanges) == 0 &&
		u.Assets[dex.DefaultAssetID] != nil &&
		u.Assets[dex.DefaultAssetID].Wallet.Open &&
		!pg.showConfirmRegister {
		return pg.addNewDexModal(gtx, common)
	}

	if pg.showConfirmRegister {
		return pg.confirmRegisterModal(gtx, common)
	}

	// Show add wallet from market page
	if pg.showAddWallet {
		return pg.createNewWalletModal(gtx, common)
	}

	// Show unlock wallet from market page
	if pg.showUnlockWallet {
		return pg.unlockWalletModal(gtx, common)
	}

	return dims
}

const (
	navDrawerWidth          = 120
	navDrawerMinimizedWidth = 72
)

func (pg *marketsPage) navDrawerLayout(gtx layout.Context, c pageCommon) layout.Dimensions {
	width := navDrawerWidth
	gtx.Constraints.Min.X = int(gtx.Metric.PxPerDp) * width

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			inset := layout.Inset{Left: values.MarginPadding5}

			return list.Layout(gtx, len(pg.drawerNavItems), func(gtx C, i int) D {
				host := pg.drawerNavItems[i].host

				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						txt := pg.theme.Label(values.TextSize14, host)
						return inset.Layout(gtx, func(gtx C) D {
							return txt.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						listMkt := layout.List{Axis: layout.Vertical}

						return listMkt.Layout(gtx, len(pg.drawerNavItems[i].elements), func(gtx C, mktIndex int) D {
							element := pg.drawerNavItems[i].elements[mktIndex]
							click := element.evt
							pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
							click.Add(gtx.Ops)
							pg.navItemHandler(gtx, &element)
							return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											img := element.imageLeft
											if img == nil {
												return layout.Dimensions{}
											}
											img.Scale = 0.13
											return inset.Layout(gtx, func(gtx C) D {
												return img.Layout(gtx)
											})
										}),
										layout.Rigid(func(gtx C) D {
											img := element.imageRight
											if img == nil {
												return layout.Dimensions{}
											}
											img.Scale = 0.13
											return inset.Layout(gtx, func(gtx C) D {
												return img.Layout(gtx)
											})
										}),
									)
								}),
								layout.Rigid(func(gtx C) D {
									return inset.Layout(gtx, func(gtx C) D {
										txt := pg.theme.Label(values.TextSize16, element.asset.name)
										return txt.Layout(gtx)
									})
								}),
							)
						})
					}),
				)
			})
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.SE.Layout(gtx, func(gtx C) D {
				return pg.toWallet.Layout(gtx)
			})
		}),
	)
}

func (pg *marketsPage) marketsLayout(gtx layout.Context, c pageCommon) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.theme.H5("Dex page").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return c.UniformPadding(gtx, func(gtx C) D {
				return pg.marketBalancesLayout(gtx, &c)
			})
		}),
	)
}

func (pg *marketsPage) marketBalancesLayout(gtx layout.Context, c *pageCommon) layout.Dimensions {
	border := widget.Border{Color: c.theme.Color.Gray1, CornerRadius: values.MarginPadding8, Width: values.MarginPadding1}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return border.Layout(gtx, func(gtx C) D {
		u := (*pg.user).Info

		col := func(gtx layout.Context, ic *widget.Image, market string, wallState *core.WalletState) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Alignment: layout.Baseline, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
										return ic.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									txt := pg.theme.Label(values.TextSize16, strings.ToUpper(market))
									return txt.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							if wallState == nil {
								return layout.Dimensions{}
							}
							var ic *widget.Icon
							if wallState.Open {
								ic = c.icons.lockOpen
								ic.Color = pg.theme.Color.Success
							} else {
								ic = c.icons.lock
								pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
								pg.walletActionWidgets[market].evt.Add(gtx.Ops)
								pg.unlockWalletHandler(gtx, pg.walletActionWidgets[market])
							}
							return ic.Layout(gtx, values.MarginPadding15)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					if wallState == nil {
						b := widget.Border{Color: c.theme.Color.Gray1, CornerRadius: values.MarginPadding8, Width: values.MarginPadding1}
						gtx.Constraints.Max.X = 120
						pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
						pg.walletActionWidgets[market].evt.Add(gtx.Ops)
						pg.addWalletHandler(gtx, pg.walletActionWidgets[market])
						return b.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
								return pg.theme.
									Label(values.TextSize14, fmt.Sprintf("Add a %s wallet", strings.ToUpper(market))).
									Layout(gtx)
							})
						})
					}
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							bal := dcrutil.Amount(wallState.Balance.Available).ToCoin()
							return pg.theme.Label(values.TextSize14, fmt.Sprintf("%v", bal)).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							bal := dcrutil.Amount(wallState.Balance.Locked).ToCoin()
							return pg.theme.Label(values.TextSize14, fmt.Sprintf("%v", bal)).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							bal := dcrutil.Amount(wallState.Balance.Immature).ToCoin()
							return pg.theme.Label(values.TextSize14, fmt.Sprintf("%v", bal)).Layout(gtx)
						}),
					)
				}),
			)
		}

		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Flexed(.2, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.theme.Label(values.TextSize16, "Balances").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return pg.theme.Label(values.TextSize14, "available").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return pg.theme.Label(values.TextSize14, "locked").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return pg.theme.Label(values.TextSize14, "immature").Layout(gtx)
						}),
					)
				}),
				layout.Flexed(.4, func(gtx C) D {
					ic := coinImageBySymbol(&c.icons, (*pg.selectedMaket).marketBase)
					if ic == nil {
						return layout.Dimensions{}
					}
					return col(gtx, ic, (*pg.selectedMaket).marketBase, u.Assets[(*pg.selectedMaket).marketBaseID].Wallet)
				}),
				layout.Flexed(.4, func(gtx C) D {
					ic := coinImageBySymbol(&c.icons, (*pg.selectedMaket).marketQuote)
					if ic == nil {
						return layout.Dimensions{}
					}
					return col(gtx, ic, (*pg.selectedMaket).marketQuote, u.Assets[(*pg.selectedMaket).marketQuoteID].Wallet)
				}),
			)
		})
	})
}

func (pg *marketsPage) initAppPasswordModal(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.pageModal.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return pg.theme.Label(values.TextSize20, "Set App Password").Layout(gtx)
		},
		func(gtx C) D {
			return pg.theme.Label(values.TextSize14, "Set your app password. This password will protect your DEX account keys and connected wallets.").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return pg.appPassword.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return pg.appPasswordAgain.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return pg.createPassword.Layout(gtx)
		},
	}, 900)
}

func (pg *marketsPage) loginModal(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.pageModal.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return pg.theme.Label(values.TextSize20, "Login").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return pg.appPassword.Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.login.Layout(gtx)
		},
	}, 900)
}

func (pg *marketsPage) createNewWalletModal(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.pageModal.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.theme.Label(values.TextSize20, "Add a").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						ic := pg.walletActionInfo.image
						ic.Scale = 0.2
						return pg.walletActionInfo.image.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.theme.Label(values.TextSize20, fmt.Sprintf("%s Wallet", pg.walletActionInfo.coinName)).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if pg.walletActionInfo.coinID == dex.DefaultAssetID {
						return pg.theme.Label(values.TextSize14, "Your Decred wallet is required to pay registration fees.").Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return pg.accountName.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return pg.walletPassword.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return pg.appPassword.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return pg.createNewWallet.Layout(gtx)
		},
	}, 900)
}

func (pg *marketsPage) unlockWalletModal(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.pageModal.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.theme.Label(values.TextSize20, "Unlock").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						ic := pg.walletActionInfo.image
						ic.Scale = 0.2
						return pg.walletActionInfo.image.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.theme.Label(values.TextSize20, fmt.Sprintf("%s Wallet", pg.walletActionInfo.coinName)).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return pg.theme.Label(values.TextSize14, `App Password
Your app password is always required when performing sensitive wallet operations.`).Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return pg.appPassword.Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.unlockWallet.Layout(gtx)
		},
	}, 900)
}

func (pg *marketsPage) addNewDexModal(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.pageModal.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return pg.theme.Label(values.TextSize20, "Add a dex").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return pg.dexServerAddress.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return pg.addCertFile.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return pg.addDexServer.Layout(gtx)
		},
	}, 900)
}

func (pg *marketsPage) confirmRegisterModal(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.pageModal.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return pg.theme.Label(values.TextSize20, "Confirm Registration").Layout(gtx)
		},
		func(gtx C) D {
			return pg.theme.Label(values.TextSize14, "Enter your app password to confirm DEX registration. When you submit this form, 1.000 DCR will be spent from your Decred wallet to pay registration fees.").Layout(gtx)
		},
		func(gtx C) D {
			return pg.theme.Label(values.TextSize14, "The DCR lot size for this DEX is 1.000 DCR. All trades are in multiples of this lot size. This is the minimum possible trade amount in DCR.").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return pg.appPassword.Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.register.Layout(gtx)
		},
	}, 900)
}

func (pg *marketsPage) initDrawerNavItems(c *pageCommon) {
	if len(pg.drawerNavItems) == len((*pg.user).Info.Exchanges) {
		return
	}

	pg.drawerNavItems = make([]*drawerNav, 0)
	for h, ex := range (*pg.user).Info.Exchanges {
		dn := &drawerNav{host: h}
		for _, mkt := range ex.Markets {
			dn.elements = append(dn.elements, navItem{
				evt:        &gesture.Click{},
				imageLeft:  coinImageBySymbol(&c.icons, mkt.BaseSymbol),
				imageRight: coinImageBySymbol(&c.icons, mkt.QuoteSymbol),
				asset: &selectedMaket{
					host:          h,
					name:          strings.ToUpper(mkt.BaseSymbol + "-" + mkt.QuoteSymbol),
					marketBase:    mkt.BaseSymbol,
					marketBaseID:  mkt.BaseID,
					marketQuote:   mkt.QuoteSymbol,
					marketQuoteID: mkt.QuoteID,
				},
			})
		}

		pg.drawerNavItems = append(pg.drawerNavItems, dn)
	}
}

func (pg *marketsPage) addWalletHandler(gtx layout.Context, wg *walletActionWidget) {
	for _, e := range wg.evt.Events(gtx) {
		if e.Type == gesture.TypeClick {
			pg.walletActionInfo = &walletActionInfo{
				image:    wg.info.image,
				coin:     wg.info.coin,
				coinName: wg.info.coinName,
			}
			pg.showAddWallet = true
		}
	}
}

func (pg *marketsPage) unlockWalletHandler(gtx layout.Context, wg *walletActionWidget) {
	for _, e := range wg.evt.Events(gtx) {
		if e.Type == gesture.TypeClick {
			pg.walletActionInfo = &walletActionInfo{
				image:    wg.info.image,
				coin:     wg.info.coin,
				coinName: wg.info.coinName,
			}
			pg.showUnlockWallet = true
		}
	}
}

func (pg *marketsPage) navItemHandler(gtx layout.Context, navItem *navItem) {
	for _, e := range navItem.evt.Events(gtx) {
		if e.Type == gesture.TypeClick {
			(*pg.selectedMaket) = navItem.asset
		}
	}
}

func (pg *marketsPage) handle(common pageCommon) {
	pg.initDrawerNavItems(&common)

	if pg.createPassword.Button.Clicked() {
		if pg.appPasswordAgain.Editor.Text() != pg.appPassword.Editor.Text() {
			return
		}

		common.dexc.InitializeClient(pg.appPassword.Editor.Text(), pg.errInitappChan)
	}

	if pg.login.Button.Clicked() {
		common.dexc.Login(pg.appPassword.Editor.Text(), pg.errLoginChan)
	}

	if pg.createNewWallet.Button.Clicked() {
		coinID := pg.walletActionInfo.coinID
		config := common.dexc.AutoWalletConfig(coinID)

		for assetID, supportedAsset := range (*pg.user).Info.Assets {
			if assetID == coinID {
				for _, cfgOpt := range supportedAsset.Info.ConfigOpts {
					if key := cfgOpt.Key; key == "fallbackfee" ||
						key == "feeratelimit" ||
						key == "redeemconftarget" ||
						key == "rpcbind" ||
						key == "rpcport" ||
						key == "txsplit" {
						config[key] = fmt.Sprintf("%v", cfgOpt.DefaultValue)
					}
				}
			}
		}

		// Bitcoin
		config["walletname"] = pg.accountName.Editor.Text()

		// Decred
		config["account"] = pg.accountName.Editor.Text()
		config["password"] = pg.walletPassword.Editor.Text()
		config["username"] = config["rpcuser"]
		config["rpcport"] = "18332"

		form := &dex.NewWalletForm{
			AssetID: coinID,
			Config:  config,
			Pass:    []byte(pg.walletPassword.Editor.Text()),
			AppPW:   []byte(pg.appPassword.Editor.Text()),
		}
		common.dexc.AddNewWallet(form, pg.errCreateWalletChan)
	}

	if pg.unlockWallet.Button.Clicked() {
		common.dexc.UnlockWallet(pg.walletActionInfo.coinID, []byte(pg.appPassword.Editor.Text()), pg.errUnlockWallChan)
	}

	if pg.addCertFile.Button.Clicked() {
		go func() {
			filename, err := dialog.File().Filter("Select TLS Certificate", "cert").Load()

			if err != nil {
				log.Error(err)
			}

			content, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Error(err)
			}
			pg.cert = string(content)
		}()
	}

	if pg.addDexServer.Button.Clicked() {
		common.dexc.GetDEXConfig(pg.dexServerAddress.Editor.Text(), pg.cert, pg.errCreateWalletChan, pg.responseGetDexChan)
	}

	if pg.register.Button.Clicked() {
		common.dexc.Register(pg.appPassword.Editor.Text(), pg.responseGetDex.Host, pg.responseGetDex.Fee.Amt, pg.cert, pg.errRegisterChan)
	}

	if pg.toWallet.Button.Clicked() {
		*common.switchView = values.WalletView
	}

	select {
	case err := <-pg.errInitappChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		pg.appPassword.Editor.SetText("")
		pg.isLoggedIn = true
		common.dexc.GetUser()
	case err := <-pg.errLoginChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		common.dexc.GetUser()
		pg.isLoggedIn = true
		pg.appPassword.Editor.SetText("")

	case err := <-pg.errUnlockWallChan:
		pg.appPassword.ClearError()
		if err != nil {
			pg.appPassword.SetError(err.Error())
			common.notify(err.Error(), false)
			return
		}
		pg.showUnlockWallet = false
		common.dexc.GetUser()
		pg.appPassword.Editor.SetText("")
	case err := <-pg.errCreateWalletChan:
		if err != nil {
			log.Error(err)
			common.notify(err.Error(), false)
			return
		}
		pg.showAddWallet = false
		common.dexc.GetUser()
		pg.appPassword.Editor.SetText("")
	case resp := <-pg.responseGetDexChan:
		pg.responseGetDex = resp
		pg.showConfirmRegister = true
		pg.appPassword.Editor.SetText("")
	case err := <-pg.errRegisterChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		pg.showConfirmRegister = false
		common.dexc.GetUser()
		pg.appPassword.Editor.SetText("")
	default:
	}
}
