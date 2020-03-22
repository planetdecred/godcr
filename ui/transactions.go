package ui

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/wallet"
)

var (
	txsList = layout.List{Axis: layout.Vertical}
)

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
			layout.Flexed(.2, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(win.gtx,
					layout.Rigid(func() {
						layout.Inset{Right: unit.Dp(10)}.Layout(win.gtx, func() {
							win.combined.transactionSort.Layout(win.gtx, func() {
							})
						})
					}),
					layout.Rigid(func() {
						win.combined.transactionStatus.Layout(win.gtx, func() {
						})
					}),
				)
			}),
			layout.Flexed(.6, func() {
				layout.UniformInset(unit.Dp(30)).Layout(win.gtx, func() {
					txsList.Layout(win.gtx, len(win.combined.transactions), func(index int) {
						layout.Inset{Bottom: unit.Dp(15)}.Layout(win.gtx, func() {
							renderTxsRow(win, index)
						})
					})
				})
			}),
		)
	}

	win.TabbedPage(bd)
}

func renderTxsRow(win *Window, index int) {
	transaction := &win.combined.transactions[index]
	info := transaction.data.(*wallet.TransactionInfo)

	if info == nil {
		return
	}

	indexOfDot := strings.Index(info.Amount, ".")

	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(win.gtx,
		layout.Rigid(func() {
			layout.Stack{}.Layout(win.gtx,
				layout.Stacked(func() {
					layout.Inset{Top: unit.Dp(10)}.Layout(win.gtx, func() {
						background(win.gtx, transaction.iconDirection.innerColor)
					})
				}),
				layout.Stacked(func() {
					layout.Inset{Left: unit.Dp(5)}.Layout(win.gtx, func() {
						background(win.gtx, transaction.iconDirection.backgroundColor)
						layout.Inset{Left: unit.Dp(10), Top: unit.Dp(-5)}.Layout(win.gtx, func() {
							txt := win.theme.Label(unit.Dp(16), transaction.iconDirection.icon)
							txt.Color = win.theme.Color.InvText
							txt.Layout(win.gtx)
						})
					})
				}),
			)
		}),
		layout.Flexed(1, func() {
			layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
				layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(22), info.Amount[0:indexOfDot+3]).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(14), info.Amount[indexOfDot+3:]).Layout(win.gtx)
					}),
				)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(10)}.Layout(win.gtx, func() {
				txt := "Pending"
				if info.Status == "confirmed" {
					txt = info.Datetime
				}
				win.theme.Body1(txt).Layout(win.gtx)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Bottom: unit.Dp(15), Left: unit.Dp(8)}.Layout(win.gtx, func() {
				transaction.iconStatus.Layout(win.gtx, unit.Dp(16))
			})
		}),
	)

	pointer.Rect(image.Rectangle{Max: win.gtx.Dimensions.Size}).Add(win.gtx.Ops)
	transaction.gesture.Add(win.gtx.Ops)
}

func background(gtx *layout.Context, col color.RGBA) {
	d := image.Point{X: 55, Y: 25}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
