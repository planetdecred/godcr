package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

type SyncProgress struct {
	SyncStatus chan wallet.SyncStatusUpdate
}

func NewSyncProgress(syncStatus chan wallet.SyncStatusUpdate) *SyncProgress {
	return &SyncProgress{
		SyncStatus: syncStatus,
	}
}

func (sp *SyncProgress) OnSyncStarted(wasRestarted bool) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncStarted,
	})
}

func (sp *SyncProgress) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.PeersConnected,
		ConnectedPeers: numberOfConnectedPeers,
	})
}

func (sp *SyncProgress) OnCFiltersFetchProgress(cfiltersFetchProgress *dcrlibwallet.CFiltersFetchProgressReport) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.CfiltersFetchProgress,
		ProgressReport: cfiltersFetchProgress,
	})
}

func (sp *SyncProgress) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersFetchProgress,
		ProgressReport: wallet.SyncHeadersFetchProgress{
			Progress: headersFetchProgress,
		},
	})
}

func (sp *SyncProgress) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.AddressDiscoveryProgress,
		ProgressReport: wallet.SyncAddressDiscoveryProgress{
			Progress: addressDiscoveryProgress,
		},
	})
}

func (sp *SyncProgress) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersRescanProgress,
		ProgressReport: wallet.SyncHeadersRescanProgress{
			Progress: headersRescanProgress,
		},
	})
}
func (sp *SyncProgress) OnSyncCompleted() {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCompleted,
	})
}

func (sp *SyncProgress) OnSyncCanceled(willRestart bool) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCanceled,
	})
}
func (sp *SyncProgress) OnSyncEndedWithError(err error)          {}
func (sp *SyncProgress) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

func (sp *SyncProgress) sendNotification(signal wallet.SyncStatusUpdate) {
	if sp.SyncStatus != nil {
		sp.SyncStatus <- signal
	}
}
