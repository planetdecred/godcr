package wallet

import (
	"github.com/raedahgroup/dcrlibwallet"
)

// SyncProgressStage represents the spv sync stage at which the multiwallet is currently
type SyncProgressStage int

const (
	// SyncStarted signifies that spv sync has started
	SyncStarted SyncProgressStage = iota

	// SyncCanceled is a pseudo stage that represents a canceled sync
	SyncCanceled

	// SyncCompleted signifies that spv sync has been completed
	SyncCompleted
)

// SyncStatusUpdate represents information about the status of the multiwallet spv sync
type SyncStatusUpdate struct {
	Stage SyncProgressStage
}

type listener struct {
	Send chan<- SyncStatusUpdate
}

// // SyncCompleted is sent when the sync is completed
// type SyncCompleted struct{}

// // SyncEndedWithError is sent when the sync ends with and error
// type SyncEndedWithError struct {
// 	Error error
// }

// // SyncCanceled is sent when the sync is canceled
// type SyncCanceled struct {
// 	WillRestart bool
// }

// // SyncPeersChanged is sent when the amount of connected peers changes during sync
// type SyncPeersChanged struct {
// 	ConnectedPeers int32
// }

// // SyncHeadersFetchProgress is sent whenever syncing makes any progress in fetching headers
// type SyncHeadersFetchProgress struct {
// 	Progress *dcrlibwallet.HeadersFetchProgressReport
// }

// // SyncAddressDiscoveryProgress is sent whenever syncing makes any progress in discovering addresses
// type SyncAddressDiscoveryProgress struct {
// 	Progress *dcrlibwallet.AddressDiscoveryProgressReport
// }

// // SyncHeadersRescanProgress is sent whenever syncing makes any progress in rescanning headers
// type SyncHeadersRescanProgress struct {
// 	Progress *dcrlibwallet.HeadersRescanProgressReport
// }

func (l *listener) Debug(info *dcrlibwallet.DebugInfo) {
	log.Trace(info)
}

func (l *listener) OnSyncStarted(restarted bool) {
	l.Send <- SyncStatusUpdate{
		Stage: SyncStarted,
	}
}

func (l *listener) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	// l.Send <- SyncPeersChanged{
	// 	ConnectedPeers: numberOfConnectedPeers,
	// }
}

func (l *listener) OnHeadersFetchProgress(progress *dcrlibwallet.HeadersFetchProgressReport) {
	// l.Send <- SyncHeadersFetchProgress{
	// 	Progress: progress,
	// }
}
func (l *listener) OnAddressDiscoveryProgress(progress *dcrlibwallet.AddressDiscoveryProgressReport) {
	// l.Send <- SyncAddressDiscoveryProgress{
	// 	Progress: progress,
	// }
}

func (l *listener) OnHeadersRescanProgress(progress *dcrlibwallet.HeadersRescanProgressReport) {
	// l.Send <- SyncHeadersRescanProgress{
	// 	Progress: progress,
	// }
}

func (l *listener) OnSyncCompleted() {
	l.Send <- SyncStatusUpdate{
		Stage: SyncCompleted,
	}
}

func (l *listener) OnSyncCanceled(willRestart bool) {
	l.Send <- SyncStatusUpdate{
		Stage: SyncCanceled,
	}
	// l.Send <- SyncCanceled{
	// 	WillRestart: willRestart,
	// }
}

func (l *listener) OnSyncEndedWithError(err error) {
	// todo: create custom sync error
	// l.Send <- SyncEndedWithError{
	// 	Error: err,
	// }
}
