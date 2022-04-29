package components

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type ProposalItem struct {
	Proposal     dcrlibwallet.Proposal
	tooltip      *decredmaterial.Tooltip
	tooltipLabel decredmaterial.Label
	voteBar      *VoteBar
}

func ProposalsList(gtx C, theme *decredmaterial.Theme, prop *ProposalItem) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		proposal := prop.Proposal
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layoutAuthorAndDate(gtx, theme, prop)
			}),
			layout.Rigid(func(gtx C) D {
				return layoutTitle(gtx, theme, proposal)
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

func layoutAuthorAndDate(gtx C, theme *decredmaterial.Theme, item *ProposalItem) D {
	proposal := item.Proposal
	grayCol := theme.Color.GrayText2

	nameLabel := theme.Body2(proposal.Username)
	nameLabel.Color = grayCol

	dotLabel := theme.H4(" . ")
	dotLabel.Color = grayCol

	versionLabel := theme.Body2("Version " + proposal.Version)
	versionLabel.Color = grayCol

	stateLabel := theme.Body2(fmt.Sprintf("%v /2", proposal.VoteStatus))
	stateLabel.Color = grayCol

	timeAgoLabel := theme.Body2(TimeAgo(proposal.Timestamp))
	timeAgoLabel.Color = grayCol

	var categoryLabel decredmaterial.Label
	var categoryLabelColor color.NRGBA
	switch proposal.Category {
	case dcrlibwallet.ProposalCategoryApproved:
		categoryLabel = theme.Body2("Approved")
		categoryLabelColor = theme.Color.Success
	case dcrlibwallet.ProposalCategoryActive:
		categoryLabel = theme.Body2("Voting")
		categoryLabelColor = theme.Color.Primary
	case dcrlibwallet.ProposalCategoryRejected:
		categoryLabel = theme.Body2("Rejected")
		categoryLabelColor = theme.Color.Danger
	case dcrlibwallet.ProposalCategoryAbandoned:
		categoryLabel = theme.Body2("Abandoned")
		categoryLabelColor = grayCol
	case dcrlibwallet.ProposalCategoryPre:
		categoryLabel = theme.Body2("In discussion")
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
							if item.Proposal.Category == dcrlibwallet.ProposalCategoryPre {
								return layout.Inset{
									Right: values.MarginPadding4,
								}.Layout(gtx, stateLabel.Layout)
							}
							return D{}
						}),
						layout.Rigid(func(gtx C) D {
							if item.Proposal.Category == dcrlibwallet.ProposalCategoryActive {
								return layout.Inset{
									Right: values.MarginPadding4,
									Top:   values.MarginPadding3,
								}.Layout(gtx, theme.Icons.TimerIcon.Layout12dp)
							}
							return D{}
						}),
						layout.Rigid(timeAgoLabel.Layout),
						layout.Rigid(func(gtx C) D {
							if item.Proposal.Category == dcrlibwallet.ProposalCategoryPre {
								return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									rect := image.Rectangle{
										Min: gtx.Constraints.Min,
										Max: gtx.Constraints.Max,
									}
									rect.Max.Y = 20
									layoutInfoTooltip(gtx, rect, *item)

									infoIcon := decredmaterial.NewIcon(theme.Icons.ActionInfo)
									infoIcon.Color = theme.Color.GrayText2
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

func layoutTitle(gtx C, theme *decredmaterial.Theme, proposal dcrlibwallet.Proposal) D {
	lbl := theme.H6(proposal.Name)
	lbl.Font.Weight = text.SemiBold
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func layoutProposalVoteBar(gtx C, item *ProposalItem) D {
	proposal := item.Proposal
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

func layoutInfoTooltip(gtx C, rect image.Rectangle, item ProposalItem) {
	inset := layout.Inset{Top: values.MarginPadding20, Left: values.MarginPaddingMinus195}
	item.tooltip.Layout(gtx, rect, inset, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Px(values.MarginPadding195)
		gtx.Constraints.Max.X = gtx.Px(values.MarginPadding195)
		return item.tooltipLabel.Layout(gtx)
	})
}

func LayoutNoProposalsFound(gtx C, theme *decredmaterial.Theme, syncing bool, category int32) D {
	var selectedCategory string
	switch category {
	case dcrlibwallet.ProposalCategoryApproved:
		selectedCategory = "approved"
	case dcrlibwallet.ProposalCategoryRejected:
		selectedCategory = "rejected"
	case dcrlibwallet.ProposalCategoryAbandoned:
		selectedCategory = "abandoned"
	default:
		selectedCategory = "under review"
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	text := theme.Body1(fmt.Sprintf("No proposals %s ", selectedCategory))
	text.Color = theme.Color.GrayText3
	if syncing {
		text = theme.Body1("Fetching proposals...")
	}

	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top:    values.MarginPadding10,
			Bottom: values.MarginPadding10,
		}.Layout(gtx, text.Layout)
	})
}

func LoadProposals(category int32, newestFirst bool, mw *dcrlibwallet.MultiWallet, theme *decredmaterial.Theme) []*ProposalItem {
	proposalItems := make([]*ProposalItem, 0)

	proposals, err := mw.Politeia.GetProposalsRaw(category, 0, 0, newestFirst)
	if err == nil {
		for i := 0; i < len(proposals); i++ {
			proposal := proposals[i]
			item := &ProposalItem{
				Proposal: proposals[i],
				voteBar:  NewVoteBar(theme),
			}

			if proposal.Category == dcrlibwallet.ProposalCategoryPre {
				tooltipLabel := theme.Caption("")
				tooltipLabel.Color = theme.Color.GrayText2
				if proposal.VoteStatus == 1 {
					tooltipLabel.Text = "Waiting for author to authorize voting"
				} else if proposal.VoteStatus == 2 {
					tooltipLabel.Text = "Waiting for admin to trigger the start of voting"
				}

				item.tooltip = theme.Tooltip()
				item.tooltipLabel = tooltipLabel
			}

			proposalItems = append(proposalItems, item)
		}
	}
	return proposalItems
}
