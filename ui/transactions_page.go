package ui

import (
	"fmt"
	"image"
	"sort"
	"strconv"
	"time"

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
	"github.com/raedahgroup/godcr/ui/values"
	"github.com/raedahgroup/godcr/wallet"
)

const PageTransactions = "transactions"

type transactionWdg struct {
	status, direction *widget.Icon
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
	page.filterDirectionW.Value = page.defaultFilterDirection
	page.filterSortW.Value = page.defaultFilterSorter

	txFilterDirection := []string{"All", "Sent", "Received", "Transfer"}
	txFilterSorts := []string{"Newest", "Oldest"}

	for i := 0; i < len(txFilterDirection); i++ {
		page.filterDirection = append(
			page.filterDirection,
			common.theme.RadioButton(page.filterDirectionW, fmt.Sprint(i), txFilterDirection[i]))
		page.filterDirection[i].Size = values.MarginPadding20
	}

	for i := 0; i < len(txFilterSorts); i++ {
		page.filterSort = append(page.filterSort,
			common.theme.RadioButton(page.filterSortW, fmt.Sprint(i), txFilterSorts[i]))
		page.filterSort[i].Size = values.MarginPadding20
	}

	return func(gtx C) D {
		page.Handle(common)
		return page.Layout(common)
	}
}

func (page *transactionsPage) Layout(common pageCommon) layout.Dimensions {
	gtx := common.gtx
	container := func(gtx C) D {
		return page.container.Layout(gtx,
			layout.Rigid(page.txsFilters(&common)),
			layout.Flexed(1, func(gtx C) D {
				return layout.Inset{
					Left:  values.MarginPadding15,
					Right: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Top:    values.MarginPadding15,
								Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
								return page.txnRowHeader(&common)
							})
						}),
						layout.Flexed(1, func(gtx C) D {
							walletID := common.info.Wallets[*common.selectedWallet].ID
							walTxs := (*page.walletTransactions).Txs[walletID]
							page.updateTotransactionDetailsButtons(&walTxs)

							if len(walTxs) == 0 {
								txt := common.theme.Body1("No transactions")
								txt.Alignment = text.Middle
								return txt.Layout(gtx)
							}

							directionFilter, _ := strconv.Atoi(page.filterDirectionW.Value)
							return page.txsList.Layout(gtx, len(walTxs), func(gtx C, index int) D {
								if directionFilter != 0 && walTxs[index].Txn.Direction != int32(directionFilter-1) {
									return layout.Dimensions{}
								}

								click := page.toTxnDetails[index]
								pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
								click.Add(gtx.Ops)
								page.goToTxnDetails(&common, &walTxs[index], click)
								return page.txnRowInfo(&common, walTxs[index])
							})
						}),
					)
				})
			}),
		)
	}
	return common.LayoutWithWallets(gtx, container)
}

func (page *transactionsPage) txsFilters(common *pageCommon) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{
			Top:    values.MarginPadding15,
			Left:   values.MarginPadding15,
			Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return (&layout.List{Axis: layout.Horizontal}).
						Layout(gtx, len(page.filterSort), func(gtx C, index int) D {
							return layout.Inset{Right: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
								return page.filterSort[index].Layout(gtx)
							})
						})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left:  values.MarginPadding35,
						Right: values.MarginPadding35,
						Top:   values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						dims := image.Point{X: 1, Y: 35}
						rect := f32.Rectangle{Max: layout.FPt(dims)}
						rect.Size()
						op.TransformOp{}.Add(gtx.Ops)
						paint.ColorOp{Color: common.theme.Color.Hint}.Add(gtx.Ops)
						paint.PaintOp{Rect: rect}.Add(gtx.Ops)
						return layout.Dimensions{Size: dims}
					})
				}),
				layout.Rigid(func(gtx C) D {
					return (&layout.List{Axis: layout.Horizontal}).
						Layout(gtx, len(page.filterDirection), func(gtx C, index int) D {
							return layout.Inset{Right: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
								return page.filterDirection[index].Layout(gtx)
							})
						})
				}),
			)
		})
	}
}

func (page *transactionsPage) txnRowHeader(common *pageCommon) layout.Dimensions {
	gtx := common.gtx
	txt := common.theme.Label(values.MarginPadding15, "#")
	txt.Color = common.theme.Color.Hint

	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding60)
			txt.Layout(gtx)
		}),
		layout.Rigid(func() {
			txt.Alignment = text.Middle
			txt.Text = "Date (UTC)"
			gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding120)
			txt.Layout(gtx)
		}),
		layout.Rigid(func() {
			txt.Text = "Status"
			gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding120)
			txt.Layout(gtx)
		}),
		layout.Rigid(func() {
			txt.Text = "Amount"
			gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding150)
			txt.Layout(gtx)
		}),
		layout.Rigid(func() {
			txt.Text = "Fee"
			gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding150)
			txt.Layout(gtx)
		}),
	)
}

func (page *transactionsPage) txnRowInfo(common *pageCommon, transaction wallet.Transaction) layout.Dimensions {
	gtx := common.gtx
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, &transaction, &txnWidgets)

	layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func() {
				gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding60)
				layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func() {
					txnWidgets.direction.Layout(gtx, values.MarginPadding15)
				})
			}),
			layout.Rigid(func() {
				txnWidgets.time.Alignment = text.Middle
				gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding120)
				txnWidgets.time.Layout(gtx)
			}),
			layout.Rigid(func() {
				txt := common.theme.Body1(transaction.Status)
				txt.Alignment = text.Middle
				gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding120)
				txt.Layout(gtx)
			}),
			layout.Rigid(func() {
				txnWidgets.amount.Alignment = text.End
				gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding120)
				txnWidgets.amount.Layout(gtx)
			}),
			layout.Rigid(func() {
				txt := common.theme.Body1(dcrutil.Amount(transaction.Txn.Fee).String())
				txt.Alignment = text.End
				gtx.Constraints.Width.Min = gtx.Px(values.MarginPadding150)
				txt.Layout(gtx)
			}),
		)
	})
}

func (page *transactionsPage) Handle(common pageCommon) {
	if page.filterSorter != page.filterSortW.Value {
		page.filterSorter = page.filterSortW.Value
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
