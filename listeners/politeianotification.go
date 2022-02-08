package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

type PoliteiaNotification struct {
	PoliteiaNotifCh chan wallet.Proposal
}

func NewPoliteiaNotification(notifCh chan wallet.Proposal) *PoliteiaNotification {
	return &PoliteiaNotification{
		PoliteiaNotifCh: notifCh,
	}
}

func (pn *PoliteiaNotification) OnProposalsSynced() {
	pn.sendNotification(wallet.Proposal{
		ProposalStatus: wallet.Synced,
	})
}

func (pn *PoliteiaNotification) OnNewProposal(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.NewProposalFound,
		Proposal:       proposal,
	}
	pn.sendNotification(update)
}

func (pn *PoliteiaNotification) OnProposalVoteStarted(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.VoteStarted,
		Proposal:       proposal,
	}
	pn.sendNotification(update)
}
func (pn *PoliteiaNotification) OnProposalVoteFinished(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.VoteFinished,
		Proposal:       proposal,
	}
	pn.sendNotification(update)
}

func (pn *PoliteiaNotification) sendNotification(signal wallet.Proposal) {
	pn.PoliteiaNotifCh <- signal
}
