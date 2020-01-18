package main

import (
	"fmt"

	app "gioui.org/app"
	"gioui.org/font/gofont"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/page"
	"github.com/raedahgroup/godcr-gio/wallet"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	walrecieve := make(chan event.Event) // chan the wallet recieves from
	walsend := make(chan event.Event)    // chan the wallet sends events to
	wal := &wallet.Wallet{
		Root:        cfg.HomeDir,
		Network:     cfg.Network,
		SendChan:    walsend,
		ReceiveChan: walrecieve,
	}

	go wal.Sync()

	gofont.Register() // IMPORTANT

	win, err := createWindow(page.LoadingID, walsend, walrecieve)
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}
	go func(win *window) {
		win.loop()
	}(win)

	app.Main()
}
