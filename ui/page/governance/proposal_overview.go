package governance

import (
	"context"
	"sync"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const GovernancePageID = "Governance"

type GovernancePage struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	fetchProposals decredmaterial.Button
	toProposals    decredmaterial.TextAndIconButton
	proposalsList  *decredmaterial.ClickableList
	listContainer  *widget.List

	proposalItems []*proposalItem
	proposalMu    sync.Mutex

	lastSyncedTime      string
	showSyncedCompleted bool
	isSyncing           bool
}

func NewGovernancePage(l *load.Load) *GovernancePage {
	pg := &GovernancePage{
		Load: l,
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		fetchProposals: l.Theme.Button("Fetch proposals"),
	}

	pg.toProposals = pg.Theme.TextAndIconButton(values.String(values.StrSeeAll), pg.Icons.NavigationArrowForward)
	pg.toProposals.Color = pg.Theme.Color.Primary
	pg.toProposals.BackgroundColor = pg.Theme.Color.Surface

	pg.proposalsList = pg.Theme.NewClickableList(layout.Vertical)

	return pg
}

func (pg *GovernancePage) ID() string {
	return GovernancePageID
}

func (pg *GovernancePage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	// pg.listenForSyncNotifications()

	proposalItems := loadProposals(dcrlibwallet.ProposalCategoryAll, true, pg.Load)
	pg.proposalMu.Lock()
	pg.proposalItems = proposalItems
	pg.proposalMu.Unlock()

	pg.isSyncing = pg.WL.MultiWallet.Politeia.IsSyncing()
}

func (pg *GovernancePage) Layout(gtx C) D {
	if pg.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
		return components.UniformPadding(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						body := func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									txt := pg.Theme.Label(values.TextSize14, "Available Treasury Balance")
									txt.Font.Weight = text.SemiBold
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									// Todo get available treasury balance
									return components.LayoutBalanceSize(gtx, pg.Load, "678,678.687654 DCR", values.TextSize14)
								}),
							)
						}
						return topNavLayout(gtx, pg.Load, "Governance", body)
					})
				}),
				layout.Rigid(func(gtx C) D {
					pg.proposalMu.Lock()
					proposalItems := pg.proposalItems
					pg.proposalMu.Unlock()
					if len(proposalItems) == 0 {
						return layoutNoProposalsFound(gtx, pg.Load)
					}

					return pg.Theme.List(pg.listContainer).Layout(gtx, 1, func(gtx C, i int) D {
						return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Card().Layout(gtx, func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
											return layout.Flex{}.Layout(gtx,
												layout.Rigid(pg.Theme.Label(values.TextSize14, "Recent Proposals").Layout),
												layout.Flexed(1, func(gtx C) D {
													return layout.E.Layout(gtx, pg.toProposals.Layout)
												}),
											)
										})
									}),
									layout.Rigid(pg.Theme.Separator().Layout),
									layout.Rigid(func(gtx C) D {
										return pg.proposalsList.Layout(gtx, len(proposalItems), func(gtx C, i int) D {
											return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													return proposalsList(gtx, pg.Load, proposalItems[i])
												}),
												layout.Rigid(pg.Theme.Separator().Layout),
											)
										})
									}),
								)
							})
						})
					})
				}),
			)
		})
	}
	return components.UniformPadding(gtx, pg.splashScreenLayout)
}

func (pg *GovernancePage) Handle() {
	for pg.fetchProposals.Clicked() {
		go pg.WL.MultiWallet.Politeia.Sync()
		pg.WL.Wallet.SaveConfigValueForKey(load.FetchProposalConfigKey, true)
		pg.isSyncing = pg.WL.MultiWallet.Politeia.IsSyncing()
	}

	for pg.toProposals.Button.Clicked() {
		pg.ChangeFragment(NewProposalsPage(pg.Load))
	}

	if clicked, selectedItem := pg.proposalsList.ItemClicked(); clicked {
		pg.proposalMu.Lock()
		selectedProposal := pg.proposalItems[selectedItem].proposal
		pg.proposalMu.Unlock()

		pg.ChangeFragment(newProposalDetailsPage(pg.Load, &selectedProposal))
	}
}

func (pg *GovernancePage) OnClose() {
	pg.ctxCancel()
}
