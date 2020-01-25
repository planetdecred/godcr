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
	wallets []*dcrlibwallet.Wallet
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

func (wal *Wallet) reloadWallets() error {
	if wal.multi == nil {
		return &InternalWalletError{
			Message: "No MultiWallet loaded",
		}
	}

	count := int(wal.multi.LoadedWalletsCount())
	wallets := make([]*dcrlibwallet.Wallet, count)

	for i, j := range wal.multi.OpenedWalletIDsRaw() {
		w := wal.multi.WalletWithID(j)
		if w == nil {
			return &InternalWalletError{
				Message: "Invalid Wallet ID",
			}
		}
		wallets[i] = w
	}

	wal.wallets = wallets
	return nil
}
