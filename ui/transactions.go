package ui

import (
	"fmt"
	"sort"
	"strconv"
	"time"

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
	PageTransactions                            = "txs"
	defaultFilterSorter, defaultFilterDirection = "0", "0"

	rowDirectionWidth = .04
	rowDateWidth      = .2
	rowStatusWidth    = .2
	rowAmountWidth    = .3
	rowFeeWidth       = .26

	txsRowLabelSize    = 16
	txsPageInsetTop    = 15
	txsPageInsetLeft   = 15
	txsPageInsetRight  = 15
	txsPageInsetBottom = 15
)

type transactionWdg struct {
	status, direction *decredmaterial.Icon
	amount, time      decredmaterial.Label
}

type transactionsPage struct {
	container layout.Flex
	txsList   layout.List

	toSend, toReceive               decredmaterial.IconButton
	toSendW, toReceiveW, toFiltersW widget.Button
	walletTransactions              **wallet.Transactions
	filterSorter                    string
	filterSortW, filterDirectionW   *widget.Enum
	filterDirection, filterSort     []decredmaterial.RadioButton
}

func (win *Window) TransactionsPage(common pageCommon) layout.Widget {
	page := transactionsPage{
		container:          layout.Flex{Axis: layout.Vertical},
		txsList:            layout.List{Axis: layout.Vertical},
		toSend:             common.theme.PlainIconButton(common.icons.contentSend),
		toReceive:          common.theme.PlainIconButton(common.icons.contentAddBox),
		walletTransactions: &win.walletTransactions,
		filterSorter:       defaultFilterSorter,
		filterDirectionW:   new(widget.Enum),
		filterSortW:        new(widget.Enum),
	}
	page.toSend.Size, page.toReceive.Size = unit.Dp(40), unit.Dp(40)
	page.filterDirectionW.SetValue(defaultFilterDirection)
	page.filterSortW.SetValue(defaultFilterSorter)

	txFilterDirection := []string{"All", "Sent", "Received", "Transfer"}
	txFilterSorts := []string{"Newest", "Oldest"}

	for i := 0; i < len(txFilterDirection); i++ {
		page.filterDirection = append(
			page.filterDirection,
			common.theme.RadioButton(fmt.Sprint(i), txFilterDirection[i]))
		page.filterDirection[i].Color = common.theme.Color.Success
		page.filterDirection[i].Size = unit.Dp(20)
	}

	for i := 0; i < len(txFilterSorts); i++ {
		page.filterSort = append(page.filterSort,
			common.theme.RadioButton(fmt.Sprint(i), txFilterSorts[i]))
		page.filterSort[i].Size = unit.Dp(20)
	}

	return func() {
		page.Layout(common)
		page.Handle(common)
	}
}

func (page *transactionsPage) Layout(common pageCommon) {
	gtx := common.gtx

	container := func() {
		page.container.Layout(gtx,
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(txsPageInsetTop)}.Layout(gtx, func() {
					layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func() {
							layout.Inset{Left: unit.Dp(txsPageInsetLeft)}.Layout(gtx, func() {
								page.renderFiltererButtons(&common)
							})
						}),
						layout.Rigid(func() {
							page.toSend.Layout(gtx, &page.toSendW)
						}),
						layout.Rigid(func() {
							layout.Inset{Right: unit.Dp(txsPageInsetRight)}.Layout(gtx, func() {
								page.toReceive.Layout(gtx, &page.toReceiveW)
							})
						}),
					)
				})
			}),
			layout.Flexed(1, func() {
				layout.Inset{Left: unit.Dp(txsPageInsetLeft), Right: unit.Dp(txsPageInsetRight)}.Layout(gtx, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							layout.Inset{Top: unit.Dp(txsPageInsetTop), Bottom: unit.Dp(txsPageInsetBottom)}.Layout(gtx, func() {
								page.txnRowHeader(&common)
							})
						}),
						layout.Flexed(1, func() {
							walletID := common.info.Wallets[*common.selectedWallet].ID
							walTxs := (*page.walletTransactions).Txs[walletID]

							if len(walTxs) == 0 {
								txt := common.theme.Body1("No transactions")
								txt.Alignment = text.Middle
								txt.Layout(gtx)
								return
							}
							directionFilter, _ := strconv.Atoi(page.filterDirectionW.Value(gtx))
							page.txsList.Layout(gtx, len(walTxs), func(index int) {
								if directionFilter != 0 && walTxs[index].Txn.Direction != int32(directionFilter-1) {
									return
								}
								page.txnRowInfo(&common, walTxs[index])
							})
						}),
					)
				})
			}),
		)
	}
	common.LayoutWithWallets(gtx, container)
}

func (page *transactionsPage) txnRowHeader(common *pageCommon) {
	gtx := common.gtx
	txt := common.theme.Label(unit.Dp(txsRowLabelSize), "#")
	txt.Color = common.theme.Color.Hint
	txt.Alignment = text.Middle

	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(rowDirectionWidth, func() {
			txt.Layout(gtx)
		}),
		layout.Flexed(rowDateWidth, func() {
			txt.Text = "Date (UTC)"
			txt.Layout(gtx)
		}),
		layout.Flexed(rowStatusWidth, func() {
			txt.Text = "Status"
			txt.Layout(gtx)
		}),
		layout.Flexed(rowAmountWidth, func() {
			txt.Text = "Amount"
			txt.Layout(gtx)
		}),
		layout.Flexed(rowFeeWidth, func() {
			txt.Text = "Fee"
			txt.Layout(gtx)
		}),
	)
}

func (page *transactionsPage) txnRowInfo(common *pageCommon, transaction wallet.Transaction) {
	gtx := common.gtx
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, &transaction, &txnWidgets)

	layout.Inset{Bottom: unit.Dp(txsPageInsetBottom)}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(rowDirectionWidth, func() {
				layout.Inset{Top: unit.Dp(3)}.Layout(gtx, func() {
					txnWidgets.direction.Layout(gtx, unit.Dp(16))
				})
			}),
			layout.Flexed(rowDateWidth, func() {
				txnWidgets.time.Alignment = text.Middle
				txnWidgets.time.Layout(gtx)
			}),
			layout.Flexed(rowStatusWidth, func() {
				txt := common.theme.Body1(transaction.Status)
				txt.Alignment = text.Middle
				txt.Layout(gtx)
			}),
			layout.Flexed(rowAmountWidth, func() {
				txnWidgets.amount.Alignment = text.End
				txnWidgets.amount.Layout(gtx)
			}),
			layout.Flexed(rowFeeWidth, func() {
				txt := common.theme.Body1(dcrutil.Amount(transaction.Txn.Fee).String())
				txt.Alignment = text.End
				txt.Layout(gtx)
			}),
		)
	})
}

func (page *transactionsPage) renderFiltererButtons(common *pageCommon) {
	gtx := common.gtx
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			(&layout.List{Axis: layout.Vertical}).
				Layout(gtx, len(page.filterSort), func(index int) {
					page.filterSort[index].Layout(gtx, page.filterSortW)
				})
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(35)}.Layout(gtx, func() {})
		}),
		layout.Rigid(func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					(&layout.List{Axis: layout.Vertical}).
						Layout(gtx, 2, func(index int) {
							page.filterDirection[index].Layout(gtx, page.filterDirectionW)
						})
				}),
				layout.Rigid(func() {
					(&layout.List{Axis: layout.Vertical}).
						Layout(gtx, 2, func(index int) {
							page.filterDirection[index+2].Layout(gtx, page.filterDirectionW)
						})
				}),
			)
		}),
	)
}

func (page *transactionsPage) Handle(common pageCommon) {
	if page.toReceiveW.Clicked(common.gtx) {
		*common.page = PageReceive
		return
	}

	if page.toSendW.Clicked(common.gtx) {
		*common.page = PageSend
		return
	}

	if page.filterSorter != page.filterSortW.Value(common.gtx) {
		page.filterSorter = page.filterSortW.Value(common.gtx)
		page.sortTransactions(&common)
	}
}

func initTxnWidgets(common *pageCommon,
	transaction *wallet.Transaction, txWidgets *transactionWdg) {
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
		txWidgets.direction = common.icons.contentAdd
		txWidgets.direction.Color = common.theme.Color.Success
	}
}

func (page *transactionsPage) sortTransactions(common *pageCommon) {
	newestFirst := page.filterSorter == defaultFilterSorter

	for _, wal := range common.info.Wallets {
		transactions := (*page.walletTransactions).Txs[wal.ID]
		sort.SliceStable(transactions, func(i, j int) bool {
			backTime := time.Unix(transactions[j].Txn.Timestamp, 0)
			frontTime := time.Unix(transactions[i].Txn.Timestamp, 0)
			if newestFirst {
				return backTime.Before(frontTime)
			}
			return frontTime.Before(backTime)
		})
	}
}
