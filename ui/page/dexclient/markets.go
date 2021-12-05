package dexclient

import (
	"context"
	"fmt"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/client/db"
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
	ctx        context.Context
	ctxCancel  context.CancelFunc
	login      decredmaterial.Button
	initialize decredmaterial.Button
	addDex     decredmaterial.Button
	sync       decredmaterial.Button
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
		addDex:     l.Theme.Button("Add a dex"),
		sync:       l.Theme.Button("Start sync to continue"),
	}

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *Page) ID() string {
	return MarketPageID
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *Page) Layout(gtx C) D {
	var body func(gtx C) D

	switch {
	case !pg.WL.MultiWallet.IsConnectedToDecredNetwork():
		body = func(gtx C) D {
			return pg.pageSections(gtx, func(gtx C) D {
				return pg.welcomeLayout(gtx, pg.sync)
			})
		}
	case !pg.Dexc().Initialized():
		body = func(gtx C) D {
			return pg.pageSections(gtx, func(gtx C) D {
				return pg.welcomeLayout(gtx, pg.initialize)
			})
		}
	case !pg.Dexc().IsLoggedIn():
		body = func(gtx C) D {
			return pg.pageSections(gtx, func(gtx C) D {
				return pg.welcomeLayout(gtx, pg.login)
			})
		}
	case len(pg.Dexc().DEXServers()) == 0:
		body = func(gtx C) D {
			return pg.pageSections(gtx, func(gtx C) D {
				return pg.welcomeLayout(gtx, pg.addDex)
			})
		}
	default:
		body = func(gtx C) D {
			dex := pg.dex()
			if !dex.Connected {
				return D{}
			}

			if dex.PendingFee == nil && pg.miniTradeFormWdg != nil {
				return pg.pageSections(gtx, pg.miniTradeFormWdg.layout)
			}

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

func (pg *Page) welcomeLayout(gtx C, button decredmaterial.Button) D {
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				description := "Trade crypto peer-to-peer."
				return layout.Inset{Bottom: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, pg.Theme.H5(description).Layout)
				})
			}),
			layout.Rigid(button.Layout),
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
	if !dex.Connected {
		// TODO: render error or UI to connect to dex
		return pg.Theme.Label(values.TextSize14, fmt.Sprintf("%s not connected yet", dex.Host)).Layout(gtx)
	}

	if dex.PendingFee == nil {
		// TODO: render trade UI
		return pg.Theme.Label(values.TextSize14, "Registration fee payment successful!").Layout(gtx)
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
		// Todo: this for test purposes, will remove later
		layout.Rigid(pg.addWallet.Layout),
	)
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *Page) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	go pg.readNotifications()
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *Page) OnNavigatedFrom() {
	pg.ctxCancel()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *Page) HandleUserInteractions() {
	if pg.sync.Button.Clicked() {
		err := pg.WL.MultiWallet.SpvSync()
		if err != nil {
			pg.Toast.NotifyError(err.Error())
		}
	}

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
					// Check if there is no dex registered, show modal to register one
					if len(pg.Dexc().DEXServers()) == 0 {
						newAddDexModal(pg.Load).Show()
					}
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
					// Check if there is no dex registered, show modal to register one
					if len(pg.Dexc().DEXServers()) == 0 {
						newAddDexModal(pg.Load).Show()
					}
				}()
				return false
			}).Show()
	}

	if pg.addDex.Button.Clicked() {
		newAddDexModal(pg.Load).Show()
	}

	if pg.miniTradeFormWdg != nil {
		pg.miniTradeFormWdg.handle()
	}

	dex := pg.dex()
	if dex != nil &&
		dex.Connected &&
		dex.PendingFee == nil &&
		pg.miniTradeFormWdg == nil {
		pg.miniTradeFormWdg = newMiniTradeFormWidget(pg.Load)
	}
}

// readNotifications reads from the Core notification channel.
func (pg *Page) readNotifications() {
	ch := pg.Dexc().Core().NotificationFeed()
	for {
		select {
		case n := <-ch:
			if n.Type() == core.NoteTypeFeePayment {
				pg.RefreshWindow()
			}

			if n.Severity() > db.Success {
				pg.Toast.NotifyError(n.Details())
			}

		case <-pg.ctx.Done():
			return
		}
	}
}
