package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
	"github.com/raedahgroup/godcr-gio/wallet"
)

var (
	walletTabItem = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.InfoShort) {
		layouts.FillWithColor(gtx, theme.Background)
		theme.Label(theme.TextSize, info.Balance.String())
	}
	walletTabs = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.MultiWalletInfo) {

	}
)

type tabBody func(*layout.Context, *materialplus.Theme, *wallet.InfoShort)

func (win *Window) tabbedBody(w WalletPage) {
	log.Debugf("Wallets %d", len(win.walletInfo.Wallets))
	log.Debugf("Buttons %d", len(win.buttons.tabs))
	if len(win.walletInfo.Wallets) == 0 {
		loading(win.gtx, win.theme, nil)
		return
	}
	//layouts.FillWithColor(win.gtx, win.theme.Background)
	layouts.Tabs{
		Selected: func() {
			walletTabItem(win.gtx, win.theme, &win.walletInfo.Wallets[win.selected])
		},
		Item: func(i int) {
			walletTabItem(win.gtx, win.theme, &win.walletInfo.Wallets[i])
		},
		Body: func() {
			w(win.gtx, win.theme, &win.walletInfo.Wallets[win.selected])
		},
		List: win.tabsList,
		Flex: layout.Flex{
			Axis: layout.Horizontal,
		},
		Size:       .3,
		ButtonSize: .2,
	}.Layout(win.gtx, &win.selected, win.buttons.tabs)
}
