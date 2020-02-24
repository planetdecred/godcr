package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

var Loading = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.MultiWalletInfo) {
	layout.Center.Layout(gtx, func() {
		theme.Icons.Loading.Layout(gtx, unit.Dp(100))
	})
}
