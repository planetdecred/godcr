package ui

import (
	"image"
	"image/color"
	"strings"
	"time"
	//"fmt"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/ararog/timeago"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageProposals = "Proposals"
const (
	categoryStateFetching = iota
	categoryStateFetched
	categoryStateError
)

type proposalItem struct {
	btn      *widget.Clickable
	proposal dcrlibwallet.Proposal
	voteBar  decredmaterial.VoteBar
}

type tab struct {
	title     string
	btn       *widget.Clickable
	category  int32
	state     int
	proposals []proposalItem
}

type tabs struct {
	tabs     []tab
	selected int
}

type proposalsPage struct {
	theme                     *decredmaterial.Theme
	wallet                    *wallet.Wallet
	proposalsList             *layout.List
	scrollContainer           *decredmaterial.ScrollContainer
	tabs                      tabs
	tabCard                   decredmaterial.Card
	itemCard                  decredmaterial.Card
	notify                    func(string, bool)
	hasFetchedInitialProposal bool
	legendIcon                *widget.Icon
	infoIcon                  *widget.Icon
}

var (
	proposalCategoryTitles = []string{"In discussion", "Voting", "Approved", "Rejected", "Abandoned"}
	proposalCategories     = []int32{
		dcrlibwallet.ProposalCategoryPre,
		dcrlibwallet.ProposalCategoryActive,
		dcrlibwallet.ProposalCategoryApproved,
		dcrlibwallet.ProposalCategoryRejected,
		dcrlibwallet.ProposalCategoryAbandoned,
	}
)

func (win *Window) ProposalsPage(common pageCommon) layout.Widget {
	pg := &proposalsPage{
		theme:           common.theme,
		wallet:          win.wallet,
		proposalsList:   &layout.List{Axis: layout.Vertical},
		scrollContainer: common.theme.ScrollContainer(),
		tabCard:         common.theme.Card(),
		itemCard:        common.theme.Card(),
		notify:          common.Notify,
		legendIcon:      common.icons.imageBrightness1,
		infoIcon:        common.icons.actionInfo,
	}
	pg.tabCard.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}

	for i := range proposalCategoryTitles {
		pg.tabs.tabs = append(pg.tabs.tabs,
			tab{
				title:    proposalCategoryTitles[i],
				btn:      new(widget.Clickable),
				category: proposalCategories[i],
				state:    categoryStateFetching,
			},
		)
	}

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *proposalsPage) handle(common pageCommon) {
	for i := range pg.tabs.tabs {
		if pg.tabs.tabs[i].btn.Clicked() {
			pg.tabs.selected = i
			pg.fetchProposalsForCategory()
		}

		for k := range pg.tabs.tabs[i].proposals {
			if pg.tabs.tabs[i].proposals[k].btn.Clicked() {
				// TODO goto proposal details page
			}
		}
	}
}

func (pg *proposalsPage) onfetchSuccess(proposals []dcrlibwallet.Proposal) {
	pg.tabs.tabs[pg.tabs.selected].proposals = make([]proposalItem, len(proposals))
	for i := range proposals {
		pg.tabs.tabs[pg.tabs.selected].proposals[i] = proposalItem{
			btn:      new(widget.Clickable),
			proposal: proposals[i],
			voteBar:  pg.theme.VoteBar(pg.infoIcon, pg.legendIcon),
		}
	}
	pg.tabs.tabs[pg.tabs.selected].state = categoryStateFetched
	if !pg.hasFetchedInitialProposal {
		pg.hasFetchedInitialProposal = true
	}
}

func (pg *proposalsPage) onFetchError(err error) {
	pg.tabs.tabs[pg.tabs.selected].state = categoryStateError
	if !pg.hasFetchedInitialProposal {
		pg.hasFetchedInitialProposal = true
	}
	pg.notify(err.Error(), false)
}

func (pg *proposalsPage) fetchProposalsForCategory() {
	selected := pg.tabs.tabs[pg.tabs.selected]
	pg.wallet.GetProposals(selected.category, pg.onfetchSuccess, pg.onFetchError)
}

func (pg *proposalsPage) layoutTabs(gtx C) D {
	width := float32(gtx.Constraints.Max.X-20) / float32(len(pg.tabs.tabs))

	return pg.tabCard.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Left:  values.MarginPadding10,
			Right: values.MarginPadding10,
		}.Layout(gtx, func(gtx C) D {
			return (&layout.List{}).Layout(gtx, len(pg.tabs.tabs), func(gtx C, i int) D {
				gtx.Constraints.Min.X = int(width)
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						return material.Clickable(gtx, pg.tabs.tabs[i].btn, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
								return layout.Center.Layout(gtx, func(gtx C) D {
									lbl := pg.theme.Body1(pg.tabs.tabs[i].title)
									if pg.tabs.selected == i {
										lbl.Color = pg.theme.Color.Primary
									}
									return lbl.Layout(gtx)
								})
							})
						})
					}),
					layout.Stacked(func(gtx C) D {
						if pg.tabs.selected != i {
							return D{}
						}
						tabHeight := gtx.Px(unit.Dp(2))
						tabRect := image.Rect(0, 0, int(width), tabHeight)
						paint.FillShape(gtx.Ops, pg.theme.Color.Primary, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: int(width), Y: tabHeight},
						}
					}),
				)
			})
		})
	})
}

func (pg *proposalsPage) layoutFetchingState(gtx C) D {
	str := "Fetching " + strings.ToLower(proposalCategoryTitles[pg.tabs.selected]) + " proposals..."

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, func(gtx C) D {
		return pg.theme.Body1(str).Layout(gtx)
	})
}

func (pg *proposalsPage) layoutNoProposalsFound(gtx C) D {
	str := "No " + strings.ToLower(proposalCategoryTitles[pg.tabs.selected]) + " proposals"

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, func(gtx C) D {
		return pg.theme.Body1(str).Layout(gtx)
	})
}

func (pg *proposalsPage) layoutAuthorAndDate(gtx C, proposal dcrlibwallet.Proposal) D {
	nameLabel := pg.theme.Body2(proposal.Username)
	dotLabel := pg.theme.H4(" . ")
	versionLabel := pg.theme.Body2("Version " + proposal.Version)

	timeAgoLabel := pg.theme.Body2(timeAgo(proposal.Timestamp))

	var categoryLabel decredmaterial.Label
	var categoryLabelColor color.NRGBA
	switch proposal.Category {
	case dcrlibwallet.ProposalCategoryApproved:
		categoryLabel = pg.theme.Body2("Approved")
		categoryLabelColor = pg.theme.Color.Success
	case dcrlibwallet.ProposalCategoryActive:
		categoryLabel = pg.theme.Body2("Voting")
		categoryLabelColor = pg.theme.Color.Primary
	case dcrlibwallet.ProposalCategoryRejected:
		categoryLabel = pg.theme.Body2("Rejected")
		categoryLabelColor = pg.theme.Color.Danger
	case dcrlibwallet.ProposalCategoryAbandoned:
		categoryLabel = pg.theme.Body2("Abandoned")
		categoryLabelColor = pg.theme.Color.Gray
	case dcrlibwallet.ProposalCategoryPre:
		categoryLabel = pg.theme.Body2("in discussion")
		categoryLabelColor = pg.theme.Color.Gray
	}
	categoryLabel.Color = categoryLabelColor

	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(nameLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: unit.Dp(-23)}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(versionLabel.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(categoryLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: unit.Dp(-23)}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(timeAgoLabel.Layout),
			)
		}),
	)
}

func (pg *proposalsPage) layoutTitle(gtx C, proposal dcrlibwallet.Proposal) D {
	lbl := pg.theme.H6(proposal.Name)
	lbl.Color = pg.theme.Color.Text

	return layout.Inset{
		Top:    values.MarginPadding5,
		Bottom: values.MarginPadding5,
	}.Layout(gtx, lbl.Layout)
}

func (pg *proposalsPage) layoutProposalVoteBar(gtx C, proposalItem proposalItem) D {
	yes := float32(proposalItem.proposal.YesVotes)
	no := float32(proposalItem.proposal.NoVotes)
	quorumPercent := float32(proposalItem.proposal.QuorumPercentage)
	passPercentage := float32(proposalItem.proposal.PassPercentage)
	eligibleTickets := float32(proposalItem.proposal.EligibleTickets)

	return proposalItem.voteBar.SetParams(yes, no, eligibleTickets, quorumPercent, passPercentage).LayoutWithLegend(gtx)
}

func (pg *proposalsPage) layoutProposalsList(gtx C) D {
	selected := pg.tabs.tabs[pg.tabs.selected]
	wdgs := make([]func(gtx C) D, len(selected.proposals))
	for i := range selected.proposals {
		index := i
		proposalItem := selected.proposals[index]
		wdgs[index] = func(gtx C) D {
			return layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
				Left:   values.MarginPadding15,
				Right:  values.MarginPadding15,
			}.Layout(gtx, func(gtx C) D {
				return material.Clickable(gtx, selected.proposals[index].btn, func(gtx C) D {
					return pg.itemCard.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return pg.layoutAuthorAndDate(gtx, proposalItem.proposal)
								}),
								layout.Rigid(func(gtx C) D {
									return pg.layoutTitle(gtx, proposalItem.proposal)
								}),
								layout.Rigid(func(gtx C) D {
									if proposalItem.proposal.Category == dcrlibwallet.ProposalCategoryActive ||
										proposalItem.proposal.Category == dcrlibwallet.ProposalCategoryApproved ||
										proposalItem.proposal.Category == dcrlibwallet.ProposalCategoryRejected {
										return pg.layoutProposalVoteBar(gtx, proposalItem)
									}
									return D{}
								}),
							)
						})
					})
				})
			})
		}
	}
	return pg.scrollContainer.Layout(gtx, wdgs)
}

func (pg *proposalsPage) layoutContent(gtx C) D {
	selected := pg.tabs.tabs[pg.tabs.selected]
	if selected.state == categoryStateFetching {
		return pg.layoutFetchingState(gtx)
	} else if selected.state == categoryStateFetched && len(selected.proposals) == 0 {
		return pg.layoutNoProposalsFound(gtx)
	}
	return pg.layoutProposalsList(gtx)
}

func (pg *proposalsPage) Layout(gtx C, common pageCommon) D {
	if !pg.hasFetchedInitialProposal {
		pg.fetchProposalsForCategory()
	}

	return common.LayoutWithoutPadding(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.layoutTabs),
			layout.Flexed(1, pg.layoutContent),
		)
	})
}

func timeAgo(timestamp int64) string {
	timeAgo, _ := timeago.TimeAgoWithTime(time.Now(), time.Unix(timestamp, 0))
	return timeAgo
}
