package staking

import (
	"fmt"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const reviewModalID = "ticket_review_modal"

type stakeReviewModal struct {
	*load.Load
	account     *dcrlibwallet.Account
	selectedVSP *wallet.VSPInfo

	totalCost       int64
	ticketCount     int64
	balanceLessCost int64
	isLoading       bool

	materialLoader   material.LoaderStyle
	modal            decredmaterial.Modal
	spendingPassword decredmaterial.Editor
	purchase         decredmaterial.Button
	cancelPurchase   decredmaterial.Button
	ticketsPurchased func()
}

func newStakeReviewModal(l *load.Load, account *dcrlibwallet.Account, selectedVSP *wallet.VSPInfo) *stakeReviewModal {
	m := &stakeReviewModal{
		Load:             l,
		account:          account,
		selectedVSP:      selectedVSP,
		modal:            *l.Theme.ModalFloatTitle(),
		spendingPassword: l.Theme.EditorPassword(new(widget.Editor), "Spending password"),
		purchase:         l.Theme.Button("Purchase ticket"),
		cancelPurchase:   l.Theme.OutlineButton("Cancel"),
	}

	th := material.NewTheme(gofont.Collection())
	m.materialLoader = material.Loader(th)
	return m
}

func (t *stakeReviewModal) Layout(gtx layout.Context) layout.Dimensions {
	l := []layout.Widget{
		func(gtx C) D {
			return t.Theme.Label(values.TextSize20, "Confirm to purchase tickets").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					tleft := t.Theme.Label(values.TextSize14, "Amount")
					tleft.Color = t.Theme.Color.Gray2
					tright := t.Theme.Label(values.TextSize14, fmt.Sprintf("%d", t.ticketCount))
					return components.EndToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := t.Theme.Label(values.TextSize14, "Total cost")
					tleft.Color = t.Theme.Color.Gray2
					tright := t.Theme.Label(values.TextSize14, dcrutil.Amount(t.totalCost).String())
					return components.EndToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, t.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := t.Theme.Label(values.TextSize14, "Account")
					tleft.Color = t.Theme.Color.Gray2
					selectedWallet := t.WL.MultiWallet.WalletWithID(t.account.WalletID)
					tright := t.Theme.Label(values.TextSize14, selectedWallet.Name)
					return components.EndToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := t.Theme.Label(values.TextSize14, "Remaining")
					tleft.Color = t.Theme.Color.Gray2
					tright := t.Theme.Label(values.TextSize14, dcrutil.Amount(t.balanceLessCost).String())
					return components.EndToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, t.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := t.Theme.Label(values.TextSize14, "VSP")
					tleft.Color = t.Theme.Color.Gray2
					tright := t.Theme.Label(values.TextSize14, t.selectedVSP.Host)
					return components.EndToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(t.spendingPassword.Layout),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if t.isLoading {
							return D{}
						}

						return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, t.cancelPurchase.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if t.isLoading {
							return t.materialLoader.Layout(gtx)
						}

						return t.purchase.Layout(gtx)
					}),
				)
			})
		},
	}

	return t.modal.Layout(gtx, l, 850)
}

func (t *stakeReviewModal) ModalID() string {
	return reviewModalID
}
func (t *stakeReviewModal) OnDismiss() {}

func (t *stakeReviewModal) Show() {
	t.ShowModal(t)
}

func (t *stakeReviewModal) Dismiss() {
	t.DismissModal(t)
}

func (t *stakeReviewModal) OnResume() {}

func (t *stakeReviewModal) Handle() {
	for t.cancelPurchase.Clicked() {
		if !t.isLoading {
			t.Dismiss()
		}
	}

	for t.purchase.Clicked() {
		t.purchaseTickets()
	}
}

func (t *stakeReviewModal) purchaseTickets() {
	if t.isLoading {
		return
	}

	t.isLoading = true
	go func() {
		password := []byte(t.spendingPassword.Editor.Text())

		wal := t.WL.MultiWallet.WalletWithID(t.account.WalletID)

		defer func() {
			t.isLoading = false
		}()

		vsp, err := t.WL.MultiWallet.NewVSPClient(t.selectedVSP.Host, t.account.WalletID, uint32(t.account.Number))
		if err != nil {
			t.Toast.NotifyError(err.Error())
			return
		}

		err = vsp.PurchaseTickets(int32(t.ticketCount), wal.GetBestBlock()+256, password)
		if err != nil {
			t.Toast.NotifyError(err.Error())
			return
		}

		t.ticketsPurchased()
		t.Dismiss()
		t.Toast.Notify(fmt.Sprintf("%v ticket(s) purchased successfully", t.ticketCount))
	}()
}

func (t *stakeReviewModal) TicketCount(tickets int64) *stakeReviewModal {
	t.ticketCount = tickets
	if t.ticketCount > 1 {
		t.purchase.Text = fmt.Sprintf("Purchase %d tickets", t.ticketCount)
	}
	return t
}

func (t *stakeReviewModal) TotalCost(total int64) *stakeReviewModal {
	t.totalCost = total
	return t
}

func (t *stakeReviewModal) BalanceLessCost(remaining int64) *stakeReviewModal {
	t.balanceLessCost = remaining
	return t
}

func (t *stakeReviewModal) TicketPurchased(ticketsPurchased func()) *stakeReviewModal {
	t.ticketsPurchased = ticketsPurchased
	return t
}
