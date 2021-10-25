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
	}

	return pg
}

func (pg *MorePage) ID() string {
	return MorePageID
}

func (pg *MorePage) OnResume() {

}

func (pg *MorePage) Handle() {
	for _, item := range pg.morePageListItems {
		for item.clickable.Clicked() {
			item.action()
		}
	}
}

func (pg *MorePage) OnClose() {}

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
		return decredmaterial.LinearLayout{
			Orientation: layout.Horizontal,
			Width:       decredmaterial.MatchParent,
			Height:      decredmaterial.WrapContent,
			Background:  pg.Theme.Color.Surface,
			Shadow:      pg.Theme.TransparentShadow(14),
			Clickable:   pg.morePageListItems[i].clickable,
			Direction:   layout.W,
			Border:      decredmaterial.Border{Radius: radius},
			Padding:     layout.UniformInset(values.MarginPadding15),
			Margin:      layout.Inset{Bottom: values.MarginPadding8}}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.morePageListItems[i].image.Layout24dp(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Top: values.MarginPadding2,
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
