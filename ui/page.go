package ui

import (
	"image"
	"strings"

	"gioui.org/widget"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
	//"gioui.org/unit"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, contentClear, contentCreate, navigationCheck,
	contentSend, contentAddBox, contentRemove, toggleRadioButtonUnchecked,
	actionCheckCircle, contentCopy, actionInfo, navigationMore,
	navigationArrowBack, navigationArrowForward, verifyAction, actionDelete, actionLock,
	communicationComment, editorModeEdit, actionBackup, actionCheck,
	actionSwapVert, navigationCancel, notificationSync, imageBrightness1 *widget.Icon

	overviewIcon, overviewIconInactive, walletIcon, walletIconInactive, receiveIcon,
	transactionIcon, transactionIconInactive, sendIcon, syncingIcon, moreIcon, moreIconInactive,
	confirmIcon, pendingIcon, logo, redirectIcon image.Image
}

type navHandler struct {
	clickable     *widget.Clickable
	image         *widget.Image
	imageInactive *widget.Image
	page          string
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

	appBarNavItems []navHandler
}

type (
	C = layout.Context
	D = layout.Dimensions
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
		overviewIconInactive:       decredIcons["overview_inactive"],
		walletIcon:                 decredIcons["wallet"],
		walletIconInactive:         decredIcons["wallet_inactive"],
		receiveIcon:                decredIcons["receive"],
		transactionIcon:            decredIcons["transaction"],
		transactionIconInactive:    decredIcons["transaction_inactive"],
		sendIcon:                   decredIcons["send"],
		syncingIcon:                decredIcons["syncing"],
		moreIcon:                   decredIcons["more"],
		moreIconInactive:           decredIcons["more_inactive"],
		logo:                       decredIcons["logo"],
		confirmIcon:                decredIcons["confirmed"],
		pendingIcon:                decredIcons["pending"],
		redirectIcon:               decredIcons["redirect"],
	}
	win.theme.NavigationCheckIcon = ic.navigationCheck

	appBarNavItems := []navHandler{
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.sendIcon)},
			page:      PageSend,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.receiveIcon)},
			page:      PageReceive,
		},
	}

	tabs := decredmaterial.NewTabs(win.theme, true)
	tabs.SetTabs([]decredmaterial.TabItem{
		decredmaterial.NewTabItem(PageOverview, &ic.overviewIcon, &ic.overviewIconInactive),
		decredmaterial.NewTabItem(PageTransactions, &ic.transactionIcon, &ic.transactionIconInactive),
		decredmaterial.NewTabItem(PageWallet, &ic.walletIcon, &ic.walletIconInactive),
		decredmaterial.NewTabItem(PageMore, &ic.moreIcon, &ic.moreIconInactive),
	})

	common := pageCommon{
		wallet:          win.wallet,
		info:            win.walletInfo,
		selectedWallet:  &win.selected,
		selectedAccount: &win.selectedAccount,
		theme:           win.theme,
		icons:           ic,
		page:            &win.current,
		navTab:          tabs,
		walletTabs:      win.walletTabs,
		accountTabs:     win.accountTabs,
		errorChannels: map[string]chan error{
			PageSignMessage:    make(chan error),
			PageCreateRestore:  make(chan error),
			PageWallet:         make(chan error),
			PageWalletAccounts: make(chan error),
		},
		keyEvents:      win.keyEvents,
		clipboard:      win.clipboard,
		states:         &win.states,
		appBarNavItems: appBarNavItems,
	}

	win.pages = make(map[string]layout.Widget)
	win.pages[PageWallet] = win.WalletPage(common)
	win.pages[PageOverview] = win.OverviewPage(common)
	win.pages[PageTransactions] = win.TransactionsPage(common)
	win.pages[PageMore] = win.MorePage(decredIcons, common)
	win.pages[PageCreateRestore] = win.CreateRestorePage(common)
	win.pages[PageReceive] = win.ReceivePage(common)
	win.pages[PageSend] = win.SendPage(common)
	win.pages[PageTransactionDetails] = win.TransactionDetailsPage(common)
	win.pages[PageSignMessage] = win.SignMessagePage(common)
	win.pages[PageVerifyMessage] = win.VerifyMessagePage(common)
	win.pages[PageWalletPassphrase] = win.WalletPassphrasePage(common)
	win.pages[PageWalletAccounts] = win.WalletAccountPage(common)
	win.pages[PageSeedBackup] = win.BackupPage(common)
	win.pages[PageSettings] = win.SettingsPage(common)
	win.pages[PageSecurityTools] = win.SecurityToolsPage(common)
	win.pages[PagePoliteia] = win.PoliteiaPage(common)
	win.pages[PageDebug] = win.DebugPage(common)
	win.pages[PageAbout] = win.AboutPage(common)
	win.pages[PageHelp] = win.HelpPage(common)
}

func (page pageCommon) ChangePage(pg string) {
	*page.page = pg
}

func (page pageCommon) handleNavEvents() {
	for i := range page.appBarNavItems {
		for page.appBarNavItems[i].clickable.Clicked() {
			*page.page = page.appBarNavItems[i].page
		}
	}

	if page.navTab.ChangeEvent() {
		selectedID := page.navTab.SelectedID()
		if selectedID != "" {
			*page.page = selectedID
		}
	}
}

func (page pageCommon) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	page.handleNavEvents()

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return page.layoutAppBar(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return page.navTab.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return body(gtx)
				})
			})
		}),
	)
}

func (page pageCommon) layoutAppBar(gtx layout.Context) layout.Dimensions {
	return decredmaterial.Card{Color: page.theme.Color.Surface}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							img := widget.Image{Src: paint.NewImageOp(page.icons.logo)}
							img.Scale = 0.085

							return img.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									return page.layoutBalance(gtx, page.info.TotalBalance)
								})
							})
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.E.Layout(gtx, func(gtx C) D {
								list := layout.List{Axis: layout.Horizontal}
								return list.Layout(gtx, len(page.appBarNavItems), func(gtx C, i int) D {
									return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
										return decredmaterial.Clickable(gtx, page.appBarNavItems[i].clickable, func(gtx C) D {
											return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													page.appBarNavItems[i].image.Scale = 0.05

													return layout.Center.Layout(gtx, func(gtx C) D {
														return page.appBarNavItems[i].image.Layout(gtx)
													})
												}),
												layout.Rigid(func(gtx C) D {
													return layout.Inset{
														Left: values.MarginPadding10,
													}.Layout(gtx, func(gtx C) D {
														return layout.Center.Layout(gtx, func(gtx C) D {
															return page.theme.Body1(page.appBarNavItems[i].page).Layout(gtx)
														})
													})
												}),
											)
										})
									})
								})
							})
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				l := page.theme.Line()
				l.Color = page.theme.Color.Background
				l.Width = gtx.Constraints.Min.X
				l.Height = 2
				return l.Layout(gtx)
			}),
		)
	})
}

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page pageCommon) layoutBalance(gtx layout.Context, amount string) layout.Dimensions {
	mainText, subText := page.breakBalance(amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return page.theme.H5(mainText).Layout(gtx)
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
