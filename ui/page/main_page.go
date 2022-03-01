package page

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/gen2brain/beeep"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/listeners"
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
	"github.com/planetdecred/godcr/wallet"
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
	*listeners.SyncProgressListener
	*listeners.TxAndBlockNotificationListener
	*listeners.ProposalNotificationListener
	ctx       context.Context
	ctxCancel context.CancelFunc
	appBarNav components.NavDrawer
	drawerNav components.NavDrawer

	hideBalanceItem HideBalanceItem

	currentPage   load.Page
	pageBackStack []load.Page
	sendPage      *send.Page   // reuse value to keep data persistent onresume.
	receivePage   *ReceivePage // pointer to receive page. to avoid duplication.

	refreshExchangeRateBtn *decredmaterial.Clickable

	// page state variables
	usdExchangeSet         bool
	isFetchingExchangeRate bool
	dcrUsdtBittrex         load.DCRUSDTBittrex
	isBalanceHidden        bool
	totalBalance           dcrutil.Amount
	totalBalanceUSD        string
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

	mp.refreshExchangeRateBtn = mp.Theme.NewClickable(true)

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
	mp.setLanguageSetting()

	mp.ctx, mp.ctxCancel = context.WithCancel(context.TODO())
	mp.listenForNotifications()

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

	mp.UpdateBalance()
}

func (mp *MainPage) setLanguageSetting() {
	langPre := mp.WL.Wallet.ReadStringConfigValueForKey(load.LanguagePreferenceKey)
	values.SetUserLanguage(langPre)
}

func (mp *MainPage) updateExchangeSetting() {
	currencyExchangeValue := mp.WL.Wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	usdExchangeSet := currencyExchangeValue == values.USDExchangeValue
	if mp.usdExchangeSet == usdExchangeSet {
		return // nothing has changed
	}
	mp.usdExchangeSet = usdExchangeSet
	if mp.usdExchangeSet {
		go mp.fetchExchangeRate()
	}
}

func (mp *MainPage) fetchExchangeRate() {
	if mp.isFetchingExchangeRate {
		return
	}
	maxAttempts := 5
	delayBtwAttempts := 2 * time.Second
	mp.isFetchingExchangeRate = true
	desc := "for getting dcrUsdtBittrex exchange rate value"
	attempts, err := components.RetryFunc(maxAttempts, delayBtwAttempts, desc, func() error {
		return load.GetUSDExchangeValue(&mp.dcrUsdtBittrex)
	})
	if err != nil {
		log.Errorf("error fetching usd exchange rate value after %d attempts: %v", attempts, err)
	} else if mp.dcrUsdtBittrex.LastTradeRate == "" {
		log.Errorf("no error while fetching usd exchange rate in %d tries, but no rate was fetched", attempts)
	} else {
		log.Infof("exchange rate value fetched: %s", mp.dcrUsdtBittrex.LastTradeRate)
		mp.UpdateBalance()
		mp.RefreshWindow()
	}
	mp.isFetchingExchangeRate = false
}

func (mp *MainPage) UpdateBalance() {
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

	if mp.refreshExchangeRateBtn.Clicked() {
		go mp.fetchExchangeRate()
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

	mp.ctxCancel()
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
	mp.updateExchangeSetting() // the setting may have changed, leading to this window refresh, let's check
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
	if !mp.usdExchangeSet {
		return D{}
	}
	switch {
	case mp.isFetchingExchangeRate && mp.dcrUsdtBittrex.LastTradeRate == "":
		gtx.Constraints.Max.Y = gtx.Px(values.MarginPadding18)
		gtx.Constraints.Max.X = gtx.Constraints.Max.Y
		return layout.Inset{
			Top:  values.MarginPadding8,
			Left: values.MarginPadding5,
		}.Layout(gtx, func(gtx C) D {
			loader := material.Loader(mp.Theme.Base)
			return loader.Layout(gtx)
		})
	case !mp.isFetchingExchangeRate && mp.dcrUsdtBittrex.LastTradeRate == "":
		return layout.Inset{
			Top:  values.MarginPadding7,
			Left: values.MarginPadding5,
		}.Layout(gtx, func(gtx C) D {
			return mp.refreshExchangeRateBtn.Layout(gtx, func(gtx C) D {
				return mp.Icons.Restore.Layout16dp(gtx)
			})
		})
	case len(mp.totalBalanceUSD) > 0:
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
	default:
		return D{}
	}
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

// postDdesktopNotification posts notifications to the desktop.
func (mp *MainPage) postDesktopNotification(notifier interface{}) {
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

// listenForNotifications starts a goroutine to watch for sync updates
// and update the UI accordingly.
func (mp *MainPage) listenForNotifications() {
	// Return if any of the listener is not nill.
	switch {
	case  mp.SyncProgressListener != nil:
		return
	case mp.TxAndBlockNotificationListener != nil:
		return
	case mp.ProposalNotificationListener != nil:
		return
	}
	
	mp.SyncProgressListener = listeners.NewSyncProgress()
	err := mp.WL.MultiWallet.AddSyncProgressListener(mp.SyncProgressListener, MainPageID)
	if err != nil {
		log.Errorf("Error adding sync progress listener: %v", err)
		return
	}

	
	mp.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err = mp.WL.MultiWallet.AddTxAndBlockNotificationListener(mp.TxAndBlockNotificationListener, true, MainPageID)
	if err != nil {
		log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	
	mp.ProposalNotificationListener = listeners.NewProposalNotificationListener()
	err = mp.WL.MultiWallet.Politeia.AddNotificationListener(mp.ProposalNotificationListener, MainPageID)
	if err != nil {
		log.Errorf("Error adding politeia notification listener: %v", err)
		return
	}

	go func() {
		for {
			select {
			case n := <-mp.TxAndBlockNotifChan:
				switch n.Type {
				case listeners.NewTransaction:
					mp.UpdateBalance(false)
					transactionNotification := mp.WL.Wallet.ReadBoolConfigValueForKey(load.TransactionNotificationConfigKey)
					update := wallet.NewTransaction{
						Transaction: n.Transaction,
					}
					if transactionNotification {
						mp.postDesktopNotification(update)
					}
					mp.RefreshWindow()
				case listeners.BlockAttached:
					beep := mp.WL.Wallet.ReadBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey)
					if beep {
						_ = beeep.Beep(5, 1)
					}

					mp.UpdateBalance(false)
					mp.RefreshWindow()
				case listeners.TxConfirmed:
					mp.UpdateBalance(false)
					mp.RefreshWindow()

				}
			case notification := <-mp.ProposalNotifChan:
				// Don't notify on wallet synced event.
				if notification.ProposalStatus != wallet.Synced {
					mp.postDesktopNotification(notification)
				}
			case n := <-mp.SyncStatusChan:
				if n.Stage == wallet.SyncCompleted {
					mp.UpdateBalance(false)
					mp.RefreshWindow()
				}
			case <-mp.ctx.Done():
				mp.WL.MultiWallet.RemoveSyncProgressListener(MainPageID)
				mp.WL.MultiWallet.RemoveTxAndBlockNotificationListener(MainPageID)
				mp.WL.MultiWallet.Politeia.RemoveNotificationListener(MainPageID)

				close(mp.SyncStatusChan)
				close(mp.TxAndBlockNotifChan)
				close(mp.ProposalNotifChan)

				mp.SyncProgressListener == nil
				mp.TxAndBlockNotificationListener == nil
				mp.ProposalNotificationListener == nil

				return
			}
		}
	}()
}
