package ui

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageDebug = "Debug"

type debugPage struct {
	dummyText decredmaterial.Label
}

func (win *Window) DebugPage(common pageCommon) layout.Widget {
	pg := &debugPage{
		dummyText: common.theme.H5("Not yet implemented"),
	}

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

// main settings layout
func (pg *debugPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	return common.Layout(gtx, func(gtx C) D {
		return layout.Center.Layout(gtx, func(gtx C) D {
			return pg.dummyText.Layout(gtx)
		})
	})
}

func (pg *debugPage) handle(common pageCommon) {

}
