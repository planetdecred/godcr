package components

import (
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type ConsensusItem struct {
	Agenda     dcrlibwallet.Agenda
	VoteButton decredmaterial.Button
}

func AgendaItemWidget(gtx C, l *load.Load, consensusItem *ConsensusItem) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	agenda := consensusItem.Agenda
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layoutAgendaStatus(gtx, l, consensusItem.Agenda)
		}),
		layout.Rigid(layoutAgendaDetails(l, agenda.Description)),
		layout.Rigid(layoutAgendaDetails(l, "ID: #"+agenda.AgendaID)),
		layout.Rigid(layoutAgendaDetails(l, "Voting Preference: "+agenda.VotingPreference)),
		layout.Rigid(func(gtx C) D {
			return layoutAgendaVoteAction(gtx, l, consensusItem)
		}),
	)
}

func layoutAgendaStatus(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {

	var statusLabel decredmaterial.Label
	var statusIcon *decredmaterial.Icon
	var backgroundColor color.NRGBA

	switch agenda.Status() {
	case dcrlibwallet.AgendaStatusFinished:
		statusLabel = l.Theme.Label(values.MarginPadding14, agenda.Status())
		statusLabel.Color = l.Theme.Color.GreenText
		statusIcon = decredmaterial.NewIcon(l.Icons.NavigationCheck)
		statusIcon.Color = l.Theme.Color.Green500
		backgroundColor = l.Theme.Color.Green50
	case dcrlibwallet.AgendaStatusInProgress:
		clr := l.Theme.Color.Primary
		statusLabel = l.Theme.Label(values.MarginPadding14, agenda.Status())
		statusLabel.Color = clr
		statusIcon = decredmaterial.NewIcon(l.Icons.NavMoreIcon)
		statusIcon.Color = clr
		backgroundColor = l.Theme.Color.LightBlue
	case dcrlibwallet.AgendaStatusUpcoming:
		statusLabel = l.Theme.Label(values.MarginPadding14, agenda.Status())
		statusLabel.Color = l.Theme.Color.Text
		statusIcon = decredmaterial.NewIcon(l.Icons.PlayIcon)
		statusIcon.Color = l.Theme.Color.DeepBlue
		backgroundColor = l.Theme.Color.Gray2
	}

	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			lbl := l.Theme.Label(values.MarginPadding20, (strings.Title(strings.ToLower(agenda.AgendaID))))
			lbl.Font.Weight = text.SemiBold
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(lbl.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Background: backgroundColor,
				Width:      decredmaterial.WrapContent,
				Height:     decredmaterial.WrapContent,
				Direction:  layout.Center,
				Alignment:  layout.Middle,
				Border:     decredmaterial.Border{Color: backgroundColor, Width: values.MarginPadding1, Radius: decredmaterial.Radius(10)},
				Padding:    layout.Inset{Top: values.MarginPadding3, Bottom: values.MarginPadding3, Left: values.MarginPadding8, Right: values.MarginPadding8},
				Margin:     layout.Inset{Left: values.MarginPadding10},
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
						return statusIcon.Layout(gtx, values.MarginPadding16)
					})
				}),
				layout.Rigid(statusLabel.Layout))
		}),
	)
}

func layoutAgendaDetails(l *load.Load, data string) layout.Widget {
	return func(gtx C) D {
		lbl := l.Theme.Label(values.MarginPadding16, data)
		lbl.Font.Weight = text.Light
		return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, lbl.Layout)
	}
}

func layoutAgendaVoteAction(gtx C, l *load.Load, item *ConsensusItem) D {
	gtx.Constraints.Min.X, gtx.Constraints.Max.X = gtx.Px(unit.Dp(150)), gtx.Px(unit.Dp(200))
	if item.Agenda.Status() == dcrlibwallet.AgendaStatusInProgress {
		item.VoteButton.Background = l.Theme.Color.Primary
		item.VoteButton.SetEnabled(true)
	} else {
		item.VoteButton.Background = l.Theme.Color.Gray3
		item.VoteButton.SetEnabled(false)
	}
	return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, item.VoteButton.Layout)
}

func LayoutNoAgendasFound(gtx C, l *load.Load, syncing bool) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top:    values.MarginPadding10,
			Bottom: values.MarginPadding10,
		}.Layout(gtx, l.Theme.Body1("No agendas yet").Layout)
	})
}

func LoadAgendas(l *load.Load, selectedWallet *dcrlibwallet.Wallet, newestFirst bool) []*ConsensusItem {
	consensusItems := make([]*ConsensusItem, 0)
	agendas, err := selectedWallet.AllVoteAgendas("", newestFirst)

	if err == nil {
		for i := 0; i < len(agendas); i++ {
			item := &ConsensusItem{
				Agenda:     *agendas[i],
				VoteButton: l.Theme.Button("Update Preference"),
			}
			consensusItems = append(consensusItems, item)
		}
	}
	return consensusItems
}
