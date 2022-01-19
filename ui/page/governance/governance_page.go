package governance

import (
	"context"
	// "fmt"
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const GovernancePageID = "Governance"

type GovernancePage struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	multiWallet *dcrlibwallet.MultiWallet

	tabCategoryList *decredmaterial.ClickableList
	listContainer   *widget.List
	tabCard         decredmaterial.Card
	itemCard        decredmaterial.Card

	selectedCategoryIndex int

	proposalsPage *ProposalsPage
	consensusPage *ConsensusPage

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton

	enableGovernanceBtn decredmaterial.Button
}

var (
	governanceTabTitles = []string{"Proposals", "Consensus Changes"}
)

func NewGovernancePage(l *load.Load) *GovernancePage {
	pg := &GovernancePage{
		Load:                  l,
		multiWallet:           l.WL.MultiWallet,
		selectedCategoryIndex: -1,
		proposalsPage:         NewProposalsPage(l),
		consensusPage:         NewConsensusPage(l),
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	pg.initTabWidgets()
	pg.initializeWidget()

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

	// pg.consensusPage.listenForSyncNotifications()
	pg.consensusPage.fetchAgendas()

	/** end consensus page resume method */
}

func (pg *GovernancePage) OnClose() {
	// pg.ctxCancel()
	// pg.proposalsPage.ctxCancel()
	// pg.consensusPage.ctxCancel()
}

func (pg *GovernancePage) initTabWidgets() {
	pg.tabCategoryList = pg.Theme.NewClickableList(layout.Horizontal)
	pg.itemCard = pg.Theme.Card()

	radius := decredmaterial.Radius(0)
	pg.tabCard = pg.Theme.Card()
	pg.tabCard.Radius = radius
}

func (pg *GovernancePage) ID() string {
	return GovernancePageID
}

func (pg *GovernancePage) Handle() {
	for pg.enableGovernanceBtn.Clicked() {
		go pg.WL.MultiWallet.Politeia.Sync()
		pg.proposalsPage.isSyncing = pg.proposalsPage.multiWallet.Politeia.IsSyncing()
		pg.WL.Wallet.SaveConfigValueForKey(load.FetchProposalConfigKey, true)
	}

	for pg.backButton.Button.Clicked() {
		pg.PopFragment()
	}

	if clicked, selectedItem := pg.tabCategoryList.ItemClicked(); clicked {
		pg.selectedCategoryIndex = selectedItem
	}

	/** begin proposal page handles */

	for pg.proposalsPage.infoButton.Button.Clicked() {
		// pg.proposalsPage.showInfoModal()
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

	for pg.consensusPage.walletDropDown.Changed() {
		pg.consensusPage.fetchAgendas()
	}

	for pg.consensusPage.orderDropDown.Changed() {
		pg.consensusPage.fetchAgendas()
	}

	for i := range pg.consensusPage.consensusItems {
		if pg.consensusPage.consensusItems[i].VoteButton.Clicked() {
			newAgendaVoteModal(pg.Load, &pg.consensusPage.consensusItems[i].Agenda).Show()
		}
	}

	for pg.consensusPage.syncButton.Clicked() {
		pg.consensusPage.isSyncing = true

		//Todo: check after 1min if sync does not start, set isSyncing to false and cancel sync
	}

	if pg.consensusPage.syncCompleted {
		time.AfterFunc(time.Second*3, func() {
			pg.consensusPage.syncCompleted = false
			pg.RefreshWindow()
		})
	}

	/** end consensus page handles */

}

func (pg *GovernancePage) Layout(gtx C) D {
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
								return pg.switchTab(gtx, pg.selectedCategoryIndex)
							}),
						)
					})
				}),
			)
		})
	}
	return components.UniformPadding(gtx, pg.splashScreenLayout)
}

func (pg *GovernancePage) switchTab(gtx C, selectedCategoryIndex int) D {
	if selectedCategoryIndex == 0 {
		return pg.proposalsPage.Layout(gtx)
	}

	return pg.consensusPage.Layout(gtx)
}

func (pg *GovernancePage) layoutTabs(gtx C) D {
	width := float32(gtx.Constraints.Max.X-20) / float32(len(governanceTabTitles))
	// pg.proposalMu.Lock()
	selectedCategory := pg.selectedCategoryIndex
	// pg.proposalMu.Unlock()

	return pg.tabCard.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Left:  values.MarginPadding12,
			Right: values.MarginPadding12,
		}.Layout(gtx, func(gtx C) D {
			return pg.tabCategoryList.Layout(gtx, len(governanceTabTitles), func(gtx C, i int) D {
				gtx.Constraints.Min.X = int(width)
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding14).Layout(gtx, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										lbl := pg.Theme.Body1(governanceTabTitles[i])
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