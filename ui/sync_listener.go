package ui

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

func (mp *mainPage) OnSyncStarted(wasRestarted bool) {
	log.Info("Main page sync started")
	mp.syncStatusUpdate <- wallet.SyncStatusUpdate{
		Stage: wallet.SyncStarted,
	}
}

func (mp *mainPage) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	mp.syncStatusUpdate <- wallet.SyncStatusUpdate{
		Stage:          wallet.PeersConnected,
		ConnectedPeers: numberOfConnectedPeers,
	}
}

func (mp *mainPage) OnCFiltersFetchProgress(cfiltersFetchProgress *dcrlibwallet.CFiltersFetchProgressReport) {
	mp.syncStatusUpdate <- wallet.SyncStatusUpdate{
		Stage:          wallet.CfiltersFetchProgress,
		ProgressReport: cfiltersFetchProgress,
	}
}

func (mp *mainPage) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	mp.syncStatusUpdate <- wallet.SyncStatusUpdate{
		Stage: wallet.HeadersFetchProgress,
		ProgressReport: wallet.SyncHeadersFetchProgress{
			Progress: headersFetchProgress,
		},
	}
}
func (mp *mainPage) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	mp.syncStatusUpdate <- wallet.SyncStatusUpdate{
		Stage: wallet.AddressDiscoveryProgress,
		ProgressReport: wallet.SyncAddressDiscoveryProgress{
			Progress: addressDiscoveryProgress,
		},
	}
}

func (mp *mainPage) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	mp.syncStatusUpdate <- wallet.SyncStatusUpdate{
		Stage: wallet.HeadersRescanProgress,
		ProgressReport: wallet.SyncHeadersRescanProgress{
			Progress: headersRescanProgress,
		},
	}
}
func (mp *mainPage) OnSyncCompleted() {
	mp.syncStatusUpdate <- wallet.SyncStatusUpdate{
		Stage: wallet.SyncCompleted,
	}
}

func (mp *mainPage) OnSyncCanceled(willRestart bool) {
	mp.syncStatusUpdate <- wallet.SyncStatusUpdate{
		Stage: wallet.SyncCanceled,
	}
}
func (mp *mainPage) OnSyncEndedWithError(err error)          {}
func (mp *mainPage) Debug(debugInfo *dcrlibwallet.DebugInfo) {}
