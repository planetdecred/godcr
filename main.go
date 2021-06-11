package main

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/app"

	_ "net/http/pprof"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui"
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

	absoluteWdPath, err := ui.GetAbsolutePath()
	if err != nil {
		panic(err)
	}

	decredIcons := make(map[string]image.Image)
	err = filepath.Walk(filepath.Join(absoluteWdPath, "ui/assets/decredicons"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() || !strings.HasSuffix(path, ".png") {
			return nil
		}

		f, _ := os.Open(path)
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		split := strings.Split(info.Name(), ".")
		decredIcons[split[0]] = img
		return nil
	})
	if err != nil {
		log.Warn(err)
	}

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

	collection := fontCollection()

	win, appWindow, err := ui.CreateWindow(wal, decredIcons, collection, internalLog)
	if err != nil {
		fmt.Printf("Could not initialize window: %s\ns", err)
		return
	}

	// Start the ui frontend
	go win.Loop(appWindow, shutdown)
	app.Main()
}
