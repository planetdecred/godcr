package staking

import (
	"context"
	"fmt"
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const stakingModalID = "staking_modal"

type stakingModal struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	ticketPrice dcrutil.Amount

	modal          decredmaterial.Modal
	cancelPurchase decredmaterial.Button
	stakeBtn       decredmaterial.Button

	increment decredmaterial.IconButton
	decrement decredmaterial.IconButton

	spendingPassword decredmaterial.Editor
	tickets          decredmaterial.Editor

	accountSelector  *components.AccountSelector
	vspSelector      *components.VSPSelector
	materialLoader   material.LoaderStyle
	ticketsPurchased func()

	balanceError    string
	totalCost       int64
	balanceLessCost int64
	isLoading       bool
}

func newStakingModal(l *load.Load) *stakingModal {
	tp := &stakingModal{
		Load: l,

		tickets:          l.Theme.Editor(new(widget.Editor), ""),
		cancelPurchase:   l.Theme.OutlineButton("Cancel"),
		stakeBtn:         l.Theme.Button("Stake"),
		modal:            *l.Theme.ModalFloatTitle(),
		increment:        l.Theme.IconButton(l.Theme.Icons.ContentAdd),
		decrement:        l.Theme.IconButton(l.Theme.Icons.ContentRemove),
		spendingPassword: l.Theme.EditorPassword(new(widget.Editor), "Spending password"),
		materialLoader:   material.Loader(l.Theme.Base),
	}

	tp.tickets.Bordered = false
	tp.tickets.Editor.Alignment = text.Middle
	tp.tickets.Editor.SetText("1")

	tp.increment.ChangeColorStyle(&values.ColorStyle{Foreground: tp.Theme.Color.DeepBlue})
	tp.decrement.ChangeColorStyle(&values.ColorStyle{Foreground: tp.Theme.Color.Gray2})
	tp.increment.Size, tp.decrement.Size = values.TextSize18, values.TextSize18

	tp.modal.SetPadding(values.MarginPadding0)

	tp.stakeBtn.SetEnabled(false)

	return tp
}

func (tp *stakingModal) TicketPurchased(ticketsPurchased func()) *stakingModal {
	tp.ticketsPurchased = ticketsPurchased
	return tp
}

func (tp *stakingModal) OnResume() {
	tp.initializeAccountSelector()

	tp.ctx, tp.ctxCancel = context.WithCancel(context.TODO())

	tp.accountSelector.ListenForTxNotifications(tp.ctx)

	err := tp.accountSelector.SelectFirstWalletValidAccount(nil)
	if err != nil {
		tp.Toast.NotifyError(err.Error())
	}

	tp.vspSelector = components.NewVSPSelector(tp.Load).Title("Select a vsp")

	lastUsedVSP := tp.WL.MultiWallet.LastUsedVSP()
	if len(tp.WL.MultiWallet.KnownVSPs()) == 0 {
		// TODO: Does this modal need this list?
		go tp.WL.MultiWallet.ReloadVSPList(context.TODO())
	} else if components.StringNotEmpty(lastUsedVSP) {
		tp.vspSelector.SelectVSP(lastUsedVSP)
	}

	go func() {
		ticketPrice, err := tp.WL.MultiWallet.TicketPrice()
		if err != nil {
			tp.Toast.NotifyError(err.Error())
		} else {
			tp.ticketPrice = dcrutil.Amount(ticketPrice.TicketPrice)
			tp.RefreshWindow()
		}
	}()
}

func (tp *stakingModal) Layout(gtx layout.Context) layout.Dimensions {
	l := []layout.Widget{
		func(gtx C) D {
			return decredmaterial.LinearLayout{
				Orientation: layout.Vertical,
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.WrapContent,
				Padding: layout.Inset{
					Top:    values.MarginPadding24,
					Right:  values.MarginPadding24,
					Left:   values.MarginPadding24,
					Bottom: values.MarginPadding12,
				},
				Border: decredmaterial.Border{
					Radius: decredmaterial.TopRadius(14),
				},
				Direction:  layout.Center,
				Alignment:  layout.Middle,
				Background: tp.Theme.Color.Gray4,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := tp.Theme.Icons.NewStakeIcon
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
										totalLabel.Color = tp.Theme.Color.GrayText1
										return totalLabel.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										costLabel := tp.Theme.Label(values.TextSize16, dcrutil.Amount(int64(tp.ticketPrice)*tp.ticketCount()).String())
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
										Color:  tp.Theme.Color.Gray3,
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
										gtx.Constraints.Min.X, gtx.Constraints.Max.X = 90, 90
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
			return layout.Inset{
				Top:    values.MarginPadding2,
				Right:  values.MarginPadding14,
				Left:   values.MarginPadding14,
				Bottom: values.MarginPadding14,
			}.Layout(gtx, func(gtx C) D {
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
					layout.Rigid(func(gtx C) D {
						return tp.spendingPassword.Layout(gtx)
					}),
				)
			})
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

	// this is needed to generate the transaction fees before calculating
	// total ticket cost
	if tp.vspSelector.SelectedVSP() == nil {
		return false
	}

	tp.calculateTotals()

	accountBalance := tp.accountSelector.SelectedAccount().Balance.Spendable
	if accountBalance < tp.totalCost || tp.balanceLessCost < 0 {
		tp.balanceError = "Insufficient funds"
		return false
	}

	if tp.spendingPassword.Editor.Text() == "" {
		return false
	}

	return true
}

func (tp *stakingModal) ModalID() string {
	return stakingModalID
}

func (tp *stakingModal) Show() {
	tp.ShowModal(tp)
}

func (tp *stakingModal) initializeAccountSelector() {
	tp.accountSelector = components.NewAccountSelector(tp.Load, nil).
		Title("Purchasing account").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			wal := tp.WL.MultiWallet.WalletWithID(account.WalletID)

			// Imported and watch only wallet accounts are invalid for sending
			accountIsValid := account.Number != dcrlibwallet.ImportedAccountNumber && !wal.IsWatchingOnlyWallet()

			if wal.ReadBoolConfigValueForKey(load.SpendUnmixedFundsKey, false) {
				// Spending from unmixed accounts is disabled for wallet

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

	feePercentage := tp.vspSelector.SelectedVSP().FeePercentage
	total := ticketPrice.TicketPrice * tp.ticketCount()
	fee := int64((float64(total) / 100) * feePercentage)

	tp.totalCost = total + fee
	tp.balanceLessCost = account.Balance.Spendable - tp.totalCost
}

func (tp *stakingModal) Handle() {
	tp.stakeBtn.SetEnabled(tp.canPurchase())

	if tp.vspSelector.Changed() {
		tp.WL.MultiWallet.SaveLastUsedVSP(tp.vspSelector.SelectedVSP().Host)
	}

	_, isChanged := decredmaterial.HandleEditorEvents(tp.spendingPassword.Editor)
	if isChanged {
		tp.spendingPassword.SetError("")
	}

	// reselect vsp if there's a delay in fetching the VSP List
	lastUsedVSP := tp.WL.MultiWallet.LastUsedVSP()
	if len(tp.WL.MultiWallet.KnownVSPs()) > 0 && lastUsedVSP != "" {
		tp.vspSelector.SelectVSP(lastUsedVSP)
	}

	if tp.cancelPurchase.Clicked() {
		if tp.isLoading {
			return
		}

		tp.Dismiss()
	}

	if tp.modal.BackdropClicked(true) {
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
		tp.decrement.ChangeColorStyle(&values.ColorStyle{Foreground: tp.Theme.Color.Text})
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
			tp.decrement.ChangeColorStyle(&values.ColorStyle{Foreground: tp.Theme.Color.Gray2})
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
	tp.modal.SetDisabled(true)
	go func() {
		password := []byte(tp.spendingPassword.Editor.Text())

		account := tp.accountSelector.SelectedAccount()
		wal := tp.WL.MultiWallet.WalletWithID(account.WalletID)

		selectedVSP := tp.vspSelector.SelectedVSP()

		defer func() {
			tp.isLoading = false
			tp.modal.SetDisabled(false)
		}()

		vspHost, vspPubKey := selectedVSP.Host, selectedVSP.PubKey
		_, err := wal.PurchaseTickets(account.Number, int32(tp.ticketCount()), vspHost, vspPubKey, password)
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

func (tp *stakingModal) Dismiss() {
	tp.ctxCancel()
	tp.DismissModal(tp)
}
