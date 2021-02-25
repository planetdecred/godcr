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
	root, Net          string
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
		root:     root,
		Net:      net,
		Sync:     make(chan SyncStatusUpdate, 2),
		Send:     send,
		confirms: confirms,
	}

	return wal, nil
}

// LoadWallets loads the wallets for network in the root directory.
// It adds a SyncProgressListener to the multiwallet and opens the wallets if no
// startup passphrase was set.
// It is non-blocking and sends its result or any erro to wal.Send.
func (wal *Wallet) LoadWallets() {
	go func() {
		resp := Response{
			Resp: LoadedWallets{},
		}
		multiWal, err := dcrlibwallet.NewMultiWallet(wal.root, "bdb", wal.Net)
		if err != nil {
			resp.Err = err
			log.Error("Wallet not loaded. Is another process using the data directory?")
			wal.Send <- resp
			return
		}
		wal.multi = multiWal
		l := &listener{
			Send: wal.Sync,
		}
		err = wal.multi.AddSyncProgressListener(l, syncID)
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

		wal.multi.SetAccountMixerNotification(l)

		startupPassSet := wal.multi.IsStartupSecuritySet()
		if !startupPassSet {
			err = wal.multi.OpenWallets(nil)
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}
		}

		resp.Resp = LoadedWallets{
			Count:              wal.multi.LoadedWalletsCount(),
			StartUpSecuritySet: startupPassSet,
		}
		wal.Send <- resp
	}()
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
func (wal *Wallet) GetUSDExchangeValues(target interface{}) {
	url := "https://api.bittrex.com/v3/markets/DCR-USDT/ticker"
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(target)
}
