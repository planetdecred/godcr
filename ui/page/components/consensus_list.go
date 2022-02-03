package components

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

var canVote bool

type ConsensusItem struct {
	Agenda     dcrlibwallet.Agenda
	VoteButton decredmaterial.Button
}

func AgendasList(gtx C, l *load.Load, consensusItem *ConsensusItem) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		agenda := consensusItem.Agenda
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layoutAgendaStatus(gtx, l, consensusItem.Agenda)
			}),
			layout.Rigid(func(gtx C) D {
				return layoutAgendaDescription(gtx, l, agenda)
			}),
			layout.Rigid(func(gtx C) D {
				return layoutAgendaID(gtx, l, agenda)
			}),
			layout.Rigid(func(gtx C) D {
				return layoutAgendaVotingPreference(gtx, l, agenda)
			}),
			layout.Rigid(func(gtx C) D {
				return layoutAgendaVoteAction(gtx, l, consensusItem)
			}),
		)
	})
}

func layoutAgendaStatus(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H5(agenda.AgendaID)
	lbl.Font.Weight = text.SemiBold

	var statusLabel decredmaterial.Label
	var statusLabelColor color.NRGBA
	var statusIcon *decredmaterial.Icon
	var backgroundColor color.NRGBA

	switch agenda.Status {
	case "Finished":
		statusLabel = l.Theme.Label(values.MarginPadding14, agenda.Status)
		statusLabelColor = l.Theme.Color.Success
		statusIcon = decredmaterial.NewIcon(l.Icons.NavigationCheck)
		statusIcon.Color = statusLabelColor
		backgroundColor = l.Theme.Color.Green50
		canVote = false
	case "In progress":
		statusLabel = l.Theme.Label(values.MarginPadding14, agenda.Status)
		statusLabelColor = l.Theme.Color.Primary
		statusIcon = decredmaterial.NewIcon(l.Icons.NavMoreIcon)
		statusIcon.Color = statusLabelColor
		backgroundColor = l.Theme.Color.LightBlue
		canVote = true
	case "Upcoming":
		statusLabel = l.Theme.Label(values.MarginPadding14, agenda.Status)
		statusLabelColor = l.Theme.Color.DeepBlue
		statusIcon = decredmaterial.NewIcon(l.Icons.PlayIcon)
		statusIcon.Color = statusLabelColor
		backgroundColor = l.Theme.Color.Gray2
		canVote = false
	}

	statusLabel.Color = statusLabelColor
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(lbl.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Background: backgroundColor,
				Width:     decredmaterial.WrapContent,
				Height:    decredmaterial.WrapContent,
				Direction: layout.Center,
				Alignment: layout.Middle,
				Border:    decredmaterial.Border{Color: backgroundColor, Width: values.MarginPadding1, Radius: decredmaterial.Radius(10)},
				Padding:   layout.Inset{Top: values.MarginPadding3, Bottom: values.MarginPadding3, Left: values.MarginPadding8, Right: values.MarginPadding8},
				Margin:    layout.Inset{Left: values.MarginPadding10},
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
						return statusIcon.Layout(gtx, values.MarginPadding16)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return statusLabel.Layout(gtx)
				}))
		}),
	)
}

func layoutAgendaDescription(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H6(agenda.Description)
	lbl.Font.Weight = text.Light
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func layoutAgendaID(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H6("ID: #" + agenda.AgendaID)
	lbl.Font.Weight = text.Light
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func layoutAgendaVotingPreference(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H6("Voting Preference: " + agenda.VotingPreference)
	lbl.Font.Weight = text.Light
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func layoutAgendaVoteAction(gtx C, l *load.Load, item *ConsensusItem) D {
	gtx.Constraints.Min.X, gtx.Constraints.Max.X = 150, 150
	if canVote {
		item.VoteButton.Background = l.Theme.Color.Primary
		item.VoteButton.SetEnabled(true)
	} else {
		item.VoteButton.Background = l.Theme.Color.Gray3
		item.VoteButton.SetEnabled(false)
	}
	return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
		return item.VoteButton.Layout(gtx)
	})
}

func LayoutNoAgendasFound(gtx C, l *load.Load, syncing bool) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	text := l.Theme.Body1("No agendas yet")

	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top:    values.MarginPadding10,
			Bottom: values.MarginPadding10,
		}.Layout(gtx, text.Layout)
	})
}

func LoadAgendas(l *load.Load, selectedWallet *dcrlibwallet.Wallet, newestFirst bool) []*ConsensusItem {
	consensusItems := make([]*ConsensusItem, 0)
	agendasResponse, err := l.WL.MultiWallet.Consensus.GetAllAgendasForWallet(selectedWallet.ID, newestFirst)

	if err == nil {
		for i := 0; i < len(agendasResponse.Agendas); i++ {
			item := &ConsensusItem{
				Agenda:     *agendasResponse.Agendas[i],
				VoteButton: l.Theme.Button("Change Vote"),
			}
			consensusItems = append(consensusItems, item)
		}
	}
	return consensusItems
}
