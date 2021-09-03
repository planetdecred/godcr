package page

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const MorePageID = "More"

type morePageHandler struct {
	clickable *widget.Clickable
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
			clickable: new(widget.Clickable),
			image:     l.Icons.SettingsIcon,
			page:      SettingsPageID,
			action: func() {
				l.ChangeFragment(NewSettingsPage(l))
			},
		},
		{
			clickable: new(widget.Clickable),
			image:     l.Icons.SecurityIcon,
			page:      SecurityToolsPageID,
			action: func() {
				l.ChangeFragment(NewSecurityToolsPage(l))
			},
		},
		{
			clickable: new(widget.Clickable),
			image:     l.Icons.HelpIcon,
			page:      HelpPageID,
			action: func() {
				l.ChangeFragment(NewHelpPage(l))
			},
		},
		{
			clickable: new(widget.Clickable),
			image:     l.Icons.AboutIcon,
			page:      AboutPageID,
			action: func() {
				l.ChangeFragment(NewAboutPage(l))
			},
		},
		{
			clickable: new(widget.Clickable),
			image:     l.Icons.DebugIcon,
			page:      DebugPageID,
			action: func() {
				l.ChangeFragment(NewDebugPage(l))
			},
		},
	}

	for i := range morePageListItems {
		morePageListItems[i].image.Scale = 1
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
		return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			return decredmaterial.Clickable(gtx, pg.morePageListItems[i].clickable, func(gtx C) D {
				return decredmaterial.LinearLayout{Orientation: layout.Horizontal,
					Width:      decredmaterial.MatchParent,
					Height:     decredmaterial.WrapContent,
					Background: pg.Theme.Color.Surface,
					Border:     decredmaterial.Border{Radius: decredmaterial.Radius(14)},
					Padding:    layout.UniformInset(values.MarginPadding15)}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Center.Layout(gtx, pg.morePageListItems[i].image.Layout24dp)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: values.MarginPadding15,
							Top:  values.MarginPadding2,
						}.Layout(gtx, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								page := pg.morePageListItems[i].page
								if page == SecurityToolsPageID {
									page = "Security Tools"
								}
								return pg.Theme.Body1(page).Layout(gtx)
							})
						})
					}),
				)
			})
		})
	})
}
