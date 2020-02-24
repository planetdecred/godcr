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
	Send      chan Response
	confirms  int32
}

// NewWallet initializies an new Wallet instance.
// The Wallet is not loaded until LoadWallets is called.
func NewWallet(root string, net string, send chan Response, confirms int32) (*Wallet, error) {
	if root == "" || net == "" { // This should really be handled by dcrlibwallet
		return nil, fmt.Errorf(`root directory or network cannot be ""`)
	}
	wal := &Wallet{
		root:     root,
		net:      net,
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
	go func(send chan<- Response, wal *Wallet) {
		var resp Response
		multiWal, err := dcrlibwallet.NewMultiWallet(wal.root, "bdb", wal.net)
		if err != nil {
			resp.Err = err
			send <- resp
			return
		}

		wal.multi = multiWal
		l := &listener{
			Send: wal.Send,
		}
		err = wal.multi.AddSyncProgressListener(l, syncID)
		if err != nil {
			resp.Err = err
			send <- resp
			return
		}

		err = wal.multi.AddTxAndBlockNotificationListener(l, syncID)
		if err != nil {
			resp.Err = err
			send <- resp
			return
		}

		startupPassSet := wal.multi.IsStartupSecuritySet()
		if !startupPassSet {
			err = wal.multi.OpenWallets(nil)
			if err != nil {
				resp.Err = err
				send <- resp
				return
			}
		}

		resp.Resp = LoadedWallets{
			Count:              wal.multi.LoadedWalletsCount(),
			StartUpSecuritySet: startupPassSet,
		}
		send <- resp
	}(wal.Send, wal)
}

// wallets returns an up-to-date map of all opened wallets
func (wal *Wallet) wallets() (map[int]*dcrlibwallet.Wallet, error) {
	if wal.multi == nil {
		return nil, MultiWalletError{
			Message: "No MultiWallet loaded",
		}
	}

	wallets := make(map[int]*dcrlibwallet.Wallet, len(wal.multi.OpenedWalletIDsRaw()))

	for _, j := range wal.multi.OpenedWalletIDsRaw() {
		w := wal.multi.WalletWithID(j)
		if w == nil {
			return nil, InternalWalletError{
				Message:  "Invalid Wallet ID",
				Err:      ErrIDNotExist,
				Affected: []int{j},
			}
		}
		wallets[j] = w
	}
	return wallets, nil
}

// Shutdown shutsdown the multiwallet
func (wal *Wallet) Shutdown() {
	wal.multi.Shutdown()
}
