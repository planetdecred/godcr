package page

import (
	"fmt"
	"strconv"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/dexclient"
	"github.com/planetdecred/godcr/ui/page/governance"
	"github.com/planetdecred/godcr/ui/page/overview"
	"github.com/planetdecred/godcr/ui/page/send"
	"github.com/planetdecred/godcr/ui/page/staking"
	"github.com/planetdecred/godcr/ui/page/transaction"
	"github.com/planetdecred/godcr/ui/page/wallets"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	MainPageID = "Main"
)

var (
	NavDrawerWidth          = unit.Value{U: unit.UnitDp, V: 160}
	NavDrawerMinimizedWidth = unit.Value{U: unit.UnitDp, V: 72}
)

type HideBalanceItem struct {
	hideBalanceButton decredmaterial.IconButton
	tooltip           *decredmaterial.Tooltip
	tooltipLabel      decredmaterial.Label
}

type NavHandler struct {
	Clickable     *widget.Clickable
	Image         *decredmaterial.Image
	ImageInactive *decredmaterial.Image
	Title         string
	PageID        string
}

type MainPage struct {
	*load.Load
	appBarNav components.NavDrawer
	drawerNav components.NavDrawer

	hideBalanceItem HideBalanceItem

	currentPage   load.Page
	pageBackStack []load.Page
	sendPage      *send.Page   // reuse value to keep data persistent onresume.
	receivePage   *ReceivePage // pointer to receive page. to avoid duplication.

	// page state variables
	usdExchangeSet  bool
	dcrUsdtBittrex  load.DCRUSDTBittrex
	isBalanceHidden bool
	totalBalance    dcrutil.Amount
	totalBalanceUSD string
}

func NewMainPage(l *load.Load) *MainPage {
	mp := &MainPage{
		Load: l,
	}

	mp.hideBalanceItem.hideBalanceButton = mp.Theme.IconButton(mp.Icons.ConcealIcon)
	mp.hideBalanceItem.hideBalanceButton.Size = unit.Dp(19)
	mp.hideBalanceItem.hideBalanceButton.Inset = layout.UniformInset(values.MarginPadding4)
	mp.hideBalanceItem.tooltip = mp.Theme.Tooltip()

	// init shared page functions
	toggleSync := func() {
		if mp.WL.MultiWallet.IsConnectedToDecredNetwork() {
			mp.WL.MultiWallet.CancelSync()
		} else {
			mp.StartSyncing()
		}
	}
	l.ToggleSync = toggleSync
	l.ChangeFragment = mp.changeFragment
	l.PopFragment = mp.popFragment
	l.PopToFragment = mp.popToFragment

	mp.initNavItems()

	// Show a seed backup prompt if necessary.
	mp.WL.Wallet.SaveConfigValueForKey(load.SeedBackupNotificationConfigKey, false)

	mp.drawerNav.DrawerToggled(false)

	return mp
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (mp *MainPage) ID() string {
	return MainPageID
}

func (mp *MainPage) initNavItems() {
	mp.appBarNav = components.NavDrawer{
		Load:        mp.Load,
		CurrentPage: mp.currentPageID(),
		AppBarNavItems: []components.NavHandler{
			{
				Clickable: mp.Theme.NewClickable(true),
				Image:     mp.Icons.SendIcon,
				Title:     values.String(values.StrSend),
				PageID:    send.PageID,
			},
			{
				Clickable: mp.Theme.NewClickable(true),
				Image:     mp.Icons.ReceiveIcon,
				Title:     values.String(values.StrReceive),
				PageID:    ReceivePageID,
			},
		},
	}

	mp.drawerNav = components.NavDrawer{
		Load:        mp.Load,
		CurrentPage: mp.currentPageID(),
		DrawerNavItems: []components.NavHandler{
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Icons.OverviewIcon,
				ImageInactive: mp.Icons.OverviewIconInactive,
				Title:         values.String(values.StrOverview),
				PageID:        overview.OverviewPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Icons.TransactionsIcon,
				ImageInactive: mp.Icons.TransactionsIconInactive,
				Title:         values.String(values.StrTransactions),
				PageID:        transaction.TransactionsPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Icons.WalletIcon,
				ImageInactive: mp.Icons.WalletIconInactive,
				Title:         values.String(values.StrWallets),
				PageID:        wallets.WalletPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Icons.StakeIcon,
				ImageInactive: mp.Icons.StakeIconInactive,
				Title:         values.String(values.StrStaking),
				PageID:        staking.OverviewPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Icons.GovernanceActiveIcon,
				ImageInactive: mp.Icons.GovernanceInactiveIcon,
				Title:         "Governance",
				PageID:        governance.ProposalsPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Icons.DexIcon,
				ImageInactive: mp.Icons.DexIconInactive,
				Title:         values.String(values.StrDex),
				PageID:        dexclient.MarketPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Icons.MoreIcon,
				ImageInactive: mp.Icons.MoreIconInactive,
				Title:         values.String(values.StrMore),
				PageID:        MorePageID,
			},
		},
		MinimizeNavDrawerButton: mp.Theme.IconButton(mp.Icons.NavigationArrowBack),
		MaximizeNavDrawerButton: mp.Theme.IconButton(mp.Icons.NavigationArrowForward),
	}
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (mp *MainPage) OnNavigatedTo() {
	// register for notifications, unregister when the page disappears
	mp.WL.MultiWallet.AddAccountMixerNotificationListener(mp, MainPageID)
	mp.WL.MultiWallet.Politeia.AddNotificationListener(mp, MainPageID)
	mp.WL.MultiWallet.AddTxAndBlockNotificationListener(mp, true, MainPageID) // notification methods will be invoked asynchronously to prevent potential deadlocks
	mp.WL.MultiWallet.AddSyncProgressListener(mp, MainPageID)
	mp.WL.MultiWallet.SetBlocksRescanProgressListener(mp)

	mp.setLanguageSetting()

	if mp.currentPage == nil {
		mp.currentPage = overview.NewOverviewPage(mp.Load)
	}
	mp.currentPage.OnNavigatedTo()

	if mp.sendPage != nil {
		mp.sendPage.OnNavigatedTo()
	}
	if mp.receivePage != nil {
		mp.receivePage.OnNavigatedTo()
	}

	if mp.WL.Wallet.ReadBoolConfigValueForKey(load.AutoSyncConfigKey) {
		mp.StartSyncing()
		if mp.WL.Wallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey) {
			go mp.WL.MultiWallet.Politeia.Sync()
		}
	}
}

func (mp *MainPage) setLanguageSetting() {
	langPre := mp.WL.Wallet.ReadStringConfigValueForKey(load.LanguagePreferenceKey)
	values.SetUserLanguage(langPre)
}

func (mp *MainPage) UpdateBalance() {
	go load.GetUSDExchangeValue(&mp.dcrUsdtBittrex)
	currencyExchangeValue := mp.WL.Wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	mp.usdExchangeSet = currencyExchangeValue == values.USDExchangeValue

	totalBalance, err := mp.CalculateTotalWalletsBalance()
	if err == nil {
		mp.totalBalance = totalBalance

		if mp.usdExchangeSet && mp.dcrUsdtBittrex.LastTradeRate != "" {
			usdExchangeRate, err := strconv.ParseFloat(mp.dcrUsdtBittrex.LastTradeRate, 64)
			if err == nil {
				balanceInUSD := totalBalance.ToCoin() * usdExchangeRate
				mp.totalBalanceUSD = load.FormatUSDBalance(mp.Printer, balanceInUSD)
			}
		}
	}
}

func (mp *MainPage) CalculateTotalWalletsBalance() (dcrutil.Amount, error) {
	totalBalance := int64(0)
	wallets := mp.WL.SortedWalletList()
	if len(wallets) == 0 {
		return 0, nil
	}

	for _, wallet := range wallets {
		accountsResult, err := wallet.GetAccountsRaw()
		if err != nil {
			return 0, err
		}

		for _, account := range accountsResult.Acc {
			totalBalance += account.TotalBalance
		}
	}

	return dcrutil.Amount(totalBalance), nil
}

func (mp *MainPage) StartSyncing() {
	for _, wal := range mp.WL.SortedWalletList() {
		if !wal.HasDiscoveredAccounts && wal.IsLocked() {
			mp.UnlockWalletForSyncing(wal)
			return
		}
	}

	err := mp.WL.MultiWallet.SpvSync()
	if err != nil {
		// show error dialog
		log.Info("Error starting sync:", err)
	}
}

func (mp *MainPage) UnlockWalletForSyncing(wal *dcrlibwallet.Wallet) {
	modal.NewPasswordModal(mp.Load).
		Title(values.String(values.StrResumeAccountDiscoveryTitle)).
		Hint("Spending password").
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(values.String(values.StrUnlock), func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := mp.WL.MultiWallet.UnlockWallet(wal.ID, []byte(password))
				if err != nil {
					errText := err.Error()
					if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
						errText = "Invalid password"
					}
					pm.SetError(errText)
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
				mp.StartSyncing()
			}()

			return false
		}).Show()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (mp *MainPage) HandleUserInteractions() {
	if mp.currentPage != nil {
		mp.currentPage.HandleUserInteractions()
	}

	mp.drawerNav.CurrentPage = mp.currentPageID()
	mp.appBarNav.CurrentPage = mp.currentPageID()

	for mp.drawerNav.MinimizeNavDrawerButton.Button.Clicked() {
		mp.drawerNav.DrawerToggled(true)
	}

	for mp.drawerNav.MaximizeNavDrawerButton.Button.Clicked() {
		mp.drawerNav.DrawerToggled(false)
	}

	for i, item := range mp.appBarNav.AppBarNavItems {
		for item.Clickable.Clicked() {
			var pg load.Page
			if i == 0 {
				if mp.sendPage == nil {
					mp.sendPage = send.NewSendPage(mp.Load)
				}
				pg = mp.sendPage
			} else {
				if mp.receivePage == nil {
					mp.receivePage = NewReceivePage(mp.Load)
				}
				pg = mp.receivePage
			}

			if pg.ID() == mp.currentPageID() {
				continue
			}

			mp.ChangeFragment(pg)
		}
	}

	for _, item := range mp.drawerNav.DrawerNavItems {
		for item.Clickable.Clicked() {
			var pg load.Page
			switch item.PageID {
			case overview.OverviewPageID:
				pg = overview.NewOverviewPage(mp.Load)
			case transaction.TransactionsPageID:
				pg = transaction.NewTransactionsPage(mp.Load)
			case wallets.WalletPageID:
				pg = wallets.NewWalletPage(mp.Load)
			case staking.OverviewPageID:
				pg = staking.NewStakingPage(mp.Load)
			case governance.ProposalsPageID:
				pg = governance.NewProposalsPage(mp.Load)
			case dexclient.MarketPageID:
				_, err := mp.WL.MultiWallet.StartDexClient() // does nothing if already started
				if err != nil {
					mp.Toast.NotifyError(fmt.Sprintf("Unable to start DEX client: %v", err))
				} else {
					pg = dexclient.NewMarketPage(mp.Load)
				}
			case MorePageID:
				pg = NewMorePage(mp.Load)
			}

			if pg == nil || pg.ID() == mp.currentPageID() {
				continue
			}

			// clear stack
			mp.changeFragment(pg)
		}
	}

	mp.isBalanceHidden = mp.WL.MultiWallet.ReadBoolConfigValueForKey(load.HideBalanceConfigKey, false)
	for mp.hideBalanceItem.hideBalanceButton.Button.Clicked() {
		mp.isBalanceHidden = !mp.isBalanceHidden
		mp.WL.MultiWallet.SetBoolConfigValueForKey(load.HideBalanceConfigKey, mp.isBalanceHidden)
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (mp *MainPage) OnNavigatedFrom() {
	// Also disappear all child pages.
	if mp.currentPage != nil {
		mp.currentPage.OnNavigatedFrom()
	}
	if mp.sendPage != nil {
		mp.sendPage.OnNavigatedFrom()
	}
	if mp.receivePage != nil {
		mp.receivePage.OnNavigatedFrom()
	}

	mp.WL.MultiWallet.RemoveAccountMixerNotificationListener(MainPageID)
	mp.WL.MultiWallet.Politeia.RemoveNotificationListener(MainPageID)
	mp.WL.MultiWallet.RemoveTxAndBlockNotificationListener(MainPageID)
	mp.WL.MultiWallet.RemoveSyncProgressListener(MainPageID)
}

func (mp *MainPage) currentPageID() string {
	if mp.currentPage != nil {
		return mp.currentPage.ID()
	}

	return ""
}

func (mp *MainPage) changeFragment(page load.Page) {

	// If Page is the last in back stack return.
	if mp.currentPageID() == page.ID() {
		return
	}

	// Maintain one pointer to Page in backstack slice.
	for i := len(mp.pageBackStack) - 1; i >= 0; i-- {
		if mp.pageBackStack[i].ID() == page.ID() {
			var mPages []load.Page
			mPagesf, mPagesb := mp.pageBackStack[:i], mp.pageBackStack[i+1:]
			mPages = append(mPages, mPagesf...)
			mPages = append(mPages, mPagesb...)
			mp.pageBackStack = mPages
		}
	}

	if mp.currentPage != nil {
		mp.currentPage.OnNavigatedFrom() // TODO: Unload unless it is possible that this page will be revisited.
		mp.pageBackStack = append(mp.pageBackStack, mp.currentPage)
	}

	page.OnNavigatedTo()
	mp.currentPage = page
}

// popFragment goes back to the previous page
func (mp *MainPage) popFragment() {
	if len(mp.pageBackStack) > 0 {
		// get and remove last page
		previousPage := mp.pageBackStack[len(mp.pageBackStack)-1]
		mp.pageBackStack = mp.pageBackStack[:len(mp.pageBackStack)-1]

		mp.currentPage.OnNavigatedFrom()
		previousPage.OnNavigatedTo()
		mp.currentPage = previousPage
	}
}

func (mp *MainPage) popToFragment(pageID string) {
	// close current page and all pages before `pageID`
	if mp.currentPage != nil {
		mp.currentPage.OnNavigatedFrom()
	}

	for i := len(mp.pageBackStack) - 1; i >= 0; i-- {
		if mp.pageBackStack[i].ID() == pageID {
			var closedPages []load.Page
			mp.pageBackStack, closedPages = mp.pageBackStack[:i+1], mp.pageBackStack[i+1:]

			for j := len(closedPages) - 1; j >= 0; j-- {
				closedPages[j].OnNavigatedFrom()
			}
			break
		}
	}

	if len(mp.pageBackStack) > 0 {
		// set curent page to `pageID`
		mp.currentPage = mp.pageBackStack[len(mp.pageBackStack)-1]
		// remove current page from backstack history
		mp.pageBackStack = mp.pageBackStack[:len(mp.pageBackStack)-1]
	} else {
		mp.currentPage = nil
	}
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (mp *MainPage) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.MatchParent,
				Orientation: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(mp.LayoutTopBar),
				layout.Rigid(func(gtx C) D {
					return decredmaterial.LinearLayout{
						Width:       decredmaterial.MatchParent,
						Height:      decredmaterial.MatchParent,
						Orientation: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(mp.drawerNav.LayoutNavDrawer),
						layout.Rigid(func(gtx C) D {
							if mp.currentPage == nil {
								return layout.Dimensions{}
							}
							return mp.currentPage.Layout(gtx)
						}),
					)
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			// TODO: hidden balance tip hover layout
			// if mp.hideBalanceItem.hideBalanceButton.Button.Hovered() {
			// 	lm := values.MarginPadding280
			// 	if mp.hideBalanceItem.tooltipLabel.Text == "Show Balance" {
			// 		lm = values.MarginPadding168
			// 	}

			// 	return layout.Inset{Top: values.MarginPadding50, Left: lm}.Layout(gtx, func(gtx C) D {
			// 		card := mp.Theme.Card()
			// 		card.Color = mp.Theme.Color.Surface
			// 		card.Border = true
			// 		card.Radius = decredmaterial.Radius(5)
			// 		card.BorderParam.CornerRadius = values.MarginPadding5
			// 		return card.Layout(gtx, func(gtx C) D {
			// 			return components.Container{
			// 				Padding: layout.UniformInset(values.MarginPadding5),
			// 			}.Layout(gtx, mp.hideBalanceItem.tooltipLabel.Layout)
			// 		})
			// 	})
			// }

			// global toasts. Stack toast on all pages and contents
			//TODO: show toasts here
			return D{}

		}),
	)
}

func (mp *MainPage) LayoutUSDBalance(gtx layout.Context) layout.Dimensions {
	mp.UpdateBalance()
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if mp.usdExchangeSet && mp.dcrUsdtBittrex.LastTradeRate != "" && len(mp.totalBalanceUSD) > 0 {
				inset := layout.Inset{
					Top:  values.MarginPadding3,
					Left: values.MarginPadding8,
				}
				border := widget.Border{Color: mp.Theme.Color.Gray2, CornerRadius: unit.Dp(8), Width: unit.Dp(0.5)}
				return inset.Layout(gtx, func(gtx C) D {
					padding := layout.Inset{
						Top:    values.MarginPadding3,
						Bottom: values.MarginPadding3,
						Left:   values.MarginPadding6,
						Right:  values.MarginPadding6,
					}
					return border.Layout(gtx, func(gtx C) D {
						return padding.Layout(gtx, func(gtx C) D {
							return mp.Theme.Body2(mp.totalBalanceUSD).Layout(gtx)
						})
					})
				})
			}
			return D{}
		}),
	)
}

func (mp *MainPage) totalDCRBalance(gtx layout.Context) layout.Dimensions {
	if mp.isBalanceHidden {
		hiddenBalanceText := mp.Theme.Label(values.TextSize18.Scale(0.8), "**********DCR")
		return layout.Inset{Bottom: values.MarginPadding0, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			return hiddenBalanceText.Layout(gtx)
		})
	}
	return components.LayoutBalance(gtx, mp.Load, mp.totalBalance.String())
}

func (mp *MainPage) LayoutTopBar(gtx layout.Context) layout.Dimensions {
	return decredmaterial.LinearLayout{
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Background:  mp.Theme.Color.Surface,
		Orientation: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.WrapContent,
				Background:  mp.Theme.Color.Surface,
				Orientation: layout.Horizontal,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						h := values.MarginPadding24
						v := values.MarginPadding14
						// Balance container
						return components.Container{Padding: layout.Inset{Right: h, Left: h, Top: v, Bottom: v}}.Layout(gtx,
							func(gtx C) D {
								return decredmaterial.LinearLayout{
									Width:       decredmaterial.WrapContent,
									Height:      decredmaterial.WrapContent,
									Background:  mp.Theme.Color.Surface,
									Orientation: layout.Horizontal,
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										img := mp.Icons.Logo
										if mp.WL.MultiWallet.ReadBoolConfigValueForKey(load.DarkModeConfigKey, false) {
											img = mp.Icons.LogoDarkMode
										}

										return layout.Inset{Right: values.MarginPadding16}.Layout(gtx,
											func(gtx C) D {
												return img.Layout24dp(gtx)
											})
									}),
									layout.Rigid(func(gtx C) D {
										return mp.totalDCRBalance(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										if !mp.isBalanceHidden {
											return mp.LayoutUSDBalance(gtx)
										}
										return layout.Dimensions{}
									}),
									layout.Rigid(func(gtx C) D {
										mp.hideBalanceItem.tooltipLabel = mp.Theme.Caption("Hide Balance")
										mp.hideBalanceItem.hideBalanceButton.Icon = mp.Icons.RevealIcon
										if mp.isBalanceHidden {
											mp.hideBalanceItem.tooltipLabel.Text = "Show Balance"
											mp.hideBalanceItem.hideBalanceButton.Icon = mp.Icons.ConcealIcon
										}
										return layout.Inset{
											Top:  values.MarginPadding1,
											Left: values.MarginPadding9,
										}.Layout(gtx, mp.hideBalanceItem.hideBalanceButton.Layout)
									}),
								)
							})
					})
				}),
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return mp.appBarNav.LayoutTopBar(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return mp.Theme.Separator().Layout(gtx)
		}),
	)
}
