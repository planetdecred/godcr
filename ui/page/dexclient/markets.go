package dexclient

import (
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const MarketPageID = "Markets"

type Page struct {
	*load.Load
	user             *core.User
	miniTradeFormWdg *miniTradeFormWidget
	initializeModal  bool
	orderBook        *core.OrderBook

	addBTCWallet decredmaterial.Button
}

func NewMarketPage(l *load.Load) *Page {
	pg := &Page{
		Load:             l,
		user:             new(core.User),
		miniTradeFormWdg: newMiniTradeFormWidget(l),
		initializeModal:  false,
		orderBook:        new(core.OrderBook),

		addBTCWallet: l.Theme.Button(new(widget.Clickable), "Add BTC wallet"),
	}

	return pg
}

func (pg *Page) ID() string {
	return MarketPageID
}

func (pg *Page) OnResume() {
	pg.refreshUser()
}

func (pg *Page) Layout(gtx C) D {
	dims := components.UniformPadding(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.pageSections(gtx, func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding20,
						Bottom: values.MarginPadding20,
					}.Layout(gtx, func(gtx C) D {
						return pg.miniTradeFormWdg.layout(gtx)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.E.Layout(gtx, func(gtx C) D {
					return pg.addBTCWallet.Layout(gtx)
				})
			}),
		)
	})

	return dims
}

func (pg *Page) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.Inset{
		Bottom: values.MarginPadding8,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, body)
		})
	})
}

func (pg *Page) Handle() {
	if !pg.initializeModal {
		pg.initializeModal = true
		pg.handleModals()
	}

	pg.miniTradeFormWdg.handle()

	if pg.addBTCWallet.Button.Clicked() {
		asset := "btc"
		wallInfo := &walletInfoWidget{
			image:    coinImageBySymbol(&pg.Load.Icons, asset),
			coin:     asset,
			coinName: "Bitcoin",
			coinID:   0,
		}
		md := newCreateWalletModal(pg.Load, wallInfo)
		md.walletCreated = func() {
			pg.refreshUser()
		}
		md.Show()
	}
}

func (pg *Page) OnClose() {
	pg.initializeModal = false
}

func (pg *Page) handleModals() {
	u := pg.user
	if !u.Initialized {
		md := newPasswordModal(pg.Load)
		md.appInitiated = func() {
			pg.refreshUser()
			pg.DL.IsLoggedIn = true
		}
		md.Show()
		return
	}

	if !pg.DL.IsLoggedIn && u.Initialized {
		md := newloginModal(pg.Load)
		md.loggedIn = func() {
			pg.refreshUser()
			pg.DL.IsLoggedIn = true
		}
		md.Show()
		return
	}

	// Show add wallet from initialize
	if len(u.Exchanges) == 0 &&
		u.Initialized &&
		u.Assets[dexc.DefaultAssetID].Wallet == nil {
		wallInfo := &walletInfoWidget{
			image:    coinImageBySymbol(&pg.Load.Icons, dexc.DefaultAsset),
			coin:     dexc.DefaultAsset,
			coinName: "Decred",
			coinID:   dexc.DefaultAssetID,
		}
		md := newCreateWalletModal(pg.Load, wallInfo)
		md.walletCreated = func() {
			pg.refreshUser()
		}
		md.Show()
		return
	}

	if u.Assets[dexc.DefaultAssetID] != nil &&
		u.Assets[dexc.DefaultAssetID].Wallet != nil &&
		!u.Assets[dexc.DefaultAssetID].Wallet.Open {
		md := newUnlockWalletModal(pg.Load)
		md.unlocked = func(password []byte) {
			pg.refreshUser()
			pg.connectDex(testDexHost, password)
		}
		md.Show()
		return
	}

	if len(u.Exchanges) == 0 {
		md := newAddDexModal(pg.Load)
		md.created = func(cert []byte, ce *core.Exchange) {
			cfModal := newConfirmRegisterModal(md.Load)
			cfModal.updateCertAndExchange(cert, ce)
			cfModal.confirmed = func(password []byte) {
				pg.refreshUser()
				pg.connectDex(ce.Host, password)
			}
			cfModal.Show()
		}
		md.Show()
		return
	}
}

func (pg *Page) refreshUser() {
	pg.user = pg.DL.Core.User()
	pg.initializeModal = false

	if pg.user.Initialized && pg.DL.IsLoggedIn {
		pg.connectDex(testDexHost, []byte("123"))
	}
}
