package wallet

import (
	"github.com/raedahgroup/dcrlibwallet"
)

type progressListener struct {
	Send chan<- Response
}

// SyncStarted is sent when sync starts
type SyncStarted struct {
	WasRestarted bool
}

// SyncCompleted is sent when the sync is completed
type SyncCompleted struct{}

// SyncEndedWithError is sent when the sync ends with and error
type SyncEndedWithError struct {
	Error error
}

// SyncCanceled is sent when the sync is canceled
type SyncCanceled struct {
	WillRestart bool
}

// SyncPeersChanged is sent when the amount of connected peers changes during sync
type SyncPeersChanged struct {
	ConnectedPeers int32
}

// SyncHeadersFetchProgress is sent whenever syncing makes any progress in fetching headers
type SyncHeadersFetchProgress struct {
	Progress *dcrlibwallet.HeadersFetchProgressReport
}

// SyncAddressDiscoveryProgress is sent whenever syncing makes any progress in discovering addresses
type SyncAddressDiscoveryProgress struct {
	Progress *dcrlibwallet.AddressDiscoveryProgressReport
}

// SyncHeadersRescanProgress is sent whenever syncing makes any progress in rescanning headers
type SyncHeadersRescanProgress struct {
	Progress *dcrlibwallet.HeadersRescanProgressReport
}

func (listener *progressListener) Debug(info *dcrlibwallet.DebugInfo) {
	// Log Traces
}

func (listener *progressListener) OnSyncStarted(restarted bool) {
	listener.Send <- Response{
		Resp: SyncStarted{
			WasRestarted: restarted,
		},
	}
}

func (listener *progressListener) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	listener.Send <- Response{
		Resp: SyncPeersChanged{
			ConnectedPeers: numberOfConnectedPeers,
		},
	}
}

func (listener *progressListener) OnHeadersFetchProgress(progress *dcrlibwallet.HeadersFetchProgressReport) {
	listener.Send <- Response{
		Resp: SyncHeadersFetchProgress{
			Progress: progress,
		},
	}
}
func (listener *progressListener) OnAddressDiscoveryProgress(progress *dcrlibwallet.AddressDiscoveryProgressReport) {
	listener.Send <- Response{
		Resp: SyncAddressDiscoveryProgress{
			Progress: progress,
		},
	}
}

func (listener *progressListener) OnHeadersRescanProgress(progress *dcrlibwallet.HeadersRescanProgressReport) {
	listener.Send <- Response{
		Resp: SyncHeadersRescanProgress{
			Progress: progress,
		},
	}
}

func (listener *progressListener) OnSyncCompleted() {
	listener.Send <- Response{
		Resp: SyncCompleted{},
	}
}

func (listener *progressListener) OnSyncCanceled(willRestart bool) {
	listener.Send <- Response{
		Resp: SyncCanceled{
			WillRestart: willRestart,
		},
	}
}

func (listener *progressListener) OnSyncEndedWithError(err error) {
	// todo: create custom sync error
	listener.Send <- Response{
		Resp: SyncEndedWithError{
			Error: err,
		},
	}
}
