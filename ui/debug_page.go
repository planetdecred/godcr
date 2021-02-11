package ui

import (
	"sync"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageDebug = "Debug"

type debugPage struct {
	theme *decredmaterial.Theme

	logTitle   decredmaterial.Label
	debugList  layout.List
	logLabels  []decredmaterial.Label
	labelsLock sync.Mutex
}

func (win *Window) DebugPage(common pageCommon, internalLog chan string) layout.Widget {
	pg := &debugPage{
		theme:    common.theme,
		logTitle: common.theme.H5("Session log entries"),
		debugList: layout.List{
			Axis: layout.Vertical,
		},
		logLabels: make([]decredmaterial.Label, 0),
	}

	go pg.watchLogs(internalLog)

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *debugPage) newLogWidget(str string) decredmaterial.Label {
	w := pg.theme.Body1(str)
	return w
}

func (pg *debugPage) watchLogs(internalLog chan string) {
	for b := range internalLog {
		pg.labelsLock.Lock()
		pg.logLabels = append(pg.logLabels, pg.newLogWidget(b))
		pg.labelsLock.Unlock()
	}
}

// main settings layout
func (pg *debugPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	return common.Layout(gtx, func(gtx C) D {
		return layout.Center.Layout(gtx, func(gtx C) D {
			return pg.debugList.Layout(gtx, len(pg.logLabels), func(gtx layout.Context, i int) layout.Dimensions {
				pg.labelsLock.Lock()
				defer pg.labelsLock.Unlock()
				return pg.logLabels[i].Layout(gtx)
			})
		})
	})
}

func (pg *debugPage) handle(common pageCommon) {

}
