package ui

import (
	"errors"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

// Window represents the app window (and UI in general). There should only be one.
// Window uses an internal state of booleans to determine what the window is currently displaying.
type Window struct {
	ops        *op.Ops
	invalidate chan struct{}

	wallet               *wallet.Wallet
	walletInfo           *wallet.MultiWalletInfo
	walletSyncStatus     *wallet.SyncStatus
	walletTransactions   *wallet.Transactions
	walletTransaction    *wallet.Transaction
	walletAccount        *wallet.Account
	vspInfo              *wallet.VSP
	proposals            *wallet.Proposals
	selectedProposal     *dcrlibwallet.Proposal
	proposal             chan *wallet.Proposal
	walletUnspentOutputs *wallet.UnspentOutputs

	load *load.Load

	modalMutex sync.Mutex
	modals     []load.Modal

	currentPage   load.Page
	pageBackStack []load.Page

	signatureResult *wallet.Signature

	selectedAccount int
	txAuthor        dcrlibwallet.TxAuthor
	broadcastResult wallet.Broadcast

	selected int
	states   states

	err string

	keyEvents             map[string]chan *key.Event
	sysDestroyWithSync    bool
	walletAcctMixerStatus chan *wallet.AccountMixer
	internalLog           chan string
}

type (
	C = layout.Context
	D = layout.Dimensions
)
type WriteClipboard struct {
	Text string
}

// CreateWindow creates and initializes a new window with start
// as the first page displayed.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func CreateWindow(wal *wallet.Wallet) (*Window, *app.Window, error) {
	win := new(Window)
	var netType string
	if wal.Net == dcrlibwallet.Testnet3 {
		netType = "testnet"
	} else {
		netType = wal.Net
	}
	appWindow := app.NewWindow(app.MinSize(values.AppWidth, values.AppHeight), app.Title(values.StringF(values.StrAppTitle, netType)))
	win.ops = &op.Ops{}

	win.walletInfo = new(wallet.MultiWalletInfo)
	win.walletSyncStatus = new(wallet.SyncStatus)
	win.walletTransactions = new(wallet.Transactions)
	win.walletUnspentOutputs = new(wallet.UnspentOutputs)
	win.walletAcctMixerStatus = make(chan *wallet.AccountMixer)
	win.vspInfo = new(wallet.VSP)
	win.proposals = new(wallet.Proposals)
	win.proposal = make(chan *wallet.Proposal)
	win.invalidate = make(chan struct{}, 2)

	win.wallet = wal
	win.states.loading = false

	win.keyEvents = make(map[string]chan *key.Event)

	l, err := win.NewLoad()
	if err != nil {
		return nil, nil, err
	}

	win.load = l

	return win, appWindow, nil
}

func (win *Window) NewLoad() (*load.Load, error) {
	l, err := load.NewLoad()
	if err != nil {
		return nil, err
	}

	l.WL = &load.WalletLoad{
		MultiWallet:     win.wallet.GetMultiWallet(),
		Wallet:          win.wallet,
		Account:         win.walletAccount,
		Info:            win.walletInfo,
		SyncStatus:      win.walletSyncStatus,
		Transactions:    win.walletTransactions,
		UnspentOutputs:  win.walletUnspentOutputs,
		VspInfo:         win.vspInfo,
		BroadcastResult: win.broadcastResult,
		Proposals:       win.proposals,

		SelectedProposal: win.selectedProposal,
		TxAuthor:         win.txAuthor,
	}

	l.Receiver = &load.Receiver{
		KeyEvents:           win.keyEvents,
		AcctMixerStatus:     win.walletAcctMixerStatus,
		InternalLog:         win.internalLog,
		SyncedProposal:      win.proposal,
		NotificationsUpdate: make(chan interface{}, 10),
		WalletRestored:      make(chan struct{}),
		AllWalletsDeleted:   make(chan struct{}),
	}

	l.SelectedWallet = &win.selected
	l.RefreshWindow = win.refreshWindow
	l.ShowModal = win.showModal
	l.DismissModal = win.dismissModal
	l.PopWindowPage = win.popPage
	l.ChangeWindowPage = win.changePage
	l.SubscribeKeyEvent = win.SubscribeKeyEvent
	l.UnsubscribeKeyEvent = win.UnsubscribeKeyEvent

	return l, nil
}

func (win *Window) changePage(page load.Page, keepBackStack bool) {
	if win.currentPage != nil && keepBackStack {
		win.currentPage.WillDisappear() // TODO: Unload() if not keeping in backstack.
		win.pageBackStack = append(win.pageBackStack, win.currentPage)
	}

	win.currentPage = page
	win.currentPage.WillAppear() // callers shouldn't need to trigger WillAppear(), page is changing, WillAppear naturally should be triggered here
	win.refreshWindow()          // TODO: Ensure no caller of this method also triggers refreshWindow!
}

// popPage goes back to the previous page
// returns true if page was popped.
func (win *Window) popPage() bool {
	if len(win.pageBackStack) == 0 {
		return false
	}

	// get and remove last page
	previousPageIndex := len(win.pageBackStack) - 1
	previousPage := win.pageBackStack[previousPageIndex]
	win.pageBackStack = win.pageBackStack[:previousPageIndex]

	// close the current page and display the previous page
	win.currentPage.WillDisappear() // Use Unload() for page that's being closed.
	previousPage.WillAppear()
	win.currentPage = previousPage
	win.refreshWindow() // TODO: Ensure no caller of this method also triggers refreshWindow!

	return true
}

func (win *Window) refreshWindow() {
	win.invalidate <- struct{}{}
}

func (win *Window) showModal(modal load.Modal) {
	modal.OnResume() // setup display data
	win.modalMutex.Lock()
	win.modals = append(win.modals, modal)
	win.modalMutex.Unlock()
}

func (win *Window) dismissModal(modal load.Modal) {
	win.modalMutex.Lock()
	defer win.modalMutex.Unlock()
	for i, m := range win.modals {
		if m.ModalID() == modal.ModalID() {
			modal.OnDismiss() // do garbage collection in modal
			win.modals = append(win.modals[:i], win.modals[i+1:]...)
			win.refreshWindow()
		}
	}
}

func (win *Window) unloaded(w *app.Window) {
	for {
		e := <-w.Events()
		switch evt := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			gtx := layout.NewContext(win.ops, evt)
			lbl := win.load.Theme.H3("Multiwallet not loaded\nIs another instance open?")
			lbl.Layout(gtx)
			evt.Frame(win.ops)
		}
	}
}

// SubscribeKeyEvent subscribes pages for key events.
func (win *Window) SubscribeKeyEvent(eventChan chan *key.Event, pageID string) {
	win.keyEvents[pageID] = eventChan
}

// UnsubscribeKeyEvent unsubscribe a page with {pageID} from receiving key events.
func (win *Window) UnsubscribeKeyEvent(pageID string) error {
	if _, ok := win.keyEvents[pageID]; ok {
		delete(win.keyEvents, pageID)
		return nil
	}

	return errors.New("Page not subscribed for key events")
}

// Loop runs main event handling and page rendering loop
func (win *Window) Loop(w *app.Window, shutdown chan int) {
	for {
		select {
		case <-win.invalidate:
			w.Invalidate()
		case e := <-win.wallet.Send:
			if e.Err != nil {
				err := e.Err.Error()
				log.Error("Wallet Error: " + err)
				if err == dcrlibwallet.ErrWalletDatabaseInUse {
					close(shutdown)
					win.unloaded(w) // This method starts an infite loop. Check.
					return
				}
				win.err = err
				if win.states.loading {
					log.Warn("Attemping to get multiwallet info")
					win.wallet.GetMultiWalletInfo()
				}

				op.InvalidateOp{}.Add(win.ops)
				break
			}

			// win.updateStates(e.Resp)

		case update := <-win.wallet.Sync:
			switch update.Stage {
			case wallet.SyncCompleted:
				if win.sysDestroyWithSync {
					close(shutdown)
					return
				}
				win.updateSyncStatus(false, true)
			case wallet.SyncStarted:
				// dcrlibwallet triggers the SyncStart method several times
				// without sending a SyncComplete signal when sync is done.
				if !win.walletInfo.Synced {
					win.updateSyncStatus(true, false)
				}
			case wallet.SyncCanceled:
				if win.sysDestroyWithSync {
					close(shutdown)
					return
				}
				win.updateSyncStatus(false, false)
			case wallet.HeadersFetchProgress:
				win.updateSyncProgress(update.ProgressReport)
			case wallet.AddressDiscoveryProgress:
				win.updateSyncProgress(update.ProgressReport)
			case wallet.HeadersRescanProgress:
				win.updateSyncProgress(update.ProgressReport)
			case wallet.PeersConnected:
				win.updateConnectedPeers(update.ConnectedPeers)
			case wallet.BlockAttached:
				if win.walletInfo.Synced {
					win.wallet.GetMultiWalletInfo()
					win.updateSyncProgress(update.BlockInfo)
				}
			case wallet.BlockConfirmed:
				win.updateSyncProgress(update.ConfirmedTxn)
			case wallet.AccountMixerStarted, wallet.AccountMixerEnded:
				go func() {
					win.walletAcctMixerStatus <- &update.AcctMixerInfo
				}()
			case wallet.ProposalAdded, wallet.ProposalVoteFinished, wallet.ProposalVoteStarted, wallet.ProposalSynced:
				win.wallet.GetAllProposals()
				go func() {
					win.proposal <- &update.Proposal
				}()
			}
			op.InvalidateOp{}.Add(win.ops)

		case e := <-w.Events():
			switch evt := e.(type) {

			case system.StageEvent:
				if evt.Stage == system.StageRunning {
					err := win.wallet.InitMultiWallet()
					if err != nil {
						if err.Error() == dcrlibwallet.ErrWalletDatabaseInUse {
							close(shutdown)
							win.unloaded(w)
							return
						}
					}
				}
			case system.DestroyEvent:
				if win.currentPage != nil {
					win.currentPage.WillDisappear()
					win.currentPage = nil
				}
				if win.walletInfo.Syncing || win.walletInfo.Synced {
					win.sysDestroyWithSync = true
					win.wallet.CancelSync()
				} else {
					close(shutdown)
				}

			case system.FrameEvent:
				ts := int64(time.Since(time.Unix(win.walletInfo.BestBlockTime, 0)).Seconds())
				win.walletInfo.LastSyncTime = wallet.SecondsToDays(ts) // TODO: Investigate
				win.displayWindow(evt)

			case key.Event:
				go func() {
					for _, c := range win.keyEvents {
						c <- &evt
					}
				}()
			case nil:
				// Ignore
			default:
				log.Tracef("Unhandled window event %+v\n", e)
			}
		case <-win.load.Receiver.WalletRestored:
			win.changePage(page.NewMainPage(win.load), false)
		case <-win.load.Receiver.AllWalletsDeleted:
			if win.currentPage != nil {
				win.currentPage.WillDisappear()
			}

			win.currentPage = nil
		}
	}
}

// displayWindow is called when a FrameEvent is received by the active window.
// Since user actions such as button clicks also trigger FrameEvents, this
// method first checks for pending user actions before displaying the UI
// elements. This ensures that the proper interface is displayed to the user
// based on their last performed action where applicable.
func (win *Window) displayWindow(evt system.FrameEvent) {
	// Set up the StartPage the first time a FrameEvent is received.
	if win.currentPage == nil {
		win.currentPage = page.NewStartPage(win.load)
		win.currentPage.WillAppear()
	}

	// A FrameEvent may be generated because of a user interaction
	// with the current page such as a button click. First handle
	// any such user interaction before rendering the page.
	win.currentPage.HandleUserInteractions()
	for _, modal := range win.modals {
		modal.Handle() // TODO: Just the top-most modal should do.
	}

	// Draw the window's UI components into an op.Ops.
	gtx := layout.NewContext(win.ops, evt)
	win.drawWindowUI(gtx)

	// Render the window's UI components on screen.
	evt.Frame(gtx.Ops)
}

// layout draws the window UI components into the provided graphical context,
// preparing the context for rendering on screen
func (win *Window) drawWindowUI(gtx C) {
	// Create a base view holder to hold all the following UI components
	// one on top the other. Components that do not take up the entire
	// window will be aligned to the top of the window.
	viewsHolder := layout.Stack{Alignment: layout.N}

	background := layout.Expanded(func(gtx C) D {
		return decredmaterial.Fill(gtx, win.load.Theme.Color.Gray4)
	})

	modals := layout.Stacked(func(gtx C) D {
		modals := win.modals
		if len(modals) == 0 {
			return layout.Dimensions{}
		}

		modalLayouts := make([]layout.StackChild, 0)
		for _, modal := range modals {
			widget := modal.Layout(gtx)
			l := layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return widget
			})
			modalLayouts = append(modalLayouts, l)
		}

		return layout.Stack{Alignment: layout.Center}.Layout(gtx, modalLayouts...)
	})

	viewsHolder.Layout(
		gtx,
		background,
		layout.Stacked(win.currentPage.Layout),
		modals,
		layout.Stacked(win.load.Toast.Layout),
	)
}
