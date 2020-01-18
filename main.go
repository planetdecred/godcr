package main

import (
	"fmt"

	app "gioui.org/app"
	"gioui.org/font/gofont"

	"github.com/raedahgroup/godcr-gio/ui/page"
	"github.com/raedahgroup/godcr-gio/wallet"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return
	}
	syncChan := make(chan int)
	go func(c chan int) {
		_, atleastone, err := wallet.LoadWallets(cfg.HomeDir, cfg.Network)
		if err != nil {
			fmt.Println(err)
			return
		}
		sig := 0
		if atleastone {
			sig = 1
		}
		c <- sig
		close(c)
	}(syncChan)

	gofont.Register() // IMPORTANT

	win, err := createWindow(page.LoadingID, syncChan)
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}
	go func(win *window) {
		win.loop()
	}(win)

	app.Main()
}
