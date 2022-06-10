package governance

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type ticketSelector struct {
	*load.Load

	dialogTitle string

	changed         bool
	showTicketModal *decredmaterial.Clickable
	selectedTicket  *dcrlibwallet.Transaction

	liveTickets []*dcrlibwallet.Transaction
}

func newTicketSelector(l *load.Load, lv []*dcrlibwallet.Transaction) *ticketSelector {
	ts := &ticketSelector{
		Load:            l,
		showTicketModal: l.Theme.NewClickable(true),
		liveTickets:     lv,
	}
	return ts
}

func (ts *ticketSelector) Title(title string) *ticketSelector {
	ts.dialogTitle = title
	return ts
}

func (ts *ticketSelector) Changed() bool {
	changed := ts.changed
	ts.changed = false
	return changed
}

func (ts *ticketSelector) SelectTicket(ticketHash string) {
	for _, liveTicket := range ts.liveTickets {
		if liveTicket.Hash == ticketHash {
			ts.changed = true
			ts.selectedTicket = liveTicket
			break
		}
	}
}

func (ts *ticketSelector) SelectedTicket() *dcrlibwallet.Transaction {
	return ts.selectedTicket
}

func (ts *ticketSelector) handle(window app.WindowNavigator) {
	if ts.showTicketModal.Clicked() {
		ticketSelectorModal := newTicketSelectorModal(ts.Load, ts.liveTickets).
			title(values.String(values.StrSelectTicket)).
			ticketSelected(func(ticket *dcrlibwallet.Transaction) {
				ts.SelectTicket(ticket.Hash)
			})
		window.ShowModal(ticketSelectorModal)
	}
}

func (ts *ticketSelector) Layout(gtx layout.Context, window app.WindowNavigator) layout.Dimensions {
	ts.handle(window)

	border := widget.Border{
		Color:        ts.Theme.Color.Gray2,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}

	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return ts.showTicketModal.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if ts.selectedTicket == nil {
							txt := ts.Theme.Label(values.TextSize16, values.String(values.StrSelectTicket))
							txt.Color = ts.Theme.Color.GrayText3
							return txt.Layout(gtx)
						}
						return ts.Theme.Label(values.TextSize16, ts.selectedTicket.Hash).Layout(gtx)
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							inset := layout.Inset{
								Left: values.MarginPadding15,
							}
							return inset.Layout(gtx, func(gtx C) D {
								ic := decredmaterial.NewIcon(ts.Theme.Icons.DropDownIcon)
								ic.Color = ts.Theme.Color.Gray1
								return ic.Layout(gtx, values.MarginPadding20)
							})
						})
					}),
				)
			})
		})
	})
}

type ticketSelectorModal struct {
	*load.Load
	*decredmaterial.Modal

	dialogTitle string

	liveTickets    []*dcrlibwallet.Transaction
	selectedTicket *dcrlibwallet.Transaction
	ticketList     *decredmaterial.ClickableList

	ticketSelectedCallback func(*dcrlibwallet.Transaction)
}

func newTicketSelectorModal(l *load.Load, lv []*dcrlibwallet.Transaction) *ticketSelectorModal {
	tsm := &ticketSelectorModal{
		Load:  l,
		Modal: l.Theme.ModalFloatTitle("TicketSelectorModal"),

		liveTickets: lv,
		ticketList:  l.Theme.NewClickableList(layout.Vertical),
	}

	return tsm
}

func (tsm *ticketSelectorModal) OnResume() {}

func (tsm *ticketSelectorModal) Handle() {
	if tsm.Modal.BackdropClicked(true) {
		tsm.Dismiss()
	}

	if clicked, selectedItem := tsm.ticketList.ItemClicked(); clicked {
		tsm.selectedTicket = tsm.liveTickets[selectedItem]
		tsm.ticketSelectedCallback(tsm.liveTickets[selectedItem])
		tsm.Dismiss()
	}
}

func (tsm *ticketSelectorModal) title(title string) *ticketSelectorModal {
	tsm.dialogTitle = title
	return tsm
}

func (tsm *ticketSelectorModal) ticketSelected(callback func(*dcrlibwallet.Transaction)) *ticketSelectorModal {
	tsm.ticketSelectedCallback = callback
	tsm.Dismiss()
	return tsm
}

func (tsm *ticketSelectorModal) OnDismiss() {}

func (tsm *ticketSelectorModal) Layout(gtx layout.Context) layout.Dimensions {
	return tsm.Modal.Layout(gtx, []layout.Widget{
		func(gtx C) D {
			title := tsm.Theme.Label(values.TextSize20, tsm.dialogTitle)
			title.Font.Weight = text.SemiBold
			return title.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := tsm.Theme.Label(values.TextSize14, values.String(values.StrHash))
					txt.Color = tsm.Theme.Color.GrayText2
					return txt.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					listTickets := tsm.liveTickets
					return tsm.ticketList.Layout(gtx, len(listTickets), func(gtx C, i int) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Flexed(0.8, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding12, Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
									return tsm.Theme.Label(values.TextSize16, listTickets[i].Hash).Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								if tsm.selectedTicket != nil || tsm.selectedTicket != listTickets[i] {
									return layout.Dimensions{}
								}
								ic := decredmaterial.NewIcon(tsm.Theme.Icons.NavigationCheck)
								return ic.Layout(gtx, values.MarginPadding20)
							}),
						)
					})
				}),
			)
		},
	})
}
