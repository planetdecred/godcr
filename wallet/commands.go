package wallet

import (
	"errors"

	"github.com/raedahgroup/godcr-gio/event"
)

var (
	// ErrInvalidArguments is returned when a wallet command is send with invalid arguments.
	ErrInvalidArguments = errors.New("invalid command arguments")

	// ErrNotFound is returned when a wallet command is given that does not exist or is not
	// implemented.
	ErrNotFound = errors.New("command not found or not implemented")

	// ErrNoSuchWallet is returned with the wallet requested by the given id does not exist
	ErrNoSuchWallet = errors.New("no such wallet with id")

	// ErrNoSuchAcct is returned when the given account number cannot be found
	ErrNoSuchAcct = errors.New("no such account")

	// ErrCreateTx is returned when a tx author cannot be created
	ErrCreateTx = errors.New("can not create transaction")
)

var cmdMap = map[string]func(*Wallet, *event.ArgumentQueue) error{
	event.CreateCmd:   createCmd,
	event.RestoreCmd:  restoreCmd,
	event.InfoCmd:     infoCmd,
	event.CreateTxCmd: createTxCmd,
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

func createTxCmd(wal *Wallet, arguments *event.ArgumentQueue) error {
	wallets, err := wal.wallets()
	if err != nil {
		return err
	}

	walletID, err := arguments.PopInt()
	if err != nil {
		return ErrInvalidArguments
	}

	if walletID > len(wallets) || walletID < 0 {
		return ErrNoSuchWallet
	}

	acct, err := arguments.PopInt()
	if err != nil {
		return ErrInvalidArguments
	}

	confirms, err := arguments.PopInt()
	if err != nil {
		return ErrInvalidArguments
	}

	if _, err := wallets[walletID].GetAccount(int32(acct), int32(confirms)); err != nil {
		return ErrNoSuchAcct
	}

	txAuthor := wallets[walletID].NewUnsignedTx(int32(acct), int32(confirms))
	if txAuthor == nil {
		return ErrCreateTx
	}

	wal.Send <- &event.WalletResponse{
		Resp: event.CreatedTxResp,
		Results: &event.ArgumentQueue{
			Queue: []interface{}{txAuthor},
		},
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

	wal.Send <- &event.WalletInfo{
		LoadedWallets:   len(wallets),
		TotalBalance:    completeTotal,
		BestBlockHeight: best.Height,
		BestBlockTime:   best.Timestamp,
		Synced:          wal.multi.IsSynced(),
	}
	return nil
}
