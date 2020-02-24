package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

var landing = func(gtx *layout.Context, theme *materialplus.Theme, info *wallet.InfoShort) {
	theme.Button("CreateWallet").Layout(gtx, new(widget.Button))
}
