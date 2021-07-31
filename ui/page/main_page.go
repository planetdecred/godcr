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
)

var (
	NavDrawerWidth          = unit.Value{U: unit.UnitDp, V: 160}
	NavDrawerMinimizedWidth = unit.Value{U: unit.UnitDp, V: 72}
)

type NavHandler struct {
	Clickable     *widget.Clickable
	Image         *widget.Image
	ImageInactive *widget.Image
	Page          string
}

type MainPage struct {
	*load.Load
	AppBarNavItems          []NavHandler
	DrawerNavItems          []NavHandler
	IsNavDrawerMinimized    bool
	MinimizeNavDrawerButton decredmaterial.IconButton
	MaximizeNavDrawerButton decredmaterial.IconButton

	AutoSync bool

	Current, Previous string
	Pages             map[string]load.Page
	SendPage          *send.Page // reuse value to keep data persistent onresume.

	// page state variables
	DcrUsdtBittrex  load.DCRUSDTBittrex
	UsdExchangeSet  bool
	TotalBalance    dcrutil.Amount
	TotalBalanceUSD string
}

func NewMainPage(l *load.Load) *MainPage {

	mp := &MainPage{
		Load:     l,
		AutoSync: true,
		Pages:    make(map[string]load.Page),

		MinimizeNavDrawerButton: l.Theme.PlainIconButton(new(widget.Clickable), l.Icons.NavigationArrowBack),
		MaximizeNavDrawerButton: l.Theme.PlainIconButton(new(widget.Clickable), l.Icons.NavigationArrowForward),
	}

	// init shared page functions
	// todo: common methods will be removed when all pages have been moved to the page package
	l.ChangeFragment = mp.changeFragment
	l.SetReturnPage = mp.setReturnPage
	l.ReturnPage = &mp.Previous
	l.Page = &mp.Current

	toggleSync := func() {
		if mp.WL.MultiWallet.IsConnectedToDecredNetwork() {
			mp.WL.MultiWallet.CancelSync()
		} else {
			mp.StartSyncing()
		}
	}
	// todo: to be removed when all pages have been migrated
	l.ToggleSync = toggleSync
	l.ToggleSync = toggleSync

	iconColor := l.Theme.Color.Gray3
	mp.MinimizeNavDrawerButton.Color, mp.MaximizeNavDrawerButton.Color = iconColor, iconColor

	mp.initNavItems()

	mp.OnResume()

	return mp
}

func (mp *MainPage) initNavItems() {
	mp.AppBarNavItems = []NavHandler{
		{
			Clickable: new(widget.Clickable),
			Image:     mp.Icons.SendIcon,
			Page:      values.String(values.StrSend),
		},
		{
			Clickable: new(widget.Clickable),
			Image:     mp.Icons.ReceiveIcon,
			Page:      values.String(values.StrReceive),
		},
	}

	mp.DrawerNavItems = []NavHandler{
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Icons.OverviewIcon,
			ImageInactive: mp.Icons.OverviewIconInactive,
			Page:          values.String(values.StrOverview),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Icons.TransactionIcon,
			ImageInactive: mp.Icons.TransactionIconInactive,
			Page:          values.String(values.StrTransactions),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Icons.WalletIcon,
			ImageInactive: mp.Icons.WalletIconInactive,
			Page:          values.String(values.StrWallets),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Icons.TicketIcon,
			ImageInactive: mp.Icons.TicketIconInactive,
			Page:          values.String(values.StrTickets),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Icons.ProposalIconActive,
			ImageInactive: mp.Icons.ProposalIconInactive,
			Page:          values.String(values.StrProposal),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Icons.MoreIcon,
			ImageInactive: mp.Icons.MoreIconInactive,
			Page:          values.String(values.StrMore),
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

	if pg, ok := mp.Pages[mp.Current]; ok {
		pg.OnResume()
	} else {
		mp.ChangeFragment(NewOverviewPage(mp.Load), OverviewPageID)
	}

	if mp.AutoSync {
		mp.AutoSync = false
		mp.StartSyncing()
		go mp.WL.MultiWallet.Politeia.Sync()
	}
}

func (mp *MainPage) UpdateBalance() {
	currencyExchangeValue := mp.WL.Wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	mp.UsdExchangeSet = currencyExchangeValue == components.USDExchangeValue

	totalBalance, err := mp.CalculateTotalWalletsBalance()
	if err == nil {
		mp.TotalBalance = totalBalance

		if mp.UsdExchangeSet && mp.DcrUsdtBittrex.LastTradeRate != "" {
			usdExchangeRate, err := strconv.ParseFloat(mp.DcrUsdtBittrex.LastTradeRate, 64)
			if err == nil {
				balanceInUSD := totalBalance.ToCoin() * usdExchangeRate
				mp.TotalBalanceUSD = load.FormatUSDBalance(mp.Printer, balanceInUSD)
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
	for mp.MinimizeNavDrawerButton.Button.Clicked() {
		mp.IsNavDrawerMinimized = true
	}

	for mp.MaximizeNavDrawerButton.Button.Clicked() {
		mp.IsNavDrawerMinimized = false
	}

	for i := range mp.AppBarNavItems {
		for mp.AppBarNavItems[i].Clickable.Clicked() {
			var pg load.Page
			var id string
			if i == 0 {
				if mp.SendPage == nil {
					mp.SendPage = send.NewSendPage(mp.Load)
				}

				pg = mp.SendPage
				id = send.PageID
			} else {
				pg = NewReceivePage(mp.Load)
				id = ReceivePageID
			}

			mp.SetReturnPage(mp.Current)
			mp.ChangeFragment(pg, id)
		}
	}

	for i := range mp.DrawerNavItems {
		for mp.DrawerNavItems[i].Clickable.Clicked() {
			if i == OverviewNavID {
				mp.ChangeFragment(NewOverviewPage(mp.Load), OverviewPageID)
			} else if i == TransactionsNavID {
				mp.ChangeFragment(NewTransactionsPage(mp.Load), TransactionsPageID)
			} else if i == WalletsNavID {
				mp.ChangeFragment(NewWalletPage(mp.Load), WalletPageID)
			} else if i == TicketsNavID {
				mp.ChangeFragment(tickets.NewTicketPage(mp.Load), tickets.PageID)
			} else if i == ProposalsNavID {
				mp.ChangeFragment(proposal.NewProposalsPage(mp.Load), proposal.ProposalsPageID)
			} else {
				mp.ChangeFragment(NewMorePage(mp.Load), MorePageID)
			}
		}
	}
}

func (mp *MainPage) OnClose() {
	if pg, ok := mp.Pages[mp.Current]; ok {
		pg.OnClose()
	}
	mp.WL.MultiWallet.RemoveAccountMixerNotificationListener(MainPageID)
	mp.WL.MultiWallet.Politeia.RemoveNotificationListener(MainPageID)
	mp.WL.MultiWallet.RemoveTxAndBlockNotificationListener(MainPageID)
	mp.WL.MultiWallet.RemoveSyncProgressListener(MainPageID)
}

func (mp *MainPage) changeFragment(page load.Page, id string) {
	mp.Pages[id] = page
	mp.changePage(id)
}

func (mp *MainPage) changePage(page string) {
	if pg, ok := mp.Pages[mp.Current]; ok {
		pg.OnClose()
	}

	if pg, ok := mp.Pages[page]; ok {
		pg.OnResume()
		mp.Current = page
	}
}

func (mp *MainPage) setReturnPage(from string) {
	mp.Previous = from
}

func (mp *MainPage) Layout(gtx layout.Context) layout.Dimensions {
	mp.Pages[mp.Current].Handle()

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			// fill the entire window with a color if a user has no wallet created
			if mp.Current == CreateRestorePageID {
				return decredmaterial.Fill(gtx, mp.Theme.Color.Surface)
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(mp.LayoutTopBar),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							card := mp.Theme.Card()
							card.Radius = decredmaterial.CornerRadius{}
							return card.Layout(gtx, mp.LayoutNavDrawer)
						}),
						layout.Rigid(mp.Pages[mp.Current].Layout),
					)
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			// stack the page content on the entire window if a user has no wallet
			if mp.Current == CreateRestorePageID {
				return mp.Pages[mp.Current].Layout(gtx)
			}
			return layout.Dimensions{}
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
	card.Radius = decredmaterial.CornerRadius{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							h := values.MarginPadding24
							v := values.MarginPadding16
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
												return components.LayoutBalance(gtx, mp.Load, mp.TotalBalance.String())
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
								return list.Layout(gtx, len(mp.AppBarNavItems), func(gtx C, i int) D {
									// header buttons container
									return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
										return decredmaterial.Clickable(gtx, mp.AppBarNavItems[i].Clickable, func(gtx C) D {
											return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
														func(gtx C) D {
															return layout.Center.Layout(gtx, func(gtx C) D {
																img := mp.AppBarNavItems[i].Image
																img.Scale = 1.0
																return mp.AppBarNavItems[i].Image.Layout(gtx)
															})
														})
												}),
												layout.Rigid(func(gtx C) D {
													return layout.Inset{
														Left: values.MarginPadding0,
													}.Layout(gtx, func(gtx C) D {
														return layout.Center.Layout(gtx, func(gtx C) D {
															return mp.Theme.Body1(mp.AppBarNavItems[i].Page).Layout(gtx)
														})
													})
												}),
											)
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
			if mp.UsdExchangeSet && mp.DcrUsdtBittrex.LastTradeRate != "" {
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
							return mp.Theme.Body2(mp.TotalBalanceUSD).Layout(gtx)
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
			return list.Layout(gtx, len(mp.DrawerNavItems), func(gtx C, i int) D {
				background := mp.Theme.Color.Surface
				if mp.DrawerNavItems[i].Page == mp.Current {
					background = mp.Theme.Color.ActiveGray
				}
				txt := mp.Theme.Label(values.TextSize16, mp.DrawerNavItems[i].Page)
				return decredmaterial.Clickable(gtx, mp.DrawerNavItems[i].Clickable, func(gtx C) D {
					card := mp.Theme.Card()
					card.Color = background
					card.Radius = decredmaterial.Radius(0)
					return card.Layout(gtx, func(gtx C) D {
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
							if mp.IsNavDrawerMinimized {
								axis = layout.Vertical
								txt.TextSize = values.TextSize10
								leftInset = values.MarginPadding0
								width = NavDrawerMinimizedWidth
							}

							gtx.Constraints.Min.X = gtx.Px(width)
							return layout.Flex{Axis: axis}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									img := mp.DrawerNavItems[i].ImageInactive
									if mp.DrawerNavItems[i].Page == mp.Current {
										img = mp.DrawerNavItems[i].Image
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
										if mp.DrawerNavItems[i].Page == mp.Current {
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
				btn := mp.MinimizeNavDrawerButton
				if mp.IsNavDrawerMinimized {
					btn = mp.MaximizeNavDrawerButton
				}
				return btn.Layout(gtx)
			})
		}),
	)
}
