package governance

import (
	"context"
	"sync"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const ProposalsPageID = "Proposals"

type (
	C = layout.Context
	D = layout.Dimensions
)

type ProposalsPage struct {
	*load.Load

	*listeners.ProposalNotificationListener
	ctx        context.Context // page context
	ctxCancel  context.CancelFunc
	proposalMu sync.Mutex

	multiWallet       *dcrlibwallet.MultiWallet
	listContainer     *widget.List
	orderDropDown     *decredmaterial.DropDown
	categoryDropDown  *decredmaterial.DropDown
	proposalsList     *decredmaterial.ClickableList
	syncButton        *widget.Clickable
	searchEditor      decredmaterial.Editor
	fetchProposalsBtn decredmaterial.Button

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton

	updatedIcon *decredmaterial.Icon

	proposalItems []*components.ProposalItem

	syncCompleted bool
	isSyncing     bool
}

func NewProposalsPage(l *load.Load) *ProposalsPage {
	pg := &ProposalsPage{
		Load:        l,
		multiWallet: l.WL.MultiWallet,
		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
	pg.searchEditor = l.Theme.IconEditor(new(widget.Editor), "Search", l.Icons.SearchIcon, true)
	pg.searchEditor.Editor.SingleLine, pg.searchEditor.Editor.Submit, pg.searchEditor.Bordered = true, true, false

	pg.updatedIcon = decredmaterial.NewIcon(pg.Icons.NavigationCheck)
	pg.updatedIcon.Color = pg.Theme.Color.Success

	pg.syncButton = new(widget.Clickable)

	pg.proposalsList = pg.Theme.NewClickableList(layout.Vertical)

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	// orderDropDown is the first dropdown when page is laid out. Its
	// position should be 0 for consistent backdrop.
	pg.orderDropDown = components.CreateOrderDropDown(l, values.ProposalDropdownGroup, 0)
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
	}, values.ProposalDropdownGroup, 1)

	pg.initializeWidget()

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *ProposalsPage) ID() string {
	return ProposalsPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *ProposalsPage) OnNavigatedTo() {
	pg.ProposalNotificationListener = listeners.NewProposalNotificationListener()
	pg.WL.MultiWallet.Politeia.AddNotificationListener(pg.ProposalNotificationListener, ProposalsPageID)

	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.listenForSyncNotifications()
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

	proposalItems := components.LoadProposals(proposalFilter, newestFirst, pg.Load)

	// group 'In discussion' and 'Active' proposals into under review
	listItems := make([]*components.ProposalItem, 0)
	for _, item := range proposalItems {
		if item.Proposal.Category == dcrlibwallet.ProposalCategoryPre ||
			item.Proposal.Category == dcrlibwallet.ProposalCategoryActive {
			listItems = append(listItems, item)
		}
	}

	pg.proposalMu.Lock()
	pg.proposalItems = proposalItems
	if proposalFilter == dcrlibwallet.ProposalCategoryAll {
		pg.proposalItems = listItems
	}
	pg.proposalMu.Unlock()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *ProposalsPage) HandleUserInteractions() {
	for pg.fetchProposalsBtn.Clicked() {
		go pg.WL.MultiWallet.Politeia.Sync()
		pg.isSyncing = pg.multiWallet.Politeia.IsSyncing()
		pg.WL.Wallet.SaveConfigValueForKey(load.FetchProposalConfigKey, true)
	}

	for pg.infoButton.Button.Clicked() {
		pg.showInfoModal()
	}

	for pg.backButton.Button.Clicked() {
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
		selectedProposal := pg.proposalItems[selectedItem].Proposal
		pg.proposalMu.Unlock()

		pg.ChangeFragment(NewProposalDetailsPage(pg.Load, &selectedProposal))
	}

	for pg.syncButton.Clicked() {
		go pg.multiWallet.Politeia.Sync()
		pg.isSyncing = true

		//Todo: check after 1min if sync does not start, set isSyncing to false and cancel sync
	}

	if pg.syncCompleted {
		time.AfterFunc(time.Second*3, func() {
			pg.syncCompleted = false
			pg.RefreshWindow()
		})
	}

	decredmaterial.DisplayOneDropdown(pg.orderDropDown, pg.categoryDropDown)
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *ProposalsPage) OnNavigatedFrom() {
	pg.ctxCancel()

	pg.WL.MultiWallet.Politeia.RemoveNotificationListener(ProposalsPageID)
	close(pg.ProposalNotifChan)
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *ProposalsPage) Layout(gtx C) D {
	if pg.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
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
											} else if pg.syncCompleted {
												text = "Updated"
											} else {
												text = "Upated " + components.TimeAgo(pg.multiWallet.Politeia.GetLastSyncedTimeStamp())
											}

											lastUpdatedInfo := pg.Theme.Label(values.TextSize10, text)
											lastUpdatedInfo.Color = pg.Theme.Color.GrayText2
											if pg.syncCompleted {
												lastUpdatedInfo.Color = pg.Theme.Color.Success
											}

											return layout.Inset{Top: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
												return lastUpdatedInfo.Layout(gtx)
											})
										}),
									)
								}

								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												txt := pg.Theme.Label(values.TextSize20, "Proposals")
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
							return pg.orderDropDown.Layout(gtx, 45, true)
						}),
						layout.Expanded(func(gtx C) D {
							return pg.categoryDropDown.Layout(gtx, pg.orderDropDown.Width+41, true)
						}),
					)
				}),
			)
		})
	}
	return components.UniformPadding(gtx, pg.splashScreenLayout)
}

func (pg *ProposalsPage) layoutContent(gtx C) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			pg.proposalMu.Lock()
			proposalItems := pg.proposalItems
			pg.proposalMu.Unlock()

			return pg.Theme.List(pg.listContainer).Layout(gtx, 1, func(gtx C, i int) D {
				return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
					return pg.Theme.Card().Layout(gtx, func(gtx C) D {
						if len(proposalItems) == 0 {
							return components.LayoutNoProposalsFound(gtx, pg.Load, pg.isSyncing, int32(pg.categoryDropDown.SelectedIndex()))
						}
						return pg.proposalsList.Layout(gtx, len(proposalItems), func(gtx C, i int) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return components.ProposalsList(gtx, pg.Load, proposalItems[i])
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
	if pg.isSyncing {
		return pg.layoutIsSyncingSection(gtx)
	} else if pg.syncCompleted {
		return pg.updatedIcon.Layout(gtx, values.MarginPadding20)
	}
	return pg.layoutStartSyncSection(gtx)
}

func (pg *ProposalsPage) layoutIsSyncingSection(gtx C) D {
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding24)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	loader := material.Loader(pg.Theme.Base)
	loader.Color = pg.Theme.Color.Gray1
	return loader.Layout(gtx)
}

func (pg *ProposalsPage) layoutStartSyncSection(gtx C) D {
	return material.Clickable(gtx, pg.syncButton, func(gtx C) D {
		return pg.Icons.Restore.Layout24dp(gtx)
	})
}

func (pg *ProposalsPage) listenForSyncNotifications() {

	go func() {
		for {
			select {
			case n := <-pg.ProposalNotifChan:
				if n.ProposalStatus == wallet.Synced {
					pg.syncCompleted = true
					pg.isSyncing = false

					pg.fetchProposals()
					pg.RefreshWindow()
				}
			case <-pg.ctx.Done():
				return
			}
		}
	}()
}
