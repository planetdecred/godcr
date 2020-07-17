package main

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "net/http/pprof"

	"gioui.org/app"
	"github.com/markbates/pkger"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui"
	"github.com/raedahgroup/godcr/wallet"
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

	decredIcons := make(map[string]image.Image)
	err = pkger.Walk("/ui/assets/decredicons", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() || !strings.HasSuffix(path, ".png") {
			return nil
		}

		f, _ := pkger.Open(path)
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

	wal.LoadWallets()

	var wg sync.WaitGroup
	shutdown := make(chan int)
	wg.Add(1)
	go func() {
		<-shutdown
		wal.Shutdown()
		wg.Done()
	}()

	win, err := ui.CreateWindow(wal, decredIcons)
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
