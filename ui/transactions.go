package ui

import (
	"image"
	"strconv"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
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
	txsList         = layout.List{Axis: layout.Vertical}
	txsToTxsDetails map[string]*gesture.Click
)

const (
	rowDirectionWidth = .04
	rowDateWidth      = .2
	rowStatusWidth    = .2
	rowAmountWidth    = .3
	rowFeeWidth       = .26

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
	initTotransactionDetailsButtons(win)

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

							directionFilter, _ := strconv.Atoi(win.inputs.transactionFilterDirection.Value(win.gtx))
							txsList.Layout(win.gtx, len(walTxs), func(index int) {
								if directionFilter != 0 && walTxs[index].Txn.Direction != int32(directionFilter-1) {
									return
								}
								renderTxsRow(win, walTxs[index])

								click := txsToTxsDetails[walTxs[index].Txn.Hash]
								pointer.Rect(image.Rectangle{Max: win.gtx.Dimensions.Size}).Add(win.gtx.Ops)
								click.Add(win.gtx.Ops)
								toDetails(win, &walTxs[index], click)
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
	txt.Alignment = text.Middle

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

func renderTxsRow(win *Window, transaction wallet.Transaction) {
	cbn := win.combined
	initTxnWidgets(win, &transaction, &cbn)
	layout.Inset{Bottom: unit.Dp(txnPageInsetBottom)}.Layout(win.gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(win.gtx,
			layout.Flexed(rowDirectionWidth, func() {
				layout.Inset{Top: unit.Dp(3)}.Layout(win.gtx, func() {
					cbn.transaction.direction.Layout(win.gtx, unit.Dp(16))
				})
			}),
			layout.Flexed(rowDateWidth, func() {
				cbn.transaction.time.Alignment = text.Middle
				cbn.transaction.time.Layout(win.gtx)
			}),
			layout.Flexed(rowStatusWidth, func() {
				txt := win.theme.Body1(transaction.Status)
				txt.Alignment = text.Middle
				txt.Layout(win.gtx)
			}),
			layout.Flexed(rowAmountWidth, func() {
				cbn.transaction.amount.Alignment = text.End
				cbn.transaction.amount.Layout(win.gtx)
			}),
			layout.Flexed(rowFeeWidth, func() {
				txt := win.theme.Body1(dcrutil.Amount(transaction.Txn.Fee).String())
				txt.Alignment = text.End
				txt.Layout(win.gtx)
			}),
		)
	})
}

func initTxnWidgets(win *Window, transaction *wallet.Transaction, cb *combined) {
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

func initTotransactionDetailsButtons(win *Window) {
	if win.walletTransactions.Total != len(txsToTxsDetails) {
		txsToTxsDetails = make(map[string]*gesture.Click, win.walletTransactions.Total)

		for _, walTxs := range win.walletTransactions.Txs {
			for _, txn := range walTxs {
				txsToTxsDetails[txn.Txn.Hash] = &gesture.Click{}
			}
		}
	}
}

func toDetails(win *Window, txn *wallet.Transaction, click *gesture.Click) {
	for _, e := range click.Events(win.gtx) {
		if e.Type == gesture.TypeClick {
			win.walletTransaction = txn
			win.current = PageTransactionDetails
		}
	}
}
