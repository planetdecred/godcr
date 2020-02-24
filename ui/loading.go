package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
	"github.com/raedahgroup/godcr-gio/wallet"
)

var (
	loading = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.MultiWalletInfo) {
		layout.Center.Layout(gtx, func() {
			theme.Icon.Logo.Layout(gtx, unit.Dp(100))
		})
	}
	blank = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.MultiWalletInfo) {
		layouts.FillWithColor(gtx, theme.Background)
	}
)
