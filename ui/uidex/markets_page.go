package uidex

import (
	"fmt"
	"io/ioutil"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/dex"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/sqweek/dialog"
)

const PageMarkets = "MarketsPage"

type navItem struct {
	imageLeft  *widget.Image
	imageRight *widget.Image
	mkt        string
}

type drawerNav struct {
	host     string
	elements []navItem
}

type marketsPage struct {
	dexc      *dex.Dex
	theme     *decredmaterial.Theme
	pageModal *decredmaterial.Modal
	exchange  layout.List

	supportedAsset      []*core.SupportedAsset
	user                **dex.User
	defaultWalletConfig map[string]string
	cert                string
	certName            string
	drawerNavItems      []*drawerNav

	appPassword      decredmaterial.Editor
	appPasswordAgain decredmaterial.Editor
	accountName      decredmaterial.Editor
	walletPassword   decredmaterial.Editor
	dexServerAddress decredmaterial.Editor

	createPassword  decredmaterial.Button
	createNewWallet decredmaterial.Button
	unlockWallet    decredmaterial.Button
	login           decredmaterial.Button
	cancel          decredmaterial.Button
	addCertFile     decredmaterial.Button
	addDexServer    decredmaterial.Button
	register        decredmaterial.Button
	toWallet        decredmaterial.IconButton

	isAppInitialized    bool
	isLoggedIn          bool
	showConfirmRegister bool
	hideAllModal        bool
	errWalletChan       chan error
	errLoginChan        chan error
	errInitappChan      chan error
	errRegisterChan     chan error
	responseGetDex      *core.Exchange
	responseGetDexChan  chan *core.Exchange
}

func (d *Dex) MarketsPage(common pageCommon) layout.Widget {
	pg := &marketsPage{
		dexc:               d.dexc,
		theme:              common.theme,
		pageModal:          common.theme.Modal(),
		exchange:           layout.List{Axis: layout.Vertical},
		user:               &d.userInfo,
		drawerNavItems:     make([]*drawerNav, 0),
		errWalletChan:      make(chan error),
		errInitappChan:     make(chan error),
		errLoginChan:       make(chan error),
		errRegisterChan:    make(chan error),
		responseGetDexChan: make(chan *core.Exchange),

		createPassword:  d.theme.Button(new(widget.Clickable), "Create password"),
		login:           d.theme.Button(new(widget.Clickable), "Login"),
		cancel:          d.theme.Button(new(widget.Clickable), "Cancel"),
		createNewWallet: d.theme.Button(new(widget.Clickable), "Add"),
		unlockWallet:    d.theme.Button(new(widget.Clickable), "Unlock"),
		addCertFile:     d.theme.Button(new(widget.Clickable), "Add a file"),
		addDexServer:    d.theme.Button(new(widget.Clickable), "Submit"),
		register:        d.theme.Button(new(widget.Clickable), "Register"),
		toWallet:        d.theme.PlainIconButton(new(widget.Clickable), common.icons.cached),

		appPassword:      d.theme.EditorPassword(new(widget.Editor), "Password"),
		appPasswordAgain: d.theme.EditorPassword(new(widget.Editor), "Password Again"),
		accountName:      d.theme.Editor(new(widget.Editor), "Account Name"),
		walletPassword:   d.theme.EditorPassword(new(widget.Editor), "Wallet Password"),
		dexServerAddress: d.theme.Editor(new(widget.Editor), "DEX Address"),
	}

	pg.toWallet.Color = d.theme.Color.Black
	pg.dexServerAddress.Editor.SetText("http://127.0.0.1:7232")

	pg.appPassword.Editor.SetText("")
	pg.appPassword.Editor.SingleLine = true
	pg.appPasswordAgain.Editor.SetText("")
	pg.appPasswordAgain.Editor.SingleLine = true

	// Get initial values
	pg.isAppInitialized = pg.dexc.IsInitialized()
	for _, v := range pg.dexc.SupportedAsset() {
		pg.supportedAsset = append(pg.supportedAsset, v)
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
					return common.UniformPadding(gtx, func(gtx C) D {
						return pg.marketsLayout(gtx, common)
					})
				}),
			)
		})
	})

	// TODO: For testing purposes, will remove this
	if pg.hideAllModal {
		return dims
	}

	if !pg.isAppInitialized {
		return pg.initAppPasswordModal(gtx, common)
	}

	if !pg.isLoggedIn && pg.isAppInitialized {
		return pg.loginModal(gtx, common)
	}

	u := ((*pg.user).Info)
	if len(u.Exchanges) == 0 && u.Initialized && u.Assets[dex.DefaultAssert].Wallet == nil {
		return pg.createNewWalletModal(gtx, common)
	}

	if u.Assets[dex.DefaultAssert] != nil &&
		u.Assets[dex.DefaultAssert].Wallet != nil &&
		!u.Assets[dex.DefaultAssert].Wallet.Open {
		return pg.unlockWalletModal(gtx, common)
	}

	if len(u.Exchanges) == 0 &&
		u.Assets[dex.DefaultAssert] != nil &&
		u.Assets[dex.DefaultAssert].Wallet.Open &&
		!pg.showConfirmRegister {
		return pg.addNewDexModal(gtx, common)
	}

	if pg.showConfirmRegister {
		return pg.confirmRegisterModal(gtx, common)
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
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						txt := pg.theme.Label(values.TextSize14, pg.drawerNavItems[i].host)
						return inset.Layout(gtx, func(gtx C) D {
							return txt.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						listMkt := layout.List{Axis: layout.Vertical}

						return listMkt.Layout(gtx, len(pg.drawerNavItems[i].elements), func(gtx C, mktIndex int) D {
							element := pg.drawerNavItems[i].elements[mktIndex]

							return decredmaterial.Clickable(gtx, new(widget.Clickable), func(gtx C) D {
								return layout.Flex{Alignment: layout.Baseline}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												img := element.imageLeft
												img.Scale = 0.11
												return inset.Layout(gtx, func(gtx C) D {
													return img.Layout(gtx)
												})
											}),
											layout.Rigid(func(gtx C) D {
												img := element.imageRight
												img.Scale = 0.11
												return inset.Layout(gtx, func(gtx C) D {
													return img.Layout(gtx)
												})
											}),
										)
									}),
									layout.Rigid(func(gtx C) D {
										return inset.Layout(gtx, func(gtx C) D {
											t := strings.Join(strings.Split(element.mkt, "_"), "-")
											txt := pg.theme.Label(values.TextSize16, t)
											return txt.Layout(gtx)
										})
									}),
								)
							})
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
	)
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
		func(gtx C) D {
			return pg.cancel.Layout(gtx)
		},
	}, 900)
}

func (pg *marketsPage) createNewWalletModal(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.pageModal.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return pg.theme.Label(values.TextSize20, "Add a decred wallet").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.theme.Label(values.TextSize14, "Your Decred wallet is required to pay registration fees.").Layout(gtx)
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
			return pg.theme.Label(values.TextSize20, "Unlock wallet").Layout(gtx)
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
		for mkt := range ex.Markets {
			dn.elements = append(dn.elements, navItem{
				imageLeft:  c.icons.btc,
				imageRight: c.icons.dcr,
				mkt:        mkt,
			})
		}

		pg.drawerNavItems = append(pg.drawerNavItems, dn)
	}
}

func (pg *marketsPage) handle(common pageCommon) {
	pg.initDrawerNavItems(&common)

	if pg.createPassword.Button.Clicked() {
		if pg.appPasswordAgain.Editor.Text() != pg.appPassword.Editor.Text() {
			return
		}

		pg.dexc.InitializeClient(pg.appPassword.Editor.Text(), pg.errInitappChan)
	}

	if pg.login.Button.Clicked() {
		pg.dexc.Login(pg.appPassword.Editor.Text(), pg.errLoginChan)
	}

	if pg.createNewWallet.Button.Clicked() {
		config := pg.defaultWalletConfig

		for assetID, supportedAsset := range (*pg.user).Info.Assets {
			if assetID == dex.DefaultAssert {
				for _, cfgOpt := range supportedAsset.Info.ConfigOpts {
					if key := cfgOpt.Key; key == "fallbackfee" ||
						key == "feeratelimit" ||
						key == "redeemconftarget" ||
						key == "txsplit" {
						config[key] = fmt.Sprintf("%s", cfgOpt.DefaultValue)
					}
				}
			}
		}

		config["account"] = pg.accountName.Editor.Text()
		config["password"] = pg.walletPassword.Editor.Text()
		config["username"] = "song"
		form := &dex.NewWalletForm{
			AssetID: dex.DefaultAssert, // Asset of dcr
			Config:  config,
			Pass:    []byte(pg.walletPassword.Editor.Text()),
			AppPW:   []byte(pg.appPassword.Editor.Text()),
		}

		pg.dexc.AddNewWallet(form, pg.errWalletChan)
		pg.appPassword.Editor.SetText("")
	}

	if pg.unlockWallet.Button.Clicked() {
		pg.dexc.UnlockWallet(dex.DefaultAssert, []byte(pg.appPassword.Editor.Text()))
		pg.dexc.GetUser()
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
		pg.dexc.GetDEXConfig(pg.dexServerAddress.Editor.Text(), pg.cert, pg.errWalletChan, pg.responseGetDexChan)
	}

	if pg.register.Button.Clicked() {
		pg.dexc.Register(pg.appPassword.Editor.Text(), pg.responseGetDex.Host, pg.responseGetDex.Fee.Amt, pg.cert, pg.errRegisterChan)
	}

	if pg.toWallet.Button.Clicked() {
		*common.switchView = values.WalletView
	}

	if pg.cancel.Button.Clicked() {
		pg.hideAllModal = true
	}

	// for i := range pg.drawerNavItems {
	// 	for pg.drawerNavItems[i].clickable.Clicked() {
	// 		log.Info("Switch market...")
	// 	}
	// }

	select {
	case err := <-pg.errInitappChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		pg.isAppInitialized = true
		pg.appPassword.Editor.SetText("")
		pg.isLoggedIn = true
		pg.dexc.GetUser()
		pg.defaultWalletConfig = pg.dexc.GetDefaultWalletConfig()
	case err := <-pg.errLoginChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		pg.dexc.GetUser()
		pg.isLoggedIn = true
		df := pg.dexc.GetDefaultWalletConfig()
		pg.defaultWalletConfig = df
		pg.appPassword.Editor.SetText("")

	case err := <-pg.errWalletChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		pg.dexc.GetUser()
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
		pg.dexc.GetUser()
		pg.appPassword.Editor.SetText("")
	default:
	}
}
