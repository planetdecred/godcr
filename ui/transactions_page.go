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
	pg := transactionsPage{
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

	pg.filterSorter = pg.defaultFilterSorter
	pg.filterDirectionW.Value = pg.defaultFilterDirection
	pg.filterSortW.Value = pg.defaultFilterSorter

	txFilterDirection := []string{"All", "Sent", "Received", "Transfer"}
	txFilterSorts := []string{"Newest", "Oldest"}

	for i := 0; i < len(txFilterDirection); i++ {
		pg.filterDirection = append(
			pg.filterDirection,
			common.theme.RadioButton(pg.filterDirectionW, fmt.Sprint(i), txFilterDirection[i]))
		pg.filterDirection[i].Size = values.MarginPadding20
	}

	for i := 0; i < len(txFilterSorts); i++ {
		pg.filterSort = append(pg.filterSort,
			common.theme.RadioButton(pg.filterSortW, fmt.Sprint(i), txFilterSorts[i]))
		pg.filterSort[i].Size = values.MarginPadding20
	}

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *transactionsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	container := func(gtx C) D {
		return pg.container.Layout(gtx,
			layout.Rigid(pg.txsFilters(&common)),
			layout.Flexed(1, func(gtx C) D {
				return layout.Inset{
					Left:  values.MarginPadding15,
					Right: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Top:    values.MarginPadding15,
								Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
								return pg.txnRowHeader(gtx, &common)
							})
						}),
						layout.Flexed(1, func(gtx C) D {
							walletID := common.info.Wallets[*common.selectedWallet].ID
							walTxs := (*pg.walletTransactions).Txs[walletID]
							pg.updateTotransactionDetailsButtons(&walTxs)

							if len(walTxs) == 0 {
								txt := common.theme.Body1("No transactions")
								txt.Alignment = text.Middle
								return txt.Layout(gtx)
							}

							directionFilter, _ := strconv.Atoi(pg.filterDirectionW.Value)
							return pg.txsList.Layout(gtx, len(walTxs), func(gtx C, index int) D {
								if directionFilter != 0 && walTxs[index].Txn.Direction != int32(directionFilter-1) {
									return layout.Dimensions{}
								}

								click := pg.toTxnDetails[index]
								pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
								click.Add(gtx.Ops)
								pg.goToTxnDetails(gtx, &common, &walTxs[index], click)
								return pg.txnRowInfo(gtx, &common, walTxs[index])
							})
						}),
					)
				})
			}),
		)
	}
	return common.LayoutWithWallets(gtx, container)
}

func (pg *transactionsPage) txsFilters(common *pageCommon) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{
			Top:    values.MarginPadding15,
			Left:   values.MarginPadding15,
			Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return (&layout.List{Axis: layout.Horizontal}).
						Layout(gtx, len(pg.filterSort), func(gtx C, index int) D {
							return layout.Inset{Right: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
								return pg.filterSort[index].Layout(gtx)
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
						Layout(gtx, len(pg.filterDirection), func(gtx C, index int) D {
							return layout.Inset{Right: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
								return pg.filterDirection[index].Layout(gtx)
							})
						})
				}),
			)
		})
	}
}

func (pg *transactionsPage) txnRowHeader(gtx layout.Context, common *pageCommon) layout.Dimensions {
	txt := common.theme.Label(values.MarginPadding15, "#")
	txt.Color = common.theme.Color.Hint

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding60)
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt.Alignment = text.Middle
			txt.Text = "Date (UTC)"
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt.Text = "Status"
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt.Text = "Amount"
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt.Text = "Fee"
			return txt.Layout(gtx)
		}),
	)
}

func (pg *transactionsPage) txnRowInfo(gtx layout.Context, common *pageCommon, transaction wallet.Transaction) layout.Dimensions {
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, &transaction, &txnWidgets)

	return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Px(values.MarginPadding60)
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return txnWidgets.direction.Layout(gtx, values.MarginPadding15)
				})
			}),
			layout.Rigid(func(gtx C) D {
				txnWidgets.time.Alignment = text.Middle
				return txnWidgets.time.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := common.theme.Body1(transaction.Status)
				txt.Alignment = text.Middle
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txnWidgets.amount.Alignment = text.End
				return txnWidgets.amount.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := common.theme.Body1(dcrutil.Amount(transaction.Txn.Fee).String())
				txt.Alignment = text.End
				return txt.Layout(gtx)
			}),
		)
	})
}

func (pg *transactionsPage) Handle(common pageCommon) {
	if pg.filterSorter != pg.filterSortW.Value {
		pg.filterSorter = pg.filterSortW.Value
		pg.sortTransactions(&common)
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

func (pg *transactionsPage) sortTransactions(common *pageCommon) {
	newestFirst := pg.filterSorter == pg.defaultFilterSorter

	for _, wal := range common.info.Wallets {
		transactions := (*pg.walletTransactions).Txs[wal.ID]
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

func (pg *transactionsPage) updateTotransactionDetailsButtons(walTxs *[]wallet.Transaction) {
	if len(*walTxs) != len(pg.toTxnDetails) {
		pg.toTxnDetails = make([]*gesture.Click, len(*walTxs))
		for index := range *walTxs {
			pg.toTxnDetails[index] = &gesture.Click{}
		}
	}
}

func (pg *transactionsPage) goToTxnDetails(gtx layout.Context, c *pageCommon, txn *wallet.Transaction, click *gesture.Click) {
	for _, e := range click.Events(gtx) {
		if e.Type == gesture.TypeClick {
			*pg.walletTransaction = txn
			*c.page = PageTransactionDetails
		}
	}
}
