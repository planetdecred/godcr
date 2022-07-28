package dexclient

import (
	"context"
	"fmt"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/client/db"
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/app"
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
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	ctx            context.Context
	ctxCancel      context.CancelFunc
	addDexBtn      decredmaterial.Button
	syncBtn        decredmaterial.Button
	materialLoader material.LoaderStyle
}

func NewMarketPage(l *load.Load) *Page {
	pg := &Page{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(MarketPageID),
		addDexBtn:        l.Theme.Button(strAddADex),
		syncBtn:          l.Theme.Button(strStartSyncToUse),
		materialLoader:   material.Loader(l.Theme.Base),
	}

	return pg
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *Page) Layout(gtx C) D {
	body := func(gtx C) D {
		switch {
		case !pg.WL.MultiWallet.IsConnectedToDecredNetwork():
			return pg.pageSections(gtx, pg.welcomeLayout(&pg.syncBtn))
		case pg.isLoadingDexClient(): // Need start DEX client
			return pg.pageSections(gtx, pg.welcomeLayout(nil))
		case pg.dexServer() == nil:
			return pg.pageSections(gtx, pg.welcomeLayout(&pg.addDexBtn))
		default:
			d := pg.dexServer()
			if !d.Connected {
				return pg.pageSections(gtx,
					pg.Theme.Label(values.TextSize16, fmt.Sprintf(nStrConnHostError, d.Host)).Layout)
			}
			if d.PendingFee != nil {
				return pg.pageSections(gtx, pg.registrationStatusLayout())
			}

			// TODO: remove this and render trade UI
			return pg.pageSections(gtx, pg.Theme.Label(values.TextSize14, "Registration fee payment successful!").Layout)
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

func (pg *Page) welcomeLayout(button *decredmaterial.Button) layout.Widget {
	return func(gtx C) D {
		return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					description := "Trade crypto peer-to-peer."
					return layout.Inset{Bottom: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, pg.Theme.H5(description).Layout)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if pg.isLoadingDexClient() {
						return layout.Center.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = 50
							return pg.materialLoader.Layout(gtx)
						})
					}
					if button == nil {
						return D{}
					}
					return button.Layout(gtx)
				}),
			)
		})
	}
}

func (pg *Page) registrationStatusLayout() layout.Widget {
	return func(gtx C) D {
		txtLabel := func(txt string) layout.Widget {
			return pg.Theme.Label(values.TextSize14, txt).Layout
		}
		d := pg.dexServer()
		reqConfirms, currentConfs := d.Fee.Confs, d.PendingFee.Confs
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(txtLabel(strWaitingConfirms)),
			layout.Rigid(txtLabel(fmt.Sprintf(nStrConfirmationsStatus, d.Host, reqConfirms))),
			layout.Rigid(txtLabel(fmt.Sprintf("%d/%d", currentConfs, reqConfirms))),
		)
	}
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *Page) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	if pg.Dexc().Core() == nil {
		go pg.startDexClient()
	} else {
		go pg.readNotifications()
	}
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
	if pg.syncBtn.Button.Clicked() {
		err := pg.WL.MultiWallet.SpvSync()
		if err != nil {
			pg.Toast.NotifyError(err.Error())
		}
	}

	if pg.addDexBtn.Button.Clicked() {
		newAddDexModal := NewAddDexModal(pg.Load).OnDexAdded(func() {
			pg.ParentWindow().Reload()
		})
		pg.ParentWindow().ShowModal(newAddDexModal)
	}
}

// isLoadingDexClient check for Dexc start, initialized, loggedin status,
// since Dex client UI not required for app password, IsInitialized and IsLoggedIn should be done at dcrlibwallet.
func (pg *Page) isLoadingDexClient() bool {
	return pg.Dexc().Core() == nil || !pg.Dexc().Core().IsInitialized() || !pg.Dexc().IsLoggedIn()
}

// startDexClient do start DEX client,
// initialize and login to DEX,
// since Dex client UI not required for app password, initialize and login should be done at dcrlibwallet.
func (pg *Page) startDexClient() {
	_, err := pg.WL.MultiWallet.StartDexClient()
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
	}

	// TODO: move to dcrlibwallet sine bypass Dex password by DEXClientPass
	if !pg.Dexc().Initialized() {
		err = pg.Dexc().InitializeWithPassword([]byte(DEXClientPass))
		if err != nil {
			pg.Toast.NotifyError(err.Error())
			return
		}
	}

	if !pg.Dexc().IsLoggedIn() {
		err := pg.Dexc().Login([]byte(DEXClientPass))
		if err != nil {
			pg.Toast.NotifyError(err.Error())
			return
		}
	}

	pg.readNotifications()
}

// readNotifications reads from the Core notification channel.
func (pg *Page) readNotifications() {
	ch := pg.Dexc().Core().NotificationFeed()
	for {
		select {
		case n := <-ch:
			if n.Type() == core.NoteTypeFeePayment || n.Type() == core.NoteTypeConnEvent {
				pg.ParentWindow().Reload()
			}

			if n.Severity() > db.Success {
				pg.Toast.NotifyError(n.Details())
			}

		case <-pg.ctx.Done():
			return
		}
	}
}

// dexServer return first Dex
func (pg *Page) dexServer() *core.Exchange {
	exchanges := sortServers(pg.Dexc().DEXServers())
	if len(exchanges) == 0 {
		return nil
	}
	return exchanges[0]
}
