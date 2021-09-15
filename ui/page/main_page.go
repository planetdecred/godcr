package page

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/proposal"
	"github.com/planetdecred/godcr/ui/page/send"
	"github.com/planetdecred/godcr/ui/page/tickets"
	"github.com/planetdecred/godcr/ui/values"
)

const MainPageID = "Main"

const (
	OverviewNavID = iota
	TransactionsNavID
	WalletsNavID
	TicketsNavID
	ProposalsNavID
	MoreNavID
)

var (
	NavDrawerWidth          = unit.Value{U: unit.UnitDp, V: 160}
	NavDrawerMinimizedWidth = unit.Value{U: unit.UnitDp, V: 72}
)

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

	autoSync bool

	currentPage   load.Page
	pageBackStack []load.Page
	sendPage      *send.Page // reuse value to keep data persistent onresume.

	// page state variables
	dcrUsdtBittrex  load.DCRUSDTBittrex
	usdExchangeSet  bool
	totalBalance    dcrutil.Amount
	totalBalanceUSD string
}

func NewMainPage(l *load.Load) *MainPage {

	mp := &MainPage{
		Load:     l,
		autoSync: true,
	}

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

	mp.drawerNav.DrawerToggled(false)

	mp.OnResume()

	return mp
}

func (mp *MainPage) ID() string {
	return MainPageID
}

func (mp *MainPage) initNavItems() {
	mp.appBarNav = components.NavDrawer{
		Load:        mp.Load,
		CurrentPage: mp.currentPageID(),
		AppBarNavItems: []components.NavHandler{
			{
				Clickable: new(widget.Clickable),
				Image:     mp.Icons.SendIcon,
				Title:     values.String(values.StrSend),
				PageID:    send.PageID,
			},
			{
				Clickable: new(widget.Clickable),
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
				Clickable:     new(widget.Clickable),
				Image:         mp.Icons.OverviewIcon,
				ImageInactive: mp.Icons.OverviewIconInactive,
				Title:         values.String(values.StrOverview),
				PageID:        OverviewPageID,
			},
			{
				Clickable:     new(widget.Clickable),
				Image:         mp.Icons.TransactionsIcon,
				ImageInactive: mp.Icons.TransactionsIconInactive,
				Title:         values.String(values.StrTransactions),
				PageID:        TransactionsPageID,
			},
			{
				Clickable:     new(widget.Clickable),
				Image:         mp.Icons.WalletIcon,
				ImageInactive: mp.Icons.WalletIconInactive,
				Title:         values.String(values.StrWallets),
				PageID:        WalletPageID,
			},
			{
				Clickable:     new(widget.Clickable),
				Image:         mp.Icons.TicketIcon,
				ImageInactive: mp.Icons.TicketIconInactive,
				Title:         values.String(values.StrTickets),
				PageID:        tickets.OverviewPageID,
			},
			{
				Clickable:     new(widget.Clickable),
				Image:         mp.Icons.ProposalIconActive,
				ImageInactive: mp.Icons.ProposalIconInactive,
				Title:         values.String(values.StrProposal),
				PageID:        proposal.ProposalsPageID,
			},
			{
				Clickable:     new(widget.Clickable),
				Image:         mp.Icons.MoreIcon,
				ImageInactive: mp.Icons.MoreIconInactive,
				Title:         values.String(values.StrMore),
				PageID:        MorePageID,
			},
		},
		MinimizeNavDrawerButton: mp.Theme.PlainIconButton(new(widget.Clickable), mp.Icons.NavigationArrowBack),
		MaximizeNavDrawerButton: mp.Theme.PlainIconButton(new(widget.Clickable), mp.Icons.NavigationArrowForward),
	}
}

func (mp *MainPage) OnResume() {
	// register for notifications
	mp.WL.MultiWallet.AddAccountMixerNotificationListener(mp, MainPageID)
	mp.WL.MultiWallet.Politeia.AddNotificationListener(mp, MainPageID)
	mp.WL.MultiWallet.AddTxAndBlockNotificationListener(mp, MainPageID)
	mp.WL.MultiWallet.AddSyncProgressListener(mp, MainPageID)

	mp.getSetting()
	mp.UpdateBalance()

	if mp.currentPage != nil {
		mp.currentPage.OnResume()
	} else {
		mp.ChangeFragment(NewOverviewPage(mp.Load))
	}

	if mp.autoSync {
		mp.autoSync = false
		mp.StartSyncing()
		go mp.WL.MultiWallet.Politeia.Sync()
	}
}

func (mp *MainPage) getSetting() {
	langPre := mp.WL.Wallet.ReadStringConfigValueForKey(languagePreferenceKey)
	if components.StringNotEmpty(langPre) {
		langPre = values.DefaultLangauge
		mp.WL.Wallet.SaveConfigValueForKey(languagePreferenceKey, langPre)
	}
	values.SetUserLanguage(langPre)

	currencyPre := mp.WL.Wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if components.StringNotEmpty(currencyPre) {
		mp.WL.Wallet.SaveConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey, DefaultExchangeValue)
	}
}

func (mp *MainPage) UpdateBalance() {
	currencyExchangeValue := mp.WL.Wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	mp.usdExchangeSet = currencyExchangeValue == components.USDExchangeValue

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
	for _, wallet := range mp.WL.SortedWalletList() {
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
		Hint(wal.Name+" Spending password").
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(values.String(values.StrUnlock), func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := mp.WL.MultiWallet.UnlockWallet(wal.ID, []byte(password))
				if err != nil {
					errText := err.Error()
					if err.Error() == "invalid_passphrase" {
						errText = "Invalid passphrase"
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

func (mp *MainPage) Handle() {
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
				pg = NewReceivePage(mp.Load)
			}

			mp.ChangeFragment(pg)
		}
	}

	for i, item := range mp.drawerNav.DrawerNavItems {
		for item.Clickable.Clicked() {
			var pg load.Page
			if i == OverviewNavID {
				pg = NewOverviewPage(mp.Load)
			} else if i == TransactionsNavID {
				pg = NewTransactionsPage(mp.Load)
			} else if i == WalletsNavID {
				pg = NewWalletPage(mp.Load)
			} else if i == TicketsNavID {
				pg = tickets.NewTicketPage(mp.Load)
			} else if i == ProposalsNavID {
				pg = proposal.NewProposalsPage(mp.Load)
			} else if i == MoreNavID {
				pg = NewMorePage(mp.Load)
			} else {
				continue
			}

			if pg.ID() == mp.currentPageID() {
				continue
			}

			// clear stack
			mp.changeFragment(pg)
		}
	}
}

func (mp *MainPage) OnClose() {
	if mp.currentPage != nil {
		mp.currentPage.OnClose()
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
	if mp.currentPage != nil {
		mp.currentPage.OnClose()
		mp.pageBackStack = append(mp.pageBackStack, mp.currentPage)
	}

	page.OnResume()
	mp.currentPage = page
}

// popFragment goes back to the previous page
func (mp *MainPage) popFragment() {
	if len(mp.pageBackStack) > 0 {
		// get and remove last page
		previousPage := mp.pageBackStack[len(mp.pageBackStack)-1]
		mp.pageBackStack = mp.pageBackStack[:len(mp.pageBackStack)-1]

		mp.currentPage.OnClose()
		previousPage.OnResume()
		mp.currentPage = previousPage
	}
}

func (mp *MainPage) popToFragment(pageID string) {

	// close current page and all pages before `pageID`
	if mp.currentPage != nil {
		mp.currentPage.OnClose()
	}

	for i := len(mp.pageBackStack) - 1; i >= 0; i-- {
		if mp.pageBackStack[i].ID() == pageID {
			var closedPages []load.Page
			mp.pageBackStack, closedPages = mp.pageBackStack[:i+1], mp.pageBackStack[i+1:]

			for j := len(closedPages) - 1; j >= 0; j-- {
				closedPages[j].OnClose()
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

func (mp *MainPage) Layout(gtx layout.Context) layout.Dimensions {
	if mp.currentPage != nil {
		mp.currentPage.Handle()
	}

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
			// global toasts. Stack toast on all pages and contents
			//TODO: show toasts here
			return layout.Dimensions{}

		}),
	)
}

func (mp *MainPage) LayoutUSDBalance(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if mp.usdExchangeSet && mp.dcrUsdtBittrex.LastTradeRate != "" {
				inset := layout.Inset{
					Top:  values.MarginPadding3,
					Left: values.MarginPadding8,
				}
				border := widget.Border{Color: mp.Theme.Color.Gray, CornerRadius: unit.Dp(8), Width: unit.Dp(0.5)}
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
										return layout.Inset{Right: values.MarginPadding16}.Layout(gtx,
											func(gtx C) D {
												return img.Layout24dp(gtx)
											})
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Center.Layout(gtx, func(gtx C) D {
											return components.LayoutBalance(gtx, mp.Load, mp.totalBalance.String())
										})
									}),
									layout.Rigid(func(gtx C) D {
										return mp.LayoutUSDBalance(gtx)
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
