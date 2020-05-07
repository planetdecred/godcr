package ui

import (
	"fmt"
	"image"
	"sort"
	"strconv"
	"time"

	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
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

	toFilters       map[string]decredmaterial.IconButton
	isShowFilters   bool
	filterSort      string
	filterDirection string
	keyEvent        **key.Event
	containerWidth  int
}

type transactionsFiltersSide struct {
	filterSortW, filterDirectionW *widget.Enum
	filterDirection, filterSort   []decredmaterial.RadioButton
}

func (win *Window) TransactionsPage(common pageCommon) layout.Widget {
	page := transactionsPage{
		container:          layout.Flex{Axis: layout.Vertical},
		txsList:            layout.List{Axis: layout.Vertical},
		toSend:             common.theme.PlainIconButton(common.icons.contentSend),
		toReceive:          common.theme.PlainIconButton(common.icons.contentAddBox),
		walletTransactions: &win.walletTransactions,
		filterSort:         defaultFilterSorter,
		filterDirection:    defaultFilterDirection,
		keyEvent:           &win.keyEvt,
	}

	fitlerSide := transactionsFiltersSide{
		filterDirectionW: new(widget.Enum),
		filterSortW:      new(widget.Enum),
	}
	fitlerSide.filterDirectionW.SetValue(defaultFilterDirection)
	fitlerSide.filterSortW.SetValue(defaultFilterSorter)

	txFilterDirection := []string{"All", "Sent", "Received", "Transfer"}
	txFilterSorts := []string{"Newest", "Oldest"}
	page.toFilters = make(map[string]decredmaterial.IconButton, len(txFilterSorts))

	for i := 0; i < len(txFilterDirection); i++ {
		fitlerSide.filterDirection = append(
			fitlerSide.filterDirection,
			common.theme.RadioButton(fmt.Sprint(i), txFilterDirection[i]))
	}

	for i := 0; i < len(txFilterSorts); i++ {
		if i == 0 {
			page.toFilters[fmt.Sprint(i)] = common.theme.IconButton(
				mustIcon(decredmaterial.NewIcon(icons.ContentFilterList)))
		} else {
			page.toFilters[fmt.Sprint(i)] = common.theme.IconButton(
				mustIcon(decredmaterial.NewIcon(icons.ContentSort)))
		}

		fitlerSide.filterSort = append(fitlerSide.filterSort,
			common.theme.RadioButton(fmt.Sprint(i), txFilterSorts[i]))
	}

	return func() {
		page.layout(common)
		page.handle(common)
		fitlerSide.layout(&common, &page)
		fitlerSide.handle(&common, &page)
	}
}

func (page *transactionsPage) layout(common pageCommon) {
	gtx := common.gtx

	container := func() {
		page.containerWidth = gtx.Constraints.Width.Max
		page.container.Layout(gtx,
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(txsPageInsetTop)}.Layout(gtx, func() {
					layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func() {
							layout.Inset{Left: unit.Dp(txsPageInsetLeft)}.Layout(gtx, func() {
								page.renderFiltererButton(&common)
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
							directionFilter, _ := strconv.Atoi(page.filterDirection)
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

func (page *transactionsPage) renderFiltererButton(common *pageCommon) {
	button := page.toFilters[page.filterSort]

	switch page.filterDirection {
	case "0":
		button.Background = common.theme.Color.Primary
	case "1":
		button.Background = common.theme.Color.Danger
	case "2":
		button.Background = common.theme.Color.Success
	case "3":
		button.Background = common.theme.Color.Hint
	default:
		button.Background = common.theme.Color.Hint
	}
	button.Layout(common.gtx, &page.toFiltersW)
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

func (fitlerSide *transactionsFiltersSide) layout(common *pageCommon, page *transactionsPage) {
	if !page.isShowFilters {
		return
	}

	gtx := common.gtx
	w := gtx.Constraints.Width.Max - page.containerWidth - 23

	body := func() {
		gtx.Constraints.Width.Min, gtx.Constraints.Width.Max = w, w

		{
			d := image.Point{X: w, Y: gtx.Constraints.Height.Max}
			dr := f32.Rectangle{
				Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
			}
			paint.ColorOp{Color: common.theme.Color.Surface}.Add(gtx.Ops)
			paint.PaintOp{Rect: dr}.Add(gtx.Ops)
			gtx.Dimensions = layout.Dimensions{Size: d}
		}

		layout.Inset{Left: unit.Dp(txsPageInsetLeft)}.Layout(gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(txsPageInsetTop)}.Layout(gtx, func() {})
				}),
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							common.theme.H5("Order").Layout(gtx)
						}),
						layout.Rigid(func() {
							(&layout.List{Axis: layout.Vertical}).
								Layout(gtx, len(fitlerSide.filterSort), func(index int) {
									fitlerSide.filterSort[index].Layout(gtx, fitlerSide.filterSortW)
								})
						}),
					)
				}),
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(txsPageInsetTop)}.Layout(gtx, func() {})
				}),
				layout.Rigid(func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							common.theme.H5("Direction").Layout(gtx)
						}),
						layout.Rigid(func() {
							(&layout.List{Axis: layout.Vertical}).
								Layout(gtx, len(fitlerSide.filterDirection), func(index int) {
									fitlerSide.filterDirection[index].Layout(gtx, fitlerSide.filterDirectionW)
								})
						}),
					)
				}),
			)
		})
	}

	layout.Inset{}.Layout(gtx, body)
}

func (page *transactionsPage) handle(common pageCommon) {
	if page.toReceiveW.Clicked(common.gtx) {
		*common.page = PageReceive
		return
	}

	if page.toSendW.Clicked(common.gtx) {
		*common.page = PageSend
		return
	}

	if page.toFiltersW.Clicked(common.gtx) {
		page.isShowFilters = !page.isShowFilters
	}

	if *page.keyEvent != nil && (*page.keyEvent).Name == key.NameEscape && page.isShowFilters {
		page.isShowFilters = false
	}
	*page.keyEvent = nil
}

func (fitlerSide *transactionsFiltersSide) handle(common *pageCommon, page *transactionsPage) {
	page.filterDirection = fitlerSide.filterDirectionW.Value(common.gtx)
	if page.filterSort != fitlerSide.filterSortW.Value(common.gtx) {
		page.filterSort = fitlerSide.filterSortW.Value(common.gtx)
		page.sortTransactions(common)
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
	newestFirst := page.filterSort == defaultFilterSorter

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
