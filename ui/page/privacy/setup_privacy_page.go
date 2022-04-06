package privacy

import (
	"context"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const SetupPrivacyPageID = "SetupPrivacy"

type (
	C = layout.Context
	D = layout.Dimensions
)

type SetupPrivacyPage struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	wallet         *dcrlibwallet.Wallet
	pageContainer  layout.List
	toPrivacySetup decredmaterial.Button

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton
}

func NewSetupPrivacyPage(l *load.Load, wallet *dcrlibwallet.Wallet) *SetupPrivacyPage {
	pg := &SetupPrivacyPage{
		Load:           l,
		wallet:         wallet,
		pageContainer:  layout.List{Axis: layout.Vertical},
		toPrivacySetup: l.Theme.Button("Set up mixer for this wallet"),
	}
	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg

}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *SetupPrivacyPage) ID() string {
	return SetupPrivacyPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SetupPrivacyPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *SetupPrivacyPage) Layout(gtx layout.Context) layout.Dimensions {
	d := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "StakeShuffle",
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			InfoButton: pg.infoButton,
			Back: func() {
				pg.PopFragment()
			},
			InfoTemplate: modal.PrivacyInfoTemplate,
			Body: func(gtx layout.Context) layout.Dimensions {
				return pg.privacyIntroLayout(gtx)
			},
		}
		return sp.Layout(gtx)
	}
	return components.UniformPadding(gtx, d)
}

func (pg *SetupPrivacyPage) privacyIntroLayout(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Bottom: values.MarginPadding24,
							}.Layout(gtx, func(gtx C) D {
								return pg.Theme.Icons.PrivacySetup.LayoutSize(gtx, values.MarginPadding280)
							})
						}),
						layout.Rigid(func(gtx C) D {
							txt := pg.Theme.H6("How does StakeShuffle enhance your privacy?")
							txt2 := pg.Theme.Body1("StakeShuffle can mix your coins through coinjoin transactions.")
							txt3 := pg.Theme.Body1("Using mixed coins protects you from exposing your financial activities to")
							txt4 := pg.Theme.Body1("the public (e.g. how much you own, who pays you).")
							txt.Alignment, txt2.Alignment, txt3.Alignment, txt4.Alignment = text.Middle, text.Middle, text.Middle, text.Middle

							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(txt.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, txt2.Layout)
								}),
								layout.Rigid(txt3.Layout),
								layout.Rigid(txt4.Layout),
							)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, pg.toPrivacySetup.Layout)
			}),
		)
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SetupPrivacyPage) HandleUserInteractions() {
	if pg.toPrivacySetup.Clicked() {
		pg.ChangeFragment(NewSetupMixerAccountsPage(pg.Load, pg.wallet))
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *SetupPrivacyPage) OnNavigatedFrom() {
	pg.ctxCancel()
}
