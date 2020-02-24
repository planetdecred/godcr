package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
	"github.com/raedahgroup/godcr-gio/wallet"
)

var (
	walletTabItem = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.MultiWalletInfo) {
		layout.Center.Layout(gtx, func() {
			theme.Label(theme.TextSize, info.TotalBalance.String())
		})
	}
	walletTabs = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.MultiWalletInfo) {

	}
)

func (win *Window) tabbedBody(gtx *layout.Context, theme *materialplus.Theme, info *wallet.MultiWalletInfo, w layout.Widget) layout.Widget {
	return func() {
		layouts.Tabs{
			Selected: func() {},
			Item: func(i int) {

			},
			Body: w,
			List: &layout.List{Axis: layout.Vertical},
			Flex: layout.Flex{
				Axis: layout.Horizontal,
			},
			Size:       .3,
			ButtonSize: .2,
		}.Layout(gtx, &win.selected, win.buttons.tabs)
	}
}
