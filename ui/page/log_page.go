package page

import (
	"fmt"
	"os"
	"runtime"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/nxadm/tail"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	LogPageID = "Log"
	LogOffset = 24000
)

type LogPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	tail *tail.Tail

	copyLog    *decredmaterial.Clickable
	copyIcon   *decredmaterial.Image
	backButton decredmaterial.IconButton

	logList *widget.List
	fullLog string
}

func NewLogPage(l *load.Load) *LogPage {
	pg := &LogPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(LogPageID),
		logList: &widget.List{
			List: layout.List{
				Axis:        layout.Vertical,
				ScrollToEnd: true,
			},
		},
		copyLog: l.Theme.NewClickable(true),
	}

	pg.copyIcon = pg.Theme.Icons.CopyIcon

	pg.backButton, _ = components.SubpageHeaderButtons(l)
	pg.watchLogs()
	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *LogPage) OnNavigatedTo() {
	pg.watchLogs()
}

func (pg *LogPage) copyLogEntries(gtx C) {
	go func() {
		clipboard.WriteOp{Text: pg.fullLog}.Add(gtx.Ops)
	}()
}

func (pg *LogPage) watchLogs() {
	go func() {
		logPath := pg.Load.WL.Wallet.LogFile()

		fi, err := os.Stat(logPath)
		if err != nil {
			pg.fullLog = fmt.Sprintf("unable to open log file: %v", err)
			return
		}

		size := fi.Size()

		var offset int64
		if size > LogOffset*2 {
			offset = size - LogOffset
		}

		pollLogs := runtime.GOOS == "windows"
		t, err := tail.TailFile(logPath, tail.Config{Follow: true, Poll: pollLogs, Location: &tail.SeekInfo{Offset: offset}})
		if err != nil {
			pg.fullLog = fmt.Sprintf("unable to tail log file: %v", err)
			return
		}
		pg.tail = t

		if offset > 0 {
			// skip the first line because it might be truncated.
			<-t.Lines
		}
		for line := range t.Lines {
			pg.fullLog += line.Text + "\n"
			pg.ParentWindow().Reload()
		}
	}()
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *LogPage) Layout(gtx C) D {
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *LogPage) layoutDesktop(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrWalletLog),
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			ExtraItem: pg.copyLog,
			Extra: func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return pg.copyLog.Layout(gtx, func(gtx C) D {
						return pg.copyIcon.Layout24dp(gtx)
					})

				})
			},
			HandleExtra: func() {
				pg.copyLogEntries(gtx)
				pg.Toast.Notify(values.String(values.StrCopied))
			},
			Body: func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
				return pg.Theme.List(pg.logList).Layout(gtx, 1, func(gtx C, index int) D {
					return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
						return pg.Theme.Card().Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
								return pg.Theme.Body1(pg.fullLog).Layout(gtx)
							})
						})
					})
				})
			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	return components.UniformPadding(gtx, container)
}

func (pg *LogPage) layoutMobile(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrWalletLog),
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			ExtraItem: pg.copyLog,
			Extra: func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return pg.copyLog.Layout(gtx, func(gtx C) D {
						return pg.copyIcon.Layout24dp(gtx)
					})

				})
			},
			HandleExtra: func() {
				pg.copyLogEntries(gtx)
				pg.Toast.Notify(values.String(values.StrCopied))
			},
			Body: func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
				return pg.Theme.List(pg.logList).Layout(gtx, 1, func(gtx C, index int) D {
					return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
						return pg.Theme.Card().Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
								return pg.Theme.Body1(pg.fullLog).Layout(gtx)
							})
						})
					})
				})
			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	return components.UniformMobile(gtx, false, false, container)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *LogPage) HandleUserInteractions() {}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *LogPage) OnNavigatedFrom() {
	if pg.tail != nil {
		pg.tail.Stop()
	}
}
