package components

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

var (
	navDrawerMaximizedWidth = values.Size180
	navDrawerMinimizedWidth = values.MarginPadding100
)

type NavHandler struct {
	Clickable     *decredmaterial.Clickable
	Image         *decredmaterial.Image
	ImageInactive *decredmaterial.Image
	Title         string
	PageID        string
}

type NavDrawer struct {
	Theme *decredmaterial.Theme

	AppBarNavItems []NavHandler
	DrawerNavItems []NavHandler
	CurrentPage    string

	axis      layout.Axis
	textSize  unit.Value
	leftInset unit.Value
	width     unit.Value
	alignment layout.Alignment
	direction layout.Direction

	MinimizeNavDrawerButton decredmaterial.IconButton
	MaximizeNavDrawerButton decredmaterial.IconButton
	activeDrawerBtn         decredmaterial.IconButton
}

func (nd *NavDrawer) LayoutNavDrawer(gtx layout.Context) layout.Dimensions {
	return decredmaterial.LinearLayout{
		Width:       gtx.Px(nd.width),
		Height:      decredmaterial.MatchParent,
		Orientation: layout.Vertical,
		Background:  nd.Theme.Color.Surface,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(nd.DrawerNavItems), func(gtx C, i int) D {

				background := nd.Theme.Color.Surface
				if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
					background = nd.Theme.Color.Gray5
				}
				return decredmaterial.LinearLayout{
					Orientation: nd.axis,
					Width:       decredmaterial.MatchParent,
					Height:      decredmaterial.WrapContent,
					Padding:     layout.UniformInset(values.MarginPadding15),
					Alignment:   nd.alignment,
					Direction:   nd.direction,
					Background:  background,
					Clickable:   nd.DrawerNavItems[i].Clickable,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						img := nd.DrawerNavItems[i].ImageInactive

						if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
							img = nd.DrawerNavItems[i].Image
						}

						return img.Layout24dp(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: nd.leftInset,
						}.Layout(gtx, func(gtx C) D {
							textColor := nd.Theme.Color.GrayText1
							if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
								textColor = nd.Theme.Color.DeepBlue
							}
							txt := nd.Theme.Label(nd.textSize, nd.DrawerNavItems[i].Title)
							txt.Color = textColor
							return txt.Layout(gtx)
						})
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.SE.Layout(gtx, func(gtx C) D {
				return nd.activeDrawerBtn.Layout(gtx)
			})
		}),
	)
}

func (nd *NavDrawer) LayoutTopBar(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.E.Layout(gtx, func(gtx C) D {
		return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
			list := layout.List{Axis: layout.Horizontal}
			return list.Layout(gtx, len(nd.AppBarNavItems), func(gtx C, i int) D {
				background := nd.Theme.Color.Surface
				if nd.AppBarNavItems[i].PageID == nd.CurrentPage {
					background = nd.Theme.Color.Gray5
				}
				return decredmaterial.LinearLayout{
					Width:       decredmaterial.WrapContent,
					Height:      decredmaterial.WrapContent,
					Orientation: layout.Horizontal,
					Background:  background,
					Padding:     layout.UniformInset(values.MarginPadding16),
					Clickable:   nd.AppBarNavItems[i].Clickable,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
							func(gtx C) D {
								return layout.Center.Layout(gtx, func(gtx C) D {
									return nd.AppBarNavItems[i].Image.Layout24dp(gtx)
								})
							})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: values.MarginPadding0,
						}.Layout(gtx, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return nd.Theme.Body1(nd.AppBarNavItems[i].Title).Layout(gtx)
							})
						})
					}),
				)
			})
		})
	})
}

func (nd *NavDrawer) DrawerToggled(min bool) {
	if min {
		nd.axis = layout.Vertical
		nd.textSize = values.TextSize12
		nd.leftInset = values.MarginPadding0
		nd.width = navDrawerMinimizedWidth
		nd.activeDrawerBtn = nd.MaximizeNavDrawerButton
		nd.alignment = layout.Middle
		nd.direction = layout.Center
	} else {
		nd.axis = layout.Horizontal
		nd.textSize = values.TextSize16
		nd.leftInset = values.MarginPadding15
		nd.width = navDrawerMaximizedWidth
		nd.activeDrawerBtn = nd.MinimizeNavDrawerButton
		nd.alignment = layout.Start
		nd.direction = layout.W
	}
}
