// Package wallet provides functions and types for interacting
// with the dcrlibwallet backend.
package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/planetdecred/dcrlibwallet"
)

const (
	syncID    = "godcr"
	DevBuild  = "dev"
	ProdBuild = "prod"
)

// Wallet represents the wallet back end of the app
type Wallet struct {
	multi     *dcrlibwallet.MultiWallet
	Root, Net string
	buildDate time.Time
	version   string
	logFile   string
}

// NewWallet initializies an new Wallet instance.
// The Wallet is not loaded until LoadWallets is called.
func NewWallet(root, net, version, logFile string, buildDate time.Time) (*Wallet, error) {
	if root == "" || net == "" { // This should really be handled by dcrlibwallet
		return nil, fmt.Errorf(`root directory or network cannot be ""`)
	}

	wal := &Wallet{
		Root:      root,
		Net:       net,
		buildDate: buildDate,
		version:   version,
		logFile:   logFile,
	}

	return wal, nil
}

func (wal *Wallet) BuildDate() time.Time {
	return wal.buildDate
}

func (wal *Wallet) Version() string {
	return wal.version
}

func (wal *Wallet) LogFile() string {
	return wal.logFile
}

func (wal *Wallet) StartupTime() time.Time {
	return time.Time{}
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

func (wal *Wallet) hdPrefix() string {
	switch wal.Net {
	case dcrlibwallet.Testnet3:
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
	case dcrlibwallet.Testnet3:
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
