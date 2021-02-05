package ui

import (
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
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageTransactions = "Transactions"

type transactionWdg struct {
	statusIcon   *widget.Image
	direction    *widget.Image
	time, status decredmaterial.Label
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
	line                                        *decredmaterial.Line

	orderDropDown  *decredmaterial.DropDown
	txTypeDropDown *decredmaterial.DropDown
	walletDropDown *decredmaterial.DropDown
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
		line:                   common.theme.Line(),
	}
	pg.line.Color = common.theme.Color.Background
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

	walletDropDownItems := []decredmaterial.DropDownItem{}
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
		walTxs := (*pg.walletTransactions).Txs[walletID]
		pg.updateTotransactionDetailsButtons(&walTxs)

		directionFilter := pg.txTypeDropDown.SelectedIndex()

		return layout.Stack{Alignment: layout.N}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Inset{
					Top: values.MarginPadding60,
				}.Layout(gtx, func(gtx C) D {
					return common.theme.Card().Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
							if len(walTxs) == 0 {
								txt := common.theme.Body1("No transactions")
								txt.Alignment = text.Middle
								return txt.Layout(gtx)
							}

							return pg.txsList.Layout(gtx, len(walTxs), func(gtx C, index int) D {
								if directionFilter != 0 && walTxs[index].Txn.Direction != int32(directionFilter-1) {
									return layout.Dimensions{}
								}

								click := pg.toTxnDetails[index]
								pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
								click.Add(gtx.Ops)
								pg.goToTxnDetails(gtx, &common, &walTxs[index], click)
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return pg.txnRowInfo(gtx, &common, walTxs[index])
									}),
									layout.Rigid(func(gtx C) D {
										pg.line.Width = gtx.Constraints.Max.X
										if index < len(walTxs)-1 {
											return layout.Inset{
												Top:    values.MarginPadding10,
												Bottom: values.MarginPadding10,
											}.Layout(gtx, func(gtx C) D {
												return pg.line.Layout(gtx)
											})
										}

										return layout.Dimensions{}
									}),
								)
							})
						})
					})
				})
			}),
			layout.Stacked(func(gtx C) D {
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
			}),
		)
	}
	return common.Layout(gtx, container)
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

func (pg *transactionsPage) txnRowInfo(gtx layout.Context, common *pageCommon, transaction wallet.Transaction) layout.Dimensions {
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, &transaction, &txnWidgets)

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					sz := gtx.Constraints.Max.X

					icon := txnWidgets.direction
					icon.Scale = float32(sz) / float32(gtx.Px(unit.Dp(float32(sz))))
					return icon.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding15, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return layoutBalance(gtx, transaction.Balance, *common)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return txnWidgets.status.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.TextSize12, Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						sz := gtx.Constraints.Max.X

						icon := txnWidgets.statusIcon
						icon.Scale = float32(sz) / float32(gtx.Px(unit.Dp(float32(sz))))
						return icon.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *transactionsPage) Handle(common pageCommon) {
	sortSelection := strconv.Itoa(pg.orderDropDown.SelectedIndex())

	if pg.filterSorter != sortSelection {
		pg.filterSorter = sortSelection
		pg.sortTransactions(&common)
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
			c.ChangePage(PageTransactionDetails)
		}
	}
}

func initTxnWidgets(common *pageCommon, transaction *wallet.Transaction, txWidgets *transactionWdg) {
	t := time.Unix(transaction.Txn.Timestamp, 0).UTC()
	txWidgets.time = common.theme.Body1(t.Format(time.UnixDate))
	txWidgets.status = common.theme.Body1("")

	if transaction.Status == "confirmed" {
		txWidgets.status.Text = formatDateOrTime(transaction.Txn.Timestamp)
		txWidgets.statusIcon = common.icons.confirmIcon
	} else {
		txWidgets.status.Text = transaction.Status
		txWidgets.status.Color = common.theme.Color.Gray
		txWidgets.statusIcon = common.icons.pendingIcon
	}

	if transaction.Txn.Direction == dcrlibwallet.TxDirectionSent {
		txWidgets.direction = common.icons.sendIcon
	} else {
		txWidgets.direction = common.icons.receiveIcon
	}
}
