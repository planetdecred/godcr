package page

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/ui/values"
	"github.com/raedahgroup/godcr-gio/ui/widgets"
)

const OverviewID = "overview"

type Overview struct {
	syncButtonWidget *widget.Button
	progressBar      *widgets.ProgressBar

	balance      material.Label
	statusTitle  material.Label
	syncStatus   material.Label
	onlineStatus material.Label
	syncButton   material.Button
	progressPercentage material.Label
	timeLeft material.Label

	column layout.Flex
	row    layout.Flex
}

// Init initializes all widgets to be used on the page
func (page *Overview) Init(theme *material.Theme) {
	page.row = layout.Flex{Axis: layout.Horizontal}
	page.column = layout.Flex{Axis: layout.Vertical}

	page.balance = theme.H5("154.0928281 DCR")
	page.statusTitle = theme.Caption("Wallet Status")
	page.syncStatus = theme.H6("Syncing...")
	page.onlineStatus = theme.Caption("Online")
	page.syncButtonWidget = new(widget.Button)
	page.syncButton = theme.Button("Cancel")
	page.progressBar = widgets.NewProgressBar()
	page.progressPercentage = theme.Caption("25%")
	page.timeLeft = theme.Caption("6 min left")
}

// Draw adds all the widgets to the stored layout context
func (page *Overview) Draw(gtx *layout.Context, _ event.Event) (evt event.Event) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			layout.UniformInset(units.ContainerPadding).Layout(gtx, func() {
				page.content(gtx)
			})
		}),
	)
	return
}

// content lays out the entire content for overview page
func (page *Overview) content(gtx *layout.Context) {
	page.column.Layout(gtx,
		layout.Rigid(func() {
			page.balance.Layout(gtx)
		}),
		layout.Rigid(func() {
			page.syncStatusColumn(gtx)
		}),
	)
}

// syncStatusColumn lays out content for displaying sync status
func (page *Overview) syncStatusColumn(gtx *layout.Context) {
	in := layout.Inset{Top: units.FlexInset}
	uniform := layout.UniformInset(units.Padding)
	in.Layout(gtx, func() {
		page.column.Layout(gtx,
			layout.Rigid(func() {
				page.syncBoxTitleRow(gtx, uniform)
			}),
			layout.Rigid(func() {
				page.syncStatusTextRow(gtx, uniform)
			}),
			layout.Rigid(func() {
				page.progressBarRow(gtx, uniform)
			}),
			layout.Rigid(func() {
				page.progressStatusRow(gtx, uniform)
			}),
		)
	})
}

func (page *Overview) syncBoxTitleRow(gtx *layout.Context, inset layout.Inset) {
	inset.Layout(gtx, func() {
		page.row.Layout(gtx,
			layout.Rigid(func() {
				page.statusTitle.Layout(gtx)
			}),
			layout.Flexed(values.EntireSpace, func() {
				layout.Align(layout.E).Layout(gtx, func() {
					page.onlineStatus.Layout(gtx)
				})
			}),
		)
	})
}

func (page *Overview) syncStatusTextRow(gtx *layout.Context, inset layout.Inset) {
	inset.Layout(gtx, func() {
		page.row.Layout(gtx,
			layout.Rigid(func() {
				page.syncStatus.Layout(gtx)
			}),
			layout.Flexed(values.EntireSpace, func() {
				layout.Align(layout.E).Layout(gtx, func() {
					page.syncButton.Layout(gtx, page.syncButtonWidget)
				})
			}),
		)
	})
}

func (page *Overview) progressBarRow(gtx *layout.Context, inset layout.Inset) {
	inset.Layout(gtx, func() {
		page.progressBar.Layout(gtx, 25)
	})
}

func (page *Overview) progressStatusRow(gtx *layout.Context, inset layout.Inset) {
	inset.Layout(gtx, func() {
		page.row.Layout(gtx,
			layout.Rigid(func() {
				page.progressPercentage.Layout(gtx)
			}),
			layout.Flexed(values.EntireSpace, func() {
				layout.Align(layout.E).Layout(gtx, func() {
					page.timeLeft.Layout(gtx)
				})
			}),
		)
	})
}
