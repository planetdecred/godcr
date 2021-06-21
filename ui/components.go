// components contain layout code that are shared by multiple pages but aren't widely used enough to be defined as
// widgets

package ui

import (
	"time"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	purchasingAccountTitle = "Purchasing account"
	sendingAccountTitle    = "Sending account"
	receivingAccountTitle  = "Receiving account"
	ticketAge              = "Ticket age"
)

type (
	TransactionRow struct {
		transaction dcrlibwallet.Transaction
		index       int
		showBadge   bool
	}
	toast struct {
		text    string
		success bool
		timer   *time.Timer
	}
	tooltips struct {
		statusTooltip     *decredmaterial.Tooltip
		walletNameTooltip *decredmaterial.Tooltip
		dateTooltip       *decredmaterial.Tooltip
		daysBehindTooltip *decredmaterial.Tooltip
	}
)

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page *pageCommon) layoutBalance(gtx layout.Context, amount string) layout.Dimensions {
	// todo: make "DCR" symbols small when there are no decimals in the balance
	mainText, subText := breakBalance(page.printer, amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			label := page.theme.Label(values.TextSize20, mainText)
			label.Color = page.theme.Color.DeepBlue
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			label := page.theme.Label(values.TextSize14, subText)
			label.Color = page.theme.Color.DeepBlue
			return label.Layout(gtx)
		}),
	)
}

var (
	navDrawerWidth          = unit.Value{U: unit.UnitDp, V: 160}
	navDrawerMinimizedWidth = unit.Value{U: unit.UnitDp, V: 72}
)

// transactionRow is a single transaction row on the transactions and overview page. It lays out a transaction's
// direction, balance, status.
func transactionRow(gtx layout.Context, common *pageCommon, row TransactionRow) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	directionIconTopMargin := values.MarginPadding16

	if row.index == 0 && row.showBadge {
		directionIconTopMargin = values.MarginPadding14
	} else if row.index == 0 {
		// todo: remove top margin from container
		directionIconTopMargin = values.MarginPadding0
	}

	wal := common.multiWallet.WalletWithID(row.transaction.WalletID)

	return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				icon := common.icons.receiveIcon
				if row.transaction.Direction == dcrlibwallet.TxDirectionSent {
					icon = common.icons.sendIcon
				}
				icon.Scale = 1.0

				return layout.Inset{Top: directionIconTopMargin}.Layout(gtx, func(gtx C) D {
					return icon.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if row.index == 0 {
							return layout.Dimensions{}
						}
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						separator := common.theme.Separator()
						separator.Width = gtx.Constraints.Max.X - gtx.Px(unit.Dp(16))
						return layout.E.Layout(gtx, func(gtx C) D {
							// Todo: add comment
							marginBottom := values.MarginPadding16
							if row.showBadge {
								marginBottom = values.MarginPadding5
							}
							return layout.Inset{Bottom: marginBottom}.Layout(gtx,
								func(gtx C) D {
									return separator.Layout(gtx)
								})
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Inset{}.Layout(gtx, func(gtx C) D {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Spacing:   layout.SpaceBetween,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
										return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												return common.layoutBalance(gtx, dcrutil.Amount(row.transaction.Amount).String())
											}),
											layout.Rigid(func(gtx C) D {
												if row.showBadge {
													return walletLabel(gtx, common, wal.Name)
												}
												return layout.Dimensions{}
											}),
										)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
												func(gtx C) D {
													status := common.theme.Body1("pending")
													if txConfirmations(common, row.transaction) <= 1 {
														status.Color = common.theme.Color.Gray5
													} else {
														status.Color = common.theme.Color.Gray4
														status.Text = formatDateOrTime(row.transaction.Timestamp)
													}
													status.Alignment = text.Middle
													return status.Layout(gtx)
												})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
												statusIcon := common.icons.confirmIcon
												if txConfirmations(common, row.transaction) <= 1 {
													statusIcon = common.icons.pendingIcon
												}
												statusIcon.Scale = 1.0
												return statusIcon.Layout(gtx)
											})
										}),
									)
								}),
							)
						})
					}),
				)
			}),
		)
	})
}

func txConfirmations(common *pageCommon, transaction dcrlibwallet.Transaction) int32 {
	if transaction.BlockHeight != -1 {
		// TODO
		return (common.multiWallet.WalletWithID(transaction.WalletID).GetBestBlock() - transaction.BlockHeight) + 1
	}

	return 0
}

// walletLabel displays the wallet which a transaction belongs to. It is only displayed on the overview page when there
// are transactions from multiple wallets
func walletLabel(gtx layout.Context, c *pageCommon, walletName string) D {
	return decredmaterial.Card{
		Color: c.theme.Color.LightGray,
	}.Layout(gtx, func(gtx C) D {
		return Container{
			layout.Inset{
				Left:  values.MarginPadding4,
				Right: values.MarginPadding4,
			}}.Layout(gtx, func(gtx C) D {
			name := c.theme.Label(values.TextSize12, walletName)
			name.Color = c.theme.Color.Gray
			return name.Layout(gtx)
		})
	})
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func endToEndRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(leftWidget),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, rightWidget)
		}),
	)
}

func (page *pageCommon) Modal(gtx layout.Context, body layout.Dimensions, modal layout.Dimensions) layout.Dimensions {
	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return body
		}),
		layout.Stacked(func(gtx C) D {
			return modal
		}),
	)
	return dims
}

func (page *pageCommon) initSelectAccountWidget(wallAcct map[int][]walletAccount, windex int) {
	if _, ok := wallAcct[windex]; !ok {
		accounts := page.info.Wallets[windex].Accounts
		if len(accounts) != len(wallAcct[windex]) {
			wallAcct[windex] = make([]walletAccount, len(accounts))
			for aindex := range accounts {
				if accounts[aindex].Name == "imported" {
					continue
				}

				wallAcct[windex][aindex] = walletAccount{
					walletIndex:  windex,
					accountIndex: aindex,
					evt:          &gesture.Click{},
					accountName:  accounts[aindex].Name,
					totalBalance: accounts[aindex].TotalBalance,
					spendable:    dcrutil.Amount(accounts[aindex].SpendableBalance).String(),
					number:       accounts[aindex].Number,
				}
			}
		}
	}
}

func displayToast(th *decredmaterial.Theme, gtx layout.Context, n *toast) layout.Dimensions {
	color := th.Color.Success
	if !n.success {
		color = th.Color.Danger
	}

	card := th.Card()
	card.Color = color
	return card.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top: values.MarginPadding7, Bottom: values.MarginPadding7,
			Left: values.MarginPadding15, Right: values.MarginPadding15,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			t := th.Body1(n.text)
			t.Color = th.Color.Surface
			return t.Layout(gtx)
		})
	})
}

func (page *pageCommon) handleToast() {
	if (*page.toast) == nil {
		return
	}

	if (*page.toast).timer == nil {
		(*page.toast).timer = time.NewTimer(time.Second * 3)
	}

	select {
	case <-(*page.toast).timer.C:
		*page.toast = nil
	default:
	}
}

func (page *pageCommon) handler() {
	page.handleToast()
}
