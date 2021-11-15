package governance

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type proposalItem struct {
	proposal     dcrlibwallet.Proposal
	tooltip      *decredmaterial.Tooltip
	tooltipLabel decredmaterial.Label
	voteBar      *VoteBar
}

func topNavLayout(gtx C, l *load.Load, title string, content layout.Widget) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := l.Theme.Label(values.TextSize20, title)
					txt.Font.Weight = text.SemiBold
					return txt.Layout(gtx)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, content)
		}),
	)
}

func proposalsList(gtx C, l *load.Load, prop *proposalItem) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		proposal := prop.proposal
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layoutAuthorAndDate(gtx, l, prop)
			}),
			layout.Rigid(func(gtx C) D {
				return layoutTitle(gtx, l, proposal)
			}),
			layout.Rigid(func(gtx C) D {
				if proposal.Category == dcrlibwallet.ProposalCategoryActive ||
					proposal.Category == dcrlibwallet.ProposalCategoryApproved ||
					proposal.Category == dcrlibwallet.ProposalCategoryRejected {
					return layoutProposalVoteBar(gtx, prop)
				}
				return D{}
			}),
		)
	})
}

func layoutAuthorAndDate(gtx C, l *load.Load, item *proposalItem) D {
	proposal := item.proposal
	grayCol := l.Theme.Color.Gray

	nameLabel := l.Theme.Body2(proposal.Username)
	nameLabel.Color = grayCol

	dotLabel := l.Theme.H4(" . ")
	dotLabel.Color = grayCol

	versionLabel := l.Theme.Body2("Version " + proposal.Version)
	versionLabel.Color = grayCol

	stateLabel := l.Theme.Body2(fmt.Sprintf("%v /2", proposal.VoteStatus))
	stateLabel.Color = grayCol

	timeAgoLabel := l.Theme.Body2(components.TimeAgo(proposal.Timestamp))
	timeAgoLabel.Color = grayCol

	var categoryLabel decredmaterial.Label
	var categoryLabelColor color.NRGBA
	switch proposal.Category {
	case dcrlibwallet.ProposalCategoryApproved:
		categoryLabel = l.Theme.Body2("Approved")
		categoryLabelColor = l.Theme.Color.Success
	case dcrlibwallet.ProposalCategoryActive:
		categoryLabel = l.Theme.Body2("Voting")
		categoryLabelColor = l.Theme.Color.Primary
	case dcrlibwallet.ProposalCategoryRejected:
		categoryLabel = l.Theme.Body2("Rejected")
		categoryLabelColor = l.Theme.Color.Danger
	case dcrlibwallet.ProposalCategoryAbandoned:
		categoryLabel = l.Theme.Body2("Abandoned")
		categoryLabelColor = grayCol
	case dcrlibwallet.ProposalCategoryPre:
		categoryLabel = l.Theme.Body2("In discussion")
		categoryLabelColor = grayCol
	}
	categoryLabel.Color = categoryLabelColor

	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(nameLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(versionLabel.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(categoryLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							if item.proposal.Category == dcrlibwallet.ProposalCategoryPre {
								return layout.Inset{
									Right: values.MarginPadding4,
								}.Layout(gtx, stateLabel.Layout)
							}
							return D{}
						}),
						layout.Rigid(func(gtx C) D {
							if item.proposal.Category == dcrlibwallet.ProposalCategoryActive {
								return layout.Inset{
									Right: values.MarginPadding4,
									Top:   values.MarginPadding3,
								}.Layout(gtx, l.Icons.TimerIcon.Layout12dp)
							}
							return D{}
						}),
						layout.Rigid(timeAgoLabel.Layout),
						layout.Rigid(func(gtx C) D {
							if item.proposal.Category == dcrlibwallet.ProposalCategoryPre {
								return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									rect := image.Rectangle{
										Min: gtx.Constraints.Min,
										Max: gtx.Constraints.Max,
									}
									rect.Max.Y = 20
									layoutInfoTooltip(gtx, rect, *item)

									infoIcon := decredmaterial.NewIcon(l.Icons.ActionInfo)
									infoIcon.Color = l.Theme.Color.Gray
									return infoIcon.Layout(gtx, values.MarginPadding20)
								})
							}
							return D{}
						}),
					)
				}),
			)
		}),
	)
}

func layoutTitle(gtx C, l *load.Load, proposal dcrlibwallet.Proposal) D {
	lbl := l.Theme.H6(proposal.Name)
	lbl.Font.Weight = text.SemiBold
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func layoutProposalVoteBar(gtx C, item *proposalItem) D {
	proposal := item.proposal
	yes := float32(proposal.YesVotes)
	no := float32(proposal.NoVotes)
	quorumPercent := float32(proposal.QuorumPercentage)
	passPercentage := float32(proposal.PassPercentage)
	eligibleTickets := float32(proposal.EligibleTickets)

	return item.voteBar.
		SetYesNoVoteParams(yes, no).
		SetVoteValidityParams(eligibleTickets, quorumPercent, passPercentage).
		SetProposalDetails(proposal.NumComments, proposal.PublishedAt, proposal.Token).
		Layout(gtx)
}

func layoutInfoTooltip(gtx C, rect image.Rectangle, item proposalItem) {
	inset := layout.Inset{Top: values.MarginPadding20, Left: values.MarginPaddingMinus195}
	item.tooltip.Layout(gtx, rect, inset, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Px(values.MarginPadding195)
		gtx.Constraints.Max.X = gtx.Px(values.MarginPadding195)
		return item.tooltipLabel.Layout(gtx)
	})
}

func layoutNoProposalsFound(gtx C, l *load.Load, syncing bool) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	text := l.Theme.Body1("No proposals yet")
	if syncing {
		text = l.Theme.Body1("Fetching proposals...")
	}
	return layout.Center.Layout(gtx, text.Layout)
}

func loadProposals(category int32, newestFirst bool, l *load.Load) []*proposalItem {
	proposalItems := make([]*proposalItem, 0)

	proposals, err := l.WL.MultiWallet.Politeia.GetProposalsRaw(category, 0, 0, newestFirst)
	if err == nil {
		for i := 0; i < len(proposals); i++ {
			proposal := proposals[i]
			item := &proposalItem{
				proposal: proposals[i],
				voteBar:  NewVoteBar(l),
			}

			if proposal.Category == dcrlibwallet.ProposalCategoryPre {
				tooltipLabel := l.Theme.Caption("")
				tooltipLabel.Color = l.Theme.Color.Gray
				if proposal.VoteStatus == 1 {
					tooltipLabel.Text = "Waiting for author to authorize voting"
				} else if proposal.VoteStatus == 2 {
					tooltipLabel.Text = "Waiting for admin to trigger the start of voting"
				}

				item.tooltip = l.Theme.Tooltip()
				item.tooltipLabel = tooltipLabel
			}

			proposalItems = append(proposalItems, item)
		}
	}
	return proposalItems
}
