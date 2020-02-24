package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

var walletInfo = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.InfoShort) {
	theme.Label(theme.TextSize, info.Name).Layout(gtx)
}
