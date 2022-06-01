package page

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/governance"
	"github.com/planetdecred/godcr/ui/page/staking"
	"github.com/planetdecred/godcr/ui/values"
)

const MorePageID = "More"

type morePageHandler struct {
	clickable *decredmaterial.Clickable
	image     *decredmaterial.Image
	page      string
	action    func()
}

type MorePage struct {
	*load.Load
	container                layout.Flex
	shadowBox                *decredmaterial.Shadow
	morePageListItemsDesktop []morePageHandler
	morePageListItemsMobile  []morePageHandler
}

func NewMorePage(l *load.Load) *MorePage {
	pg := &MorePage{
		container: layout.Flex{Axis: layout.Vertical},
		Load:      l,
		shadowBox: l.Theme.Shadow(),
	}
	pg.initPageItems()

	return pg
}

func (pg *MorePage) initPageItems() {
	pg.morePageListItemsDesktop = []morePageHandler{
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.SettingsIcon,
			page:      SettingsPageID,
			action: func() {
				pg.ChangeFragment(NewSettingsPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.SecurityIcon,
			page:      SecurityToolsPageID,
			action: func() {
				pg.ChangeFragment(NewSecurityToolsPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.HelpIcon,
			page:      HelpPageID,
			action: func() {
				pg.ChangeFragment(NewHelpPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.AboutIcon,
			page:      AboutPageID,
			action: func() {
				pg.ChangeFragment(NewAboutPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.DebugIcon,
			page:      DebugPageID,
			action: func() {
				pg.ChangeFragment(NewDebugPage(pg.Load))
			},
		},
	}

	pg.morePageListItemsMobile = []morePageHandler{
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.SettingsIcon,
			page:      SettingsPageID,
			action: func() {
				pg.ChangeFragment(NewSettingsPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.SecurityIcon,
			page:      SecurityToolsPageID,
			action: func() {
				pg.ChangeFragment(NewSecurityToolsPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.StakeIcon,
			page:      values.String(values.StrStaking),
			action: func() {
				pg.ChangeFragment(staking.NewStakingPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.GovernanceActiveIcon,
			page:      "Governance",
			action: func() {
				pg.ChangeFragment(governance.NewGovernancePage(pg.Load))
			},
		},
		// Temp disabling. Will uncomment after release
		// {
		// 	clickable:     pg.Theme.NewClickable(true),
		// 	image:         pg.Theme.Icons.DexIcon,
		// 	page:         values.String(values.StrDex),
		// 	action: func() {
		// 		_, err := pg.WL.MultiWallet.StartDexClient() // does nothing if already started
		// 		if err != nil {
		// 			pg.Toast.NotifyError(fmt.Sprintf("Unable to start DEX client: %v", err))
		// 		} else {
		// 			pg = dexclient.NewMarketPage(pg.Load)
		// 		}
		// 	},
		// },
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.HelpIcon,
			page:      HelpPageID,
			action: func() {
				pg.ChangeFragment(NewHelpPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.AboutIcon,
			page:      AboutPageID,
			action: func() {
				pg.ChangeFragment(NewAboutPage(pg.Load))
			},
		},
		{
			clickable: pg.Theme.NewClickable(true),
			image:     pg.Theme.Icons.DebugIcon,
			page:      DebugPageID,
			action: func() {
				pg.ChangeFragment(NewDebugPage(pg.Load))
			},
		},
	}
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *MorePage) ID() string {
	return MorePageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *MorePage) OnNavigatedTo() {
	pg.initPageItems() //re-initialize the nav options to reflect the changes if theme was changed.
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *MorePage) HandleUserInteractions() {
	for _, item := range pg.morePageListItemsDesktop {
		for item.clickable.Clicked() {
			item.action()
		}
	}

	for _, item := range pg.morePageListItemsMobile {
		for item.clickable.Clicked() {
			item.action()
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
func (pg *MorePage) OnNavigatedFrom() {}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *MorePage) Layout(gtx layout.Context) layout.Dimensions {
	if pg.Load.GetCurrentAppWidth() <= gtx.Px(values.StartMobileView) {
		container := func(gtx C) D {
			pg.layoutMoreItemsMobile(gtx)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}
		return components.UniformPadding(gtx, container)
	}

	container := func(gtx C) D {
		pg.layoutMoreItemsDesktop(gtx)
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}
	return components.UniformPadding(gtx, container)
}

func (pg *MorePage) layoutMoreItemsDesktop(gtx layout.Context) layout.Dimensions {

	list := layout.List{Axis: layout.Vertical}
	return list.Layout(gtx, len(pg.morePageListItemsDesktop), func(gtx C, i int) D {
		radius := decredmaterial.Radius(14)
		pg.shadowBox.SetShadowRadius(14)
		pg.shadowBox.SetShadowElevation(5)
		return decredmaterial.LinearLayout{
			Orientation: layout.Horizontal,
			Width:       decredmaterial.MatchParent,
			Height:      decredmaterial.WrapContent,
			Background:  pg.Theme.Color.Surface,
			Clickable:   pg.morePageListItemsDesktop[i].clickable,
			Direction:   layout.W,
			Shadow:      pg.shadowBox,
			Border:      decredmaterial.Border{Radius: radius},
			Padding:     layout.UniformInset(values.MarginPadding15),
			Margin:      layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.morePageListItemsDesktop[i].image.Layout24dp(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Top:  values.MarginPadding2,
					Left: values.MarginPadding18,
				}.Layout(gtx, func(gtx C) D {
					page := pg.morePageListItemsDesktop[i].page
					if page == SecurityToolsPageID {
						page = "Security Tools"
					}
					return pg.Theme.Body1(page).Layout(gtx)
				})
			}),
		)
	})
}

func (pg *MorePage) layoutMoreItemsMobile(gtx layout.Context) layout.Dimensions {

	list := layout.List{Axis: layout.Vertical}
	return list.Layout(gtx, len(pg.morePageListItemsMobile), func(gtx C, i int) D {
		radius := decredmaterial.Radius(14)
		pg.shadowBox.SetShadowRadius(14)
		pg.shadowBox.SetShadowElevation(5)
		return decredmaterial.LinearLayout{
			Orientation: layout.Horizontal,
			Width:       decredmaterial.MatchParent,
			Height:      decredmaterial.WrapContent,
			Background:  pg.Theme.Color.Surface,
			Clickable:   pg.morePageListItemsMobile[i].clickable,
			Direction:   layout.W,
			Shadow:      pg.shadowBox,
			Border:      decredmaterial.Border{Radius: radius},
			Padding:     layout.UniformInset(values.MarginPadding15),
			Margin:      layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.morePageListItemsMobile[i].image.Layout24dp(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Top:  values.MarginPadding2,
					Left: values.MarginPadding18,
				}.Layout(gtx, func(gtx C) D {
					page := pg.morePageListItemsMobile[i].page
					if page == SecurityToolsPageID {
						page = "Security Tools"
					}
					return pg.Theme.Body1(page).Layout(gtx)
				})
			}),
		)
	})
}
