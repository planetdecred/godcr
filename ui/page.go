package ui

import (
	"image"

	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/io/key"
	"gioui.org/layout"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
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
	// gtx             layout.Context
	theme         *decredmaterial.Theme
	icons         pageIcons
	page          *string
	navTab        *decredmaterial.Tabs
	walletsTab    *decredmaterial.Tabs
	accountsTab   *decredmaterial.Tabs
	errorChannels map[string]chan error
	keyEvents     chan *key.Event
	states        *states
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
		overviewIcon:               decredIcons["overview"],
		walletIcon:                 decredIcons["wallet"],
		receiveIcon:                decredIcons["receive"],
		transactionIcon:            decredIcons["transaction"],
		sendIcon:                   decredIcons["send"],
		syncingIcon:                decredIcons["syncing"],
	}

	tabs := decredmaterial.NewTabs(win.theme)
	tabs.SetTabs([]decredmaterial.TabItem{
		{
			Label: win.theme.Body1("Overview"),
			Icon:  ic.overviewIcon,
		},
		{
			Label: win.theme.Body1("Wallets"),
			Icon:  ic.walletIcon,
		},
		{
			Label: win.theme.Body1("Send"),
			Icon:  ic.sendIcon,
		},
		{
			Label: win.theme.Body1("Receive"),
			Icon:  ic.receiveIcon,
		},
		{
			Label: win.theme.Body1("Transactions"),
			Icon:  ic.transactionIcon,
		},
	})

	accountsTab := decredmaterial.NewTabs(win.theme)
	accountsTab.Position = decredmaterial.Top
	accountsTab.Separator = false
	common := pageCommon{
		wallet:          win.wallet,
		info:            win.walletInfo,
		selectedWallet:  &win.selected,
		selectedAccount: &win.selectedAccount,
		theme:           win.theme,
		icons:           ic,
		page:            &win.current,
		navTab:          tabs,
		walletsTab:      decredmaterial.NewTabs(win.theme),
		accountsTab:     accountsTab,
		errorChannels: map[string]chan error{
			PageSignMessage:    make(chan error),
			PageCreateRestore:  make(chan error),
			PageWallet:         make(chan error),
			PageWalletAccounts: make(chan error),
		},
		keyEvents: win.keyEvents,
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
	for range page.navTab.ChangeEvent(gtx) {
		*page.page = navs[page.navTab.Selected]
	}

	//return layout.Inset{Right: values.MarginPadding10, Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
	//	return body(gtx)
	//})

	page.navTab.Separator = true
	return page.navTab.Layout(gtx, func(gtx C) D {
		p := values.MarginPadding10
		return layout.Inset{Top: p, Left: p, Right: p}.Layout(gtx, func(gtx C) D {
			return body(gtx)
		})
	})
}

func (page pageCommon) LayoutWithWallets(gtx layout.Context, body layout.Widget) layout.Dimensions {
	wallets := make([]decredmaterial.TabItem, len(page.info.Wallets))
	for i := range page.info.Wallets {
		wallets[i] = decredmaterial.TabItem{
			Label: page.theme.Body1(page.info.Wallets[i].Name),
		}
	}
	page.walletsTab.SetTabs(wallets)
	page.walletsTab.Position = decredmaterial.Top
	for range page.accountsTab.ChangeEvent(gtx) {
		*page.selectedAccount = page.accountsTab.Selected
	}

	accounts := make([]decredmaterial.TabItem, len(page.info.Wallets[*page.selectedWallet].Accounts))
	for i, acct := range page.info.Wallets[*page.selectedWallet].Accounts {
		if acct.Name == "imported" {
			continue
		}
		accounts[i] = decredmaterial.TabItem{
			Label: page.theme.Body1(page.info.Wallets[*page.selectedWallet].Accounts[i].Name),
		}
	}
	page.accountsTab.SetTabs(accounts)
	for range page.accountsTab.ChangeEvent(gtx) {
		*page.selectedAccount = page.accountsTab.Selected
	}
	page.accountsTab.Separator = false

	bd := func(gtx C) D {
		for range page.walletsTab.ChangeEvent(gtx) {
			*page.selectedWallet = page.walletsTab.Selected
			*page.selectedAccount = 0
			page.accountsTab.Selected = 0
		}
		if *page.selectedWallet == 0 {
			page.walletsTab.Selected = *page.selectedWallet
		}
		page.walletsTab.Separator = false
		return page.walletsTab.Layout(gtx, body)
	}
	return page.Layout(gtx, bd)
}

func (page pageCommon) LayoutWithAccounts(gtx layout.Context, body layout.Widget) layout.Dimensions {
	accounts := make([]decredmaterial.TabItem, len(page.info.Wallets[*page.selectedWallet].Accounts))
	for i, account := range page.info.Wallets[*page.selectedWallet].Accounts {
		if account.Name == "imported" {
			continue
		}
		accounts[i] = decredmaterial.TabItem{
			Label: page.theme.Body1(page.info.Wallets[*page.selectedWallet].Accounts[i].Name),
		}
	}

	page.accountsTab.SetTitle(page.theme.Label(values.TextSize18, "Accounts:"))

	page.accountsTab.SetTabs(accounts)
	for range page.accountsTab.ChangeEvent(gtx) {
		*page.selectedAccount = page.accountsTab.Selected
	}

	return page.LayoutWithWallets(gtx, func(gtx C) D {
		return page.accountsTab.Layout(gtx, body)
	})
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
