package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, contentClear, contentCreate *decredmaterial.Icon
}

type pageCommon struct {
	wallet         *wallet.Wallet
	info           *wallet.MultiWalletInfo
	selectedWallet *int
	gtx            *layout.Context
	theme          *decredmaterial.Theme
	icons          pageIcons
	page           *string
}

func (win *Window) addPages() {
	icons := pageIcons{
		contentAdd:    mustIcon(decredmaterial.NewIcon(icons.ContentAdd)),
		contentClear:  mustIcon(decredmaterial.NewIcon(icons.ContentClear)),
		contentCreate: mustIcon(decredmaterial.NewIcon(icons.ContentCreate)),
	}
	common := pageCommon{
		wallet:         win.wallet,
		info:           win.walletInfo,
		selectedWallet: &win.selected,
		gtx:            win.gtx,
		theme:          win.theme,
		icons:          icons,
	}

	win.pages = make(map[string]layout.Widget)

	win.pages[PageWallet] = WalletPage(common)
	win.pages[PageOverview] = win.OverviewPage
	win.pages[PageTransactions] = win.TransactionsPage
	win.pages[PageReceive] = win.Receive
	win.pages[PageRestore] = win.RestorePage
	win.pages[PageSend] = win.SendPage
}
