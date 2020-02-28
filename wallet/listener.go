package wallet

import (
	"github.com/raedahgroup/dcrlibwallet"
)

const (
	// FetchHeadersStep is the first step when a wallet is syncing.
	FetchHeadersStep = 1
	// AddressDiscoveryStep is the third step when a wallet is syncing.
	AddressDiscoveryStep = 2
	// RescanHeadersStep is the second step when a wallet is syncing.
	RescanHeadersStep = 3
	// TotalSyncSteps is the total number of steps to complete a sync process
	TotalSyncSteps = 3
)

type listener struct {
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

func (l *listener) Debug(info *dcrlibwallet.DebugInfo) {
	log.Trace(info)
}

func (l *listener) OnSyncStarted(restarted bool) {
	l.Send <- Response{
		Resp: SyncStarted{
			WasRestarted: restarted,
		},
	}
}

func (l *listener) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	l.Send <- Response{
		Resp: SyncPeersChanged{
			ConnectedPeers: numberOfConnectedPeers,
		},
	}
}

func (l *listener) OnHeadersFetchProgress(progress *dcrlibwallet.HeadersFetchProgressReport) {
	l.Send <- Response{
		Resp: SyncHeadersFetchProgress{
			Progress: progress,
		},
	}
}
func (l *listener) OnAddressDiscoveryProgress(progress *dcrlibwallet.AddressDiscoveryProgressReport) {
	l.Send <- Response{
		Resp: SyncAddressDiscoveryProgress{
			Progress: progress,
		},
	}
}

func (l *listener) OnHeadersRescanProgress(progress *dcrlibwallet.HeadersRescanProgressReport) {
	l.Send <- Response{
		Resp: SyncHeadersRescanProgress{
			Progress: progress,
		},
	}
}

func (l *listener) OnSyncCompleted() {
	l.Send <- Response{
		Resp: SyncCompleted{},
	}
}

func (l *listener) OnSyncCanceled(willRestart bool) {
	l.Send <- Response{
		Resp: SyncCanceled{
			WillRestart: willRestart,
		},
	}
}

func (l *listener) OnSyncEndedWithError(err error) {
	// todo: create custom sync error
	l.Send <- Response{
		Resp: SyncEndedWithError{
			Error: err,
		},
	}
}
