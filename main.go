package main

import (
	"fmt"

	app "gioui.org/app"
	"gioui.org/font/gofont"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/page"
	"github.com/raedahgroup/godcr-gio/ui/window"
	"github.com/raedahgroup/godcr-gio/wallet"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	fmt.Printf("godcr v%s\n", Version())
	walrecieve := make(chan event.Event) // chan the wallet recieves commands from
	walsend := make(chan event.Event)    // chan the wallet sends events to
	wal := &wallet.Wallet{
		Root:    cfg.HomeDir,
		Network: cfg.Network,
		Duplex: event.Duplex{
			Send:    walsend,
			Receive: walrecieve,
		},
	}

	// Start up the wallet backend
	go wal.Sync()

	defer func(c chan<- event.Event) {
		c <- event.WalletCmd{Cmd: event.ShutdownCmd}
	}(walrecieve)

	gofont.Register() // IMPORTANT
	win, err := window.CreateWindow(page.LoadingID, event.Duplex{Receive: walsend, Send: walrecieve})
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}
	// Start the ui frontend
	go win.Loop()

	app.Main()
	// TODO: wait for the wallet to finish shutting down
}
