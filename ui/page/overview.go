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

	balance                material.Label
	statusTitle            material.Label
	syncStatus             material.Label
	onlineStatus           material.Label
	syncButton             material.Button
	progressPercentage     material.Label
	timeLeft               material.Label
	syncSteps 			   material.Label

	transactionColmunTitle material.Label
	transactionIcon        material.Label
	transactionAmount      material.Label
	transactionWallet      material.Label
	transactionStatus      material.Label
	transactionDate        material.Label

	column       layout.Flex
	columnMargin layout.Inset
	row          layout.Flex
	list         layout.List
}

// Init initializes all widgets to be used on the page
func (page *Overview) Init(theme *material.Theme) {
	page.row = layout.Flex{Axis: layout.Horizontal}
	page.column = layout.Flex{Axis: layout.Vertical}
	page.columnMargin = layout.Inset{Top: units.ColumnMargin}
	page.list = layout.List{Axis: layout.Vertical}

	page.balance = theme.H5("154.0928281 DCR")
	page.statusTitle = theme.Caption("Wallet Status")
	page.syncStatus = theme.H6("Syncing...")
	page.onlineStatus = theme.Caption("Online")
	page.syncButtonWidget = new(widget.Button)
	page.syncButton = theme.Button("Cancel")
	page.progressBar = widgets.NewProgressBar()
	page.progressPercentage = theme.Caption("25%")
	page.timeLeft = theme.Caption("6 min left")
	page.transactionColmunTitle = theme.Caption("Recent Transactions")
	page.transactionIcon = theme.Caption("icon")
	page.transactionAmount = theme.Caption("34.17458878 DCR")
	page.transactionWallet = theme.Caption("Default")
	page.transactionDate = theme.Caption("01/12/2020")
	page.transactionStatus = theme.Caption("Pending")
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
			page.recentTransactionsColumn(gtx)
		}),
		layout.Rigid(func() {
			page.syncStatusColumn(gtx)
		}),
	)
}

// recentTransactionsColumn lays out the list of recent transactions
func (page *Overview) recentTransactionsColumn(gtx *layout.Context) {
	var transactionRows []func()
	for i:=0; i<5; i++ {
		transactionRows = append(transactionRows, func() {
			page.recentTransactionRow(gtx)
		})
	}

	page.columnMargin.Layout(gtx, func() {
		page.column.Layout(gtx,
			layout.Rigid(func() {
				page.transactionColmunTitle.Layout(gtx)
			}),
			layout.Rigid(func() {
				page.list.Layout(gtx, len(transactionRows), func(i int) {
					layout.Inset{Top: units.Padding}.Layout(gtx, transactionRows[i])
				})
			}),
		)
	})
}

func (page *Overview) recentTransactionRow(gtx *layout.Context) {
	margin := layout.UniformInset(units.TransactionsRowMargin)
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			margin.Layout(gtx, func() {
				page.transactionIcon.Layout(gtx)
			})
		}),
		layout.Rigid(func() {
			margin.Layout(gtx, func() {
				page.transactionAmount.Layout(gtx)
			})
		}),
		layout.Rigid(func() {
			margin.Layout(gtx, func() {
				page.transactionWallet.Layout(gtx)
			})
		}),
		layout.Rigid(func() {
			margin.Layout(gtx, func() {
				page.transactionDate.Layout(gtx)
			})
		}),
		layout.Flexed(1, func() {
			layout.Align(layout.E).Layout(gtx, func() {
				margin.Layout(gtx, func() {
					page.transactionStatus.Layout(gtx)
				})
			})
		}),
	)
}

// syncStatusColumn lays out content for displaying sync status
func (page *Overview) syncStatusColumn(gtx *layout.Context) {
	uniform := layout.UniformInset(units.Padding)
	page.columnMargin.Layout(gtx, func() {
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

// syncBoxTitleRow lays out widgets in the title row inside the sync status box
func (page *Overview) syncBoxTitleRow(gtx *layout.Context, inset layout.Inset) {
	page.endToEndRow(gtx, inset, page.statusTitle, page.onlineStatus)
}

// syncBoxTitleRow lays out sync status text and sync button
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

// syncBoxTitleRow lays out the progress bar
func (page *Overview) progressBarRow(gtx *layout.Context, inset layout.Inset) {
	inset.Layout(gtx, func() {
		page.progressBar.Layout(gtx, 25)
	})
}

// syncBoxTitleRow lays out the progress bar
func (page *Overview) progressStatusRow(gtx *layout.Context, inset layout.Inset) {
	page.endToEndRow(gtx, inset, page.progressPercentage, page.timeLeft)
}

func (page *Overview) syncDetailsColumn(gtx *layout.Context, inset layout.Inset) {
	uniform := layout.UniformInset(units.Padding)
	page.columnMargin.Layout(gtx, func() {
		page.column.Layout(gtx,
			layout.Rigid(func() {
				page.syncBoxTitleRow(gtx, uniform)
			}),
		)
	})
}

func (page *Overview) endToEndRow(gtx *layout.Context, inset layout.Inset, leftLabel, rightLabel material.Label) {
	inset.Layout(gtx, func() {
		page.row.Layout(gtx,
			layout.Rigid(func() {
				leftLabel.Layout(gtx)
			}),
			layout.Flexed(values.EntireSpace, func() {
				layout.Align(layout.E).Layout(gtx, func() {
					rightLabel.Layout(gtx)
				})
			}),
		)
	})
}
