package components

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

var (
	NavDrawerWidth          = unit.Value{U: unit.UnitDp, V: 180}
	NavDrawerMinimizedWidth = unit.Value{U: unit.UnitDp, V: 100}
)

type NavHandler struct {
	Clickable     *widget.Clickable
	Image         *widget.Image
	ImageInactive *widget.Image
	Title         string
	PageID        string
}

type NavDrawer struct {
	*load.Load

	AppBarNavItems []NavHandler
	DrawerNavItems []NavHandler
	CurrentPage    string

	Axis      layout.Axis
	TextSize  unit.Value
	LeftInset unit.Value
	Width     unit.Value
	Alignment layout.Alignment
	Direction layout.Direction

	MinimizeNavDrawerButton decredmaterial.IconButton
	MaximizeNavDrawerButton decredmaterial.IconButton
	ActiveDrawerBtn         decredmaterial.IconButton
}

func (nd *NavDrawer) LayoutNavDrawer(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:       gtx.Px(nd.Width),
				Height:      decredmaterial.MatchParent,
				Background:  nd.Theme.Color.Surface,
				Orientation: nd.Axis,
			}.Layout2(gtx, func(gtx C) D {
				list := layout.List{Axis: layout.Vertical}
				return list.Layout(gtx, len(nd.DrawerNavItems), func(gtx C, i int) D {
					background := nd.Theme.Color.Surface
					if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
						background = nd.Theme.Color.ActiveGray
					}

					return nd.layoutNavRow(gtx, background, nd.DrawerNavItems[i].Clickable, func(gtx C) D {
						txt := nd.Theme.Label(nd.TextSize, nd.DrawerNavItems[i].Title)

						gtx.Constraints.Min.X = gtx.Px(nd.Width)
						return decredmaterial.Clickable(gtx, nd.DrawerNavItems[i].Clickable, func(gtx C) D {
							return decredmaterial.LinearLayout{
								Orientation: nd.Axis,
								Width:       decredmaterial.MatchParent,
								Height:      decredmaterial.WrapContent,
								Padding:     layout.UniformInset(values.MarginPadding15),
								Alignment:   nd.Alignment,
								Direction:   nd.Direction,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									img := nd.DrawerNavItems[i].ImageInactive

									if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
										img = nd.DrawerNavItems[i].Image
									}

									img.Scale = 1.0
									return img.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Left: nd.LeftInset,
										Top:  values.MarginPadding4,
									}.Layout(gtx, func(gtx C) D {
										textColor := nd.Theme.Color.Gray4
										if nd.DrawerNavItems[i].PageID == nd.CurrentPage {
											textColor = nd.Theme.Color.DeepBlue
										}
										txt.Color = textColor
										return txt.Layout(gtx)
									})
								}),
							)
						})
					})
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.SE.Layout(gtx, func(gtx C) D {
				btn := nd.ActiveDrawerBtn
				btn.Color = nd.Theme.Color.Gray3

				return btn.Layout(gtx)
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
					background = nd.Theme.Color.ActiveGray
				}

				// header buttons container
				return nd.layoutNavRow(gtx, background, nd.AppBarNavItems[i].Clickable, func(gtx C) D {
					return decredmaterial.Clickable(gtx, nd.AppBarNavItems[i].Clickable, func(gtx C) D {
						return Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
							return decredmaterial.LinearLayout{
								Width:       decredmaterial.WrapContent,
								Height:      decredmaterial.WrapContent,
								Orientation: layout.Horizontal,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
										func(gtx C) D {
											return layout.Center.Layout(gtx, func(gtx C) D {
												img := nd.AppBarNavItems[i].Image
												img.Scale = 1.0
												return nd.AppBarNavItems[i].Image.Layout(gtx)
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
			})
		})
	})
}

func (nd *NavDrawer) layoutNavRow(gtx layout.Context, background color.NRGBA, Clickable *widget.Clickable, body layout.Widget) layout.Dimensions {
	card := nd.Theme.Card()
	card.Color = background
	card.Radius = decredmaterial.Radius(0)
	return card.HovarableLayout(gtx, Clickable, body)
}
