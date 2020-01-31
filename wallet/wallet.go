// Package wallet provides functions and types for interacting
// with the dcrlibwallet backend.
package wallet

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
)

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

// NewWallet initializies an new wallet instance
func NewWallet(root string, net string, send chan interface{}) (*Wallet, error) {
	wal := &Wallet{
		root: root,
		net:  net,
		Send: send,
	}
	if root == "" || net == "" { // This should really be handled by dcrlibwallet
		return nil, fmt.Errorf(`root directory or network cannot be ""`)
	}

	return wal, nil
}

// LoadWallets loads the wallets for network in the root directory and returns
// an error if it occurs.
func (wal *Wallet) LoadWallets() {
	go func(send chan<- interface{}, wal *Wallet) {
		multiWal, err := dcrlibwallet.NewMultiWallet(wal.root, "bdb", wal.net)
		if err != nil {
			send <- err
			return
		}
		wal.multi = multiWal
		send <- &LoadedWallets{
			Count: wal.multi.LoadedWalletsCount(),
		}
	}(wal.Send, wal)
}

// wallets returns an up-to-date slice of loaded wallets
func (wal *Wallet) wallets() ([]*dcrlibwallet.Wallet, error) {
	if wal.multi == nil {
		return nil, &InternalWalletError{
			Message: "No MultiWallet loaded",
		}
	}

	count := int(wal.multi.LoadedWalletsCount())
	wallets := make([]*dcrlibwallet.Wallet, count)

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
