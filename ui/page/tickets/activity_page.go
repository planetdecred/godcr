package tickets

import (
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/wallet"
	"github.com/planetdecred/godcr/ui/values"
)

const ActivityPageID = "TicketsActivity"

type ActivityPage struct {
	*load.Load
	tickets      []load.Ticket
	ticketsList  layout.List
	filterSorter int

	orderDropDown      *decredmaterial.DropDown
	ticketTypeDropDown *decredmaterial.DropDown
	walletDropDown     *decredmaterial.DropDown

	wallets []*dcrlibwallet.Wallet

	backButton decredmaterial.IconButton
}

func newTicketActivityPage(l *load.Load) *ActivityPage {
	pg := &ActivityPage{
		Load:        l,
		ticketsList: layout.List{Axis: layout.Vertical},
	}
	pg.orderDropDown = components.CreateOrderDropDown(l)
	pg.ticketTypeDropDown = pg.Theme.DropDown([]decredmaterial.DropDownItem{
		{Text: "All"},
		{Text: "Unmined"},
		{Text: "Immature"},
		{Text: "Live"},
		{Text: "Voted"},
		{Text: "Missed"},
		{Text: "Expired"},
		{Text: "Revoked"},
	}, 1)

	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)

	return pg
}

func (pg *ActivityPage) ID() string {
	return ActivityPageID
}

func (pg *ActivityPage) OnResume() {
	pg.wallets = pg.WL.SortedWalletList()
	components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets)
}

func (pg *ActivityPage) Layout(gtx layout.Context) layout.Dimensions {
	components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets)
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "Ticket activity",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return layout.Stack{Alignment: layout.N}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Card().Layout(gtx, func(gtx C) D {
								gtx.Constraints.Min = gtx.Constraints.Max
								if len(pg.tickets) == 0 {
									txt := pg.Theme.Body1("No tickets yet")
									txt.Color = pg.Theme.Color.Gray2
									txt.Alignment = text.Middle
									return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D { return txt.Layout(gtx) })
								}
								return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
									return pg.ticketsList.Layout(gtx, len(pg.tickets), func(gtx C, index int) D {
										return ticketActivityRow(gtx, pg.Load, pg.tickets[index], index)
									})
								})
							})
						})
					}),
					layout.Stacked(func(gtx C) D {
						return pg.dropDowns(gtx)
					}),
				)
			},
		}
		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *ActivityPage) dropDowns(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.walletDropDown.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left: values.MarginPadding5,
					}.Layout(gtx, func(gtx C) D {
						return pg.ticketTypeDropDown.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left: values.MarginPadding5,
					}.Layout(gtx, func(gtx C) D {
						return pg.orderDropDown.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *ActivityPage) Handle() {

}

func (pg *ActivityPage) OnClose() {}
