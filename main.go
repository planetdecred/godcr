package main

import (
	"fmt"
	"net/http"
	"os"

	"gioui.org/app"

	_ "net/http/pprof"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui"
	_ "github.com/planetdecred/godcr/ui/assets"
	"github.com/planetdecred/godcr/wallet"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	if cfg.Profile > 0 {
		go func() {
			log.Info(fmt.Sprintf("Starting profiling server on port %d", cfg.Profile))
			log.Error(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", cfg.Profile), nil))
		}()
	}

	dcrlibwallet.SetLogLevels(cfg.DebugLevel)

	var confirms int32 = dcrlibwallet.DefaultRequiredConfirmations

	if cfg.SpendUnconfirmed {
		confirms = 0
	}

	wal, err := wallet.NewWallet(cfg.HomeDir, cfg.Network, make(chan wallet.Response, 3), confirms)
	if err != nil {
		log.Error(err)
		return
	}

	shutdown := make(chan int)
	go func() {
		<-shutdown
		wal.Shutdown()
		os.Exit(0)
	}()

	win, appWindow, err := ui.CreateWindow(wal, internalLog)
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}

	// Start the ui frontend
	go win.Loop(appWindow, shutdown)

	app.Main()
}
