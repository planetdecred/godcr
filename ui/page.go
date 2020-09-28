package ui

import (
	"image"

	"gioui.org/widget"

	"gioui.org/io/key"
	"gioui.org/layout"

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
	navigationArrowBack, verifyAction, actionDelete, actionLock,
	communicationComment, editorModeEdit, actionBackup, actionCheck,
	actionSwapVert, navigationCancel, notificationSync, imageBrightness1 *widget.Icon

	overviewIcon, walletIcon, receiveIcon, transactionIcon, sendIcon, syncingIcon image.Image
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
	}

	tabs := decredmaterial.NewTabs(win.theme)
	tabs.SetTabs([]decredmaterial.TabItem{
		decredmaterial.NewTabItem("Overview", &ic.overviewIcon),
		decredmaterial.NewTabItem("Transactions", &ic.transactionIcon),
		decredmaterial.NewTabItem("Wallets", &ic.walletIcon),
		decredmaterial.NewTabItem("More", &ic.walletIcon),

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
		keyEvents: win.keyEvents,
		clipboard: win.clipboard,
		states:    &win.states,
	}

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

func (page pageCommon) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	navs := []string{PageOverview, PageWallet, PageSend, PageReceive, PageTransactions}
	if page.navTab.ChangeEvent() {
		*page.page = navs[page.navTab.Selected]
	}

	page.navTab.Separator = true
	return page.navTab.Layout(gtx, func(gtx C) D {
		return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			return body(gtx)
		})
	})
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
