package ui

import (
	"image"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
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
	*pageCommon
	pageClosing  chan bool
	container    layout.Flex
	txsList      layout.List
	toTxnDetails []*gesture.Click
	separator    decredmaterial.Line
	theme        *decredmaterial.Theme

	orderDropDown  *decredmaterial.DropDown
	txTypeDropDown *decredmaterial.DropDown
	walletDropDown *decredmaterial.DropDown

	transactions []dcrlibwallet.Transaction
	wallets      []*dcrlibwallet.Wallet
}

func TransactionsPage(common *pageCommon) Page {
	pg := &transactionsPage{
		pageCommon:  common,
		pageClosing: make(chan bool, 1),
		container:   layout.Flex{Axis: layout.Vertical},
		txsList:     layout.List{Axis: layout.Vertical},
		separator:   common.theme.Separator(),
		theme:       common.theme,
	}

	pg.orderDropDown = createOrderDropDown(common)
	pg.txTypeDropDown = common.theme.DropDown([]decredmaterial.DropDownItem{
		{
			Text: values.String(values.StrAll),
		},
		{
			Text: values.String(values.StrSent),
		},
		{
			Text: values.String(values.StrReceived),
		},
		{
			Text: values.String(values.StrYourself),
		},
		{
			Text: values.String(values.StrStaking),
		},
	}, 1)

	return pg
}

func (pg *transactionsPage) OnResume() {
	pg.wallets = pg.sortedWalletList()
	pg.createOrUpdateWalletDropDown(&pg.walletDropDown, pg.wallets)
	pg.listenForTxNotifications()
	pg.loadTransactions()
}

func (pg *transactionsPage) loadTransactions() {
	selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
	newestFirst := pg.orderDropDown.SelectedIndex() == 0

	txFilter := dcrlibwallet.TxFilterAll
	switch pg.txTypeDropDown.SelectedIndex() {
	case 1:
		txFilter = dcrlibwallet.TxFilterSent
	case 2:
		txFilter = dcrlibwallet.TxFilterReceived
	case 3:
		txFilter = dcrlibwallet.TxFilterTransferred
	case 4:
		txFilter = dcrlibwallet.TxFilterStaking
	}

	wallTxs, err := selectedWallet.GetTransactionsRaw(0, 0, txFilter, newestFirst) //TODO
	if err != nil {
		log.Error("Error loading transactions:", err)
	} else {
		pg.transactions = wallTxs
	}
}

func (pg *transactionsPage) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.pageCommon

	container := func(gtx C) D {
		wallTxs := pg.transactions
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
									txt := common.theme.Body1(values.String(values.StrNoTransactionsYet))
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
									pg.goToTxnDetails(click.Events(gtx), &wallTxs[index])
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
			layout.Stacked(pg.dropDowns),
		)
	}
	return common.UniformPadding(gtx, container)
}

func (pg *transactionsPage) dropDowns(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Bottom: values.MarginPadding10,
	}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(pg.walletDropDown.Layout),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: values.MarginPadding5,
						}.Layout(gtx, pg.txTypeDropDown.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: values.MarginPadding5,
						}.Layout(gtx, pg.orderDropDown.Layout)
					}),
				)
			}),
		)
	})
}

func (pg *transactionsPage) handle() {
	for pg.txTypeDropDown.Changed() {
		pg.loadTransactions()
	}

	for pg.orderDropDown.Changed() {
		pg.loadTransactions()
	}

	for pg.walletDropDown.Changed() {
		pg.loadTransactions()
	}
}

func (pg *transactionsPage) goToTxnDetails(events []gesture.ClickEvent, txn *dcrlibwallet.Transaction) {
	for _, e := range events {
		if e.Type == gesture.TypeClick {
			pg.setReturnPage(PageTransactions)
			pg.changeFragment(TransactionDetailsPage(pg.pageCommon, txn), "txdetails")
		}
	}
}

func (pg *transactionsPage) listenForTxNotifications() {
	go func() {
		for {
			var notification interface{}

			select {
			case notification = <-pg.notificationsUpdate:
			case <-pg.pageClosing:
				return
			}

			switch n := notification.(type) {
			case wallet.NewTransaction:
				selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
				if selectedWallet.ID == n.Transaction.WalletID {
					pg.loadTransactions()
					pg.refreshWindow()
				}
			}
		}
	}()
}

func (pg *transactionsPage) onClose() {
	pg.pageClosing <- true
}

func initTxnWidgets(common *pageCommon, transaction *dcrlibwallet.Transaction) transactionWdg {

	var txn transactionWdg
	t := time.Unix(transaction.Timestamp, 0).UTC()
	txn.time = common.theme.Body1(t.Format(time.UnixDate))
	txn.status = common.theme.Body1("")
	txn.wallet = common.theme.Body2(common.multiWallet.WalletWithID(transaction.WalletID).Name)

	if txConfirmations(common, *transaction) > 1 {
		txn.status.Text = formatDateOrTime(transaction.Timestamp)
		txn.statusIcon = common.icons.confirmIcon
	} else {
		txn.status.Text = "pending"
		txn.status.Color = common.theme.Color.Gray
		txn.statusIcon = common.icons.pendingIcon
	}

	if transaction.Direction == dcrlibwallet.TxDirectionSent {
		txn.direction = common.icons.sendIcon
	} else {
		txn.direction = common.icons.receiveIcon
	}

	return txn
}
