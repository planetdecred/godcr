package main

import (
	"fmt"
	"sync"

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

	var wg sync.WaitGroup

	dup := event.NewDuplexBase()
	gofont.Register() // IMPORTANT
	win, err := window.CreateWindow(page.LoadingID, dup.Reverse())
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}
	// Start the ui frontend
	// Does not need to be added to the WaitGroup, app.Main() handles that
	go win.Loop()

	wal := wallet.NewWallet(cfg.HomeDir, cfg.Network, dup.Duplex())

	// Start up the wallet backend
	wg.Add(1)
	go wal.Sync(&wg)

	app.Main()

	wg.Wait()
}
