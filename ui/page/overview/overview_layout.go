package overview

import (
	"context"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	tPage "github.com/planetdecred/godcr/ui/page/transaction"
	wPage "github.com/planetdecred/godcr/ui/page/wallets"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const OverviewPageID = "Overview"

type (
	C = layout.Context
	D = layout.Dimensions
)

// walletSyncDetails contains sync data for each wallet when a sync
// is in progress.
type walletSyncDetails struct {
	name               decredmaterial.Label
	status             decredmaterial.Label
	blockHeaderFetched decredmaterial.Label
	syncingProgress    decredmaterial.Label
}

type AppOverviewPage struct {
	*load.Load
	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	allWallets   []*dcrlibwallet.Wallet
	transactions []dcrlibwallet.Transaction
	bestBlock    *dcrlibwallet.BlockInfo
	rescanUpdate *wallet.RescanUpdate

	listContainer  *widget.List
	walletSyncList *layout.List

	syncClickable                *decredmaterial.Clickable
	transactionsList             *decredmaterial.ClickableList
	syncingIcon                  *decredmaterial.Image
	walletStatusIcon, cachedIcon *decredmaterial.Icon
	syncedIcon, notSyncedIcon    *decredmaterial.Icon

	toTransactions    decredmaterial.TextAndIconButton
	sync              decredmaterial.Label
	toggleSyncDetails decredmaterial.Button
	checkBox          decredmaterial.CheckBoxStyle

	walletSyncing         bool
	walletSynced          bool
	isConnnected          bool
	isBackupModalOpened   bool
	rescanningBlocks      bool
	syncDetailsVisibility bool

	remainingSyncTime    string
	headersToFetchOrScan int32
	headerFetchProgress  int32
	connectedPeers       int32
	syncProgress         int
	syncStep             int
}

func NewOverviewPage(l *load.Load) *AppOverviewPage {
	pg := &AppOverviewPage{
		Load:       l,
		allWallets: l.WL.SortedWalletList(),

		listContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		checkBox:  l.Theme.CheckBox(new(widget.Bool), "I am aware of the risk"),
		bestBlock: l.WL.MultiWallet.GetBestBlock(),
	}

	pg.initRecentTxWidgets()
	pg.initWalletStatusWidgets()
	pg.initSyncDetailsWidgets()

	return pg
}

func (pg *AppOverviewPage) ID() string {
	return OverviewPageID
}

func (pg *AppOverviewPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.walletSyncing = pg.WL.MultiWallet.IsSyncing()
	pg.walletSynced = pg.WL.MultiWallet.IsSynced()
	pg.isConnnected = pg.WL.MultiWallet.IsConnectedToDecredNetwork()
	pg.connectedPeers = pg.WL.MultiWallet.ConnectedPeers()
	pg.bestBlock = pg.WL.MultiWallet.GetBestBlock()

	pg.loadTransactions()
	pg.listenForSyncNotifications()
}

// Layout lays out the entire content for overview pg.
func (pg *AppOverviewPage) Layout(gtx layout.Context) layout.Dimensions {
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.recentTransactionsSection(gtx)
		},
		func(gtx C) D {
			return pg.syncStatusSection(gtx)
		},
	}

	return components.UniformPadding(gtx, func(gtx C) D {
		return pg.Theme.List(pg.listContainer).Layout(gtx, len(pageContent), func(gtx C, i int) D {
			return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, pageContent[i])
		})
	})
}

func (pg *AppOverviewPage) showBackupInfo() {
	modal.NewInfoModal(pg.Load).
		SetupWithTemplate(modal.WalletBackupInfoTemplate).
		SetCancelable(false).
		SetContentAlignment(layout.W, layout.Center).
		CheckBox(pg.checkBox).
		NegativeButton("Backup later", func() {
			pg.WL.Wallet.SaveConfigValueForKey("seedBackupNotification", true)
		}).
		PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
		PositiveButton("Backup now", func() {
			pg.WL.Wallet.SaveConfigValueForKey("seedBackupNotification", true)
			pg.ChangeFragment(wPage.NewWalletPage(pg.Load))
		}).Show()
}

func (pg *AppOverviewPage) Handle() {
	backupLater := pg.WL.Wallet.ReadBoolConfigValueForKey("seedBackupNotification")
	for _, wal := range pg.allWallets {
		if len(wal.EncryptedSeed) > 0 {
			if !backupLater && !pg.isBackupModalOpened {
				pg.showBackupInfo()
				pg.isBackupModalOpened = true
			}
		}
	}

	if pg.syncClickable.Clicked() {
		if pg.rescanningBlocks {
			pg.WL.MultiWallet.CancelRescan()
		} else {
			go pg.ToggleSync()
		}
	}

	if pg.toTransactions.Button.Clicked() {
		pg.ChangeFragment(tPage.NewTransactionsPage(pg.Load))
	}

	if clicked, selectedItem := pg.transactionsList.ItemClicked(); clicked {
		pg.ChangeFragment(tPage.NewTransactionDetailsPage(pg.Load, &pg.transactions[selectedItem]))
	}

	if pg.toggleSyncDetails.Clicked() {
		pg.syncDetailsVisibility = !pg.syncDetailsVisibility
		if pg.syncDetailsVisibility {
			pg.toggleSyncDetails.Text = values.String(values.StrHideDetails)
		} else {
			pg.toggleSyncDetails.Text = values.String(values.StrShowDetails)
		}
	}
}

func (pg *AppOverviewPage) listenForSyncNotifications() {
	go func() {
		for {
			var notification interface{}

			select {
			case notification = <-pg.Receiver.NotificationsUpdate:
			case <-pg.ctx.Done():
				return
			}

			switch n := notification.(type) {
			case wallet.NewTransaction:
				pg.loadTransactions()
			case wallet.RescanUpdate:
				pg.rescanningBlocks = n.Stage != wallet.RescanEnded
				pg.rescanUpdate = &n
			case wallet.SyncStatusUpdate:
				switch t := n.ProgressReport.(type) {
				case wallet.SyncHeadersFetchProgress:
					pg.headerFetchProgress = t.Progress.HeadersFetchProgress
					pg.headersToFetchOrScan = t.Progress.TotalHeadersToFetch
					pg.syncProgress = int(t.Progress.TotalSyncProgress)
					pg.remainingSyncTime = components.TimeFormat(int(t.Progress.TotalTimeRemainingSeconds), true)
					pg.syncStep = wallet.FetchHeadersSteps
				case wallet.SyncAddressDiscoveryProgress:
					pg.syncProgress = int(t.Progress.TotalSyncProgress)
					pg.remainingSyncTime = components.TimeFormat(int(t.Progress.TotalTimeRemainingSeconds), true)
					pg.syncStep = wallet.AddressDiscoveryStep
				case wallet.SyncHeadersRescanProgress:
					pg.headersToFetchOrScan = t.Progress.TotalHeadersToScan
					pg.syncProgress = int(t.Progress.TotalSyncProgress)
					pg.remainingSyncTime = components.TimeFormat(int(t.Progress.TotalTimeRemainingSeconds), true)
					pg.syncStep = wallet.RescanHeadersStep
				}

				switch n.Stage {
				case wallet.PeersConnected:
					pg.connectedPeers = n.ConnectedPeers
				case wallet.SyncStarted:
					fallthrough
				case wallet.SyncCanceled:
					fallthrough
				case wallet.SyncCompleted:
					pg.loadTransactions()
					pg.walletSyncing = pg.WL.MultiWallet.IsSyncing()
					pg.walletSynced = pg.WL.MultiWallet.IsSynced()
					pg.isConnnected = pg.WL.MultiWallet.IsConnectedToDecredNetwork()
				case wallet.BlockAttached:
					pg.bestBlock = pg.WL.MultiWallet.GetBestBlock()
				}
			}

			pg.RefreshWindow()
		}
	}()
}

func (pg *AppOverviewPage) OnClose() {
	pg.ctxCancel()
}
