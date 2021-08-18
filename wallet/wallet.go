// Package wallet provides functions and types for interacting
// with the dcrlibwallet backend.
package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/planetdecred/dcrlibwallet"
)

const syncID = "godcr"

// Wallet represents the wallet back end of the app
type Wallet struct {
	multi              *dcrlibwallet.MultiWallet
	Root, Net          string
	Send               chan Response
	Sync               chan SyncStatusUpdate
	confirms           int32
	OverallBlockHeight int32
}

// NewWallet initializies an new Wallet instance.
// The Wallet is not loaded until LoadWallets is called.
func NewWallet(root string, net string, send chan Response, confirms int32) (*Wallet, error) {
	if root == "" || net == "" { // This should really be handled by dcrlibwallet
		return nil, fmt.Errorf(`root directory or network cannot be ""`)
	}
	wal := &Wallet{
		Root:     root,
		Net:      net,
		Sync:     make(chan SyncStatusUpdate, 2),
		Send:     send,
		confirms: confirms,
	}

	return wal, nil
}

func (wal *Wallet) InitMultiWallet() error {
	politeiaHost := dcrlibwallet.PoliteiaMainnetHost
	if wal.Net == dcrlibwallet.Testnet3 {
		politeiaHost = dcrlibwallet.PoliteiaTestnetHost
	}
	multiWal, err := dcrlibwallet.NewMultiWallet(wal.Root, "bdb", wal.Net, politeiaHost)
	if err != nil {
		return err
	}

	wal.multi = multiWal
	return nil
}

func (wal *Wallet) SetupListeners() {
	resp := Response{
		Resp: LoadedWallets{},
	}
	l := &listener{
		Send: wal.Sync,
	}
	err := wal.multi.AddSyncProgressListener(l, syncID)
	if err != nil {
		resp.Err = err
		wal.Send <- resp
		return
	}

	err = wal.multi.AddTxAndBlockNotificationListener(l, syncID)
	if err != nil {
		resp.Err = err
		wal.Send <- resp
		return
	}

	wal.multi.AddAccountMixerNotificationListener(l, syncID)

	wal.multi.Politeia.AddNotificationListener(l, syncID)

	startupPassSet := wal.multi.IsStartupSecuritySet()

	resp.Resp = LoadedWallets{
		Count:              wal.multi.LoadedWalletsCount(),
		StartUpSecuritySet: startupPassSet,
	}
	wal.Send <- resp
}

// wallets returns an up-to-date map of all opened wallets
func (wal *Wallet) wallets() ([]dcrlibwallet.Wallet, error) {
	if wal.multi == nil {
		return nil, MultiWalletError{
			Message: "No MultiWallet loaded",
		}
	}

	wallets := []dcrlibwallet.Wallet{}
	for _, j := range wal.multi.OpenedWalletIDsRaw() {
		w := wal.multi.WalletWithID(j)
		if w == nil {
			return nil, InternalWalletError{
				Message:  "Invalid Wallet ID",
				Err:      ErrIDNotExist,
				Affected: []int{j},
			}
		}
		wallets = append(wallets, *w)
	}

	// sort wallet by ids
	if len(wallets) > 0 {
		sort.SliceStable(wallets, func(i, j int) bool {
			return wallets[i].ID < wallets[j].ID
		})
	}

	return wallets, nil
}

func (wal *Wallet) hdPrefix() string {
	switch wal.Net {
	case "testnet3": // should use a constant
		return dcrlibwallet.TestnetHDPath
	case "mainnet":
		return dcrlibwallet.MainnetHDPath
	default:
		return ""
	}
}

// Shutdown shutsdown the multiwallet
func (wal *Wallet) Shutdown() {
	if wal.multi != nil {
		wal.multi.Shutdown()
	}
}

// GetBlockExplorerURL accept transaction hash,
// return the block explorer URL with respect to the network
func (wal *Wallet) GetBlockExplorerURL(txnHash string) string {
	switch wal.Net {
	case "testnet3": // should use a constant
		return "https://testnet.dcrdata.org/tx/" + txnHash
	case "mainnet":
		return "https://explorer.dcrdata.org/tx/" + txnHash
	default:
		return ""
	}
}

//GetUSDExchangeValues gets the exchange rate of DCR - USDT from a specified endpoint
func (wal *Wallet) GetUSDExchangeValues(target interface{}) error {
	url := "https://api.bittrex.com/v3/markets/DCR-USDT/ticker"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(target)
	return nil
}
