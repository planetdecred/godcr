package wallet

import (
	"github.com/raedahgroup/dcrlibwallet"
)

// Response represents a discriminated union for wallet responses.
// Either Resp or Err must be nil.
type Response struct {
	Resp interface{}
	Err  error
}

// ResponseError wraps err in a Response
func ResponseError(err error) Response {
	return Response{
		Err: err,
	}
}

// ResponseResp wraps resp in a Response
func ResponseResp(resp interface{}) Response {
	return Response{
		Resp: resp,
	}
}

// MultiWalletInfo represents bulk information about the wallets returned by the wallet backend
type MultiWalletInfo struct {
	LoadedWallets   int
	TotalBalance    string
	Wallets         []InfoShort
	BestBlockHeight int32
	BestBlockTime   int64
	Synced          bool
	Syncing         bool
}

// InfoShort represents basic information about a wallet
type InfoShort struct {
	ID              int
	Name            string
	Balance         string
	Accounts        []Account
	BestBlockHeight int32
	BlockTimestamp  int64
	IsWaiting       bool
}

// Account represents information about a wallet's account
type Account struct {
	Number    string
	Name      string
	Spendable string
	Keys      struct {
		Internal, External, Imported string
	}
	HDPath         string
	TotalBalance   string
	CurrentAddress string
}

// AddedAccount is sent when the wallet is done adding an account
type AddedAccount struct {
	ID int32
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
