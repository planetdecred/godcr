package ui

import (
	"errors"
	"fmt"
	"image"
	"strings"
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

// Window represents the app window (and UI in general). There should only be one.
// Window uses an internal state of booleans to determine what the window is currently displaying.
type Window struct {
	window *app.Window
	theme  *decredmaterial.Theme
	ops    *op.Ops

	wallet             *wallet.Wallet
	walletInfo         *wallet.MultiWalletInfo
	walletSyncStatus   *wallet.SyncStatus
	walletTransactions *wallet.Transactions
	walletTransaction  *wallet.Transaction
	walletAccount      *wallet.Account
	walletTickets      *wallet.Tickets

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

	internalLog chan string
}

type WriteClipboard struct {
	Text string
}

// CreateWindow creates and initializes a new window with start
// as the first page displayed.
// Should never be called more than once as it calls
// app.NewWindow() which does not support being called more
// than once.
func CreateWindow(wal *wallet.Wallet, decredIcons map[string]image.Image, collection []text.FontFace, internalLog chan string) (*Window, error) {
	win := new(Window)
	var netType string
	if strings.Contains(wal.Net, "testnet") {
		netType = "testnet"
	} else {
		netType = wal.Net
	}
	win.window = app.NewWindow(app.Title(fmt.Sprintf("%s (%s)", "godcr", netType)))
	theme := decredmaterial.NewTheme(collection, decredIcons)
	if theme == nil {
		return nil, errors.New("Unexpected error while loading theme")
	}
	win.theme = theme
	win.ops = &op.Ops{}

	win.walletInfo = new(wallet.MultiWalletInfo)
	win.walletSyncStatus = new(wallet.SyncStatus)
	win.walletTransactions = new(wallet.Transactions)
	win.walletUnspentOutputs = new(wallet.UnspentOutputs)
	win.walletAcctMixerStatus = make(chan *wallet.AccountMixer)

	win.wallet = wal
	win.states.loading = true
	win.current = PageOverview
	win.keyEvents = make(chan *key.Event)
	win.clipboard = make(chan interface{}, 2)
	win.toast = make(chan *toast)
	win.modal = make(chan *modalLoad)
	win.theme.ReadClipboard = win.clipboard

	win.walletTabs, win.accountTabs = decredmaterial.NewTabs(win.theme), decredmaterial.NewTabs(win.theme)
	win.walletTabs.Position, win.accountTabs.Position = decredmaterial.Top, decredmaterial.Top
	win.walletTabs.Separator, win.walletTabs.Separator = false, false
	win.accountTabs.SetTitle(win.theme.Label(values.TextSize18, "Accounts:"))
	win.walletTabs.SetTabs([]decredmaterial.TabItem{})
	win.accountTabs.SetTabs([]decredmaterial.TabItem{})

	win.internalLog = internalLog

	win.addPages(decredIcons)

	return win, nil
}

func (win *Window) changePage(page string) {
	win.current = page
	win.refresh()
}

func (win *Window) refresh() {
	win.window.Invalidate()
}

func (win *Window) setReturnPage(from string) {
	win.previous = from
	win.window.Invalidate()
}

func (win *Window) unloaded() {
	lbl := win.theme.H3("Multiwallet not loaded\nIs another instance open?")
	for {
		e := <-win.window.Events()
		switch evt := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			gtx := layout.NewContext(win.ops, evt)
			lbl.Layout(gtx)
			evt.Frame(win.ops)
		}
	}
}

// Loop runs main event handling and page rendering loop
func (win *Window) Loop(shutdown chan int) {
	for {
		select {
		case e := <-win.wallet.Send:
			if e.Err != nil {
				err := e.Err.Error()
				log.Error("Wallet Error: " + err)
				if err == dcrlibwallet.ErrWalletDatabaseInUse {
					close(shutdown)
					win.unloaded()
					return
				}
				win.err = err
				if win.states.loading {
					log.Warn("Attemping to get multiwallet info")
					win.wallet.GetMultiWalletInfo()
				}

				win.window.Invalidate()
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
					win.wallet.GetAllTickets()
					win.wallet.GetMultiWalletInfo()
					win.updateSyncProgress(update.BlockInfo)
				}
			case wallet.BlockConfirmed:
				win.updateSyncProgress(update.ConfirmedTxn)
			case wallet.AccountMixerStarted, wallet.AccountMixerEnded:
				go func() {
					win.walletAcctMixerStatus <- &update.AcctMixerInfo
				}()
			}
			win.window.Invalidate()
		case e := <-win.window.Events():
			switch evt := e.(type) {
			case system.DestroyEvent:
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
				s := win.states
				if win.walletInfo.LoadedWallets == 0 {
					win.changePage(PageCreateRestore)
				}

				if s.loading {
					win.Loading(gtx)
				} else {
					win.theme.Background(gtx, win.pages[win.current])
				}

				evt.Frame(win.ops)
			case key.Event:
				go func() {
					win.keyEvents <- &evt
				}()
			case clipboard.Event:
				go func() {
					win.theme.Clipboard <- evt.Text
				}()
			case nil:
				// Ignore
			default:
				log.Tracef("Unhandled window event %+v\n", e)
			}
		case e := <-win.clipboard:
			switch c := e.(type) {
			case decredmaterial.ReadClipboard:
				win.window.ReadClipboard()
			case WriteClipboard:
				go func() {
					win.window.WriteClipboard(c.Text)
					win.toast <- &toast{
						text:    "copied",
						success: true,
					}
				}()
			}
		}
	}
}
