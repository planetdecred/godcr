package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/wallet"
)

const (
	transactionStatusConfirmed = "confirmed"
	transactionStatusPending   = "pending"

	walletStatusWaiting = "waiting for other wallets"
	walletStatusSyncing = "syncing..."
	walletStatusSynced  = "synced"

	errWalletIDNotFound = "wallet ID not found"
)

// Balance takes the balance as an integer and returns its string equivalent
func Balance(balance int64) string {
	return dcrutil.Amount(balance).String()
}

// BreakBalance takes the balance string and returns it in two slices
func BreakBalance(balance string) (b1, b2 string) {
	balanceParts := strings.Split(balance, ".")
	if len(balanceParts) == 1 {
		return balanceParts[0], ""
	}
	b1 = balanceParts[0]
	b2 = balanceParts[1]
	b1 = b1 + "." + b2[:2]
	b2 = b2[2:]
	return
}

// TransactionStatus accepts the bestBlockHeight, transactionBlockHeight and returns a transaction status
// which could be confirmed or pending
func TransactionStatus(bestBlockHeight, txnBlockHeight int32) string {
	confirmations := bestBlockHeight - txnBlockHeight + 1
	if txnBlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		return transactionStatusConfirmed
	}
	return transactionStatusPending
}

// divMod divides a numerator by a denominator and returns its quotient and remainder.
func divMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

// RemainingSyncTime takes time on int64 and returns its string equivalent.
func RemainingSyncTime(totalTimeLeft int64) string {
	q, r := divMod(totalTimeLeft, 24*60*60)
	timeLeft := time.Duration(r) * time.Second
	if q > 0 {
		return fmt.Sprintf("%dd%s", q, timeLeft.String())
	}
	return timeLeft.String()
}

// WalletSyncStatus returns the sync status of a single wallet
func WalletSyncStatus(info wallet.InfoShort, bestBlockHeight int32) string {
	if info.IsWaiting {
		return walletStatusWaiting
	}
	if info.BestBlockHeight < bestBlockHeight {
		return walletStatusSyncing
	}

	return walletStatusSynced
}

// WalletSyncProgressTime returns the sync time in days
func WalletSyncProgressTime(timestamp int64) string {
	return fmt.Sprintf("%s behind", dcrlibwallet.CalculateDaysBehind(timestamp))
}

// LastBlockSync returns how long ago the current block was attached
func LastBlockSync(timestamp int64) string {
	return truncateTime(time.Since(time.Unix(timestamp, 0)).String(), 0)
}

// truncateTime takes a time duration in string and chops off decimal places in the string
// by the specified number of places.
func truncateTime(duration string, place int) string {
	var durationCharacter string
	durationSlice := strings.Split(duration, ".")
	if len(durationSlice) == 1 {
		return duration
	}

	secondsDecimals := durationSlice[1]
	if place > len(secondsDecimals) {
		return duration
	}

	secondLastCharacter := secondsDecimals[len(secondsDecimals)-2 : len(secondsDecimals)-1]
	_, err := strconv.Atoi(secondLastCharacter)
	if err != nil {
		durationCharacter = secondsDecimals[len(secondsDecimals)-2:]
	} else {
		durationCharacter = secondsDecimals[len(secondsDecimals)-1:]
	}
	if place == 0 {
		return durationSlice[0] + durationCharacter
	}
	return durationSlice[0] + "." + secondsDecimals[0:place] + durationCharacter
}

// WalletNameFromID gets the id of a wallet by using its name
func WalletNameFromID(id int, walletInfo []wallet.InfoShort) (string, error) {
	for _, info := range walletInfo {
		if info.ID == id {
			return info.Name, nil
		}
	}
	return "", fmt.Errorf(errWalletIDNotFound)
}
