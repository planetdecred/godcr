package ui

import (
	"fmt"
	"time"

	"github.com/raedahgroup/godcr-gio/wallet"
)

// updateSyncStatus updates the sync status in the walletInfo state.
func (win Window) updateSyncStatus(syncing, synced bool) {
	win.walletInfo.Syncing = syncing
	win.walletInfo.Synced = synced
}

func (win Window) updateSyncProgress(report interface{}) {
	status := win.walletSyncStatus
	switch t := report.(type) {
	case wallet.SyncHeadersFetchProgress:
		status.HeadersFetchProgress = t.Progress.HeadersFetchProgress
		status.HeadersToFetch = t.Progress.TotalHeadersToFetch
		status.Progress = t.Progress.TotalSyncProgress
		status.RemainingTime = remainingSyncTime(t.Progress.TotalTimeRemainingSeconds)
		status.TotalSteps = wallet.TotalSyncSteps
		status.Steps = wallet.FetchHeadersSteps
		status.CurrentBlockHeight = t.Progress.CurrentHeaderHeight
		win.wallet.OverallBlockHeight = t.Progress.TotalHeadersToFetch
		win.wallet.GetMultiWalletInfo()
	case wallet.SyncAddressDiscoveryProgress:
		status.RescanHeadersProgress = t.Progress.AddressDiscoveryProgress
		status.Progress = t.Progress.TotalSyncProgress
		status.RemainingTime = remainingSyncTime(t.Progress.TotalTimeRemainingSeconds)
		status.TotalSteps = wallet.TotalSyncSteps
		status.Steps = wallet.AddressDiscoveryStep
	case wallet.SyncHeadersRescanProgress:
		status.RescanHeadersProgress = t.Progress.RescanProgress
		status.Progress = t.Progress.TotalSyncProgress
		status.RemainingTime = remainingSyncTime(t.Progress.TotalTimeRemainingSeconds)
		status.TotalSteps = wallet.TotalSyncSteps
		status.Steps = wallet.RescanHeadersStep
	}
}

// updateConnectedPeers updates connected peers in the SyncStatus state
func (win Window) updateConnectedPeers(peers int32) {
	win.walletSyncStatus.ConnectedPeers = peers
}

// remainingSyncTime takes time on int64 and returns its string equivalent.
func remainingSyncTime(totalTimeLeft int64) string {
	q, r := divMod(totalTimeLeft, 24*60*60)
	timeLeft := time.Duration(r) * time.Second
	if q > 0 {
		return fmt.Sprintf("%dd%s", q, timeLeft.String())
	}
	return timeLeft.String()
}

// divMod divides a numerator by a denominator and returns its quotient and remainder.
func divMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
