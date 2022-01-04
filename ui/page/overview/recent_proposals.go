package overview

import (
	"time"

	"gioui.org/layout"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *AppOverviewPage) initializeProposalsWidget() *AppOverviewPage {
	pg.toProposals = pg.Theme.TextAndIconButton(values.String(values.StrSeeAll), pg.Icons.NavigationArrowForward)
	pg.toProposals.Color = pg.Theme.Color.Primary
	pg.toProposals.BackgroundColor = pg.Theme.Color.Surface

	pg.proposalsList = pg.Theme.NewClickableList(layout.Vertical)
	return pg
}

func (pg *AppOverviewPage) loadRecentProposals() {
	proposalItems := components.LoadProposals(dcrlibwallet.ProposalCategoryAll, true, pg.Load)

	// get only proposals within the last week
	listItems := make([]*components.ProposalItem, 0)
	for _, item := range proposalItems {
		utcTime := time.Unix(item.Proposal.Timestamp, 0).UTC()
		if time.Now().UTC().Sub(utcTime).Hours() <= 24*8 {
			listItems = append(listItems, item)
		}
	}

	pg.proposalMu.Lock()
	pg.proposalItems = listItems
	pg.proposalMu.Unlock()
}

func (pg *AppOverviewPage) recentProposalsSection(gtx C) D {
	pg.proposalMu.Lock()
	proposalItems := pg.proposalItems
	pg.proposalMu.Unlock()

	return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								title := pg.Theme.Label(values.TextSize14, "Recent Proposals")
								title.Color = pg.Theme.Color.GrayText1
								return title.Layout(gtx)
							}),
							layout.Flexed(1, func(gtx C) D {
								if len(proposalItems) == 0 {
									return D{}
								}
								return layout.E.Layout(gtx, pg.toProposals.Layout)
							}),
						)
					})
				}),
				layout.Rigid(pg.Theme.Separator().Layout),
				layout.Rigid(func(gtx C) D {
					if len(proposalItems) == 0 {
						return components.LayoutNoProposalsFound(gtx, pg.Load, pg.WL.MultiWallet.Politeia.IsSyncing(), dcrlibwallet.ProposalCategoryAll)
					}
					return pg.proposalsList.Layout(gtx, len(proposalItems), func(gtx C, i int) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return components.ProposalsList(gtx, pg.Load, proposalItems[i])
							}),
							layout.Rigid(pg.Theme.Separator().Layout),
						)
					})
				}),
			)
		})
	})
}
