package wallet

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
)

// LoadWallets loads the wallets and returns the wallet
func LoadWallets(rootdir, network string) (*dcrlibwallet.MultiWallet, bool, error) {
	// fmt.Printf("%s ; %s \n", rootdir, network)
	// return nil, false, nil
	if rootdir == "" || network == "" {
		return nil, false, fmt.Errorf("Dir cannot be empty")
	}
	wal, err := dcrlibwallet.NewMultiWallet(rootdir, "", network)
	if err != nil {
		return nil, false, err
	}
	return wal, wal.LoadedWalletsCount() > 0, err
}
