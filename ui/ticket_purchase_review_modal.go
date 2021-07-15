package ui

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const TicketReviewModalID = "ticket_review_modal"

type ticketReviewModal struct {
	*pageCommon

	randomID        string
	totalCost       int64
	vspHost         string
	ticketCount     int64
	balanceLessCost int64

	modal            decredmaterial.Modal
	spendingPassword decredmaterial.Editor
	purchase         decredmaterial.Button
	cancelPurchase   decredmaterial.Button

	account *dcrlibwallet.Account
}

func newTicketReviewModal(common *pageCommon) *ticketReviewModal {
	m := &ticketReviewModal{
		pageCommon:       common,
		randomID:         fmt.Sprintf("%s-%d", TicketReviewModalID, generateRandomNumber()),
		modal:            *common.theme.ModalFloatTitle(),
		spendingPassword: common.theme.EditorPassword(new(widget.Editor), "Spending password"),
		purchase:         common.theme.Button(new(widget.Clickable), "Purchase ticket"),
		cancelPurchase:   common.theme.Button(new(widget.Clickable), "Cancel"),
	}

	m.cancelPurchase.Background, m.cancelPurchase.Color = color.NRGBA{}, common.theme.Color.Primary
	return m
}

func (t *ticketReviewModal) Layout(gtx layout.Context) layout.Dimensions {
	l := []layout.Widget{
		func(gtx C) D {
			return t.theme.Label(values.TextSize20, "Confirm to purchase tickets").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					tleft := t.theme.Label(values.TextSize14, "Amount")
					tleft.Color = t.theme.Color.Gray2
					tright := t.theme.Label(values.TextSize14, fmt.Sprintf("%d", t.ticketCount))
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := t.theme.Label(values.TextSize14, "Total cost")
					tleft.Color = t.theme.Color.Gray2
					tright := t.theme.Label(values.TextSize14, dcrutil.Amount(t.totalCost).String())
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, t.theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := t.theme.Label(values.TextSize14, "Account")
					tleft.Color = t.theme.Color.Gray2
					selectedWallet := t.multiWallet.WalletWithID(t.account.WalletID)
					tright := t.theme.Label(values.TextSize14, selectedWallet.Name)
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := t.theme.Label(values.TextSize14, "Remaining")
					tleft.Color = t.theme.Color.Gray2
					tright := t.theme.Label(values.TextSize14, dcrutil.Amount(t.balanceLessCost).String())
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, t.theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := t.theme.Label(values.TextSize14, "VSP")
					tleft.Color = t.theme.Color.Gray2
					tright := t.theme.Label(values.TextSize14, t.vspHost)
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
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

func (t *ticketReviewModal) modalID() string {
	return t.randomID
}
func (t *ticketReviewModal) OnDismiss() {}

func (t *ticketReviewModal) Show() {
	t.showModal(t)
}

func (t *ticketReviewModal) Dismiss() {
	t.dismissModal(t)
}

func (t *ticketReviewModal) OnResume() {}

func (t *ticketReviewModal) handle() {
	for t.cancelPurchase.Button.Clicked() {
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
