package page

import (
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

	copyLog    *widget.Clickable
	copyIcon   *decredmaterial.Image
	backButton decredmaterial.IconButton

	logList layout.List
	fullLog string
}

func (pg *LogPage) ID() string {
	return LogPageID
}

func NewLogPage(l *load.Load) *LogPage {
	pg := &LogPage{
		Load:    l,
		copyLog: new(widget.Clickable),
		logList: layout.List{Axis: layout.Vertical},
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
		clipboard.WriteOp{Text: pg.fullLog}.Add(gtx.Ops)
	}()
}

func (pg *LogPage) watchLogs() {
	logPath := pg.Load.WL.Wallet.LogFile()
	t, _ := tail.TailFile(logPath, tail.Config{Follow: true})
	for line := range t.Lines {
		pg.fullLog += line.Text + "\n"
		pg.RefreshWindow()
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
						return pg.logList.Layout(gtx, 1, func(gtx C, index int) D {
							return pg.Theme.Body1(pg.fullLog).Layout(gtx)
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
