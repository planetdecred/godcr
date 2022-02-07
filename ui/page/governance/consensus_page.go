package governance

import (
	"context"
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
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
	LiveTickets    []*dcrlibwallet.Transaction

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

	return pg
}

func (pg *ConsensusPage) ID() string {
	return ConsensusPageID
}


func (pg *ConsensusPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.FetchAgendas()
}

func (pg *ConsensusPage) OnNavigatedFrom() {
	pg.ctxCancel()
}

func (pg *ConsensusPage) HandleUserInteractions() {

}

func (pg *ConsensusPage) FetchAgendas() {
	newestFirst := pg.orderDropDown.SelectedIndex() == 0

	selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
	consensusItems := components.LoadAgendas(pg.Load, selectedWallet, newestFirst)

	listItems := make([]*components.ConsensusItem, 0)
	for _, item := range consensusItems {
		listItems = append(listItems, item)
	}

	pg.agendaMu.Lock()
	pg.consensusItems = listItems
	pg.agendaMu.Unlock()
	pg.RefreshWindow()
}

func (pg *ConsensusPage) Layout(gtx C) D {
	if pg.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
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
			consensusItems := pg.consensusItems

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
