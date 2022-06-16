package page

import (
	// "context"
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"

	// "github.com/decred/dcrd/dcrutil/v4"
	// "github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	// "github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	// "github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWalletID = "create_wallet"

type CreateWallet struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	listLock        sync.Mutex
	scrollContainer *widget.List

	wallectSelected func()
}

func NewCreateWallet(l *load.Load) *CreateWallet {
	pg := &CreateWallet{
		GenericPageModal: app.NewGenericPageModal(CreateWalletID),
		scrollContainer: &widget.List{
			List: layout.List{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			},
		},

		Load: l,
	}

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *CreateWallet) OnNavigatedTo() {}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *CreateWallet) HandleUserInteractions() {
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *CreateWallet) OnNavigatedFrom() {}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *CreateWallet) Layout(gtx C) D {
	pageContent := []func(gtx C) D{
		pg.Theme.Label(values.TextSize20, values.String(values.StrSelectWalletToOpen)).Layout,
	}

	return decredmaterial.LinearLayout{
		Width:      decredmaterial.MatchParent,
		Height:     decredmaterial.MatchParent,
		Background: pg.Theme.Color.Success,
		Alignment:  layout.Middle,
		Direction:  layout.Center,
	}.Layout2(gtx, func(gtx C) D {
		return components.UniformPadding(gtx, func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:      gtx.Dp(values.MarginPadding550),
				Height:     decredmaterial.MatchParent,
				Background: pg.Theme.Color.Primary,
				Alignment:  layout.Middle,
			}.Layout2(gtx, func(gtx C) D {
				list := &layout.List{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}

				return list.Layout(gtx, len(pageContent), func(gtx C, i int) D {
					return layout.Inset{Top: values.MarginPadding26}.Layout(gtx, func(gtx C) D {
						return pageContent[i](gtx)
					})
				})
			})
		})
	})
}
