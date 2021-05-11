package wallet

import (
	"github.com/planetdecred/dcrlibwallet"
)

type ProposalStatus int

const (
	Synced ProposalStatus = iota
	VoteStarted
	NewProposalFound
	VoteFinished
)

type Proposal struct {
	Proposal       *dcrlibwallet.Proposal
	ProposalStatus ProposalStatus
}

func (l *listener) OnNewProposal(proposal *dcrlibwallet.Proposal) {
	l.Send <- SyncStatusUpdate{
		Stage: ProposalAdded,
		Proposal: Proposal{
			Proposal:       proposal,
			ProposalStatus: NewProposalFound,
		},
	}
}

func (l *listener) OnProposalVoteStarted(proposal *dcrlibwallet.Proposal) {
	l.Send <- SyncStatusUpdate{
		Stage: ProposalVoteStarted,
		Proposal: Proposal{
			Proposal:       proposal,
			ProposalStatus: VoteStarted,
		},
	}
}

func (l *listener) OnProposalVoteFinished(proposal *dcrlibwallet.Proposal) {
	l.Send <- SyncStatusUpdate{
		Stage: ProposalVoteFinished,
		Proposal: Proposal{
			Proposal:       proposal,
			ProposalStatus: VoteFinished,
		},
	}
}

func (l *listener) OnProposalsSynced() {
	l.Send <- SyncStatusUpdate{
		Stage: ProposalSynced,
		Proposal: Proposal{
			ProposalStatus: Synced,
		},
	}
}
