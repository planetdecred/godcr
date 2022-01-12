package dexclient

import (
	"context"
	"fmt"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const OrdersHistoryPageID = "OrdersHistory"

type OrdersHistoryPage struct {
	*load.Load
	ctx        context.Context
	ctxCancel  context.CancelFunc
	list       *widget.List
	backButton decredmaterial.IconButton
	host       string
	orders     []*core.Order
}

func NewOrdersHistoryPage(l *load.Load, host string) *OrdersHistoryPage {
	pg := &OrdersHistoryPage{
		Load: l,
		host: host,
		list: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}

	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)

	return pg
}

func (pg *OrdersHistoryPage) ID() string {
	return OrdersHistoryPageID
}

func (pg *OrdersHistoryPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	go pg.readNotifications()
	pg.getOrdersHistory()
}

func (pg *OrdersHistoryPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      strOrderHistory,
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Flexed(0.125, pg.layoutLabel(strTrade)),
									layout.Flexed(0.125, pg.layoutLabel(strSide)),
									layout.Flexed(0.125, pg.layoutLabel(strRate)),
									layout.Flexed(0.125, pg.layoutLabel(strQuantity)),
									layout.Flexed(0.125, pg.layoutLabel(strFilled)),
									layout.Flexed(0.125, pg.layoutLabel(strSettled)),
									layout.Flexed(0.125, pg.layoutLabel(strStatus)),
									layout.Flexed(0.125, pg.layoutLabel(strTime)),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Top:    values.MarginPadding8,
									Bottom: values.MarginPadding8}.Layout(gtx, pg.Theme.Separator().Layout)
							}),
							layout.Rigid(func(gtx C) D {
								return pg.Theme.List(pg.list).Layout(gtx, len(pg.orders), func(gtx C, i int) D {
									order := pg.orders[i]
									return layout.Flex{}.Layout(gtx,
										layout.Flexed(0.125, pg.layoutLabel(typeString(order))),
										layout.Flexed(0.125, pg.layoutLabel(sellString(order))),
										layout.Flexed(0.125, pg.layoutLabel(rateString(order))),
										layout.Flexed(0.125, pg.layoutLabel(formatCoinValue(order.Qty))),
										layout.Flexed(0.125, pg.layoutLabel(fmt.Sprintf("%.1f%%", (float64(order.Filled)/float64(order.Qty))*100))),
										layout.Flexed(0.125, pg.layoutLabel(fmt.Sprintf("%.1f%%", settled(order)/float64(order.Qty)*100))),
										layout.Flexed(0.125, pg.layoutLabel(statusString(order))),
										layout.Flexed(0.125, pg.layoutLabel(timeSince(order.Stamp))),
									)
								})
							}),
						)
					})
				})
			},
		}

		return sp.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *OrdersHistoryPage) layoutLabel(text string) layout.Widget {
	return pg.Theme.Label(values.TextSize14, text).Layout
}

func (pg *OrdersHistoryPage) Handle() {
}

func (pg *OrdersHistoryPage) OnClose() {
	pg.ctxCancel()
}

func (pg *OrdersHistoryPage) readNotifications() {
	ch := pg.Dexc().Core().NotificationFeed()
	for {
		select {
		case n := <-ch:
			switch n.Type() {
			case core.NoteTypeOrder:
				orderNote := n.(*core.OrderNote)
				for i, order := range pg.orders {
					if orderNote.Order.ID.String() == order.ID.String() {
						pg.orders[i] = orderNote.Order
						pg.RefreshWindow()
						break
					}
				}
			default:
			}
		case <-pg.ctx.Done():
			return
		}
	}
}

func (pg *OrdersHistoryPage) getOrdersHistory() {
	pg.orders = make([]*core.Order, 0)
	ords, err := pg.Dexc().OrderHistory(pg.host)
	if err != nil {
		fmt.Println("Orders error: %w", err)
		return
	}
	pg.orders = ords
}
