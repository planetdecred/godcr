package ui

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/planetdecred/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/wallet"
)

const PageTransactionDetails = "TransactionDetails"

type transactionPage struct {
	transactionPageContainer    layout.List
	transactionInputsContainer  layout.List
	transactionOutputsContainer layout.List
	backButton                  decredmaterial.IconButton
	txnInfo                     **wallet.Transaction
	viewTxnOnDcrdata            decredmaterial.Button

	outputsCollapsible *decredmaterial.Collapsible
	inputsCollapsible  *decredmaterial.Collapsible
}

func (win *Window) TransactionPage(common pageCommon) layout.Widget {
	pg := &transactionPage{
		transactionPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		transactionInputsContainer: layout.List{
			Axis: layout.Vertical,
		},
		transactionOutputsContainer: layout.List{
			Axis: layout.Vertical,
		},
		txnInfo: &win.walletTransaction,

		outputsCollapsible: common.theme.Collapsible(),
		inputsCollapsible:  common.theme.Collapsible(),

		backButton:       common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		viewTxnOnDcrdata: common.theme.Button(new(widget.Clickable), "View on dcrdata"),
	}
	pg.backButton.Color = common.theme.Color.Hint
	pg.backButton.Size = values.MarginPadding30

	return func(gtx C) D {
		pg.Handler(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *transactionPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	margin := values.MarginPadding20

	widgets := []func(gtx C) D{
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return pg.txnBalanceAndStatus(gtx, &common)
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return pg.txnTypeAndID(gtx, &common)
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return pg.txnInputs(gtx, &common)
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return pg.txnOutputs(gtx, &common)
			})
		},
	}

	return common.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.header(gtx, &common)
			}),
			layout.Flexed(1, func(gtx C) D {
				if *pg.txnInfo == nil {
					return layout.Dimensions{}
				}
				return pg.transactionPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
					return layout.Inset{}.Layout(gtx, widgets[i])
				})
			}),
			layout.Rigid(func(gtx C) D {
				if *pg.txnInfo == nil {
					return layout.Dimensions{}
				}
				return pg.viewTxnOnDcrdata.Layout(gtx)
			}),
		)
	})
}

func (pg *transactionPage) header(gtx layout.Context, common *pageCommon) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.W.Layout(gtx, func(gtx C) D {
				return pg.backButton.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			txt := common.theme.H4("")
			if *pg.txnInfo != nil {
				txt.Text = dcrlibwallet.TransactionDirectionName((*pg.txnInfo).Txn.Direction)
			} else {
				txt.Text = "Not found"
			}

			txt.Alignment = text.Middle
			return txt.Layout(gtx)
		}),
	)
}

func (pg *transactionPage) txnBalanceAndStatus(gtx layout.Context, common *pageCommon) layout.Dimensions {
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, *pg.txnInfo, &txnWidgets)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return txnWidgets.direction.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
						txnWidgets.amount.TextSize = values.TextSize28
						return txnWidgets.amount.Layout(gtx)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return txnWidgets.time.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return txnWidgets.status.Layout(gtx, values.MarginPadding15)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						txt := common.theme.Body1((*pg.txnInfo).Status)
						txt.Color = txnWidgets.status.Color
						return txt.Layout(gtx)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			txt := common.theme.Body1(fmt.Sprintf("%d confirmations", (*pg.txnInfo).Confirmations))
			txt.Color = common.theme.Color.Primary
			return txt.Layout(gtx)
		}),
	)
}

func (pg *transactionPage) txnTypeAndID(gtx layout.Context, common *pageCommon) layout.Dimensions {
	transaction := *pg.txnInfo

	column := func(gtx layout.Context, label string, t decredmaterial.Label) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				lb := common.theme.Body1(label)
				lb.Color = common.theme.Color.Hint
				return lb.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return t.Layout(gtx)
			}),
		)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
				return column(gtx, "Transaction ID", common.theme.Body1(transaction.Txn.Hash))
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return column(gtx, "To", common.theme.H6(transaction.WalletName))
				}),
				layout.Rigid(func(gtx C) D {
					txt := ""
					if transaction.Txn.BlockHeight != -1 {
						txt = fmt.Sprintf("%d", transaction.Txn.BlockHeight)
					}
					return column(gtx, "Included in block", common.theme.H6(txt))
				}),
				layout.Rigid(func(gtx C) D {
					return column(gtx, "Type", common.theme.H6(transaction.Txn.Type))
				}),
			)
		}),
	)
}

func (pg *transactionPage) txnInputs(gtx layout.Context, common *pageCommon) layout.Dimensions {
	transaction := *pg.txnInfo

	collapsibleHeader := func(gtx C) D {
		return common.theme.Body1(fmt.Sprintf("%d Inputs consumed", len(transaction.Txn.Inputs))).Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		return pg.transactionInputsContainer.Layout(gtx, len(transaction.Txn.Inputs), func(gtx C, i int) D {
			return pg.txnIORow(gtx, common, dcrutil.Amount(transaction.Txn.Inputs[i].Amount).String(),
				transaction.Txn.Inputs[i].PreviousOutpoint)
		})
	}

	return pg.inputsCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
}

func (pg *transactionPage) txnOutputs(gtx layout.Context, common *pageCommon) layout.Dimensions {
	transaction := *pg.txnInfo

	collapsibleHeader := func(gtx C) D {
		return common.theme.Body1(fmt.Sprintf("%d Outputs created", len(transaction.Txn.Outputs))).Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		return pg.transactionOutputsContainer.Layout(gtx, len(transaction.Txn.Outputs), func(gtx C, i int) D {
			return pg.txnIORow(gtx, common, dcrutil.Amount(transaction.Txn.Outputs[i].Amount).String(),
				transaction.Txn.Outputs[i].Address)
		})
	}

	return pg.outputsCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
}

func (pg *transactionPage) txnIORow(gtx layout.Context, common *pageCommon, amount string, hash string) layout.Dimensions {

	return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return common.theme.Label(values.MarginPadding15, amount).Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := common.theme.Label(values.MarginPadding15, hash)
				txt.Color = common.theme.Color.Primary
				return txt.Layout(gtx)
			}),
		)
	})
}

func (pg *transactionPage) viewTxnOnBrowser(common *pageCommon) {
	var err error
	url := common.wallet.GetBlockExplorerURL((*pg.txnInfo).Txn.Hash)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Error(err)
	}
}

func (pg *transactionPage) Handler(common pageCommon) {
	if pg.viewTxnOnDcrdata.Button.Clicked() {
		pg.viewTxnOnBrowser(&common)
	}
	if pg.backButton.Button.Clicked() {
		switch common.navTab.Selected {
		case 0:
			*common.page = PageOverview
			return
		default:
			*common.page = PageTransactions
			return
		}
	}
}
