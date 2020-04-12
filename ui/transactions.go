package ui

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"github.com/raedahgroup/godcr-gio/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	txsList = layout.List{Axis: layout.Vertical}
)

const (
	rowDirectionWidth  = .1
	rowDateWidth       = .2
	rowStatusWidth     = .2
	rowAmountWidth     = .3
	rowFeeWidth        = .2
	txnRowLabelSize    = 16
	txnPageInsetTop    = 15
	txnPageInsetLeft   = 15
	txnPageInsetRight  = 15
	txnPageInsetBottom = 15
)

func (win *Window) TransactionsPage() {
	if win.walletInfo.LoadedWallets == 0 {
		win.Page(func() {
			win.outputs.noWallet.Layout(win.gtx)
		})
		return
	}
	bd := func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(txnPageInsetTop)}.Layout(win.gtx, func() {
					layout.Flex{Spacing: layout.SpaceBetween}.Layout(win.gtx,
						layout.Rigid(func() {
							layout.Inset{Left: unit.Dp(txnPageInsetLeft)}.Layout(win.gtx, func() {
								renderFiltererButton(win)
							})
						}),
						layout.Rigid(func() {
							win.outputs.toSend.Layout(win.gtx, &win.inputs.toSend)
						}),
						layout.Rigid(func() {
							layout.Inset{Right: unit.Dp(txnPageInsetRight)}.Layout(win.gtx, func() {
								win.outputs.toReceive.Layout(win.gtx, &win.inputs.toReceive)
							})
						}),
					)
				})
			}),
			layout.Flexed(1, func() {
				layout.Inset{Left: unit.Dp(txnPageInsetLeft), Right: unit.Dp(txnPageInsetRight)}.Layout(win.gtx, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
						layout.Rigid(func() {
							layout.Inset{Top: unit.Dp(txnPageInsetTop), Bottom: unit.Dp(txnPageInsetBottom)}.Layout(win.gtx, func() {
								rowHeader(win)
							})
						}),
						layout.Flexed(1, func() {
							walletID := win.walletInfo.Wallets[win.selected].ID
							walTxs := win.walletTransactions.Txs[walletID]

							if len(walTxs) == 0 {
								txt := win.theme.Body1("No transactions")
								txt.Alignment = text.Middle
								txt.Layout(win.gtx)
								return
							}

							var txs []wallet.TransactionInfo
							directionFilter, _ := strconv.Atoi(win.inputs.transactionFilterDirection.Value(win.gtx))

							for _, txn := range walTxs {
								if directionFilter != 0 && txn.Txn.Direction != int32(directionFilter-1) {
									continue
								}
								txs = append(txs, txn)
							}

							txsList.Layout(win.gtx, len(txs), func(index int) {
								layout.Inset{Bottom: unit.Dp(txnPageInsetBottom)}.Layout(win.gtx, func() {
									renderTxsRow(win, txs[index])
								})
							})
						}),
					)
				})
			}),
		)
	}

	win.TabbedPage(bd)
}

func renderFiltererButton(win *Window) {
	var button decredmaterial.IconButton

	if win.inputs.transactionFilterSort.Value(win.gtx) == "0" {
		button = win.outputs.toTransactionsFilters.sortNewest
	} else {
		button = win.outputs.toTransactionsFilters.sortOldest
	}

	switch win.inputs.transactionFilterDirection.Value(win.gtx) {
	case "0":
		button.Background = win.theme.Color.Primary
	case "1":
		button.Background = win.theme.Color.Danger
	case "2":
		button.Background = win.theme.Color.Success
	case "3":
		button.Background = win.theme.Color.Hint
	default:
		button.Background = win.theme.Color.Hint
	}

	button.Layout(win.gtx, &win.inputs.toTransactionsFilters)
}

func rowHeader(win *Window) {
	txt := win.theme.Label(unit.Dp(txnRowLabelSize), "#")
	txt.Color = win.theme.Color.Hint
	layout.Flex{Axis: layout.Horizontal}.Layout(win.gtx,
		layout.Flexed(rowDirectionWidth, func() {
			txt.Layout(win.gtx)
		}),
		layout.Flexed(rowDateWidth, func() {
			txt.Text = "Date (UTC)"
			txt.Layout(win.gtx)
		}),
		layout.Flexed(rowStatusWidth, func() {
			txt.Text = "Status"
			txt.Layout(win.gtx)
		}),
		layout.Flexed(rowAmountWidth, func() {
			txt.Text = "Amount"
			txt.Layout(win.gtx)
		}),
		layout.Flexed(rowFeeWidth, func() {
			txt.Text = "Fee"
			txt.Layout(win.gtx)
		}),
	)
}

func renderTxsRow(win *Window, transaction wallet.TransactionInfo) {
	cbn := win.combined
	initTxnWidgets(win, &transaction, &cbn)
	layout.Flex{Axis: layout.Horizontal}.Layout(win.gtx,
		layout.Flexed(rowDirectionWidth, func() {
			layout.Inset{Top: unit.Dp(3)}.Layout(win.gtx, func() {
				cbn.transaction.direction.Layout(win.gtx, unit.Dp(16))
			})
		}),
		layout.Flexed(rowDateWidth, func() {
			cbn.transaction.time.Layout(win.gtx)
		}),
		layout.Flexed(rowStatusWidth, func() {
			win.theme.Body1(transaction.Status).Layout(win.gtx)
		}),
		layout.Flexed(rowAmountWidth, func() {
			cbn.transaction.amount.Layout(win.gtx)
		}),
		layout.Flexed(rowFeeWidth, func() {
			win.theme.Body1(dcrutil.Amount(transaction.Txn.Fee).String()).Layout(win.gtx)
		}),
	)
}

func initTxnWidgets(win *Window, transaction *wallet.TransactionInfo, cb *combined) {
	txWidgets := &cb.transaction
	txWidgets.amount = win.theme.Label(unit.Dp(16), transaction.Balance)
	txWidgets.time = win.theme.Body1("Pending")

	if transaction.Status == "confirmed" {
		txWidgets.time.Text = dcrlibwallet.ExtractDateOrTime(transaction.Txn.Timestamp)
		txWidgets.status, _ = decredmaterial.NewIcon(icons.ActionCheckCircle)
		txWidgets.status.Color = win.theme.Color.Success
	} else {
		txWidgets.status, _ = decredmaterial.NewIcon(icons.ToggleRadioButtonUnchecked)
	}

	if transaction.Txn.Direction == dcrlibwallet.TxDirectionSent {
		txWidgets.direction, _ = decredmaterial.NewIcon(icons.ContentRemove)
		txWidgets.direction.Color = win.theme.Color.Danger
	} else {
		txWidgets.direction, _ = decredmaterial.NewIcon(icons.ContentAdd)
		txWidgets.direction.Color = win.theme.Color.Success
	}
}
