package staking

import (
	"fmt"
	"strconv"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const purchaseModalID = "ticket_purchase_modal"

type stakingModal struct {
	*load.Load

	balanceError     string
	ticketPrice      dcrutil.Amount
	totalCost        int64
	balanceLessCost  int64
	vspIsFetched     bool
	isLoading        bool
	ticketsPurchased func()

	modal          decredmaterial.Modal
	cancelPurchase decredmaterial.Button
	stakeBtn       decredmaterial.Button

	increment decredmaterial.IconButton
	decrement decredmaterial.IconButton

	spendingPassword decredmaterial.Editor
	tickets          decredmaterial.Editor
	materialLoader   material.LoaderStyle

	accountSelector *components.AccountSelector
	vspSelector     *vspSelector
}

func newStakingModal(l *load.Load) *stakingModal {
	tp := &stakingModal{
		Load: l,

		tickets:          l.Theme.Editor(new(widget.Editor), ""),
		cancelPurchase:   l.Theme.OutlineButton("Cancel"),
		stakeBtn:         l.Theme.Button("Stake"),
		modal:            *l.Theme.ModalFloatTitle(),
		increment:        l.Theme.PlainIconButton(l.Icons.ContentAdd),
		decrement:        l.Theme.PlainIconButton(l.Icons.ContentRemove),
		spendingPassword: l.Theme.EditorPassword(new(widget.Editor), "Spending password"),
		materialLoader:   material.Loader(material.NewTheme(gofont.Collection())),
	}

	tp.tickets.Bordered = false
	tp.tickets.Editor.Alignment = text.Middle
	tp.tickets.Editor.SetText("1")

	tp.increment.Color, tp.decrement.Color = l.Theme.Color.Text, l.Theme.Color.InactiveGray
	tp.increment.Size, tp.decrement.Size = values.TextSize18, values.TextSize18

	tp.modal.SetPadding(values.MarginPadding0)

	tp.stakeBtn.SetEnabled(false)

	tp.vspIsFetched = len((*l.WL.VspInfo).List) > 0

	return tp
}

func (tp *stakingModal) TicketPurchased(ticketsPurchased func()) *stakingModal {
	tp.ticketsPurchased = ticketsPurchased
	return tp
}

func (tp *stakingModal) OnResume() {
	tp.initializeAccountSelector()
	err := tp.accountSelector.SelectFirstWalletValidAccount()
	if err != nil {
		tp.Toast.NotifyError(err.Error())
	}

	tp.vspSelector = newVSPSelector(tp.Load).title("Select a vsp")
	tp.ticketPrice = dcrutil.Amount(tp.WL.TicketPrice())

	if tp.vspIsFetched && components.StringNotEmpty(tp.WL.GetRememberVSP()) {
		tp.vspSelector.selectVSP(tp.WL.GetRememberVSP())
	}
}

func (tp *stakingModal) Layout(gtx layout.Context) layout.Dimensions {
	l := []layout.Widget{
		func(gtx C) D {
			return decredmaterial.LinearLayout{
				Orientation: layout.Vertical,
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.WrapContent,
				Padding:     layout.UniformInset(values.MarginPadding16),
				Border: decredmaterial.Border{
					Radius: decredmaterial.TopRadius(14),
				},
				Direction:  layout.Center,
				Alignment:  layout.Middle,
				Background: tp.Theme.Color.LightGray,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := tp.Icons.NewStakeIcon
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
										totalLabel := tp.Theme.Label(values.TextSize14, "Total")
										totalLabel.Color = tp.Theme.Color.Gray3
										return totalLabel.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										costLabel := tp.Theme.Label(values.TextSize16, dcrutil.Amount(int64(tp.ticketPrice)*tp.ticketCount()).String())
										costLabel.Color = tp.Theme.Color.Gray6
										return costLabel.Layout(gtx)
									}),
								)
							}),
							layout.Flexed(.5, func(gtx C) D {
								return decredmaterial.LinearLayout{
									Orientation: layout.Horizontal,
									Width:       decredmaterial.WrapContent,
									Height:      decredmaterial.WrapContent,
									Border: decredmaterial.Border{
										Radius: decredmaterial.Radius(10),
										Color:  tp.Theme.Color.InactiveGray,
										Width:  values.MarginPadding1,
									},
									Direction:  layout.E,
									Alignment:  layout.Middle,
									Background: tp.Theme.Color.Surface,
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return tp.decrement.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										gtx.Constraints.Min.X, gtx.Constraints.Max.X = 100, 100
										return tp.tickets.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return tp.increment.Layout(gtx)
									}),
								)
							}),
						)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(tp.accountSelector.Layout),
				layout.Rigid(func(gtx C) D {
					if tp.balanceError == "" {
						return D{}
					}

					label := tp.Theme.Caption(tp.balanceError)
					label.Color = tp.Theme.Color.Danger
					return label.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, func(gtx C) D {
						return tp.vspSelector.Layout(gtx)
					})
				}),
				layout.Rigid(tp.spendingPassword.Layout),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if tp.isLoading {
							return D{}
						}
						return layout.Inset{
							Right:  values.MarginPadding4,
							Bottom: values.MarginPadding15,
						}.Layout(gtx, tp.cancelPurchase.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if tp.isLoading {
							return layout.Inset{
								Top:    values.MarginPadding10,
								Bottom: values.MarginPadding15,
							}.Layout(gtx, tp.materialLoader.Layout)
						}
						return tp.stakeBtn.Layout(gtx)
					}),
				)
			})
		},
	}

	return tp.modal.Layout(gtx, l)
}

func (tp *stakingModal) ticketCount() int64 {
	ticketCount, err := strconv.ParseInt(tp.tickets.Editor.Text(), 10, 64)
	if err != nil {
		return 0
	}

	return ticketCount
}

func (tp *stakingModal) canPurchase() bool {
	tp.balanceError = ""
	if tp.ticketCount() < 1 {
		return false
	}

	if tp.vspSelector.selectedVSP == nil {
		return false
	}

	if tp.spendingPassword.Editor.Text() == "" {
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

func (tp *stakingModal) ModalID() string {
	return purchaseModalID
}

func (tp *stakingModal) Show() {
	tp.ShowModal(tp)
}

func (tp *stakingModal) Dismiss() {
	tp.DismissModal(tp)
}

func (tp *stakingModal) initializeAccountSelector() {
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

func (tp *stakingModal) OnDismiss() {}

func (tp *stakingModal) calculateTotals() {
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

func (tp *stakingModal) Handle() {
	tp.stakeBtn.SetEnabled(tp.canPurchase())

	if tp.vspSelector.Changed() {
		tp.WL.RememberVSP(tp.vspSelector.selectedVSP.Host)
	}

	_, isChanged := decredmaterial.HandleEditorEvents(tp.spendingPassword.Editor)
	if isChanged {
		tp.spendingPassword.SetError("")
	}

	// reselect vsp if there's a delay in fetching the VSP List
	if !tp.vspIsFetched && len((*tp.WL.VspInfo).List) > 0 {
		if tp.WL.GetRememberVSP() != "" {
			tp.vspSelector.selectVSP(tp.WL.GetRememberVSP())
			tp.vspIsFetched = true
		}
	}

	if tp.cancelPurchase.Clicked() {
		if tp.isLoading {
			return
		}

		tp.Dismiss()
	}

	if tp.canPurchase() && tp.stakeBtn.Clicked() {
		tp.purchaseTickets()
	}

	// increment the ticket value
	if tp.increment.Button.Clicked() {
		tp.decrement.Color = tp.Theme.Color.Text
		value, err := strconv.Atoi(tp.tickets.Editor.Text())
		if err != nil {
			return
		}
		value++
		tp.tickets.Editor.SetText(fmt.Sprintf("%d", value))
	}

	// decrement the ticket value
	if tp.decrement.Button.Clicked() {
		value, err := strconv.Atoi(tp.tickets.Editor.Text())
		if err != nil {
			return
		}
		value--
		if value < 1 {
			tp.decrement.Color = tp.Theme.Color.InactiveGray
			return
		}
		tp.tickets.Editor.SetText(fmt.Sprintf("%d", value))
	}
}

func (tp *stakingModal) purchaseTickets() {

	if tp.isLoading {
		return
	}

	tp.isLoading = true
	go func() {
		password := []byte(tp.spendingPassword.Editor.Text())

		account := tp.accountSelector.SelectedAccount()
		wal := tp.WL.MultiWallet.WalletWithID(account.WalletID)

		selectedVSP := tp.vspSelector.SelectedVSP()

		defer func() {
			tp.isLoading = false
		}()

		vsp, err := tp.WL.MultiWallet.NewVSPClient(selectedVSP.Host, account.WalletID, uint32(account.Number))
		if err != nil {
			tp.Toast.NotifyError(err.Error())
			return
		}

		err = vsp.PurchaseTickets(int32(tp.ticketCount()), wal.GetBestBlock()+256, password)
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				tp.spendingPassword.SetError("Invalid password")
			} else {
				tp.Toast.NotifyError(err.Error())
			}
			return
		}

		tp.ticketsPurchased()
		tp.Dismiss()
	}()
}
