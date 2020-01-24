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
	wal.Send <- event.WalletResponse{
		Resp: event.CreatedResp,
		Results: &event.ArgumentQueue{
			Queue: []interface{}{wall.Seed},
		},
	}
	return err
}
