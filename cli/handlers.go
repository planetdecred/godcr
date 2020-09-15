package main

import (
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/wallet"
)

type Cli struct {
	wall            *wallet.Wallet
	signatureResult *wallet.Signature
	txAuthor        dcrlibwallet.TxAuthor
	broadcastResult wallet.Broadcast
	handlers        map[string]handler
	errChan         chan error
}

func NewCli(wall *wallet.Wallet) *Cli {
	return &Cli{
		wall:     wall,
		handlers: handlers,
		errChan:  make(chan error),
	}
}

func (c *Cli) loop(shouldSync bool, command string) {
out:
	for {
		select {
		case err := <-c.errChan:
			log.Error(err)
			close(beginShutdown)
		case e := <-c.wall.Send:
			switch resp := e.Resp.(type) {
			case wallet.LoadedWallets:
				log.Infof("Loaded %d wallets...", resp.Count)
				if resp.Count > 0 && shouldSync {
					log.Infof("Syncing %d wallets...", resp.Count)
					c.wall.StartSync()
					break out
				}

				handler := c.handlers[command]
				_, err := handler(c, nil)
				if err != nil {
					log.Error(err)
				}

			case wallet.MultiWalletInfo:
				log.Info("Total balance: ", resp.TotalBalance)
				for _, w := range resp.Wallets {
					log.Info("Wallet: ", w.Name, "|", w.Balance)
					for _, acct := range w.Accounts {
						log.Info(acct.Name, "|", acct.TotalBalance)
					}
				}
				close(beginShutdown)
			case wallet.CreatedSeed:
				log.Info("Wallet created")
				close(beginShutdown)
			case wallet.MultiWalletError:
				log.Error(resp.Message)
			case wallet.Transactions:
			}

		case update := <-c.wall.Sync:
			if update.Stage == wallet.SyncCompleted {
				// if cfg.Wallet.Balance {
				// 	goto balance
				// }
				// if cfg.Wallet.Send != "" {
				// 	goto send
				// }
				// if cfg.Wallet.Receive != "" {
				// 	goto receive
				// }
			}
		}
	}
}
