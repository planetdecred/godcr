package page

import (
	"image/color"
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
	Image         *widget.Image
	ImageInactive *widget.Image
	Title         string
	PageID        string
}

type MainPage struct {
	*load.Load
	appBarNavItems          []NavHandler
	drawerNavItems          []NavHandler
	isNavDrawerMinimized    bool
	minimizeNavDrawerButton decredmaterial.IconButton
	maximizeNavDrawerButton decredmaterial.IconButton

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

		minimizeNavDrawerButton: l.Theme.PlainIconButton(new(widget.Clickable), l.Icons.NavigationArrowBack),
		maximizeNavDrawerButton: l.Theme.PlainIconButton(new(widget.Clickable), l.Icons.NavigationArrowForward),
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

	iconColor := l.Theme.Color.Gray3
	mp.minimizeNavDrawerButton.Color, mp.maximizeNavDrawerButton.Color = iconColor, iconColor

	mp.initNavItems()

	mp.OnResume()

	return mp
}

func (mp *MainPage) ID() string {
	return MainPageID
}

func (mp *MainPage) initNavItems() {
	mp.appBarNavItems = []NavHandler{
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
	}

	mp.drawerNavItems = []NavHandler{
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Icons.OverviewIcon,
			ImageInactive: mp.Icons.OverviewIconInactive,
			Title:         values.String(values.StrOverview),
			PageID:        OverviewPageID,
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Icons.TransactionIcon,
			ImageInactive: mp.Icons.TransactionIconInactive,
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
	}
}

func (mp *MainPage) OnResume() {
	// register for notifications
	mp.WL.MultiWallet.AddAccountMixerNotificationListener(mp, MainPageID)
	mp.WL.MultiWallet.Politeia.AddNotificationListener(mp, MainPageID)
	mp.WL.MultiWallet.AddTxAndBlockNotificationListener(mp, MainPageID)
	mp.WL.MultiWallet.AddSyncProgressListener(mp, MainPageID)

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
	for mp.minimizeNavDrawerButton.Button.Clicked() {
		mp.isNavDrawerMinimized = true
	}

	for mp.maximizeNavDrawerButton.Button.Clicked() {
		mp.isNavDrawerMinimized = false
	}

	for i := range mp.appBarNavItems {
		for mp.appBarNavItems[i].Clickable.Clicked() {
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

	for i := range mp.drawerNavItems {
		for mp.drawerNavItems[i].Clickable.Clicked() {
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
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(mp.LayoutTopBar),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							card := mp.Theme.Card()
							card.Radius = decredmaterial.Radius(0)
							return card.Layout(gtx, mp.LayoutNavDrawer)
						}),
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

func (mp *MainPage) LayoutTopBar(gtx layout.Context) layout.Dimensions {
	card := mp.Theme.Card()
	card.Radius = decredmaterial.Radius(0)
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							h := values.MarginPadding24
							v := values.MarginPadding14
							// Balance container
							return components.Container{Padding: layout.Inset{Right: h, Left: h, Top: v, Bottom: v}}.Layout(gtx,
								func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											img := mp.Icons.Logo
											img.Scale = 1.0
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx,
												func(gtx C) D {
													return img.Layout(gtx)
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
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
								list := layout.List{Axis: layout.Horizontal}
								return list.Layout(gtx, len(mp.appBarNavItems), func(gtx C, i int) D {
									background := mp.Theme.Color.Surface
									if mp.appBarNavItems[i].PageID == mp.currentPageID() {
										background = mp.Theme.Color.ActiveGray
									}
									// header buttons container
									return decredmaterial.Clickable(gtx, mp.appBarNavItems[i].Clickable, func(gtx C) D {
										return mp.layoutCard(gtx, background, func(gtx C) D {
											return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
												return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
													layout.Rigid(func(gtx C) D {
														return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
															func(gtx C) D {
																return layout.Center.Layout(gtx, func(gtx C) D {
																	img := mp.appBarNavItems[i].Image
																	img.Scale = 1.0
																	return mp.appBarNavItems[i].Image.Layout(gtx)
																})
															})
													}),
													layout.Rigid(func(gtx C) D {
														return layout.Inset{
															Left: values.MarginPadding0,
														}.Layout(gtx, func(gtx C) D {
															return layout.Center.Layout(gtx, func(gtx C) D {
																return mp.Theme.Body1(mp.appBarNavItems[i].Title).Layout(gtx)
															})
														})
													}),
												)
											})
										})
									})
								})
							})
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return mp.Theme.Separator().Layout(gtx)
			}),
		)
	})
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

func (mp *MainPage) LayoutNavDrawer(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(mp.drawerNavItems), func(gtx C, i int) D {
				background := mp.Theme.Color.Surface
				if mp.drawerNavItems[i].PageID == mp.currentPageID() {
					background = mp.Theme.Color.ActiveGray
				}
				txt := mp.Theme.Label(values.TextSize16, mp.drawerNavItems[i].Title)
				return decredmaterial.Clickable(gtx, mp.drawerNavItems[i].Clickable, func(gtx C) D {
					return mp.layoutCard(gtx, background, func(gtx C) D {
						return components.Container{
							Padding: layout.Inset{
								Top:    values.MarginPadding16,
								Right:  values.MarginPadding24,
								Bottom: values.MarginPadding16,
								Left:   values.MarginPadding24,
							},
						}.Layout(gtx, func(gtx C) D {
							axis := layout.Horizontal
							leftInset := values.MarginPadding15
							width := NavDrawerWidth
							if mp.isNavDrawerMinimized {
								axis = layout.Vertical
								txt.TextSize = values.TextSize10
								leftInset = values.MarginPadding0
								width = NavDrawerMinimizedWidth
							}

							gtx.Constraints.Min.X = gtx.Px(width)
							return layout.Flex{Axis: axis}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									img := mp.drawerNavItems[i].ImageInactive
									if mp.drawerNavItems[i].PageID == mp.currentPageID() {
										img = mp.drawerNavItems[i].Image
									}
									return layout.Center.Layout(gtx, func(gtx C) D {
										img.Scale = 1.0
										return img.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Left: leftInset,
										Top:  values.MarginPadding4,
									}.Layout(gtx, func(gtx C) D {
										textColor := mp.Theme.Color.Gray4
										if mp.drawerNavItems[i].PageID == mp.currentPageID() {
											textColor = mp.Theme.Color.DeepBlue
										}
										txt.Color = textColor
										return layout.Center.Layout(gtx, txt.Layout)
									})
								}),
							)
						})
					})
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.SE.Layout(gtx, func(gtx C) D {
				btn := mp.minimizeNavDrawerButton
				if mp.isNavDrawerMinimized {
					btn = mp.maximizeNavDrawerButton
				}
				return btn.Layout(gtx)
			})
		}),
	)
}

func (mp *MainPage) layoutCard(gtx layout.Context, background color.NRGBA, body layout.Widget) layout.Dimensions {
	card := mp.Theme.Card()
	card.Color = background
	card.Radius = decredmaterial.Radius(0)
	return card.Layout(gtx, body)
}
