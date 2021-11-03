package dexclient

import (
	"context"
	"fmt"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const MarketPageID = "Markets"

type Page struct {
	*load.Load
	ctx             context.Context // page context
	ctxCancel       context.CancelFunc
	initializeModal bool
	addBTCWallet    decredmaterial.Button
	advancedTrade   decredmaterial.Button
}

// TODO: Add collapsible button to select a market.
// Use mktName="DCR-BTC" in the meantime.
func (pg *Page) dexMarket(mktName string) *core.Market {
	dex := pg.dex()
	if dex == nil {
		return nil
	}
	return dex.Markets[mktName]
}

func NewMarketPage(l *load.Load) *Page {
	pg := &Page{
		Load:            l,
		initializeModal: false,
	}

	return pg
}

func (pg *Page) ID() string {
	return MarketPageID
}

func (pg *Page) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	go pg.readNotifications()
}

func (pg *Page) Layout(gtx C) D {
	return pg.registrationStatusLayout(gtx)
}

func (pg *Page) dex() *core.Exchange {
	// TODO: Should ideally pick a DEX by host, but this currently
	// picks the first DEX in the map, if one has been previously
	// connected. This is okay because there's only one DEX that
	// can be connected for now.
	for _, dex := range pg.Dexc.Core.User().Exchanges {
		return dex
	}
	return nil
}

func (pg *Page) registrationStatusLayout(gtx C) D {
	dex := pg.dex()
	if dex == nil || !dex.Connected || dex.PendingFee == nil { // TODO: We should probably show the status even if !dex.Connected but dex.PendingFee is set.
		return D{}
	}
	reqConfirms, currentConfs := dex.Fee.Confs, dex.PendingFee.Confs
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.Theme.Label(values.TextSize14, "Waiting for confirmations...").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			t := fmt.Sprintf("In order to trade at %s, the registration fee payment needs %d confirmations.", dex.Host, reqConfirms)
			return pg.Theme.Label(values.TextSize14, t).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			t := fmt.Sprintf("%d/%d", currentConfs, reqConfirms)
			return pg.Theme.Label(values.TextSize14, t).Layout(gtx)
		}),
	)
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

	switch {
	case !u.Initialized: // Must initialize to proceed.
		md := newPasswordModal(pg.Load)
		md.appInitiated = func() {
			pg.initializeModal = false
			pg.Dexc.IsLoggedIn = true
		}
		md.Show()

	case !pg.Dexc.IsLoggedIn: // Initialized client must be logged in.
		md := newloginModal(pg.Load)
		md.loggedIn = func() {
			pg.initializeModal = false
			pg.Dexc.IsLoggedIn = true
		}
		md.Show()

	case len(u.Exchanges) == 0: // Connect a DEX to proceed. May require adding a wallet to pay the fee.
		md := newAddDexModal(pg.Load)
		md.done = func() {
			pg.initializeModal = false
		}
		md.Show()
	}
}
