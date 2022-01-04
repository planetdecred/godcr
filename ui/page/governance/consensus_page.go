package governance

import (
	"context"
	"sync"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const ConsensusPageID = "Consensus"

type ConsensusPage struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc
	agendaMu  sync.Mutex

	multiWallet       *dcrlibwallet.MultiWallet
	listContainer     *widget.List
	walletDropDown    *decredmaterial.DropDown
	orderDropDown     *decredmaterial.DropDown
	consensusList     *decredmaterial.ClickableList
	syncButton        *widget.Clickable
	searchEditor      decredmaterial.Editor
	fetchProposalsBtn decredmaterial.Button

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton
	voteButton decredmaterial.Button

	updatedIcon *decredmaterial.Icon

	consensusItems []*components.ConsensusItem
	wallets        []*dcrlibwallet.Wallet

	syncCompleted bool
	isSyncing     bool
}

func NewConsensusPage(l *load.Load) *ConsensusPage {
	pg := &ConsensusPage{
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

	pg.consensusList = pg.Theme.NewClickableList(layout.Vertical)

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	pg.voteButton = l.Theme.Button("Change Vote")

	pg.wallets = pg.WL.SortedWalletList()
	pg.walletDropDown = components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets, values.TxDropdownGroup, 0)
	pg.orderDropDown = components.CreateOrderDropDown(l, values.ConsensusDropdownGroup, 0)

	pg.initLayoutWidgets()

	return pg
}

func (pg *ConsensusPage) initLayoutWidgets() {
	//categoryList to be removed with new update to UI.
	// pg.consensusList = pg.Theme.NewClickableList(layout.Horizontal)
	// pg.itemCard = pg.Theme.Card()

}

func (pg *ConsensusPage) ID() string {
	return ConsensusPageID
}

func (pg *ConsensusPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.listenForSyncNotifications()
	pg.fetchAgendas()
	pg.isSyncing = pg.multiWallet.Consensus.IsSyncing()
}

func (pg *ConsensusPage) OnClose() {
	pg.ctxCancel()
}

func (pg *ConsensusPage) Handle() {

}

func (pg *ConsensusPage) fetchAgendas() {
	newestFirst := pg.orderDropDown.SelectedIndex() == 0

	selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
	consensusItems := components.LoadAgendas(pg.Load, selectedWallet, newestFirst)

	// group 'In discussion' and 'Active' proposals into under review
	listItems := make([]*components.ConsensusItem, 0)
	for _, item := range consensusItems {
		listItems = append(listItems, item)
	}

	pg.agendaMu.Lock()
	pg.consensusItems = listItems
	pg.agendaMu.Unlock()
}

func (pg *ConsensusPage) Layout(gtx C) D {
	if pg.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					// layout.Rigid(pg.backButton.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
							body := func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										var text string
										if pg.isSyncing {
											text = "Syncing..."
										} else if pg.syncCompleted {
											text = "Updated"
										} else {
											text = "Upated " + components.TimeAgo(pg.multiWallet.Consensus.GetLastSyncedTimeStamp())
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
								layout.Flexed(1, func(gtx C) D {
									return layout.E.Layout(gtx, body)
								}),
							)
						})
					}),
				)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
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
							return pg.walletDropDown.Layout(gtx, pg.orderDropDown.Width+41, true)
						}),
					)
				})
			}),
		)
	}
	return D{}
}

func (pg *ConsensusPage) layoutContent(gtx C) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			// pg.proposalMu.Lock()
			consensusItems := pg.consensusItems
			// pg.proposalMu.Unlock()

			return pg.Theme.List(pg.listContainer).Layout(gtx, 1, func(gtx C, i int) D {
				return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
					return pg.Theme.Card().Layout(gtx, func(gtx C) D {
						if len(consensusItems) == 0 {
							return components.LayoutNoAgendasFound(gtx, pg.Load, pg.isSyncing)
						}
						return pg.consensusList.Layout(gtx, len(consensusItems), func(gtx C, i int) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return components.AgendasList(gtx, pg.Load, consensusItems[i])
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

func (pg *ConsensusPage) layoutSyncSection(gtx C) D {
	if pg.isSyncing {
		return pg.layoutIsSyncingSection(gtx)
	} else if pg.syncCompleted {
		return pg.updatedIcon.Layout(gtx, values.MarginPadding20)
	}
	return pg.layoutStartSyncSection(gtx)
}

func (pg *ConsensusPage) layoutIsSyncingSection(gtx C) D {
	th := material.NewTheme(gofont.Collection())
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding24)
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	loader := material.Loader(th)
	loader.Color = pg.Theme.Color.Gray1
	return loader.Layout(gtx)
}

func (pg *ConsensusPage) layoutStartSyncSection(gtx C) D {
	return material.Clickable(gtx, pg.syncButton, func(gtx C) D {
		return pg.Icons.Restore.Layout24dp(gtx)
	})
}

func (pg *ConsensusPage) listenForSyncNotifications() {
	go func() {
		for {
			var notification interface{}

			select {
			case notification = <-pg.Receiver.NotificationsUpdate:
			case <-pg.ctx.Done():
				return
			}

			switch n := notification.(type) {
			case wallet.Agenda:
				if n.AgendaStatus == wallet.SyncedAgenda {
					pg.syncCompleted = true
					pg.isSyncing = false

					pg.fetchAgendas()
					pg.RefreshWindow()
				}
			}
		}
	}()
}