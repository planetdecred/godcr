package ui

import (
	"context"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

// Window represents the app window (and UI in general). There should only be one.
// Window uses an internal state of booleans to determine what the window is currently displaying.
type Window struct {
	appCtx     context.Context
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

	keyEvents             chan *key.Event
	sysDestroyWithSync    bool
	walletAcctMixerStatus chan *wallet.AccountMixer
	internalLog           chan string

	dexc *dexc.Dexc
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
func CreateWindow(wal *wallet.Wallet, dc *dexc.Dexc, appCtx context.Context) (*Window, *app.Window, error) {
	win := new(Window)
	var netType string
	if wal.Net == "testnet3" {
		netType = "testnet"
	} else {
		netType = wal.Net
	}
	appWindow := app.NewWindow(app.Size(values.AppWidth, values.AppHeight), app.Title(values.StringF(values.StrAppTitle, netType)))
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

	win.dexc = dc

	win.states.loading = false

	win.keyEvents = make(chan *key.Event)

	l, err := win.NewLoad(appCtx)
	if err != nil {
		return nil, nil, err
	}

	win.load = l

	return win, appWindow, nil
}

func (win *Window) NewLoad(appCtx context.Context) (*load.Load, error) {
	l, err := load.NewLoad()
	if err != nil {
		return nil, err
	}

	l.AppCtx = appCtx
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
	}

	l.SelectedWallet = &win.selected
	l.RefreshWindow = win.refreshWindow
	l.ShowModal = win.showModal
	l.DismissModal = win.dismissModal
	l.PopWindowPage = win.popPage
	l.ChangeWindowPage = win.changePage

	l.DL = &load.DexcLoad{
		Core: win.dexc.Core,
		Dexc: win.dexc,
	}
	return l, nil
}

func (win *Window) Start() {
	if win.currentPage == nil {
		sp := page.NewStartPage(win.load)
		sp.OnResume()
		win.currentPage = sp
	}
}

func (win *Window) changePage(page load.Page, keepBackStack bool) {
	if win.currentPage != nil && keepBackStack {
		win.currentPage.OnClose()
		win.pageBackStack = append(win.pageBackStack, win.currentPage)
	}

	win.currentPage = page
	win.refreshWindow()
}

// popPage goes back to the previous page
// returns true if page was popped.
func (win *Window) popPage() bool {
	if len(win.pageBackStack) > 0 {
		// get and remove last page
		previousPage := win.pageBackStack[len(win.pageBackStack)-1]
		win.pageBackStack = win.pageBackStack[:len(win.pageBackStack)-1]

		win.currentPage.OnClose()

		previousPage.OnResume()
		win.currentPage = previousPage
		win.refreshWindow()

		return true
	}

	return false
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

func (win *Window) layoutPage(gtx C, page load.Page) {
	layout.Stack{
		Alignment: layout.N,
	}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return decredmaterial.Fill(gtx, win.load.Theme.Color.LightGray)
		}),
		layout.Stacked(func(gtx C) D {
			page.Handle()
			return page.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			modals := win.modals

			if len(modals) > 0 {
				modalLayouts := make([]layout.StackChild, 0)
				for _, modal := range modals {
					modal.Handle()
					widget := modal.Layout(gtx)
					l := layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return widget
					})
					modalLayouts = append(modalLayouts, l)
				}

				return layout.Stack{Alignment: layout.Center}.Layout(gtx, modalLayouts...)
			}

			return layout.Dimensions{}
		}),
		layout.Stacked(func(gtx C) D {
			return win.load.Toast.Layout(gtx)
		}),
	)
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
					win.unloaded(w)
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

			win.updateStates(e.Resp)
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
					win.Start()
				}

			case system.DestroyEvent:
				if win.currentPage != nil {
					win.currentPage.OnClose()
				}
				if win.walletInfo.Syncing || win.walletInfo.Synced {
					win.sysDestroyWithSync = true
					win.wallet.CancelSync()
				} else {
					close(shutdown)
				}
			case system.FrameEvent:
				gtx := layout.NewContext(win.ops, evt)
				ts := int64(time.Since(time.Unix(win.walletInfo.BestBlockTime, 0)).Seconds())
				win.walletInfo.LastSyncTime = wallet.SecondsToDays(ts)

				if win.currentPage != nil {
					win.layoutPage(gtx, win.currentPage)
				} else {
					win.Loading(gtx)
				}

				evt.Frame(gtx.Ops)
			case key.Event:
				go func() {
					win.keyEvents <- &evt
				}()
			case nil:
				// Ignore
			default:
				log.Tracef("Unhandled window event %+v\n", e)
			}
		}
	}
}
