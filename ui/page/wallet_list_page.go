package page

import (
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletListID = "wallet_list"

type walletListItem struct {
	wal          *dcrlibwallet.Wallet
	totalBalance string
}

type WalletList struct {
	*load.Load
	listLock        sync.Mutex
	scrollContainer *widget.List
	listItems       []*walletListItem
	shadowBox       *decredmaterial.Shadow

	walletsList     *decredmaterial.ClickableList
	addWalClickable *decredmaterial.Clickable
}

func NewWalletList(l *load.Load) *WalletList {
	pg := &WalletList{
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
	pg.walletsList.DividerHeight = values.MarginPadding8

	pg.addWalClickable = l.Theme.NewClickable(false)
	pg.addWalClickable.Radius = decredmaterial.Radius(10)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *WalletList) ID() string {
	return WalletListID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *WalletList) OnNavigatedTo() {
	pg.loadWallets()
}

func (pg *WalletList) loadWallets() {
	wallets := pg.WL.SortedWalletList()
	listItems := make([]*walletListItem, 0)

	for _, wal := range wallets {
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

func (pg *WalletList) walletList(gtx C, item *walletListItem) D {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding10,
				Left:  values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				return pg.Theme.Icons.WalletIcon.LayoutSize(gtx, values.MarginPadding26)
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

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *WalletList) Layout(gtx C) D {
	gtx.Constraints.Min = gtx.Constraints.Max // use maximum height & width

	pg.listLock.Lock()
	listItems := pg.listItems
	pg.listLock.Unlock()
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

	return layout.Center.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding550)

		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.walletsList.Layout(gtx, len(listItems), func(gtx C, i int) D {
					return decredmaterial.LinearLayout{
						Width:      decredmaterial.WrapContent,
						Height:     decredmaterial.WrapContent,
						Padding:    layout.UniformInset(values.MarginPadding9),
						Background: pg.Theme.Color.Surface,
						Clickable:  pg.addWalClickable,
						Shadow:     pg.shadowBox,
						Alignment:  layout.Middle,
						Direction:  layout.Center,
						Border:     decredmaterial.Border{Radius: decredmaterial.Radius(14)},
					}.Layout2(gtx, func(gtx C) D {
						return pg.walletList(gtx, listItems[i])
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, pg.layoutAddWalletSection)
						})
					}),
				)
			}),
		)
	})
}

func (pg *WalletList) layoutAddWalletSection(gtx C) D {
	// gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	// return layout.SE.Layout(gtx, func(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
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
	// })
}
