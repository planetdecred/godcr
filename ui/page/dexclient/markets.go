package dexclient

import (
	"context"
	"fmt"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
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
	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	addDexWidget *addDexWidget
	login        decredmaterial.Button
	initialize   decredmaterial.Button
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
		Load:       l,
		login:      l.Theme.Button("Login"),
		initialize: l.Theme.Button("Start using now"),
	}

	return pg
}

func (pg *Page) ID() string {
	return MarketPageID
}

func (pg *Page) Layout(gtx C) D {
	var body func(gtx C) D
	switch {
	case !pg.Dexc().Initialized() || !pg.Dexc().IsLoggedIn():
		body = func(gtx C) D {
			return pg.pageSections(gtx, pg.welcomeLayout)
		}
	case pg.addDexWidget != nil:
		body = func(gtx C) D {
			return pg.pageSections(gtx, pg.addDexWidget.layout)
		}
	default:
		body = func(gtx C) D {
			return pg.pageSections(gtx, pg.registrationStatusLayout)
		}
	}

	return components.UniformPadding(gtx, body)
}

func (pg *Page) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding16).Layout(gtx, body)
	})
}

func (pg *Page) welcomeLayout(gtx C) D {
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				description := "Trade crypto peer-to-peer.\nNo trading fees. No KYC."
				return layout.Center.Layout(gtx, pg.Theme.H5(description).Layout)
			}),
			layout.Rigid(func(gtx C) D {
				if !pg.Dexc().Initialized() {
					return pg.initialize.Layout(gtx)
				}
				return pg.login.Layout(gtx)
			}),
		)
	})
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
		// TODO: render another UI by dex and wallet status
		return D{}
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

	if len(pg.Dexc().DEXServers()) == 0 { // Connect a DEX to proceed. May require adding a wallet to pay the fee.
		pg.addDexWidget = newAddDexWidget(pg.Load)
	}
}

func (pg *Page) OnClose() {
	pg.ctxCancel()
}

func (pg *Page) Handle() {
	if pg.login.Button.Clicked() {
		modal.NewPasswordModal(pg.Load).
			Title("Login").
			Hint("App password").
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton("Login", func(password string, pm *modal.PasswordModal) bool {
				go func() {
					err := pg.Dexc().Login([]byte(password))
					if err != nil {
						pm.SetError(err.Error())
						pm.SetLoading(false)
						return
					}
					pm.Dismiss()
				}()
				return false
			}).Show()
	}

	if pg.initialize.Button.Clicked() {
		modal.NewCreatePasswordModal(pg.Load).
			Title("Set App Password").
			SetDescription("Set your app password. This password will protect your DEX account keys and connected wallets.").
			EnableName(false).
			PasswordHint("Password").
			ConfirmPasswordHint("Confirm password").
			PasswordCreated(func(walletName, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					err := pg.Dexc().InitializeWithPassword([]byte(password))
					if err != nil {
						m.SetError(err.Error())
						m.SetLoading(false)
						return
					}
					pg.Toast.Notify("App password created")
					m.Dismiss()
				}()
				return false
			}).Show()
	}

	if pg.addDexWidget != nil {
		pg.addDexWidget.handle()
	}
}
