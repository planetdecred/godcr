package ui

import (
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
		status.RemainingTime = wallet.SecondsToDays(t.Progress.TotalTimeRemainingSeconds)
		status.TotalSteps = wallet.TotalSyncSteps
		status.Steps = wallet.FetchHeadersSteps
		status.CurrentBlockHeight = t.Progress.CurrentHeaderHeight
		win.wallet.OverallBlockHeight = t.Progress.TotalHeadersToFetch
		win.wallet.GetMultiWalletInfo()
	case wallet.SyncAddressDiscoveryProgress:
		status.RescanHeadersProgress = t.Progress.AddressDiscoveryProgress
		status.Progress = t.Progress.TotalSyncProgress
		status.RemainingTime = wallet.SecondsToDays(t.Progress.TotalTimeRemainingSeconds)
		status.TotalSteps = wallet.TotalSyncSteps
		status.Steps = wallet.AddressDiscoveryStep
	case wallet.SyncHeadersRescanProgress:
		status.RescanHeadersProgress = t.Progress.RescanProgress
		status.Progress = t.Progress.TotalSyncProgress
		status.RemainingTime = wallet.SecondsToDays(t.Progress.TotalTimeRemainingSeconds)
		status.TotalSteps = wallet.TotalSyncSteps
		status.Steps = wallet.RescanHeadersStep
	case wallet.TxConfirmed:
		win.wallet.GetAllTransactions(0, 0, 0)
	}
}

// updateConnectedPeers updates connected peers in the SyncStatus state
func (win Window) updateConnectedPeers(peers int32) {
	win.walletSyncStatus.ConnectedPeers = peers
}
