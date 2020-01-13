package pages

import (
	"sort"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/decred/dcrd/dcrutil"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/helper"
	"github.com/raedahgroup/godcr-gio/pages/common"
	"github.com/raedahgroup/godcr-gio/widgets"
)

type (
	OverviewPage struct {
		multiWallet *helper.MultiWallet
		syncer      *common.Syncer

		err          error
		totalBalance string

		transactions             []dcrlibwallet.Transaction
		seeAllTransactionsButton *widgets.Button
	}
)

func NewOverviewPage() *OverviewPage {
	return &OverviewPage{
		seeAllTransactionsButton: widgets.NewButton("See All", nil).SetBackgroundColor(helper.WhiteColor).SetColor(helper.DecredLightBlueColor),
	}
}

func (o *OverviewPage) BeforeRender(syncer *common.Syncer, multiWallet *helper.MultiWallet) {
	o.syncer = syncer
	o.multiWallet = multiWallet

	o.totalBalance, o.err = o.multiWallet.TotalBalance()

	// fetch recent transactions from all wallets
	transactions := []dcrlibwallet.Transaction{}

	for _, id := range o.multiWallet.WalletIDs {
		txns, err := o.multiWallet.WalletWithID(id).GetTransactionsRaw(0, 10, 0, true)
		if err != nil {
			o.err = err
			return
		}
		transactions = append(transactions, txns...)
	}
	sort.SliceStable(transactions, func(i, j int) bool {
		backTime := time.Unix(transactions[j].Timestamp, 0)
		frontTime := time.Unix(transactions[i].Timestamp, 0)
		return backTime.Before(frontTime)
	})

	if len(transactions) > 3 {
		transactions = transactions[:3]
	}
}

func (o *OverviewPage) Render(ctx *layout.Context, changePageFunc func(string)) {
	if o.err != nil {
		inset := layout.Inset{
			Left: unit.Dp(15),
		}
		inset.Layout(ctx, func() {
			widgets.NewErrorLabel(o.err.Error()).Draw(ctx)
		})
		return
	}

	inset := layout.Inset{
		Left: unit.Dp(15),
	}
	inset.Layout(ctx, func() {
		widgets.NewLabel(o.totalBalance).
			SetColor(helper.BlackColor).
			SetWeight(text.Bold).
			SetSize(7).
			Draw(ctx)
	})

	inset = layout.Inset{
		Top:  unit.Dp(40),
		Left: unit.Dp(15),
	}
	inset.Layout(ctx, func() {
		widgets.NewLabel("Current Total Balance").
			SetSize(5).
			SetColor(helper.GrayColor).
			Draw(ctx)
	})

	var nextTopInset float32 = 65

	inset = layout.Inset{
		Top: unit.Dp(nextTopInset),
	}
	inset.Layout(ctx, func() {
		if len(o.transactions) == 0 {
			o.drawNoTransactionsCard(ctx)
			nextTopInset += 85
		} else {
			o.drawRecentTransactionsCard(ctx, changePageFunc)
			nextTopInset += 195
		}
	})

	inset = layout.Inset{
		Top: unit.Dp(nextTopInset),
	}
	inset.Layout(ctx, func() {
		o.syncer.Render(ctx)
	})
}

func (o *OverviewPage) drawNoTransactionsCard(ctx *layout.Context) {
	helper.PaintArea(ctx, helper.WhiteColor, ctx.Constraints.Width.Max, 80)

	inset := layout.UniformInset(unit.Dp(15))
	inset.Layout(ctx, func() {
		widgets.NewLabel("Recent Transactions").
			SetSize(5).
			SetColor(helper.BlackColor).
			SetWeight(text.Bold).
			Draw(ctx)

		inset := layout.Inset{
			Top: unit.Dp(40),
		}
		inset.Layout(ctx, func() {
			widgets.NewLabel("No transactions yet").
				SetSize(4).
				SetColor(helper.GrayColor).
				SetWeight(text.Bold).
				Draw(ctx)
		})
	})
}

func (o *OverviewPage) drawRecentTransactionsCard(ctx *layout.Context, changePageFunc func(string)) {
	helper.PaintArea(ctx, helper.WhiteColor, ctx.Constraints.Width.Max, 190)

	inset := layout.UniformInset(unit.Dp(15))
	inset.Layout(ctx, func() {
		widgets.NewLabel("Recent Transactions").
			SetSize(5).
			SetColor(helper.BlackColor).
			SetWeight(text.Bold).
			Draw(ctx)

		inset := layout.Inset{
			Top: unit.Dp(5),
		}
		inset.Layout(ctx, func() {
			inset := layout.Inset{
				Top: unit.Dp(20),
			}
			inset.Layout(ctx, func() {
				(&layout.List{Axis: layout.Vertical}).Layout(ctx, len(o.transactions), func(i int) {
					inset := layout.Inset{
						Top: unit.Dp(5),
					}
					inset.Layout(ctx, func() {
						layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
							layout.Rigid(func() {
								transactionImage(ctx, o.transactions[i].Direction)
							}),
							layout.Rigid(func() {
								inset := layout.UniformInset(unit.Dp(6))
								inset.Layout(ctx, func() {
									widgets.NewLabel(dcrutil.Amount(o.transactions[i].Amount).String()).
										SetSize(4).
										Draw(ctx)
								})
							}),
							layout.Flexed(1, func() {
								layout.Align(layout.NE).Layout(ctx, func() {
									widgets.NewLabel(dcrlibwallet.ExtractDateOrTime(o.transactions[i].Timestamp)).
										SetColor(helper.GrayColor).
										Draw(ctx)
								})
							}),
						)
					})
				})
			})
		})

		inset = layout.Inset{
			Top: unit.Dp(137),
		}
		inset.Layout(ctx, func() {
			ctx.Constraints.Height.Min = 35
			o.seeAllTransactionsButton.Draw(ctx, func() {
				changePageFunc("transactions")
			})
		})
	})
}

func transactionImage(ctx *layout.Context, direction int32) {
	switch direction {
	case dcrlibwallet.TxDirectionSent:
		helper.SendImage.Layout(ctx)
	case dcrlibwallet.TxDirectionReceived:
		helper.ReceiveImage.Layout(ctx)
	case dcrlibwallet.TxDirectionTransferred:
		helper.ReceiveImage.Layout(ctx)
	default:
		helper.InfoImage.Layout(ctx)
	}
}
