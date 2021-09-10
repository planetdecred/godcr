package tickets

import (
	"fmt"
	"image/color"
	"strconv"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const purchaseModalID = "ticket_purchase_modal"

type ticketPurchaseModal struct {
	*load.Load

	ticketPrice       string
	totalCost         int64
	balanceLessCost   int64
	vspIsFetched      bool
	isPurchaseLoading bool

	modal          decredmaterial.Modal
	tickets        decredmaterial.Editor
	rememberVSP    decredmaterial.CheckBoxStyle
	selectVSP      []*gesture.Click
	cancelPurchase decredmaterial.Button
	reviewPurchase decredmaterial.Button

	accountSelector *components.AccountSelector
	vspSelector     *vspSelector

	vsp *dcrlibwallet.VSP
}

func newTicketPurchaseModal(l *load.Load) *ticketPurchaseModal {
	tp := &ticketPurchaseModal{
		Load: l,

		tickets:        l.Theme.Editor(new(widget.Editor), ""),
		rememberVSP:    l.Theme.CheckBox(new(widget.Bool), "Remember VSP"),
		cancelPurchase: l.Theme.Button(new(widget.Clickable), "Cancel"),
		reviewPurchase: l.Theme.Button(new(widget.Clickable), "Review purchase"),
		modal:          *l.Theme.ModalFloatTitle(),
	}

	tp.cancelPurchase.Background = color.NRGBA{}
	tp.cancelPurchase.Color = l.Theme.Color.Primary
	tp.vspIsFetched = len((*l.WL.VspInfo).List) > 0

	tp.tickets.Editor.SetText("1")
	return tp
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
									return components.LayoutBalance(gtx, tp.Load, tp.ticketPrice)
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(.5, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									tit := tp.Theme.Label(values.TextSize14, "Total")
									tit.Color = tp.Theme.Color.Gray2
									return tit.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return tp.Theme.Label(values.TextSize16, tp.ticketPrice).Layout(gtx)
								}),
							)
						}),
						layout.Flexed(.5, tp.tickets.Layout),
					)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return tp.accountSelector.Layout(gtx)
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
						if tp.canPurchase() {
							tp.reviewPurchase.Background = tp.Theme.Color.Primary
						} else {
							tp.reviewPurchase.Background = tp.Theme.Color.Hint
						}
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
	if tp.vspSelector.selectedVSP.Info == nil {
		return false
	}

	tp.calculateTotals()
	accountBalance := tp.accountSelector.SelectedAccount().Balance.Spendable
	if accountBalance < tp.totalCost || tp.balanceLessCost < 0 {
		return false
	}

	if tp.vspSelector.selectedVSP.Host == "" {
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

func (tp *ticketPurchaseModal) OnResume() {
	tp.initializeAccountSelector()
	err := tp.accountSelector.SelectFirstWalletValidAccount()
	if err != nil {
		tp.Toast.NotifyError(err.Error())
	}

	tp.vspSelector = newVSPSelector(tp.Load).title("Select a vsp")
	tp.ticketPrice = dcrutil.Amount(tp.WL.TicketPrice()).String()

	if tp.vspIsFetched && tp.WL.GetRememberVSP() != "" {
		tp.vspSelector.selectVSP(tp.WL.GetRememberVSP())
		tp.rememberVSP.CheckBox.Value = true
	}
}

func (tp *ticketPurchaseModal) initializeAccountSelector() {
	tp.accountSelector = components.NewAccountSelector(tp.Load).
		Title("Purchasing account").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			wal := tp.WL.MultiWallet.WalletWithID(account.WalletID)

			// Imported and watch only wallet accounts are invalid for sending
			accountIsValid := account.Number != maxInt32 && !wal.IsWatchingOnlyWallet()

			if wal.ReadBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false) {
				// privacy is enabled for selected wallet

				accountIsValid = account.Number == wal.MixedAccountNumber()
			}
			return accountIsValid
		})
}

func (tp *ticketPurchaseModal) OnDismiss() {}

func (tp *ticketPurchaseModal) calculateTotals() {
	accountBalance := tp.accountSelector.SelectedAccount().Balance.Spendable
	feePercentage := tp.vspSelector.selectedVSP.Info.FeePercentage
	total := tp.WL.TicketPrice() * tp.ticketCount()
	fee := int64((float64(total) / 100) * feePercentage)
	tp.totalCost = total + fee
	tp.balanceLessCost = accountBalance - tp.totalCost
}

func (tp *ticketPurchaseModal) createNewVSPD() {
	selectedAccount := tp.accountSelector.SelectedAccount()
	selectedVSP := tp.vspSelector.SelectedVSP()
	vspd, err := tp.WL.NewVSPD(selectedVSP.Host, selectedAccount.WalletID, selectedAccount.Number)
	if err != nil {
		tp.Toast.NotifyError(err.Error())
	}
	tp.vsp = vspd
}

func (tp *ticketPurchaseModal) purchaseTickets(password []byte) {
	tp.Dismiss()
	tp.Toast.Notify(fmt.Sprintf("attempting to purchase %v ticket(s)", tp.ticketCount()))

	go func() {
		account := tp.accountSelector.SelectedAccount()
		err := tp.WL.PurchaseTicket(account.WalletID, uint32(tp.ticketCount()), password, tp.vsp)
		if err != nil {
			tp.Toast.NotifyError(err.Error())
			return
		}
		tp.Toast.Notify(fmt.Sprintf("%v ticket(s) purchased successfully", tp.ticketCount()))
	}()
}

func (tp *ticketPurchaseModal) Handle() {
	// reselect vsp if there's a delay in fetching the VSP List
	if !tp.vspIsFetched && len((*tp.WL.VspInfo).List) > 0 {
		if tp.WL.GetRememberVSP() != "" {
			tp.vspSelector.selectVSP(tp.WL.GetRememberVSP())
			tp.vspIsFetched = true
		}
	}

	if tp.cancelPurchase.Button.Clicked() {
		tp.Dismiss()
	}

	if tp.reviewPurchase.Button.Clicked() && tp.canPurchase() {
		go tp.createNewVSPD()

		if tp.vspSelector.Changed() && tp.rememberVSP.CheckBox.Value {
			tp.WL.RememberVSP(tp.vspSelector.selectedVSP.Host)
		} else if !tp.rememberVSP.CheckBox.Value {
			tp.WL.RememberVSP("")
		}

		newTicketReviewModal(tp.Load).
			Account(tp.accountSelector.SelectedAccount()).
			VSPHost(tp.vspSelector.selectedVSP.Host).
			TicketCount(tp.ticketCount()).
			TotalCost(tp.totalCost).
			BalanceLessCost(tp.balanceLessCost).
			TicketPurchase(tp.purchaseTickets).
			Show()
	}
}

func (tp *ticketPurchaseModal) editorsNotEmpty(btn *decredmaterial.Button, editors ...*widget.Editor) bool {
	btn.Color = tp.Theme.Color.Surface
	for _, e := range editors {
		if e.Text() == "" {
			btn.Background = tp.Theme.Color.Hint
			return false
		}
	}

	btn.Background = tp.Theme.Color.Primary
	return true
}
