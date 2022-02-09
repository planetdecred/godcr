package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

// ProposalNotificationListener satisfies dcrlibwallet ProposalNotificationListener interface contract.
type ProposalNotificationListener struct {
	ProposalNotifChan chan wallet.Proposal
}

func NewProposalNotificationListener(notifCh chan wallet.Proposal) *ProposalNotificationListener {
	return &ProposalNotificationListener{
		ProposalNotifChan: notifCh,
	}
}

func (pn *ProposalNotificationListener) OnProposalsSynced() {
	pn.sendNotification(wallet.Proposal{
		ProposalStatus: wallet.Synced,
	})
}

func (pn *ProposalNotificationListener) OnNewProposal(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.NewProposalFound,
		Proposal:       proposal,
	}
	pn.sendNotification(update)
}

func (pn *ProposalNotificationListener) OnProposalVoteStarted(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.VoteStarted,
		Proposal:       proposal,
	}
	pn.sendNotification(update)
}
func (pn *ProposalNotificationListener) OnProposalVoteFinished(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.VoteFinished,
		Proposal:       proposal,
	}
	pn.sendNotification(update)
}

func (pn *ProposalNotificationListener) sendNotification(signal wallet.Proposal) {
	pn.ProposalNotifChan <- signal
}
