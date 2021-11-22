package dexclient

import (
	"context"
	"fmt"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"

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
	ctx       context.Context // page context
	ctxCancel context.CancelFunc
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
		Load: l,
	}

	return pg
}

func (pg *Page) ID() string {
	return MarketPageID
}

func (pg *Page) Layout(gtx C) D {
	return pg.registrationStatusLayout(gtx)
}

func (pg *Page) dex() *core.Exchange {
	// TODO: Should ideally pick a DEX by host, but this currently
	// picks the first DEX in the map, if one has been previously
	// connected. This is okay because there's only one DEX that
	// can be connected for now.
	for _, dex := range pg.Dexc().DEXServers() {
		return dex
	}
	return nil
}

func (pg *Page) registrationStatusLayout(gtx C) D {
	dex := pg.dex()
	if dex == nil || dex.PendingFee == nil {
		return pg.Theme.Label(values.TextSize14, "Ready to trade").Layout(gtx)
	}
	reqConfirms, currentConfs := dex.Fee.Confs, dex.PendingFee.Confs
	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(pg.Theme.Label(values.TextSize14, "Waiting for confirmations...").Layout),
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

func (pg *Page) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	go pg.readNotifications()
}

func (pg *Page) Handle() {
	switch {
	case !pg.Dexc().Initialized(): // Must initialize to proceed.
		pg.ChangeFragment(NewDexPasswordPage(pg.Load))
	case !pg.Dexc().IsLoggedIn(): // Initialized client must log in.
		pg.ChangeFragment(NewDexLoginPage(pg.Load))
	case len(pg.Dexc().DEXServers()) == 0: // Connect a DEX to proceed. May require adding a wallet to pay the fee.
		pg.ChangeFragment(NewAddDexPage(pg.Load))
	}
}

func (pg *Page) OnClose() {
	pg.ctxCancel()
}
