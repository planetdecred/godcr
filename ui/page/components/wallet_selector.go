package components

import (
	"context"
	"fmt"
	"sync"

	"gioui.org/layout"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletSelectorID = "wallet_selector"

type badWalletListItem struct {
	*dcrlibwallet.Wallet
	deleteBtn decredmaterial.Button
}

type WalletSelector struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	*listeners.TxAndBlockNotificationListener
	ctx context.Context

	listLock sync.Mutex

	mainWalletList      []*load.WalletItem
	watchOnlyWalletList []*load.WalletItem
	badWalletsList      []*badWalletListItem

	shadowBox            *decredmaterial.Shadow
	walletsList          *decredmaterial.ClickableList
	watchOnlyWalletsList *decredmaterial.ClickableList

	wallectSelected func()
}

func NewWalletSelector(l *load.Load, onWalletSelected func()) *WalletSelector {
	ws := &WalletSelector{
		GenericPageModal: app.NewGenericPageModal(WalletSelectorID),
		wallectSelected:  onWalletSelected,
		Load:             l,
		shadowBox:        l.Theme.Shadow(),
	}

	ws.walletsList = l.Theme.NewClickableList(layout.Vertical)
	ws.watchOnlyWalletsList = l.Theme.NewClickableList(layout.Vertical)

	return ws
}

func (ws *WalletSelector) Expose(ctx context.Context) {
	ws.ctx = ctx
	ws.listenForTxNotifications()
	ws.loadWallets()

	if ws.WL.MultiWallet.ReadBoolConfigValueForKey(load.AutoSyncConfigKey, false) {
		ws.startSyncing()
	}
}

func (ws *WalletSelector) loadWallets() {
	wallets := ws.WL.SortedWalletList()
	mainWalletList := make([]*load.WalletItem, 0)
	watchOnlyWalletList := make([]*load.WalletItem, 0)

	for _, wal := range wallets {
		accountsResult, err := wal.GetAccountsRaw()
		if err != nil {
			continue
		}

		var totalBalance int64
		for _, acc := range accountsResult.Acc {
			totalBalance += acc.TotalBalance
		}

		// sort wallets into normal wallet and watchonly wallets
		if wal.IsWatchingOnlyWallet() {
			listItem := &load.WalletItem{
				Wallet:       wal,
				TotalBalance: dcrutil.Amount(totalBalance).String(),
			}

			watchOnlyWalletList = append(watchOnlyWalletList, listItem)
		} else {
			listItem := &load.WalletItem{
				Wallet:       wal,
				TotalBalance: dcrutil.Amount(totalBalance).String(),
			}

			mainWalletList = append(mainWalletList, listItem)
		}
	}

	ws.listLock.Lock()
	ws.mainWalletList = mainWalletList
	ws.watchOnlyWalletList = watchOnlyWalletList
	ws.listLock.Unlock()

	ws.loadBadWallets()
}

func (ws *WalletSelector) loadBadWallets() {
	badWallets := ws.WL.MultiWallet.BadWallets()
	ws.badWalletsList = make([]*badWalletListItem, 0, len(badWallets))
	for _, badWallet := range badWallets {
		listItem := &badWalletListItem{
			Wallet:    badWallet,
			deleteBtn: ws.Theme.OutlineButton(values.String(values.StrDeleted)),
		}
		listItem.deleteBtn.Color = ws.Theme.Color.Danger
		listItem.deleteBtn.Inset = layout.Inset{}
		ws.badWalletsList = append(ws.badWalletsList, listItem)
	}
}

func (ws *WalletSelector) deleteBadWallet(badWalletID int) {
	warningModal := modal.NewInfoModal(ws.Load).
		Title(values.String(values.StrRemoveWallet)).
		Body(values.String(values.StrWalletRestoreMsg)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButtonStyle(ws.Load.Theme.Color.Surface, ws.Load.Theme.Color.Danger).
		PositiveButton(values.String(values.StrRemove), func(isChecked bool) bool {
			go func() {
				err := ws.WL.MultiWallet.DeleteBadWallet(badWalletID)
				if err != nil {
					errorModal := modal.NewErrorModal(ws.Load, err.Error(), func(isChecked bool) bool {
						return true
					})
					ws.ParentWindow().ShowModal(errorModal)
					return
				}
				infoModal := modal.NewSuccessModal(ws.Load, values.String(values.StrWalletRemoved), func(isChecked bool) bool {
					return true
				})
				ws.ParentWindow().ShowModal(infoModal)
				ws.loadBadWallets() // refresh bad wallets list
				ws.ParentWindow().Reload()
			}()
			return true
		})
	ws.ParentWindow().ShowModal(warningModal)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (ws *WalletSelector) HandleUserInteractions() {

	ws.listLock.Lock()
	mainWalletList := ws.mainWalletList
	watchOnlyWalletList := ws.watchOnlyWalletList
	ws.listLock.Unlock()

	if ok, selectedItem := ws.walletsList.ItemClicked(); ok {
		ws.WL.SelectedWallet = mainWalletList[selectedItem]
		ws.wallectSelected()
	}

	if ok, selectedItem := ws.watchOnlyWalletsList.ItemClicked(); ok {
		ws.WL.SelectedWallet = watchOnlyWalletList[selectedItem]
		ws.wallectSelected()
	}

	for _, badWallet := range ws.badWalletsList {
		if badWallet.deleteBtn.Clicked() {
			ws.deleteBadWallet(badWallet.ID)
		}
	}
}

func (ws *WalletSelector) syncStatusIcon(gtx C) D {
	var (
		syncStatusIcon *decredmaterial.Image
		syncStatus     string
	)

	switch {
	case ws.WL.MultiWallet.IsSynced():
		syncStatusIcon = ws.Theme.Icons.SuccessIcon
		syncStatus = values.String(values.StrSynced)
	case ws.WL.MultiWallet.IsSyncing():
		syncStatusIcon = ws.Theme.Icons.SyncingIcon
		syncStatus = values.String(values.StrSyncingState)
	default:
		syncStatusIcon = ws.Theme.Icons.NotSynced
		syncStatus = values.String(values.StrWalletNotSynced)
	}

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(syncStatusIcon.Layout16dp),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Left: values.MarginPadding5,
			}.Layout(gtx, ws.Theme.Label(values.TextSize16, syncStatus).Layout)
		}),
	)
}

func (ws *WalletSelector) WalletListLayout(gtx C) D {
	walletSections := []func(gtx C) D{
		ws.walletList,
	}

	if len(ws.watchOnlyWalletList) != 0 {
		walletSections = append(walletSections, ws.watchOnlyWalletSection)
	}

	if len(ws.badWalletsList) != 0 {
		walletSections = append(walletSections, ws.badWalletsSection)
	}
	list := &layout.List{
		Axis: layout.Vertical,
	}
	return list.Layout(gtx, len(walletSections), func(gtx C, i int) D {
		return walletSections[i](gtx)
	})
}

func (ws *WalletSelector) walletList(gtx C) D {
	ws.listLock.Lock()
	mainWalletList := ws.mainWalletList
	ws.listLock.Unlock()

	return ws.walletsList.Layout(gtx, len(mainWalletList), func(gtx C, i int) D {
		return ws.walletWrapper(gtx, mainWalletList[i], false)
	})
}

func (ws *WalletSelector) watchOnlyWalletSection(gtx C) D {
	ws.listLock.Lock()
	watchOnlyWalletList := ws.watchOnlyWalletList
	ws.listLock.Unlock()

	return ws.watchOnlyWalletsList.Layout(gtx, len(watchOnlyWalletList), func(gtx C, i int) D {
		return ws.walletWrapper(gtx, watchOnlyWalletList[i], true)
	})
}

func (ws *WalletSelector) badWalletsSection(gtx layout.Context) layout.Dimensions {
	m20 := values.MarginPadding20
	m10 := values.MarginPadding10

	layoutBadWallet := func(gtx C, badWallet *badWalletListItem, lastItem bool) D {
		return layout.Inset{Top: m10, Bottom: m10}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(ws.Theme.Body2(badWallet.Name).Layout),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, badWallet.deleteBtn.Layout)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					if lastItem {
						return D{}
					}
					return layout.Inset{Top: values.MarginPadding10, Left: values.MarginPadding38, Right: values.MarginPaddingMinus10}.Layout(gtx, func(gtx C) D {
						return ws.Theme.Separator().Layout(gtx)
					})
				}),
			)
		})
	}

	card := ws.Theme.Card()
	card.Color = ws.Theme.Color.Surface
	card.Radius = decredmaterial.Radius(10)

	sectionTitleLabel := ws.Theme.Body1("Bad Wallets") // TODO: localize string
	sectionTitleLabel.Color = ws.Theme.Color.GrayText2

	return card.Layout(gtx, func(gtx C) D {
		return layout.Inset{Top: m20, Left: m20}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(sectionTitleLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: m10, Bottom: m10}.Layout(gtx, ws.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return ws.Theme.NewClickableList(layout.Vertical).Layout(gtx, len(ws.badWalletsList), func(gtx C, i int) D {
							return layoutBadWallet(gtx, ws.badWalletsList[i], i == len(ws.badWalletsList)-1)
						})
					})
				}),
			)
		})
	})
}

func (ws *WalletSelector) walletWrapper(gtx C, item *load.WalletItem, isWatchingOnlyWallet bool) D {
	ws.shadowBox.SetShadowRadius(14)
	return decredmaterial.LinearLayout{
		Width:      decredmaterial.WrapContent,
		Height:     decredmaterial.WrapContent,
		Padding:    layout.UniformInset(values.MarginPadding9),
		Background: ws.Theme.Color.Surface,
		Alignment:  layout.Middle,
		Shadow:     ws.shadowBox,
		Margin:     layout.UniformInset(values.MarginPadding5),
		Border:     decredmaterial.Border{Radius: decredmaterial.Radius(14)},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding10,
				Left:  values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				if isWatchingOnlyWallet {
					return ws.Theme.Icons.DcrWatchOnly.Layout36dp(gtx)
				}
				return ws.Theme.Icons.DecredSymbol2.LayoutSize(gtx, values.MarginPadding30)
			})
		}),
		layout.Rigid(ws.Theme.Label(values.TextSize16, item.Wallet.Name).Layout),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if len(item.Wallet.EncryptedSeed) > 0 {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(ws.Theme.Icons.RedAlert.Layout16dp),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Right: values.MarginPadding10,
									}.Layout(gtx, ws.Theme.Label(values.TextSize16, values.String(values.StrNotBackedUp)).Layout)
								}),
							)
						}
						return D{}
					}),
					layout.Rigid(ws.syncStatusIcon),
				)
			})
		}),
	)
}

func (ws *WalletSelector) listenForTxNotifications() {
	if ws.TxAndBlockNotificationListener != nil {
		return
	}
	ws.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err := ws.WL.MultiWallet.AddTxAndBlockNotificationListener(ws.TxAndBlockNotificationListener, true, WalletSelectorID)
	if err != nil {
		// log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	go func() {
		for {
			select {
			case n := <-ws.TxAndBlockNotifChan:
				switch n.Type {
				case listeners.BlockAttached:
					// refresh wallet account and balance on every new block
					// only if sync is completed.
					if ws.WL.MultiWallet.IsSynced() {
						ws.updateAccountBalance()
						ws.ParentWindow().Reload()
					}
				case listeners.NewTransaction:
					// refresh wallets when new transaction is received
					ws.updateAccountBalance()
					ws.ParentWindow().Reload()
				}
			case <-ws.ctx.Done():
				ws.WL.MultiWallet.RemoveTxAndBlockNotificationListener(WalletSelectorID)
				close(ws.TxAndBlockNotifChan)
				ws.TxAndBlockNotificationListener = nil
				return
			}
		}
	}()
}

func (ws *WalletSelector) updateAccountBalance() {
	ws.listLock.Lock()
	defer ws.listLock.Unlock()

	// update main wallets balance
	for _, item := range ws.mainWalletList {
		wal := ws.WL.MultiWallet.WalletWithID(item.Wallet.ID)
		if wal != nil {
			accountsResult, err := wal.GetAccountsRaw()
			if err != nil {
				continue
			}

			var totalBalance int64
			for _, acc := range accountsResult.Acc {
				totalBalance += acc.TotalBalance
			}

			item.TotalBalance = dcrutil.Amount(totalBalance).String()
		}
	}

	// update watch only wallets balance
	for _, item := range ws.watchOnlyWalletList {
		wal := ws.WL.MultiWallet.WalletWithID(item.Wallet.ID)
		if wal != nil {
			accountsResult, err := wal.GetAccountsRaw()
			if err != nil {
				continue
			}

			var totalBalance int64
			for _, acc := range accountsResult.Acc {
				totalBalance += acc.TotalBalance
			}

			item.TotalBalance = dcrutil.Amount(totalBalance).String()
		}
	}
}

func (ws *WalletSelector) startSyncing() {
	for _, wal := range ws.WL.SortedWalletList() {
		if !wal.HasDiscoveredAccounts && wal.IsLocked() {
			ws.UnlockWalletForSyncing(wal)
			return
		}
	}

	err := ws.WL.MultiWallet.SpvSync()
	if err != nil {
		// show error dialog
		log.Info("Error starting sync:", err)
	}
}

func (ws *WalletSelector) UnlockWalletForSyncing(wal *dcrlibwallet.Wallet) {
	spendingPasswordModal := modal.NewPasswordModal(ws.Load).
		Title(values.String(values.StrResumeAccountDiscoveryTitle)).
		Hint(values.String(values.StrSpendingPassword)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(values.String(values.StrUnlock), func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := ws.WL.MultiWallet.UnlockWallet(wal.ID, []byte(password))
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
				ws.startSyncing()
			}()

			return false
		})
	fmt.Println("==================>", ws.ParentWindow())
	fmt.Println("==================>", spendingPasswordModal)
	ws.ParentWindow().ShowModal(spendingPasswordModal)
}
