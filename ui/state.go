package ui

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/wallet"
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

		if e.LoadedWallets > 0 {
			// set wallets and accounts tab when wallet info is updated
			go func() {
				wallets := make([]decredmaterial.TabItem, len(e.Wallets))
				for i := range e.Wallets {
					wallets[i] = decredmaterial.TabItem{
						Title: e.Wallets[i].Name,
					}
				}
				win.walletTabs.SetTabs(wallets)

				accounts := make([]decredmaterial.TabItem, len(e.Wallets[win.selected].Accounts))
				for i, account := range e.Wallets[win.selected].Accounts {
					if account.Name == "imported" {
						continue
					}
					accounts[i] = decredmaterial.TabItem{
						Title: e.Wallets[win.selected].Accounts[i].Name,
					}
				}
				win.accountTabs.SetTabs(accounts)
			}()
		}
		return
	case *wallet.Transactions:
		win.walletTransactions = e
		return
	case *wallet.Transaction:
		win.walletTransaction = e
		return
	case *wallet.UnspentOutputs:
		win.walletUnspentOutputs = e
		return
	case wallet.CreatedSeed:
		win.current = PageWallet
		win.states.creating = false
		go func() {
			win.toast <- &toast{
				text:    "Wallet created",
				success: true,
			}
		}()

		go func() {
			win.modal <- &modalLoad{}
		}()
		win.window.Invalidate()
	case wallet.Renamed:
		go func() {
			win.toast <- &toast{
				text:    "Wallet renamed",
				success: true,
			}
		}()

		go func() {
			win.modal <- &modalLoad{}
		}()
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
	case wallet.UpdatedAccount:
		go func() {
			win.toast <- &toast{
				text:    "Account renamed",
				success: true,
			}
		}()

		go func() {
			win.modal <- &modalLoad{}
		}()
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
