package dexclient

import (
	"context"
	"errors"
	"fmt"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/client/db"
	"gioui.org/layout"
	"gioui.org/widget"

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
	dex              *core.Exchange

	loginBtn          decredmaterial.Button
	initializeBtn     decredmaterial.Button
	addDexBtn         decredmaterial.Button
	syncBtn           decredmaterial.Button
	walletSettingsBtn *decredmaterial.Clickable
	ordersHistoryBtn  *decredmaterial.Clickable
	dexSettingsBtn    *decredmaterial.Clickable
	dexSelectBtn      *decredmaterial.Clickable
}

var defaultMarketID = "dcr_btc"

func NewMarketPage(l *load.Load) *Page {
	clickable := func() *decredmaterial.Clickable {
		cl := l.Theme.NewClickable(true)
		cl.ChangeStyle(&values.ClickableStyle{HoverColor: l.Theme.Color.Surface})
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}

	pg := &Page{
		Load:              l,
		loginBtn:          l.Theme.Button("Login"),
		initializeBtn:     l.Theme.Button("Start using now"),
		addDexBtn:         l.Theme.Button("Add a dex"),
		syncBtn:           l.Theme.Button("Start sync to continue"),
		ordersHistoryBtn:  clickable(),
		walletSettingsBtn: clickable(),
		dexSettingsBtn:    clickable(),
		dexSelectBtn:      clickable(),
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
	container := func(gtx C) D {
		switch {
		case !pg.WL.MultiWallet.IsConnectedToDecredNetwork():
			return pg.pageSections(gtx, pg.welcomeLayout(pg.syncBtn))
		case !pg.Dexc().Initialized():
			return pg.pageSections(gtx, pg.welcomeLayout(pg.initializeBtn))
		case !pg.Dexc().IsLoggedIn():
			return pg.pageSections(gtx, pg.welcomeLayout(pg.loginBtn))
		case pg.dex == nil:
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.headerLayout()),
				layout.Rigid(func(gtx C) D {
					return pg.pageSections(gtx, pg.welcomeLayout(pg.addDexBtn))
				}),
			)
		default:
			if !pg.dex.Connected {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(pg.headerLayout()),
					layout.Rigid(func(gtx C) D {
						return pg.pageSections(gtx,
							pg.Theme.Label(values.TextSize16, fmt.Sprintf("Connection to dex server %s failed. You can close app and try again later or wait for it to reconnect", pg.dex.Host)).Layout)
					}),
				)
			}

			if pg.dex.PendingFee != nil {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(pg.headerLayout()),
					layout.Rigid(func(gtx C) D {
						return pg.pageSections(gtx, pg.registrationStatusLayout())
					}),
				)
			}

			pg.initMiniTradeForm()

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.headerLayout()),
				layout.Rigid(func(gtx C) D {
					return pg.pageSections(gtx, pg.miniTradeFormWdg.layout)
				}),
			)
		}
	}

	return components.UniformPadding(gtx, container)
}

func (pg *Page) headerLayout() layout.Widget {
	return func(gtx C) D {
		bottom := layout.Inset{Bottom: values.MarginPadding15}
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
		return bottom.Layout(gtx, func(gtx C) D {
			dexIc := pg.Load.Icons.DexIcon
			orderHistoryIc := pg.Load.Icons.TimerIcon
			walletIc := pg.Load.Icons.WalletIcon
			dexSettingIc := pg.Load.Icons.SettingsIcon
			dexIc.Scale, orderHistoryIc.Scale, walletIc.Scale, dexSettingIc.Scale = .1, .5, .3, .3

			return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if pg.dex == nil {
						return D{}
					}
					return layout.W.Layout(gtx, btn(pg.dexSelectBtn, pg.dex.Host, dexIc))
				}),
				layout.Rigid(func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								if pg.dex == nil {
									return D{}
								}
								return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, btn(pg.ordersHistoryBtn, "Order History", orderHistoryIc))
							}),
							layout.Rigid(btn(pg.walletSettingsBtn, "Wallets Settings", walletIc)),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, btn(pg.dexSettingsBtn, "Dex Settings", dexSettingIc))
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

func (pg *Page) welcomeLayout(button decredmaterial.Button) layout.Widget {
	return func(gtx C) D {
		return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					description := "Trade crypto peer-to-peer.\nNo trading fees. No KYC."
					return layout.Center.Layout(gtx, pg.Theme.H5(description).Layout)
				}),
				layout.Rigid(button.Layout),
			)
		})
	}
}

func (pg *Page) registrationStatusLayout() layout.Widget {
	return func(gtx C) D {
		reqConfirms, currentConfs := pg.dex.Fee.Confs, pg.dex.PendingFee.Confs
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(pg.Theme.Label(values.TextSize14, "Waiting for confirmations...").Layout),
			layout.Rigid(func(gtx C) D {
				description := "Trade crypto peer-to-peer."
				return layout.Inset{Bottom: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, pg.Theme.H5(description).Layout)
				})
			}),
		)
	}
}

func (pg *Page) initMiniTradeForm() {
	mkt := pg.dex.Markets[defaultMarketID]
	if pg.miniTradeFormWdg == nil ||
		pg.miniTradeFormWdg.host != pg.dex.Host ||
		pg.miniTradeFormWdg.mkt.Name != mkt.Name {
		pg.miniTradeFormWdg = newMiniTradeFormWidget(pg.Load, pg.dex.Host, mkt)
	}

	if dex.PendingFee == nil {
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
	)
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *Page) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	go pg.readNotifications()

	if err := pg.initOrSetDefaultDex(); err == nil {
		mkt := pg.dex.Markets[defaultMarketID]
		if pg.Dexc().IsLoggedIn() && mkt != nil {
			go pg.listenerMessages(pg.dex.Host, mkt.BaseID, mkt.QuoteID)
		}
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
	if pg.sync.Button.Clicked() {
		err := pg.WL.MultiWallet.SpvSync()
		if err != nil {
			pg.Toast.NotifyError(err.Error())
		}
	}

	if pg.loginBtn.Button.Clicked() {
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
						return
					}

					if err := pg.initOrSetDefaultDex(); err == nil {
						if mkt := pg.dex.Markets[defaultMarketID]; mkt != nil {
							go pg.listenerMessages(pg.dex.Host, mkt.BaseID, mkt.QuoteID)
						}
					}
				}()
				return false
			}).Show()
	}

	if pg.initializeBtn.Button.Clicked() {
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

	if pg.addDexBtn.Button.Clicked() {
		newAddDexModal(pg.Load).Show()
	}

	if pg.miniTradeFormWdg != nil {
		pg.miniTradeFormWdg.handle(pg.orderBook)
	}

	if pg.walletSettingsBtn.Clicked() {
		pg.ChangeFragment(NewDexWalletsPage(pg.Load))
	}

	if pg.ordersHistoryBtn.Clicked() {
		pg.ChangeFragment(NewOrdersHistoryPage(pg.Load, pg.dex.Host))
	}

	if pg.dexSettingsBtn.Clicked() {
		pg.ChangeFragment(NewDexSettingsPage(pg.Load))
	}

	if pg.dexSelectBtn.Clicked() {
		newSelectorDexModal(pg.Load, pg.dex.Host).
			OnDexSelected(func(dex *core.Exchange) {
				pg.setDex(dex)
				mkt := dex.Markets[defaultMarketID]
				if mkt == nil {
					return
				}
				go pg.listenerMessages(pg.dex.Host, mkt.BaseID, mkt.QuoteID)
			}).Show()
	}
}

// readNotifications reads from the Core notification channel.
func (pg *Page) readNotifications() {
	ch := pg.Dexc().Core().NotificationFeed()
	for {
		select {
		case n := <-ch:
			pg.updateDexState()
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
		<-bookFeed.Next()
		pg.updateDexState()
		pg.getOrderBook(host, baseID, quoteID)
		pg.RefreshWindow()
	}
}

func (pg *Page) getOrderBook(host string, baseID, quoteID uint32) {
	orderBoook, err := pg.Dexc().Core().Book(host, baseID, quoteID)
	if err != nil {
		return
	}

	pg.orderBook = orderBoook
}

func (pg *Page) initOrSetDefaultDex() error {
	valueOut := pg.WL.MultiWallet.ReadStringConfigValueForKey(DexHostConfigKey)
	if valueOut != "" {
		if dex, ok := pg.Dexc().DEXServers()[valueOut]; ok {
			pg.dex = dex
			return nil
		}
		return errors.New("No dex server")
	}
	exchanges := sliceExchanges(pg.Dexc().DEXServers())
	if len(exchanges) == 0 {
		return errors.New("No dex server")
	}
	pg.setDex(exchanges[0])
	return nil
}

func (pg *Page) setDex(dex *core.Exchange) {
	pg.dex = dex
	pg.WL.MultiWallet.SetStringConfigValueForKey(DexHostConfigKey, dex.Host)
}

func (pg *Page) updateDexState() {
	if pg.dex == nil {
		return
	}
	if dex := pg.Dexc().DEXServers()[pg.dex.Host]; dex != nil {
		pg.dex = dex
	}
}
