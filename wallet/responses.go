package wallet

import (
	"github.com/raedahgroup/dcrlibwallet"
)

// MultiWalletInfo represents bulk information about the wallets returned by the wallet backend
type MultiWalletInfo struct {
	LoadedWallets   int
	TotalBalance    int64
	Balances        []int64
	BestBlockHeight int32
	BestBlockTime   int64
	Synced          bool
}

// LoadedWallets is sent when then the Wallet is done loading wallets
type LoadedWallets struct {
	Count int32
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
