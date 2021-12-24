package governance

import (
	"context"
	// "fmt"
	"image"
	// "image/color"
	// "strconv"
	// "strings"
	// "sync"
	"time"

	// "gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	// "gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	// "github.com/planetdecred/godcr/wallet"
)

const GovernancePageID = "Governance"

type GovernancePage struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	multiWallet *dcrlibwallet.MultiWallet

	//categoryList to be removed with new update to UI.
	categoryList  *decredmaterial.ClickableList
	listContainer *widget.List
	tabCard       decredmaterial.Card
	itemCard      decredmaterial.Card

	selectedCategoryIndex int

	proposalsPage *ProposalsPage
	consensusPage *ConsensusPage

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton
	// legendIcon    *decredmaterial.Icon
	// infoIcon      *decredmaterial.Icon
	// updatedIcon   *decredmaterial.Icon
	// syncButton    *widget.Clickable
	// startSyncIcon *decredmaterial.Image
	// timerIcon     *decredmaterial.Image

	// showSyncedCompleted bool
	// isSyncing           bool
}

var (
	proposalCategoryTitles = []string{"Proposals", "Consensus Changes"}
	proposalCategories     = []int32{
		dcrlibwallet.ProposalCategoryPre,
		dcrlibwallet.ProposalCategoryActive,
		dcrlibwallet.ProposalCategoryApproved,
		dcrlibwallet.ProposalCategoryRejected,
		dcrlibwallet.ProposalCategoryAbandoned,
	}
)

func NewGovernancePage(l *load.Load) *GovernancePage {
	pg := &GovernancePage{
		Load:                  l,
		multiWallet:           l.WL.MultiWallet,
		selectedCategoryIndex: -1,
		proposalsPage: NewProposalsPage(l),
		consensusPage: NewConsensusPage(l),
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)
	// pg.consensusPage.consensusItems[0].VoteButton = l.Theme.Button("Change Vote")

	pg.initLayoutWidgets()

	return pg
}

func (pg *GovernancePage) OnResume() {
	selectedCategory := pg.selectedCategoryIndex

 	if selectedCategory == -1 {
		 pg.selectedCategoryIndex = 0
 	}

	 /** begin proposal page resume method */

	 pg.proposalsPage.ctx, pg.proposalsPage.ctxCancel = context.WithCancel(context.TODO())
	 pg.proposalsPage.listenForSyncNotifications()
	 pg.proposalsPage.fetchProposals()
	 pg.proposalsPage.isSyncing = pg.proposalsPage.multiWallet.Politeia.IsSyncing()

	 /** begin proposal page resume method */

	 /** begin consensus page resume method */

	 pg.consensusPage.fetchAgendas()

	 /** end consensus page resume method */
}

func (pg *GovernancePage) OnClose() {
	// pg.ctxCancel()
}


func (pg *GovernancePage) initLayoutWidgets() {
	//categoryList to be removed with new update to UI.
	pg.categoryList = pg.Theme.NewClickableList(layout.Horizontal)
	pg.itemCard = pg.Theme.Card()

	radius := decredmaterial.Radius(0)
	pg.tabCard = pg.Theme.Card()
	pg.tabCard.Radius = radius
}

func (pg *GovernancePage) ID() string {
	return GovernancePageID
}

func (pg *GovernancePage) Handle() {
	for pg.backButton.Button.Clicked() {
		pg.PopFragment()
	}

	//categoryList to be removed with new update to UI.
	if clicked, selectedItem := pg.categoryList.ItemClicked(); clicked {
		pg.selectedCategoryIndex = selectedItem
		// go pg.switchTab(gtx, selectedItem)÷\
	}

	/** begin proposal page handles */
	for pg.proposalsPage.fetchProposalsBtn.Clicked() {
		go pg.WL.MultiWallet.Politeia.Sync()
		pg.proposalsPage.isSyncing = pg.proposalsPage.multiWallet.Politeia.IsSyncing()
		pg.WL.Wallet.SaveConfigValueForKey(load.FetchProposalConfigKey, true)
	}

	for pg.infoButton.Button.Clicked() {
		pg.proposalsPage.showInfoModal()
	}

	for pg.proposalsPage.categoryDropDown.Changed() {
		pg.proposalsPage.fetchProposals()
	}

	for pg.proposalsPage.orderDropDown.Changed() {
		pg.proposalsPage.fetchProposals()
	}

	pg.proposalsPage.searchEditor.EditorIconButtonEvent = func() {
		//TODO: Proposals search functionality
	}

	if clicked, selectedItem := pg.proposalsPage.proposalsList.ItemClicked(); clicked {
		pg.proposalsPage.proposalMu.Lock()
		selectedProposal := pg.proposalsPage.proposalItems[selectedItem].Proposal
		pg.proposalsPage.proposalMu.Unlock()

		pg.ChangeFragment(NewProposalDetailsPage(pg.Load, &selectedProposal))
	}

	for pg.proposalsPage.syncButton.Clicked() {
		go pg.multiWallet.Politeia.Sync()
		pg.proposalsPage.isSyncing = true

		//Todo: check after 1min if sync does not start, set isSyncing to false and cancel sync
	}

	if pg.proposalsPage.syncCompleted {
		time.AfterFunc(time.Second*3, func() {
			pg.proposalsPage.syncCompleted = false
			pg.RefreshWindow()
		})
	}

	decredmaterial.DisplayOneDropdown(pg.proposalsPage.orderDropDown, pg.proposalsPage.categoryDropDown)
	
	/** end proposal page handles */

	/** begin consensus page handles */
	
	for pg.consensusPage.orderDropDown.Changed() {
		pg.consensusPage.fetchAgendas()
	}

	// for i := range pg.consensusPage.consensusItems {
	// 	for pg.consensusPage.consensusItems[i].VoteButton.Clicked() {
	// 		newAgendaVoteModal(pg.Load).Show()
	// 	}
	// }


	// for pg.consensusPage.consensusItems[0].VoteButton.Clicked() {
	// 	newAgendaVoteModal(pg.Load).Show()
	// }

	/** end consensus page handles */

}

func (pg *GovernancePage) Layout(gtx C) D {
	// border := widget.Border{Color: pg.Theme.Color.Primary, CornerRadius: values.MarginPadding0, Width: values.MarginPadding1}
	// borderLayout := func(gtx layout.Context, body layout.Widget) layout.Dimensions {
	// 	return border.Layout(gtx, body)
	// }

	// return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
	// 	layout.Rigid(func(gtx C) D {
	// 		return layout.Flex{}.Layout(gtx,
	// 			layout.Flexed(1, func(gtx C) D {
	// 				return borderLayout(gtx, pg.layoutTabs)
	// 			}),
	// 		)
	// 	}),
	// )
	if pg.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
		return components.UniformPadding(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(pg.backButton.Layout),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
								body := func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.End}.Layout(gtx,
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
									)
								}

								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												txt := pg.Theme.Label(values.TextSize20, "Governance")
												txt.Font.Weight = text.SemiBold
												return txt.Layout(gtx)
											}),
										)
									}),
									layout.Flexed(1, func(gtx C) D {
										return layout.E.Layout(gtx, body)
									}),
								)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.layoutTabs(gtx)
						}),
					)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return layout.Stack{}.Layout(gtx,
						layout.Expanded(func(gtx C) D {
							// return pg.proposalsPage.Layout(gtx)
							return pg.switchTab(gtx, pg.selectedCategoryIndex)
							// return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, pg.proposalsPage.layoutContent)
						}),
						// layout.Expanded(func(gtx C) D {
						// 	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding150)
						// 	gtx.Constraints.Min.X = gtx.Constraints.Max.X

						// 	card := pg.Theme.Card()
						// 	card.Radius = decredmaterial.Radius(8)
						// 	return card.Layout(gtx, func(gtx C) D {
						// 		return layout.Inset{
						// 			Left:   values.MarginPadding10,
						// 			Right:  values.MarginPadding10,
						// 			Top:    values.MarginPadding2,
						// 			Bottom: values.MarginPadding2,
						// 		}.Layout(gtx, pg.proposalsPage.searchEditor.Layout)
						// 	})
						// }),
						// layout.Expanded(func(gtx C) D {
						// 	gtx.Constraints.Min.X = gtx.Constraints.Max.X
						// 	return layout.E.Layout(gtx, func(gtx C) D {
						// 		card := pg.Theme.Card()
						// 		card.Radius = decredmaterial.Radius(8)
						// 		return card.Layout(gtx, func(gtx C) D {
						// 			return layout.UniformInset(values.MarginPadding8).Layout(gtx, func(gtx C) D {
						// 				return pg.proposalsPage.layoutSyncSection(gtx)
						// 			})
						// 		})
						// 	})
						// }),
						// layout.Expanded(func(gtx C) D {
						// 	return pg.proposalsPage.orderDropDown.Layout(gtx, 45, true)
						// }),
						// layout.Expanded(func(gtx C) D {
						// 	return pg.proposalsPage.categoryDropDown.Layout(gtx, pg.proposalsPage.orderDropDown.Width+41, true)
						// }),
					)
					})
				}),
			)
		})
	}
	return D{}
}

func (pg *GovernancePage) switchTab(gtx C, selectedCategoryIndex int) D {
	if selectedCategoryIndex == 0 {
		return pg.proposalsPage.Layout(gtx)
	}

	return pg.consensusPage.Layout(gtx)
}

func (pg *GovernancePage) layoutTabs(gtx C) D {
	width := float32(gtx.Constraints.Max.X-20) / float32(len(proposalCategoryTitles))
	// pg.proposalMu.Lock()
	selectedCategory := pg.selectedCategoryIndex
	// pg.proposalMu.Unlock()

	return pg.tabCard.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Left:  values.MarginPadding12,
			Right: values.MarginPadding12,
		}.Layout(gtx, func(gtx C) D {
			// categoryList to be removed with new update to UI.
			return pg.categoryList.Layout(gtx, len(proposalCategoryTitles), func(gtx C, i int) D {
				gtx.Constraints.Min.X = int(width)
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding14).Layout(gtx, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										lbl := pg.Theme.Body1(proposalCategoryTitles[i])
										lbl.Color = pg.Theme.Color.Gray3
										if selectedCategory == i {
											lbl.Color = pg.Theme.Color.Primary
										}
										return lbl.Layout(gtx)
									}),
								)
							})
						})
					}),
					layout.Stacked(func(gtx C) D {
						if selectedCategory != i {
							return D{}
						}
						tabHeight := gtx.Px(unit.Dp(2))
						tabRect := image.Rect(0, 0, int(width), tabHeight)
						paint.FillShape(gtx.Ops, pg.Theme.Color.Primary, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: int(width), Y: tabHeight},
						}
					}),
				)
			})
		})
	})
}

// func (pg *ProposalsPage) layoutContent(gtx C) D {
// 	return layout.Stack{}.Layout(gtx,
// 		layout.Expanded(func(gtx C) D {
// 			pg.proposalMu.Lock()
// 			proposalItems := pg.proposalItems
// 			pg.proposalMu.Unlock()

// 			return pg.Theme.List(pg.listContainer).Layout(gtx, 1, func(gtx C, i int) D {
// 				return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
// 					return pg.Theme.Card().Layout(gtx, func(gtx C) D {
// 						if len(proposalItems) == 0 {
// 							return components.LayoutNoProposalsFound(gtx, pg.Load, pg.isSyncing, int32(pg.categoryDropDown.SelectedIndex()))
// 						}
// 						return pg.proposalsList.Layout(gtx, len(proposalItems), func(gtx C, i int) D {
// 							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 								layout.Rigid(func(gtx C) D {
// 									return components.ProposalsList(gtx, pg.Load, proposalItems[i])
// 								}),
// 								layout.Rigid(func(gtx C) D {
// 									return pg.Theme.Separator().Layout(gtx)
// 								}),
// 							)
// 						})
// 					})
// 				})
// 			})
// 		}),
// 	)
// }

// func (pg *GovernancePage) layoutSyncSection(gtx C) D {
// 	if pg.isSyncing {
// 		return pg.layoutIsSyncingSection(gtx)
// 	} else if pg.syncCompleted {
// 		return pg.updatedIcon.Layout(gtx, values.MarginPadding20)
// 	}
// 	return pg.layoutStartSyncSection(gtx)
// }

// func (pg *ProposalsPage) layoutIsSyncingSection(gtx C) D {
// 	th := material.NewTheme(gofont.Collection())
// 	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding24)
// 	gtx.Constraints.Min.X = gtx.Constraints.Max.X
// 	loader := material.Loader(th)
// 	loader.Color = pg.Theme.Color.Gray1
// 	return loader.Layout(gtx)
// }

// func (pg *ProposalsPage) layoutStartSyncSection(gtx C) D {
// 	return material.Clickable(gtx, pg.syncButton, func(gtx C) D {
// 		return pg.Icons.Restore.Layout24dp(gtx)
// 	})
// }