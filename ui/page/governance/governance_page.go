package governance

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const GovernancePageID = "Governance"

type Page struct {
	*load.Load

	multiWallet *dcrlibwallet.MultiWallet

	tabCategoryList        *decredmaterial.ClickableList
	splashScreenInfoButton decredmaterial.IconButton
	enableGovernanceBtn    decredmaterial.Button

	proposalsPage *ProposalsPage
	consensusPage *ConsensusPage

	selectedCategoryIndex int
	changed               bool
}

var governanceTabTitles = []string{
	values.String(values.StrProposal),
	values.String(values.StrConsensusChange),
}

func NewGovernancePage(l *load.Load) *Page {
	pg := &Page{
		Load:                  l,
		multiWallet:           l.WL.MultiWallet,
		selectedCategoryIndex: -1,
		proposalsPage:         NewProposalsPage(l),
		consensusPage:         NewConsensusPage(l),
		tabCategoryList:       l.Theme.NewClickableList(layout.Horizontal),
	}

	pg.tabCategoryList.IsHoverable = false

	pg.initSplashScreenWidgets()

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *Page) OnNavigatedTo() {
	selectedCategory := pg.selectedCategoryIndex

	if selectedCategory == -1 {
		pg.selectedCategoryIndex = 0
	}

	if pg.selectedCategoryIndex == 1 {
		pg.consensusPage.OnNavigatedTo()
	} else {
		pg.proposalsPage.OnNavigatedTo()
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
	pg.consensusPage.OnNavigatedFrom()
	pg.proposalsPage.OnNavigatedFrom()
}

func (pg *Page) ID() string {
	return GovernancePageID
}

func (pg *Page) HandleUserInteractions() {
	for pg.splashScreenInfoButton.Button.Clicked() {
		pg.showInfoModal()
	}

	for pg.enableGovernanceBtn.Clicked() {
		go pg.consensusPage.FetchAgendas()
		go pg.WL.MultiWallet.Politeia.Sync()
		pg.proposalsPage.isSyncing = pg.multiWallet.Politeia.IsSyncing()
		pg.WL.MultiWallet.SaveUserConfigValue(load.FetchProposalConfigKey, true)
	}

	if clicked, selectedItem := pg.tabCategoryList.ItemClicked(); clicked {
		if pg.selectedCategoryIndex != selectedItem {
			pg.selectedCategoryIndex = selectedItem
			pg.changed = true
		}

		// call selected page OnNavigatedTo() only once
		if pg.changed && pg.selectedCategoryIndex == 0 {
			pg.proposalsPage.OnNavigatedTo()
		} else if pg.changed && pg.selectedCategoryIndex == 1 {
			pg.consensusPage.OnNavigatedTo()
		}
		pg.changed = false
	}

	// handle individual page user interactions
	if pg.selectedCategoryIndex == 0 {
		pg.proposalsPage.HandleUserInteractions()
	} else {
		pg.consensusPage.HandleUserInteractions()
	}
}

func (pg *Page) Layout(gtx C) D {
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *Page) layoutDesktop(gtx layout.Context) layout.Dimensions {
	if pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey, false) {
		return components.UniformPadding(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.layoutPageTopNav),
				layout.Rigid(pg.layoutTabs),
				layout.Rigid(pg.Theme.Separator().Layout),
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return pg.switchTab(gtx, pg.selectedCategoryIndex)
					})
				}),
			)
		})
	}
	return components.UniformPadding(gtx, pg.splashScreenLayout)
}

func (pg *Page) layoutMobile(gtx layout.Context) layout.Dimensions {
	if pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey, false) {
		return components.UniformMobile(gtx, true, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.layoutPageTopNav),
				layout.Rigid(pg.layoutTabs),
				layout.Rigid(pg.Theme.Separator().Layout),
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return pg.switchTab(gtx, pg.selectedCategoryIndex)
					})
				}),
			)
		})
	}
	return components.UniformMobile(gtx, true, pg.splashScreenLayout)
}

func (pg *Page) switchTab(gtx C, selectedCategoryIndex int) D {
	if selectedCategoryIndex == 0 {
		return pg.proposalsPage.Layout(gtx)
	}

	return pg.consensusPage.Layout(gtx)
}

func (pg *Page) layoutTabs(gtx C) D {
	var dims layout.Dimensions

	return layout.Inset{
		Top: values.MarginPadding20,
	}.Layout(gtx, func(gtx C) D {
		return pg.tabCategoryList.Layout(gtx, len(governanceTabTitles), func(gtx C, i int) D {
			return layout.Stack{Alignment: layout.S}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					return layout.Inset{
						Right:  values.MarginPadding24,
						Bottom: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							lbl := pg.Theme.Label(values.TextSize16, governanceTabTitles[i])
							lbl.Color = pg.Theme.Color.GrayText1
							if pg.selectedCategoryIndex == i {
								lbl.Color = pg.Theme.Color.Primary
								dims = lbl.Layout(gtx)
							}

							return lbl.Layout(gtx)
						})
					})
				}),
				layout.Stacked(func(gtx C) D {
					if pg.selectedCategoryIndex != i {
						return D{}
					}

					tabHeight := gtx.Dp(values.MarginPadding2)
					tabRect := image.Rect(0, 0, dims.Size.X, tabHeight)

					return layout.Inset{
						Left: values.MarginPaddingMinus22,
					}.Layout(gtx, func(gtx C) D {
						paint.FillShape(gtx.Ops, pg.Theme.Color.Primary, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: dims.Size.X, Y: tabHeight},
						}
					})
				}),
			)
		})
	})
}

func (pg *Page) layoutPageTopNav(gtx C) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(pg.Theme.Icons.GovernanceActiveIcon.Layout24dp),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Left: values.MarginPadding20,
			}.Layout(gtx, func(gtx C) D {
				txt := pg.Theme.Label(values.TextSize20, values.String(values.StrGovernance))
				txt.Font.Weight = text.SemiBold
				return txt.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return D{}
				//TODO: governance syncing functionality.
				//TODO: Split wallet sync from governance
			})
		}),
	)
}
