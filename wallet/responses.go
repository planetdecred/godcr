package wallet

import (
	"github.com/raedahgroup/dcrlibwallet"
)

// MultiWalletInfo represents bulk information about the wallets returned by the wallet backend
type MultiWalletInfo struct {
	LoadedWallets   int
	TotalBalance    int64
	BestBlockHeight int32
	BestBlockTime   int64
	Synced          bool
}

type LoadedWallets struct {
	Count int32
}

type Restored struct{}

type CreatedSeed struct {
	Seed string
}

type Transactions struct {
	Txs [][]dcrlibwallet.Transaction
}

// SyncEvent represents sync events
type SyncEvent struct {
	Event   string
	Payload interface{}
}
