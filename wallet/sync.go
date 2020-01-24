package wallet

import (
	"fmt"
	"sync"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/event"
)

var cmdMap = map[string]func(*Wallet, *event.ArgumentQueue) error{
	event.CreateCmd: createCmd,
}

// Sync is the main wallet sync loop
func (wal *Wallet) Sync(wg *sync.WaitGroup) {
	defer wg.Done()
	if wal.multi == nil {
		return
	}

	defer wal.multi.Shutdown()

	for {
		e := <-wal.Receive
		if cmd, ok := e.(event.WalletCmd); ok {
			switch cmd.Cmd {
			case event.LoadedWalletsCmd:
				wal.Send <- event.Loaded{
					WalletsLoadedCount: wal.multi.LoadedWalletsCount(),
				}
			case event.ShutdownCmd:
				return
			default:
				if fun, ok := cmdMap[cmd.Cmd]; ok {
					err := fun(wal, cmd.Arguments)
					if err != nil {
						wal.Send <- err
					}
				}
			}
		} else {
			fmt.Printf("Not a wallet command %+v\n", e)
		}
	}
}

type progressListener struct {
	send chan<- event.Event
}

func (listener progressListener) Debug(info *dcrlibwallet.DebugInfo) {

}

func (listener progressListener) OnSyncStarted() {
	listener.send <- event.SyncEvent{Event: event.SyncStart}
}
