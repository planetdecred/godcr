package uiwallet

import (
	"github.com/gen2brain/beeep"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

// updateSyncStatus updates the sync status in the walletInfo state.
func (w Wallet) updateSyncStatus(syncing, synced bool) {
	w.walletInfo.Syncing = syncing
	w.walletInfo.Synced = synced
}

func (w Wallet) updateSyncProgress(report interface{}) {
	status := w.walletSyncStatus
	switch t := report.(type) {
	case wallet.SyncHeadersFetchProgress:
		status.HeadersFetchProgress = t.Progress.HeadersFetchProgress
		status.HeadersToFetch = t.Progress.TotalHeadersToFetch
		status.Progress = t.Progress.TotalSyncProgress
		status.RemainingTime = wallet.SecondsToDays(t.Progress.TotalTimeRemainingSeconds)
		status.TotalSteps = wallet.TotalSyncSteps
		status.Steps = wallet.FetchHeadersSteps
		status.CurrentBlockHeight = t.Progress.CurrentHeaderHeight
		w.wallet.OverallBlockHeight = t.Progress.TotalHeadersToFetch
		w.wallet.GetMultiWalletInfo()
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
	case wallet.NewBlock:
		for _, info := range w.walletInfo.Wallets {
			if info.ID == t.WalletID {
				w.wallet.GetAllTransactions(0, 0, 0)
				break
			}
		}

		beep := w.wallet.ReadBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey)
		if beep {
			err := beeep.Beep(5, 1)
			if err != nil {
				log.Error(err.Error())
			}
		}
	case wallet.TxConfirmed:
		if t.Hash != "" {
			w.wallet.GetAllTransactions(0, 0, 0)
			if w.walletTransaction != nil &&
				w.walletTransaction.Txn.Hash == t.Hash {
				w.wallet.GetTransaction(t.WalletID, t.Hash)
			}
		}
	}
}

// updateConnectedPeers updates connected peers in the SyncStatus state
func (w Wallet) updateConnectedPeers(peers int32) {
	w.walletSyncStatus.ConnectedPeers = peers
}
