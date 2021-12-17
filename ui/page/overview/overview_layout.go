package overview

import (
	"context"
	"sync"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	gPage "github.com/planetdecred/godcr/ui/page/governance"
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
	mixerWallets []*dcrlibwallet.Wallet

	transactions  []dcrlibwallet.Transaction
	proposalItems []*components.ProposalItem
	proposalMu    sync.Mutex

	bestBlock    *dcrlibwallet.BlockInfo
	rescanUpdate *wallet.RescanUpdate

	scrollContainer        *widget.List
	proposalsListContainer *widget.List
	walletSyncList         *layout.List
	listContainer          *layout.List
	listMixer              *layout.List

	syncClickable    *decredmaterial.Clickable
	transactionsList *decredmaterial.ClickableList
	proposalsList    *decredmaterial.ClickableList

	syncingIcon                  *decredmaterial.Image
	walletStatusIcon, cachedIcon *decredmaterial.Icon
	syncedIcon, notSyncedIcon    *decredmaterial.Icon

	toTransactions decredmaterial.TextAndIconButton
	toProposals    decredmaterial.TextAndIconButton
	toMixer        decredmaterial.IconButton

	sync              decredmaterial.Label
	toggleSyncDetails decredmaterial.Button
	checkBox          decredmaterial.CheckBoxStyle

	walletSyncing         bool
	walletSynced          bool
	isConnnected          bool
	isBackupModalOpened   bool
	rescanningBlocks      bool
	syncDetailsVisibility bool
	syncCompleted         bool

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

		listContainer: &layout.List{Axis: layout.Vertical},
		listMixer:     &layout.List{Axis: layout.Vertical},
		scrollContainer: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		checkBox:  l.Theme.CheckBox(new(widget.Bool), "I am aware of the risk"),
		bestBlock: l.WL.MultiWallet.GetBestBlock(),
	}

	pg.toMixer = l.Theme.IconButton(l.Icons.NavigationArrowForward)
	pg.toMixer.Size = values.MarginPadding24
	pg.toMixer.Inset = layout.UniformInset(values.MarginPadding4)

	pg.initRecentTxWidgets()
	pg.initWalletStatusWidgets()
	pg.initSyncDetailsWidgets()
	pg.initializeProposalsWidget()

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *AppOverviewPage) ID() string {
	return OverviewPageID
}

// WillAppear is called when the page is about to displayed and may
// be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *AppOverviewPage) WillAppear() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.walletSyncing = pg.WL.MultiWallet.IsSyncing()
	pg.walletSynced = pg.WL.MultiWallet.IsSynced()
	pg.isConnnected = pg.WL.MultiWallet.IsConnectedToDecredNetwork()
	pg.connectedPeers = pg.WL.MultiWallet.ConnectedPeers()
	pg.bestBlock = pg.WL.MultiWallet.GetBestBlock()

	pg.getMixerWallets()
	pg.loadTransactions()
	pg.listenForSyncNotifications()
	pg.loadRecentProposals()
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
// Layout lays out the entire content for overview pg.
func (pg *AppOverviewPage) Layout(gtx layout.Context) layout.Dimensions {
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			if len(pg.mixerWallets) == 0 {
				return D{}
			}

			return components.MixerInfoLayout(gtx, pg.Load, true, pg.toMixer.Layout, func(gtx C) D {
				return pg.listMixer.Layout(gtx, len(pg.mixerWallets), func(gtx C, i int) D {
					return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						accounts, _ := pg.mixerWallets[i].GetAccountsRaw()
						var unmixedBalance string
						for _, acct := range accounts.Acc {
							if acct.Number == pg.mixerWallets[i].UnmixedAccountNumber() {
								unmixedBalance = dcrutil.Amount(acct.TotalBalance).String()
							}
						}

						return components.MixerInfoContentWrapper(gtx, pg.Load, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									txt := pg.Theme.Label(values.TextSize14, pg.mixerWallets[i].Name)
									txt.Font.Weight = text.Medium

									return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, txt.Layout)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											t := pg.Theme.Label(values.TextSize14, "Unmixed balance")
											t.Color = pg.Theme.Color.GrayText2
											return t.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return components.LayoutBalanceSize(gtx, pg.Load, unmixedBalance, values.TextSize20)
										}),
									)
								}),
							)
						})
					})
				})
			})
		},
		func(gtx C) D {
			return pg.recentTransactionsSection(gtx)
		},
		func(gtx C) D {
			if pg.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
				return pg.recentProposalsSection(gtx)
			}
			return D{}
		},
		func(gtx C) D {
			return pg.syncStatusSection(gtx)
		},
	}

	return components.UniformPadding(gtx, func(gtx C) D {
		return pg.Theme.List(pg.scrollContainer).Layout(gtx, len(pageContent), func(gtx C, i int) D {
			m := values.MarginPadding5
			if i == len(pageContent) {
				// remove padding after the last item
				m = values.MarginPadding0
			}
			return layout.Inset{
				Right:  values.MarginPadding2,
				Bottom: m,
			}.Layout(gtx, pageContent[i])
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
			pg.WL.Wallet.SaveConfigValueForKey(load.SeedBackupNotificationConfigKey, true)
		}).
		PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
		PositiveButton("Backup now", func() {
			pg.WL.Wallet.SaveConfigValueForKey(load.SeedBackupNotificationConfigKey, true)
			pg.ChangeFragment(wPage.NewWalletPage(pg.Load))
		}).Show()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *AppOverviewPage) HandleUserInteractions() {
	backupLater := pg.WL.Wallet.ReadBoolConfigValueForKey(load.SeedBackupNotificationConfigKey)
	for _, wal := range pg.allWallets {
		if len(wal.EncryptedSeed) > 0 {
			if !backupLater && !pg.isBackupModalOpened {
				pg.showBackupInfo()
				pg.isBackupModalOpened = true
			}
		}
	}

	if pg.toMixer.Button.Clicked() {
		if len(pg.mixerWallets) == 1 {
			pg.ChangeFragment(wPage.NewPrivacyPage(pg.Load, pg.mixerWallets[0]))
		}
		pg.ChangeFragment(wPage.NewWalletPage(pg.Load))
	}

	if pg.syncClickable.Clicked() {
		if pg.rescanningBlocks {
			pg.WL.MultiWallet.CancelRescan()
		} else {
			// If connected to the Decred network disable button. Prevents multiple clicks.
			if pg.isConnnected {
				pg.syncClickable.SetEnabled(false, nil)
			}

			// On exit update button state.
			go func() {
				pg.ToggleSync()
				if !pg.syncClickable.Enabled() {
					pg.syncClickable.SetEnabled(true, nil)
				}
			}()
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

	for pg.toProposals.Button.Clicked() {
		pg.ChangeFragment(gPage.NewProposalsPage(pg.Load))
	}

	if clicked, selectedItem := pg.proposalsList.ItemClicked(); clicked {
		pg.proposalMu.Lock()
		selectedProposal := pg.proposalItems[selectedItem].Proposal
		pg.proposalMu.Unlock()

		pg.ChangeFragment(gPage.NewProposalDetailsPage(pg.Load, &selectedProposal))
	}

	if pg.syncCompleted {
		pg.syncCompleted = false
		pg.loadRecentProposals()
		pg.RefreshWindow()
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
			case wallet.Proposal:
				if n.ProposalStatus == wallet.Synced {
					pg.syncCompleted = true
					pg.RefreshWindow()
				}
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

func (pg *AppOverviewPage) getMixerWallets() {
	wallets := make([]*dcrlibwallet.Wallet, 0)

	for _, wal := range pg.allWallets {
		if wal.IsAccountMixerActive() {
			wallets = append(wallets, wal)
		}
	}

	pg.mixerWallets = wallets
}

// WillDisappear is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// WillAppear() will be called again. This method should not destroy UI
// components unless they'll be recreated in the WillAppear() method.
// Part of the load.Page interface.
func (pg *AppOverviewPage) WillDisappear() {
	pg.ctxCancel()
}
