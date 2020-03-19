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
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	txsList        = layout.List{Axis: layout.Vertical}
	icoReceived, _ = decredmaterial.NewIcon(icons.ContentAdd)
	icoSent, _     = decredmaterial.NewIcon(icons.ContentRemove)
	icoConfirm, _  = decredmaterial.NewIcon(icons.ActionCheckCircle)
	icoPending, _  = decredmaterial.NewIcon(icons.ToggleRadioButtonUnchecked)
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
			layout.Flexed(.8, func() {
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
	transaction := win.combined.transactions[index].data.(*dcrlibwallet.Transaction)
	if transaction == nil {
		return
	}

	amountStr := dcrutil.Amount(transaction.Amount).String()
	indexOfDot := strings.Index(amountStr, ".")

	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(win.gtx,
		layout.Rigid(func() {
			renderTxsDirection(win, transaction.Direction)
		}),
		layout.Flexed(1, func() {
			layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
				layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(22), amountStr[0:indexOfDot+3]).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(14), amountStr[indexOfDot+3:]).Layout(win.gtx)
					}),
				)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(10)}.Layout(win.gtx, func() {
				win.theme.Body1(dcrlibwallet.ExtractDateOrTime(transaction.Timestamp)).Layout(win.gtx)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Bottom: unit.Dp(15), Left: unit.Dp(8)}.Layout(win.gtx, func() {
				renderTxsStatus(win, transaction)
			})
		}),
	)

	pointer.Rect(image.Rectangle{Max: win.gtx.Dimensions.Size}).Add(win.gtx.Ops)
	click := &win.combined.transactions[index]
	click.gesture.Add(win.gtx.Ops)
}

func renderTxsStatus(win *Window, txn *dcrlibwallet.Transaction) {
	confirmations := win.walletInfo.BestBlockHeight - txn.BlockHeight + 1
	if txn.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		icoConfirm.Color = win.theme.Color.Success
		icoConfirm.Layout(win.gtx, unit.Dp(16))

		return
	}

	icoPending.Layout(win.gtx, unit.Dp(16))
}

func renderTxsDirection(win *Window, direction int32) {
	innerColor := color.RGBA{254, 209, 198, 255}
	frontColor := win.theme.Color.Danger
	icon := icoSent

	if direction == 1 || direction == 2 {
		innerColor = color.RGBA{198, 236, 203, 255}
		frontColor = win.theme.Color.Success
		icon = icoReceived
	}

	icon.Color = win.theme.Color.InvText

	layout.Stack{}.Layout(win.gtx,
		layout.Stacked(func() {
			layout.Inset{Top: unit.Dp(10)}.Layout(win.gtx, func() {
				background(win.gtx, innerColor)
			})
		}),
		layout.Stacked(func() {
			layout.Inset{Left: unit.Dp(5)}.Layout(win.gtx, func() {
				background(win.gtx, frontColor)
				layout.Inset{Left: unit.Dp(5), Top: unit.Dp(-2)}.Layout(win.gtx, func() {
					icon.Layout(win.gtx, unit.Dp(16))
				})
			})
		}),
	)
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
