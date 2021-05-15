package uiwallet

import (
	"errors"
	"image"
	"time"

	"gioui.org/app"
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

// Wallet represents the app window (and UI in general). There should only be one.
// Wallet uses an internal state of booleans to determine what the window is currently displaying.
type Wallet struct {
	ops                op.Ops
	theme              *decredmaterial.Theme
	wallet             *wallet.Wallet
	walletInfo         *wallet.MultiWalletInfo
	walletSyncStatus   *wallet.SyncStatus
	walletTransactions *wallet.Transactions
	walletTransaction  *wallet.Transaction
	walletAccount      *wallet.Account
	walletTickets      *wallet.Tickets
	vspInfo            *wallet.VSP
	proposals          *wallet.Proposals
	selectedProposal   *dcrlibwallet.Proposal
	proposal           chan *wallet.Proposal

	walletUnspentOutputs *wallet.UnspentOutputs

	current, previous string

	signatureResult *wallet.Signature

	selectedAccount int
	txAuthor        dcrlibwallet.TxAuthor
	broadcastResult wallet.Broadcast

	selected int
	states

	err string

	pages                   map[string]layout.Widget
	walletTabs, accountTabs *decredmaterial.Tabs
	keyEvents               chan *key.Event
	clipboard               chan interface{}
	toast                   chan *toast
	modal                   chan *modalLoad
	sysDestroyWithSync      bool
	walletAcctMixerStatus   chan *wallet.AccountMixer
	internalLog             chan string

	// Toggle between wallet and dex view mode
	switchView *int

	refreshWindow func()
}

type WriteClipboard struct {
	Text string
}

// NewWalletUI creates and initializes a new walletUI with start
func NewWalletUI(wal *wallet.Wallet, decredIcons map[string]image.Image, collection []text.FontFace,
	internalLog chan string, v *int, invalidate func()) (*Wallet, error) {
	wall := new(Wallet)
	theme := decredmaterial.NewTheme(collection, decredIcons)
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	wall.ops = op.Ops{}
	wall.theme = theme

	wall.walletInfo = new(wallet.MultiWalletInfo)
	wall.walletSyncStatus = new(wallet.SyncStatus)
	wall.walletTransactions = new(wallet.Transactions)
	wall.walletUnspentOutputs = new(wallet.UnspentOutputs)
	wall.walletAcctMixerStatus = make(chan *wallet.AccountMixer)
	wall.walletTickets = new(wallet.Tickets)
	wall.vspInfo = new(wallet.VSP)
	wall.proposals = new(wallet.Proposals)
	wall.proposal = make(chan *wallet.Proposal)

	wall.wallet = wal
	wall.states.loading = true
	wall.current = PageOverview
	wall.keyEvents = make(chan *key.Event)
	wall.clipboard = make(chan interface{}, 2)
	wall.toast = make(chan *toast)
	wall.modal = make(chan *modalLoad)
	wall.switchView = v
	wall.theme.ReadClipboard = wall.clipboard

	wall.walletTabs, wall.accountTabs = decredmaterial.NewTabs(wall.theme), decredmaterial.NewTabs(wall.theme)
	wall.walletTabs.Position, wall.accountTabs.Position = decredmaterial.Top, decredmaterial.Top
	wall.walletTabs.Separator, wall.walletTabs.Separator = false, false
	wall.accountTabs.SetTitle(wall.theme.Label(values.TextSize18, "Accounts:"))
	wall.walletTabs.SetTabs([]decredmaterial.TabItem{})
	wall.accountTabs.SetTabs([]decredmaterial.TabItem{})

	wall.internalLog = internalLog

	wall.addPages(decredIcons)

	wall.refreshWindow = invalidate

	return wall, nil
}

func (wall *Wallet) Ops() *op.Ops {
	return &wall.ops
}

func (wall *Wallet) changePage(page string) {
	wall.current = page
	wall.refresh()
}

func (wall *Wallet) refresh() {
	wall.refreshWindow()
}

func (wall *Wallet) setReturnPage(from string) {
	wall.previous = from
	wall.refresh()
}

func (wall *Wallet) unloaded(w *app.Window) {
	lbl := wall.theme.H3("Multiwallet not loaded\nIs another instance open?")
	var ops op.Ops

	for {
		e := <-w.Events()
		switch evt := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, evt)
			lbl.Layout(gtx)
			evt.Frame(gtx.Ops)
		}
	}
}

// Loop runs main event handling and page rendering loop
func (wall *Wallet) Run(shutdown chan int, w *app.Window) {
	for {
		select {
		case e := <-wall.wallet.Send:
			if e.Err != nil {
				err := e.Err.Error()
				log.Error("Wallet Error: " + err)
				if err == dcrlibwallet.ErrWalletDatabaseInUse {
					close(shutdown)
					wall.unloaded(w)
					return
				}
				wall.err = err
				if wall.states.loading {
					log.Warn("Attemping to get multiwallet info")
					wall.wallet.GetMultiWalletInfo()
				}

				w.Invalidate()
				break
			}

			wall.updateStates(e.Resp)

		case update := <-wall.wallet.Sync:
			switch update.Stage {
			case wallet.SyncCompleted:
				if wall.sysDestroyWithSync {
					close(shutdown)
					return
				}
				wall.updateSyncStatus(false, true)
			case wallet.SyncStarted:
				// dcrlibwallet triggers the SyncStart method several times
				// without sending a SyncComplete signal when sync is done.
				if !wall.walletInfo.Synced {
					wall.updateSyncStatus(true, false)
				}
			case wallet.SyncCanceled:
				if wall.sysDestroyWithSync {
					close(shutdown)
					return
				}
				wall.updateSyncStatus(false, false)
			case wallet.HeadersFetchProgress:
				wall.updateSyncProgress(update.ProgressReport)
			case wallet.AddressDiscoveryProgress:
				wall.updateSyncProgress(update.ProgressReport)
			case wallet.HeadersRescanProgress:
				wall.updateSyncProgress(update.ProgressReport)
			case wallet.PeersConnected:
				wall.updateConnectedPeers(update.ConnectedPeers)
			case wallet.BlockAttached:
				if wall.walletInfo.Synced {
					wall.wallet.GetAllTickets()
					wall.wallet.GetMultiWalletInfo()
					wall.updateSyncProgress(update.BlockInfo)
				}
			case wallet.BlockConfirmed:
				wall.updateSyncProgress(update.ConfirmedTxn)
			case wallet.AccountMixerStarted, wallet.AccountMixerEnded:
				go func() {
					wall.walletAcctMixerStatus <- &update.AcctMixerInfo
				}()
			case wallet.ProposalAdded, wallet.ProposalVoteFinished, wallet.ProposalVoteStarted, wallet.ProposalSynced:
				wall.wallet.GetAllProposals()
				go func() {
					wall.proposal <- &update.Proposal
				}()
			}
			w.Invalidate()

		case e := <-wall.clipboard:
			switch c := e.(type) {
			case decredmaterial.ReadClipboard:
				w.ReadClipboard()
			case WriteClipboard:
				go func() {
					w.WriteClipboard(c.Text)
					wall.toast <- &toast{
						text:    "copied",
						success: true,
					}
				}()
			}
		}
	}
}

func (wall *Wallet) HandlerDestroy(shutdown chan int) {
	if wall.walletInfo.Syncing || wall.walletInfo.Synced {
		wall.sysDestroyWithSync = true
		wall.wallet.CancelSync()
	} else {
		close(shutdown)
	}
}

func (wall *Wallet) HandlerPages(gtx layout.Context) {
	ts := int64(time.Since(time.Unix(wall.walletInfo.BestBlockTime, 0)).Seconds())
	wall.walletInfo.LastSyncTime = wallet.SecondsToDays(ts)
	s := wall.states
	if wall.walletInfo.LoadedWallets == 0 {
		wall.changePage(PageCreateRestore)
	}

	if s.loading {
		wall.Loading(gtx)
	} else {
		wall.theme.Background(gtx, wall.pages[wall.current])
	}
}

func (wall *Wallet) HandlerKeyEvents(evt *key.Event) {
	go func() {
		wall.keyEvents <- evt
	}()
}

func (wall *Wallet) HandlerClipboard(evt *clipboard.Event) {
	go func() {
		wall.theme.Clipboard <- evt.Text
	}()
}
