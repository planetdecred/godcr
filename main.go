package main

import (
	"fmt"
	"sync"

	app "gioui.org/app"
	"gioui.org/font/gofont"

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
	wal, dup, err := wallet.New(cfg.HomeDir, cfg.Network)

	if err != nil {
		fmt.Println("Error loading wallet") // TODO: Show error on frontend
		return
	}

	// Start up the wallet backend
	wg.Add(1)
	go wal.Sync(&wg)

	gofont.Register() // IMPORTANT
	win, err := window.CreateWindow(page.LoadingID, dup.Reverse())
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}
	// Start the ui frontend
	// Does not need to be added to the WaitGroup, app.Main() handles that
	go win.Loop()

	app.Main()

	wg.Wait()
}
