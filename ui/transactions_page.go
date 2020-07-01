package ui

import (
	"fmt"
	"image"
	"sort"
	"strconv"
	"time"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

const PageTransactions = "transactions"

type transactionWdg struct {
	status, direction *decredmaterial.Icon
	amount, time      decredmaterial.Label
}

type transactionsPage struct {
	container                                   layout.Flex
	txsList                                     layout.List
	walletTransactions                          **wallet.Transactions
	walletTransaction                           **wallet.Transaction
	filterSorter                                string
	filterSortW, filterDirectionW               *widget.Enum
	filterDirection, filterSort                 []decredmaterial.RadioButton
	defaultFilterSorter, defaultFilterDirection string
	toTxnDetails                                []*gesture.Click

	rowDirectionWidth,
	rowDateWidth,
	rowStatusWidth,
	rowAmountWidth,
	rowFeeWidth float32
}

func (win *Window) TransactionsPage(common pageCommon) layout.Widget {
	page := transactionsPage{
		container:              layout.Flex{Axis: layout.Vertical},
		txsList:                layout.List{Axis: layout.Vertical},
		walletTransactions:     &win.walletTransactions,
		walletTransaction:      &win.walletTransaction,
		filterDirectionW:       new(widget.Enum),
		filterSortW:            new(widget.Enum),
		defaultFilterSorter:    "0",
		defaultFilterDirection: "0",
		rowDirectionWidth:      .04,
		rowDateWidth:           .2,
		rowStatusWidth:         .2,
		rowAmountWidth:         .3,
		rowFeeWidth:            .26,
	}

	page.filterSorter = page.defaultFilterSorter
	page.filterDirectionW.SetValue(page.defaultFilterDirection)
	page.filterSortW.SetValue(page.defaultFilterSorter)

	txFilterDirection := []string{"All", "Sent", "Received", "Transfer"}
	txFilterSorts := []string{"Newest", "Oldest"}

	for i := 0; i < len(txFilterDirection); i++ {
		page.filterDirection = append(
			page.filterDirection,
			common.theme.RadioButton(fmt.Sprint(i), txFilterDirection[i]))
		page.filterDirection[i].Size = values.MarginPadding20
	}

	for i := 0; i < len(txFilterSorts); i++ {
		page.filterSort = append(page.filterSort,
			common.theme.RadioButton(fmt.Sprint(i), txFilterSorts[i]))
		page.filterSort[i].Size = values.MarginPadding20
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
			layout.Rigid(page.txsFilters(&common)),
			layout.Flexed(1, func() {
				layout.Inset{
					Left:  values.MarginPadding15,
					Right: values.MarginPadding15}.Layout(gtx, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							layout.Inset{
								Top:    values.MarginPadding15,
								Bottom: values.MarginPadding15}.Layout(gtx, func() {
								page.txnRowHeader(&common)
							})
						}),
						layout.Flexed(1, func() {
							walletID := common.info.Wallets[*common.selectedWallet].ID
							walTxs := (*page.walletTransactions).Txs[walletID]
							page.updateTotransactionDetailsButtons(&walTxs)

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

								click := page.toTxnDetails[index]
								pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
								click.Add(gtx.Ops)
								page.goToTxnDetails(&common, &walTxs[index], click)
							})
						}),
					)
				})
			}),
		)
	}
	common.LayoutWithWallets(gtx, container)
}

func (page *transactionsPage) txsFilters(common *pageCommon) layout.Widget {
	gtx := common.gtx
	return func() {
		layout.Inset{
			Top:    values.MarginPadding15,
			Left:   values.MarginPadding15,
			Bottom: values.MarginPadding15}.Layout(gtx, func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					(&layout.List{Axis: layout.Horizontal}).
						Layout(gtx, len(page.filterSort), func(index int) {
							layout.Inset{Right: values.MarginPadding15}.Layout(gtx, func() {
								page.filterSort[index].Layout(gtx, page.filterSortW)
							})
						})
				}),
				layout.Rigid(func() {
					layout.Inset{
						Left:  values.MarginPadding35,
						Right: values.MarginPadding35,
						Top:   values.MarginPadding5}.Layout(gtx, func() {
						rect := f32.Rectangle{Max: f32.Point{X: 1, Y: 35}}
						op.TransformOp{}.Offset(f32.Point{X: 0, Y: 0}).Add(gtx.Ops)
						paint.ColorOp{Color: common.theme.Color.Hint}.Add(gtx.Ops)
						paint.PaintOp{Rect: rect}.Add(gtx.Ops)
					})
				}),
				layout.Rigid(func() {
					(&layout.List{Axis: layout.Horizontal}).
						Layout(gtx, len(page.filterDirection), func(index int) {
							layout.Inset{Right: values.MarginPadding15}.Layout(gtx, func() {
								page.filterDirection[index].Layout(gtx, page.filterDirectionW)
							})
						})
				}),
			)
		})
	}
}

func (page *transactionsPage) txnRowHeader(common *pageCommon) {
	gtx := common.gtx
	txt := common.theme.Label(values.MarginPadding15, "#")
	txt.Color = common.theme.Color.Hint

	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(page.rowDirectionWidth, func() {
			txt.Layout(gtx)
		}),
		layout.Flexed(page.rowDateWidth, func() {
			txt.Alignment = text.Middle
			txt.Text = "Date (UTC)"
			txt.Layout(gtx)
		}),
		layout.Flexed(page.rowStatusWidth, func() {
			txt.Text = "Status"
			txt.Layout(gtx)
		}),
		layout.Flexed(page.rowAmountWidth, func() {
			txt.Text = "Amount"
			txt.Layout(gtx)
		}),
		layout.Flexed(page.rowFeeWidth, func() {
			txt.Text = "Fee"
			txt.Layout(gtx)
		}),
	)
}

func (page *transactionsPage) txnRowInfo(common *pageCommon, transaction wallet.Transaction) {
	gtx := common.gtx
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, &transaction, &txnWidgets)

	layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(page.rowDirectionWidth, func() {
				layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func() {
					txnWidgets.direction.Layout(gtx, values.MarginPadding15)
				})
			}),
			layout.Flexed(page.rowDateWidth, func() {
				txnWidgets.time.Alignment = text.Middle
				txnWidgets.time.Layout(gtx)
			}),
			layout.Flexed(page.rowStatusWidth, func() {
				txt := common.theme.Body1(transaction.Status)
				txt.Alignment = text.Middle
				txt.Layout(gtx)
			}),
			layout.Flexed(page.rowAmountWidth, func() {
				txnWidgets.amount.Alignment = text.End
				txnWidgets.amount.Layout(gtx)
			}),
			layout.Flexed(page.rowFeeWidth, func() {
				txt := common.theme.Body1(dcrutil.Amount(transaction.Txn.Fee).String())
				txt.Alignment = text.End
				txt.Layout(gtx)
			}),
		)
	})
}

func (page *transactionsPage) Handle(common pageCommon) {
	if page.filterSorter != page.filterSortW.Value(common.gtx) {
		page.filterSorter = page.filterSortW.Value(common.gtx)
		page.sortTransactions(&common)
	}
}

func initTxnWidgets(common *pageCommon, transaction *wallet.Transaction, txWidgets *transactionWdg) {
	txWidgets.amount = common.theme.Label(values.MarginPadding15, transaction.Balance)
	txWidgets.time = common.theme.Body1(transaction.DateTime)

	if transaction.Status == "confirmed" {
		txWidgets.status = common.icons.actionCheckCircle
		txWidgets.status.Color = common.theme.Color.Success
	} else {
		txWidgets.status = common.icons.toggleRadioButtonUnchecked
	}

	if transaction.Txn.Direction == dcrlibwallet.TxDirectionSent {
		txWidgets.direction = common.icons.contentRemove
		txWidgets.direction.Color = common.theme.Color.Danger
	} else {
		txWidgets.direction = common.icons.contentAdd
		txWidgets.direction.Color = common.theme.Color.Success
	}
}

func (page *transactionsPage) sortTransactions(common *pageCommon) {
	newestFirst := page.filterSorter == page.defaultFilterSorter

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

func (page *transactionsPage) updateTotransactionDetailsButtons(walTxs *[]wallet.Transaction) {
	if len(*walTxs) != len(page.toTxnDetails) {
		page.toTxnDetails = make([]*gesture.Click, len(*walTxs))
		for index := range *walTxs {
			page.toTxnDetails[index] = &gesture.Click{}
		}
	}
}

func (page *transactionsPage) goToTxnDetails(c *pageCommon, txn *wallet.Transaction, click *gesture.Click) {
	for _, e := range click.Events(c.gtx) {
		if e.Type == gesture.TypeClick {
			*page.walletTransaction = txn
			*c.page = PageTransactionDetails
		}
	}
}
