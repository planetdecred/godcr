package overview

import (
	"gioui.org/layout"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *AppOverviewPage) initRecentTxWidgets() {
	pg.transactionsList = pg.Theme.NewClickableList(layout.Vertical)
	pg.transactionsList.Radius = decredmaterial.CornerRadius{
		BottomRight: values.MarginPadding14.V,
		BottomLeft:  values.MarginPadding14.V,
	}
	pg.toTransactions = pg.Theme.TextAndIconButton(values.String(values.StrSeeAll), pg.Icons.NavigationArrowForward)
	pg.toTransactions.Color = pg.Theme.Color.Primary
	pg.toTransactions.BackgroundColor = pg.Theme.Color.Surface
}

func (pg *AppOverviewPage) loadTransactions() {
	transactions, err := pg.WL.MultiWallet.GetTransactionsRaw(0, 5, dcrlibwallet.TxFilterAll, true)
	if err != nil {
		return
	}

	pg.transactions = transactions
}

// recentTransactionsSection lays out the list of recent transactions.
func (pg *AppOverviewPage) recentTransactionsSection(gtx layout.Context) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		padding := values.MarginPadding15
		return components.Container{Padding: layout.Inset{Top: padding}}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					title := pg.Theme.Body2(values.String(values.StrRecentTransactions))
					title.Color = pg.Theme.Color.GrayText1

					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					titlePadding := values.MarginPadding15
					return components.Container{Padding: layout.Inset{
						Left:   titlePadding,
						Right:  titlePadding,
						Bottom: titlePadding,
					}}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Rigid(title.Layout),
							layout.Rigid(func(gtx C) D {
								if len(pg.transactions) > 0 {
									return pg.toTransactions.Layout(gtx)
								}
								return D{}
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					if len(pg.transactions) == 0 {
						message := pg.Theme.Body1(values.String(values.StrNoTransactionsYet))
						message.Color = pg.Theme.Color.GrayText3
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Center.Layout(gtx, func(gtx C) D {
							return layout.Inset{
								Top:    values.MarginPadding10,
								Bottom: values.MarginPadding10,
							}.Layout(gtx, message.Layout)
						})
					}

					return pg.transactionsList.Layout(gtx, len(pg.transactions), func(gtx C, i int) D {
						var row = components.TransactionRow{
							Transaction: pg.transactions[i],
							Index:       i,
							ShowBadge:   len(pg.allWallets) > 1,
						}

						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return components.LayoutTransactionRow(gtx, pg.Load, row)
							}),
							layout.Rigid(func(gtx C) D {
								// No divider for last row
								if row.Index == len(pg.transactions)-1 {
									return layout.Dimensions{}
								}

								separator := pg.Theme.Separator()
								return layout.E.Layout(gtx, func(gtx C) D {
									// Show bottom divider for all rows except last
									return layout.Inset{Left: values.MarginPadding56}.Layout(gtx, separator.Layout)
								})
							}),
						)
					})
				}),
			)
		})
	})
}
