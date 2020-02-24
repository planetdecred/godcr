package window

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
)

var Loading = func(theme *materialplus.Theme, gtx *layout.Context) {
	layout.Center.Layout(gtx, func() {
		theme.Icons.Loading.Layout(gtx, unit.Dp(100))
	})
}
