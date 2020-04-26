package ui

import (
	"image"

	"gioui.org/layout"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, contentClear, contentCreate, navigationCheck,
	contentSend, contentAddBox, contentRemove, toggleRadioButtonUnchecked,
	actionCheckCircle *decredmaterial.Icon
	overviewIcon, walletIcon image.Image
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
		contentAdd:                 mustIcon(decredmaterial.NewIcon(icons.ContentAdd)),
		contentClear:               mustIcon(decredmaterial.NewIcon(icons.ContentClear)),
		contentCreate:              mustIcon(decredmaterial.NewIcon(icons.ContentCreate)),
		navigationCheck:            mustIcon(decredmaterial.NewIcon(icons.NavigationCheck)),
		overviewIcon:               mustDcrIcon(decredmaterial.NewDcrIcon(decredmaterial.OverviewIcon)),
		walletIcon:                 mustDcrIcon(decredmaterial.NewDcrIcon(decredmaterial.WalletIcon)),
		contentSend:                mustIcon(decredmaterial.NewIcon(icons.ContentSend)),
		contentAddBox:              mustIcon(decredmaterial.NewIcon(icons.ContentAddBox)),
		contentRemove:              mustIcon(decredmaterial.NewIcon(icons.ContentRemove)),
		toggleRadioButtonUnchecked: mustIcon(decredmaterial.NewIcon(icons.ToggleRadioButtonUnchecked)),
		actionCheckCircle:          mustIcon(decredmaterial.NewIcon(icons.ActionCheckCircle)),
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
			Label: win.theme.Body1("Settings"),
			Icon:  icons.overviewIcon,
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

	win.pages[PageWallet] = win.WalletPage(common)
	win.pages[PageOverview] = win.OverviewPage
	win.pages[PageTransactions] = win.TransactionsPage(common)
	win.pages[PageReceive] = win.Receive
	win.pages[PageRestore] = win.RestorePage
	win.pages[PageSend] = win.SendPage
	win.pages[PageSignMessage] = win.SignMessagePage
	win.pages[PageTransactionDetails] = TransactionPage(common, &win.walletTransaction)

}

func (page pageCommon) Layout(gtx *layout.Context, body layout.Widget) {
	navs := []string{PageOverview, PageWallet, PageTransactions, PageOverview}
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
	bd := func() {
		if page.walletsTab.Changed() {
			*page.selectedWallet = page.walletsTab.Selected
		}
		page.walletsTab.Separator = false
		page.walletsTab.Layout(gtx, body)
	}
	page.Layout(gtx, bd)
}
