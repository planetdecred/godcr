// components contain layout code that are shared by multiple pages but aren't widely used enough to be defined as
// widgets

package uidex

import (
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

// layoutTopBar is the top horizontal bar on every page of the app. It lays out the wallet balance, receive and send
// buttons.
func (page pageCommon) layoutTopBar(gtx layout.Context) layout.Dimensions {
	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							h := values.MarginPadding16
							v := values.MarginPadding10
							return Container{padding: layout.Inset{Right: h, Left: h, Top: v, Bottom: v}}.Layout(gtx,
								func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											img := page.icons.logo
											img.Scale = .3
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx,
												func(gtx C) D {
													return img.Layout(gtx)
												})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Center.Layout(gtx, func(gtx C) D {
												return page.theme.H5("DCRDEX").Layout(gtx)
											})
										}),
									)
								})
						})
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
								list := layout.List{Axis: layout.Horizontal}
								return list.Layout(gtx, len(page.appBarNavItems), func(gtx C, i int) D {
									// header buttons container
									return Container{layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
										return decredmaterial.Clickable(gtx, page.appBarNavItems[i].clickable, func(gtx C) D {
											return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
														func(gtx C) D {
															return layout.Center.Layout(gtx, func(gtx C) D {
																img := page.appBarNavItems[i].image
																img.Scale = 1.0
																return page.appBarNavItems[i].image.Layout(gtx)
															})
														})
												}),
												layout.Rigid(func(gtx C) D {
													return layout.Inset{
														Left: values.MarginPadding0,
													}.Layout(gtx, func(gtx C) D {
														return layout.Center.Layout(gtx, func(gtx C) D {
															return page.theme.Body1(page.appBarNavItems[i].page).Layout(gtx)
														})
													})
												}),
											)
										})
									})
								})
							})
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return page.theme.Separator().Layout(gtx)
			}),
		)
	})
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func endToEndRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return leftWidget(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return rightWidget(gtx)
			})
		}),
	)
}
