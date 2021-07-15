package ui

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
	"github.com/planetdecred/godcr/ui/values"
)

const TicketPurchaseModalID = "ticket_purchase_modal"

type ticketPurchaseModal struct {
	*pageCommon

	randomID        string
	ticketPrice     string
	totalCost       int64
	balanceLessCost int64
	vspIsFetched    bool

	modal          decredmaterial.Modal
	tickets        decredmaterial.Editor
	rememberVSP    decredmaterial.CheckBoxStyle
	selectVSP      []*gesture.Click
	cancelPurchase decredmaterial.Button
	reviewPurchase decredmaterial.Button

	accountSelector *accountSelector
	vspSelector     *vspSelector

	vsp *dcrlibwallet.VSP
}

func newTicketPurchaseModal(common *pageCommon) *ticketPurchaseModal {
	tp := &ticketPurchaseModal{
		pageCommon: common,

		randomID:       fmt.Sprintf("%s-%d", TicketPurchaseModalID, generateRandomNumber()),
		tickets:        common.theme.Editor(new(widget.Editor), ""),
		rememberVSP:    common.theme.CheckBox(new(widget.Bool), "Remember VSP"),
		cancelPurchase: common.theme.Button(new(widget.Clickable), "Cancel"),
		reviewPurchase: common.theme.Button(new(widget.Clickable), "Review purchase"),
		modal:          *common.theme.ModalFloatTitle(),
	}

	tp.cancelPurchase.Background = color.NRGBA{}
	tp.cancelPurchase.Color = common.theme.Color.Primary
	tp.vspIsFetched = len((*common.vspInfo).List) > 0

	tp.tickets.Editor.SetText("1")
	return tp
}

func (tp *ticketPurchaseModal) createNewVSPD() {
	selectedAccount := tp.accountSelector.selectedAccount
	selectedVSP := tp.vspSelector.SelectedVSP()
	vspd, err := tp.wallet.NewVSPD(selectedVSP.Host, selectedAccount.WalletID, selectedAccount.Number)
	if err != nil {
		tp.notify(err.Error(), false)
	}
	tp.vsp = vspd
}

func (tp *ticketPurchaseModal) Layout(gtx layout.Context) layout.Dimensions {
	l := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := tp.icons.ticketPurchasedIcon
								ic.Scale = 1.2
								return ic.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
									return tp.layoutBalance(gtx, tp.ticketPrice, true)
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
									tit := tp.theme.Label(values.TextSize14, "Total")
									tit.Color = tp.theme.Color.Gray2
									return tit.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return tp.theme.Label(values.TextSize16, tp.ticketPrice).Layout(gtx)
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
							tp.reviewPurchase.Background = tp.theme.Color.Primary
						} else {
							tp.reviewPurchase.Background = tp.theme.Color.Hint
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
	accountBalance := tp.accountSelector.selectedAccount.Balance.Spendable
	if accountBalance < tp.totalCost || tp.balanceLessCost < 0 {
		return false
	}

	if tp.vspSelector.selectedVSP.Host == "" {
		return false
	}

	return true
}

func (tp *ticketPurchaseModal) modalID() string {
	return tp.randomID
}

func (tp *ticketPurchaseModal) Show() {
	tp.showModal(tp)
}

func (tp *ticketPurchaseModal) Dismiss() {
	tp.dismissModal(tp)
}

func (tp *ticketPurchaseModal) OnResume() {
	tp.initializeAccountSelector()
	err := tp.accountSelector.selectFirstWalletValidAccount()
	if err != nil {
		log.Error(err)
	}

	tp.vspSelector = newVSPSelector(tp.pageCommon).title("Select a vsp")
	tp.ticketPrice = dcrutil.Amount(tp.wallet.TicketPrice()).String()

	if tp.vspIsFetched && tp.wallet.GetRememberVSP() != "" {
		tp.vspSelector.selectVSP(tp.wallet.GetRememberVSP())
		tp.rememberVSP.CheckBox.Value = true
	}
}

func (tp *ticketPurchaseModal) initializeAccountSelector() {
	tp.accountSelector = newAccountSelector(tp.pageCommon).
		title("Purchasing account").
		accountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		accountValidator(func(account *dcrlibwallet.Account) bool {
			wal := tp.multiWallet.WalletWithID(account.WalletID)

			// Imported and watch only wallet accounts are invalid for sending
			accountIsValid := account.Number != MaxInt32 && !wal.IsWatchingOnlyWallet()

			if wal.ReadBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false) {
				// privacy is enabled for selected wallet

				accountIsValid = account.Number == wal.MixedAccountNumber()
			}
			return accountIsValid
		})
}

func (tp *ticketPurchaseModal) OnDismiss() {}

func (tp *ticketPurchaseModal) calculateTotals() {
	accountBalance := tp.accountSelector.selectedAccount.Balance.Spendable
	feePercentage := tp.vspSelector.selectedVSP.Info.FeePercentage
	total := tp.wallet.TicketPrice() * tp.ticketCount()
	fee := int64((float64(total) / 100) * feePercentage)
	tp.totalCost = total + fee
	tp.balanceLessCost = accountBalance - tp.totalCost
}

func (tp *ticketPurchaseModal) handle() {
	// reselect vsp if there's a delay in fetching the VSP List
	if !tp.vspIsFetched && len((*tp.vspInfo).List) > 0 {
		if tp.wallet.GetRememberVSP() != "" {
			tp.vspSelector.selectVSP(tp.wallet.GetRememberVSP())
			tp.vspIsFetched = true
		}
	}

	if tp.cancelPurchase.Button.Clicked() {
		tp.Dismiss()
	}

	if tp.reviewPurchase.Button.Clicked() && tp.canPurchase() {
		go tp.createNewVSPD()

		if tp.vspSelector.Changed() && tp.rememberVSP.CheckBox.Value {
			tp.wallet.RememberVSP(tp.vspSelector.selectedVSP.Host)
		} else if !tp.rememberVSP.CheckBox.Value {
			tp.wallet.RememberVSP("")
		}

		newTicketReviewModal(tp.pageCommon).
			Account(tp.accountSelector.selectedAccount).
			VSPHost(tp.vspSelector.selectedVSP.Host).
			TicketCount(tp.ticketCount()).
			TotalCost(tp.totalCost).
			BalanceLessCost(tp.balanceLessCost).
			Show()
	}
}

func (tp *ticketPurchaseModal) editorsNotEmpty(btn *decredmaterial.Button, editors ...*widget.Editor) bool {
	btn.Color = tp.theme.Color.Surface
	for _, e := range editors {
		if e.Text() == "" {
			btn.Background = tp.theme.Color.Hint
			return false
		}
	}

	btn.Background = tp.theme.Color.Primary
	return true
}
