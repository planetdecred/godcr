package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/wallet"
)

type Cli struct {
	wallet             *wallet.Wallet
	walletInfo         *wallet.MultiWalletInfo
	walletSyncStatus   *wallet.SyncStatus
	walletTransactions *wallet.Transactions
	walletTransaction  *wallet.Transaction
	signatureResult    *wallet.Signature
	selectedAccount    int
	txAuthor           dcrlibwallet.TxAuthor
	broadcastResult    wallet.Broadcast
}

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
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Err() != nil {
		// handle error.
	}

loop:
	for {
		select {
		case e := <-wal.Send:
			switch resp := e.Resp.(type) {
			case wallet.LoadedWallets:
				if resp.Count > 0 {
					log.Infof("Syncing %d wallets...", resp.Count)
					wal.StartSync()
				} else {
					goto end
				}
			case wallet.MultiWalletInfo:
				log.Info(">>>>> Your TotalBalance", resp.TotalBalance)
				goto end
			}
		case update := <-wal.Sync:
			if update.Stage == wallet.SyncCompleted {
				log.Info("Wallets synced")
				if cfg.Wallet.Balance {
					goto balance
				}
				if cfg.Wallet.Send != "" {
					goto send
				}
				if cfg.Wallet.Receive != "" {
					goto receive
				}
			}
		}
	}

balance:
	{
		wal.GetMultiWalletInfo()
		goto loop
	}

send:
	{
		fmt.Print("Enter amount to send: ")
		// for scanner.Scan() {
		scanner.Scan()
		line := scanner.Text()
		log.Info(line)
		scanner.Scan()
		fmt.Print("Enter passphare: ")
		line2 := scanner.Text()
		log.Info(line2)
		// wal.BroadcastTransaction()
		goto end
	}

receive:
	{
		fmt.Print("Enter amount to receve: ")
		scanner.Scan()
		line := scanner.Text()
		fmt.Print(line)
		goto end
	}

end:
	log.Info("Bye!")
}
