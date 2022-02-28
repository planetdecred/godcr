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

type Page struct {
	*load.Load
	ctx       context.Context
	ctxCancel context.CancelFunc

	loginBtn          decredmaterial.Button
	initializeBtn     decredmaterial.Button
	addDexBtn         decredmaterial.Button
	syncBtn           decredmaterial.Button
	materialLoader    material.LoaderStyle
	walletSettingsBtn *decredmaterial.Clickable
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
		loginBtn:          l.Theme.Button(strLogin),
		initializeBtn:     l.Theme.Button(strStartUseDex),
		addDexBtn:         l.Theme.Button(strAddADex),
		syncBtn:           l.Theme.Button(strStartSyncToUse),
		materialLoader:    material.Loader(material.NewTheme(gofont.Collection())),
		walletSettingsBtn: clickable(),
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
		case pg.dexServer() == nil:
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.headerLayout()),
				layout.Rigid(func(gtx C) D {
					return pg.pageSections(gtx, pg.welcomeLayout(&pg.addDexBtn))
				}),
			)
		default:
			d := pg.dexServer()
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.headerLayout()),
				layout.Rigid(func(gtx C) D {
					if !d.Connected {
						return pg.pageSections(gtx,
							pg.Theme.Label(values.TextSize16, fmt.Sprintf(nStrConnHostError, d.Host)).Layout)
					}
					if d.PendingFee != nil {
						return pg.pageSections(gtx, pg.registrationStatusLayout())
					}

					// TODO: remove this and render trade UI
					return pg.pageSections(gtx, pg.Theme.Label(values.TextSize14, "Registration fee payment successful!").Layout)
				}),
			)
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
					if pg.Dexc().Core() == nil {
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

func (pg *Page) headerLayout() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Inset{
			Bottom: values.MarginPadding15,
		}.Layout(gtx, func(gtx C) D {
			walletIc := pg.Theme.Icons.WalletIcon
			walletIc.Scale = .3
			return layout.E.Layout(gtx, func(gtx C) D {
				return widget.Border{
					Color:        pg.Theme.Color.Gray2,
					CornerRadius: values.MarginPadding0,
					Width:        values.MarginPadding1,
				}.Layout(gtx, func(gtx C) D {
					return pg.walletSettingsBtn.Layout(gtx, func(gtx C) D {
						return layout.Inset{
							Top:    values.MarginPadding4,
							Bottom: values.MarginPadding4,
							Left:   values.MarginPadding8,
							Right:  values.MarginPadding8,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(pg.Theme.Label(values.MarginPadding12, strWalletSetting).Layout),
								layout.Rigid(walletIc.Layout),
							)
						})
					})
				})
			})
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
		return
	}
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
					pm.Dismiss()
					// Check if there is no dex registered, show modal to register one
					if len(pg.Dexc().DEXServers()) == 0 {
						newAddDexModal(pg.Load).WithAppPassword(password).
							OnDexAdded(func() {
								pg.RefreshWindow()
							}).Show()
						return
					}
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
					// Check if there is no dex registered, show modal to register one
					if len(pg.Dexc().DEXServers()) == 0 {
						newAddDexModal(pg.Load).
							WithAppPassword(password).
							OnDexAdded(func() {
								pg.RefreshWindow()
							}).Show()
						return
					}
				}()
				return false
			}).Show()
	}

	if pg.addDexBtn.Button.Clicked() {
		newAddDexModal(pg.Load).OnDexAdded(func() {
			pg.RefreshWindow()
		}).Show()
	}

	if pg.walletSettingsBtn.Clicked() {
		pg.ChangeFragment(NewDexWalletsPage(pg.Load))
	}
}

func (pg *Page) startDexClient() {
	_, err := pg.WL.MultiWallet.StartDexClient()
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
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

// dexServer return first Dex
func (pg *Page) dexServer() *core.Exchange {
	exchanges := sortServers(pg.Dexc().DEXServers())
	if len(exchanges) == 0 {
		return nil
	}
	return exchanges[0]
}
