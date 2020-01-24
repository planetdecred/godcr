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
	multi   *dcrlibwallet.MultiWallet
	root    string // root directory for all wallet data
	retwork string
	event.Duplex
}

// New loads a new wallet instance
func New(rootdir string, network string) (*Wallet, event.DuplexBase, error) {
	wal := new(Wallet)
	duplexB := event.NewDuplexBase()

	err := wal.loadWallets(rootdir, network)
	if err != nil {
		return nil, duplexB, err
	}

	wal.Duplex = duplexB.Duplex()
	return wal, duplexB, nil
}

// loadWallets loads the wallets for network in the rootdir and returns the wallet,
// the number of wallets loaded or an error if it occurs.
func (wal *Wallet) loadWallets(root string, net string) error {
	if root == "" || net == "" { // This should really be handled by dcrlibwallet
		return fmt.Errorf(`root directory or network cannot be ""`)
	}
	multiWal, err := dcrlibwallet.NewMultiWallet(root, "bdb", net)
	if err != nil {
		return err
	}
	wal.multi = multiWal
	return err
}
