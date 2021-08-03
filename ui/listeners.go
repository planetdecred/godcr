package ui

import (
	"encoding/json"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/page"
	"github.com/planetdecred/godcr/wallet"
)

// Transaction notifications

func (mp *page.MainPage) OnTransaction(transaction string) {
	mp.UpdateBalance()

	// beeep send notification

	var tx dcrlibwallet.Transaction
	err := json.Unmarshal([]byte(transaction), &tx)
	if err == nil {
		mp.updateNotification(wallet.NewTransaction{
			Transaction: &tx,
		})
	}
}

func (mp *page.MainPage) OnBlockAttached(walletID int, blockHeight int32) {
	mp.updateBalance()
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.BlockAttached,
	})
}

func (mp *page.MainPage) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	mp.updateBalance()
}

// Account mixer
func (mp *page.MainPage) OnAccountMixerStarted(walletID int) {}
func (mp *page.MainPage) OnAccountMixerEnded(walletID int)   {}

// Politeia notifications
func (mp *page.MainPage) OnProposalsSynced() {
	mp.notificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.Synced,
	}
}

func (mp *page.MainPage) OnNewProposal(proposal *dcrlibwallet.Proposal) {
	mp.notificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.NewProposalFound,
		Proposal:       proposal,
	}
}

func (mp *page.MainPage) OnProposalVoteStarted(proposal *dcrlibwallet.Proposal) {
	mp.notificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.VoteStarted,
		Proposal:       proposal,
	}
}
func (mp *page.MainPage) OnProposalVoteFinished(proposal *dcrlibwallet.Proposal) {
	mp.notificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.VoteFinished,
		Proposal:       proposal,
	}
}

// Sync notifications

func (mp *page.MainPage) OnSyncStarted(wasRestarted bool) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncStarted,
	})
}

func (mp *page.MainPage) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.PeersConnected,
		ConnectedPeers: numberOfConnectedPeers,
	})
}

func (mp *page.MainPage) OnCFiltersFetchProgress(cfiltersFetchProgress *dcrlibwallet.CFiltersFetchProgressReport) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.CfiltersFetchProgress,
		ProgressReport: cfiltersFetchProgress,
	})
}

func (mp *page.MainPage) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersFetchProgress,
		ProgressReport: wallet.SyncHeadersFetchProgress{
			Progress: headersFetchProgress,
		},
	})
}
func (mp *page.MainPage) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.AddressDiscoveryProgress,
		ProgressReport: wallet.SyncAddressDiscoveryProgress{
			Progress: addressDiscoveryProgress,
		},
	})
}

func (mp *page.MainPage) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersRescanProgress,
		ProgressReport: wallet.SyncHeadersRescanProgress{
			Progress: headersRescanProgress,
		},
	})
}
func (mp *page.MainPage) OnSyncCompleted() {
	mp.updateBalance()
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCompleted,
	})
}

func (mp *page.MainPage) OnSyncCanceled(willRestart bool) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCanceled,
	})
}
func (mp *page.MainPage) OnSyncEndedWithError(err error)          {}
func (mp *page.MainPage) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

// todo: this will be removed when all pages have been moved to the page package
// updateNotification sends notification to the notification channel depending on which channel the page uses
func (mp *page.MainPage) updateNotification(signal interface{}) {
	switch *mp.page {
	case page.OverviewPageID, page.TransactionsPageID:
		mp.Load.Receiver.NotificationsUpdate <- signal
	default:
		mp.notificationsUpdate <- signal
	}
}
