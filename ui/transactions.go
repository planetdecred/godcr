package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
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
				const numRows = 80
				txsList.Layout(win.gtx, numRows, func(index int) {
					transactionRow(win, index)
				})
			}),
		)
	}
	win.TabbedPage(bd)
}

func transactionRow(win *Window, index int) {
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
						txt := fmt.Sprintf("34.17")
						win.theme.Label(unit.Dp(22), txt).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(14), "34888243 DCR").Layout(win.gtx)
					}),
				)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(10)}.Layout(win.gtx, func() {
				txt := fmt.Sprintf("Yesterday")
				win.theme.Body1(txt).Layout(win.gtx)
			})
		}),
	)
}
