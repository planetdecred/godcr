package page

import (
	"encoding/json"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

// Transaction notifications

func (mp *MainPage) OnTransaction(transaction string) {
	mp.UpdateBalance()

	// beeep send notification

	var tx dcrlibwallet.Transaction
	err := json.Unmarshal([]byte(transaction), &tx)
	if err == nil {
		mp.UpdateNotification(wallet.NewTransaction{
			Transaction: &tx,
		})
	}
}

func (mp *MainPage) OnBlockAttached(walletID int, blockHeight int32) {
	mp.UpdateBalance()
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.BlockAttached,
	})
}

func (mp *MainPage) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	mp.UpdateBalance()
}

// Account mixer
func (mp *MainPage) OnAccountMixerStarted(walletID int) {}
func (mp *MainPage) OnAccountMixerEnded(walletID int)   {}

// Politeia notifications
func (mp *MainPage) OnProposalsSynced() {
	mp.Load.Receiver.NotificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.Synced,
	}
}

func (mp *MainPage) OnNewProposal(proposal *dcrlibwallet.Proposal) {
	mp.Load.Receiver.NotificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.NewProposalFound,
		Proposal:       proposal,
	}
}

func (mp *MainPage) OnProposalVoteStarted(proposal *dcrlibwallet.Proposal) {
	mp.Load.Receiver.NotificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.VoteStarted,
		Proposal:       proposal,
	}
}
func (mp *MainPage) OnProposalVoteFinished(proposal *dcrlibwallet.Proposal) {
	mp.Load.Receiver.NotificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.VoteFinished,
		Proposal:       proposal,
	}
}

// Sync notifications

func (mp *MainPage) OnSyncStarted(wasRestarted bool) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncStarted,
	})
}

func (mp *MainPage) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.PeersConnected,
		ConnectedPeers: numberOfConnectedPeers,
	})
}

func (mp *MainPage) OnCFiltersFetchProgress(cfiltersFetchProgress *dcrlibwallet.CFiltersFetchProgressReport) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.CfiltersFetchProgress,
		ProgressReport: cfiltersFetchProgress,
	})
}

func (mp *MainPage) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersFetchProgress,
		ProgressReport: wallet.SyncHeadersFetchProgress{
			Progress: headersFetchProgress,
		},
	})
}
func (mp *MainPage) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.AddressDiscoveryProgress,
		ProgressReport: wallet.SyncAddressDiscoveryProgress{
			Progress: addressDiscoveryProgress,
		},
	})
}

func (mp *MainPage) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersRescanProgress,
		ProgressReport: wallet.SyncHeadersRescanProgress{
			Progress: headersRescanProgress,
		},
	})
}
func (mp *MainPage) OnSyncCompleted() {
	mp.UpdateBalance()
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCompleted,
	})
}

func (mp *MainPage) OnSyncCanceled(willRestart bool) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCanceled,
	})
}
func (mp *MainPage) OnSyncEndedWithError(err error)          {}
func (mp *MainPage) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

// UpdateNotification sends notification to the notification channel
func (mp *MainPage) UpdateNotification(signal interface{}) {
	mp.Load.Receiver.NotificationsUpdate <- signal
}
