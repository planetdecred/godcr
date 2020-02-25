package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
)

// TabbedWallets layouts a layout.Tabs
func (win *Window) TabbedPage(body layout.Widget) {
	items := make([]materialplus.TabItem, win.walletInfo.LoadedWallets)
	for i := 0; i < win.walletInfo.LoadedWallets; i++ {
		items[i] = materialplus.TabItem{
			Button: win.theme.Button(win.walletInfo.Wallets[i].Name),
		}
	}
	win.tabs.Layout(win.gtx, &win.selected, win.inputs.tabs, items, body)
}

// Header lays out the window header
func (win *Window) Header() {
	win.theme.Label(unit.Dp(50), "GoDcr").Layout(win.gtx)
}
