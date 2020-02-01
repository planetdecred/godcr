// Package wallet provides functions and types for interacting
// with the dcrlibwallet backend.
package wallet

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
)

const syncID = "godcr"

// Wallet represents the wallet back end of the app
type Wallet struct {
	multi     *dcrlibwallet.MultiWallet
	root, net string
	Send      chan interface{}
}

// InternalWalletError represents errors generated during the handling of the multiwallet
// and connected wallets
type InternalWalletError struct {
	Message         string
	AffectedWallets []int
}

func (err *InternalWalletError) Error() string {
	return err.Message
}

// NewWallet initializies an new Wallet instance.
// The Wallet is not loaded until LoadWallets is called.
func NewWallet(root string, net string, send chan interface{}) (*Wallet, error) {
	if root == "" || net == "" { // This should really be handled by dcrlibwallet
		return nil, fmt.Errorf(`root directory or network cannot be ""`)
	}
	wal := &Wallet{
		root: root,
		net:  net,
		Send: send,
	}

	return wal, nil
}

// LoadWallets loads the wallets for network in the root directory.
// It adds a SyncProgressListener to the multiwallet and opens the wallets if no
// startup passphrase was set.
// It is non-blocking and sends its result or any erro to wal.Send.
func (wal *Wallet) LoadWallets() {
	go func(send chan<- interface{}, wal *Wallet) {
		multiWal, err := dcrlibwallet.NewMultiWallet(wal.root, "bdb", wal.net)
		if err != nil {
			send <- err
			return
		}
		wal.multi = multiWal
		wal.multi.AddSyncProgressListener(&progressListener{
			Send: wal.Send,
		}, syncID)
		startupPassSet := wal.multi.IsStartupSecuritySet()
		if !startupPassSet {
			err = wal.multi.OpenWallets(nil)
			if err != nil {
				send <- err
				return
			}
		}

		send <- &LoadedWallets{
			Count:              wal.multi.LoadedWalletsCount(),
			StartUpSecuritySet: startupPassSet,
		}
	}(wal.Send, wal)
}

// wallets returns an up-to-date slice of all opened wallets
func (wal *Wallet) wallets() ([]*dcrlibwallet.Wallet, error) {
	if wal.multi == nil {
		return nil, &InternalWalletError{
			Message: "No MultiWallet loaded",
		}
	}

	wallets := make([]*dcrlibwallet.Wallet, len(wal.multi.OpenedWalletIDsRaw()))

	for i, j := range wal.multi.OpenedWalletIDsRaw() {
		w := wal.multi.WalletWithID(j)
		if w == nil {
			return nil, &InternalWalletError{
				Message: "Invalid Wallet ID",
			}
		}
		wallets[i] = w
	}
	return wallets, nil
}

// Shutdown shutsdown the multiwallet
func (wal *Wallet) Shutdown() {
	wal.multi.Shutdown()
}
