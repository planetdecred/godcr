package components

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

var (
	bottomNavigationBarHeight = values.MarginPadding100
)

type BottomNavigationBarHandler struct {
	Clickable     *decredmaterial.Clickable
	Image         *decredmaterial.Image
	ImageInactive *decredmaterial.Image
	Title         string
	PageID        string
}

type BottomNavigationBar struct {
	*load.Load

	AppBarNavItems       []BottomNavigationBarHandler
	BottomNaigationItems []BottomNavigationBarHandler
	CurrentPage          string

	axis        layout.Axis
	textSize    unit.Value
	bottomInset unit.Value
	height      unit.Value
	alignment   layout.Alignment
	direction   layout.Direction
}

func (bottomNavigationbar *BottomNavigationBar) LayoutBottomNavigationBar(gtx layout.Context) layout.Dimensions {
	return layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return UniformPadding(gtx, func(gtx C) D {
				return decredmaterial.LinearLayout{
					Width:       decredmaterial.WrapContent,
					Height:      decredmaterial.WrapContent,
					Orientation: layout.Horizontal,
					Background:  bottomNavigationbar.Theme.Color.Surface,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						list := layout.List{Axis: layout.Horizontal}
						return list.Layout(gtx, len(bottomNavigationbar.BottomNaigationItems), func(gtx C, i int) D {

							background := bottomNavigationbar.Theme.Color.Surface
							if bottomNavigationbar.BottomNaigationItems[i].PageID == bottomNavigationbar.CurrentPage {
								background = bottomNavigationbar.Theme.Color.Gray5
							}
							return decredmaterial.LinearLayout{
								Orientation: bottomNavigationbar.axis,
								Width:       (gtx.Px(values.AppWidth) * 100 / len(bottomNavigationbar.BottomNaigationItems)) / 100, // Divide each cell equally
								Height:      decredmaterial.WrapContent,
								Padding:     layout.UniformInset(values.MarginPadding15),
								Alignment:   bottomNavigationbar.alignment,
								Direction:   bottomNavigationbar.direction,
								Background:  background,
								Clickable:   bottomNavigationbar.BottomNaigationItems[i].Clickable,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									img := bottomNavigationbar.BottomNaigationItems[i].ImageInactive

									if bottomNavigationbar.BottomNaigationItems[i].PageID == bottomNavigationbar.CurrentPage {
										img = bottomNavigationbar.BottomNaigationItems[i].Image
									}

									return img.Layout24dp(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Bottom: bottomNavigationbar.bottomInset,
									}.Layout(gtx, func(gtx C) D {
										textColor := bottomNavigationbar.Theme.Color.GrayText1
										if bottomNavigationbar.BottomNaigationItems[i].PageID == bottomNavigationbar.CurrentPage {
											textColor = bottomNavigationbar.Theme.Color.DeepBlue
										}
										txt := bottomNavigationbar.Theme.Label(bottomNavigationbar.textSize, bottomNavigationbar.BottomNaigationItems[i].Title)
										txt.Color = textColor
										return txt.Layout(gtx)
									})
								}),
							)
						})
					}),
				)
			})
		}),
	)
}

// func (bottomNavigationbar *BottomNavigationBar) LayoutTopBar(gtx layout.Context) layout.Dimensions {
// 	gtx.Constraints.Min.X = gtx.Constraints.Max.X
// 	return layout.E.Layout(gtx, func(gtx C) D {
// 		return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
// 			list := layout.List{Axis: layout.Horizontal}
// 			return list.Layout(gtx, len(bottomNavigationbar.AppBarNavItems), func(gtx C, i int) D {
// 				background := bottomNavigationbar.Theme.Color.Surface
// 				if bottomNavigationbar.AppBarNavItems[i].PageID == bottomNavigationbar.CurrentPage {
// 					background = bottomNavigationbar.Theme.Color.Gray5
// 				}
// 				return decredmaterial.LinearLayout{
// 					Width:       decredmaterial.WrapContent,
// 					Height:      decredmaterial.WrapContent,
// 					Orientation: layout.Horizontal,
// 					Background:  background,
// 					Padding:     layout.UniformInset(values.MarginPadding16),
// 					Clickable:   bottomNavigationbar.AppBarNavItems[i].Clickable,
// 				}.Layout(gtx,
// 					layout.Rigid(func(gtx C) D {
// 						return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
// 							func(gtx C) D {
// 								return layout.Center.Layout(gtx, func(gtx C) D {
// 									return bottomNavigationbar.AppBarNavItems[i].Image.Layout24dp(gtx)
// 								})
// 							})
// 					}),
// 					layout.Rigid(func(gtx C) D {
// 						return layout.Inset{
// 							Left: values.MarginPadding0,
// 						}.Layout(gtx, func(gtx C) D {
// 							return layout.Center.Layout(gtx, func(gtx C) D {
// 								return bottomNavigationbar.Theme.Body1(bottomNavigationbar.AppBarNavItems[i].Title).Layout(gtx)
// 							})
// 						})
// 					}),
// 				)
// 			})
// 		})
// 	})
// }

func (bottomNavigationbar *BottomNavigationBar) OnViewCreated() {
	bottomNavigationbar.axis = layout.Vertical
	bottomNavigationbar.textSize = values.TextSize12
	bottomNavigationbar.bottomInset = values.MarginPadding0
	bottomNavigationbar.height = bottomNavigationBarHeight
	bottomNavigationbar.alignment = layout.Middle
	bottomNavigationbar.direction = layout.Center
}
