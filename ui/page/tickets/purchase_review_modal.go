package tickets

import (
	"fmt"
	"image/color"

	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const reviewModalID = "ticket_review_modal"

type ticketReviewModal struct {
	*load.Load

	totalCost       int64
	vspHost         string
	ticketCount     int64
	balanceLessCost int64

	modal            decredmaterial.Modal
	spendingPassword decredmaterial.Editor
	purchase         decredmaterial.Button
	cancelPurchase   decredmaterial.Button
	purchaseTickets  func(password []byte)

	account *dcrlibwallet.Account
}

func newTicketReviewModal(l *load.Load) *ticketReviewModal {
	m := &ticketReviewModal{
		Load:             l,
		modal:            *l.Theme.ModalFloatTitle(),
		spendingPassword: l.Theme.EditorPassword(new(widget.Editor), "Spending password"),
		purchase:         l.Theme.Button(new(widget.Clickable), "Purchase ticket"),
		cancelPurchase:   l.Theme.Button(new(widget.Clickable), "Cancel"),
	}

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
					tright := t.Theme.Label(values.TextSize14, t.vspHost)
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
						return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, t.cancelPurchase.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if t.ticketCount > 1 {
							t.purchase.Text = fmt.Sprintf("Purchase %d tickets", t.ticketCount)
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
	for t.cancelPurchase.Button.Clicked() {
		t.Dismiss()
	}

	for t.purchase.Button.Clicked() {
		t.purchaseTickets([]byte(t.spendingPassword.Editor.Text()))
		t.Dismiss()
	}
}

func (t *ticketReviewModal) VSPHost(host string) *ticketReviewModal {
	t.vspHost = host
	return t
}

func (t *ticketReviewModal) Account(account *dcrlibwallet.Account) *ticketReviewModal {
	t.account = account
	return t
}

func (t *ticketReviewModal) TicketCount(tickets int64) *ticketReviewModal {
	t.ticketCount = tickets
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

func (t *ticketReviewModal) TicketPurchase(purchaseTickets func([]byte)) *ticketReviewModal {
	t.purchaseTickets = purchaseTickets
	return t
}
