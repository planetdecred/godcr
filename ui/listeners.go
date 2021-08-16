package ui

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gen2brain/beeep"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/page"
	"github.com/planetdecred/godcr/wallet"
)

// Transaction notifications
func (mp *mainPage) OnTransaction(transaction string) {
	mp.updateBalance()

	var tx dcrlibwallet.Transaction
	err := json.Unmarshal([]byte(transaction), &tx)
	if err == nil {
		mp.updateNotification(wallet.NewTransaction{
			Transaction: &tx,
		})
	}
	mp.desktopNotifier(tx)
}

func (mp *mainPage) OnBlockAttached(walletID int, blockHeight int32) {
	mp.updateBalance()
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.BlockAttached,
	})
}

func (mp *mainPage) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	mp.updateBalance()
}

// Account mixer
func (mp *mainPage) OnAccountMixerStarted(walletID int) {}
func (mp *mainPage) OnAccountMixerEnded(walletID int)   {}

// Politeia notifications
func (mp *mainPage) OnProposalsSynced() {
	mp.notificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.Synced,
	}
}

func (mp *mainPage) OnNewProposal(proposal *dcrlibwallet.Proposal) {
	mp.notificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.NewProposalFound,
		Proposal:       proposal,
	}
	mp.desktopNotifier(prop)
}

func (mp *mainPage) OnProposalVoteStarted(proposal *dcrlibwallet.Proposal) {
	mp.notificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.VoteStarted,
		Proposal:       proposal,
	}
	mp.desktopNotifier(prop)
}
func (mp *mainPage) OnProposalVoteFinished(proposal *dcrlibwallet.Proposal) {
	mp.notificationsUpdate <- wallet.Proposal{
		ProposalStatus: wallet.VoteFinished,
		Proposal:       proposal,
	}
	mp.desktopNotifier(prop)
}

// Sync notifications

func (mp *mainPage) OnSyncStarted(wasRestarted bool) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncStarted,
	})
}

func (mp *mainPage) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.PeersConnected,
		ConnectedPeers: numberOfConnectedPeers,
	})
}

func (mp *mainPage) OnCFiltersFetchProgress(cfiltersFetchProgress *dcrlibwallet.CFiltersFetchProgressReport) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.CfiltersFetchProgress,
		ProgressReport: cfiltersFetchProgress,
	})
}

func (mp *mainPage) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersFetchProgress,
		ProgressReport: wallet.SyncHeadersFetchProgress{
			Progress: headersFetchProgress,
		},
	})
}

func (mp *mainPage) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.AddressDiscoveryProgress,
		ProgressReport: wallet.SyncAddressDiscoveryProgress{
			Progress: addressDiscoveryProgress,
		},
	})
}

func (mp *mainPage) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersRescanProgress,
		ProgressReport: wallet.SyncHeadersRescanProgress{
			Progress: headersRescanProgress,
		},
	})
}

func (mp *mainPage) OnSyncCompleted() {
	mp.updateBalance()
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCompleted,
	})
}

func (mp *mainPage) OnSyncCanceled(willRestart bool) {
	mp.updateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCanceled,
	})
}

func (mp *mainPage) OnSyncEndedWithError(err error) {}

func (mp *mainPage) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

// todo: this will be removed when all pages have been moved to the page package
// updateNotification sends notification to the notification channel depending on which channel the page uses
func (mp *mainPage) updateNotification(signal interface{}) {
	switch *mp.page {
	case page.OverviewPageID, page.Transactions:
		mp.load.Receiver.NotificationsUpdate <- signal
	default:
		mp.notificationsUpdate <- signal
		}
}

func (mp *mainPage) desktopNotifier(notifier interface{}) {
	var notification string
	switch t := notifier.(type) {
	case dcrlibwallet.Transaction:
		// remove trailing zeros from amount and convert to string
		amount := strconv.FormatFloat(dcrlibwallet.AmountCoin(t.Amount), 'f', -1, 64)

		defaultNotification := "Transaction notification of %s DCR"
		notificationText := "You have %s %s DCR "

		wallet := mp.pageCommon.multiWallet.WalletWithID(t.WalletID)
		if wallet == nil {
			notification = fmt.Sprintf(defaultNotification, amount)
			return
		}

		getAccount := func(acct int32) string {
			var account string
			accountName, err := wallet.AccountName(acct)
			if err != nil {
				log.Error(err)
			} else {
				account = accountName
			}
			return account
		}

		// get source account
		var txSourceAccount string
		if t.Direction == dcrlibwallet.TxDirectionSent ||
			t.Direction == dcrlibwallet.TxDirectionTransferred {
			for _, input := range t.Inputs {
				if input.AccountNumber != -1 {
					txSourceAccount = getAccount(input.AccountNumber)
				}
			}
		}

		//get distination account
		var txDestAccount, txDestinationAddress string
		if t.Direction == dcrlibwallet.TxDirectionTransferred ||
			t.Direction == dcrlibwallet.TxDirectionReceived ||
			t.Direction == dcrlibwallet.TxDirectionSent {
			for _, output := range t.Outputs {
				if output.AccountNumber != -1 {
					txDestAccount = getAccount(output.AccountNumber)
				}
				if output.AccountNumber == -1 {
					txDestinationAddress = output.Address
				}
			}
		}

		switch {
		case t.Direction == dcrlibwallet.TxDirectionReceived:
			notification = fmt.Sprintf(notificationText+"to %s account in %s wallet.", "received", amount, txDestAccount, wallet.Name)
		case t.Direction == dcrlibwallet.TxDirectionSent:
			notification = fmt.Sprintf(notificationText+"from %s account in %s wallet to %s.", "sent", amount, wallet.Name, txSourceAccount, txDestinationAddress)
		case t.Direction == dcrlibwallet.TxDirectionTransferred:
			notification = fmt.Sprintf(notificationText+"from %s account, to %s account in %s wallet.", "transferred", amount, txSourceAccount, txDestAccount, wallet.Name)
		default:
			notification = fmt.Sprintf(defaultNotification+" to/from (%s wallet) account %s", amount, wallet.Name, txSourceAccount)
		}
	case wallet.Proposal:
		fmt.Println(t.ProposalStatus)
		switch {
		case t.ProposalStatus == wallet.NewProposalFound:
			notification = fmt.Sprintf("A new proposal has been added Token: %s", t.Proposal.Token)
		case t.ProposalStatus == wallet.VoteStarted:
			notification = fmt.Sprintf("Voting has started for proposal with Token: %s", t.Proposal.Token)
		case t.ProposalStatus == wallet.VoteFinished:
			notification = fmt.Sprintf("Voting has ended for proposal with Token: %s", t.Proposal.Token)
		default:
			notification = fmt.Sprintf("New update for proposal with Token", t.Proposal.Token)
		}
	}
	fmt.Println(t)

	err := beeep.Notify("Decred Godcr Wallet", notification, "assets/information.png")
	if err != nil {
		log.Info("could not initiate desktop notification, reason:", err.Error())
	}
}
