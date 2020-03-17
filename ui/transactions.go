package ui

import (
	"fmt"
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrlibwallet"
)

var txsList = layout.List{Axis: layout.Vertical}

func (win *Window) TransactionsPage() {
	bd := func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Flexed(.2, func() {
				layout.Flex{Spacing: layout.SpaceBetween}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.H3("Transactions").Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.outputs.toSend.Layout(win.gtx, &win.inputs.toSend)
					}),
				)
			}),
			layout.Flexed(.8, func() {
				layout.UniformInset(unit.Dp(30)).Layout(win.gtx, func() {
					txsList.Layout(win.gtx, len(win.combined.transactions), func(index int) {
						layout.Inset{Bottom: unit.Dp(8)}.Layout(win.gtx, func() {
							transactionRow(win, index)
						})
					})
				})
			}),
		)
	}

	win.TabbedPage(bd)
}

func transactionRow(win *Window, index int) {
	transaction := win.combined.transactions[index].data.(*dcrlibwallet.Transaction)

	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(win.gtx,
		layout.Rigid(func() {
			layout.Inset{Right: unit.Dp(10)}.Layout(win.gtx, func() {
				txt := fmt.Sprintf("%d", index)
				win.theme.Body1(txt).Layout(win.gtx)
			})
		}),
		layout.Flexed(1, func() {
			layout.Inset{Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(win.gtx, func() {
				layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(22), fmt.Sprintf("%d", transaction.Size)).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(14), fmt.Sprintf("%d DCR", transaction.Amount)).Layout(win.gtx)
					}),
				)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(10)}.Layout(win.gtx, func() {
				win.theme.Body1("Yesterday").Layout(win.gtx)
			})
		}),
	)

	pointer.Rect(image.Rectangle{Max: win.gtx.Dimensions.Size}).Add(win.gtx.Ops)
	click := &win.combined.transactions[index]
	click.gesture.Add(win.gtx.Ops)
}
