package governance

import (
	"context"
	// "image"
	// "image/color"
	// "strconv"
	// "strings"
	"sync"
	"time"
	// "fmt"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	// "gioui.org/op/clip"
	// "gioui.org/op/paint"
	"gioui.org/text"
	// "gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	// "github.com/planetdecred/godcr/wallet"
)

const ProposalsListPageID = "ProposalsList"

type ProposalsPage struct {
	*load.Load

	ctx        context.Context // page context
	ctxCancel  context.CancelFunc
	proposalMu sync.Mutex

	multiWallet *dcrlibwallet.MultiWallet

	orderDropDown    *decredmaterial.DropDown
	categoryDropDown *decredmaterial.DropDown

	//categoryList to be removed with new update to UI.
	categoryList   *decredmaterial.ClickableList
	proposalsList  *decredmaterial.ClickableList
	listContainer  *widget.List
	tabCard        decredmaterial.Card
	itemCard       decredmaterial.Card
	syncCard       decredmaterial.Card
	updatedLabel   decredmaterial.Label
	lastSyncedTime string

	searchEditor decredmaterial.Editor

	proposalItems         []*proposalItem
	proposalCount         []int
	selectedCategoryIndex int

	legendIcon    *decredmaterial.Icon
	infoIcon      *decredmaterial.Icon
	updatedIcon   *decredmaterial.Icon
	syncButton    *widget.Clickable
	startSyncIcon *decredmaterial.Image
	timerIcon     *decredmaterial.Image

	backButton decredmaterial.IconButton

	showSyncedCompleted bool
	isSyncing           bool
}

func NewProposalsPage(l *load.Load) *ProposalsPage {
	pg := &ProposalsPage{
		Load:                  l,
		multiWallet:           l.WL.MultiWallet,
		selectedCategoryIndex: -1,
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
	pg.searchEditor = l.Theme.IconEditor(new(widget.Editor), "Search", l.Icons.SearchIcon, true)
	pg.searchEditor.Editor.SingleLine, pg.searchEditor.Editor.Submit, pg.searchEditor.Bordered = true, true, false

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

	pg.initLayoutWidgets()
	return pg
}

func (pg *ProposalsPage) ID() string {
	return ProposalsOverviewPageID
}

func (pg *ProposalsPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	// pg.listenForSyncNotifications()
	pg.fetchProposals()

	pg.isSyncing = pg.multiWallet.Politeia.IsSyncing()
}

func (pg *ProposalsPage) fetchProposals() {
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

func (pg *ProposalsPage) Handle() {
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
		pg.isSyncing = true
		go pg.multiWallet.Politeia.Sync()
	}

	if pg.showSyncedCompleted {
		time.AfterFunc(time.Second*3, func() {
			pg.showSyncedCompleted = false
		})
	}
}

// func (pg *ProposalsPage) listenForSyncNotifications() {
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
// 						pg.loadProposals(selectedCategory)
// 					}
// 				}
// 			}
// 		}
// 	}()
// }

func (pg *ProposalsPage) OnClose() {
	pg.ctxCancel()
}

// - Layout

func (pg *ProposalsPage) initLayoutWidgets() {
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
}

func (pg *ProposalsPage) Layout(gtx C) D {
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
										if pg.isSyncing {
											text = "Syncing..."
										} else {
											text = components.TimeAgo(pg.multiWallet.Politeia.GetLastSyncedTimeStamp())
										}

										lastUpdatedInfo := pg.Theme.Label(values.TextSize10, text)
										lastUpdatedInfo.Color = pg.Theme.Color.Gray
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

func (pg *ProposalsPage) layoutContent(gtx C) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			pg.proposalMu.Lock()
			proposalItems := pg.proposalItems
			pg.proposalMu.Unlock()

			if len(proposalItems) == 0 {
				return layoutNoProposalsFound(gtx, pg.Load)
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

func (pg *ProposalsPage) layoutSyncSection(gtx C) D {
	if pg.multiWallet.Politeia.IsSyncing() {
		return pg.layoutIsSyncingSection(gtx)
	}
	return pg.layoutStartSyncSection(gtx)
}

// func (pg *ProposalsPage) layoutProposalsList(gtx C) D {
// 	pg.proposalMu.Lock()
// 	proposalItems := pg.proposalItems
// 	pg.proposalMu.Unlock()
// 	return pg.proposalsList.Layout(gtx, len(proposalItems), func(gtx C, i int) D {
// 		return layout.Inset{
// 			Top:    values.MarginPadding2,
// 			Bottom: values.MarginPadding2,
// 		}.Layout(gtx, func(gtx C) D {
// 			return pg.itemCard.Layout(gtx, func(gtx C) D {
// 				gtx.Constraints.Min.X = gtx.Constraints.Max.X
// 				return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
// 					item := proposalItems[i]
// 					proposal := item.proposal
// 					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 						layout.Rigid(func(gtx C) D {
// 							return pg.layoutAuthorAndDate(gtx, item)
// 						}),
// 						layout.Rigid(func(gtx C) D {
// 							return pg.layoutTitle(gtx, proposal)
// 						}),
// 						layout.Rigid(func(gtx C) D {
// 							if proposal.Category == dcrlibwallet.ProposalCategoryActive ||
// 								proposal.Category == dcrlibwallet.ProposalCategoryApproved ||
// 								proposal.Category == dcrlibwallet.ProposalCategoryRejected {
// 								return pg.layoutProposalVoteBar(gtx, item)
// 							}
// 							return D{}
// 						}),
// 					)
// 				})
// 			})
// 		})
// 	})
// }

func (pg *ProposalsPage) layoutIsSyncingSection(gtx C) D {
	th := material.NewTheme(gofont.Collection())
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding24)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	loader := material.Loader(th)
	loader.Color = pg.Theme.Color.Gray
	return loader.Layout(gtx)
}

func (pg *ProposalsPage) layoutStartSyncSection(gtx C) D {
	return material.Clickable(gtx, pg.syncButton, func(gtx C) D {
		return pg.startSyncIcon.Layout24dp(gtx)
	})
}
