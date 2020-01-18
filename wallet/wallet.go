// Package wallet provides functions and types for interacting
// with the dcrlibwallet backend.
package wallet

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
)

// LoadWallets loads the wallets for network in the rootdir and returns the wallet
// a boolean representing if there is at least one wallet loaded and an error if
// it occurs.
func LoadWallets(rootdir, network string) (*dcrlibwallet.MultiWallet, bool, error) {
	if rootdir == "" || network == "" {
		return nil, false, fmt.Errorf(`root directory or network cannot be ""`)
	}
	wal, err := dcrlibwallet.NewMultiWallet(rootdir, "", network)
	if err != nil {
		return nil, false, err
	}
	return wal, wal.LoadedWalletsCount() > 0, err
}
