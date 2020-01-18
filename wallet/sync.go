package wallet

import (
	"github.com/raedahgroup/godcr-gio/event"
)

// Sync is the main wallet sync loop
func (wal *Wallet) Sync() {
	loaded, err := wal.loadWallets()
	if err != nil {
		wal.SendChan <- err
		return
	}
	defer wal.multi.Shutdown()

	wal.SendChan <- event.Loaded{
		WalletsLoadedCount: loaded,
	}
	close(wal.SendChan)
}
