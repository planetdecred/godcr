package load

import (
	"errors"
	"sort"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

// ErrIDNotExist is returned when a given ID does not exist
var ErrIDNotExist = errors.New("ID does not exist")

type WalletLoad struct {
	MultiWallet      *dcrlibwallet.MultiWallet
	TxAuthor         dcrlibwallet.TxAuthor
	SelectedProposal *dcrlibwallet.Proposal

	Proposals       *wallet.Proposals
	SyncStatus      *wallet.SyncStatus
	Transactions    *wallet.Transactions
	Transaction     *wallet.Transaction
	BroadcastResult wallet.Broadcast
	Tickets         **wallet.Tickets
	VspInfo         *wallet.VSP
	UnspentOutputs  *wallet.UnspentOutputs
	Wallet          *wallet.Wallet
	Account         *wallet.Account
	Info            *wallet.MultiWalletInfo

	SelectedWallet  *int
	SelectedAccount *int
}

func (wl *WalletLoad) SortedWalletList() []*dcrlibwallet.Wallet {
	wallets := wl.MultiWallet.AllWallets()

	sort.Slice(wallets, func(i, j int) bool {
		return wallets[i].ID < wallets[j].ID
	})

	return wallets
}

func (wl *WalletLoad) HDPrefix() string {
	switch wl.Wallet.Net {
	case "testnet3": // should use a constant
		return dcrlibwallet.TestnetHDPath
	case "mainnet":
		return dcrlibwallet.MainnetHDPath
	default:
		return ""
	}
}
