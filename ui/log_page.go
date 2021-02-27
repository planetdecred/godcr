package ui

import (
	"sync"

	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageLog = "Log"

type logPage struct {
	theme *decredmaterial.Theme

	entriesList layout.List
	logEntries  []decredmaterial.Label
	entriesLock sync.Mutex
}

func (win *Window) LogPage(common pageCommon, internalLog chan string) layout.Widget {
	pg := &logPage{
		theme: common.theme,
		entriesList: layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: true,
		},
		logEntries: make([]decredmaterial.Label, 0, 20),
	}

	go pg.watchLogs(internalLog)

	return func(gtx C) D {
		//pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *logPage) watchLogs(internalLog chan string) {
	for l := range internalLog {
		entry := l[:len(l)-1]
		pg.entriesLock.Lock()
		pg.logEntries = append(pg.logEntries, pg.theme.Body1(entry))
		pg.entriesLock.Unlock()
	}
}

func (pg *logPage) Layout(gtx C, common pageCommon) D {
	container := func(gtx C) D {
		page := SubPage{
			title: "Wallet log",
			back: func() {
				*common.page = PageDebug
			},
			body: func(gtx C) D {
				background := common.theme.Color.Surface
				card := common.theme.Card()
				card.Color = background
				return card.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						return pg.entriesList.Layout(gtx, len(pg.logEntries), func(gtx C, i int) D {
							pg.entriesLock.Lock()
							defer pg.entriesLock.Unlock()
							return pg.logEntries[i].Layout(gtx)
						})
					})

				})
			},
		}
		return common.SubPageLayout(gtx, page)
	}
	return common.Layout(gtx, container)
}
