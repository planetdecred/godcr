package wallet

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/planetdecred/dcrlibwallet"
)

// transactionStatus accepts the bestBlockHeight, transactionBlockHeight returns a transaction status
// which could be confirmed/pending and confirmations count
func transactionStatus(bestBlockHeight, txnBlockHeight int32) (string, int32) {
	confirmations := bestBlockHeight - txnBlockHeight + 1
	if txnBlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		return "confirmed", confirmations
	}
	return "pending", confirmations
}

// WalletSyncStatus returns the sync status of a single wallet
func walletSyncStatus(isWaiting bool, walletBestBlock, bestBlockHeight int32) string {
	if isWaiting {
		return "waiting for other wallets"
	}
	if walletBestBlock < bestBlockHeight {
		return "syncing..."
	}

	return "synced"
}

func (wal *Wallet) GetMultiWallet() *dcrlibwallet.MultiWallet {
	return wal.multi
}

func (wal *Wallet) UnlockWallet(walletID int, password []byte) error {
	return wal.multi.UnlockWallet(walletID, password)
}

// CurrentAddress returns the next address for the specified wallet account.
func (wal *Wallet) CurrentAddress(walletID int, accountID int32) (string, error) {
	wall := wal.multi.WalletWithID(walletID)
	if wall == nil {
		return "", ErrIDNotExist
	}
	return wall.CurrentAddress(accountID)
}

// NextAddress returns the next address for the specified wallet account.
func (wal *Wallet) NextAddress(walletID int, accountID int32) (string, error) {
	wall := wal.multi.WalletWithID(walletID)
	if wall == nil {
		return "", ErrIDNotExist
	}
	return wall.NextAddress(accountID)
}

// IsAddressValid checks if the given address is valid for the multiwallet network
func (wal *Wallet) IsAddressValid(address string) (bool, error) {
	return wal.multi.IsAddressValid(address), nil
}

// HaveAddress checks if the given address is valid for the wallet
func (wal *Wallet) HaveAddress(address string) (bool, string) {
	for _, wallet := range wal.multi.AllWallets() {
		result := wallet.HaveAddress(address)
		if result {
			return true, wallet.Name
		}
	}
	return false, ""
}

// VerifyMessage checks if the given message matches the signature for the address.
func (wal *Wallet) VerifyMessage(address string, message string, signature string) (bool, error) {
	return wal.multi.VerifyMessage(address, message, signature)
}

// StartSync starts the multiwallet SPV sync
func (wal *Wallet) StartSync() error {
	return wal.multi.SpvSync()
}

// RescanBlocks rescans the multiwallet
func (wal *Wallet) RescanBlocks(walletID int) error {
	return wal.multi.RescanBlocks(walletID)
}

func (wal *Wallet) IsSyncingProposals() bool {
	return wal.multi.Politeia.IsSyncing()
}

func (wal *Wallet) GetWalletSeedPhrase(walletID int, password []byte) (string, error) {
	return wal.multi.WalletWithID(walletID).DecryptSeed(password)
}

func (wal *Wallet) VerifyWalletSeedPhrase(walletID int, seedPhrase string, privpass []byte) error {
	_, err := wal.multi.VerifySeedForWallet(walletID, seedPhrase, privpass)
	return err
}

func (wal *Wallet) SaveConfigValueForKey(key string, value interface{}) {
	wal.multi.SaveUserConfigValue(key, value)
}

func (wal *Wallet) ReadBoolConfigValueForKey(key string) bool {
	return wal.multi.ReadBoolConfigValueForKey(key, false)
}

func (wal *Wallet) ReadStringConfigValueForKey(key string) string {
	return wal.multi.ReadStringConfigValueForKey(key)
}

func calculateDaysBehind(lastHeaderTime int64) string {
	diff := time.Since(time.Unix(lastHeaderTime, 0))
	daysBehind := int(math.Round(diff.Hours() / 24))
	if daysBehind < 1 {
		return "<1 day"
	} else if daysBehind == 1 {
		return "1 day"
	} else {
		return fmt.Sprintf("%d days", daysBehind)
	}
}

// SecondsToDays takes time in seconds and returns its string equivalent in the format ddhhmm.
func SecondsToDays(totalTimeLeft int64) string {
	q, r := divMod(totalTimeLeft, 24*60*60)
	timeLeft := time.Duration(r) * time.Second
	if q > 0 {
		return fmt.Sprintf("%dd%s", q, timeLeft.String())
	}
	return timeLeft.String()
}

// GetRawBalance gets the balance in int64, formats it and returns a string while also leaving out the "DCR" suffix
func GetRawBalance(balance int64, AmountUnit int) string {
	return strconv.FormatFloat(float64(balance)/math.Pow10(AmountUnit+8), 'f', -(AmountUnit + 8), 64)
}

// divMod divides a numerator by a denominator and returns its quotient and remainder.
func divMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
