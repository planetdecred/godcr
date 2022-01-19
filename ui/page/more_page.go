package page

import (
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"

	"github.com/planetdecred/godcr/ui/decredmaterial"
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
	container         layout.Flex
	shadowBox         *decredmaterial.Shadow
	morePageListItems []morePageHandler
}

func NewMorePage(l *load.Load) *MorePage {
	morePageListItems := []morePageHandler{
		{
			clickable: l.Theme.NewClickable(true),
			image:     l.Icons.SettingsIcon,
			page:      SettingsPageID,
			action: func() {
				l.ChangeFragment(NewSettingsPage(l))
			},
		},
		{
			clickable: l.Theme.NewClickable(true),
			image:     l.Icons.SecurityIcon,
			page:      SecurityToolsPageID,
			action: func() {
				l.ChangeFragment(NewSecurityToolsPage(l))
			},
		},
		{
			clickable: l.Theme.NewClickable(true),
			image:     l.Icons.HelpIcon,
			page:      HelpPageID,
			action: func() {
				l.ChangeFragment(NewHelpPage(l))
			},
		},
		{
			clickable: l.Theme.NewClickable(true),
			image:     l.Icons.AboutIcon,
			page:      AboutPageID,
			action: func() {
				l.ChangeFragment(NewAboutPage(l))
			},
		},
		{
			clickable: l.Theme.NewClickable(true),
			image:     l.Icons.DebugIcon,
			page:      DebugPageID,
			action: func() {
				l.ChangeFragment(NewDebugPage(l))
			},
		},
	}

	pg := &MorePage{
		container:         layout.Flex{Axis: layout.Vertical},
		morePageListItems: morePageListItems,
		Load:              l,
		shadowBox:         l.Theme.Shadow(),
	}

	return pg
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

}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *MorePage) HandleUserInteractions() {
	for _, item := range pg.morePageListItems {
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

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *MorePage) Layout(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		pg.layoutMoreItems(gtx)
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}
	return components.UniformPadding(gtx, container)
}

func (pg *MorePage) layoutMoreItems(gtx layout.Context) layout.Dimensions {

	list := layout.List{Axis: layout.Vertical}
	return list.Layout(gtx, len(pg.morePageListItems), func(gtx C, i int) D {
		radius := decredmaterial.Radius(14)
		pg.shadowBox.SetShadowRadius(14)
		pg.shadowBox.SetShadowElevation(5)
		return decredmaterial.LinearLayout{
			Orientation: layout.Horizontal,
			Width:       decredmaterial.MatchParent,
			Height:      decredmaterial.WrapContent,
			Background:  pg.Theme.Color.Surface,
			Clickable:   pg.morePageListItems[i].clickable,
			Direction:   layout.W,
			Shadow:      pg.shadowBox,
			Border:      decredmaterial.Border{Radius: radius},
			Padding:     layout.UniformInset(values.MarginPadding15),
			Margin:      layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.morePageListItems[i].image.Layout24dp(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Top:  values.MarginPadding2,
					Left: values.MarginPadding18,
				}.Layout(gtx, func(gtx C) D {
					page := pg.morePageListItems[i].page
					if page == SecurityToolsPageID {
						page = "Security Tools"
					}
					return pg.Theme.Body1(page).Layout(gtx)
				})
			}),
		)
	})
}
