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
}

// updateStates changes the wallet state based on the received update
func (win *Window) updateStates(update interface{}) {
	switch e := update.(type) {
	case wallet.MultiWalletInfo:
		if win.walletInfo.LoadedWallets == 0 && e.LoadedWallets > 0 {
			win.changePage(PageOverview)
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
	case *wallet.Tickets:
		win.walletTickets = e
		return
	case wallet.CreatedSeed:
		win.notifyOnSuccess("Wallet created")
		win.changePage(PageWallet)
	case wallet.Renamed:
		win.notifyOnSuccess("Wallet renamed")
	case wallet.Restored:
		win.changePage(PageWallet)
		win.states.creating = false
		win.window.Invalidate()
	case wallet.DeletedWallet:
		win.selected = 0
		win.changePage(PageWallet)
		win.notifyOnSuccess("Wallet removed")
	case wallet.AddedAccount:
		win.notifyOnSuccess("Account created")
	case wallet.UpdatedAccount:
		win.notifyOnSuccess("Account renamed")
	case *wallet.Signature:
		win.notifyOnSuccess("Message signed")
		win.signatureResult = update.(*wallet.Signature)
	case *dcrlibwallet.TxAuthor:
		txAuthor := update.(*dcrlibwallet.TxAuthor)
		win.txAuthor = *txAuthor
	case *wallet.Broadcast:
		broadcastResult := update.(*wallet.Broadcast)
		win.broadcastResult = *broadcastResult
	case *wallet.ChangePassword:
		win.notifyOnSuccess("Spending password changed")
	case *wallet.StartupPassphrase:
		win.notifyOnSuccess(update.(*wallet.StartupPassphrase).Msg)
	case wallet.OpenWallet:
		go func() {
			win.modal <- &modalLoad{}
		}()
	case wallet.SetupAccountMixer:
		win.notifyOnSuccess("Mixer setup completed")
	}

	win.states.loading = true
	win.wallet.GetMultiWalletInfo()
	win.wallet.GetAllTransactions(0, 0, 0)
	win.wallet.GetAllTickets()

	log.Debugf("Updated with multiwallet info: %+v\n and window state %+v", win.walletInfo, win.states)
}

func (win *Window) notifyOnSuccess(text string) {
	go func() {
		win.toast <- &toast{
			text:    text,
			success: true,
		}
	}()

	go func() {
		win.modal <- &modalLoad{}
	}()
}
