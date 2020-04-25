package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, contentClear, contentCreate, navigationCheck *decredmaterial.Icon
}

type pageCommon struct {
	wallet         *wallet.Wallet
	info           *wallet.MultiWalletInfo
	selectedWallet *int
	gtx            *layout.Context
	theme          *decredmaterial.Theme
	icons          pageIcons
	page           *string
	navTab         *decredmaterial.Tabs
	walletsTab     *decredmaterial.Tabs
}

func (win *Window) addPages() {
	icons := pageIcons{
		contentAdd:      mustIcon(decredmaterial.NewIcon(icons.ContentAdd)),
		contentClear:    mustIcon(decredmaterial.NewIcon(icons.ContentClear)),
		contentCreate:   mustIcon(decredmaterial.NewIcon(icons.ContentCreate)),
		navigationCheck: mustIcon(decredmaterial.NewIcon(icons.NavigationCheck)),
	}
	tabs := decredmaterial.NewTabs()
	tabs.SetTabs([]decredmaterial.TabItem{
		{
			Button: win.theme.Button("Overview"),
		},
		{
			Button: win.theme.Button("Wallets"),
		},
		{
			Button: win.theme.Button("Transactions"),
		},
		{
			Button: win.theme.Button("Settings"),
		},
	})

	common := pageCommon{
		wallet:         win.wallet,
		info:           win.walletInfo,
		selectedWallet: &win.selected,
		gtx:            win.gtx,
		theme:          win.theme,
		icons:          icons,
		page:           &win.current,
		navTab:         tabs,
		walletsTab:     decredmaterial.NewTabs(),
		//cancelDialogW:  win.theme.PlainIconButton(icons.contentClear),
	}

	win.pages = make(map[string]layout.Widget)

	win.pages[PageWallet] = WalletPage(common)
	win.pages[PageOverview] = win.OverviewPage
	win.pages[PageTransactions] = win.TransactionsPage
	win.pages[PageReceive] = win.Receive
	win.pages[PageRestore] = win.RestorePage
	win.pages[PageSend] = win.SendPage

}

func (page pageCommon) Layout(gtx *layout.Context, body layout.Widget) {
	navs := []string{PageOverview, PageWallet, PageTransactions, PageOverview}
	toMax(gtx)
	page.navTab.Layout(gtx, body)

	if page.navTab.Changed() {
		*page.page = navs[page.navTab.Selected]
	}
}

func (page pageCommon) LayoutWithWallets(gtx *layout.Context, body layout.Widget) {
	wallets := make([]decredmaterial.TabItem, len(page.info.Wallets))
	for i := range page.info.Wallets {
		wallets[i] = decredmaterial.TabItem{
			Button: page.theme.Button(page.info.Wallets[i].Name),
		}
	}
	page.walletsTab.SetTabs(wallets)
	bd := func() {
		page.walletsTab.Layout(gtx, body)
		if page.walletsTab.Changed() {
			*page.selectedWallet = page.walletsTab.Selected
		}
	}
	page.Layout(gtx, bd)
}
