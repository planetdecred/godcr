package page

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/gen2brain/beeep"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/dexclient"
	"github.com/planetdecred/godcr/ui/page/governance"
	"github.com/planetdecred/godcr/ui/page/overview"
	"github.com/planetdecred/godcr/ui/page/privacy"
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

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	NavDrawerWidth          = unit.Dp(160)
	NavDrawerMinimizedWidth = unit.Dp(72)
)

type HideBalanceItem struct {
	hideBalanceButton decredmaterial.IconButton
	tooltip           *decredmaterial.Tooltip
}

type NavHandler struct {
	Clickable     *widget.Clickable
	Image         *decredmaterial.Image
	ImageInactive *decredmaterial.Image
	Title         string
	PageID        string
}

type MainPage struct {
	*app.MasterPage

	*load.Load
	*listeners.SyncProgressListener
	*listeners.TxAndBlockNotificationListener
	*listeners.ProposalNotificationListener
	ctx                  context.Context
	ctxCancel            context.CancelFunc
	drawerNav            components.NavDrawer
	bottomNavigationBar  components.BottomNavigationBar
	floatingActionButton components.BottomNavigationBar

	hideBalanceItem HideBalanceItem

	sendPage    *send.Page   // reuse value to keep data persistent onresume.
	receivePage *ReceivePage // pointer to receive page. to avoid duplication.

	refreshExchangeRateBtn *decredmaterial.Clickable
	darkmode               *decredmaterial.Clickable
	openWalletSelector     *decredmaterial.Clickable

	// page state variables
	dcrUsdtBittrex load.DCRUSDTBittrex
	totalBalance   dcrutil.Amount

	usdExchangeSet         bool
	isFetchingExchangeRate bool
	isBalanceHidden        bool
	isNavExpanded          bool
	setNavExpanded         func()
	totalBalanceUSD        string
}

func NewMainPage(l *load.Load) *MainPage {
	mp := &MainPage{
		Load:       l,
		MasterPage: app.NewMasterPage(MainPageID),
	}

	mp.hideBalanceItem.hideBalanceButton = mp.Theme.IconButton(mp.Theme.Icons.ConcealIcon)
	mp.hideBalanceItem.hideBalanceButton.Size = unit.Dp(19)
	mp.hideBalanceItem.hideBalanceButton.Inset = layout.UniformInset(values.MarginPadding4)
	mp.hideBalanceItem.tooltip = mp.Theme.Tooltip()

	mp.darkmode = mp.Theme.NewClickable(false)
	mp.openWalletSelector = mp.Theme.NewClickable(true)
	mp.openWalletSelector.Radius = decredmaterial.Radius(8)

	// init shared page functions
	toggleSync := func() {
		if mp.WL.MultiWallet.IsConnectedToDecredNetwork() {
			mp.WL.MultiWallet.CancelSync()
		} else {
			mp.StartSyncing()
		}
	}
	l.ToggleSync = toggleSync

	mp.setLanguageSetting()

	mp.initNavItems()

	mp.refreshExchangeRateBtn = mp.Theme.NewClickable(true)

	// Show a seed backup prompt if necessary.
	mp.WL.MultiWallet.SaveUserConfigValue(load.SeedBackupNotificationConfigKey, false)

	mp.setNavExpanded = func() {
		mp.drawerNav.DrawerToggled(mp.isNavExpanded)
	}

	mp.bottomNavigationBar.OnViewCreated()

	return mp
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (mp *MainPage) ID() string {
	return MainPageID
}

func (mp *MainPage) initNavItems() {
	mp.drawerNav = components.NavDrawer{
		Load:        mp.Load,
		CurrentPage: mp.CurrentPageID(),
		DrawerNavItems: []components.NavHandler{
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.OverviewIcon,
				ImageInactive: mp.Theme.Icons.OverviewIconInactive,
				Title:         values.String(values.StrOverview),
				PageID:        overview.OverviewPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.SendIcon,
				Title:         values.String(values.StrSend),
				ImageInactive: mp.Theme.Icons.SendInactiveIcon,
				PageID:        send.SendPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.ReceiveIcon,
				ImageInactive: mp.Theme.Icons.ReceiveInactiveIcon,
				Title:         values.String(values.StrReceive),
				PageID:        ReceivePageID,
			},
			{
				// TODO -- deprectated in v2 layout
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.WalletIcon,
				ImageInactive: mp.Theme.Icons.WalletIconInactive,
				Title:         values.String(values.StrWallets),
				PageID:        wallets.WalletPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.TransactionsIcon,
				ImageInactive: mp.Theme.Icons.TransactionsIconInactive,
				Title:         values.String(values.StrTransactions),
				PageID:        transaction.TransactionsPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.Mixer,
				ImageInactive: mp.Theme.Icons.MixerInactive,
				Title:         values.String(values.StrStakeShuffle),
				PageID:        privacy.AccountMixerPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.StakeIcon,
				ImageInactive: mp.Theme.Icons.StakeIconInactive,
				Title:         values.String(values.StrStaking),
				PageID:        staking.OverviewPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.GovernanceActiveIcon,
				ImageInactive: mp.Theme.Icons.GovernanceInactiveIcon,
				Title:         values.String(values.StrGovernance),
				PageID:        governance.GovernancePageID,
			},
			// Temp disabling. Will uncomment after release
			// {
			// 	Clickable:     mp.Theme.NewClickable(true),
			// 	Image:         mp.Theme.Icons.DexIcon,
			// 	ImageInactive: mp.Theme.Icons.DexIconInactive,
			// 	Title:         values.String(values.StrDex),
			// 	PageID:        dexclient.MarketPageID,
			// },
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.MoreIcon,
				ImageInactive: mp.Theme.Icons.MoreIconInactive,
				Title:         values.String(values.StrMore),
				PageID:        MorePageID,
			},
		},
		MinimizeNavDrawerButton: mp.Theme.IconButton(mp.Theme.Icons.NavigationArrowBack),
		MaximizeNavDrawerButton: mp.Theme.IconButton(mp.Theme.Icons.NavigationArrowForward),
	}

	mp.bottomNavigationBar = components.BottomNavigationBar{
		Load:        mp.Load,
		CurrentPage: mp.CurrentPageID(),
		BottomNaigationItems: []components.BottomNavigationBarHandler{
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.OverviewIcon,
				ImageInactive: mp.Theme.Icons.OverviewIconInactive,
				Title:         values.String(values.StrOverview),
				PageID:        overview.OverviewPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.TransactionsIcon,
				ImageInactive: mp.Theme.Icons.TransactionsIconInactive,
				Title:         values.String(values.StrTransactions),
				PageID:        transaction.TransactionsPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.StakeIcon,
				ImageInactive: mp.Theme.Icons.StakeIconInactive,
				Title:         values.String(values.StrStaking),
				PageID:        staking.OverviewPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.WalletIcon,
				ImageInactive: mp.Theme.Icons.WalletIconInactive,
				Title:         values.String(values.StrWallets),
				PageID:        wallets.WalletPageID,
			},
			{
				Clickable:     mp.Theme.NewClickable(true),
				Image:         mp.Theme.Icons.MoreIcon,
				ImageInactive: mp.Theme.Icons.MoreIconInactive,
				Title:         values.String(values.StrMore),
				PageID:        MorePageID,
			},
		},
	}

	mp.floatingActionButton = components.BottomNavigationBar{
		Load:        mp.Load,
		CurrentPage: mp.CurrentPageID(),
		FloatingActionButton: []components.BottomNavigationBarHandler{
			{
				Clickable: mp.Theme.NewClickable(true),
				Image:     mp.Theme.Icons.SendIcon,
				Title:     values.String(values.StrSend),
				PageID:    send.SendPageID,
			},
			{
				Clickable: mp.Theme.NewClickable(true),
				Image:     mp.Theme.Icons.ReceiveIcon,
				Title:     values.String(values.StrReceive),
				PageID:    ReceivePageID,
			},
		},
	}
	mp.floatingActionButton.FloatingActionButton[0].Clickable.Hoverable = false
	mp.floatingActionButton.FloatingActionButton[1].Clickable.Hoverable = false
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (mp *MainPage) OnNavigatedTo() {
	mp.setNavExpanded()

	mp.ctx, mp.ctxCancel = context.WithCancel(context.TODO())
	mp.listenForNotifications()

	if mp.CurrentPage() == nil {
		mp.Display(overview.NewOverviewPage(mp.Load)) // TODO: Should pagestack have a start page?
	}
	mp.CurrentPage().OnNavigatedTo()

	if mp.sendPage != nil {
		mp.sendPage.OnNavigatedTo()
	}
	if mp.receivePage != nil {
		mp.receivePage.OnNavigatedTo()
	}

	if mp.WL.MultiWallet.ReadBoolConfigValueForKey(load.AutoSyncConfigKey, false) {
		mp.StartSyncing()
		if mp.WL.MultiWallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey, false) {
			go mp.WL.MultiWallet.Politeia.Sync()
		}
	}

	mp.updateBalance()

}

func (mp *MainPage) setLanguageSetting() {
	langPre := mp.WL.MultiWallet.ReadStringConfigValueForKey(load.LanguagePreferenceKey)
	if langPre == "" {
		mp.WL.MultiWallet.SaveUserConfigValue(load.LanguagePreferenceKey, values.DefaultLangauge)
	}
	values.SetUserLanguage(langPre)
}

func (mp *MainPage) updateExchangeSetting() {
	currencyExchangeValue := mp.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if currencyExchangeValue == "" {
		mp.WL.MultiWallet.SaveUserConfigValue(dcrlibwallet.CurrencyConversionConfigKey, values.DefaultExchangeValue)
	}

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
		mp.updateBalance()
		mp.ParentWindow().Reload()
	}
	mp.isFetchingExchangeRate = false
}

func (mp *MainPage) updateBalance() {
	totalBalance, err := components.CalculateTotalWalletsBalance(mp.Load)
	if err == nil {
		mp.totalBalance = totalBalance.Total

		if mp.usdExchangeSet && mp.dcrUsdtBittrex.LastTradeRate != "" {
			usdExchangeRate, err := strconv.ParseFloat(mp.dcrUsdtBittrex.LastTradeRate, 64)
			if err == nil {
				balanceInUSD := totalBalance.Total.ToCoin() * usdExchangeRate
				mp.totalBalanceUSD = load.FormatUSDBalance(mp.Printer, balanceInUSD)
			}
		}
	}
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
	spendingPasswordModal := modal.NewPasswordModal(mp.Load).
		Title(values.String(values.StrResumeAccountDiscoveryTitle)).
		Hint(values.String(values.StrSpendingPassword)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(values.String(values.StrUnlock), func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := mp.WL.MultiWallet.UnlockWallet(wal.ID, []byte(password))
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
				mp.StartSyncing()
			}()

			return false
		})
	mp.ParentWindow().ShowModal(spendingPasswordModal)
}

// OnDarkModeChanged is triggered whenever the dark mode setting is changed
// to enable restyling UI elements where necessary.
// Satisfies the load.AppSettingsChangeHandler interface.
func (mp *MainPage) OnDarkModeChanged(isDarkModeOn bool) {
	// TODO: currentPage will likely be the Settings page when this method
	// is called. If that page implements the AppSettingsChangeHandler interface,
	// the following code will trigger the OnDarkModeChanged method of that
	// page.
	if currentPage, ok := mp.CurrentPage().(load.AppSettingsChangeHandler); ok {
		currentPage.OnDarkModeChanged(isDarkModeOn)
	}

	mp.initNavItems()
	mp.setNavExpanded()
}

func (mp *MainPage) OnLanguageChanged() {
	mp.setLanguageSetting()
}

func (mp *MainPage) OnCurrencyChanged() {
	mp.updateExchangeSetting()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (mp *MainPage) HandleUserInteractions() {
	if mp.CurrentPage() != nil {
		mp.CurrentPage().HandleUserInteractions()
	}

	if mp.refreshExchangeRateBtn.Clicked() {
		go mp.fetchExchangeRate()
	}

	// darkmode settings
	for mp.darkmode.Clicked() {
		isDarkModeOn := mp.WL.MultiWallet.ReadBoolConfigValueForKey(load.DarkModeConfigKey, false)
		if isDarkModeOn {
			mp.WL.MultiWallet.SaveUserConfigValue(load.DarkModeConfigKey, false)
		} else {
			mp.WL.MultiWallet.SaveUserConfigValue(load.DarkModeConfigKey, true)
		}

		mp.RefreshTheme(mp.ParentWindow())
	}

	for mp.openWalletSelector.Clicked() {
		onWalSelected := func() {
			mp.ParentNavigator().CloseCurrentPage()
		}
		mp.ParentWindow().Display(NewWalletList(mp.Load, onWalSelected))
	}

	mp.drawerNav.CurrentPage = mp.CurrentPageID()
	mp.bottomNavigationBar.CurrentPage = mp.CurrentPageID()
	mp.floatingActionButton.CurrentPage = mp.CurrentPageID()

	for mp.drawerNav.MinimizeNavDrawerButton.Button.Clicked() {
		mp.isNavExpanded = true
		mp.setNavExpanded()
	}

	for mp.drawerNav.MaximizeNavDrawerButton.Button.Clicked() {
		mp.isNavExpanded = false
		mp.setNavExpanded()
	}

	for _, item := range mp.drawerNav.DrawerNavItems {
		for item.Clickable.Clicked() {
			var pg app.Page
			switch item.PageID {
			case overview.OverviewPageID:
				pg = overview.NewOverviewPage(mp.Load) // todo :New wallet ui --- current overview page is deprecated.
			case send.SendPageID:
				pg = send.NewSendPage(mp.Load)
			case ReceivePageID:
				pg = NewReceivePage(mp.Load)
			case wallets.WalletPageID:
				pg = wallets.NewWalletPage(mp.Load)
			case transaction.TransactionsPageID:
				pg = transaction.NewTransactionsPage(mp.Load)
			case privacy.AccountMixerPageID:
				pg = privacy.NewAccountMixerPage(mp.Load, mp.WL.SelectedWallet.Wallet) // todo implement new staking ui
			case staking.OverviewPageID:
				pg = staking.NewStakingPage(mp.Load)
			case governance.GovernancePageID:
				pg = governance.NewGovernancePage(mp.Load)
			case dexclient.MarketPageID:
				_, err := mp.WL.MultiWallet.StartDexClient() // does nothing if already started
				if err != nil {
					mp.Toast.NotifyError(values.StringF(values.StrDexStartupErr, err))
				} else {
					pg = dexclient.NewMarketPage(mp.Load)
				}
			case MorePageID:
				pg = NewMorePage(mp.Load)
			}

			if pg == nil || mp.ID() == mp.CurrentPageID() {
				continue
			}

			// check if wallet is synced and clear stack
			if mp.ID() == send.SendPageID || mp.ID() == ReceivePageID {
				if mp.WL.MultiWallet.IsSynced() {
					mp.Display(pg)
				} else if mp.WL.MultiWallet.IsSyncing() {
					mp.Toast.NotifyError(values.String(values.StrNotConnected))
				} else {
					mp.Toast.NotifyError(values.String(values.StrWalletSyncing))
				}
			} else {
				mp.Display(pg)
			}
		}
	}

	for _, item := range mp.bottomNavigationBar.BottomNaigationItems {
		for item.Clickable.Clicked() {
			var pg app.Page
			switch item.PageID {
			case overview.OverviewPageID:
				pg = overview.NewOverviewPage(mp.Load)
			case transaction.TransactionsPageID:
				pg = transaction.NewTransactionsPage(mp.Load)
			case staking.OverviewPageID:
				pg = staking.NewStakingPage(mp.Load)
			case wallets.WalletPageID:
				pg = wallets.NewWalletPage(mp.Load)
			case MorePageID:
				pg = NewMorePage(mp.Load)
			}

			if pg == nil || mp.ID() == mp.CurrentPageID() {
				continue
			}

			// clear stack
			mp.Display(pg)
		}
	}

	for i, item := range mp.floatingActionButton.FloatingActionButton {
		for item.Clickable.Clicked() {
			var pg app.Page
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

			if mp.ID() == mp.CurrentPageID() {
				continue
			}

			if mp.WL.MultiWallet.IsSynced() {
				mp.Display(pg)
			} else if mp.WL.MultiWallet.IsSyncing() {
				mp.Toast.NotifyError(values.String(values.StrWalletSyncing))
			} else {
				mp.Toast.NotifyError(values.String(values.StrNotConnected))
			}
		}
	}

	mp.isBalanceHidden = mp.WL.MultiWallet.ReadBoolConfigValueForKey(load.HideBalanceConfigKey, false)
	for mp.hideBalanceItem.hideBalanceButton.Button.Clicked() {
		mp.isBalanceHidden = !mp.isBalanceHidden
		mp.WL.MultiWallet.SetBoolConfigValueForKey(load.HideBalanceConfigKey, mp.isBalanceHidden)
	}
}

// HandleKeyEvent is called when a key is pressed on the current window.
// Satisfies the load.KeyEventHandler interface for receiving key events.
func (mp *MainPage) HandleKeyEvent(evt *key.Event) {
	if currentPage := mp.CurrentPage(); currentPage != nil {
		if keyEvtHandler, ok := currentPage.(load.KeyEventHandler); ok {
			keyEvtHandler.HandleKeyEvent(evt)
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
func (mp *MainPage) OnNavigatedFrom() {
	// Also disappear all child pages.
	if mp.CurrentPage() != nil {
		mp.CurrentPage().OnNavigatedFrom()
	}
	if mp.sendPage != nil {
		mp.sendPage.OnNavigatedFrom()
	}
	if mp.receivePage != nil {
		mp.receivePage.OnNavigatedFrom()
	}

	mp.ctxCancel()
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (mp *MainPage) Layout(gtx layout.Context) layout.Dimensions {
	mp.Load.SetCurrentAppWidth(gtx.Constraints.Max.X)
	if mp.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return mp.layoutMobile(gtx)
	}
	return mp.layoutDesktop(gtx)
}

func (mp *MainPage) layoutDesktop(gtx layout.Context) layout.Dimensions {
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
							if mp.CurrentPage() == nil {
								return D{}
							}
							return mp.CurrentPage().Layout(gtx)
						}),
					)
				}),
			)
		}),
	)
}

func (mp *MainPage) layoutMobile(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(0.08, mp.LayoutTopBar),
		layout.Flexed(0.795, func(gtx C) D {
			return layout.Stack{Alignment: layout.N}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					currentPage := mp.CurrentPage()
					if currentPage == nil {
						return layout.Dimensions{}
					}
					return currentPage.Layout(gtx)
				}),
				layout.Stacked(func(gtx C) D {
					return mp.floatingActionButton.LayoutSendReceive(gtx)
				}),
			)
		}),
		layout.Flexed(0.125, mp.bottomNavigationBar.LayoutBottomNavigationBar),
	)
}

func (mp *MainPage) LayoutUSDBalance(gtx layout.Context) layout.Dimensions {
	if !mp.usdExchangeSet {
		return D{}
	}
	switch {
	case mp.isFetchingExchangeRate && mp.dcrUsdtBittrex.LastTradeRate == "":
		gtx.Constraints.Max.Y = gtx.Dp(values.MarginPadding18)
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
				return mp.Theme.Icons.Restore.Layout16dp(gtx)
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

func (mp *MainPage) totalDCRBalance(gtx C) D {
	if mp.isBalanceHidden {
		hiddenBalanceText := mp.Theme.Label(values.TextSize18*0.8, "**********DCR")
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
			h := values.MarginPadding24
			v := values.MarginPadding8
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.WrapContent,
				Orientation: layout.Horizontal,
				Alignment:   layout.Middle,
				Padding: layout.Inset{
					Right:  h,
					Left:   values.MarginPadding10,
					Top:    v,
					Bottom: v,
				},
			}.GradientLayout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return decredmaterial.LinearLayout{
							Width:       decredmaterial.WrapContent,
							Height:      decredmaterial.WrapContent,
							Orientation: layout.Horizontal,
							Alignment:   layout.Middle,
						}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return decredmaterial.LinearLayout{
									Width:  gtx.Dp(values.Size180),
									Height: decredmaterial.WrapContent,
									Padding: layout.Inset{
										Right:  values.MarginPadding18,
										Left:   values.MarginPadding18,
										Top:    values.MarginPadding10,
										Bottom: values.MarginPadding10,
									},
									Alignment: layout.Middle,
									Clickable: mp.openWalletSelector,
									Shadow:    mp.Theme.Shadow(),
									Border: decredmaterial.Border{
										Radius: mp.openWalletSelector.Radius,
										Width:  values.MarginPadding2,
										Color:  mp.Theme.Color.Gray3,
									},
								}.GradientLayout(gtx,
									layout.Rigid(mp.Theme.Icons.WalletIcon.Layout24dp),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left: values.MarginPadding10,
										}.Layout(gtx, func(gtx C) D {
											txt := mp.Theme.Body1(mp.WL.SelectedWallet.Wallet.Name)
											txt.Font.Weight = text.Bold
											return txt.Layout(gtx)
										})
									}),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Right: values.MarginPadding16,
									Left:  values.MarginPadding24,
								}.Layout(gtx,
									func(gtx C) D {
										return mp.Theme.Icons.Logo.Layout24dp(gtx)
									})
							}),
							layout.Rigid(func(gtx C) D {
								return mp.totalDCRBalance(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								if !mp.isBalanceHidden {
									return mp.LayoutUSDBalance(gtx)
								}
								return D{}
							}),
							layout.Rigid(func(gtx C) D {
								mp.hideBalanceItem.hideBalanceButton.Icon = mp.Theme.Icons.RevealIcon
								if mp.isBalanceHidden {
									mp.hideBalanceItem.hideBalanceButton.Icon = mp.Theme.Icons.ConcealIcon
								}
								return layout.Inset{
									Top:  values.MarginPadding1,
									Left: values.MarginPadding9,
								}.Layout(gtx, mp.hideBalanceItem.hideBalanceButton.Layout)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.E.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								//todo -- dex functionality
								return mp.Theme.Icons.DexIcon.Layout24dp(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Right: values.MarginPadding24,
									Left:  values.MarginPadding24,
								}.Layout(gtx, func(gtx C) D {
									//todo -- app level settings functionality
									return mp.Theme.Icons.HeaderSettingsIcon.Layout24dp(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return mp.darkmode.Layout(gtx, mp.Theme.Icons.DarkmodeIcon.Layout24dp)
							}),
						)
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return mp.Theme.Separator().Layout(gtx)
		}),
	)
}

// postDdesktopNotification posts notifications to the desktop.
func (mp *MainPage) postDesktopNotification(notifier interface{}) {
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
			notification = values.StringF(values.StrDcrReceived, amount)
		case dcrlibwallet.TxTypeVote:
			reward := strconv.FormatFloat(dcrlibwallet.AmountCoin(t.Transaction.VoteReward), 'f', -1, 64)
			notification = values.StringF(values.StrTicektVoted, reward)
		case dcrlibwallet.TxTypeRevocation:
			notification = values.String(values.StrTicketRevoked)
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
		proposalNotification := mp.WL.MultiWallet.ReadBoolConfigValueForKey(load.ProposalNotificationConfigKey, false)
		if !proposalNotification {
			return
		}
		switch {
		case t.ProposalStatus == wallet.NewProposalFound:
			notification = values.StringF(values.StrProposalAddedNotif, t.Proposal.Name)
		case t.ProposalStatus == wallet.VoteStarted:
			notification = values.StringF(values.StrVoteStartedNotif, t.Proposal.Name)
		case t.ProposalStatus == wallet.VoteFinished:
			notification = values.StringF(values.StrVoteEndedNotif, t.Proposal.Name)
		default:
			notification = values.StringF(values.StrNewProposalUpdate, t.Proposal.Name)
		}
		initializeBeepNotification(notification)
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

// listenForNotifications starts a goroutine to watch for notifications
// and update the UI accordingly.
func (mp *MainPage) listenForNotifications() {
	// Return if any of the listener is not nil.
	switch {
	case mp.SyncProgressListener != nil:
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
					mp.updateBalance()
					transactionNotification := mp.WL.MultiWallet.ReadBoolConfigValueForKey(load.TransactionNotificationConfigKey, false)
					if transactionNotification {
						update := wallet.NewTransaction{
							Transaction: n.Transaction,
						}
						mp.postDesktopNotification(update)
					}
					mp.ParentWindow().Reload()
				case listeners.BlockAttached:
					beep := mp.WL.MultiWallet.ReadBoolConfigValueForKey(dcrlibwallet.BeepNewBlocksConfigKey, false)
					if beep {
						err := beeep.Beep(5, 1)
						if err != nil {
							log.Error(err.Error)
						}
					}

					mp.updateBalance()
					mp.ParentWindow().Reload()
				case listeners.TxConfirmed:
					mp.updateBalance()
					mp.ParentWindow().Reload()

				}
			case notification := <-mp.ProposalNotifChan:
				// Post desktop notification for all events except the synced event.
				if notification.ProposalStatus != wallet.Synced {
					mp.postDesktopNotification(notification)
				}
			case n := <-mp.SyncStatusChan:
				if n.Stage == wallet.SyncCompleted {
					mp.updateBalance()
					mp.ParentWindow().Reload()
				}
			case <-mp.ctx.Done():
				mp.WL.MultiWallet.RemoveSyncProgressListener(MainPageID)
				mp.WL.MultiWallet.RemoveTxAndBlockNotificationListener(MainPageID)
				mp.WL.MultiWallet.Politeia.RemoveNotificationListener(MainPageID)

				close(mp.SyncStatusChan)
				close(mp.TxAndBlockNotifChan)
				close(mp.ProposalNotifChan)

				mp.SyncProgressListener = nil
				mp.TxAndBlockNotificationListener = nil
				mp.ProposalNotificationListener = nil

				return
			}
		}
	}()
}
