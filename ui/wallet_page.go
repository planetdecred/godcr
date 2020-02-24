package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
)

func (win *Window) WalletsPage() layout.Widget {
	log.Debug("On Wallets")
	return func() {
		win.TabbedWallets(
			func() {
				layouts.FillWithColor(win.gtx, win.theme.Background)
				win.theme.Label(unit.Dp(100), "Selected")
			},
			func() {
				win.theme.Label(unit.Dp(100), "Body")
			},
			func(i int) {
				win.theme.Label(unit.Dp(100), "Item")
			},
		)
	}
}
