package uiwallet

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
func (w *Wallet) updateStates(update interface{}) {
	switch e := update.(type) {
	case wallet.MultiWalletInfo:
		if w.walletInfo.LoadedWallets == 0 && e.LoadedWallets > 0 {
			w.changePage(PageOverview)
		}
		*w.walletInfo = e
		w.states.loading = false

		if e.LoadedWallets > 0 {
			// set wallets and accounts tab when wallet info is updated
			go func() {
				wallets := make([]decredmaterial.TabItem, len(e.Wallets))
				for i := range e.Wallets {
					wallets[i] = decredmaterial.TabItem{
						Title: e.Wallets[i].Name,
					}
				}
				w.walletTabs.SetTabs(wallets)

				accounts := make([]decredmaterial.TabItem, len(e.Wallets[w.selected].Accounts))
				for i, account := range e.Wallets[w.selected].Accounts {
					if account.Name == "imported" {
						continue
					}
					accounts[i] = decredmaterial.TabItem{
						Title: e.Wallets[w.selected].Accounts[i].Name,
					}
				}
				w.accountTabs.SetTabs(accounts)
			}()
		}
		return
	case *wallet.Transactions:
		w.states.loading = false
		w.walletTransactions = e
		return
	case *wallet.Transaction:
		w.walletTransaction = e
		return
	case *wallet.UnspentOutputs:
		w.walletUnspentOutputs = e
	case *wallet.Tickets:
		w.states.loading = false
		w.walletTickets = e
		return
	case *wallet.VSPInfo:
		w.states.loading = false
		w.vspInfo.List = append(w.vspInfo.List, *e)
		w.refresh()
		return
	case *wallet.VSP:
		w.vspInfo = e
		w.refresh()
		return
	case *wallet.Proposals:
		w.states.loading = false
		w.proposals = e
		return
	case wallet.CreatedSeed:
		w.notifyOnSuccess("Wallet created")
		w.changePage(PageWallet)
	case wallet.Renamed:
		w.notifyOnSuccess("Wallet renamed")
	case wallet.Restored:
		w.changePage(PageWallet)
		w.states.creating = false
		// w.window.Invalidate()
	case wallet.DeletedWallet:
		w.selected = 0
		w.changePage(PageWallet)
		w.notifyOnSuccess("Wallet removed")
	case wallet.AddedAccount:
		w.notifyOnSuccess("Account created")
	case wallet.UpdatedAccount:
		w.notifyOnSuccess("Account renamed")
	case *wallet.Signature:
		w.notifyOnSuccess("Message signed")
		w.signatureResult = update.(*wallet.Signature)
	case *dcrlibwallet.TxAuthor:
		txAuthor := update.(*dcrlibwallet.TxAuthor)
		w.txAuthor = *txAuthor
	case *wallet.Broadcast:
		broadcastResult := update.(*wallet.Broadcast)
		w.broadcastResult = *broadcastResult
	case *wallet.ChangePassword:
		w.notifyOnSuccess("Spending password changed")
	case *wallet.StartupPassphrase:
		w.notifyOnSuccess(update.(*wallet.StartupPassphrase).Msg)
	case wallet.OpenWallet:
		go func() {
			w.modal <- &modalLoad{}
		}()
	case wallet.SetupAccountMixer:
		w.notifyOnSuccess("Mixer setup completed")
	case *wallet.TicketPurchase:
		w.notifyOnSuccess("Ticket(s) purchased, attempting to pay fee")
	}

	w.states.loading = true
	w.wallet.GetMultiWalletInfo()
	w.wallet.GetAllTransactions(0, 0, 0)
	w.wallet.GetAllTickets()
	w.wallet.GetAllProposals()
	// w.window.Invalidate()
	log.Debugf("Updated with multiwallet info: %+v\n and window state %+v", w.walletInfo, w.states)
}

func (w *Wallet) notifyOnSuccess(text string) {
	go func() {
		w.toast <- &toast{
			text:    text,
			success: true,
		}
	}()

	go func() {
		w.modal <- &modalLoad{}
	}()
}
