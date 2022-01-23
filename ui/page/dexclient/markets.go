package dexclient

import (
	"context"
	"fmt"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/client/db"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

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

const DexHostConfigKey = "dex_host"

type Page struct {
	*load.Load
	ctx              context.Context
	ctxCancel        context.CancelFunc
	miniTradeFormWdg *miniTradeFormWidget
	orderBook        *core.OrderBook
	dexServer        *core.Exchange
	market           *core.Market
	shouldStartDex   bool
	materialLoader   material.LoaderStyle

	loginBtn          decredmaterial.Button
	initializeBtn     decredmaterial.Button
	addDexBtn         decredmaterial.Button
	syncBtn           decredmaterial.Button
	walletSettingsBtn *decredmaterial.Clickable
	ordersHistoryBtn  *decredmaterial.Clickable
	dexSettingsBtn    *decredmaterial.Clickable
	dexSelectBtn      *decredmaterial.Clickable
}

func NewMarketPage(l *load.Load) *Page {
	clickable := func() *decredmaterial.Clickable {
		cl := l.Theme.NewClickable(true)
		cl.ChangeStyle(&values.ClickableStyle{HoverColor: l.Theme.Color.Surface})
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}

	pg := &Page{
		Load:              l,
		shouldStartDex:    l.Dexc().Core() == nil,
		materialLoader:    material.Loader(material.NewTheme(gofont.Collection())),
		loginBtn:          l.Theme.Button(strLogin),
		initializeBtn:     l.Theme.Button(strStartUseDex),
		addDexBtn:         l.Theme.Button(strAddADex),
		syncBtn:           l.Theme.Button(strStartSyncToUse),
		ordersHistoryBtn:  clickable(),
		walletSettingsBtn: clickable(),
		dexSettingsBtn:    clickable(),
		dexSelectBtn:      clickable(),
		miniTradeFormWdg:  newMiniTradeFormWidget(l),
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
	body := func(gtx C) D {
		switch {
		case pg.Dexc().Core() == nil: // Need start DEX client
			return pg.pageSections(gtx, pg.welcomeLayout(nil))
		case !pg.WL.MultiWallet.IsConnectedToDecredNetwork():
			return pg.pageSections(gtx, pg.welcomeLayout(&pg.syncBtn))
		case !pg.Dexc().Initialized():
			return pg.pageSections(gtx, pg.welcomeLayout(&pg.initializeBtn))
		case !pg.Dexc().IsLoggedIn():
			return pg.pageSections(gtx, pg.welcomeLayout(&pg.loginBtn))
		case pg.dexServer == nil:
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.headerLayout()),
				layout.Rigid(func(gtx C) D {
					return pg.pageSections(gtx, pg.welcomeLayout(&pg.addDexBtn))
				}),
			)
		default:
			if !pg.dexServer.Connected {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(pg.headerLayout()),
					layout.Rigid(func(gtx C) D {
						return pg.pageSections(gtx,
							pg.Theme.Label(values.TextSize16, fmt.Sprintf(nStrConnHostError, pg.dexServer.Host)).Layout)
					}),
				)
			}

			if pg.dexServer.PendingFee != nil {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(pg.headerLayout()),
					layout.Rigid(func(gtx C) D {
						return pg.pageSections(gtx, pg.registrationStatusLayout())
					}),
				)
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.headerLayout()),
				layout.Rigid(func(gtx C) D {
					return pg.pageSections(gtx, pg.miniTradeFormWdg.setHostAndMarket(pg.dexServer.Host, pg.market).layout)
				}),
			)
		}
	}

	return components.UniformPadding(gtx, body)
}

func (pg *Page) headerLayout() layout.Widget {
	return func(gtx C) D {
		btn := func(btn *decredmaterial.Clickable, textBtn string, ic *decredmaterial.Image) layout.Widget {
			return func(gtx C) D {
				return widget.Border{
					Color:        pg.Theme.Color.Gray2,
					CornerRadius: values.MarginPadding0,
					Width:        values.MarginPadding1,
				}.Layout(gtx, func(gtx C) D {
					return btn.Layout(gtx, func(gtx C) D {
						return layout.Inset{
							Top:    values.MarginPadding4,
							Bottom: values.MarginPadding4,
							Left:   values.MarginPadding8,
							Right:  values.MarginPadding8,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(pg.Theme.Label(values.MarginPadding12, textBtn).Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, ic.Layout)
								}),
							)
						})
					})
				})
			}
		}
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
			dexIc := pg.Icons.DexIcon
			orderHistoryIc := pg.Icons.TimerIcon
			walletIc := pg.Icons.WalletIcon
			dexSettingIc := pg.Icons.SettingsIcon
			dexIc.Scale, orderHistoryIc.Scale, walletIc.Scale, dexSettingIc.Scale = .1, .5, .3, .3

			return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if pg.dexServer == nil {
						return D{}
					}
					return layout.W.Layout(gtx, btn(pg.dexSelectBtn, pg.dexServer.Host, dexIc))
				}),
				layout.Rigid(func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								if pg.dexServer == nil {
									return D{}
								}
								return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, btn(pg.ordersHistoryBtn, strOrderHistory, orderHistoryIc))
							}),
							layout.Rigid(btn(pg.walletSettingsBtn, strWalletSetting, walletIc)),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, btn(pg.dexSettingsBtn, strDexSetting, dexSettingIc))
							}),
						)
					})
				}),
			)
		})
	}
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
					return layout.Center.Layout(gtx, pg.Theme.H5("Trade crypto peer-to-peer.\n").Layout)
				}),
				layout.Rigid(func(gtx C) D {
					if pg.shouldStartDex {
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
		reqConfirms, currentConfs := pg.dexServer.Fee.Confs, pg.dexServer.PendingFee.Confs
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(txtLabel(strWaitingConfirms)),
			layout.Rigid(txtLabel(fmt.Sprintf(nStrConfirmationsStatus, pg.dexServer.Host, reqConfirms))),
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
		return
	}

	go pg.readNotifications()
	if pg.initDexServer() && pg.initMarket() && pg.Dexc().IsLoggedIn() {
		go pg.listenerMessages(pg.dexServer.Host, pg.market.BaseID, pg.market.QuoteID)
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

func (pg *Page) startDexClient() {
	_, err := pg.WL.MultiWallet.StartDexClient()
	pg.shouldStartDex = false
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
	}

	pg.readNotifications()
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

	if pg.loginBtn.Button.Clicked() {
		modal.NewPasswordModal(pg.Load).
			Title(strLogin).
			Hint(strAppPassword).
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton(strLogin, func(password string, pm *modal.PasswordModal) bool {
				go func() {
					err := pg.Dexc().Login([]byte(password))
					if err != nil {
						pm.SetError(err.Error())
						pm.SetLoading(false)
						return
					}
					// Check if there is no dexServer registered, show modal to register one
					if len(pg.Dexc().DEXServers()) == 0 {
						pm.Dismiss()
						newAddDexModal(pg.Load).WithAppPassword(password).
							OnDexAdded(func(dexServer *core.Exchange) {
								pg.dexCreated(dexServer)
							}).Show()
						return
					}

					if pg.initDexServer() && pg.initMarket() {
						go pg.listenerMessages(pg.dexServer.Host, pg.market.BaseID, pg.market.QuoteID)
					}
					pm.Dismiss()
				}()
				return false
			}).Show()
	}

	if pg.initializeBtn.Button.Clicked() {
		modal.NewCreatePasswordModal(pg.Load).
			Title(strSetAppPassword).
			SetDescription(strInitDexPasswordDesc).
			EnableName(false).
			PasswordHint(strAppPassword).
			ConfirmPasswordHint(strConfirmPassword).
			PasswordCreated(func(_, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					err := pg.Dexc().InitializeWithPassword([]byte(password))
					if err != nil {
						m.SetError(err.Error())
						m.SetLoading(false)
						return
					}
					pg.Toast.Notify(strSuccessful)
					m.Dismiss()
					// Check if there is no dexServer registered, show modal to register one
					if len(pg.Dexc().DEXServers()) == 0 {
						newAddDexModal(pg.Load).WithAppPassword(password).
							OnDexAdded(func(dexServer *core.Exchange) {
								pg.dexCreated(dexServer)
							}).Show()
					}
				}()
				return false
			}).Show()
	}

	if pg.addDexBtn.Button.Clicked() {
		newAddDexModal(pg.Load).OnDexAdded(func(dexServer *core.Exchange) {
			pg.dexCreated(dexServer)
		}).Show()
	}

	pg.miniTradeFormWdg.handle(pg.orderBook)

	if pg.walletSettingsBtn.Clicked() {
		pg.ChangeFragment(NewDexWalletsPage(pg.Load))
	}

	if pg.ordersHistoryBtn.Clicked() {
		pg.ChangeFragment(NewOrdersHistoryPage(pg.Load, pg.dexServer.Host))
	}

	if pg.dexSettingsBtn.Clicked() {
		pg.ChangeFragment(NewDexSettingsPage(pg.Load))
	}

	if pg.dexSelectBtn.Clicked() {
		newSelectorDexModal(pg.Load, pg.dexServer.Host).
			OnDexSelected(func(dexServer *core.Exchange) {
				pg.selectDexServer(dexServer)
				if pg.initMarket() {
					go pg.listenerMessages(pg.dexServer.Host, pg.market.BaseID, pg.market.QuoteID)
				}
			}).Show()
	}
}

// readNotifications reads from the Core notification channel.
func (pg *Page) readNotifications() {
	ch := pg.Dexc().Core().NotificationFeed()
	for {
		select {
		case n := <-ch:
			pg.updateDexMarketState()
			if n.Type() == core.NoteTypeFeePayment || n.Type() == core.NoteTypeConnEvent {
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

func (pg *Page) listenerMessages(host string, baseID, quoteID uint32) {
	bookFeed, err := pg.Dexc().Core().SyncBook(host, baseID, quoteID)
	if err != nil {
		return
	}

	for {
		select {
		case <-bookFeed.Next():
			pg.updateDexMarketState()
			pg.getOrderBook(host, baseID, quoteID)
			pg.RefreshWindow()
		case <-pg.ctx.Done():
			return
		}
	}
}

func (pg *Page) getOrderBook(host string, baseID, quoteID uint32) {
	orderBoook, err := pg.Dexc().Core().Book(host, baseID, quoteID)
	if err != nil {
		return
	}
	pg.orderBook = orderBoook
}

// initDexServer initialize Page's dexServer, check for value of DexHostConfigKey in storage,
// if Dex exist set to Page otherwise choose first Dex on the slice.
func (pg *Page) initDexServer() bool {
	valueOut := pg.WL.MultiWallet.ReadStringConfigValueForKey(DexHostConfigKey)
	if valueOut != "" {
		if dexServer, ok := pg.Dexc().DEXServers()[valueOut]; ok {
			pg.selectDexServer(dexServer)
			return true
		}
	}
	exchanges := sliceExchanges(pg.Dexc().DEXServers())
	if len(exchanges) == 0 {
		pg.selectDexServer(nil)
		return false
	}
	pg.selectDexServer(exchanges[0])
	return true
}

// selectDexServer set Page's Dex and save to storage last choose.
func (pg *Page) selectDexServer(dexServer *core.Exchange) {
	pg.dexServer = dexServer
	value := ""
	if dexServer != nil {
		value = dexServer.Host
	}
	pg.WL.MultiWallet.SetStringConfigValueForKey(DexHostConfigKey, value)
}

// initMarket initialize Page's Market, choose first Market on the slice.
func (pg *Page) initMarket() bool {
	pg.market = nil
	markets := sliceMarkets(pg.dexServer.Markets)
	if len(markets) == 0 {
		return false
	}

	// TODO: select the first Market to use instead check for supported market.
	for _, mkt := range markets {
		if supportedMarket(mkt) {
			pg.selectMarket(mkt)
			return true
		}
	}

	// pg.selectMarket(markets[0])
	return true
}

func (pg *Page) selectMarket(market *core.Market) {
	pg.market = market
}

func (pg *Page) updateDexMarketState() {
	if pg.dexServer == nil {
		return
	}

	dexServer := pg.Dexc().DEXServers()[pg.dexServer.Host]
	if dexServer == nil {
		return
	}

	pg.selectDexServer(dexServer)

	if pg.market == nil {
		return
	}

	market := dexServer.Markets[pg.market.Name]
	if market == nil {
		return
	}

	pg.selectMarket(market)
}

func (pg *Page) dexCreated(dexServer *core.Exchange) {
	pg.selectDexServer(dexServer)
	if pg.initMarket() {
		go pg.listenerMessages(pg.dexServer.Host, pg.market.BaseID, pg.market.QuoteID)
	}
	pg.RefreshWindow()
}
