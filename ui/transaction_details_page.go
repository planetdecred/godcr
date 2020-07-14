package ui

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

const PageTransactionDetails = "transactiondetails"

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
	page := &transactionPage{
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
	page.backButton.Color = common.theme.Color.Hint
	page.backButton.Size = values.MarginPadding30

	return func(gtx C) D {
		page.Handler(common)
		return page.Layout(common)
	}
}

func (page *transactionPage) Layout(common pageCommon) layout.Dimensions {
	gtx := common.gtx
	margin := values.MarginPadding20
	widgets := []func(gtx C) D{
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return page.txnBalanceAndStatus(&common)
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return page.txnTypeAndID(&common)
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return page.txnInputs(&common)
			})
		},
		func(gtx C) D {
			return layout.Inset{Top: margin}.Layout(gtx, func(gtx C) D {
				return page.txnOutputs(&common)
			})
		},
	}

	return common.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return page.header(&common)
			}),
			layout.Flexed(1, func(gtx C) D {
				if *page.txnInfo == nil {
					return layout.Dimensions{}
				}
				return page.transactionPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
					return layout.Inset{}.Layout(gtx, widgets[i])
				})
			}),
			layout.Rigid(func(gtx C) D {
				if *page.txnInfo == nil {
					return layout.Dimensions{}
				}
				return page.viewTxnOnDcrdata.Layout(gtx)
			}),
		)
	})
}

func (page *transactionPage) header(common *pageCommon) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(common.gtx,
		layout.Rigid(func(gtx C) D {
			return layout.W.Layout(common.gtx, func(gtx C) D {
				return page.backButton.Layout(common.gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			txt := common.theme.H4("")
			if *page.txnInfo != nil {
				txt.Text = dcrlibwallet.TransactionDirectionName((*page.txnInfo).Txn.Direction)
			} else {
				txt.Text = "Not found"
			}

			txt.Alignment = text.Middle
			return txt.Layout(common.gtx)
		}),
	)
}

func (page *transactionPage) txnBalanceAndStatus(common *pageCommon) layout.Dimensions {
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, *page.txnInfo, &txnWidgets)

	return layout.Flex{Axis: layout.Vertical}.Layout(common.gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding5}.Layout(common.gtx, func(gtx C) D {
				return txnWidgets.direction.Layout(common.gtx, values.MarginPadding30)
			})
			return layout.Inset{Left: values.MarginPadding30}.Layout(common.gtx, func(gtx C) D {
				txnWidgets.amount.TextSize = values.TextSize28
				return txnWidgets.amount.Layout(common.gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return txnWidgets.time.Layout(common.gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding5}.Layout(common.gtx, func(gtx C) D {
				return txnWidgets.status.Layout(common.gtx, values.MarginPadding15)
			})
			return layout.Inset{Left: values.MarginPadding20}.Layout(common.gtx, func(gtx C) D {
				txt := common.theme.Body1((*page.txnInfo).Status)
				txt.Color = txnWidgets.status.Color
				return txt.Layout(common.gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			txt := common.theme.Body1(fmt.Sprintf("%d confirmations", (*page.txnInfo).Confirmations))
			txt.Color = common.theme.Color.Primary
			return txt.Layout(common.gtx)
		}),
	)
}

func (page *transactionPage) txnTypeAndID(common *pageCommon) layout.Dimensions {
	transaction := *page.txnInfo
	gtx := common.gtx
	row := func(label string, t decredmaterial.Label) layout.Dimensions {
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
				return row("Transaction ID", common.theme.Body1(transaction.Txn.Hash))
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return row("To", common.theme.H6(transaction.WalletName))
				}),
				layout.Rigid(func(gtx C) D {
					txt := ""
					if transaction.Txn.BlockHeight != -1 {
						txt = fmt.Sprintf("%d", transaction.Txn.BlockHeight)
					}
					return row("Included in block", common.theme.H6(txt))
				}),
				layout.Rigid(func(gtx C) D {
					return row("Type", common.theme.H6(transaction.Txn.Type))
				}),
			)
		}),
	)
}

func (page *transactionPage) txnInputs(common *pageCommon) layout.Dimensions {
	transaction := *page.txnInfo
	collapsibleHeader := func(gtx C) D {
		return common.theme.Body1(fmt.Sprintf("%d Inputs consumed", len(transaction.Txn.Inputs))).Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		return page.transactionInputsContainer.Layout(common.gtx, len(transaction.Txn.Inputs), func(gtx C, i int) D {
			return page.txnIORow(common, dcrutil.Amount(transaction.Txn.Inputs[i].Amount).String(),
				transaction.Txn.Inputs[i].PreviousOutpoint)
		})
	}

	return page.inputsCollapsible.Layout(common.gtx, collapsibleHeader, collapsibleBody)
}

func (page *transactionPage) txnOutputs(common *pageCommon) layout.Dimensions {
	transaction := *page.txnInfo
	collapsibleHeader := func(gtx C) D {
		return common.theme.Body1(fmt.Sprintf("%d Outputs created", len(transaction.Txn.Outputs))).Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		return page.transactionOutputsContainer.Layout(common.gtx, len(transaction.Txn.Outputs), func(gtx C, i int) D {
			return page.txnIORow(common, dcrutil.Amount(transaction.Txn.Outputs[i].Amount).String(),
				transaction.Txn.Outputs[i].Address)
		})
	}

	return page.outputsCollapsible.Layout(common.gtx, collapsibleHeader, collapsibleBody)
}

func (page *transactionPage) txnIORow(common *pageCommon, amount string, hash string) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding5}.Layout(common.gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(common.gtx,
			layout.Rigid(func(gtx C) D {
				return common.theme.Label(values.MarginPadding15, amount).Layout(common.gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := common.theme.Label(values.MarginPadding15, hash)
				txt.Color = common.theme.Color.Primary
				return txt.Layout(common.gtx)
			}),
		)
	})
}

func (page *transactionPage) viewTxnOnBrowser(common *pageCommon) {
	var err error
	url := common.wallet.GetBlockExplorerURL((*page.txnInfo).Txn.Hash)

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

func (page *transactionPage) Handler(common pageCommon) {
	if page.viewTxnOnDcrdata.Button.Clicked() {
		page.viewTxnOnBrowser(&common)
	}
	if page.backButton.Button.Clicked() {
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
