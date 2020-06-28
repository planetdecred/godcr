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
	backButtonW                 decredmaterial.IconButton
	txnInfo                     **wallet.Transaction
	viewTxnOnDcrdataW,
	backButton widget.Button
	viewTxnOnDcrdata decredmaterial.Button

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

		backButtonW:      common.theme.PlainIconButton(common.icons.navigationArrowBack),
		viewTxnOnDcrdata: common.theme.Button("View on dcrdata"),
	}
	page.backButtonW.Color = common.theme.Color.Hint
	page.backButtonW.Size = values.MarginPadding30

	return func() {
		page.Layout(common)
		page.Handler(common)
	}
}

func (page *transactionPage) Layout(common pageCommon) {
	gtx := common.gtx
	margin := values.MarginPadding20
	widgets := []func(){
		func() {
			layout.Inset{Top: margin}.Layout(gtx, func() {
				page.txnBalanceAndStatus(&common)
			})
		},
		func() {
			layout.Inset{Top: margin}.Layout(gtx, func() {
				page.txnTypeAndID(&common)
			})
		},
		func() {
			layout.Inset{Top: margin}.Layout(gtx, func() {
				page.txnInputs(&common)
			})
		},
		func() {
			layout.Inset{Top: margin}.Layout(gtx, func() {
				page.txnOutputs(&common)
			})
		},
	}

	common.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				page.header(&common)
			}),
			layout.Flexed(1, func() {
				if *page.txnInfo == nil {
					return
				}
				page.transactionPageContainer.Layout(gtx, len(widgets), func(i int) {
					layout.Inset{}.Layout(gtx, widgets[i])
				})
			}),
			layout.Rigid(func() {
				if *page.txnInfo == nil {
					return
				}
				page.viewTxnOnDcrdata.Layout(gtx, &page.viewTxnOnDcrdataW)
			}),
		)
	})
}

func (page *transactionPage) header(common *pageCommon) {
	layout.W.Layout(common.gtx, func() {
		page.backButtonW.Layout(common.gtx, &page.backButton)
	})
	txt := common.theme.H4("")
	if *page.txnInfo != nil {
		txt.Text = dcrlibwallet.TransactionDirectionName((*page.txnInfo).Txn.Direction)
	} else {
		txt.Text = "Not found"
	}

	txt.Alignment = text.Middle
	txt.Layout(common.gtx)
}

func (page *transactionPage) txnBalanceAndStatus(common *pageCommon) {
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, *page.txnInfo, &txnWidgets)

	layout.Flex{Axis: layout.Vertical}.Layout(common.gtx,
		layout.Rigid(func() {
			layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding5}.Layout(common.gtx, func() {
				txnWidgets.direction.Layout(common.gtx, values.MarginPadding30)
			})
			layout.Inset{Left: values.MarginPadding30}.Layout(common.gtx, func() {
				txnWidgets.amount.TextSize = values.TextSize28
				txnWidgets.amount.Layout(common.gtx)
			})
		}),
		layout.Rigid(func() {
			txnWidgets.time.Layout(common.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Top: values.MarginPadding5}.Layout(common.gtx, func() {
				txnWidgets.status.Layout(common.gtx, values.MarginPadding15)
			})
			layout.Inset{Left: values.MarginPadding20}.Layout(common.gtx, func() {
				txt := common.theme.Body1((*page.txnInfo).Status)
				txt.Color = txnWidgets.status.Color
				txt.Layout(common.gtx)
			})
		}),
		layout.Rigid(func() {
			txt := common.theme.Body1(fmt.Sprintf("%d confirmations", (*page.txnInfo).Confirmations))
			txt.Color = common.theme.Color.Primary
			txt.Layout(common.gtx)
		}),
	)
}

func (page *transactionPage) txnTypeAndID(common *pageCommon) {
	transaction := *page.txnInfo
	gtx := common.gtx
	row := func(label string, t decredmaterial.Label) {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				lb := common.theme.Body1(label)
				lb.Color = common.theme.Color.Hint
				lb.Layout(gtx)
			}),
			layout.Rigid(func() {
				t.Layout(gtx)
			}),
		)
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, func() {
				row("Transaction ID", common.theme.Body1(transaction.Txn.Hash))
			})
		}),
		layout.Rigid(func() {
			layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func() {
					row("To", common.theme.H6(transaction.WalletName))
				}),
				layout.Rigid(func() {
					txt := ""
					if transaction.Txn.BlockHeight != -1 {
						txt = fmt.Sprintf("%d", transaction.Txn.BlockHeight)
					}
					row("Included in block", common.theme.H6(txt))
				}),
				layout.Rigid(func() {
					row("Type", common.theme.H6(transaction.Txn.Type))
				}),
			)
		}),
	)
}

func (page *transactionPage) txnInputs(common *pageCommon) {
	transaction := *page.txnInfo
	collapsibleHeader := func(gtx *layout.Context) {
		common.theme.Body1(fmt.Sprintf("%d Inputs consumed", len(transaction.Txn.Inputs))).Layout(gtx)
	}

	collapsibleBody := func(gtx *layout.Context) {
		page.transactionInputsContainer.Layout(common.gtx, len(transaction.Txn.Inputs), func(i int) {
			page.txnIORow(common, dcrutil.Amount(transaction.Txn.Inputs[i].Amount).String(),
				transaction.Txn.Inputs[i].PreviousOutpoint)
		})
	}

	page.inputsCollapsible.Layout(common.gtx, collapsibleHeader, collapsibleBody)
}

func (page *transactionPage) txnOutputs(common *pageCommon) {
	transaction := *page.txnInfo
	collapsibleHeader := func(gtx *layout.Context) {
		common.theme.Body1(fmt.Sprintf("%d Outputs created", len(transaction.Txn.Outputs))).Layout(gtx)
	}

	collapsibleBody := func(gtx *layout.Context) {
		page.transactionOutputsContainer.Layout(common.gtx, len(transaction.Txn.Outputs), func(i int) {
			page.txnIORow(common, dcrutil.Amount(transaction.Txn.Outputs[i].Amount).String(),
				transaction.Txn.Outputs[i].Address)
		})
	}

	page.outputsCollapsible.Layout(common.gtx, collapsibleHeader, collapsibleBody)
}

func (page *transactionPage) txnIORow(common *pageCommon, amount string, hash string) {
	layout.Inset{Bottom: values.MarginPadding5}.Layout(common.gtx, func() {
		layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(common.gtx,
			layout.Rigid(func() {
				common.theme.Label(values.MarginPadding15, amount).Layout(common.gtx)
			}),
			layout.Rigid(func() {
				txt := common.theme.Label(values.MarginPadding15, hash)
				txt.Color = common.theme.Color.Primary
				txt.Layout(common.gtx)
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
	if page.viewTxnOnDcrdataW.Clicked(common.gtx) {
		page.viewTxnOnBrowser(&common)
	}
	if page.backButton.Clicked(common.gtx) {
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
