package components

import (
	// "fmt"
	"image"
	"image/color"
	"time"

	"fmt"

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
	tooltip      *decredmaterial.Tooltip
	tooltipLabel decredmaterial.Label
	VoteButton decredmaterial.Button
	// voteBar      *VoteBar
}

func AgendasList(gtx C, l *load.Load, consensusItem *ConsensusItem) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		agenda := consensusItem.Agenda
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// layout.Rigid(func(gtx C) D {
			// 	return layoutAuthorAndDate(gtx, l, prop)
			// }),
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
			// layout.Rigid(func(gtx C) D {
			// 	if proposal.Category == dcrlibwallet.ProposalCategoryActive ||
			// 		proposal.Category == dcrlibwallet.ProposalCategoryApproved ||
			// 		proposal.Category == dcrlibwallet.ProposalCategoryRejected {
			// 		return layoutProposalVoteBar(gtx, prop)
			// 	}
			// 	return D{}
			// }),
		)
	})
}

// func layoutAuthorAndDate(gtx C, l *load.Load, item *ProposalItem) D {
// 	proposal := item.Proposal
// 	grayCol := l.Theme.Color.GrayText2

// 	nameLabel := l.Theme.Body2(proposal.Username)
// 	nameLabel.Color = grayCol

// 	dotLabel := l.Theme.H4(" . ")
// 	dotLabel.Color = grayCol

// 	versionLabel := l.Theme.Body2("Version " + proposal.Version)
// 	versionLabel.Color = grayCol

// 	stateLabel := l.Theme.Body2(fmt.Sprintf("%v /2", proposal.VoteStatus))
// 	stateLabel.Color = grayCol

// 	timeAgoLabel := l.Theme.Body2(TimeAgo(proposal.Timestamp))
// 	timeAgoLabel.Color = grayCol

// 	var categoryLabel decredmaterial.Label
// 	var categoryLabelColor color.NRGBA
// 	switch proposal.Category {
// 	case dcrlibwallet.ProposalCategoryApproved:
// 		categoryLabel = l.Theme.Body2("Approved")
// 		categoryLabelColor = l.Theme.Color.Success
// 	case dcrlibwallet.ProposalCategoryActive:
// 		categoryLabel = l.Theme.Body2("Voting")
// 		categoryLabelColor = l.Theme.Color.Primary
// 	case dcrlibwallet.ProposalCategoryRejected:
// 		categoryLabel = l.Theme.Body2("Rejected")
// 		categoryLabelColor = l.Theme.Color.Danger
// 	case dcrlibwallet.ProposalCategoryAbandoned:
// 		categoryLabel = l.Theme.Body2("Abandoned")
// 		categoryLabelColor = grayCol
// 	case dcrlibwallet.ProposalCategoryPre:
// 		categoryLabel = l.Theme.Body2("In discussion")
// 		categoryLabelColor = grayCol
// 	}
// 	categoryLabel.Color = categoryLabelColor

// 	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
// 		layout.Rigid(func(gtx C) D {
// 			return layout.Flex{}.Layout(gtx,
// 				layout.Rigid(nameLabel.Layout),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
// 				}),
// 				layout.Rigid(versionLabel.Layout),
// 			)
// 		}),
// 		layout.Rigid(func(gtx C) D {
// 			return layout.Flex{}.Layout(gtx,
// 				layout.Rigid(categoryLabel.Layout),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Flex{}.Layout(gtx,
// 						layout.Rigid(func(gtx C) D {
// 							if item.Proposal.Category == dcrlibwallet.ProposalCategoryPre {
// 								return layout.Inset{
// 									Right: values.MarginPadding4,
// 								}.Layout(gtx, stateLabel.Layout)
// 							}
// 							return D{}
// 						}),
// 						layout.Rigid(func(gtx C) D {
// 							if item.Proposal.Category == dcrlibwallet.ProposalCategoryActive {
// 								ic := l.Icons.TimerIcon
// 								if l.WL.MultiWallet.ReadBoolConfigValueForKey(load.DarkModeConfigKey, false) {
// 									ic = l.Icons.TimerDarkMode
// 								}
// 								return layout.Inset{
// 									Right: values.MarginPadding4,
// 									Top:   values.MarginPadding3,
// 								}.Layout(gtx, ic.Layout12dp)
// 							}
// 							return D{}
// 						}),
// 						layout.Rigid(timeAgoLabel.Layout),
// 						layout.Rigid(func(gtx C) D {
// 							if item.Proposal.Category == dcrlibwallet.ProposalCategoryPre {
// 								return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
// 									rect := image.Rectangle{
// 										Min: gtx.Constraints.Min,
// 										Max: gtx.Constraints.Max,
// 									}
// 									rect.Max.Y = 20
// 									layoutInfoTooltip(gtx, rect, *item)

// 									infoIcon := decredmaterial.NewIcon(l.Icons.ActionInfo)
// 									infoIcon.Color = l.Theme.Color.GrayText2
// 									return infoIcon.Layout(gtx, values.MarginPadding20)
// 								})
// 							}
// 							return D{}
// 						}),
// 					)
// 				}),
// 			)
// 		}),
// 	)
// }

func layoutAgendaTitle(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H5(agenda.Id)
	lbl.Font.Weight = text.SemiBold

	var categoryLabel decredmaterial.Label
	var categoryLabelColor color.NRGBA

	currentTime := time.Now().Unix()
	println("[][][][]", agenda.StartTime, currentTime, agenda.EndTime)
	if currentTime > agenda.EndTime {
		categoryLabel = l.Theme.Body2("Finished")
		categoryLabelColor = l.Theme.Color.Success
		canVote = false
	} else if currentTime > agenda.StartTime && currentTime < agenda.EndTime {
		categoryLabel = l.Theme.Body2("In progress")
		categoryLabelColor = l.Theme.Color.Primary
		canVote = true
	} else if currentTime > agenda.StartTime {
		categoryLabel = l.Theme.Body2("Upcoming")
		categoryLabelColor = l.Theme.Color.Black
		canVote = false
	}

	categoryLabel.Color = categoryLabelColor
	// return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(lbl.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(categoryLabel.Layout),
			)
		}),
	)
}

func layoutAgendaDescription(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H6(agenda.Description)
		lbl.Font.Weight = text.Light
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func layoutAgendaID(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H6("ID: #" + agenda.Id)
	lbl.Font.Weight = text.Light
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func layoutAgendaVotingPreference(gtx C, l *load.Load, agenda dcrlibwallet.Agenda) D {
	lbl := l.Theme.H6("Voting Preference: " + "Abstain")
	lbl.Font.Weight = text.Light
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func layoutAgendaVoteAction(gtx C, l *load.Load, item *ConsensusItem) D {
	gtx.Constraints.Min.X, gtx.Constraints.Max.X = 150, 150
	// var voteButton decredmaterial.Button
	item.VoteButton = l.Theme.Button("Change Vote")
	if canVote {
		item.VoteButton.Background = l.Theme.Color.Primary
	} else {
		item.VoteButton.Background = l.Theme.Color.Gray3
	}
	return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
		return item.VoteButton.Layout(gtx)
	})
}

// func layoutProposalVoteBar(gtx C, item *ProposalItem) D {
// 	proposal := item.Proposal
// 	yes := float32(proposal.YesVotes)
// 	no := float32(proposal.NoVotes)
// 	quorumPercent := float32(proposal.QuorumPercentage)
// 	passPercentage := float32(proposal.PassPercentage)
// 	eligibleTickets := float32(proposal.EligibleTickets)

// 	return item.voteBar.
// 		SetYesNoVoteParams(yes, no).
// 		SetVoteValidityParams(eligibleTickets, quorumPercent, passPercentage).
// 		SetProposalDetails(proposal.NumComments, proposal.PublishedAt, proposal.Token).
// 		Layout(gtx)
// }

func layoutAgendaInfoTooltip(gtx C, rect image.Rectangle, item ConsensusItem) {
	inset := layout.Inset{Top: values.MarginPadding20, Left: values.MarginPaddingMinus195}
	item.tooltip.Layout(gtx, rect, inset, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Px(values.MarginPadding195)
		gtx.Constraints.Max.X = gtx.Px(values.MarginPadding195)
		return item.tooltipLabel.Layout(gtx)
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

func LoadAgendas(l *load.Load) []*ConsensusItem {
	consensusItems := make([]*ConsensusItem, 0)

	agendasResponse, err := l.WL.MultiWallet.GetAllAgendas(1)
	fmt.Println("[][][] agendas", agendasResponse)
	fmt.Println("[][][] error", err)
	if err == nil {
		fmt.Println("[][][] length of agendas", len(agendasResponse.Agendas))
		for i := 0; i < len(agendasResponse.Agendas); i++ {
			// agenda := agendas[i]
			item := &ConsensusItem{
				Agenda: *agendasResponse.Agendas[i],
			}
			// agenda := &dcrlibwallet.Agenda {
			// 	Id: agendasResponse.Agendas[i].Id,

			// }

			// if proposal.Category == dcrlibwallet.ProposalCategoryPre {
			// 	tooltipLabel := l.Theme.Caption("")
			// 	tooltipLabel.Color = l.Theme.Color.GrayText2
			// 	if proposal.VoteStatus == 1 {
			// 		tooltipLabel.Text = "Waiting for author to authorize voting"
			// 	} else if proposal.VoteStatus == 2 {
			// 		tooltipLabel.Text = "Waiting for admin to trigger the start of voting"
			// 	}

			// 	item.tooltip = l.Theme.Tooltip()
			// 	item.tooltipLabel = tooltipLabel
			// }

			consensusItems = append(consensusItems, item)
		}
	}
	return consensusItems
}