package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

type SyncProgressListener struct {
	SyncStatusChan chan wallet.SyncStatusUpdate
}

func NewSyncProgress(syncStatus chan wallet.SyncStatusUpdate) *SyncProgressListener {
	return &SyncProgressListener{
		SyncStatusChan: syncStatus,
	}
}

func (sp *SyncProgressListener) OnSyncStarted(wasRestarted bool) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncStarted,
	})
}

func (sp *SyncProgressListener) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.PeersConnected,
		ConnectedPeers: numberOfConnectedPeers,
	})
}

func (sp *SyncProgressListener) OnCFiltersFetchProgress(cfiltersFetchProgress *dcrlibwallet.CFiltersFetchProgressReport) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.CfiltersFetchProgress,
		ProgressReport: cfiltersFetchProgress,
	})
}

func (sp *SyncProgressListener) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.HeadersFetchProgress,
		ProgressReport: headersFetchProgress,
	})
}

func (sp *SyncProgressListener) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.AddressDiscoveryProgress,
		ProgressReport: addressDiscoveryProgress,
	})
}

func (sp *SyncProgressListener) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.HeadersRescanProgress,
		ProgressReport: headersRescanProgress,
	})
}
func (sp *SyncProgressListener) OnSyncCompleted() {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCompleted,
	})
}

func (sp *SyncProgressListener) OnSyncCanceled(willRestart bool) {
	sp.sendNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCanceled,
	})
}
func (sp *SyncProgressListener) OnSyncEndedWithError(err error)          {}
func (sp *SyncProgressListener) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

func (sp *SyncProgressListener) sendNotification(signal wallet.SyncStatusUpdate) {
	sp.SyncStatusChan <- signal
}
