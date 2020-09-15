package main

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/wallet"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	dcrlibwallet.SetLogLevels("off")
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

	// use wait group to keep main alive until shutdown completes
	shutdownWaitGroup := new(sync.WaitGroup)
	go listenForShutdownRequests()
	go handleShutdownRequests(shutdownWaitGroup)

	cli := NewCli(wal)
	reflect.ValueOf(cfg.Wallet)
	cli.loop(cfg.Sync, "createwallet")

	// wait for handleShutdown goroutine, to finish before exiting main
	shutdownWaitGroup.Wait()
}
