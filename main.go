package main

import (
	"fmt"

	app "gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/text"

	"github.com/markbates/pkger"

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

	source, err := pkger.Open("/ui/fonts/source_sans_pro_regular.otf")
	if err != nil {
		fmt.Println("Failed to load font")
		gofont.Register()
	} else {
		stat, err := source.Stat()
		if err != nil {
			fmt.Println(err)
		}
		bytes := make([]byte, stat.Size())
		source.Read(bytes)
		fnt, err := opentype.Parse(bytes)
		if err != nil {
			fmt.Println(err)
		}
		font.Register(text.Font{}, fnt)
	}

	wal, _ := wallet.NewWallet(cfg.HomeDir, cfg.Network, make(chan interface{}))
	wal.LoadWallets()

	win, err := window.CreateWindow(page.LoadingID, wal)
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}
	// Start the ui frontend
	// Does not need to be added to the WaitGroup, app.Main() handles that
	go win.Loop()

	app.Main()
	wal.Shutdown()
}
