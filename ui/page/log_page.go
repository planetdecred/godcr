package page

import (
	"sync"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/nxadm/tail"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const LogPageID = "Log"

type LogPage struct {
	*load.Load

	internalLog chan string
	copyLog     *widget.Clickable
	copyIcon    *decredmaterial.Image
	backButton  decredmaterial.IconButton

	entriesList layout.List
	fullLog     string
	logEntries  []decredmaterial.Label
	entriesLock sync.Mutex
}

func (pg *LogPage) ID() string {
	return LogPageID
}

func NewLogPage(l *load.Load) *LogPage {
	pg := &LogPage{
		Load: l,
		entriesList: layout.List{
			Axis: layout.Vertical,
		},
		copyLog:    new(widget.Clickable),
		logEntries: make([]decredmaterial.Label, 0, 20),
	}

	pg.copyIcon = pg.Icons.CopyIcon

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	go pg.watchLogs()

	return pg
}

func (pg *LogPage) OnResume() {

}

func (pg *LogPage) copyLogEntries(gtx C) {
	go func() {
		pg.entriesLock.Lock()
		defer pg.entriesLock.Unlock()
		clipboard.WriteOp{Text: pg.fullLog}.Add(gtx.Ops)
	}()
}

func (pg *LogPage) watchLogs() {
	//TODO
	//add function to get log directory
	logPath := pg.Load.WL.MultiWallet.LogDir()
	t, _ := tail.TailFile(logPath, tail.Config{Follow: true})
	for line := range t.Lines {
		logRow := line.Text
		entry := logRow[:len(logRow)-1]
		pg.entriesLock.Lock()
		pg.fullLog += entry
		pg.logEntries = append(pg.logEntries, pg.Theme.Body1(logRow))
		pg.entriesLock.Unlock()
	}
}

func (pg *LogPage) Layout(gtx C) D {
	container := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Wallet log",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			ExtraItem: pg.copyLog,
			Extra: func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return decredmaterial.Clickable(gtx, pg.copyLog, func(gtx C) D {
						return pg.copyIcon.Layout24dp(gtx)
					})

				})
			},
			HandleExtra: func() {
				pg.copyLogEntries(gtx)
				pg.Toast.Notify("Copied")
			},
			Body: func(gtx C) D {
				background := pg.Theme.Color.Surface
				card := pg.Theme.Card()
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
		return sp.Layout(gtx)
	}
	return components.UniformPadding(gtx, container)
}

func (pg *LogPage) Handle()  {}
func (pg *LogPage) OnClose() {}
