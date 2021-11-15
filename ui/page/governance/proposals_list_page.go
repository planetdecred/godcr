package governance

import (
	"context"
	"sync"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const ProposalsListPageID = "ProposalsList"

type ProposalsListPage struct {
	*load.Load

	ctx        context.Context // page context
	ctxCancel  context.CancelFunc
	proposalMu sync.Mutex

	multiWallet      *dcrlibwallet.MultiWallet
	listContainer    *widget.List
	orderDropDown    *decredmaterial.DropDown
	categoryDropDown *decredmaterial.DropDown
	proposalsList    *decredmaterial.ClickableList
	syncButton       *widget.Clickable
	backButton       decredmaterial.IconButton
	searchEditor     decredmaterial.Editor

	proposalItems []*proposalItem

	syncCompleted bool
}

func NewProposalsPage(l *load.Load) *ProposalsListPage {
	pg := &ProposalsListPage{
		Load:        l,
		multiWallet: l.WL.MultiWallet,
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
	pg.searchEditor = l.Theme.IconEditor(new(widget.Editor), "Search", l.Icons.SearchIcon, true)
	pg.searchEditor.Editor.SingleLine, pg.searchEditor.Editor.Submit, pg.searchEditor.Bordered = true, true, false

	pg.syncButton = new(widget.Clickable)

	pg.proposalsList = pg.Theme.NewClickableList(layout.Vertical)

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	pg.orderDropDown = components.CreateOrderDropDown(l)
	pg.categoryDropDown = l.Theme.DropDown([]decredmaterial.DropDownItem{
		{
			Text: "Under Review",
		},
		{
			Text: "Approved",
		},
		{
			Text: "Rejected",
		},
		{
			Text: "Abandoned",
		},
	}, 1)

	return pg
}

func (pg *ProposalsListPage) ID() string {
	return ProposalsOverviewPageID
}

func (pg *ProposalsListPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.fetchProposals()
}

func (pg *ProposalsListPage) fetchProposals() {
	newestFirst := pg.orderDropDown.SelectedIndex() == 0

	proposalFilter := dcrlibwallet.ProposalCategoryAll
	switch pg.categoryDropDown.SelectedIndex() {
	case 1:
		proposalFilter = dcrlibwallet.ProposalCategoryApproved
	case 2:
		proposalFilter = dcrlibwallet.ProposalCategoryRejected
	case 3:
		proposalFilter = dcrlibwallet.ProposalCategoryAbandoned
	}

	pg.proposalMu.Lock()
	pg.proposalItems = loadProposals(proposalFilter, newestFirst, pg.Load)
	pg.proposalMu.Unlock()
}

func (pg *ProposalsListPage) Handle() {
	if pg.backButton.Button.Clicked() {
		pg.PopFragment()
	}

	for pg.categoryDropDown.Changed() {
		pg.fetchProposals()
	}

	for pg.orderDropDown.Changed() {
		pg.fetchProposals()
	}

	pg.searchEditor.EditorIconButtonEvent = func() {
		//TODO: Proposals search functionality
	}

	if clicked, selectedItem := pg.proposalsList.ItemClicked(); clicked {
		pg.proposalMu.Lock()
		selectedProposal := pg.proposalItems[selectedItem].proposal
		pg.proposalMu.Unlock()

		pg.ChangeFragment(newProposalDetailsPage(pg.Load, &selectedProposal))
	}

	for pg.syncButton.Clicked() {
		go pg.multiWallet.Politeia.Sync()
	}

	if pg.syncCompleted {
		pg.syncCompleted = false
		pg.fetchProposals()
		pg.RefreshWindow()
	}
}

func (pg *ProposalsListPage) OnClose() {
	pg.ctxCancel()
}

// - Layout

func (pg *ProposalsListPage) Layout(gtx C) D {
	return components.UniformPadding(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(pg.backButton.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
							body := func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Flex{}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												txt := pg.Theme.Label(values.TextSize14, "Available Treasury Balance: ")
												txt.Font.Weight = text.SemiBold
												return txt.Layout(gtx)
											}),
											layout.Rigid(func(gtx C) D {
												// Todo get available treasury balance
												return components.LayoutBalanceSize(gtx, pg.Load, "678,678.687654 DCR", values.TextSize14)
											}),
										)
									}),
									layout.Rigid(func(gtx C) D {
										var text string
										if pg.multiWallet.Politeia.IsSyncing() {
											text = "Syncing..."
										} else if pg.syncCompleted {
											text = "Updated"
										} else {
											text = components.TimeAgo(pg.multiWallet.Politeia.GetLastSyncedTimeStamp())
										}

										lastUpdatedInfo := pg.Theme.Label(values.TextSize10, text)
										lastUpdatedInfo.Color = pg.Theme.Color.Gray
										if pg.syncCompleted {
											lastUpdatedInfo.Color = pg.Theme.Color.Success
										}

										return layout.Inset{Top: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
											return lastUpdatedInfo.Layout(gtx)
										})
									}),
								)
							}
							return topNavLayout(gtx, pg.Load, "Proposals", body)
						})
					}),
				)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, pg.layoutContent)
					}),
					layout.Expanded(func(gtx C) D {
						gtx.Constraints.Max.X = gtx.Px(values.MarginPadding150)
						gtx.Constraints.Min.X = gtx.Constraints.Max.X

						card := pg.Theme.Card()
						card.Radius = decredmaterial.Radius(8)
						return card.Layout(gtx, func(gtx C) D {
							return layout.Inset{
								Left:   values.MarginPadding10,
								Right:  values.MarginPadding10,
								Top:    values.MarginPadding2,
								Bottom: values.MarginPadding2,
							}.Layout(gtx, pg.searchEditor.Layout)
						})
					}),
					layout.Expanded(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, func(gtx C) D {
							card := pg.Theme.Card()
							card.Radius = decredmaterial.Radius(8)
							return card.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding8).Layout(gtx, func(gtx C) D {
									return pg.layoutSyncSection(gtx)
								})
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return pg.categoryDropDown.Layout(gtx, 45, true)
					}),
					layout.Expanded(func(gtx C) D {
						return pg.orderDropDown.Layout(gtx, pg.categoryDropDown.Width+39, true)
					}),
				)
			}),
		)
	})
}

func (pg *ProposalsListPage) layoutContent(gtx C) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			pg.proposalMu.Lock()
			proposalItems := pg.proposalItems
			pg.proposalMu.Unlock()

			if len(proposalItems) == 0 {
				return layoutNoProposalsFound(gtx, pg.Load, pg.multiWallet.Politeia.IsSyncing())
			}

			// group 'In discussion' and 'Active' proposals into under review
			listItems := make([]*proposalItem, 0)
			for _, item := range proposalItems {
				if item.proposal.Category == dcrlibwallet.ProposalCategoryPre ||
					item.proposal.Category == dcrlibwallet.ProposalCategoryActive {
					listItems = append(listItems, item)
				}
			}

			prop := proposalItems
			if int32(pg.categoryDropDown.SelectedIndex()+1) == dcrlibwallet.ProposalCategoryAll {
				prop = listItems
			}

			return pg.Theme.List(pg.listContainer).Layout(gtx, 1, func(gtx C, i int) D {
				return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
					return pg.Theme.Card().Layout(gtx, func(gtx C) D {
						return pg.proposalsList.Layout(gtx, len(prop), func(gtx C, i int) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return proposalsList(gtx, pg.Load, prop[i])
								}),
								layout.Rigid(func(gtx C) D {
									return pg.Theme.Separator().Layout(gtx)
								}),
							)
						})
					})
				})
			})
		}),
	)
}

func (pg *ProposalsListPage) layoutSyncSection(gtx C) D {
	if pg.multiWallet.Politeia.IsSyncing() {
		return pg.layoutIsSyncingSection(gtx)
	} else {
		return pg.layoutStartSyncSection(gtx)
	}
}

func (pg *ProposalsListPage) layoutIsSyncingSection(gtx C) D {
	th := material.NewTheme(gofont.Collection())
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding24)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	loader := material.Loader(th)
	loader.Color = pg.Theme.Color.Gray
	return loader.Layout(gtx)
}

func (pg *ProposalsListPage) layoutStartSyncSection(gtx C) D {
	return material.Clickable(gtx, pg.syncButton, func(gtx C) D {
		return pg.Icons.Restore.Layout24dp(gtx)
	})
}

func (pg *ProposalsListPage) listenForSyncNotifications() {
	go func() {
		for {
			var notification interface{}

			select {
			case notification = <-pg.Receiver.NotificationsUpdate:
			case <-pg.ctx.Done():
				return
			}

			switch n := notification.(type) {
			case wallet.Proposal:
				if n.ProposalStatus == wallet.Synced {
					pg.syncCompleted = true
					pg.RefreshWindow()
				}
			}
		}
	}()
}
