package wallet

import (
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
)

// Response represents the structure of data that the Send channel receives
type Response struct {
	Resp interface{}
	Err  error
}

// MultiWalletInfo represents bulk information about the wallets returned by the wallet backend
type MultiWalletInfo struct {
	LoadedWallets   int
	TotalBalance    dcrutil.Amount
	Wallets         []InfoShort
	BestBlockHeight int32
	BestBlockTime   int64
	Synced          bool
	Syncing         bool
}

// InfoShort represents basic information about a wallet
type InfoShort struct {
	Name     string
	Balance  dcrutil.Amount
	Accounts []dcrlibwallet.Account
}

// Account represents infomation about a wallet's account
type Account struct {
	id      int32
	name    string
	Balance dcrutil.Amount
}

// LoadedWallets is sent when then the Wallet is done loading wallets
type LoadedWallets struct {
	Count              int32
	StartUpSecuritySet bool
}

// Restored is sent when the Wallet is done restoring a wallet
type Restored struct{}

// CreatedSeed is sent when the Wallet is done creating a wallet
type CreatedSeed struct {
	Seed string
}

// DeletedWallet is sent when a wallet is deleted
type DeletedWallet struct {
	ID int
}

// Transactions is sent in response to Wallet.GetAllTransactions
type Transactions struct {
	Txs [][]dcrlibwallet.Transaction
}
