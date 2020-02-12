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

// InfoShort represents basic information about a wallet
type InfoShort struct {
	ID              int
	Name     string
	Balance  int64
	Accounts []int32
	BestBlockHeight int32
	BlockTimestamp  int64
	IsWaiting       bool
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
	Txs    [][]dcrlibwallet.Transaction
	Recent []dcrlibwallet.Transaction
}

// SyncStatus is sent when a wallet progress event is triggered.
type SyncStatus struct {
	Progress                 int32
	HeadersFetchProgress     int32
	HeadersToFetch           int32
	RescanHeadersProgress    int32
	AddressDiscoveryProgress int32
	RemainingTime            int64
	ConnectedPeers           int32
	Steps                    int32
	TotalSteps               int32
	CurrentBlockHeight       int32
}
