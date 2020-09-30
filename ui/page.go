package ui

import (
	"image"
	"strings"

	"gioui.org/widget"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/ui/values"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, contentClear, contentCreate, navigationCheck,
	contentSend, contentAddBox, contentRemove, toggleRadioButtonUnchecked,
	actionCheckCircle, contentCopy, actionInfo, navigationMore,
	navigationArrowBack, navigationArrowForward, verifyAction, actionDelete, actionLock,
	communicationComment, editorModeEdit, actionBackup, actionCheck,
	actionSwapVert, navigationCancel, notificationSync, imageBrightness1 *widget.Icon

	overviewIcon, walletIcon, receiveIcon, transactionIcon, sendIcon, syncingIcon, logo image.Image
}

type navHandler struct {
	clickable *widget.Clickable
	image     *widget.Image
	page      string
	isActive  bool
}

type pageCommon struct {
	wallet          *wallet.Wallet
	info            *wallet.MultiWalletInfo
	selectedWallet  *int
	selectedAccount *int
	theme           *decredmaterial.Theme
	icons           pageIcons
	page            *string
	navTab          *decredmaterial.Tabs
	walletTabs      *decredmaterial.Tabs
	accountTabs     *decredmaterial.Tabs
	errorChannels   map[string]chan error
	keyEvents       chan *key.Event
	clipboard       chan interface{}
	states          *states

	appBarNavItems          []navHandler
	drawerNavItems          []navHandler
	isNavDrawerMinimized    *bool
	minimizeNavDrawerButton decredmaterial.IconButton
	maximizeNavDrawerButton decredmaterial.IconButton
}

type (
	C = layout.Context
	D = layout.Dimensions
)

const (
	navDrawerWidth          = 190
	navDrawerMinimizedWidth = 115
)

func (win *Window) addPages(decredIcons map[string]image.Image) {
	ic := pageIcons{
		contentAdd:                 mustIcon(widget.NewIcon(icons.ContentAdd)),
		contentClear:               mustIcon(widget.NewIcon(icons.ContentClear)),
		contentCreate:              mustIcon(widget.NewIcon(icons.ContentCreate)),
		navigationCheck:            mustIcon(widget.NewIcon(icons.NavigationCheck)),
		contentSend:                mustIcon(widget.NewIcon(icons.ContentSend)),
		contentAddBox:              mustIcon(widget.NewIcon(icons.ContentAddBox)),
		contentRemove:              mustIcon(widget.NewIcon(icons.ContentRemove)),
		toggleRadioButtonUnchecked: mustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked)),
		actionCheckCircle:          mustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		navigationArrowBack:        mustIcon(widget.NewIcon(icons.NavigationArrowBack)),
		navigationArrowForward:     mustIcon(widget.NewIcon(icons.NavigationArrowForward)),
		contentCopy:                mustIcon(widget.NewIcon(icons.NavigationMoreVert)),
		actionInfo:                 mustIcon(widget.NewIcon(icons.ActionInfo)),
		navigationMore:             mustIcon(widget.NewIcon(icons.NavigationMoreVert)),
		actionDelete:               mustIcon(widget.NewIcon(icons.ActionDelete)),
		communicationComment:       mustIcon(widget.NewIcon(icons.CommunicationComment)),
		verifyAction:               mustIcon(widget.NewIcon(icons.ActionVerifiedUser)),
		editorModeEdit:             mustIcon(widget.NewIcon(icons.EditorModeEdit)),
		actionLock:                 mustIcon(widget.NewIcon(icons.ActionLock)),
		actionBackup:               mustIcon(widget.NewIcon(icons.ActionSettingsBackupRestore)),
		actionCheck:                mustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		actionSwapVert:             mustIcon(widget.NewIcon(icons.ActionSwapVert)),
		navigationCancel:           mustIcon(widget.NewIcon(icons.NavigationCancel)),
		notificationSync:           mustIcon(widget.NewIcon(icons.NotificationSync)),
		imageBrightness1:           mustIcon(widget.NewIcon(icons.ImageBrightness1)),
		overviewIcon:               decredIcons["overview"],
		walletIcon:                 decredIcons["wallet"],
		receiveIcon:                decredIcons["receive"],
		transactionIcon:            decredIcons["transaction"],
		sendIcon:                   decredIcons["send"],
		syncingIcon:                decredIcons["syncing"],
		logo:                       decredIcons["logo"],
	}

	appBarNavItems := []navHandler{
		{
			clickable: new(widget.Clickable),
			page:      PageSend,
		},
		{
			clickable: new(widget.Clickable),
			page:      PageReceive,
		},
	}

	drawerNavItems := []navHandler{
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.overviewIcon)},
			page:      PageOverview,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.transactionIcon)},
			page:      PageTransactions,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.walletIcon)},
			page:      PageWallet,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.walletIcon)},
			page:      PageMore,
		},
	}

	common := pageCommon{
		wallet:          win.wallet,
		info:            win.walletInfo,
		selectedWallet:  &win.selected,
		selectedAccount: &win.selectedAccount,
		theme:           win.theme,
		icons:           ic,
		page:            &win.current,
		walletTabs:      win.walletTabs,
		accountTabs:     win.accountTabs,
		errorChannels: map[string]chan error{
			PageSignMessage:    make(chan error),
			PageCreateRestore:  make(chan error),
			PageWallet:         make(chan error),
			PageWalletAccounts: make(chan error),
		},
		keyEvents:               win.keyEvents,
		clipboard:               win.clipboard,
		states:                  &win.states,
		appBarNavItems:          appBarNavItems,
		drawerNavItems:          drawerNavItems,
		minimizeNavDrawerButton: win.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowBack),
		maximizeNavDrawerButton: win.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowForward),
	}

	isNavDrawerMinimized := false
	common.isNavDrawerMinimized = &isNavDrawerMinimized
	common.minimizeNavDrawerButton.Color = common.theme.Color.Gray
	common.maximizeNavDrawerButton.Color = common.theme.Color.Gray

	win.pages = make(map[string]layout.Widget)
	win.pages[PageWallet] = win.WalletPage(common)
	win.pages[PageOverview] = win.OverviewPage(common)
	win.pages[PageTransactions] = win.TransactionsPage(common)
	win.pages[PageCreateRestore] = win.CreateRestorePage(common)
	win.pages[PageReceive] = win.ReceivePage(common)
	win.pages[PageSend] = win.SendPage(common)
	win.pages[PageTransactionDetails] = win.TransactionPage(common)
	win.pages[PageSignMessage] = win.SignMessagePage(common)
	win.pages[PageVerifyMessage] = win.VerifyMessagePage(common)
	win.pages[PageWalletPassphrase] = win.WalletPassphrasePage(common)
	win.pages[PageWalletAccounts] = win.WalletAccountPage(common)
	win.pages[PageSeedBackup] = win.BackupPage(common)
}

func (page pageCommon) ChangePage(pg string) {
	*page.page = pg
}

func (page pageCommon) handleNavEvents() {
	for page.minimizeNavDrawerButton.Button.Clicked() {
		*page.isNavDrawerMinimized = true
	}

	for page.maximizeNavDrawerButton.Button.Clicked() {
		*page.isNavDrawerMinimized = false
	}

	for i := range page.appBarNavItems {
		for page.appBarNavItems[i].clickable.Clicked() {
			*page.page = page.appBarNavItems[i].page
		}
	}

	for i := range page.drawerNavItems {
		for page.drawerNavItems[i].clickable.Clicked() {
			*page.page = page.drawerNavItems[i].page
		}
	}
}

func (page pageCommon) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	page.handleNavEvents()

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return page.layoutAppBar(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					width := navDrawerWidth
					if *page.isNavDrawerMinimized {
						width = navDrawerMinimizedWidth
					}
					gtx.Constraints.Max.X = width
					return decredmaterial.Card{Color: page.theme.Color.Surface}.Layout(gtx, func(gtx C) D {
						page.layoutNavDrawer(gtx)
						return layout.Dimensions{Size: gtx.Constraints.Max}
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx C) D {
						return body(gtx)
					})
				}),
			)
		}),
	)
}

func (page pageCommon) layoutAppBar(gtx layout.Context) layout.Dimensions {
	return decredmaterial.Card{Color: page.theme.Color.Surface}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Inset{
			Top:    unit.Dp(7),
			Bottom: unit.Dp(7),
			Left:   unit.Dp(18),
			Right:  unit.Dp(18),
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					img := widget.Image{Src: paint.NewImageOp(page.icons.logo)}
					img.Scale = 0.085

					return img.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx C) D {
						return page.layoutBalance(gtx, page.info.TotalBalance)
					})
				}),
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.E.Layout(gtx, func(gtx C) D {
						list := layout.List{Axis: layout.Horizontal}
						return list.Layout(gtx, len(page.appBarNavItems), func(gtx C, i int) D {
							return layout.Inset{
								Left: unit.Dp(10),
							}.Layout(gtx, func(gtx C) D {
								return decredmaterial.Clickable(gtx, page.appBarNavItems[i].clickable, func(gtx C) D {
									return page.theme.H6(page.appBarNavItems[i].page).Layout(gtx)
								})
							})
						})
					})
				}),
			)
		})
	})
}

func (page pageCommon) layoutNavDrawer(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(page.drawerNavItems), func(gtx C, i int) D {
				return decredmaterial.Clickable(gtx, page.drawerNavItems[i].clickable, func(gtx C) D {
					background := page.theme.Color.Surface
					if page.drawerNavItems[i].page == *page.page {
						background = page.theme.Color.Background
					}

					return decredmaterial.Card{Color: background}.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Stack{}.Layout(gtx,
							layout.Stacked(func(gtx C) D {
								return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx C) D {
									axis := layout.Horizontal
									leftInset := float32(15)
									if *page.isNavDrawerMinimized {
										axis = layout.Vertical
										leftInset = 0
									}

									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return layout.Flex{Axis: axis}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											page.drawerNavItems[i].image.Scale = 0.05

											return layout.Center.Layout(gtx, func(gtx C) D {
												return page.drawerNavItems[i].image.Layout(gtx)
											})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Left: unit.Dp(leftInset),
											}.Layout(gtx, func(gtx C) D {
												return layout.Center.Layout(gtx, func(gtx C) D {
													return page.theme.Body1(page.drawerNavItems[i].page).Layout(gtx)
												})
											})
										}),
									)
								})
							}),
						)
					})
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.SE.Layout(gtx, func(gtx C) D {
				btn := page.minimizeNavDrawerButton
				if *page.isNavDrawerMinimized {
					btn = page.maximizeNavDrawerButton
				}
				return btn.Layout(gtx)
			})
		}),
	)
}

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page pageCommon) layoutBalance(gtx layout.Context, amount string) layout.Dimensions {
	mainText, subText := page.breakBalance(amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return page.theme.H4(mainText).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return page.theme.Body1(subText).Layout(gtx)
		}),
	)
}

// breakBalance takes the balance string and returns it in two slices
func (page pageCommon) breakBalance(balance string) (b1, b2 string) {
	balanceParts := strings.Split(balance, ".")
	if len(balanceParts) == 1 {
		return balanceParts[0], ""
	}
	b1 = balanceParts[0]
	b2 = balanceParts[1]
	b1 = b1 + "." + b2[:2]
	b2 = b2[2:]
	return
}

func (page pageCommon) Modal(gtx layout.Context, body layout.Dimensions, modal layout.Dimensions) layout.Dimensions {
	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return body
		}),
		layout.Stacked(func(gtx C) D {
			return modal
		}),
	)
	return dims
}

func (page pageCommon) LayoutWithWallets(gtx layout.Context, body layout.Widget) layout.Dimensions {
	bd := func(gtx C) D {
		if page.walletTabs.ChangeEvent() {
			*page.selectedWallet = page.walletTabs.Selected
			*page.selectedAccount = 0
			page.accountTabs.Selected = 0

			accounts := make([]decredmaterial.TabItem, len(page.info.Wallets[*page.selectedWallet].Accounts))
			for i, account := range page.info.Wallets[*page.selectedWallet].Accounts {
				if account.Name == "imported" {
					continue
				}
				accounts[i] = decredmaterial.TabItem{
					Title: page.info.Wallets[*page.selectedWallet].Accounts[i].Name,
				}
			}
			page.accountTabs.SetTabs(accounts)
		}
		return page.walletTabs.Layout(gtx, body)
	}
	return page.Layout(gtx, bd)
}

func (page pageCommon) LayoutWithAccounts(gtx layout.Context, body layout.Widget) layout.Dimensions {
	if page.accountTabs.ChangeEvent() {
		*page.selectedAccount = page.accountTabs.Selected
	}

	return page.LayoutWithWallets(gtx, func(gtx C) D {
		return page.accountTabs.Layout(gtx, body)
	})
}

func (page pageCommon) SelectedAccountLayout(gtx layout.Context) layout.Dimensions {
	current := page.info.Wallets[*page.selectedWallet]
	account := page.info.Wallets[*page.selectedWallet].Accounts[*page.selectedAccount]

	selectedDetails := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return page.theme.H6(account.Name).Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return page.theme.H6(dcrutil.Amount(account.SpendableBalance).String()).Layout(gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return page.theme.Body2(current.Name).Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return page.theme.Body2(current.Balance).Layout(gtx)
							})
						}),
					)
				}),
			)
		})
	}
	return decredmaterial.Card{}.Layout(gtx, selectedDetails)
}

func toMax(gtx layout.Context) {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
