package ui

import (
	"image"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, contentClear, contentCreate, navigationCheck,
	contentSend, contentAddBox, contentRemove, toggleRadioButtonUnchecked,
	actionCheckCircle, contentCopy, actionInfo, navigationMore,
	navigationArrowBack *decredmaterial.Icon
	overviewIcon, walletIcon image.Image
}

type pageCommon struct {
	wallet          *wallet.Wallet
	info            *wallet.MultiWalletInfo
	selectedWallet  *int
	selectedAccount *int
	gtx             *layout.Context
	theme           *decredmaterial.Theme
	icons           pageIcons
	page            *string
	navTab          *decredmaterial.Tabs
	walletsTab      *decredmaterial.Tabs
	accountsTab     *decredmaterial.Tabs
	errorChannels   map[string]chan error
	keyEvents	   chan *key.Event
}

func (win *Window) addPages() {
	icons := pageIcons{
		contentAdd:                 mustIcon(decredmaterial.NewIcon(icons.ContentAdd)),
		contentClear:               mustIcon(decredmaterial.NewIcon(icons.ContentClear)),
		contentCreate:              mustIcon(decredmaterial.NewIcon(icons.ContentCreate)),
		navigationCheck:            mustIcon(decredmaterial.NewIcon(icons.NavigationCheck)),
		contentSend:                mustIcon(decredmaterial.NewIcon(icons.ContentSend)),
		contentAddBox:              mustIcon(decredmaterial.NewIcon(icons.ContentAddBox)),
		contentRemove:              mustIcon(decredmaterial.NewIcon(icons.ContentRemove)),
		toggleRadioButtonUnchecked: mustIcon(decredmaterial.NewIcon(icons.ToggleRadioButtonUnchecked)),
		actionCheckCircle:          mustIcon(decredmaterial.NewIcon(icons.ActionCheckCircle)),
		navigationArrowBack:        mustIcon(decredmaterial.NewIcon(icons.NavigationArrowBack)),
		contentCopy:                mustIcon(decredmaterial.NewIcon(icons.NavigationMoreVert)),
		actionInfo:                 mustIcon(decredmaterial.NewIcon(icons.ActionInfo)),
		navigationMore:             mustIcon(decredmaterial.NewIcon(icons.NavigationMoreVert)),
		overviewIcon:               mustDcrIcon(decredmaterial.NewDcrIcon(decredmaterial.OverviewIcon)),
		walletIcon:                 mustDcrIcon(decredmaterial.NewDcrIcon(decredmaterial.WalletIcon)),
	}
	tabs := decredmaterial.NewTabs()
	tabs.SetTabs([]decredmaterial.TabItem{
		{
			Label: win.theme.Body1("Overview"),
			Icon:  icons.overviewIcon,
		},
		{
			Label: win.theme.Body1("Wallets"),
			Icon:  icons.walletIcon,
		},
		{
			Label: win.theme.Body1("Transactions"),
			Icon:  icons.overviewIcon,
		},
		{
			Label: win.theme.Body1("Receive"),
			Icon:  icons.overviewIcon,
		},
		{
			Label: win.theme.Body1("Settings"),
			Icon:  icons.overviewIcon,
		},
	})

	accountsTab := decredmaterial.NewTabs()
	accountsTab.Position = decredmaterial.Top
	accountsTab.Separator = false
	common := pageCommon{
		wallet:          win.wallet,
		info:            win.walletInfo,
		selectedWallet:  &win.selected,
		selectedAccount: &win.selectedAccount,
		gtx:             win.gtx,
		theme:           win.theme,
		icons:           icons,
		page:            &win.current,
		navTab:          tabs,
		walletsTab:      decredmaterial.NewTabs(),
		accountsTab:     accountsTab,
		errorChannels: map[string]chan error{
			PageSignMessage: make(chan error),
		},
		keyEvents: 		win.keyEvents,
		//cancelDialogW:  win.theme.PlainIconButton(icons.contentClear),
	}

	win.pages = make(map[string]layout.Widget)
	win.pages[PageWallet] = win.WalletPage(common)
	win.pages[PageOverview] = win.OverviewPage(common)
	win.pages[PageTransactions] = win.TransactionsPage(common)
	win.pages[PageReceive] = win.Receive
	win.pages[PageCreateRestore] = win.CreateRestorePage(common)
	win.pages[PageReceive] = win.ReceivePage(common)
	win.pages[PageSend] = win.SendPage
	win.pages[PageTransactionDetails] = win.TransactionPage(common)
	win.pages[PageSignMessage] = win.SignMessagePage(common)
	win.pages[PageVerifyMessage] = win.VerifyMessagePage(common)
}

func (page pageCommon) Layout(gtx *layout.Context, body layout.Widget) {
	navs := []string{PageOverview, PageWallet, PageTransactions, PageReceive, PageOverview}
	toMax(gtx)
	page.navTab.Separator = true
	page.navTab.Layout(gtx, body)

	if page.navTab.Changed() {
		*page.page = navs[page.navTab.Selected]
	}
}

func (page pageCommon) LayoutWithWallets(gtx *layout.Context, body layout.Widget) {
	wallets := make([]decredmaterial.TabItem, len(page.info.Wallets))
	for i := range page.info.Wallets {
		wallets[i] = decredmaterial.TabItem{
			Label: page.theme.Body1(page.info.Wallets[i].Name),
		}
	}
	page.walletsTab.SetTabs(wallets)
	page.walletsTab.Position = decredmaterial.Top
	if page.accountsTab.Changed() {
		*page.selectedAccount = page.accountsTab.Selected
	}

	accounts := make([]decredmaterial.TabItem, len(page.info.Wallets[*page.selectedWallet].Accounts))
	for i := range page.info.Wallets[*page.selectedWallet].Accounts {
		accounts[i] = decredmaterial.TabItem{
			Label: page.theme.Body1(page.info.Wallets[*page.selectedWallet].Accounts[i].Name),
		}
	}
	page.accountsTab.SetTabs(accounts)
	if page.accountsTab.Changed() {
		*page.selectedAccount = page.accountsTab.Selected
	}
	page.accountsTab.Separator = false

	bd := func() {
		if page.walletsTab.Changed() {
			*page.selectedWallet = page.walletsTab.Selected
			*page.selectedAccount = 0
			page.accountsTab.Selected = 0
		}
		page.walletsTab.Separator = false
		page.walletsTab.Layout(gtx, body)
	}
	page.Layout(gtx, bd)
}

func (page pageCommon) accountTab(gtx *layout.Context, body layout.Widget) {
	accounts := make([]decredmaterial.TabItem, len(page.info.Wallets[*page.selectedWallet].Accounts))
	for i, account := range page.info.Wallets[*page.selectedWallet].Accounts {
		if account.Name == "imported" {
			continue
		}
		accounts[i] = decredmaterial.TabItem{
			Label: page.theme.Body1(page.info.Wallets[*page.selectedWallet].Accounts[i].Name),
		}
	}
	page.accountsTab.SetTabs(accounts)
	if page.accountsTab.Changed() {
		*page.selectedAccount = page.accountsTab.Selected
	}
	layout.Flex{}.Layout(gtx,
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func() {
				page.theme.H6("Accounts: ").Layout(gtx)
			})
		}),
		layout.Rigid(func() {
			page.accountsTab.Layout(gtx, body)
		}),
	)
}
