package wallet

import (
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
	TotalBalance    int64
	Wallets         []InfoShort
	BestBlockHeight int32
	BestBlockTime   int64
	Synced          bool
	Syncing         bool
}

// Account represents information about an account in a wallet
type Account struct {
	Number           int32
	Name             string
	TotalBalance     int64
	SpendableBalance int64
}

// InfoShort represents basic information about a wallet
type InfoShort struct {
	ID               int
	Name             string
	TotalBalance     int64
	SpendableBalance int64
	Accounts         []Account
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

// Transactions is sent in response to Wallet.GetAllTransactions
type Transactions struct {
	Txs [][]dcrlibwallet.Transaction
}

// TxHash is sent when the Wallet successfully broadcasts a transaction
type TxHash struct {
	Hash string
}
