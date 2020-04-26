package ui

import (
	"fmt"
	"os/exec"
	"runtime"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const (
	PageTransactionDetails = "transaction"
	pageContentInset       = 15
	rowGroupInset          = 20
)

type transactionWdg struct {
	status, direction *decredmaterial.Icon
	amount, time      decredmaterial.Label
}

type transactionPage struct {
	transactionPageContainer                       layout.List
	transactionInputsContainer                     layout.List
	transactionOutputsContainer                    layout.List
	hideTransaction                                widget.Button
	hideTransactionW, expandMore, expandLess       decredmaterial.IconButton
	details                                        **wallet.Transaction
	isTxnInputsShow, isTxnOutputsShow              bool
	expandInputs, expandOutputs, viewTxnOnDcrdataW widget.Button
	viewTxnOnDcrdata                               decredmaterial.Button
}

func TransactionPage(common pageCommon, transaction **wallet.Transaction) layout.Widget {
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
		hideTransactionW: common.theme.PlainIconButton(common.icons.contentClear),
		details:          transaction,
		expandMore: decredmaterial.IconButton{
			Icon:    mustIcon(decredmaterial.NewIcon(icons.NavigationExpandMore)),
			Size:    unit.Dp(25),
			Color:   common.theme.Color.Text,
			Padding: unit.Dp(0),
		},
		expandLess: decredmaterial.IconButton{
			Icon:    mustIcon(decredmaterial.NewIcon(icons.NavigationExpandLess)),
			Size:    unit.Dp(25),
			Color:   common.theme.Color.Text,
			Padding: unit.Dp(0),
		},
		viewTxnOnDcrdata: common.theme.Button("View on dcrdata"),
	}

	return func() {
		page.layout(common)
		page.handler(common)
	}
}

func (page *transactionPage) layout(common pageCommon) {
	gtx := common.gtx
	widgets := []func(){
		func() {
			layout.Inset{Top: unit.Dp(rowGroupInset)}.Layout(gtx, func() {
				page.txnBalanceAndStatus(&common)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(rowGroupInset)}.Layout(gtx, func() {
				page.txnTypeAndID(&common)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(rowGroupInset)}.Layout(gtx, func() {
				page.txnInputs(&common)
			})
		},
		func() {
			layout.Inset{Top: unit.Dp(rowGroupInset)}.Layout(gtx, func() {
				page.txnOutputs(&common)
			})
		},
	}

	common.LayoutWithWallets(gtx, func() {
		// common.theme.Surface(gtx, func() {
		layout.UniformInset(unit.Dp(pageContentInset)).Layout(gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					page.txnDirection(&common)
				}),
				layout.Flexed(1, func() {
					if *page.details == nil {
						return
					}
					page.transactionPageContainer.Layout(gtx, len(widgets), func(i int) {
						layout.Inset{}.Layout(gtx, widgets[i])
					})
				}),
				layout.Rigid(func() {
					if *page.details == nil {
						return
					}
					layout.Center.Layout(gtx, func() {
						layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func() {
							page.viewTxnOnDcrdata.Layout(gtx, &page.viewTxnOnDcrdataW)
						})
					})
				}),
			)
		})
		// })
	})
}

func (page *transactionPage) txnDirection(common *pageCommon) {
	txt := common.theme.H4("")
	if *page.details != nil {
		txt.Text = dcrlibwallet.TransactionDirectionName((*page.details).Txn.Direction)
	} else {
		txt.Text = "Not found"
	}

	txt.Alignment = text.Middle
	txt.Layout(common.gtx)
	layout.E.Layout(common.gtx, func() {
		page.hideTransactionW.Layout(common.gtx, &page.hideTransaction)
	})
}

func (page *transactionPage) initTxnWidgets(common *pageCommon, txWidgets *transactionWdg) {
	transaction := *page.details
	txWidgets.amount = common.theme.Label(unit.Dp(16), transaction.Balance)
	txWidgets.time = common.theme.Body1("Pending")

	if transaction.Status == "confirmed" {
		txWidgets.time.Text = dcrlibwallet.ExtractDateOrTime(transaction.Txn.Timestamp)
		txWidgets.status, _ = decredmaterial.NewIcon(icons.ActionCheckCircle)
		txWidgets.status.Color = common.theme.Color.Success
	} else {
		txWidgets.status, _ = decredmaterial.NewIcon(icons.ToggleRadioButtonUnchecked)
	}

	if transaction.Txn.Direction == dcrlibwallet.TxDirectionSent {
		txWidgets.direction, _ = decredmaterial.NewIcon(icons.ContentRemove)
		txWidgets.direction.Color = common.theme.Color.Danger
	} else {
		txWidgets.direction, _ = decredmaterial.NewIcon(icons.ContentAdd)
		txWidgets.direction.Color = common.theme.Color.Success
	}
}

func (page *transactionPage) txnBalanceAndStatus(common *pageCommon) {
	txnWidgets := transactionWdg{}
	page.initTxnWidgets(common, &txnWidgets)

	vertFlex.Layout(common.gtx,
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(-4), Top: unit.Dp(4)}.Layout(common.gtx, func() {
				txnWidgets.direction.Layout(common.gtx, unit.Dp(28))
			})
			layout.Inset{Left: unit.Dp(28)}.Layout(common.gtx, func() {
				txnWidgets.amount.TextSize = unit.Dp(28)
				txnWidgets.amount.Layout(common.gtx)
			})
		}),
		layout.Rigid(func() {
			txnWidgets.time.Layout(common.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(3)}.Layout(common.gtx, func() {
				txnWidgets.status.Layout(common.gtx, unit.Dp(16))
			})
			layout.Inset{Left: unit.Dp(18)}.Layout(common.gtx, func() {
				txt := common.theme.Body1((*page.details).Status)
				txt.Color = txnWidgets.status.Color
				txt.Layout(common.gtx)
			})
		}),
		layout.Rigid(func() {
			txt := common.theme.Body1(fmt.Sprintf("%d confirmations", (*page.details).Confirmations))
			txt.Color = common.theme.Color.Primary
			txt.Layout(common.gtx)
		}),
	)
}

func (page *transactionPage) txnTypeAndID(common *pageCommon) {
	transaction := *page.details
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
			layout.Inset{Bottom: unit.Dp(rowGroupInset)}.Layout(gtx, func() {
				row("Transaction ID", common.theme.Body1(transaction.Txn.Hash))
			})
		}),
		layout.Rigid(func() {
			layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func() {
					row("To", common.theme.H6(transaction.WalletName))
				}),
				layout.Rigid(func() {
					row("Included in block", common.theme.H6(fmt.Sprintf("%d", transaction.Txn.BlockHeight)))
				}),
				layout.Rigid(func() {
					row("Type", common.theme.H6(transaction.Txn.Type))
				}),
			)
		}),
	)
}

func (page *transactionPage) txnInputs(common *pageCommon) {
	transaction := *page.details

	layout.Flex{Axis: layout.Vertical}.Layout(common.gtx,
		layout.Rigid(func() {
			txt := fmt.Sprintf("%d Inputs consumed", len(transaction.Txn.Inputs))
			page.txnIORowHeader(common, txt, &page.expandInputs, page.isTxnInputsShow)
		}),
		layout.Rigid(func() {
			if page.isTxnInputsShow {
				page.transactionInputsContainer.Layout(common.gtx, len(transaction.Txn.Inputs), func(i int) {
					page.txnIORow(common, dcrutil.Amount(transaction.Txn.Inputs[i].Amount).String(),
						transaction.Txn.Inputs[i].PreviousOutpoint)
				})
			}
		}),
	)
}

func (page *transactionPage) txnOutputs(common *pageCommon) {
	transaction := *page.details
	layout.Flex{Axis: layout.Vertical}.Layout(common.gtx,
		layout.Rigid(func() {
			txt := fmt.Sprintf("%d Outputs created", len(transaction.Txn.Outputs))
			page.txnIORowHeader(common, txt, &page.expandOutputs, page.isTxnOutputsShow)
		}),
		layout.Rigid(func() {
			if page.isTxnOutputsShow {
				page.transactionOutputsContainer.Layout(common.gtx, len(transaction.Txn.Outputs), func(i int) {
					page.txnIORow(common, dcrutil.Amount(transaction.Txn.Outputs[i].Amount).String(),
						transaction.Txn.Outputs[i].Address)
				})
			}
		}),
	)
}

func (page *transactionPage) txnIORow(common *pageCommon, amount string, hash string) {
	layout.Inset{Bottom: unit.Dp(5)}.Layout(common.gtx, func() {
		layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(common.gtx,
			layout.Rigid(func() {
				common.theme.Body1(amount).Layout(common.gtx)
			}),
			layout.Rigid(func() {
				txt := common.theme.Body1(hash)
				txt.Color = common.theme.Color.Primary
				txt.Layout(common.gtx)
			}),
		)
	})
}

func (page *transactionPage) txnIORowHeader(common *pageCommon, str string, in *widget.Button, isShow bool) {
	layout.Flex{Spacing: layout.SpaceBetween}.Layout(common.gtx,
		layout.Rigid(func() {
			common.theme.Body1(str).Layout(common.gtx)
		}),
		layout.Rigid(func() {
			if isShow {
				page.expandLess.Layout(common.gtx, in)
				return
			}
			page.expandMore.Layout(common.gtx, in)
		}),
	)
}

func (page *transactionPage) viewTxnOnBrowser(common *pageCommon) {
	var err error
	url := common.wallet.GetBlockExplorerURL((*page.details).Txn.Hash)

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

func (page *transactionPage) handler(common pageCommon) {
	if page.expandInputs.Clicked(common.gtx) {
		page.isTxnInputsShow = !page.isTxnInputsShow
	}
	if page.expandOutputs.Clicked(common.gtx) {
		page.isTxnOutputsShow = !page.isTxnOutputsShow
	}

	if page.viewTxnOnDcrdataW.Clicked(common.gtx) {
		page.viewTxnOnBrowser(&common)
	}
}
