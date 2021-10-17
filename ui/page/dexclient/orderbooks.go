package dexclient

import (
	"fmt"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type OrderBooksWidget struct {
	sellOrders layout.List
	buyOrders  layout.List
	th         *decredmaterial.Theme
	text       decredmaterial.Label
}

func NewOrderBooksWidget(l *load.Load) *OrderBooksWidget {
	return &OrderBooksWidget{
		th:         l.Theme,
		sellOrders: layout.List{Axis: layout.Vertical},
		buyOrders:  layout.List{Axis: layout.Vertical},
		text:       l.Theme.Label(values.TextSize12, ""),
	}
}

func (ordWdg *OrderBooksWidget) Layout(gtx C, sells, buys []*core.MiniOrder) D {
	gtx.Constraints.Min.X, gtx.Constraints.Max.X = 450, 450

	return layout.Inset{
		Left:  values.MarginPadding10,
		Right: values.MarginPadding10,
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Flexed(0.4, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Quantity").Layout(gtx)
					}),
					layout.Flexed(0.4, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Rate").Layout(gtx)
					}),
					layout.Flexed(0.2, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Epoch").Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = 250, 250
						return ordWdg.orderBooksLayout(gtx, &ordWdg.buyOrders, sells)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Top:    values.MarginPadding8,
							Bottom: values.MarginPadding8,
						}.Layout(gtx, func(gtx C) D {
							return ordWdg.th.Body2("374.8770363158391").Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.Y, gtx.Constraints.Max.Y = 250, 250
						return ordWdg.orderBooksLayout(gtx, &ordWdg.sellOrders, buys)
					}),
				)
			}),
		)
	})
}

func (ordWdg *OrderBooksWidget) orderBooksLayout(gtx C, l *layout.List, orders []*core.MiniOrder) D {
	return l.Layout(gtx, len(orders), func(gtx C, i int) D {
		ord := orders[i]
		if ord.Sell {
			ordWdg.text.Color = ordWdg.th.Color.ChartSellLine
		} else {
			ordWdg.text.Color = ordWdg.th.Color.ChartBuyLine
		}

		return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Flexed(0.4, func(gtx C) D {
				ordWdg.text.Text = fmt.Sprintf("%f", ord.Qty)
				return ordWdg.text.Layout(gtx)
			}),
			layout.Flexed(0.4, func(gtx C) D {
				ordWdg.text.Text = fmt.Sprintf("%f", ord.Rate)
				return ordWdg.text.Layout(gtx)
			}),
			layout.Flexed(0.2, func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return ordWdg.th.Body1(fmt.Sprintf("%b", ord.Epoch)).Layout(gtx)
				})
			}),
		)
	})
}

type UserOrderBooksWidget struct {
	myOrders layout.List
	th       *decredmaterial.Theme
}

func NewUserOrderBooksWidget(th *decredmaterial.Theme) *UserOrderBooksWidget {
	return &UserOrderBooksWidget{
		th:       th,
		myOrders: layout.List{Axis: layout.Vertical},
	}
}

func (ordWdg *UserOrderBooksWidget) Layout(gtx C, orders []*core.Order) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return layout.Inset{
		Left:  values.MarginPadding10,
		Right: values.MarginPadding10,
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(0.125, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Type").Layout(gtx)
					}),
					layout.Flexed(0.125, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Side").Layout(gtx)
					}),
					layout.Flexed(0.125, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Age").Layout(gtx)
					}),
					layout.Flexed(0.125, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Rate").Layout(gtx)
					}),
					layout.Flexed(0.125, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Quantity").Layout(gtx)
					}),
					layout.Flexed(0.125, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Filled").Layout(gtx)
					}),
					layout.Flexed(0.125, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Settled").Layout(gtx)
					}),
					layout.Flexed(0.125, func(gtx C) D {
						return ordWdg.th.Label(values.TextSize14, "Status").Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return ordWdg.myOrders.Layout(gtx, len(orders), func(gtx C, i int) D {
					ord := orders[i]

					return layout.S.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Flexed(0.125, func(gtx C) D {
								return ordWdg.th.Label(values.TextSize14, typeString(ord)).Layout(gtx)
							}),
							layout.Flexed(0.125, func(gtx C) D {
								return ordWdg.th.Label(values.TextSize14, sellString(ord)).Layout(gtx)
							}),
							layout.Flexed(0.125, func(gtx C) D {
								return ordWdg.th.Label(values.TextSize14, timeSince(ord.Stamp)).Layout(gtx)
							}),
							layout.Flexed(0.125, func(gtx C) D {
								return ordWdg.th.Label(values.TextSize14, rateString(ord)).Layout(gtx)
							}),
							layout.Flexed(0.125, func(gtx C) D {
								return ordWdg.th.Label(values.TextSize14, formatCoinValue(ord.Qty)).Layout(gtx)
							}),
							layout.Flexed(0.125, func(gtx C) D {
								return ordWdg.th.Label(values.TextSize14, fmt.Sprintf("%.1f%%", (float64(ord.Filled)/float64(ord.Qty))*100)).Layout(gtx)
							}),
							layout.Flexed(0.125, func(gtx C) D {
								return ordWdg.th.Label(values.TextSize14, fmt.Sprintf("%.1f%%", settled(ord)/float64(ord.Qty)*100)).Layout(gtx)
							}),
							layout.Flexed(0.125, func(gtx C) D {
								return ordWdg.th.Label(values.TextSize14, statusString(ord)).Layout(gtx)
							}),
						)
					})
				})
			}),
		)
	})
}
