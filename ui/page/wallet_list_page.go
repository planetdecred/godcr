package page

import (
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletListID = "wallet_list"

type walletListItem struct {
	wal          *dcrlibwallet.Wallet
	totalBalance string
}

type WalletList struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	listLock        sync.Mutex
	scrollContainer *widget.List
	listItems       []*walletListItem
	shadowBox       *decredmaterial.Shadow

	walletsList          *decredmaterial.ClickableList
	watchOnlyWalletsList *decredmaterial.ClickableList

	addWalClickable *decredmaterial.Clickable

	hasWatchOnly bool
}

func NewWalletList(l *load.Load) *WalletList {
	pg := &WalletList{
		GenericPageModal: app.NewGenericPageModal(WalletListID),
		scrollContainer: &widget.List{
			List: layout.List{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			},
		},

		Load:      l,
		shadowBox: l.Theme.Shadow(),
	}

	pg.walletsList = l.Theme.NewClickableList(layout.Vertical)
	pg.walletsList.DividerHeight = values.MarginPadding4
	pg.walletsList.IsShadowEnabled = true

	pg.watchOnlyWalletsList = l.Theme.NewClickableList(layout.Vertical)
	pg.watchOnlyWalletsList.DividerHeight = values.MarginPadding4
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
	pg.hasWatchOnly = false
	pg.loadWallets()
}

func (pg *WalletList) loadWallets() {
	wallets := pg.WL.SortedWalletList()
	listItems := make([]*walletListItem, 0)

	for _, wal := range wallets {
		if wal.IsWatchingOnlyWallet() {
			pg.hasWatchOnly = true
		}

		accountsResult, err := wal.GetAccountsRaw()
		if err != nil {
			continue
		}

		var totalBalance int64
		for _, acc := range accountsResult.Acc {
			totalBalance += acc.TotalBalance
		}

		listItem := &walletListItem{
			wal:          wal,
			totalBalance: dcrutil.Amount(totalBalance).String(),
		}

		listItems = append(listItems, listItem)
	}

	pg.listLock.Lock()
	pg.listItems = listItems
	pg.listLock.Unlock()

	// pg.loadBadWallets()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *WalletList) HandleUserInteractions() {}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *WalletList) OnNavigatedFrom() {}

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
	var syncStatusIcon *decredmaterial.Image
	switch {
	case pg.WL.MultiWallet.IsSynced():
		syncStatusIcon = pg.Theme.Icons.SuccessIcon
	case pg.WL.MultiWallet.IsSyncing():
		syncStatusIcon = pg.Theme.Icons.SyncingIcon
	default:
		syncStatusIcon = pg.Theme.Icons.FailedIcon
	}
	if pg.WL.MultiWallet.IsSynced() {
		syncStatusIcon = pg.Theme.Icons.SuccessIcon
	}
	return syncStatusIcon.Layout24dp(gtx)
}

func (pg *WalletList) walletSection(gtx C) D {
	walletSections := []func(gtx C) D{
		pg.walletList,
	}

	if pg.hasWatchOnly {
		walletSections = append(walletSections, pg.watchOnlyWalletList)
	}

	// if len(pg.badWalletsList) != 0 {
	// 	pg = append(pg, pg.badWalletsSection)
	// }

	return pg.Theme.List(pg.scrollContainer).Layout(gtx, len(walletSections), func(gtx C, i int) D {
		return walletSections[i](gtx)
	})
}

func (pg *WalletList) walletList(gtx C) D {
	pg.listLock.Lock()
	listItems := pg.listItems
	pg.listLock.Unlock()

	return pg.walletsList.Layout(gtx, len(listItems), func(gtx C, i int) D {
		item := listItems[i]

		if item.wal.IsWatchingOnlyWallet() {
			return D{}
		}
		return pg.walletWrapper(gtx, item, false)
	})
}

func (pg *WalletList) watchOnlyWalletList(gtx C) D {
	pg.listLock.Lock()
	listItems := pg.listItems
	pg.listLock.Unlock()

	return pg.watchOnlyWalletsList.Layout(gtx, len(listItems), func(gtx C, i int) D {
		item := listItems[i]

		if !item.wal.IsWatchingOnlyWallet() {
			return D{}
		}
		return pg.walletWrapper(gtx, item, true)
	})
}

func (pg *WalletList) walletWrapper(gtx C, item *walletListItem, isWatchingOnlyWallet bool) D {
	return decredmaterial.LinearLayout{
		Width:      decredmaterial.WrapContent,
		Height:     decredmaterial.WrapContent,
		Padding:    layout.UniformInset(values.MarginPadding9),
		Background: pg.Theme.Color.Surface,
		Alignment:  layout.Middle,
		Border:     decredmaterial.Border{Radius: decredmaterial.Radius(14)},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding10,
				Left:  values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				if isWatchingOnlyWallet {
					return pg.Theme.Icons.DcrWatchOnly.Layout24dp(gtx)
				}
				return pg.Theme.Icons.DecredSymbol2.Layout24dp(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.Theme.Body1(item.wal.Name).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(pg.syncStatusIcon),
						layout.Rigid(func(gtx C) D {
							var txt decredmaterial.Label
							if len(item.wal.EncryptedSeed) > 0 {
								txt = pg.Theme.Caption(values.String(values.StrNotBackedUp))
								txt.Color = pg.Theme.Color.Danger
								return txt.Layout(gtx)
							}
							return D{}
						}),
					)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			balanceLabel := pg.Theme.Body1(item.totalBalance)
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
