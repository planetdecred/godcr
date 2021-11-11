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

	ctx            context.Context // page context
	ctxCancel      context.CancelFunc
	proposalMu     sync.Mutex
	fetchProposals decredmaterial.Button

	multiWallet *dcrlibwallet.MultiWallet

	//categoryList to be removed with new update to UI.
	categoryList   *decredmaterial.ClickableList
	proposalsList  *decredmaterial.ClickableList
	listContainer  *widget.List
	tabCard        decredmaterial.Card
	itemCard       decredmaterial.Card
	syncCard       decredmaterial.Card
	updatedLabel   decredmaterial.Label
	lastSyncedTime string

	proposalItems         []*proposalItem
	proposalCount         []int
	selectedCategoryIndex int

	legendIcon    *decredmaterial.Icon
	infoIcon      *decredmaterial.Icon
	updatedIcon   *decredmaterial.Icon
	syncButton    *widget.Clickable
	startSyncIcon *decredmaterial.Image
	timerIcon     *decredmaterial.Image
	toProposals   decredmaterial.TextAndIconButton

	showSyncedCompleted bool
	isSyncing           bool
}

func NewGovernancePage(l *load.Load) *GovernancePage {
	pg := &GovernancePage{
		Load:        l,
		multiWallet: l.WL.MultiWallet,
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		fetchProposals: l.Theme.Button("Fetch proposals"),
	}

	pg.toProposals = pg.Theme.TextAndIconButton(values.String(values.StrSeeAll), pg.Icons.NavigationArrowForward)
	pg.toProposals.Color = pg.Theme.Color.Primary
	pg.toProposals.BackgroundColor = pg.Theme.Color.Surface

	//categoryList to be removed with new update to UI.
	pg.categoryList = pg.Theme.NewClickableList(layout.Horizontal)
	pg.itemCard = pg.Theme.Card()
	pg.syncButton = new(widget.Clickable)

	pg.infoIcon = decredmaterial.NewIcon(pg.Icons.ActionInfo)
	pg.infoIcon.Color = pg.Theme.Color.Gray

	pg.legendIcon = decredmaterial.NewIcon(pg.Icons.ImageBrightness1)
	pg.legendIcon.Color = pg.Theme.Color.InactiveGray

	pg.updatedIcon = decredmaterial.NewIcon(pg.Icons.NavigationCheck)
	pg.updatedIcon.Color = pg.Theme.Color.Success

	pg.updatedLabel = pg.Theme.Body2("Updated")
	pg.updatedLabel.Color = pg.Theme.Color.Success

	radius := decredmaterial.Radius(0)
	pg.tabCard = pg.Theme.Card()
	pg.tabCard.Radius = radius

	pg.syncCard = pg.Theme.Card()
	pg.syncCard.Radius = radius

	pg.proposalsList = pg.Theme.NewClickableList(layout.Vertical)
	pg.proposalsList.DividerHeight = values.MarginPadding8

	pg.timerIcon = pg.Icons.TimerIcon

	pg.startSyncIcon = pg.Icons.Restore

	return pg
}

func (pg *GovernancePage) ID() string {
	return GovernancePageID
}

func (pg *GovernancePage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	// pg.listenForSyncNotifications()

	proposalItems := loadProposals(dcrlibwallet.ProposalCategoryAll, pg.Load)
	pg.proposalMu.Lock()
	pg.proposalItems = proposalItems
	pg.proposalMu.Unlock()

	pg.isSyncing = pg.multiWallet.Politeia.IsSyncing()
}

func (pg *GovernancePage) topNavLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := pg.Theme.Label(values.TextSize20, GovernancePageID)
					txt.Font.Weight = text.SemiBold
					return txt.Layout(gtx)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						txt := pg.Theme.Label(values.TextSize14, "Available Treasury Balance")
						txt.Font.Weight = text.SemiBold
						return txt.Layout(gtx)
					}),
					layout.Rigid(pg.Theme.Label(values.TextSize14, "636,765 DCR").Layout), // Todo get available treasury balance
				)
			})
		}),
	)
}

func (pg *GovernancePage) Layout(gtx C) D {
	if pg.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
		return components.UniformPadding(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return pg.topNavLayout(gtx)
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
											return layout.Inset{
												Top:    values.MarginPadding2,
												Bottom: values.MarginPadding2,
											}.Layout(gtx, func(gtx C) D {
												return proposalsList(gtx, pg.Load, proposalItems[i])
											})
										})
									}),
									layout.Rigid(func(gtx C) D {
										if i == len(proposalItems) {
											return D{}
										}
										return pg.Theme.Separator().Layout(gtx)
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
		pg.isSyncing = pg.multiWallet.Politeia.IsSyncing()
	}
}

func (pg *GovernancePage) OnClose() {
	pg.ctxCancel()
}

// func (pg *GovernancePage) listenForSyncNotifications() {
// 	go func() {
// 		for {
// 			var notification interface{}

// 			select {
// 			case notification = <-pg.Receiver.NotificationsUpdate:
// 			case <-pg.ctx.Done():
// 				return
// 			}

// 			switch n := notification.(type) {
// 			case wallet.Proposal:
// 				if n.ProposalStatus == wallet.Synced {
// 					pg.isSyncing = false
// 					pg.showSyncedCompleted = true

// 					pg.proposalMu.Lock()
// 					selectedCategory := pg.selectedCategoryIndex
// 					pg.proposalMu.Unlock()
// 					if selectedCategory != -1 {
// 						pg.countProposals()
// 						pg.loadProposals(selectedCategory)
// 					}
// 				}
// 			}
// 		}
// 	}()
// }
