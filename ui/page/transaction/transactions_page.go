package transaction

import (
	"context"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const TransactionsPageID = "Transactions"

type (
	C = layout.Context
	D = layout.Dimensions
)

type TransactionsPage struct {
	*load.Load
	ctx       context.Context // page context
	ctxCancel context.CancelFunc
	separator decredmaterial.Line

	orderDropDown   *decredmaterial.DropDown
	txTypeDropDown  *decredmaterial.DropDown
	walletDropDown  *decredmaterial.DropDown
	transactionList *decredmaterial.ClickableList
	container       *widget.List
	transactions    []dcrlibwallet.Transaction
	wallets         []*dcrlibwallet.Wallet
}

func NewTransactionsPage(l *load.Load) *TransactionsPage {
	pg := &TransactionsPage{
		Load: l,
		container: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		separator:       l.Theme.Separator(),
		transactionList: l.Theme.NewClickableList(layout.Vertical),
	}

	pg.transactionList.Radius = decredmaterial.Radius(values.MarginPadding14.V)

	pg.orderDropDown = components.CreateOrderDropDown(l)
	pg.wallets = pg.WL.SortedWalletList()
	pg.walletDropDown = components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets)
	pg.txTypeDropDown = l.Theme.DropDown([]decredmaterial.DropDownItem{
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
			Text: "Mixed",
		},
		{
			Text: values.String(values.StrStaking),
		},
	}, 1)

	return pg
}

func (pg *TransactionsPage) ID() string {
	return TransactionsPageID
}

func (pg *TransactionsPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.listenForTxNotifications()
	pg.loadTransactions()
}

func (pg *TransactionsPage) loadTransactions() {
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
		txFilter = dcrlibwallet.TxFilterMixed
	case 5:
		txFilter = dcrlibwallet.TxFilterStaking
	}

	wallTxs, err := selectedWallet.GetTransactionsRaw(0, 0, txFilter, newestFirst) //TODO
	if err != nil {
		// log.Error("Error loading transactions:", err)
	} else {
		pg.transactions = wallTxs
	}
}

func (pg *TransactionsPage) Layout(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		wallTxs := pg.transactions
		return layout.Stack{Alignment: layout.N}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Inset{
					Top: values.MarginPadding60,
				}.Layout(gtx, func(gtx C) D {
					return pg.Theme.List(pg.container).Layout(gtx, 1, func(gtx C, i int) D {
						return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Card().Layout(gtx, func(gtx C) D {

								// return "No transactions yet" text if there are no transactions
								if len(wallTxs) == 0 {
									padding := values.MarginPadding16
									txt := pg.Theme.Body1(values.String(values.StrNoTransactionsYet))
									txt.Color = pg.Theme.Color.Gray2
									return layout.Center.Layout(gtx, func(gtx C) D {
										return layout.Inset{Top: padding, Bottom: padding}.Layout(gtx, txt.Layout)
									})
								}

								return pg.transactionList.Layout(gtx, len(wallTxs), func(gtx C, index int) D {
									var row = components.TransactionRow{
										Transaction: wallTxs[index],
										Index:       index,
										ShowBadge:   false,
									}

									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return components.LayoutTransactionRow(gtx, pg.Load, row)
										}),
										layout.Rigid(func(gtx C) D {
											// No divider for last row
											if row.Index == len(wallTxs)-1 {
												return layout.Dimensions{}
											}

											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											separator := pg.Theme.Separator()
											return layout.E.Layout(gtx, func(gtx C) D {
												// Show bottom divider for all rows except last
												return layout.Inset{Left: values.MarginPadding56}.Layout(gtx, separator.Layout)
											})
										}),
									)
								})
							})
						})
					})
				})
			}),
			layout.Expanded(func(gtx C) D {
				return pg.walletDropDown.Layout(gtx, 0, false)
			}),
			layout.Expanded(func(gtx C) D {
				return pg.orderDropDown.Layout(gtx, 0, true)
			}),
			layout.Expanded(func(gtx C) D {
				return pg.txTypeDropDown.Layout(gtx, pg.orderDropDown.Width-4, true)
			}),
		)
	}
	return components.UniformPadding(gtx, container)

}

func (pg *TransactionsPage) Handle() {

	for pg.txTypeDropDown.Changed() {
		pg.loadTransactions()
	}

	for pg.orderDropDown.Changed() {
		pg.loadTransactions()
	}

	for pg.walletDropDown.Changed() {
		pg.loadTransactions()
	}

	if clicked, selectedItem := pg.transactionList.ItemClicked(); clicked {
		pg.ChangeFragment(NewTransactionDetailsPage(pg.Load, &pg.transactions[selectedItem]))
	}
}

func (pg *TransactionsPage) listenForTxNotifications() {
	go func() {
		for {
			var notification interface{}

			select {
			case notification = <-pg.Receiver.NotificationsUpdate:
			case <-pg.ctx.Done():
				return
			}

			switch n := notification.(type) {
			case wallet.NewTransaction:
				selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
				if selectedWallet.ID == n.Transaction.WalletID {
					pg.loadTransactions()
					pg.RefreshWindow()
				}
			}
		}
	}()
}

func (pg *TransactionsPage) OnClose() {
	pg.ctxCancel()
}
