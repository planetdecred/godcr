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
