package wallets

import (
	"context"
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/seedbackup"
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

	multiWallet *dcrlibwallet.MultiWallet

	container *widget.List

	card      decredmaterial.Card
	separator decredmaterial.Line
	accounts  []*dcrlibwallet.Account

	totalBalance string

	accountsList        *decredmaterial.ClickableList
	addAcctClickable    *decredmaterial.Clickable
	backupAcctClickable *decredmaterial.Clickable
	renameWallet        *decredmaterial.Clickable

	toggleSyncDetails            decredmaterial.Button
	syncedIcon, notSyncedIcon    *decredmaterial.Icon
	walletStatusIcon, cachedIcon *decredmaterial.Icon
	syncingIcon                  *decredmaterial.Image
	autoSyncSwitch               *decredmaterial.Switch
	walletSyncList               *layout.List
	syncClickable                *decredmaterial.Clickable
	toBackup                     decredmaterial.Button

	rescanUpdate         *wallet.RescanUpdate
	remainingSyncTime    string
	syncStepLabel        string
	headersToFetchOrScan int32
	stepFetchProgress    int32
	syncProgress         int
	syncStep             int

	syncDetailsVisibility bool
}

func NewWalletPage(l *load.Load) *AppOverviewPage {
	pg := &AppOverviewPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(OverviewPageID),
		multiWallet:      l.WL.MultiWallet,
		container: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		card:         l.Theme.Card(),
		separator:    l.Theme.Separator(),
		renameWallet: l.Theme.NewClickable(false),
	}

	pg.separator.Color = l.Theme.Color.Gray2
	pg.accountsList = pg.Theme.NewClickableList(layout.Vertical)
	pg.addAcctClickable = pg.Theme.NewClickable(false)

	backupClickable := pg.Theme.NewClickable(false)
	backupClickable.ChangeStyle(&values.ClickableStyle{Color: pg.Theme.Color.OrangeRipple})
	backupClickable.Radius = decredmaterial.CornerRadius{BottomRight: 14, BottomLeft: 14}
	pg.backupAcctClickable = backupClickable

	pg.initWalletStatusWidgets()
	pg.initSyncDetailsWidgets()

	pg.toBackup = pg.Theme.Button(values.String(values.StrBackupNow))
	pg.toBackup.Color = pg.Theme.Color.Primary
	pg.toBackup.TextSize = values.TextSize14
	pg.toBackup.Background = color.NRGBA{}
	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *AppOverviewPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.listenForNotifications()
	pg.loadWalletAccounts()

	autoSync := pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.AutoSyncConfigKey, false)
	pg.autoSyncSwitch.SetChecked(autoSync)
}

func (pg *AppOverviewPage) loadWalletAccounts() {
	accountsResult, err := pg.WL.SelectedWallet.Wallet.GetAccountsRaw()
	if err != nil {
		log.Errorf("Wallet account error: %v", err)
		return
	}

	var totalBalance int64
	for _, acc := range accountsResult.Acc {
		totalBalance += acc.TotalBalance
	}

	pg.accounts = accountsResult.Acc
	pg.totalBalance = dcrutil.Amount(totalBalance).String()
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
// Layout lays out the widgets for the main wallets pg.
func (pg *AppOverviewPage) Layout(gtx layout.Context) layout.Dimensions {
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

// func (pg *AppOverviewPage) headerLayout(gtx layout.Context) D {
// 	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
// 		layout.Rigid(pg.Theme.Icons.WalletIcon.Layout24dp),
// 		layout.Rigid(func(gtx C) D {
// 			return layout.Inset{
// 				Right: values.MarginPadding10,
// 				Left:  values.MarginPadding10,
// 			}.Layout(gtx, func(gtx C) D {
// 				return pg.Theme.Body1(pg.WL.SelectedWallet.Wallet.Name).Layout(gtx)
// 			})
// 		}),
// 		layout.Rigid(func(gtx C) D {
// 			return pg.renameWallet.Layout(gtx, pg.Theme.Icons.EditIcon.Layout24dp)
// 		}),
// 		layout.Flexed(1, func(gtx C) D {
// 			return layout.E.Layout(gtx, func(gtx C) D {
// 				balanceLabel := pg.Theme.Body1(pg.totalBalance)
// 				balanceLabel.Color = pg.Theme.Color.GrayText2
// 				return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, balanceLabel.Layout)
// 			})
// 		}),
// 	)
// }

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *AppOverviewPage) HandleUserInteractions() {
	if ok, selectedItem := pg.accountsList.ItemClicked(); ok {
		pg.ParentNavigator().Display(NewAcctDetailsPage(pg.Load, pg.accounts[selectedItem]))
	}

	if pg.syncClickable.Clicked() {
		if pg.WL.MultiWallet.IsRescanning() {
			pg.WL.MultiWallet.CancelRescan()
		} else {
			// If connected to the Decred network disable button. Prevents multiple clicks.
			if pg.WL.MultiWallet.IsConnectedToDecredNetwork() {
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

	if pg.renameWallet.Clicked() {
		textModal := modal.NewTextInputModal(pg.Load).
			Hint(values.String(values.StrWalletName)).
			PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
			PositiveButton(values.String(values.StrRename), func(newName string, tim *modal.TextInputModal) bool {
				err := pg.multiWallet.RenameWallet(pg.WL.SelectedWallet.Wallet.ID, newName)
				if err != nil {
					pg.Toast.NotifyError(err.Error())
					return false
				}
				return true
			})

		textModal.Title(values.String(values.StrRenameWalletSheetTitle)).
			NegativeButton(values.String(values.StrCancel), func() {})
		pg.ParentWindow().ShowModal(textModal)
	}

	if pg.autoSyncSwitch.Changed() {
		pg.WL.MultiWallet.SaveUserConfigValue(load.AutoSyncConfigKey, pg.autoSyncSwitch.IsChecked())
		pg.autoSyncSwitch.SetChecked(pg.autoSyncSwitch.IsChecked())
		if pg.autoSyncSwitch.IsChecked() && (!pg.WL.MultiWallet.IsSyncing() || pg.WL.MultiWallet.IsSynced()) {
			info := modal.NewInfoModal(pg.Load).
				Title(values.String(values.StrAutoSync)).
				Body(values.String(values.StrAutoSyncInfo)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(values.String(values.StrContinue), func(isChecked bool) bool {
					for _, wal := range pg.WL.SortedWalletList() {
						if !wal.HasDiscoveredAccounts && wal.IsLocked() {
							pg.unlockWalletForSyncing(wal)
							return true
						}
					}

					err := pg.WL.MultiWallet.SpvSync()
					if err != nil {
						// show error dialog
						log.Info("Error starting sync:", err)
					}
					return true
				})
			pg.ParentWindow().ShowModal(info)
		}
	}

	if pg.toggleSyncDetails.Clicked() {
		pg.syncDetailsVisibility = !pg.syncDetailsVisibility
		if pg.syncDetailsVisibility {
			pg.toggleSyncDetails.Text = values.String(values.StrHideDetails)
		} else {
			pg.toggleSyncDetails.Text = values.String(values.StrShowDetails)
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
func (pg *AppOverviewPage) listenForNotifications() {
	switch {
	case pg.SyncProgressListener != nil:
		return
	case pg.TxAndBlockNotificationListener != nil:
		return
	case pg.BlocksRescanProgressListener != nil:
		return
	}

	pg.SyncProgressListener = listeners.NewSyncProgress()
	err := pg.WL.MultiWallet.AddSyncProgressListener(pg.SyncProgressListener, OverviewPageID)
	if err != nil {
		log.Errorf("Error adding sync progress listener: %v", err)
		return
	}

	pg.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err = pg.WL.MultiWallet.AddTxAndBlockNotificationListener(pg.TxAndBlockNotificationListener, true, OverviewPageID)
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
				pg.WL.MultiWallet.RemoveSyncProgressListener(OverviewPageID)
				pg.WL.MultiWallet.RemoveTxAndBlockNotificationListener(OverviewPageID)
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

func (pg *AppOverviewPage) unlockWalletForSyncing(wal *dcrlibwallet.Wallet) {
	spendingPasswordModal := modal.NewPasswordModal(pg.Load).
		Title(values.String(values.StrResumeAccountDiscoveryTitle)).
		Hint(values.String(values.StrSpendingPassword)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(values.String(values.StrUnlock), func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := pg.WL.MultiWallet.UnlockWallet(wal.ID, []byte(password))
				if err != nil {
					errText := err.Error()
					if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
						errText = values.String(values.StrInvalidPassphrase)
					}
					pm.SetError(errText)
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()

				for _, wal := range pg.WL.SortedWalletList() {
					if !wal.HasDiscoveredAccounts && wal.IsLocked() {
						pg.unlockWalletForSyncing(wal)
						return
					}
				}

				err = pg.WL.MultiWallet.SpvSync()
				if err != nil {
					// show error dialog
					log.Info("Error starting sync:", err)
				}
			}()

			return false
		})
	pg.ParentWindow().ShowModal(spendingPasswordModal)
}

func (pg *AppOverviewPage) updateAccountBalance() {
	pg.listLock.Lock()
	defer pg.listLock.Unlock()

	wal := pg.WL.MultiWallet.WalletWithID(pg.WL.SelectedWallet.Wallet.ID)
	if wal != nil {
		accountsResult, err := wal.GetAccountsRaw()
		if err != nil {
			log.Errorf("Wallet account error: %v", err)
			return
		}

		var totalBalance int64
		for _, acc := range accountsResult.Acc {
			totalBalance += acc.TotalBalance
		}

		pg.totalBalance = dcrutil.Amount(totalBalance).String()
		pg.accounts = accountsResult.Acc
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *AppOverviewPage) OnNavigatedFrom() {
	pg.ctxCancel()
}
