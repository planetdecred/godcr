package page

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gen2brain/beeep"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/wallet"
)

// Transaction notifications

func (mp *MainPage) OnTransaction(transaction string) {
	transactionNotification := mp.WL.Wallet.ReadBoolConfigValueForKey(load.TransactionNotificationConfigKey)
	load.GetUSDExchangeValue(&mp.dcrUsdtBittrex)
	mp.UpdateBalance()

	var tx dcrlibwallet.Transaction
	err := json.Unmarshal([]byte(transaction), &tx)
	if err == nil {
		update := wallet.NewTransaction{
			Transaction: &tx,
		}
		mp.UpdateNotification(update)
		if transactionNotification {
			mp.desktopNotifier(update)
		}

	}
}

func (mp *MainPage) OnBlockAttached(walletID int, blockHeight int32) {
	beep := mp.WL.Wallet.ReadBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey)
	if beep {
		err := beeep.Beep(5, 1)
		if err != nil {
			log.Error(err.Error())
		}
	}

	load.GetUSDExchangeValue(&mp.dcrUsdtBittrex)
	mp.UpdateBalance()
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.BlockAttached,
	})
}

func (mp *MainPage) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	load.GetUSDExchangeValue(&mp.dcrUsdtBittrex)
	mp.UpdateBalance()
}

// Account mixer
func (mp *MainPage) OnAccountMixerStarted(walletID int) {
	mp.UpdateNotification(wallet.AccountMixer{
		WalletID:  walletID,
		RunStatus: wallet.MixerStarted,
	})
}

func (mp *MainPage) OnAccountMixerEnded(walletID int) {
	mp.UpdateNotification(wallet.AccountMixer{
		WalletID:  walletID,
		RunStatus: wallet.MixerEnded,
	})
}

// Politeia notifications
func (mp *MainPage) OnProposalsSynced() {
	mp.UpdateNotification(wallet.Proposal{
		ProposalStatus: wallet.Synced,
	})
}

func (mp *MainPage) OnNewProposal(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.NewProposalFound,
		Proposal:       proposal,
	}
	mp.UpdateNotification(update)
	mp.desktopNotifier(update)
}

func (mp *MainPage) OnProposalVoteStarted(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.VoteStarted,
		Proposal:       proposal,
	}
	mp.UpdateNotification(update)
	mp.desktopNotifier(update)
}
func (mp *MainPage) OnProposalVoteFinished(proposal *dcrlibwallet.Proposal) {
	update := wallet.Proposal{
		ProposalStatus: wallet.VoteFinished,
		Proposal:       proposal,
	}
	mp.UpdateNotification(update)
	mp.desktopNotifier(update)
}

// Sync notifications

func (mp *MainPage) OnSyncStarted(wasRestarted bool) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncStarted,
	})
}

func (mp *MainPage) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.PeersConnected,
		ConnectedPeers: numberOfConnectedPeers,
	})
}

func (mp *MainPage) OnCFiltersFetchProgress(cfiltersFetchProgress *dcrlibwallet.CFiltersFetchProgressReport) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage:          wallet.CfiltersFetchProgress,
		ProgressReport: cfiltersFetchProgress,
	})
}

func (mp *MainPage) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersFetchProgress,
		ProgressReport: wallet.SyncHeadersFetchProgress{
			Progress: headersFetchProgress,
		},
	})
}

func (mp *MainPage) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.AddressDiscoveryProgress,
		ProgressReport: wallet.SyncAddressDiscoveryProgress{
			Progress: addressDiscoveryProgress,
		},
	})
}

func (mp *MainPage) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.HeadersRescanProgress,
		ProgressReport: wallet.SyncHeadersRescanProgress{
			Progress: headersRescanProgress,
		},
	})
}
func (mp *MainPage) OnSyncCompleted() {
	load.GetUSDExchangeValue(&mp.dcrUsdtBittrex)
	mp.UpdateBalance()
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCompleted,
	})
}

func (mp *MainPage) OnSyncCanceled(willRestart bool) {
	mp.UpdateNotification(wallet.SyncStatusUpdate{
		Stage: wallet.SyncCanceled,
	})
}
func (mp *MainPage) OnSyncEndedWithError(err error)          {}
func (mp *MainPage) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

// Block Rescan

func (mp *MainPage) OnBlocksRescanStarted(walletID int) {
	mp.UpdateNotification(wallet.RescanUpdate{
		Stage:    wallet.RescanStarted,
		WalletID: walletID,
	})
}

func (mp *MainPage) OnBlocksRescanProgress(progress *dcrlibwallet.HeadersRescanProgressReport) {
	mp.UpdateNotification(wallet.RescanUpdate{
		Stage:          wallet.RescanProgress,
		WalletID:       progress.WalletID,
		ProgressReport: progress,
	})
}

func (mp *MainPage) OnBlocksRescanEnded(walletID int, err error) {
	mp.UpdateNotification(wallet.RescanUpdate{
		Stage:    wallet.RescanEnded,
		WalletID: walletID,
	})
}

// UpdateNotification sends notification to the notification channel depending on which channel the page uses
func (mp *MainPage) UpdateNotification(signal interface{}) {
	mp.Load.Receiver.NotificationsUpdate <- signal
}

func (mp *MainPage) desktopNotifier(notifier interface{}) {
	proposalNotification := mp.WL.Wallet.ReadBoolConfigValueForKey(load.ProposalNotificationConfigKey)
	var notification string
	switch t := notifier.(type) {
	case wallet.NewTransaction:

		switch t.Transaction.Type {
		case dcrlibwallet.TxTypeRegular:
			if t.Transaction.Direction != dcrlibwallet.TxDirectionReceived {
				return
			}
			// remove trailing zeros from amount and convert to string
			amount := strconv.FormatFloat(dcrlibwallet.AmountCoin(t.Transaction.Amount), 'f', -1, 64)
			notification = fmt.Sprintf("You have received %s DCR", amount)
		case dcrlibwallet.TxTypeVote:
			reward := strconv.FormatFloat(dcrlibwallet.AmountCoin(t.Transaction.VoteReward), 'f', -1, 64)
			notification = fmt.Sprintf("A ticket just voted\nVote reward: %s DCR", reward)
		case dcrlibwallet.TxTypeRevocation:
			notification = "A ticket was revoked"
		default:
			return
		}

		if mp.WL.MultiWallet.OpenedWalletsCount() > 1 {
			wallet := mp.WL.MultiWallet.WalletWithID(t.Transaction.WalletID)
			if wallet == nil {
				return
			}

			notification = fmt.Sprintf("[%s] %s", wallet.Name, notification)
		}

		initializeBeepNotification(notification)
	case wallet.Proposal:
		switch {
		case t.ProposalStatus == wallet.NewProposalFound:
			notification = fmt.Sprintf("A new proposal has been added Token: %s", t.Proposal.Token)
		case t.ProposalStatus == wallet.VoteStarted:
			notification = fmt.Sprintf("Voting has started for proposal with Token: %s", t.Proposal.Token)
		case t.ProposalStatus == wallet.VoteFinished:
			notification = fmt.Sprintf("Voting has ended for proposal with Token: %s", t.Proposal.Token)
		default:
			notification = fmt.Sprintf("New update for proposal with Token: %s", t.Proposal.Token)
		}
		if proposalNotification {
			initializeBeepNotification(notification)
		}
	}
}

func initializeBeepNotification(n string) {
	absoluteWdPath, err := GetAbsolutePath()
	if err != nil {
		log.Error(err.Error())
	}

	err = beeep.Notify("Decred Godcr Wallet", n, filepath.Join(absoluteWdPath, "ui/assets/decredicons/qrcodeSymbol.png"))
	if err != nil {
		log.Info("could not initiate desktop notification, reason:", err.Error())
	}
}
