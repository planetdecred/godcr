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

type MainPage struct {
	Load                    *load.Load
	AppBarNavItems          []load.NavHandler
	DrawerNavItems          []load.NavHandler
	IsNavDrawerMinimized    bool
	MinimizeNavDrawerButton decredmaterial.IconButton
	MaximizeNavDrawerButton decredmaterial.IconButton

	AutoSync bool

	Current, Previous string
	Pages             map[string]load.Page
	SendPage          *send.Page // reuse value to keep data persistent onresume.

	// page state variables
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
	l.ChangeFragment = mp.Load.ChangeFragment
	l.SetReturnPage = mp.Load.SetReturnPage
	l.ReturnPage = &mp.Previous
	l.Page = &mp.Current

	toggleSync := func() {
		if mp.Load.WL.MultiWallet.IsConnectedToDecredNetwork() {
			mp.Load.WL.MultiWallet.CancelSync()
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
	mp.AppBarNavItems = []load.NavHandler{
		{
			Clickable: new(widget.Clickable),
			Image:     mp.Load.Icons.SendIcon,
			Page:      values.String(values.StrSend),
		},
		{
			Clickable: new(widget.Clickable),
			Image:     mp.Load.Icons.ReceiveIcon,
			Page:      values.String(values.StrReceive),
		},
	}

	mp.DrawerNavItems = []load.NavHandler{
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Load.Icons.OverviewIcon,
			ImageInactive: mp.Load.Icons.OverviewIconInactive,
			Page:          values.String(values.StrOverview),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Load.Icons.TransactionIcon,
			ImageInactive: mp.Load.Icons.TransactionIconInactive,
			Page:          values.String(values.StrTransactions),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Load.Icons.WalletIcon,
			ImageInactive: mp.Load.Icons.WalletIconInactive,
			Page:          values.String(values.StrWallets),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Load.Icons.TicketIcon,
			ImageInactive: mp.Load.Icons.TicketIconInactive,
			Page:          values.String(values.StrTickets),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Load.Icons.ProposalIconActive,
			ImageInactive: mp.Load.Icons.ProposalIconInactive,
			Page:          values.String(values.StrProposal),
		},
		{
			Clickable:     new(widget.Clickable),
			Image:         mp.Load.Icons.MoreIcon,
			ImageInactive: mp.Load.Icons.MoreIconInactive,
			Page:          values.String(values.StrMore),
		},
	}
}

func (mp *MainPage) OnResume() {
	// register for notifications
	mp.Load.WL.MultiWallet.AddAccountMixerNotificationListener(mp, MainPageID)
	mp.Load.WL.MultiWallet.Politeia.AddNotificationListener(mp, MainPageID)
	mp.Load.WL.MultiWallet.AddTxAndBlockNotificationListener(mp, MainPageID)
	mp.Load.WL.MultiWallet.AddSyncProgressListener(mp, MainPageID)

	mp.UpdateBalance()

	if pg, ok := mp.Pages[mp.Current]; ok {
		pg.OnResume()
	} else {
		mp.Load.ChangeFragment(NewOverviewPage(mp.Load), OverviewPageID)
	}

	if mp.AutoSync {
		mp.AutoSync = false
		mp.StartSyncing()
	}
}

func (mp *MainPage) UpdateBalance() {
	currencyExchangeValue := mp.Load.WL.Wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	mp.UsdExchangeSet = currencyExchangeValue == components.USDExchangeValue

	totalBalance, err := mp.CalculateTotalWalletsBalance()
	if err == nil {
		mp.TotalBalance = totalBalance

		if mp.UsdExchangeSet && mp.Load.DcrUsdtBittrex.LastTradeRate != "" {
			usdExchangeRate, err := strconv.ParseFloat(mp.Load.DcrUsdtBittrex.LastTradeRate, 64)
			if err == nil {
				balanceInUSD := totalBalance.ToCoin() * usdExchangeRate
				mp.TotalBalanceUSD = load.FormatUSDBalance(mp.Load.Printer, balanceInUSD)
			}
		}

	}
}

func (mp *MainPage) CalculateTotalWalletsBalance() (dcrutil.Amount, error) {
	totalBalance := int64(0)
	for _, wallet := range mp.Load.SortedWalletList() {
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
	for _, wal := range mp.Load.SortedWalletList() {
		if !wal.HasDiscoveredAccounts && wal.IsLocked() {
			mp.UnlockWalletForSyncing(wal)
			return
		}
	}

	err := mp.Load.WL.MultiWallet.SpvSync()
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
				err := mp.Load.WL.MultiWallet.UnlockWallet(wal.ID, []byte(password))
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

			mp.Load.SetReturnPage(mp.Current)
			mp.Load.ChangeFragment(pg, id)
		}
	}

	for i := range mp.DrawerNavItems {
		for mp.DrawerNavItems[i].Clickable.Clicked() {
			if i == OverviewNavID {
				mp.Load.ChangeFragment(NewOverviewPage(mp.Load), OverviewPageID)
			} else if i == TransactionsNavID {
				mp.Load.ChangeFragment(NewTransactionsPage(mp.Load), TransactionsPageID)
			} else if i == WalletsNavID {
				mp.Load.ChangeFragment(NewWalletPage(mp.Load), WalletPageID)
			} else if i == TicketsNavID {
				mp.Load.ChangeFragment(tickets.NewTicketPage(mp.Load), tickets.PageID)
			} else if i == ProposalsNavID {
				mp.Load.ChangeFragment(proposal.NewProposalsPage(mp.Load), proposal.ProposalsPageID)
			} else {
				mp.Load.ChangeFragment(NewMorePage(mp.Load), MorePageID)
			}
		}
	}
}

func (mp *MainPage) OnClose() {
	if pg, ok := mp.Pages[mp.Current]; ok {
		pg.OnClose()
	}
	mp.Load.WL.MultiWallet.RemoveAccountMixerNotificationListener(MainPageID)
	mp.Load.WL.MultiWallet.Politeia.RemoveNotificationListener(MainPageID)
	mp.Load.WL.MultiWallet.RemoveTxAndBlockNotificationListener(MainPageID)
	mp.Load.WL.MultiWallet.RemoveSyncProgressListener(MainPageID)
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
				return decredmaterial.Fill(gtx, mp.Load.Theme.Color.Surface)
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(mp.LayoutTopBar),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							card := mp.Load.Theme.Card()
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
	card := mp.Load.Theme.Card()
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
							return load.Container{Padding: layout.Inset{Right: h, Left: h, Top: v, Bottom: v}}.Layout(gtx,
								func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											img := mp.Load.Icons.Logo
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
									return load.Container{layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
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
															return mp.Load.Theme.Body1(mp.AppBarNavItems[i].Page).Layout(gtx)
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
				return mp.Load.Theme.Separator().Layout(gtx)
			}),
		)
	})
}

func (mp *MainPage) LayoutUSDBalance(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if mp.UsdExchangeSet && mp.Load.DcrUsdtBittrex.LastTradeRate != "" {
				inset := layout.Inset{
					Top:  values.MarginPadding3,
					Left: values.MarginPadding8,
				}
				border := widget.Border{Color: mp.Load.Theme.Color.Gray, CornerRadius: unit.Dp(8), Width: unit.Dp(0.5)}
				return inset.Layout(gtx, func(gtx C) D {
					padding := layout.Inset{
						Top:    values.MarginPadding3,
						Bottom: values.MarginPadding3,
						Left:   values.MarginPadding6,
						Right:  values.MarginPadding6,
					}
					return border.Layout(gtx, func(gtx C) D {
						return padding.Layout(gtx, func(gtx C) D {
							return mp.Load.Theme.Body2(mp.TotalBalanceUSD).Layout(gtx)
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
				background := mp.Load.Theme.Color.Surface
				if mp.DrawerNavItems[i].Page == mp.Current {
					background = mp.Load.Theme.Color.ActiveGray
				}
				txt := mp.Load.Theme.Label(values.TextSize16, mp.DrawerNavItems[i].Page)
				return decredmaterial.Clickable(gtx, mp.DrawerNavItems[i].Clickable, func(gtx C) D {
					card := mp.Load.Theme.Card()
					card.Color = background
					card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
					return card.Layout(gtx, func(gtx C) D {
						return load.Container{
							layout.Inset{
								Top:    values.MarginPadding16,
								Right:  values.MarginPadding24,
								Bottom: values.MarginPadding16,
								Left:   values.MarginPadding24,
							},
						}.Layout(gtx, func(gtx C) D {
							axis := layout.Horizontal
							leftInset := values.MarginPadding15
							width := components.NavDrawerWidth
							if mp.IsNavDrawerMinimized {
								axis = layout.Vertical
								txt.TextSize = values.TextSize10
								leftInset = values.MarginPadding0
								width = components.NavDrawerMinimizedWidth
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
										textColor := mp.Load.Theme.Color.Gray4
										if mp.DrawerNavItems[i].Page == mp.Current {
											textColor = mp.Load.Theme.Color.DeepBlue
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
