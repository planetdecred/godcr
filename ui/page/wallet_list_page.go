package page

import (
	"context"
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletListID = "wallet_list"

type badWalletListItem struct {
	*dcrlibwallet.Wallet
	deleteBtn decredmaterial.Button
}

type WalletList struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	*listeners.TxAndBlockNotificationListener
	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	listLock        sync.Mutex
	scrollContainer *widget.List

	mainWalletList      []*load.WalletItem
	watchOnlyWalletList []*load.WalletItem
	badWalletsList      []*badWalletListItem

	shadowBox            *decredmaterial.Shadow
	walletsList          *decredmaterial.ClickableList
	watchOnlyWalletsList *decredmaterial.ClickableList
	addWalClickable      *decredmaterial.Clickable

	wallectSelected func()
}

func NewWalletList(l *load.Load, onWalletSelected func()) *WalletList {
	pg := &WalletList{
		GenericPageModal: app.NewGenericPageModal(WalletListID),
		scrollContainer: &widget.List{
			List: layout.List{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			},
		},
		wallectSelected: onWalletSelected,
		Load:      l,
		shadowBox: l.Theme.Shadow(),
	}

	pg.walletsList = l.Theme.NewClickableList(layout.Vertical)
	pg.walletsList.IsShadowEnabled = true

	pg.watchOnlyWalletsList = l.Theme.NewClickableList(layout.Vertical)
	pg.watchOnlyWalletsList.IsShadowEnabled = true

	pg.addWalClickable = l.Theme.NewClickable(false)
	pg.addWalClickable.Radius = decredmaterial.Radius(14)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *WalletList) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.listenForTxNotifications()
	pg.loadWallets()
}

func (pg *WalletList) loadWallets() {
	wallets := pg.WL.SortedWalletList()
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

	pg.listLock.Lock()
	pg.mainWalletList = mainWalletList
	pg.watchOnlyWalletList = watchOnlyWalletList
	pg.listLock.Unlock()

	pg.loadBadWallets()
}

func (pg *WalletList) loadBadWallets() {
	badWallets := pg.WL.MultiWallet.BadWallets()
	pg.badWalletsList = make([]*badWalletListItem, 0, len(badWallets))
	for _, badWallet := range badWallets {
		listItem := &badWalletListItem{
			Wallet:    badWallet,
			deleteBtn: pg.Theme.OutlineButton(values.String(values.StrDeleted)),
		}
		listItem.deleteBtn.Color = pg.Theme.Color.Danger
		listItem.deleteBtn.Inset = layout.Inset{}
		pg.badWalletsList = append(pg.badWalletsList, listItem)
	}
}

func (pg *WalletList) deleteBadWallet(badWalletID int) {
	warningModal := modal.NewInfoModal(pg.Load).
		Title(values.String(values.StrRemoveWallet)).
		Body(values.String(values.StrWalletRestoreMsg)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButtonStyle(pg.Load.Theme.Color.Surface, pg.Load.Theme.Color.Danger).
		PositiveButton(values.String(values.StrRemove), func(isChecked bool) bool {
			go func() {
				err := pg.WL.MultiWallet.DeleteBadWallet(badWalletID)
				if err != nil {
					pg.Toast.NotifyError(err.Error())
					return
				}
				pg.Toast.Notify(values.String(values.StrWalletRemoved))
				pg.loadBadWallets() // refresh bad wallets list
				pg.ParentWindow().Reload()
			}()
			return true
		})
	pg.ParentWindow().ShowModal(warningModal)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *WalletList) HandleUserInteractions() {

	pg.listLock.Lock()
	mainWalletList := pg.mainWalletList
	watchOnlyWalletList := pg.watchOnlyWalletList
	pg.listLock.Unlock()

	if ok, selectedItem := pg.walletsList.ItemClicked(); ok {
		pg.WL.SelectedWallet = mainWalletList[selectedItem]
		pg.wallectSelected()
	}

	if ok, selectedItem := pg.watchOnlyWalletsList.ItemClicked(); ok {
		pg.WL.SelectedWallet = watchOnlyWalletList[selectedItem]
		pg.wallectSelected()
	}

	for _, badWallet := range pg.badWalletsList {
		if badWallet.deleteBtn.Clicked() {
			pg.deleteBadWallet(badWallet.ID)
		}
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *WalletList) OnNavigatedFrom() {
	pg.ctxCancel()
}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *WalletList) Layout(gtx C) D {
	pageContent := []func(gtx C) D{
		pg.Theme.Label(values.TextSize20, values.String(values.StrSelectWalletToOpen)).Layout,
		pg.walletSection, // wallet list layout
		func(gtx C) D {
			return layout.Inset{
				Left:   values.MarginPadding5,
				Bottom: values.MarginPadding10,
			}.Layout(gtx, pg.layoutAddWalletSection)
		},
	}

	gtx.Constraints.Min = gtx.Constraints.Max
	return components.UniformPadding(gtx, func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding550)
		list := &layout.List{
			Axis: layout.Vertical,
		}

		return layout.Center.Layout(gtx, func(gtx C) D {
			return list.Layout(gtx, len(pageContent), func(gtx C, i int) D {
				return layout.Inset{Top: values.MarginPadding26}.Layout(gtx, func(gtx C) D {
					return pageContent[i](gtx)
				})
			})
		})
	})
}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *WalletList) Layout(gtx C) D {
	pageContent := []func(gtx C) D{
		pg.Theme.Label(values.TextSize20, values.String(values.StrSelectWalletToOpen)).Layout,
		pg.walletSection, // wallet list layout
		func(gtx C) D {
			return layout.Inset{
				Left:   values.MarginPadding5,
				Bottom: values.MarginPadding10,
			}.Layout(gtx, pg.layoutAddWalletSection)
		},
	}

	gtx.Constraints.Min = gtx.Constraints.Max
	return components.UniformPadding(gtx, func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding550)
		list := &layout.List{
			Axis: layout.Vertical,
		}

		return layout.Center.Layout(gtx, func(gtx C) D {
			return list.Layout(gtx, len(pageContent), func(gtx C, i int) D {
				return layout.Inset{Top: values.MarginPadding26}.Layout(gtx, func(gtx C) D {
					return pageContent[i](gtx)
				})
			})
		})
	})
}

func (pg *WalletList) syncStatusIcon(gtx C) D {
	var (
		syncStatusIcon *decredmaterial.Image
		syncStatus     string
	)

	switch {
	case pg.WL.MultiWallet.IsSynced():
		syncStatusIcon = pg.Theme.Icons.SuccessIcon
		syncStatus = values.String(values.StrSynced)
	case pg.WL.MultiWallet.IsSyncing():
		syncStatusIcon = pg.Theme.Icons.SyncingIcon
		syncStatus = values.String(values.StrSyncingState)
	default:
		syncStatusIcon = pg.Theme.Icons.FailedIcon
		syncStatus = values.String(values.StrWalletNotSynced)
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(syncStatusIcon.Layout16dp),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Left: values.MarginPadding5,
			}.Layout(gtx, pg.Theme.Caption(syncStatus).Layout)
		}),
	)
}

func (pg *WalletList) walletSection(gtx C) D {
	walletSections := []func(gtx C) D{
		pg.walletList,
	}

	if len(pg.watchOnlyWalletList) != 0 {
		walletSections = append(walletSections, pg.watchOnlyWalletSection)
	}

	if len(pg.badWalletsList) != 0 {
		walletSections = append(walletSections, pg.badWalletsSection)
	}

	return pg.Theme.List(pg.scrollContainer).Layout(gtx, len(walletSections), func(gtx C, i int) D {
		return walletSections[i](gtx)
	})
}

func (pg *WalletList) walletList(gtx C) D {
	pg.listLock.Lock()
	mainWalletList := pg.mainWalletList
	pg.listLock.Unlock()

	return pg.walletsList.Layout(gtx, len(mainWalletList), func(gtx C, i int) D {
		return pg.walletWrapper(gtx, mainWalletList[i], false)
	})
}

func (pg *WalletList) watchOnlyWalletSection(gtx C) D {
	pg.listLock.Lock()
	watchOnlyWalletList := pg.watchOnlyWalletList
	pg.listLock.Unlock()

	return pg.watchOnlyWalletsList.Layout(gtx, len(watchOnlyWalletList), func(gtx C, i int) D {
		return pg.walletWrapper(gtx, watchOnlyWalletList[i], true)
	})
}

func (pg *WalletList) badWalletsSection(gtx layout.Context) layout.Dimensions {
	m20 := values.MarginPadding20
	m10 := values.MarginPadding10

	layoutBadWallet := func(gtx C, badWallet *badWalletListItem, lastItem bool) D {
		return layout.Inset{Top: m10, Bottom: m10}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(pg.Theme.Body2(badWallet.Name).Layout),
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
						return pg.Theme.Separator().Layout(gtx)
					})
				}),
			)
		})
	}

	card := pg.Theme.Card()
	card.Color = pg.Theme.Color.Surface
	card.Radius = decredmaterial.Radius(10)

	sectionTitleLabel := pg.Theme.Body1("Bad Wallets") // TODO: localize string
	sectionTitleLabel.Color = pg.Theme.Color.GrayText2

	return card.Layout(gtx, func(gtx C) D {
		return layout.Inset{Top: m20, Left: m20}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(sectionTitleLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: m10, Bottom: m10}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return pg.Theme.NewClickableList(layout.Vertical).Layout(gtx, len(pg.badWalletsList), func(gtx C, i int) D {
							return layoutBadWallet(gtx, pg.badWalletsList[i], i == len(pg.badWalletsList)-1)
						})
					})
				}),
			)
		})
	})
}

func (pg *WalletList) walletWrapper(gtx C, item *load.WalletItem, isWatchingOnlyWallet bool) D {
	pg.shadowBox.SetShadowRadius(14)
	return decredmaterial.LinearLayout{
		Width:      decredmaterial.WrapContent,
		Height:     decredmaterial.WrapContent,
		Padding:    layout.UniformInset(values.MarginPadding9),
		Background: pg.Theme.Color.Surface,
		Alignment:  layout.Middle,
		Shadow:     pg.shadowBox,
		Margin:     layout.UniformInset(values.MarginPadding5),
		Border:     decredmaterial.Border{Radius: decredmaterial.Radius(14)},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding10,
				Left:  values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				if isWatchingOnlyWallet {
					return pg.Theme.Icons.DcrWatchOnly.Layout36dp(gtx)
				}
				return pg.Theme.Icons.DecredSymbol2.LayoutSize(gtx, values.MarginPadding30)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.Theme.Label(values.TextSize16, item.Wallet.Name).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(pg.syncStatusIcon),
						layout.Rigid(func(gtx C) D {
							if len(item.Wallet.EncryptedSeed) > 0 {
								return layout.Flex{
									Axis:      layout.Horizontal,
									Alignment: layout.Middle,
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										ic := decredmaterial.NewIcon(pg.Theme.Icons.ImageBrightness1)
										ic.Color = pg.Theme.Color.Gray1
										return layout.Inset{
											Left:  values.MarginPadding7,
											Right: values.MarginPadding7,
										}.Layout(gtx, func(gtx C) D {
											return ic.Layout(gtx, values.MarginPadding4)
										})
									}),
									layout.Rigid(pg.Theme.Icons.RedAlert.Layout16dp),
									layout.Rigid(func(gtx C) D {
										txt := pg.Theme.Caption(values.String(values.StrNotBackedUp))
										txt.Color = pg.Theme.Color.Danger
										return layout.Inset{
											Left: values.MarginPadding5,
										}.Layout(gtx, txt.Layout)
									}),
								)
							}
							return D{}
						}),
					)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			balanceLabel := pg.Theme.Body1(item.TotalBalance)
			balanceLabel.Color = pg.Theme.Color.GrayText2
			return layout.Inset{
				Right: values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				return layout.E.Layout(gtx, balanceLabel.Layout)
			})
		}),
	)
}

func (pg *WalletList) layoutAddWalletSection(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			pg.shadowBox.SetShadowRadius(14)

			return decredmaterial.LinearLayout{
				Width:      decredmaterial.WrapContent,
				Height:     decredmaterial.WrapContent,
				Padding:    layout.UniformInset(values.MarginPadding12),
				Background: pg.Theme.Color.Surface,
				Clickable:  pg.addWalClickable,
				Shadow:     pg.shadowBox,
				Border:     decredmaterial.Border{Radius: pg.addWalClickable.Radius},
			}.Layout(gtx,
				layout.Rigid(pg.Theme.Icons.NewWalletIcon.Layout24dp),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Left: values.MarginPadding4,
						Top:  values.MarginPadding2,
					}.Layout(gtx, pg.Theme.Body2(values.String(values.StrAddWallet)).Layout)
				}),
			)
		}),
	)
}

func (pg *WalletList) listenForTxNotifications() {
	if pg.TxAndBlockNotificationListener != nil {
		return
	}
	pg.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err := pg.WL.MultiWallet.AddTxAndBlockNotificationListener(pg.TxAndBlockNotificationListener, true, WalletListID)
	if err != nil {
		log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	go func() {
		for {
			select {
			case n := <-pg.TxAndBlockNotifChan:
				switch n.Type {
				case listeners.BlockAttached:
					// refresh wallet account and balance on every new block
					// only if sync is completed.
					if pg.WL.MultiWallet.IsSynced() {
						pg.updateAccountBalance()
						pg.ParentWindow().Reload()
					}
				case listeners.NewTransaction:
					// refresh wallets when new transaction is received
					pg.updateAccountBalance()
					pg.ParentWindow().Reload()
				}
			case <-pg.ctx.Done():
				pg.WL.MultiWallet.RemoveTxAndBlockNotificationListener(WalletListID)
				close(pg.TxAndBlockNotifChan)
				pg.TxAndBlockNotificationListener = nil

				return
			}
		}
	}()
}

func (pg *WalletList) updateAccountBalance() {
	pg.listLock.Lock()
	defer pg.listLock.Unlock()

	// update main wallets balance
	for _, item := range pg.mainWalletList {
		wal := pg.WL.MultiWallet.WalletWithID(item.Wallet.ID)
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
	for _, item := range pg.watchOnlyWalletList {
		wal := pg.WL.MultiWallet.WalletWithID(item.Wallet.ID)
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
