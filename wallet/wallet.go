// Package wallet provides functions and types for interacting
// with the dcrlibwallet backend.
package wallet

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/event"
)

// Wallet represents the wallet back end of the app
type Wallet struct {
	multi       *dcrlibwallet.MultiWallet
	Root        string // root directory for all wallet data
	Network     string
	SendChan    chan event.Event // chan the wallet sends events to
	ReceiveChan chan event.Event // chan the wallet recieves commands from
}

// loadWallets loads the wallets for network in the rootdir and returns the wallet,
// the number of wallets loaded or an error if it occurs.
func (wal *Wallet) loadWallets() (int32, error) {
	if wal.Root == "" || wal.Network == "" { // This should really be handled by dcrlibwallet
		return 0, fmt.Errorf(`root directory or network cannot be ""`)
	}
	multiWal, err := dcrlibwallet.NewMultiWallet(wal.Root, "", wal.Network)
	if err != nil {
		return 0, err
	}
	wal.multi = multiWal
	return multiWal.LoadedWalletsCount(), err
}
