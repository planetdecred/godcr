package info

import (
	"context"
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	// "github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/seedbackup"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const InfoID = "Info"

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

type WalletInfo struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	*listeners.SyncProgressListener
	*listeners.BlocksRescanProgressListener
	*listeners.TxAndBlockNotificationListener
	ctx       context.Context // page context
	ctxCancel context.CancelFunc
	listLock  sync.Mutex

	multiWallet  *dcrlibwallet.MultiWallet
	rescanUpdate *wallet.RescanUpdate

	container *widget.List

	walletStatusIcon *decredmaterial.Icon
	syncSwitch       *decredmaterial.Switch
	toBackup         decredmaterial.Button
	checkBox         decredmaterial.CheckBoxStyle

	remainingSyncTime    string
	syncStepLabel        string
	headersToFetchOrScan int32
	stepFetchProgress    int32
	syncProgress         int
	syncStep             int
	isBackupModalOpened  bool
}

func NewInfoPage(l *load.Load) *WalletInfo {
	pg := &WalletInfo{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(InfoID),
		multiWallet:      l.WL.MultiWallet,
		container: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		checkBox: l.Theme.CheckBox(new(widget.Bool), "I am aware of the risk"),
	}

	pg.toBackup = pg.Theme.Button(values.String(values.StrBackupNow))
	pg.toBackup.Color = pg.Theme.Color.Primary
	pg.toBackup.TextSize = values.TextSize14
	pg.toBackup.Background = color.NRGBA{}

	pg.initWalletStatusWidgets()

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *WalletInfo) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	// backupLater := pg.WL.SelectedWallet.Wallet.ReadBoolConfigValueForKey(load.SeedBackupNotificationConfigKey, false)
	// needBackup := pg.WL.SelectedWallet.Wallet.EncryptedSeed != nil
	// if needBackup && !backupLater && !pg.isBackupModalOpened {
	// 	pg.showBackupInfo()
	// 	pg.isBackupModalOpened = true
	// }

	autoSync := pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.AutoSyncConfigKey, false)
	pg.syncSwitch.SetChecked(autoSync)

	pg.listenForNotifications()
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
// Layout lays out the widgets for the main wallets pg.
func (pg *WalletInfo) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		return pg.Theme.List(pg.container).Layout(gtx, 1, func(gtx C, i int) D {
			return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Right: values.MarginPadding10,
									Left:  values.MarginPadding10,
								}.Layout(gtx, func(gtx C) D {
									txt := pg.Theme.Body1(pg.WL.SelectedWallet.Wallet.Name)
									txt.Font.Weight = text.SemiBold
									return txt.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								if len(pg.WL.SelectedWallet.Wallet.EncryptedSeed) > 0 {
									return layout.Inset{
										Top: values.MarginPadding16,
									}.Layout(gtx, func(gtx C) D {
										return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
											layout.Rigid(pg.Theme.Icons.RedAlert.Layout24dp),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{
													Left: values.MarginPadding9,
												}.Layout(gtx, pg.Theme.Body2("Backup your wallet now to avoid losses").Layout)
											}),
											layout.Rigid(pg.toBackup.Layout),
										)
									})
								}
								return D{}
							}),
							layout.Rigid(pg.syncStatusSection),
						)
					})
				})
			})
		})
	}

	return components.UniformPadding(gtx, body)
}

// func (pg *WalletInfo) showBackupInfo() {
// 	backupNowOrLaterModal := modal.NewInfoModal(pg.Load).
// 		SetupWithTemplate(modal.WalletBackupInfoTemplate).
// 		SetCancelable(false).
// 		SetContentAlignment(layout.W, layout.Center).
// 		CheckBox(pg.checkBox, true).
// 		NegativeButton(values.String(values.StrBackupLater), func() {
// 			pg.WL.SelectedWallet.Wallet.SaveUserConfigValue(load.SeedBackupNotificationConfigKey, true)
// 		}).
// 		PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
// 		PositiveButton(values.String(values.StrBackupNow), func(isChecked bool) bool {
// 			pg.WL.SelectedWallet.Wallet.SaveUserConfigValue(load.SeedBackupNotificationConfigKey, true)
// 			pg.ParentNavigator().Display(seedbackup.NewBackupInstructionsPage(pg.Load, pg.WL.SelectedWallet.Wallet))
// 			return true
// 		})
// 	pg.ParentWindow().ShowModal(backupNowOrLaterModal)
// }

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *WalletInfo) HandleUserInteractions() {
	// backupLater := pg.WL.SelectedWallet.Wallet.ReadBoolConfigValueForKey(load.SeedBackupNotificationConfigKey, false)
	// needBackup := pg.WL.SelectedWallet.Wallet.EncryptedSeed != nil
	// if needBackup && !backupLater && !pg.isBackupModalOpened {
	// 	pg.showBackupInfo()
	// 	pg.isBackupModalOpened = true
	// }

	if pg.syncSwitch.Changed() {
		if pg.WL.MultiWallet.IsRescanning() {
			pg.WL.MultiWallet.CancelRescan()
		} else {
			pg.WL.MultiWallet.SaveUserConfigValue(load.AutoSyncConfigKey, pg.syncSwitch.IsChecked())
			go func() {
				pg.ToggleSync()
			}()
		}
	}

	for pg.toBackup.Button.Clicked() {
		pg.ParentNavigator().Display(seedbackup.NewBackupInstructionsPage(pg.Load, pg.WL.SelectedWallet.Wallet))
	}
}

// listenForNotifications starts a goroutine to watch for sync updates
// and update the UI accordingly. To prevent UI lags, this method does not
// refresh the window display everytime a sync update is received. During
// active blocks sync, rescan or proposals sync, the Layout method auto
// refreshes the display every set interval. Other sync updates that affect
// the UI but occur outside of an active sync requires a display refresh.
func (pg *WalletInfo) listenForNotifications() {
	switch {
	case pg.SyncProgressListener != nil:
		return
	case pg.TxAndBlockNotificationListener != nil:
		return
	case pg.BlocksRescanProgressListener != nil:
		return
	}

	pg.SyncProgressListener = listeners.NewSyncProgress()
	err := pg.WL.MultiWallet.AddSyncProgressListener(pg.SyncProgressListener, InfoID)
	if err != nil {
		log.Errorf("Error adding sync progress listener: %v", err)
		return
	}

	pg.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err = pg.WL.MultiWallet.AddTxAndBlockNotificationListener(pg.TxAndBlockNotificationListener, true, InfoID)
	if err != nil {
		log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	pg.BlocksRescanProgressListener = listeners.NewBlocksRescanProgressListener()
	pg.WL.MultiWallet.SetBlocksRescanProgressListener(pg.BlocksRescanProgressListener)

	go func() {
		for {
			select {
			case n := <-pg.SyncStatusChan:
				// Update sync progress fields which will be displayed
				// when the next UI invalidation occurs.
				switch t := n.ProgressReport.(type) {
				case *dcrlibwallet.HeadersFetchProgressReport:
					pg.stepFetchProgress = t.HeadersFetchProgress
					pg.headersToFetchOrScan = t.TotalHeadersToFetch
					pg.syncProgress = int(t.TotalSyncProgress)
					pg.remainingSyncTime = components.TimeFormat(int(t.TotalTimeRemainingSeconds), true)
					pg.syncStep = wallet.FetchHeadersSteps
				case *dcrlibwallet.AddressDiscoveryProgressReport:
					pg.syncProgress = int(t.TotalSyncProgress)
					pg.remainingSyncTime = components.TimeFormat(int(t.TotalTimeRemainingSeconds), true)
					pg.syncStep = wallet.AddressDiscoveryStep
					pg.stepFetchProgress = t.AddressDiscoveryProgress
				case *dcrlibwallet.HeadersRescanProgressReport:
					pg.headersToFetchOrScan = t.TotalHeadersToScan
					pg.syncProgress = int(t.TotalSyncProgress)
					pg.remainingSyncTime = components.TimeFormat(int(t.TotalTimeRemainingSeconds), true)
					pg.syncStep = wallet.RescanHeadersStep
					pg.stepFetchProgress = t.RescanProgress
				}

				// We only care about sync state changes here, to
				// refresh the window display.
				switch n.Stage {
				case wallet.SyncStarted:
					fallthrough
				case wallet.SyncCanceled:
					fallthrough
				case wallet.SyncCompleted:
					pg.ParentWindow().Reload()
				}

			case n := <-pg.TxAndBlockNotifChan:
				switch n.Type {
				case listeners.NewTransaction:
					pg.ParentWindow().Reload()
				case listeners.BlockAttached:
					pg.ParentWindow().Reload()
				}
			case n := <-pg.BlockRescanChan:
				pg.rescanUpdate = &n
				if n.Stage == wallet.RescanEnded {
					pg.ParentWindow().Reload()
				}
			case <-pg.ctx.Done():
				pg.WL.MultiWallet.RemoveSyncProgressListener(InfoID)
				pg.WL.MultiWallet.RemoveTxAndBlockNotificationListener(InfoID)
				pg.WL.MultiWallet.SetBlocksRescanProgressListener(nil)

				close(pg.SyncStatusChan)
				close(pg.TxAndBlockNotifChan)
				close(pg.BlockRescanChan)

				pg.SyncProgressListener = nil
				pg.TxAndBlockNotificationListener = nil
				pg.BlocksRescanProgressListener = nil

				return
			}
		}
	}()
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *WalletInfo) OnNavigatedFrom() {
	pg.ctxCancel()
}
