package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/units"
)

const OverviewID = "overview"

type Overview struct {
	balance      material.Label
	statusTitle  material.Label
	syncStatus   material.Label
	onlineStatus material.Label

	row    layout.Flex
	column layout.Flex
}

// Init initializes all widgets to be used on the page
func (page *Overview) Init(theme *material.Theme) {
	page.row = layout.Flex{Axis: layout.Horizontal}
	page.column = layout.Flex{Axis: layout.Vertical}

	page.balance = theme.H5("154.0928281 DCR")
	page.statusTitle = theme.Caption("Wallet Status")
	page.syncStatus = theme.H6("Syncing...")
	page.onlineStatus = theme.Caption("Online")
}

// Draw adds all the widgets to the stored layout context
func (page *Overview) Draw(gtx *layout.Context, _ event.Event) (evt event.Event) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			layout.UniformInset(units.Padding).Layout(gtx, func() {
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
	in.Layout(gtx, func() {
		page.column.Layout(gtx,
			layout.Rigid(func() {
				page.row.Layout(gtx,
					layout.Rigid(func() {
						page.statusTitle.Layout(gtx)
					}),
					layout.Flexed(units.EntireSpace, func() {
						layout.Align(layout.E).Layout(gtx, func() {
							page.onlineStatus.Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func() {
				page.column.Layout(gtx,
					layout.Rigid(func() {
						page.syncStatus.Layout(gtx)
					}),
					layout.Flexed(units.EntireSpace, func() {
						// todo
						// sync trigger button
					}),
				)
			}),
		)
	})
}
