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

type marketsPage struct {
	dexc      *dex.Dex
	theme     *decredmaterial.Theme
	pageModal *decredmaterial.Modal
	exchange  layout.List

	supportedAsset      []*core.SupportedAsset
	user                *core.User
	defaultWalletConfig map[string]string
	cert                string
	certName            string

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
	toWallet        decredmaterial.Button

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

func (win *Dex) MarketsPage(common pageCommon) layout.Widget {
	pg := &marketsPage{
		dexc:               win.dexc,
		theme:              common.theme,
		pageModal:          common.theme.Modal(),
		exchange:           layout.List{Axis: layout.Vertical},
		user:               new(core.User),
		errWalletChan:      make(chan error),
		errInitappChan:     make(chan error),
		errLoginChan:       make(chan error),
		errRegisterChan:    make(chan error),
		responseGetDexChan: make(chan *core.Exchange),

		createPassword:  common.theme.Button(new(widget.Clickable), "Create password"),
		login:           common.theme.Button(new(widget.Clickable), "Login"),
		cancel:          common.theme.Button(new(widget.Clickable), "Cancel"),
		createNewWallet: common.theme.Button(new(widget.Clickable), "Add"),
		unlockWallet:    common.theme.Button(new(widget.Clickable), "Unlock"),
		addCertFile:     common.theme.Button(new(widget.Clickable), "Add a file"),
		addDexServer:    common.theme.Button(new(widget.Clickable), "Submit"),
		register:        common.theme.Button(new(widget.Clickable), "Register"),
		toWallet:        common.theme.Button(new(widget.Clickable), "To Wallet"),

		appPassword:      win.theme.EditorPassword(new(widget.Editor), "Password"),
		appPasswordAgain: win.theme.EditorPassword(new(widget.Editor), "Password Again"),
		accountName:      win.theme.Editor(new(widget.Editor), "Account Name"),
		walletPassword:   win.theme.EditorPassword(new(widget.Editor), "Wallet Password"),
		dexServerAddress: win.theme.Editor(new(widget.Editor), "DEX Address"),
	}

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
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min = gtx.Constraints.Max
			return common.UniformPadding(gtx, func(gtx C) D {
				return pg.marketsLayout(gtx, common)
			})
		})
	})

	// For testing purposes, will remove this
	if pg.hideAllModal {
		return dims
	}

	if !pg.isAppInitialized {
		return pg.initAppPasswordModal(gtx, common)
	}

	if !pg.isLoggedIn && pg.isAppInitialized {
		return pg.loginModal(gtx, common)
	}

	if len(pg.user.Exchanges) == 0 && pg.user.Initialized && pg.user.Assets[dex.DefaultAssert].Wallet == nil {
		return pg.createNewWalletModal(gtx, common)
	}

	if pg.user.Assets[dex.DefaultAssert] != nil &&
		pg.user.Assets[dex.DefaultAssert].Wallet != nil &&
		!pg.user.Assets[dex.DefaultAssert].Wallet.Open {
		return pg.unlockWalletModal(gtx, common)
	}

	if len(pg.user.Exchanges) == 0 && pg.user.Assets[dex.DefaultAssert].Wallet.Open && !pg.showConfirmRegister {
		return pg.addNewDexModal(gtx, common)
	}

	if pg.showConfirmRegister {
		return pg.confirmRegisterModal(gtx, common)
	}

	return dims
}

func (pg *marketsPage) marketsLayout(gtx layout.Context, c pageCommon) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.theme.H5("Dex page").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			var list []*core.Exchange
			for _, ex := range pg.user.Exchanges {
				list = append(list, ex)
			}
			return pg.exchange.Layout(gtx, len(list), func(gtx layout.Context, i int) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
							return pg.theme.Label(values.TextSize16, list[i].Host).Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						var elements []layout.FlexChild
						for _, mkt := range list[i].Markets {
							t := strings.Join(strings.Split(mkt.Name, "_"), "-")
							elements = append(elements, layout.Rigid(func(gtx C) D {
								return pg.theme.Label(values.TextSize16, strings.ToUpper(t)).Layout(gtx)
							}))
						}
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx, elements...)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return pg.toWallet.Layout(gtx)
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

func (pg *marketsPage) handle(common pageCommon) {
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
		for assetID, supportedAsset := range pg.user.Assets {
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
		pg.user = pg.dexc.GetUser()
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

	select {
	case err := <-pg.errInitappChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		pg.isAppInitialized = true
		pg.appPassword.Editor.SetText("")
		pg.isLoggedIn = true
		pg.user = pg.dexc.GetUser()
		pg.defaultWalletConfig = pg.dexc.GetDefaultWalletConfig()
	case err := <-pg.errLoginChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		pg.user = pg.dexc.GetUser()
		pg.isLoggedIn = true
		df := pg.dexc.GetDefaultWalletConfig()
		pg.defaultWalletConfig = df
		pg.appPassword.Editor.SetText("")

	case err := <-pg.errWalletChan:
		if err != nil {
			common.notify(err.Error(), false)
			return
		}
		pg.user = pg.dexc.GetUser()
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
		pg.user = pg.dexc.GetUser()
		pg.appPassword.Editor.SetText("")
	default:
	}
}
