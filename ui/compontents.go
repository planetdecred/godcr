package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"strconv"
	"strings"
)

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page pageCommon) layoutBalance(gtx layout.Context, amount string) layout.Dimensions {
	mainText, subText := page.breakBalance(amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return page.theme.Label(unit.Sp(20), mainText).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return page.theme.Label(unit.Sp(14), subText).Layout(gtx)
		}),
	)
}

// breakBalance takes the balance string and returns it in two slices
func (page pageCommon) breakBalance(balance string) (b1, b2 string) {
	balanceParts := strings.Split(balance, ".")
	if len(balanceParts) == 1 {
		return balanceParts[0], ""
	}

	b1 = balanceParts[0]
	if bal, err := strconv.Atoi(b1); err == nil {
		b1 = page.printer.Sprint(bal)
	}

	b2 = balanceParts[1]
	b1 = b1 + "." + b2[:2]
	b2 = b2[2:]
	return
}

func (page pageCommon) layoutTopBar(gtx layout.Context) layout.Dimensions {
	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							h := values.MarginPadding24
							v := values.MarginPadding16
							// Balance container
							return Container{padding: layout.Inset{Right: h, Left: h, Top: v, Bottom: v}}.Layout(gtx,
								func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											img := page.icons.logo
											img.Scale = 1.0
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx,
												func(gtx C) D {
													return img.Layout(gtx)
												})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Center.Layout(gtx, func(gtx C) D {
												return page.layoutBalance(gtx, page.info.TotalBalance)
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

func (page pageCommon) layoutNavDrawer(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(page.drawerNavItems), func(gtx C, i int) D {
				background := page.theme.Color.Surface
				if page.drawerNavItems[i].page == *page.page {
					background = page.theme.Color.LightGray
				}
				txt := page.theme.Label(values.TextSize16, page.drawerNavItems[i].page)
				return decredmaterial.Clickable(gtx, page.drawerNavItems[i].clickable, func(gtx C) D {
					card := page.theme.Card()
					card.Color = background
					card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
					return card.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
							axis := layout.Horizontal
							leftInset := values.MarginPadding15
							width := navDrawerWidth
							if *page.isNavDrawerMinimized {
								axis = layout.Vertical
								txt.TextSize = values.TextSize10
								leftInset = values.MarginPadding0
								width = navDrawerMinimizedWidth
							}

							gtx.Constraints.Min.X = int(gtx.Metric.PxPerDp) * width
							return layout.Flex{Axis: axis}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									img := page.drawerNavItems[i].imageInactive
									if page.drawerNavItems[i].page == *page.page {
										img = page.drawerNavItems[i].image
									}
									return layout.Center.Layout(gtx, func(gtx C) D {
										img.Scale = 1.0
										return img.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Left: leftInset,
										Top:  values.MarginPadding4,
									}.Layout(gtx, func(gtx C) D {
										return layout.Center.Layout(gtx, txt.Layout)
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
				btn := page.minimizeNavDrawerButton
				if *page.isNavDrawerMinimized {
					btn = page.maximizeNavDrawerButton
				}
				return btn.Layout(gtx)
			})
		}),
	)
}

func (page pageCommon) handleNavEvents() {
	for page.minimizeNavDrawerButton.Button.Clicked() {
		*page.isNavDrawerMinimized = true
	}

	for page.maximizeNavDrawerButton.Button.Clicked() {
		*page.isNavDrawerMinimized = false
	}

	for i := range page.appBarNavItems {
		for page.appBarNavItems[i].clickable.Clicked() {
			page.changePage(page.appBarNavItems[i].page)
		}
	}

	for i := range page.drawerNavItems {
		for page.drawerNavItems[i].clickable.Clicked() {
			page.changePage(page.drawerNavItems[i].page)
		}
	}
}
