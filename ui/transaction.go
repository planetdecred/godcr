package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil/v2"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const (
	pageInset        = 220
	pageContentInset = 15
	rowGroupInset    = 20
)

var (
	transactionPageContainer    = &layout.List{Axis: layout.Vertical}
	transactionInputsContainer  = &layout.List{Axis: layout.Vertical}
	transactionOutputsContainer = &layout.List{Axis: layout.Vertical}
)

func (win *Window) TransactionPage() {
	transaction := *win.walletTransaction

	widgets := []func(){
		func() {
			layout.Inset{Top: unit.Dp(rowGroupInset)}.Layout(win.gtx, func() {
				txnBalanceAndStatus(win, &transaction)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(rowGroupInset)}.Layout(win.gtx, func() {
				txnTypeAndID(win, &transaction)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(rowGroupInset)}.Layout(win.gtx, func() {
				txnInputs(win, &transaction)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(rowGroupInset)}.Layout(win.gtx, func() {
				txnOutputs(win, &transaction)
			})
		},
	}

	win.gtx.Constraints.Height.Max -= pageInset
	win.gtx.Constraints.Height.Min -= pageInset

	win.theme.Surface(win.gtx, func() {
		layout.UniformInset(unit.Dp(pageContentInset)).Layout(win.gtx, func() {
			win.gtx.Constraints.Width.Min = win.gtx.Constraints.Width.Max - pageInset
			win.gtx.Constraints.Width.Max -= pageInset

			layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
				layout.Rigid(func() {
					txnDirection(win, &transaction)
				}),
				layout.Flexed(1, func() {
					transactionPageContainer.Layout(win.gtx, len(widgets), func(i int) {
						layout.Inset{}.Layout(win.gtx, widgets[i])
					})
				}),
				layout.Rigid(func() {
					layout.Center.Layout(win.gtx, func() {
						layout.Inset{Top: unit.Dp(10)}.Layout(win.gtx, func() {
							win.outputs.viewTxnOnDcrdata.Layout(win.gtx, &win.inputs.viewTxnOnDcrdata)
						})
					})
				}),
			)
		})
	})
}

func txnDirection(win *Window, transaction *wallet.Transaction) {
	txt := win.theme.H4(dcrlibwallet.TransactionDirectionName(transaction.Txn.Direction))
	txt.Alignment = text.Middle
	txt.Layout(win.gtx)
	layout.E.Layout(win.gtx, func() {
		win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
	})
}

func txnBalanceAndStatus(win *Window, transaction *wallet.Transaction) {
	cbn := win.combined
	initTxnWidgets(win, transaction, &cbn)

	win.vFlex(
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(-4), Top: unit.Dp(4)}.Layout(win.gtx, func() {
				cbn.transaction.direction.Layout(win.gtx, unit.Dp(28))
			})
			layout.Inset{Left: unit.Dp(28)}.Layout(win.gtx, func() {
				cbn.transaction.amount.TextSize = unit.Dp(28)
				cbn.transaction.amount.Layout(win.gtx)
			})
		}),
		layout.Rigid(func() {
			cbn.transaction.time.Layout(win.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(3)}.Layout(win.gtx, func() {
				cbn.transaction.status.Layout(win.gtx, unit.Dp(16))
			})
			layout.Inset{Left: unit.Dp(18)}.Layout(win.gtx, func() {
				txt := win.theme.Body1(transaction.Status)
				txt.Color = cbn.transaction.status.Color
				txt.Layout(win.gtx)
			})
		}),
		layout.Rigid(func() {
			txt := win.theme.Body1(fmt.Sprintf("%d confirmations", transaction.Confirmations))
			txt.Color = win.theme.Color.Primary
			txt.Layout(win.gtx)
		}),
	)
}

func txnTypeAndID(win *Window, transaction *wallet.Transaction) {
	row := func(label string, t decredmaterial.Label) {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Rigid(func() {
				lb := win.theme.Body1(label)
				lb.Color = win.theme.Color.Hint
				lb.Layout(win.gtx)
			}),
			layout.Rigid(func() {
				t.Layout(win.gtx)
			}),
		)
	}

	layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
		layout.Rigid(func() {
			layout.Inset{Bottom: unit.Dp(rowGroupInset)}.Layout(win.gtx, func() {
				row("Transaction ID", win.theme.Body1(transaction.Txn.Hash))
			})
		}),
		layout.Rigid(func() {
			layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(win.gtx,
				layout.Rigid(func() {
					row("To", win.theme.H6(transaction.WalletName))
				}),
				layout.Rigid(func() {
					row("Included in block", win.theme.H6(fmt.Sprintf("%d", transaction.Txn.BlockHeight)))
				}),
				layout.Rigid(func() {
					row("Type", win.theme.H6(transaction.Txn.Type))
				}),
			)
		}),
	)
}

func txnInputs(win *Window, transaction *wallet.Transaction) {
	layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
		layout.Rigid(func() {
			txt := fmt.Sprintf("%d Inputs consumed", len(transaction.Txn.Inputs))
			txnIORowHeader(win, txt, &win.inputs.toggleTxnDetailsIOs.txnInputs, win.inputs.toggleTxnDetailsIOs.isTxnInputsShow)
		}),
		layout.Rigid(func() {
			if win.inputs.toggleTxnDetailsIOs.isTxnInputsShow {
				transactionInputsContainer.Layout(win.gtx, len(transaction.Txn.Inputs), func(i int) {
					txnIORow(win, dcrutil.Amount(transaction.Txn.Inputs[i].Amount).String(),
						transaction.Txn.Inputs[i].PreviousOutpoint)
				})
			}
		}),
	)
}

func txnOutputs(win *Window, transaction *wallet.Transaction) {
	layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
		layout.Rigid(func() {
			txt := fmt.Sprintf("%d Outputs created", len(transaction.Txn.Outputs))
			txnIORowHeader(win, txt, &win.inputs.toggleTxnDetailsIOs.txnOutputs, win.inputs.toggleTxnDetailsIOs.isTxnOutputsShow)
		}),
		layout.Rigid(func() {
			if win.inputs.toggleTxnDetailsIOs.isTxnOutputsShow {
				transactionOutputsContainer.Layout(win.gtx, len(transaction.Txn.Outputs), func(i int) {
					txnIORow(win, dcrutil.Amount(transaction.Txn.Outputs[i].Amount).String(),
						transaction.Txn.Outputs[i].Address)
				})
			}
		}),
	)
}

func txnIORow(win *Window, amount string, hash string) {
	layout.Inset{Bottom: unit.Dp(5)}.Layout(win.gtx, func() {
		layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Rigid(func() {
				win.theme.Body1(amount).Layout(win.gtx)
			}),
			layout.Rigid(func() {
				txt := win.theme.Body1(hash)
				txt.Color = win.theme.Color.Primary
				txt.Layout(win.gtx)
			}),
		)
	})
}

func txnIORowHeader(win *Window, str string, in *widget.Button, isShow bool) {
	layout.Flex{Spacing: layout.SpaceBetween}.Layout(win.gtx,
		layout.Rigid(func() {
			win.theme.Body1(str).Layout(win.gtx)
		}),
		layout.Rigid(func() {
			if isShow {
				win.outputs.toggleTxnDetailsIOs.expandLess.Layout(win.gtx, in)
				return
			}
			win.outputs.toggleTxnDetailsIOs.expandMore.Layout(win.gtx, in)
		}),
	)
}

func initTxnWidgets(win *Window, transaction *wallet.Transaction, cb *combined) {
	txWidgets := &cb.transaction
	txWidgets.amount = win.theme.Label(unit.Dp(18), transaction.Balance)
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
