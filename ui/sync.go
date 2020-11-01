package ui

import (
	"github.com/planetdecred/godcr/wallet"
)

type ProposalNotificationListeners struct {
	win *Window
}

func (p ProposalNotificationListeners) OnNewProposal(proposalID int, censorshipToken string) {
	p.win.wallet.GetProposalUpdate(censorshipToken, 1)
}

func (p ProposalNotificationListeners) OnProposalVoteStarted(proposalID int, censorshipToken string) {
	p.win.wallet.GetProposalUpdate(censorshipToken, 2)
}

func (p ProposalNotificationListeners) OnProposalVoteFinished(proposalID int, censorshipToken string) {
	p.win.wallet.GetProposalUpdate(censorshipToken, 3)
}

// SyncProposals starts syncing the multiwallet proposals
/**func (win Window) SyncProposals() {
	win.isSyncingProposals = true
	win.wallet.StartProposalsSync(ProposalNotificationListeners{win: &win})
}**/

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
	case wallet.NewBlock:
		for _, info := range win.walletInfo.Wallets {
			if info.ID == t.WalletID {
				win.wallet.GetAllTransactions(0, 0, 0)
				break
			}
		}
	case wallet.TxConfirmed:
		if t.Hash != "" {
			win.wallet.GetAllTransactions(0, 0, 0)
			if win.walletTransaction != nil &&
				win.walletTransaction.Txn.Hash == t.Hash {
				win.wallet.GetTransaction(t.WalletID, t.Hash)
			}
		}
	}
}

// updateConnectedPeers updates connected peers in the SyncStatus state
func (win Window) updateConnectedPeers(peers int32) {
	win.walletSyncStatus.ConnectedPeers = peers
}
