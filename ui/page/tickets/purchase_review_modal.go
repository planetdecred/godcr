package tickets

import (
	"fmt"
	"image/color"

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

type ticketReviewModal struct {
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

func newTicketReviewModal(l *load.Load, account *dcrlibwallet.Account, selectedVSP *wallet.VSPInfo) *ticketReviewModal {
	m := &ticketReviewModal{
		Load:             l,
		account:          account,
		selectedVSP:      selectedVSP,
		modal:            *l.Theme.ModalFloatTitle(),
		spendingPassword: l.Theme.EditorPassword(new(widget.Editor), "Spending password"),
		purchase:         l.Theme.Button(new(widget.Clickable), "Purchase ticket"),
		cancelPurchase:   l.Theme.Button(new(widget.Clickable), "Cancel"),
	}

	th := material.NewTheme(gofont.Collection())
	m.materialLoader = material.Loader(th)

	m.purchase.Background = m.Theme.Color.Primary
	m.cancelPurchase.Background, m.cancelPurchase.Color = color.NRGBA{}, l.Theme.Color.Primary
	return m
}

func (t *ticketReviewModal) Layout(gtx layout.Context) layout.Dimensions {
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

func (t *ticketReviewModal) ModalID() string {
	return reviewModalID
}
func (t *ticketReviewModal) OnDismiss() {}

func (t *ticketReviewModal) Show() {
	t.ShowModal(t)
}

func (t *ticketReviewModal) Dismiss() {
	t.DismissModal(t)
}

func (t *ticketReviewModal) OnResume() {}

func (t *ticketReviewModal) Handle() {
	for t.cancelPurchase.Clicked() {
		if !t.isLoading {
			t.Dismiss()
		}
	}

	for t.purchase.Clicked() {
		t.purchaseTickets()
	}
}

func (t *ticketReviewModal) purchaseTickets() {
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

func (t *ticketReviewModal) TicketCount(tickets int64) *ticketReviewModal {
	t.ticketCount = tickets
	if t.ticketCount > 1 {
		t.purchase.Text = fmt.Sprintf("Purchase %d tickets", t.ticketCount)
	}
	return t
}

func (t *ticketReviewModal) TotalCost(total int64) *ticketReviewModal {
	t.totalCost = total
	return t
}

func (t *ticketReviewModal) BalanceLessCost(remaining int64) *ticketReviewModal {
	t.balanceLessCost = remaining
	return t
}

func (t *ticketReviewModal) TicketPurchased(ticketsPurchased func()) *ticketReviewModal {
	t.ticketsPurchased = ticketsPurchased
	return t
}
