package wallet

import (
	"errors"

	"github.com/raedahgroup/godcr-gio/event"
)

var (
	// ErrInvalidArguments is returned when a wallet command is send with invalid arguments.
	ErrInvalidArguments = errors.New("Invalid command arguments")

	// ErrNotFound is returned when a wallet command is given that does not exist or is not
	// implemented.
	ErrNotFound = errors.New("Command not found or not implemented")
)

var cmdMap = map[string]func(*Wallet, *event.ArgumentQueue) error{
	event.CreateCmd:  createCmd,
	event.RestoreCmd: restoreCmd,
	event.InfoCmd:    infoCmd,
}

func createCmd(wal *Wallet, arguments *event.ArgumentQueue) error {
	passphrase, err := arguments.PopString()
	if err != nil {
		return ErrInvalidArguments
	}
	passtype, err := arguments.PopInt()
	if err != nil {
		return ErrInvalidArguments
	}

	wall, err := wal.multi.CreateNewWallet(passphrase, int32(passtype))
	if err != nil {
		return err
	}
	wal.Send <- event.WalletResponse{
		Resp: event.CreatedResp,
		Results: &event.ArgumentQueue{
			Queue: []interface{}{wall.Seed},
		},
	}
	return nil
}

func restoreCmd(wal *Wallet, arguments *event.ArgumentQueue) error {
	seed, err := arguments.PopString()
	if err != nil {
		return ErrInvalidArguments
	}

	passphrase, err := arguments.PopString()
	if err != nil {
		return ErrInvalidArguments
	}
	passtype, err := arguments.PopInt()
	if err != nil {
		return ErrInvalidArguments
	}

	_, err = wal.multi.RestoreWallet(seed, passphrase, int32(passtype))
	if err != nil {
		return err
	}

	wal.Send <- event.WalletResponse{
		Resp: event.RestoredResp,
	}

	return nil
}

func infoCmd(wal *Wallet, _ *event.ArgumentQueue) error {
	wallets, err := wal.wallets()
	if err != nil {
		return err
	}
	var completeTotal int64
	for _, wall := range wallets {
		iter, err := wall.AccountsIterator(2) // Placeholder
		if err != nil {
			return err
		}
		for acct := iter.Next(); acct != nil; acct = iter.Next() {
			completeTotal += acct.TotalBalance
		}
	}
	best := wal.multi.GetBestBlock()

	wal.Send <- event.WalletInfo{
		LoadedWallets:   len(wallets),
		TotalBalance:    completeTotal,
		BestBlockHeight: best.Height,
		BestBlockTime:   best.Timestamp,
	}
	return nil
}
