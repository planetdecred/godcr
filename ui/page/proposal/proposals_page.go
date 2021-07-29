package proposal

import (
	"sync"
	"time"

	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/wallet"
)

const ProposalsPageID = "Proposals"

type proposalItem struct {
	proposal     dcrlibwallet.Proposal
	voteBar      decredmaterial.VoteBar
	tooltip      *decredmaterial.Tooltip
	tooltipLabel decredmaterial.Label
}

type Page struct {
	*load.Load

	pageClosing chan bool
	proposalMu  sync.Mutex

	selectedProposal **dcrlibwallet.Proposal
	multiWallet      *dcrlibwallet.MultiWallet

	categoryList  *decredmaterial.ClickableList
	proposalsList *decredmaterial.ClickableList

	tabCard      decredmaterial.Card
	itemCard     decredmaterial.Card
	syncCard     decredmaterial.Card
	updatedLabel decredmaterial.Label

	proposalItems         []proposalItem
	proposalCount         []int
	selectedCategoryIndex int

	legendIcon    *widget.Icon
	infoIcon      *widget.Icon
	updatedIcon   *widget.Icon
	syncButton    *widget.Clickable
	startSyncIcon *widget.Image
	timerIcon     *widget.Image

	showSyncedCompleted bool
	isSyncing           bool
}

var (
	proposalCategoryTitles = []string{"In discussion", "Voting", "Approved", "Rejected", "Abandoned"}
	proposalCategories     = []int32{
		dcrlibwallet.ProposalCategoryPre,
		dcrlibwallet.ProposalCategoryActive,
		dcrlibwallet.ProposalCategoryApproved,
		dcrlibwallet.ProposalCategoryRejected,
		dcrlibwallet.ProposalCategoryAbandoned,
	}
)

func NewProposalsPage(l *load.Load) *Page {
	pg := &Page{
		Load:        l,
		pageClosing: make(chan bool, 1),
	}
	pg.initLayoutWidgets()

	return pg
}

func (pg *Page) OnResume() {
	pg.listenForSyncNotifications()

	pg.proposalMu.Lock()
	selectedCategory := pg.selectedCategoryIndex
	pg.proposalMu.Unlock()
	if selectedCategory == -1 {
		pg.countProposals()
		pg.loadProposals(0)
	}

	pg.isSyncing = pg.multiWallet.Politeia.IsSyncing()
}

func (pg *Page) countProposals() {
	proposalCount := make([]int, len(proposalCategories))
	for i, category := range proposalCategories {
		count, err := pg.multiWallet.Politeia.Count(category)
		if err == nil {
			proposalCount[i] = int(count)
		}
	}

	pg.proposalMu.Lock()
	pg.proposalCount = proposalCount
	pg.proposalMu.Unlock()
}

func (pg *Page) loadProposals(category int) {
	proposals, err := pg.multiWallet.Politeia.GetProposalsRaw(proposalCategories[category], 0, 0, true)
	if err != nil {
		pg.proposalMu.Lock()
		pg.proposalItems = make([]proposalItem, 0)
		pg.proposalMu.Unlock()
	} else {
		proposalItems := make([]proposalItem, len(proposals))
		for i := 0; i < len(proposals); i++ {
			proposal := proposals[i]
			item := proposalItem{
				proposal: proposals[i],
				voteBar:  pg.Theme.VoteBar(pg.infoIcon, pg.legendIcon),
			}

			if proposal.Category == dcrlibwallet.ProposalCategoryPre {
				tooltipLabel := pg.Theme.Caption("")
				tooltipLabel.Color = pg.Theme.Color.Gray
				if proposal.VoteStatus == 1 {
					tooltipLabel.Text = "Waiting for author to authorize voting"
				} else if proposal.VoteStatus == 2 {
					tooltipLabel.Text = "Waiting for admin to trigger the start of voting"
				}

				item.tooltip = pg.Theme.Tooltip()
				item.tooltipLabel = tooltipLabel
			}

			proposalItems[i] = item
		}
		pg.proposalMu.Lock()
		pg.selectedCategoryIndex = category
		pg.proposalItems = proposalItems
		pg.proposalMu.Unlock()
	}
}

func (pg *Page) Handle() {
	if clicked, selectedItem := pg.categoryList.ItemClicked(); clicked {
		go pg.loadProposals(selectedItem)
	}

	if clicked, selectedItem := pg.proposalsList.ItemClicked(); clicked {
		pg.proposalMu.Lock()
		selectedProposal := pg.proposalItems[selectedItem].proposal
		pg.proposalMu.Unlock()

		pg.SetReturnPage(ProposalsPageID)
		pg.ChangeFragment(newProposalDetailsPage(pg.Load, selectedProposal), PageProposalDetails)
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

func (pg *Page) listenForSyncNotifications() {
	go func() {
		for {
			var notification interface{}

			select {
			case notification = <-pg.Receiver.NotificationsUpdate:
			case <-pg.pageClosing:
				return
			}

			switch n := notification.(type) {
			case wallet.Proposal:
				if n.ProposalStatus == wallet.Synced {
					pg.isSyncing = false
					pg.showSyncedCompleted = true

					pg.proposalMu.Lock()
					selectedCategory := pg.selectedCategoryIndex
					pg.proposalMu.Unlock()
					if selectedCategory != -1 {
						pg.countProposals()
						pg.loadProposals(selectedCategory)
					}
				}
			}
		}
	}()
}

func (pg *Page) OnClose() {
	pg.pageClosing <- true
}
