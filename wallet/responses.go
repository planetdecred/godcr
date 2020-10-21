package wallet

import (
	"github.com/planetdecred/dcrlibwallet"
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
	LastSyncTime    string
	Synced          bool
	Syncing         bool
}

// InfoShort represents basic information about a wallet
type InfoShort struct {
	ID               int
	Name             string
	Balance          string
	Accounts         []Account
	TotalBalance     string
	SpendableBalance int64
	BestBlockHeight  int32
	BlockTimestamp   int64
	DaysBehind       string
	Status           string
	IsWaiting        bool
	Seed             []byte
}

// Account represents information about a wallet's account
type Account struct {
	Number           int32
	Name             string
	SpendableBalance int64
	Keys             struct {
		Internal, External, Imported string
	}
	HDPath         string
	TotalBalance   string
	CurrentAddress string
}

// AddedAccount is sent when the wallet is done adding an account
type AddedAccount struct {
	ID               int32
	Number           int32
	Name             string
	TotalBalance     string
	CurrentAddress   string
	SpendableBalance int64
}

// UpdatedAccount is sent when the wallet is done updated an account
type UpdatedAccount struct {
	ID int32
}

// LoadedWallets is sent when then the Wallet is done loading wallets
type LoadedWallets struct {
	Count              int32
	StartUpSecuritySet bool
}

// Restored is sent when the Wallet is done restoring a wallet
type Restored struct{}

// StartUpPassphrase is sent when the startup passphrase is set
type StartupPassphrase struct {
	Msg string
}

// OpenWallet is sent when the startup passphrase is set
type OpenWallet struct{}

// Renamed is sent when the Wallet is done renaming a wallet
type Renamed struct {
	ID int
}

// CreatedSeed is sent when the Wallet is done creating a wallet
type CreatedSeed struct {
	Seed string
}

// DeletedWallet is sent when a wallet is deleted
type DeletedWallet struct {
	ID int
}

// ChangePassword is sent when the Wallet password is changed
type ChangePassword struct {
	ID int
}

// Transaction wraps the dcrlibwallet Transaction type and adds processed data
type Transaction struct {
	Txn           dcrlibwallet.Transaction
	Status        string
	Balance       string
	WalletName    string
	AccountName   string
	Confirmations int32
	DateTime      string
}

// Transactions is sent in response to Wallet.GetAllTransactions
type Transactions struct {
	Total  int
	Txs    map[int][]Transaction
	Recent []Transaction
}

// SyncStatus is sent when a wallet progress event is triggered.
type SyncStatus struct {
	Progress                 int32
	HeadersFetchProgress     int32
	HeadersToFetch           int32
	RescanHeadersProgress    int32
	AddressDiscoveryProgress int32
	RemainingTime            string
	ConnectedPeers           int32
	Steps                    int32
	TotalSteps               int32
	CurrentBlockHeight       int32
}

// Signature is sent in response to Wallet.SignMessage
type Signature struct {
	Signature string
	Err       error
}

// TxHash is sent when the Wallet successfully broadcasts a transaction
type TxHash struct {
	Hash string
}

type TxAuthor struct {
	TxAuthor int
}

// Broadcast is sent when the Wallet  broadcasts a transaction
type Broadcast struct {
	TxHash string
}

type UnspentOutput struct {
	UTXO     dcrlibwallet.UnspentOutput
	Amount   string
	DateTime string
}

// UnspentOutputs wraps the dcrlibwallet UTXO type and adds processed data
type UnspentOutputs struct {
	List []*UnspentOutput
}

// SetupAccountMixer is sent when finished setup the wallet account mixer
type SetupAccountMixer struct{}

type Ticket struct {
	Info     dcrlibwallet.TicketInfo
	Fee      string
	Amount   string
	DateTime string
}
type Tickets struct {
	Total int
	List  map[int][]Ticket
}
