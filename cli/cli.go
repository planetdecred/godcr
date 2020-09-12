package main

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/wallet"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
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
	// scanner := bufio.NewScanner(os.Stdin)
	// scan := make(chan bool)

	// go func() { scan <- scanner.Scan() }()

out:
	for {
		select {
		case e := <-wal.Send:
			switch resp := e.Resp.(type) {
			case wallet.LoadedWallets:
				if resp.Count > 0 {
					wal.StartSync()
				}
			case wallet.MultiWalletInfo:
				log.Info(">>>>> TotalBalance", resp.TotalBalance)
				break out
			}
		case update := <-wal.Sync:
			if update.Stage == wallet.SyncCompleted {
				if cfg.Wallet.Balance {
					wal.GetMultiWalletInfo()
				}
			}
		}
	}

	log.Info("Bye!")
}
