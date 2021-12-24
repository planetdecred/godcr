package governance

import (
	"context"
	// "fmt"
	// "image"
	// "image/color"
	// "strconv"
	// "strings"
	"sync"
	// "time"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	// "gioui.org/op/clip"
	// "gioui.org/op/paint"
	// "gioui.org/text"
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

const ConsensusPageID = "Consensus"

type ConsensusPage struct {
	*load.Load

	ctx        context.Context // page context
	ctxCancel  context.CancelFunc
	agendaMu sync.Mutex

	multiWallet       *dcrlibwallet.MultiWallet
	listContainer     *widget.List
	orderDropDown     *decredmaterial.DropDown
	consensusList     *decredmaterial.ClickableList
	syncButton        *widget.Clickable
	searchEditor      decredmaterial.Editor
	fetchProposalsBtn decredmaterial.Button

	backButton decredmaterial.IconButton
	infoButton decredmaterial.IconButton
	voteButton               decredmaterial.Button

	updatedIcon *decredmaterial.Icon

	consensusItems []*components.ConsensusItem

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

}

func (pg *ConsensusPage) OnClose() {
	// pg.ctxCancel()
}

func (pg *ConsensusPage) Handle() {

}

func (pg *ConsensusPage) fetchAgendas() {
	// newestFirst := pg.orderDropDown.SelectedIndex() == 0

	consensusItems := components.LoadAgendas(pg.Load)

	// group 'In discussion' and 'Active' proposals into under review
	listItems := make([]*components.ConsensusItem, 0)
	for _, item := range consensusItems {
		listItems = append(listItems, item)
	}

	pg.agendaMu.Lock()
	pg.consensusItems = listItems
	// if proposalFilter == dcrlibwallet.ProposalCategoryAll {
	// 	pg.proposalItems = listItems
	// }
	pg.agendaMu.Unlock()
}

func (pg *ConsensusPage) Layout(gtx C) D {
	if pg.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
		// return components.UniformPadding(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
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
					)
					})
				}),
			)
		// })
	}
	return D{}
}

// func (pg *ConsensusPage) layoutAgendaVoteAction(gtx C, l *load.Load, item *ConsensusItem) D {
// 	gtx.Constraints.Min.X, gtx.Constraints.Max.X = 150, 150
// 	// var voteButton decredmaterial.Button
// 	// pg.VoteButton = l.Theme.Button("Change Vote")
// 	if canVote {
// 		pg.VoteButton.Background = l.Theme.Color.Primary
// 	} else {
// 		pg.VoteButton.Background = l.Theme.Color.Gray3
// 	}
// 	return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
// 		return pg.VoteButton.Layout(gtx)
// 	})
// }

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
									// txt := pg.Theme.Label(values.MarginPadding24, "Sample agenda item")
									// txt.Font.Weight = text.SemiBold
		
									// return layout.Inset{
									// 	Top:    values.MarginPadding30,
									// 	Bottom: values.MarginPadding16,
									// }.Layout(gtx, txt.Layout)
									// return pg.Theme.Label(values.MarginPadding24, "How does Governance Work?")
									// return components.ProposalsList(gtx, pg.Load, consensusItems[i])

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
