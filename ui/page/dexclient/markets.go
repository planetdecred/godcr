package dexclient

import (
	"context"
	"fmt"

	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"

	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const MarketPageID = "Markets"

const testDexHost = "dex-test.ssgen.io:7232"

type selectedMaket struct {
	host          string
	name          string
	marketBase    string
	marketQuote   string
	marketBaseID  uint32
	marketQuoteID uint32
}

type Page struct {
	*load.Load
	ctx             context.Context // page context
	ctxCancel       context.CancelFunc
	initializeModal bool
	addBTCWallet    decredmaterial.Button
	advancedTrade   decredmaterial.Button
	selectedMaket   *selectedMaket
}

// TODO: Aadd collapsible button to select a market.
var mkt = &selectedMaket{
	host:          testDexHost,
	name:          "DCR-BTC",
	marketBase:    "dcr",
	marketBaseID:  dcr.BipID,
	marketQuote:   "btc",
	marketQuoteID: btc.BipID,
}

func NewMarketPage(l *load.Load) *Page {
	pg := &Page{
		Load:            l,
		initializeModal: false,
		selectedMaket:   mkt,
	}

	return pg
}

func (pg *Page) ID() string {
	return MarketPageID
}

func (pg *Page) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.refreshUser()

	if pg.Dexc.Core.User().Initialized && pg.Dexc.IsLoggedIn {
		go pg.readNotifications()
	}
}

func (pg *Page) Layout(gtx C) D {
	return pg.registrationStatusLayout(gtx)
}

func (pg *Page) registrationStatusLayout(gtx C) D {
	dex := pg.Dexc.Core.User().Exchanges[pg.selectedMaket.host]
	if dex == nil || !dex.Connected {
		return D{}
	}
	if dex.PendingFee == nil {
		return D{}
	}
	reqConfirms, currentConfs := dex.Fee.Confs, dex.PendingFee.Confs
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.Theme.Label(values.TextSize14, "Waiting for confirmations...").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			t := fmt.Sprintf("In order to trade at %s, the registration fee payment needs %d confirmations.", pg.selectedMaket.host, reqConfirms)
			return pg.Theme.Label(values.TextSize14, t).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			t := fmt.Sprintf("%d/%d", currentConfs, reqConfirms)
			return pg.Theme.Label(values.TextSize14, t).Layout(gtx)
		}),
	)
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
}

func (pg *Page) OnClose() {
	pg.initializeModal = false
	pg.ctxCancel()
}

func (pg *Page) handleModals() {
	u := pg.Dexc.Core.User()

	// Must initialize to proceed.
	if !u.Initialized {
		md := newPasswordModal(pg.Load)
		md.appInitiated = func() {
			pg.refreshUser()
			pg.Dexc.IsLoggedIn = true
		}
		md.Show()
		return
	}

	// Initialized client must be logged in.
	if !pg.Dexc.IsLoggedIn {
		md := newloginModal(pg.Load)
		md.loggedIn = func(password []byte) {
			pg.refreshUser()
			pg.Dexc.IsLoggedIn = true
			if u.Assets[dexc.DefaultAssetID] != nil &&
				u.Assets[dexc.DefaultAssetID].Wallet != nil {
				pg.connectDex(pg.selectedMaket.host, password)
			}
		}
		md.Show()
		return
	}

	// The dcr wallet must be connected before registering with
	// a DEX.
	// TODO: Since other assets can now be used to pay the fee,
	// this shouldn't be a pre-requirement. Instead, attempt to
	// connect a DEX first and determine what wallet is required
	// for fee payment.
	if dcrWallet := u.Assets[dcr.BipID]; dcrWallet == nil || dcrWallet.Wallet == nil {
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
			pg.connectDex(pg.selectedMaket.host, password)
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

// TODO: Investigate uses of this method and where appropriate,
// do `pg.initializeModal = false` instead and remove this method.
func (pg *Page) refreshUser() {
	pg.initializeModal = false
}
