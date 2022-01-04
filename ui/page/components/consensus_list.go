package components

import (
	"image/color"
	"time"

	// "fmt"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

var canVote bool

type ConsensusItem struct {
	Agenda       dcrlibwallet.Agenda
	VoteButton   decredmaterial.Button
}

func AgendasList(gtx C, l *load.Load, consensusItem *ConsensusItem) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		agenda := consensusItem.Agenda
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layoutAgendaTitle(gtx, l, consensusItem.Agenda)
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

func layoutAgendaTitle(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H5(agenda.AgendaID)
	lbl.Font.Weight = text.SemiBold

	var categoryLabel decredmaterial.Label
	var categoryLabelColor color.NRGBA
	var categoryIcon *decredmaterial.Icon

	currentTime := time.Now().Unix()
	// println("[][][][]", agenda.StartTime, currentTime, agenda.EndTime)
	if currentTime > agenda.ExpireTime {
		categoryLabel = l.Theme.Label(values.MarginPadding14, "Finished")
		categoryLabelColor = l.Theme.Color.Success
		categoryIcon = decredmaterial.NewIcon(l.Icons.NavigationCheck)
		categoryIcon.Color = categoryLabelColor
		canVote = false
	} else if currentTime > agenda.StartTime && currentTime < agenda.ExpireTime {
		categoryLabel = l.Theme.Label(values.MarginPadding14, "In progress")
		categoryLabelColor = l.Theme.Color.Primary
		categoryIcon = decredmaterial.NewIcon(l.Icons.NavMoreIcon)
		categoryIcon.Color = categoryLabelColor
		canVote = true
	} else if currentTime > agenda.StartTime {
		categoryLabel = l.Theme.Label(values.MarginPadding14, "Upcoming")
		categoryLabelColor = l.Theme.Color.Black
		categoryIcon = decredmaterial.NewIcon(l.Icons.PlayIcon)
		categoryIcon.Color = categoryLabelColor
		canVote = false
	}

	categoryLabel.Color = categoryLabelColor
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(lbl.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			// return layout.Flex{}.Layout(gtx,
			// 	layout.Rigid(categoryLabel.Layout),
			// )
			return decredmaterial.LinearLayout{
				Width:     decredmaterial.WrapContent,
				Height:    decredmaterial.WrapContent,
				Direction: layout.Center,
				Alignment: layout.Middle,
				Border:    decredmaterial.Border{Color: categoryLabelColor, Width: values.MarginPadding1, Radius: decredmaterial.Radius(10)},
				Padding:   layout.Inset{Top: values.MarginPadding3, Bottom: values.MarginPadding3, Left: values.MarginPadding8, Right: values.MarginPadding8},
				Margin:    layout.Inset{Left: values.MarginPadding10},
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
						return categoryIcon.Layout(gtx, values.MarginPadding16)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return categoryLabel.Layout(gtx)
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
	if syncing {
		text = l.Theme.Body1("Fetching agendas...")
	}
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top:    values.MarginPadding10,
			Bottom: values.MarginPadding10,
		}.Layout(gtx, text.Layout)
	})
}

func LoadAgendas(l *load.Load, selectedWallet *dcrlibwallet.Wallet, newestFirst bool) []*ConsensusItem {
	consensusItems := make([]*ConsensusItem, 0)
	// agendasResponse, err := selectedWallet.GetAllAgendas()
	// l.WL.MultiWallet.Consensus.ClearSavedVoteChoices()
	// l.WL.MultiWallet.Consensus.ClearSavedAgendas()
	// _, err := l.WL.MultiWallet.Consensus.GetAllAgendas(selectedWallet.ID)
	
	agendas, err := l.WL.MultiWallet.Consensus.GetAgendasByWalletIDRaw(selectedWallet.ID, 0, 0, newestFirst)

	// fmt.Println("[][][] agendas", agendas)
	// fmt.Println("[][][] error", err)
	if err == nil {
		// fmt.Println("[][][] length of agendas", len(agendas))
		for i := 0; i < len(agendas); i++ {
			item := &ConsensusItem{
				Agenda:     agendas[i],
				VoteButton: l.Theme.Button("Change Vote"),
			}
			consensusItems = append(consensusItems, item)
		}
	}
	return consensusItems
}
