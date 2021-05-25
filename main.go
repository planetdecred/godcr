package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "net/http/pprof"

	"gioui.org/app"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui"
	"github.com/planetdecred/godcr/ui/values"
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

	wal.LoadWallets()
	shutdown := make(chan int)
	go func() {
		<-shutdown
		wal.Shutdown()
		os.Exit(0)
	}()

	dbPath := filepath.Join(cfg.HomeDir, cfg.Network, "dexc.db")
	dc, err := dexc.NewDex(cfg.DebugLevel, dbPath, cfg.Network, make(chan dexc.Response, 3), logWriter{})
	if err != nil {
		fmt.Printf("error creating Dex: %s", err)
		return
	}
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dc.Run(appCtx, cancel)

	var netType string
	if strings.Contains(wal.Net, "testnet") {
		netType = "testnet"
	} else {
		netType = wal.Net
	}
	w := app.NewWindow(app.Size(values.AppWidth, values.AppHeight), app.Title(fmt.Sprintf("%s (%s)", "godcr", netType)))

	// Create ui
	appui, err := ui.NewUI(w, wal, dc, internalLog)
	if err != nil {
		fmt.Printf("Could not initialize wallet UI: %s\ns", err)
		return
	}

	// Start the ui frontend
	appui.Loop(shutdown, w)
	app.Main()
}
