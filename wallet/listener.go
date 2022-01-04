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

	// AgendaVoteFinished indicates that agenda voting is finished
	AgendaVoteFinished

	// AgendaVoteStarted indicates that agenda voting has started
	AgendaVoteStarted

	// AgendaSynced indicates that agenda has finished syncing
	AgendaSynced

	// AgendaAdded indicates that a new agenda was added
	AgendaAdded
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
		Agenda			Agenda
	}

	RescanUpdate struct {
		Stage          RescanNotificationType
		WalletID       int
		ProgressReport *dcrlibwallet.HeadersRescanProgressReport
	}
)
