package privacy

import (
	"context"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
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
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	pageContainer  layout.List
	toPrivacySetup decredmaterial.Button

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton
}

func NewSetupPrivacyPage(l *load.Load) *SetupPrivacyPage {
	pg := &SetupPrivacyPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(SetupPrivacyPageID),
		pageContainer:    layout.List{Axis: layout.Vertical},
		toPrivacySetup:   l.Theme.Button(values.String(values.StrSetupStakeShuffle)),
	}
	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg

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
	return components.UniformPadding(gtx, func(gtx C) D {
		return pg.privacyIntroLayout(gtx)
	})
}

func (pg *SetupPrivacyPage) privacyIntroLayout(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Bottom: values.MarginPadding24,
								}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Left: values.MarginPadding5,
											}.Layout(gtx, pg.Theme.Icons.TransactionFingerprint.Layout48dp)
										}),
										layout.Rigid(pg.Theme.Icons.ArrowForward.Layout24dp),
										layout.Rigid(func(gtx C) D {
											return pg.Theme.Icons.Mixer.LayoutSize(gtx, values.MarginPadding120)
										}),
										layout.Rigid(pg.Theme.Icons.ArrowForward.Layout24dp),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Left: values.MarginPadding5,
											}.Layout(gtx, pg.Theme.Icons.TransactionsIcon.Layout48dp)
										}),
									)
								})
							}),
							layout.Rigid(func(gtx C) D {
								txt := pg.Theme.H6(values.String(values.StrStakeShuffle))
								txt2 := pg.Theme.Body1(values.String(values.StrSetUpPrivacy))

								txt.Alignment, txt2.Alignment = text.Middle, text.Middle

								return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(txt.Layout),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, txt2.Layout)
									}),
								)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.UniformInset(values.MarginPadding30).Layout(gtx, pg.toPrivacySetup.Layout)
				}),
			)
		})
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SetupPrivacyPage) HandleUserInteractions() {
	if pg.toPrivacySetup.Clicked() {
		accounts, err := pg.WL.SelectedWallet.Wallet.GetAccountsRaw()
		if err != nil {
			log.Error(err)
		}

		walCount := accounts.Count
		// Filter out imported account and default account.
		for _, v := range accounts.Acc {
			if v.Number == dcrlibwallet.ImportedAccountNumber || v.Number == dcrlibwallet.DefaultAccountNum {
				walCount--
			}
		}

		if walCount <= 1 {
			go showModalSetupMixerInfo(&sharedModalConfig{
				Load:          pg.Load,
				window:        pg.ParentWindow(),
				pageNavigator: pg.ParentNavigator(),
				checkBox:      pg.Theme.CheckBox(new(widget.Bool), "Automatically move funds from default to unmixed account"),
			})
		} else {
			pg.ParentNavigator().Display(NewSetupMixerAccountsPage(pg.Load))
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
func (pg *SetupPrivacyPage) OnNavigatedFrom() {
	pg.ctxCancel()
}
