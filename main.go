package main

import (
	"fmt"
	"sync"

	app "gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/font/opentype"
	"gioui.org/text"

	"github.com/markbates/pkger"

	"github.com/raedahgroup/dcrlibwallet"
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

	dcrlibwallet.SetLogLevels(cfg.DebugLevel)
	sans, err := pkger.Open("/ui/fonts/source_sans_pro_regular.otf")
	if err != nil {
		log.Warn("Failed to load font Source Sans Pro. Using gofont")
		gofont.Register()
	} else {
		stat, err := sans.Stat()
		if err != nil {
			log.Warn(err)
		}
		bytes := make([]byte, stat.Size())
		sans.Read(bytes)
		fnt, err := opentype.Parse(bytes)
		if err != nil {
			log.Warn(err)
		}
		if fnt != nil {
			font.Register(text.Font{}, fnt)
		} else {
			log.Warn("Failed to load font Source Sans Pro. Using gofont")
			gofont.Register()
		}
	}

	wal, _ := wallet.NewWallet(cfg.HomeDir, cfg.Network, make(chan wallet.Response))
	wal.LoadWallets()

	var wg sync.WaitGroup
	shutdown := make(chan int)
	wg.Add(1)
	go func() {
		<-shutdown
		wal.Shutdown()
		wg.Done()
	}()

	win, err := window.CreateWindow(page.LoadingID, wal)
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}

	// Start the ui frontend
	// Does not need to be added to the WaitGroup, app.Main() handles that
	go win.Loop(shutdown)

	app.Main()
	wg.Wait()
}
