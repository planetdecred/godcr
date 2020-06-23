package ui

import (
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/wallet"
)

// states represents a combination of booleans that determine what the wallet is displaying.
type states struct {
	loading  bool // true if the window is in the middle of an operation that cannot be stopped
	creating bool // true if a wallet is being created or restored
	deleted  bool // true if a wallet is has being deleted
}

// updateStates changes the wallet state based on the received update
func (win *Window) updateStates(update interface{}) {
	switch e := update.(type) {
	case wallet.MultiWalletInfo:
		if win.walletInfo.LoadedWallets == 0 && e.LoadedWallets > 0 {
			win.current = PageOverview
		}
		*win.walletInfo = e
		win.states.loading = false
		return
	case *wallet.Transactions:
		win.walletTransactions = e
		return
	case *wallet.Transaction:
		win.walletTransaction = e
		return
	case wallet.CreatedSeed:
		win.current = PageWallet
		win.states.creating = false
		win.window.Invalidate()
	case wallet.Restored:
		win.current = PageWallet
		win.states.creating = false
		win.window.Invalidate()
	case wallet.DeletedWallet:
		win.selected = 0
		win.current = PageWallet
		win.states.deleted = true
		win.window.Invalidate()
	case wallet.AddedAccount:
		win.current = PageWallet
		win.states.creating = false
		win.window.Invalidate()
	case *wallet.Signature:
		win.signatureResult = update.(*wallet.Signature)
	case *dcrlibwallet.TxAuthor:
		txAuthor := update.(*dcrlibwallet.TxAuthor)
		win.txAuthor = *txAuthor
	case *wallet.Broadcast:
		broadcastResult := update.(*wallet.Broadcast)
		win.broadcastResult = *broadcastResult
	}
	win.states.loading = true
	win.wallet.GetMultiWalletInfo()
	win.wallet.GetAllTransactions(0, 0, 0)

	log.Debugf("Updated with multiwallet info: %+v\n and window state %+v", win.walletInfo, win.states)
}
