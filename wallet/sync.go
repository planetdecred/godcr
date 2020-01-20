package wallet

import (
	"fmt"

	"github.com/raedahgroup/godcr-gio/event"
)

// Sync is the main wallet sync loop
func (wal *Wallet) Sync() {
	loaded, err := wal.loadWallets()
	if err != nil {
		wal.Send <- err
		return
	}
	defer wal.multi.Shutdown()

	wal.Send <- event.Loaded{
		WalletsLoadedCount: loaded,
	}

	for {
		e := <-wal.Receive
		if cmd, ok := e.(event.WalletCmd); ok {
			switch cmd.Cmd {
			case event.ShutdownCmd:
				return
			}
		} else {
			fmt.Printf("Not a wallet command %+v\n", e)
		}
	}
}
