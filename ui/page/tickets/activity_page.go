package tickets

import (
	"sort"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const ActivityPageID = "TicketsActivity"

type ActivityPage struct {
	*load.Load
	tickets      **wallet.Tickets
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
		tickets:     l.WL.Tickets,
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
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "Ticket activity",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				walletID := pg.wallets[pg.walletDropDown.SelectedIndex()].ID
				tickets := (*pg.tickets).Confirmed[walletID]
				return layout.Stack{Alignment: layout.N}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Card().Layout(gtx, func(gtx C) D {
								gtx.Constraints.Min = gtx.Constraints.Max
								if pg.ticketTypeDropDown.SelectedIndex()-1 != -1 {
									tickets = filterTickets(tickets, func(ticketStatus string) bool {
										return ticketStatus == strings.ToUpper(pg.ticketTypeDropDown.Selected())
									})
								}

								if len(tickets) == 0 {
									txt := pg.Theme.Body1("No tickets yet")
									txt.Color = pg.Theme.Color.Gray2
									txt.Alignment = text.Middle
									return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D { return txt.Layout(gtx) })
								}
								return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
									return pg.ticketsList.Layout(gtx, len(tickets), func(gtx C, index int) D {
										return ticketActivityRow(gtx, pg.Load, tickets[index], index)
									})
								})
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return pg.walletDropDown.Layout(gtx, 0, false)
					}),
					layout.Expanded(func(gtx C) D {
						return pg.orderDropDown.Layout(gtx, 0, true)
					}),
					layout.Expanded(func(gtx C) D {
						return pg.ticketTypeDropDown.Layout(gtx, pg.orderDropDown.Width+10, true)
					}),
				)
			},
		}
		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func filterTickets(tickets []wallet.Ticket, f func(string) bool) []wallet.Ticket {
	t := make([]wallet.Ticket, 0)
	for _, v := range tickets {
		if f(v.Info.Status) {
			t = append(t, v)
		}
	}
	return t
}

func (pg *ActivityPage) Handle() {

	sortSelection := pg.orderDropDown.SelectedIndex()
	if pg.filterSorter != sortSelection {
		pg.filterSorter = sortSelection
		newestFirst := pg.filterSorter == 0
		for _, wal := range pg.wallets {
			tickets := (*pg.tickets).Confirmed[wal.ID]
			sort.SliceStable(tickets, func(i, j int) bool {
				backTime := time.Unix(tickets[j].Info.Ticket.Timestamp, 0)
				frontTime := time.Unix(tickets[i].Info.Ticket.Timestamp, 0)
				if newestFirst {
					return backTime.Before(frontTime)
				}
				return frontTime.Before(backTime)
			})
		}
	}

}

func (pg *ActivityPage) OnClose() {}
