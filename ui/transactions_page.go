package ui

import (
	"image"
	"sort"
	"time"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageTransactions = "Transactions"

type transactionWdg struct {
	statusIcon           *widget.Image
	direction            *widget.Image
	time, status, wallet decredmaterial.Label
}

type transactionsPage struct {
<<<<<<< HEAD
	container                   layout.Flex
	txsList                     layout.List
	walletTransactions          **wallet.Transactions
	walletTransaction           **wallet.Transaction
	filterSorter                int
	filterDirection, filterSort []decredmaterial.RadioButton
	toTxnDetails                []*gesture.Click
	separator                   decredmaterial.Line
=======
	container                     layout.Flex
	txsList                       layout.List
	walletTransactions            **wallet.Transactions
	walletTransaction             **wallet.Transaction
	filterSorter                  int
	filterSortW, filterDirectionW *widget.Enum
	filterDirection, filterSort   []decredmaterial.RadioButton
	toTxnDetails                  []*gesture.Click
	separator                     decredmaterial.Line
	theme 	*decredmaterial.Theme
>>>>>>> Implement dark mode

	orderDropDown  *decredmaterial.DropDown
	txTypeDropDown *decredmaterial.DropDown
	walletDropDown *decredmaterial.DropDown
}

func (win *Window) TransactionsPage(common pageCommon) layout.Widget {
	pg := transactionsPage{
		container:          layout.Flex{Axis: layout.Vertical},
		txsList:            layout.List{Axis: layout.Vertical},
		walletTransactions: &win.walletTransactions,
		walletTransaction:  &win.walletTransaction,
		separator:          common.theme.Separator(),
		theme: common.theme,
	}

	pg.orderDropDown = common.theme.DropDown([]decredmaterial.DropDownItem{{Text: "Newest"}, {Text: "Oldest"}}, 1)
	pg.txTypeDropDown = common.theme.DropDown([]decredmaterial.DropDownItem{
		{
			Text: "All",
		},
		{
			Text: "Sent",
		},
		{
			Text: "Received",
		},
		{
			Text: "Yourself",
		},
		{
			Text: "Staking",
		},
	}, 1)

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *transactionsPage) setWallets(common pageCommon) {
	if len(common.info.Wallets) == 0 || pg.walletDropDown != nil {
		return
	}

	var walletDropDownItems []decredmaterial.DropDownItem
	for i := range common.info.Wallets {
		item := decredmaterial.DropDownItem{
			Text: common.info.Wallets[i].Name,
			Icon: common.icons.walletIcon,
		}
		walletDropDownItems = append(walletDropDownItems, item)
	}
	pg.walletDropDown = common.theme.DropDown(walletDropDownItems, 2)
}

func (pg *transactionsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.setWallets(common)
	container := func(gtx C) D {
		walletID := common.info.Wallets[pg.walletDropDown.SelectedIndex()].ID
		wallTxs := (*pg.walletTransactions).Txs[walletID]
		if pg.txTypeDropDown.SelectedIndex()-1 != -1 {
			wallTxs = filterTransactions(wallTxs, func(i int) bool {
				return i == pg.txTypeDropDown.SelectedIndex()-1
			})
		}

		return layout.Stack{Alignment: layout.N}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Inset{
					Top: values.MarginPadding60,
				}.Layout(gtx, func(gtx C) D {
					return common.theme.Card().Layout(gtx, func(gtx C) D {
						padding := values.MarginPadding16
						return Container{layout.Inset{Bottom: padding, Left: padding}}.Layout(gtx,
							func(gtx C) D {
								// return "No transactions yet" text if there are no transactions
								if len(wallTxs) == 0 {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									txt := common.theme.Body1("No transactions yet")
									txt.Color = common.theme.Color.Gray2
									return txt.Layout(gtx)
								}

								// update transaction row click gesture when the length of the click gesture slice and
								// transactions list are different.
								if len(wallTxs) != len(pg.toTxnDetails) {
									pg.toTxnDetails = createClickGestures(len(wallTxs))
								}

								return pg.txsList.Layout(gtx, len(wallTxs), func(gtx C, index int) D {
									click := pg.toTxnDetails[index]
									pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
									click.Add(gtx.Ops)
									pg.goToTxnDetails(click.Events(gtx), &common, &wallTxs[index])
									var row = TransactionRow{
										transaction: wallTxs[index],
										index:       index,
										showBadge:   false,
									}
									return transactionRow(gtx, common, row)
								})
							})
					})
				})
			}),
			layout.Stacked(func(gtx C) D {
				return pg.dropDowns(gtx)
			}),
		)
	}
	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, container)
	})
}

func filterTransactions(transactions []wallet.Transaction, f func(int) bool) []wallet.Transaction {
	t := make([]wallet.Transaction, 0)
	for _, v := range transactions {
		if f(int(v.Txn.Direction)) {
			t = append(t, v)
		}
	}
	return t
}

func (pg *transactionsPage) dropDowns(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Bottom: values.MarginPadding10,
	}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.walletDropDown.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: values.MarginPadding5,
						}.Layout(gtx, func(gtx C) D {
							return pg.txTypeDropDown.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: values.MarginPadding5,
						}.Layout(gtx, func(gtx C) D {
							return pg.orderDropDown.Layout(gtx)
						})
					}),
				)
			}),
		)
	})
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
						paint.Fill(gtx.Ops, common.theme.Color.Hint)
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

func (pg *transactionsPage) Handle(common pageCommon) {
	sortSelection := pg.orderDropDown.SelectedIndex()

	if pg.filterSorter != sortSelection {
		pg.filterSorter = sortSelection
		pg.sortTransactions(&common)
	}
}

func (pg *transactionsPage) sortTransactions(common *pageCommon) {
	newestFirst := pg.filterSorter == 0

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

func (pg *transactionsPage) goToTxnDetails(events []gesture.ClickEvent, common *pageCommon, txn *wallet.Transaction) {
	for _, e := range events {
		if e.Type == gesture.TypeClick {
			*pg.walletTransaction = txn

			common.setReturnPage(PageTransactions)
			common.changePage(PageTransactionDetails)
		}
	}
}

func initTxnWidgets(common pageCommon, transaction wallet.Transaction) transactionWdg {
	var txn transactionWdg
	t := time.Unix(transaction.Txn.Timestamp, 0).UTC()
	txn.time = common.theme.Body1(t.Format(time.UnixDate))
	txn.status = common.theme.Body1("")
	txn.wallet = common.theme.Body2(transaction.WalletName)

	if transaction.Status == "confirmed" {
		txn.status.Text = formatDateOrTime(transaction.Txn.Timestamp)
		txn.statusIcon = common.icons.confirmIcon
	} else {
		txn.status.Text = transaction.Status
		txn.status.Color = common.theme.Color.Gray
		txn.statusIcon = common.icons.pendingIcon
	}

	if transaction.Txn.Direction == dcrlibwallet.TxDirectionSent {
		txn.direction = common.icons.sendIcon
	} else {
		txn.direction = common.icons.receiveIcon
	}

	return txn
}
