package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"github.com/raedahgroup/godcr-gio/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	txsList = layout.List{Axis: layout.Vertical}
)

const (
	pageHeadHeight      = 0.2
	pageContainerHeight = .8
	filterButtonWidth   = 200
)

func (win *Window) TransactionsPage() {
	bd := func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Flexed(pageHeadHeight, func() {
				layout.Flex{Spacing: layout.SpaceBetween}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.H3("Transactions").Layout(win.gtx)
					}),
					layout.Rigid(func() {
						win.outputs.toSend.Layout(win.gtx, &win.inputs.toSend)
					}),
					layout.Rigid(func() {
						layout.Inset{Right: unit.Dp(20)}.Layout(win.gtx, func() {
							win.outputs.toReceive.Layout(win.gtx, &win.inputs.toReceive)
						})
					}),
				)
			}),
			layout.Flexed(pageContainerHeight, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(win.gtx,
					layout.Rigid(func() {
						win.gtx.Constraints.Width.Min = filterButtonWidth
						layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
							layout.Rigid(func() {
								win.combined.transactionSort.Layout(win.gtx, func() {
								})
							}),
							layout.Rigid(func() {
								win.combined.transactionStatus.Layout(win.gtx, func() {
								})
							}),
						)
					}),
					layout.Flexed(1, func() {
						layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(win.gtx, func() {
							txsList.Layout(win.gtx, len(*win.transactionsWallet), func(index int) {
								layout.Inset{Bottom: unit.Dp(15)}.Layout(win.gtx, func() {
									renderTxsRow(win, &(*win.transactionsWallet)[index])
								})
							})
						})
					}),
				)
			}),
		)
	}

	win.TabbedPage(bd)
}

func renderTxsRow(win *Window, transaction *wallet.TransactionInfo) {
	txWidgets := win.combined.transaction
	txWidgets.amount = win.theme.Label(unit.Dp(22), transaction.Amount)
	txWidgets.time = win.theme.Body1("Pending")

	if transaction.Status == "Confirmed" {
		txWidgets.time.Text = dcrlibwallet.ExtractDateOrTime(transaction.Timestamp)
		txWidgets.status, _ = decredmaterial.NewIcon(icons.ActionCheckCircle)
		txWidgets.status.Color = win.theme.Color.Success
	} else {
		txWidgets.status, _ = decredmaterial.NewIcon(icons.ToggleRadioButtonUnchecked)
	}

	if transaction.Direction == wallet.TxDirectionSent {
		txWidgets.direction, _ = decredmaterial.NewIcon(icons.ContentRemove)
		txWidgets.direction.Color = win.theme.Color.Danger
	} else {
		txWidgets.direction, _ = decredmaterial.NewIcon(icons.ContentAdd)
		txWidgets.direction.Color = win.theme.Color.Success
	}

	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(win.gtx,
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(5), Top: unit.Dp(-6)}.Layout(win.gtx, func() {
				txWidgets.direction.Layout(win.gtx, unit.Dp(16))
			})
		}),
		layout.Flexed(1, func() {
			layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
				txWidgets.amount.Layout(win.gtx)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(10)}.Layout(win.gtx, func() {
				txWidgets.time.Layout(win.gtx)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Bottom: unit.Dp(15), Left: unit.Dp(8)}.Layout(win.gtx, func() {
				txWidgets.status.Layout(win.gtx, unit.Dp(16))
			})
		}),
	)
}
