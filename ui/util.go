package ui

import (
	"fmt"
	"gioui.org/gesture"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
	"strings"
	"time"
)

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func checkLockWallet(c pageCommon) {
	walletsLocked := getLockWallet(c)
	if len(walletsLocked) > 0 {
		go func() {
			c.modalReceiver <- &modalLoad{
				template: UnlockWalletRestoreTemplate,
				title:    "Unlock to resume restoration",
				confirm: func(pass string) {
					err := c.wallet.UnlockWallet(walletsLocked[0].ID, []byte(pass))
					if err != nil {
						errText := err.Error()
						if err.Error() == "invalid_passphrase" {
							errText = "Invalid passphrase"
						}
						c.notify(errText, false)
					} else {
						c.closeModal()
					}
				},
				confirmText: "Unlock",
				cancel:      "",
				cancelText:  "",
			}
		}()
	}
}

// getLockWallet get all the lock wallets
func getLockWallet(c pageCommon) []*dcrlibwallet.Wallet {
	allWallets := c.wallet.AllWallets()

	var walletsLocked []*dcrlibwallet.Wallet
	for _, wl := range allWallets {
		if !wl.HasDiscoveredAccounts && wl.IsLocked() {
			walletsLocked = append(walletsLocked, wl)
		}
	}

	return walletsLocked
}

func formatDateOrTime(timestamp int64) string {
	utcTime := time.Unix(timestamp, 0).UTC()
	if time.Now().UTC().Sub(utcTime).Hours() < 168 {
		return utcTime.Weekday().String()
	}

	t := strings.Split(utcTime.Format(time.UnixDate), " ")
	t2 := t[2]
	if t[2] == "" {
		t2 = t[3]
	}
	return fmt.Sprintf("%s %s", t[1], t2)
}

// createClickGestures returns a slice of click gestures
func createClickGestures(count int) []*gesture.Click {
	var gestures = make([]*gesture.Click, count)
	for i := 0; i < count; i++ {
		gestures[i] = &gesture.Click{}
	}
	return gestures
}

// showBadge loops through a slice of recent transactions and checks if there are transaction from different wallets.
// It returns true if transactions from different wallets exists and false if they don't
func showLabel(recentTransactions []wallet.Transaction) bool {
	var name string
	for _, t := range recentTransactions {
		if name != "" && name != t.WalletName {
			return true
		}
		name = t.WalletName
	}
	return false
}
