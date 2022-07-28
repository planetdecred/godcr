package page

import (
	"context"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/dexclient"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletDexServerSelectorID = "wallet_dex_server_selector"

type WalletDexServerSelector struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	scrollContainer   *widget.List
	shadowBox         *decredmaterial.Shadow
	walletSelector    *components.WalletSelector
	dexServerSelector *components.DexServerSelector
	addWalClickable   *decredmaterial.Clickable
	addDexClickable   *decredmaterial.Clickable
	settings          *decredmaterial.Clickable
}

func NewWalletDexServerSelector(l *load.Load, onWalletSelected func(), onDexServerSelected func(server string)) *WalletDexServerSelector {
	pg := &WalletDexServerSelector{
		GenericPageModal: app.NewGenericPageModal(WalletDexServerSelectorID),
		scrollContainer: &widget.List{
			List: layout.List{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			},
		},
		Load:      l,
		shadowBox: l.Theme.Shadow(),

		walletSelector:    components.NewWalletSelector(l, onWalletSelected),
		dexServerSelector: components.NewDexServerSelector(l, onDexServerSelected),
	}

	pg.addWalClickable = l.Theme.NewClickable(false)
	pg.addWalClickable.Radius = decredmaterial.Radius(14)

	pg.addDexClickable = l.Theme.NewClickable(false)
	pg.addDexClickable.Radius = decredmaterial.Radius(14)

	pg.settings = l.Theme.NewClickable(false)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *WalletDexServerSelector) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.walletSelector.Expose(pg.ctx)
	pg.dexServerSelector.Expose()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *WalletDexServerSelector) HandleUserInteractions() {
	pg.walletSelector.HandleUserInteractions()
	pg.dexServerSelector.HandleUserInteractions()

	if pg.addWalClickable.Clicked() {
		pg.ParentNavigator().Display(NewCreateWallet(pg.Load))
	}

	if pg.addDexClickable.Clicked() {
		dm := dexclient.NewAddDexModal(pg.Load)
		dm.OnDexAdded(func() {
			// TODO: go to the trade form
			log.Info("TODO: go to the trade form")
		})
		pg.ParentWindow().ShowModal(dm)
	}

	if pg.settings.Clicked() {
		pg.ParentNavigator().Display(NewSettingsPage(pg.Load))
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *WalletDexServerSelector) OnNavigatedFrom() {
	pg.ctxCancel()
}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *WalletDexServerSelector) Layout(gtx C) D {
	pg.SetCurrentAppWidth(gtx.Constraints.Max.X)
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *WalletDexServerSelector) layoutDesktop(gtx C) D {
	return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.pageHeaderLayout),
			layout.Rigid(func(gtx C) D {
				return pg.pageContentLayout(gtx)
			}),
		)
	})
}

func (pg *WalletDexServerSelector) layoutMobile(gtx C) D {
	return components.UniformMobile(gtx, false, false, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.pageHeaderLayout),
			layout.Rigid(pg.pageContentLayout),
		)
	})
}

func (pg *WalletDexServerSelector) pageHeaderLayout(gtx C) D {
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Right:  values.MarginPadding15,
					Bottom: values.MarginPadding30,
				}.Layout(gtx, func(gtx C) D {
					return pg.settings.Layout(gtx, pg.Theme.Icons.SettingsIcon.Layout24dp)
				})
			})
		}),
	)
}

func (pg *WalletDexServerSelector) sectionTitle(title string) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, pg.Theme.Label(values.TextSize20, title).Layout)
	}
}

func (pg *WalletDexServerSelector) pageContentLayout(gtx C) D {
	pageContent := []func(gtx C) D{
		pg.sectionTitle(values.String(values.StrSelectWalletToOpen)),
		pg.walletSelector.WalletListLayout,
		pg.layoutAddMoreRowSection(pg.addWalClickable, values.String(values.StrAddWallet), pg.Theme.Icons.NewWalletIcon.Layout24dp),
		pg.sectionTitle(values.String(values.StrSelectWalletToOpen)),
		pg.dexServerSelector.DexServersLayout,
		pg.layoutAddMoreRowSection(pg.addDexClickable, values.String(values.StrAddDexServer), pg.Theme.Icons.DexIcon.Layout16dp),
	}

	return decredmaterial.LinearLayout{
		Width:     decredmaterial.MatchParent,
		Height:    decredmaterial.MatchParent,
		Direction: layout.Center,
	}.Layout2(gtx, func(gtx C) D {
		return decredmaterial.LinearLayout{
			Width:  gtx.Dp(values.MarginPadding550),
			Height: decredmaterial.MatchParent,
			Margin: layout.Inset{
				Bottom: values.MarginPadding30,
			},
		}.Layout2(gtx, func(gtx C) D {
			return pg.Theme.List(pg.scrollContainer).Layout(gtx, len(pageContent), func(gtx C, i int) D {
				return layout.Inset{
					Right: values.MarginPadding48,
				}.Layout(gtx, pageContent[i])
			})
		})
	})
}

func (pg *WalletDexServerSelector) layoutAddMoreRowSection(clk *decredmaterial.Clickable, buttonText string, ic func(gtx C) D) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{
			Left:   values.MarginPadding5,
			Top:    values.MarginPadding10,
			Bottom: values.MarginPadding48,
		}.Layout(gtx, func(gtx C) D {
			pg.shadowBox.SetShadowRadius(14)
			return decredmaterial.LinearLayout{
				Width:      decredmaterial.WrapContent,
				Height:     decredmaterial.WrapContent,
				Padding:    layout.UniformInset(values.MarginPadding12),
				Background: pg.Theme.Color.Surface,
				Clickable:  clk,
				Shadow:     pg.shadowBox,
				Border:     decredmaterial.Border{Radius: clk.Radius},
				Alignment:  layout.Middle,
			}.Layout(gtx,
				layout.Rigid(ic),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left: values.MarginPadding4,
						Top:  values.MarginPadding2,
					}.Layout(gtx, pg.Theme.Body2(buttonText).Layout)
				}),
			)
		})
	}
}
