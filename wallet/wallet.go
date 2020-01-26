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
	multi     *dcrlibwallet.MultiWallet
	root, net string
	event.Duplex
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

// NewWallet creates a new wallet instance
func NewWallet(rootdir string, network string, duplex event.Duplex) *Wallet {
	wal := new(Wallet)
	wal.root = rootdir
	wal.net = network
	wal.Duplex = duplex
	return wal
}

// loadWallets loads the wallets for network in the root directory and returns
// an error if it occurs.
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
