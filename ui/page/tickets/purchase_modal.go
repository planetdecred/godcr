package tickets

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const purchaseModalID = "ticket_purchase_modal"

type ticketPurchaseModal struct {
	*load.Load

	balanceError     string
	ticketPrice      dcrutil.Amount
	totalCost        int64
	balanceLessCost  int64
	vspIsFetched     bool
	ticketsPurchased func()

	modal           decredmaterial.Modal
	tickets         decredmaterial.Editor
	rememberVSP     decredmaterial.CheckBoxStyle
	cancelPurchase  decredmaterial.Button
	reviewPurchase  decredmaterial.Button
	accountSelector *components.AccountSelector
	vspSelector     *vspSelector
}

func newTicketPurchaseModal(l *load.Load) *ticketPurchaseModal {
	tp := &ticketPurchaseModal{
		Load: l,

		tickets:        l.Theme.Editor(new(widget.Editor), ""),
		rememberVSP:    l.Theme.CheckBox(new(widget.Bool), "Remember VSP"),
		cancelPurchase: l.Theme.OutlineButton(new(widget.Clickable), "Cancel"),
		reviewPurchase: l.Theme.Button(new(widget.Clickable), "Review purchase"),
		modal:          *l.Theme.ModalFloatTitle(),
	}

	tp.reviewPurchase.SetEnabled(false)

	tp.vspIsFetched = len((*l.WL.VspInfo).List) > 0

	tp.tickets.Editor.SetText("1")
	return tp
}

func (tp *ticketPurchaseModal) TicketPurchased(ticketsPurchased func()) *ticketPurchaseModal {
	tp.ticketsPurchased = ticketsPurchased
	return tp
}

func (tp *ticketPurchaseModal) OnResume() {
	tp.initializeAccountSelector()
	err := tp.accountSelector.SelectFirstWalletValidAccount()
	if err != nil {
		tp.Toast.NotifyError(err.Error())
	}

	tp.vspSelector = newVSPSelector(tp.Load).title("Select a vsp")
	tp.ticketPrice = dcrutil.Amount(tp.WL.TicketPrice())

	if tp.vspIsFetched && components.StringNotEmpty(tp.WL.GetRememberVSP()) {
		tp.vspSelector.selectVSP(tp.WL.GetRememberVSP())
		tp.rememberVSP.CheckBox.Value = true
	}
}

func (tp *ticketPurchaseModal) Layout(gtx layout.Context) layout.Dimensions {
	l := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := tp.Icons.TicketPurchasedIcon
								return ic.Layout48dp(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
									return components.LayoutBalanceSize(gtx, tp.Load, tp.ticketPrice.String(), values.Size28)
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Flexed(.5, func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										tit := tp.Theme.Label(values.TextSize14, "Total")
										tit.Color = tp.Theme.Color.Gray3
										return tit.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return tp.Theme.Label(values.TextSize16, dcrutil.Amount(int64(tp.ticketPrice)*tp.ticketCount()).String()).Layout(gtx)
									}),
								)
							}),
							layout.Flexed(.5, tp.tickets.Layout),
						)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return tp.accountSelector.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if tp.balanceError == "" {
						return D{}
					}

					label := tp.Theme.Body1(tp.balanceError)
					label.Color = tp.Theme.Color.Orange
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return tp.vspSelector.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return tp.rememberVSP.Layout(gtx)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, tp.cancelPurchase.Layout)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return tp.reviewPurchase.Layout(gtx)
					}),
				)
			})
		},
	}

	return tp.modal.Layout(gtx, l, 850)
}

func (tp *ticketPurchaseModal) ticketCount() int64 {
	ticketCount, err := strconv.ParseInt(tp.tickets.Editor.Text(), 10, 64)
	if err != nil {
		return 0
	}

	return ticketCount
}

func (tp *ticketPurchaseModal) canPurchase() bool {
	tp.balanceError = ""
	if tp.ticketCount() < 1 {
		return false
	}

	if tp.vspSelector.selectedVSP == nil {
		return false
	}

	tp.calculateTotals()

	accountBalance := tp.accountSelector.SelectedAccount().Balance.Spendable
	if accountBalance < tp.totalCost || tp.balanceLessCost < 0 {
		tp.balanceError = "Insufficient funds"
		return false
	}

	return true
}

func (tp *ticketPurchaseModal) ModalID() string {
	return purchaseModalID
}

func (tp *ticketPurchaseModal) Show() {
	tp.ShowModal(tp)
}

func (tp *ticketPurchaseModal) Dismiss() {
	tp.DismissModal(tp)
}

func (tp *ticketPurchaseModal) initializeAccountSelector() {
	tp.accountSelector = components.NewAccountSelector(tp.Load).
		Title("Purchasing account").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			wal := tp.WL.MultiWallet.WalletWithID(account.WalletID)

			// Imported and watch only wallet accounts are invalid for sending
			accountIsValid := account.Number != dcrlibwallet.ImportedAccountNumber && !wal.IsWatchingOnlyWallet()

			if wal.ReadBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false) {
				// privacy is enabled for selected wallet

				accountIsValid = account.Number == wal.MixedAccountNumber()
			}
			return accountIsValid
		})
}

func (tp *ticketPurchaseModal) OnDismiss() {}

func (tp *ticketPurchaseModal) calculateTotals() {
	account := tp.accountSelector.SelectedAccount()
	wal := tp.WL.MultiWallet.WalletWithID(account.WalletID)

	ticketPrice, err := wal.TicketPrice()
	if err != nil {
		tp.Toast.NotifyError(err.Error())
		return
	}

	feePercentage := tp.vspSelector.selectedVSP.Info.FeePercentage
	total := ticketPrice.TicketPrice * tp.ticketCount()
	fee := int64((float64(total) / 100) * feePercentage)

	tp.totalCost = total + fee
	tp.balanceLessCost = account.Balance.Spendable - tp.totalCost
}

func (tp *ticketPurchaseModal) Handle() {
	tp.reviewPurchase.SetEnabled(tp.canPurchase())

	// reselect vsp if there's a delay in fetching the VSP List
	if !tp.vspIsFetched && len((*tp.WL.VspInfo).List) > 0 {
		if tp.WL.GetRememberVSP() != "" {
			tp.vspSelector.selectVSP(tp.WL.GetRememberVSP())
			tp.vspIsFetched = true
		}
	}

	if tp.cancelPurchase.Clicked() {
		tp.Dismiss()
	}

	if tp.canPurchase() && tp.reviewPurchase.Clicked() {

		if tp.vspSelector.Changed() && tp.rememberVSP.CheckBox.Value {
			tp.WL.RememberVSP(tp.vspSelector.selectedVSP.Host)
		} else if !tp.rememberVSP.CheckBox.Value {
			tp.WL.RememberVSP("")
		}

		selectedVSP := tp.vspSelector.SelectedVSP()
		account := tp.accountSelector.SelectedAccount()

		newTicketReviewModal(tp.Load, account, selectedVSP).
			TicketCount(tp.ticketCount()).
			TotalCost(tp.totalCost).
			BalanceLessCost(tp.balanceLessCost).
			TicketPurchased(func() {
				tp.Dismiss()
				tp.ticketsPurchased()
			}).
			Show()
	}
}
