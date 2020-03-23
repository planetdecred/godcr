package ui

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var (
	txsList           = layout.List{Axis: layout.Vertical}
	innerColorDanger  = color.RGBA{254, 209, 198, 255}
	innerColorSuccess = color.RGBA{198, 236, 203, 255}
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
	if transaction.info == nil {
		return
	}

	info := transaction.info
	amountTxt1, amountTxt2 := parseAmount(info.Amount)

	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(win.gtx,
		layout.Rigid(func() {
			layout.Stack{}.Layout(win.gtx,
				layout.Stacked(func() {
					layout.Inset{Top: unit.Dp(10)}.Layout(win.gtx, func() {
						background(win, false, info.Direction)
					})
				}),
				layout.Stacked(func() {
					layout.Inset{Left: unit.Dp(5)}.Layout(win.gtx, func() {
						background(win, true, info.Direction)
						layout.Inset{Left: unit.Dp(5), Top: unit.Dp(-2)}.Layout(win.gtx, func() {
							transaction.iconDirection.Color = win.theme.Color.InvText
							transaction.iconDirection.Layout(win.gtx, unit.Dp(16))
						})
					})
				}),
			)
		}),
		layout.Flexed(1, func() {
			layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
				layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(22), amountTxt1).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(14), amountTxt2).Layout(win.gtx)
					}),
				)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(10)}.Layout(win.gtx, func() {
				txt := "Pending"
				if info.Status == "Confirmed" {
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
}

func parseAmount(amount string) (txt1, txt2 string) {
	splitAt := len(amount) - 4

	if strings.Index(amount, ".") > 0 {
		number := strings.Split(amount, " ")[0]
		sufixNumb := strings.Split(number, ".")[1]

		if len(sufixNumb) > 2 {
			splitAt = strings.Index(amount, ".") + 3
		}
	}

	return amount[0:splitAt], amount[splitAt:]
}

func background(win *Window, isFront bool, direction int32) {
	var col color.RGBA

	switch direction {
	case 0:
		col = innerColorDanger
		if isFront {
			col = win.theme.Color.Danger
		}
	case 1, 2:
		col = innerColorSuccess
		if isFront {
			col = win.theme.Color.Success
		}
	}

	d := image.Point{X: 55, Y: 25}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(win.gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(win.gtx.Ops)
	win.gtx.Dimensions = layout.Dimensions{Size: d}
}
