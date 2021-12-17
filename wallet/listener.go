package wallet

import (
	"github.com/planetdecred/dcrlibwallet"
)

type RescanNotificationType int

const (

	// RescanStarted indicates a block rescan start signal
	RescanStarted RescanNotificationType = iota

	// RescanProgress indicates a block rescan progress signal
	RescanProgress

	// RescanEnded indicates a block rescan end signal
	RescanEnded
)

// SyncNotificationType represents the spv sync stage at which the multiwallet is currently
type SyncNotificationType int

const (
	// SyncStarted signifies that spv sync has started
	SyncStarted SyncNotificationType = iota

	// SyncCanceled is a pseudo stage that represents a canceled sync
	SyncCanceled

	// SyncCompleted signifies that spv sync has been completed
	SyncCompleted

	// CfiltersFetchProgress indicates a cfilters fetch signal
	CfiltersFetchProgress

	// HeadersFetchProgress indicates a headers fetch signal
	HeadersFetchProgress

	// HeadersFetchProgress indicates an address discovery signal
	AddressDiscoveryProgress

	// HeadersRescanProgress indicates an address rescan signal
	HeadersRescanProgress

	// PeersConnected indicates a peer connected signal
	PeersConnected

	// BlockAttached indicates a block attached signal
	BlockAttached

	// BlockConfirmed indicates a block update signal
	BlockConfirmed

	// AccountMixerStarted indicates on account mixer started
	AccountMixerStarted

	// AccountMixerEnded indicates on account mixer ended
	AccountMixerEnded

	// ProposalVoteFinished indicates that proposal voting is finished
	ProposalVoteFinished

	// ProposalVoteStarted indicates that proposal voting has started
	ProposalVoteStarted

	// ProposalSynced indicates that proposal has finished syncing
	ProposalSynced

	// ProposalAdded indicates that a new proposal was added
	ProposalAdded
)

const (
	// FetchHeadersStep is the first step when a wallet is syncing.
	FetchHeadersSteps = iota + 1

	// AddressDiscoveryStep is the third step when a wallet is syncing.
	AddressDiscoveryStep

	// RescanHeadersStep is the second step when a wallet is syncing.
	RescanHeadersStep
)

// TotalSyncSteps is the total number of steps to complete a sync process
const TotalSyncSteps = 3

type (
	// SyncStatusUpdate represents information about the status of the multiwallet spv sync
	SyncStatusUpdate struct {
		Stage          SyncNotificationType
		ProgressReport interface{}
		ConnectedPeers int32
		BlockInfo      NewBlock
		ConfirmedTxn   TxConfirmed
		AcctMixerInfo  AccountMixer
		Proposal       Proposal
	}

	RescanUpdate struct {
		Stage          RescanNotificationType
		WalletID       int
		ProgressReport *dcrlibwallet.HeadersRescanProgressReport
	}
)

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
