package governance

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const GovernancePageID = "Governance"

type Page struct {
	*load.Load
	*app.MasterPage

	multiWallet *dcrlibwallet.MultiWallet

	tabCategoryList        *decredmaterial.ClickableList
	splashScreenInfoButton decredmaterial.IconButton
	enableGovernanceBtn    decredmaterial.Button
}

var governanceTabTitles = []string{
	values.String(values.StrProposal),
	values.String(values.StrConsensusChange),
}

func NewGovernancePage(l *load.Load) *Page {
	pg := &Page{
		Load:            l,
		MasterPage:      app.NewMasterPage(GovernancePageID),
		multiWallet:     l.WL.MultiWallet,
		tabCategoryList: l.Theme.NewClickableList(layout.Horizontal),
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
	if activeTab := pg.CurrentPage(); activeTab != nil {
		activeTab.OnNavigatedTo()
	} else if pg.isGovernanceFeatureEnabled() {
		pg.Display(NewProposalsPage(pg.Load))
	}
}

func (pg *Page) isGovernanceFeatureEnabled() bool {
	return pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey, false)
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *Page) OnNavigatedFrom() {
	if activeTab := pg.CurrentPage(); activeTab != nil {
		activeTab.OnNavigatedFrom()
	}
}

func (pg *Page) HandleUserInteractions() {
	if activeTab := pg.CurrentPage(); activeTab != nil {
		activeTab.HandleUserInteractions()
	}

	for pg.splashScreenInfoButton.Button.Clicked() {
		pg.showInfoModal()
	}

	for pg.enableGovernanceBtn.Clicked() {
		go pg.WL.MultiWallet.Politeia.Sync()
		pg.Display(NewProposalsPage(pg.Load))
		pg.WL.MultiWallet.SaveUserConfigValue(load.FetchProposalConfigKey, true)
	}

	if tabItemClicked, clickedTabIndex := pg.tabCategoryList.ItemClicked(); tabItemClicked {
		if clickedTabIndex == 0 {
			pg.Display(NewProposalsPage(pg.Load)) // Display should do nothing if the page is already displayed.
		} else if clickedTabIndex == 1 {
			pg.Display(NewConsensusPage(pg.Load))
		}
	}
}

func (pg *Page) Layout(gtx C) D {
	if !pg.isGovernanceFeatureEnabled() {
		return components.UniformPadding(gtx, pg.splashScreenLayout)
	}

	return components.UniformPadding(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.layoutPageTopNav),
			layout.Rigid(pg.layoutTabs),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Flexed(1, func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
					return pg.CurrentPage().Layout(gtx)
				})
			}),
		)
	})
}

func (pg *Page) selectedTabIndex() int {
	switch pg.CurrentPageID() {
	case ProposalsPageID:
		return 0
	case ConsensusPageID:
		return 1
	default:
		return -1
	}
}

func (pg *Page) layoutTabs(gtx C) D {
	var selectedTabDims layout.Dimensions

	return layout.Inset{
		Top: values.MarginPadding20,
	}.Layout(gtx, func(gtx C) D {
		return pg.tabCategoryList.Layout(gtx, len(governanceTabTitles), func(gtx C, i int) D {
			isSelectedTab := pg.selectedTabIndex() == i
			return layout.Stack{Alignment: layout.S}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					return layout.Inset{
						Right:  values.MarginPadding24,
						Bottom: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							lbl := pg.Theme.Label(values.TextSize16, governanceTabTitles[i])
							lbl.Color = pg.Theme.Color.GrayText1
							if isSelectedTab {
								lbl.Color = pg.Theme.Color.Primary
								selectedTabDims = lbl.Layout(gtx)
							}

							return lbl.Layout(gtx)
						})
					})
				}),
				layout.Stacked(func(gtx C) D {
					if !isSelectedTab {
						return D{}
					}

					tabHeight := gtx.Dp(values.MarginPadding2)
					tabRect := image.Rect(0, 0, selectedTabDims.Size.X, tabHeight)

					return layout.Inset{
						Left: values.MarginPaddingMinus22,
					}.Layout(gtx, func(gtx C) D {
						paint.FillShape(gtx.Ops, pg.Theme.Color.Primary, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: selectedTabDims.Size.X, Y: tabHeight},
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
