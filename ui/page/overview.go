package page

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/ui/values"
	"github.com/raedahgroup/godcr-gio/ui/widgets"
)

const OverviewID = "overview"

// Overview represents the overview page of the app.
// It is the first page the user sees on launch when a wallet exists.
type Overview struct {
	syncButtonWidget *widget.Button
	progressBar      *widgets.ProgressBar

	balance             material.Label
	statusTitle         material.Label
	syncStatus          material.Label
	onlineStatus        material.Label
	syncButton          material.Button
	progressPercentage  material.Label
	timeLeft            material.Label
	syncSteps           material.Label
	headersFetched      material.Label
	connectedPeersTitle material.Label
	connectedPeers      material.Label

	walletHeaderFetchedTitle   material.Label
	walletSyncingProgressTitle material.Label
	walletSyncDetails          walletSyncDetails

	transactionColumnTitle material.Label
	transactionIcon        material.Label
	transactionAmount      material.Label
	transactionWallet      material.Label
	transactionStatus      material.Label
	transactionDate        material.Label

	column         layout.Flex
	columnMargin   layout.Inset
	row            layout.Flex
	list           layout.List
	walletSyncList layout.List
	syncStatusList layout.List
}

// walletSyncDetails contains sync data for each wallet when a sync
// is in progress.
type walletSyncDetails struct {
	name               material.Label
	status             material.Label
	blockHeaderFetched material.Label
	syncingProgress    material.Label
}

// Init initializes all widgets to be used on the overview page.
func (page *Overview) Init(theme *material.Theme) {
	page.row = layout.Flex{Axis: layout.Horizontal}
	page.column = layout.Flex{Axis: layout.Vertical}
	page.columnMargin = layout.Inset{Top: units.ColumnMargin}
	page.list = layout.List{Axis: layout.Vertical}
	page.walletSyncList = layout.List{Axis: layout.Horizontal}
	page.syncStatusList = layout.List{Axis: layout.Vertical}

	page.balance = theme.H5("154.0928281 DCR")
	page.statusTitle = theme.Caption("Wallet Status")
	page.syncStatus = theme.H6("Syncing...")
	page.onlineStatus = theme.Caption("Online")
	page.syncButtonWidget = new(widget.Button)
	page.syncButton = theme.Button("Cancel")
	page.progressBar = widgets.NewProgressBar()
	page.progressPercentage = theme.Caption("25%")
	page.timeLeft = theme.Caption("6 min left")
	page.syncStatus = theme.H5("Syncing...")
	page.syncSteps = theme.Caption("Step 1/3")
	page.headersFetched = theme.Caption("Fetching block headers. 89%")
	page.connectedPeersTitle = theme.Caption("Connected peers count")
	page.connectedPeers = theme.Caption("16")
	page.walletHeaderFetchedTitle = theme.Caption("Block header fetched")
	page.walletSyncingProgressTitle = theme.Caption("SyncingProgress")
	page.transactionColumnTitle = theme.Caption("Recent Transactions")
	page.transactionIcon = theme.Caption("icon")
	page.transactionAmount = theme.Caption("34.17458878 DCR")
	page.transactionWallet = theme.Caption("Default")
	page.transactionDate = theme.Caption("01/12/2020")
	page.transactionStatus = theme.Caption("Pending")

	page.walletSyncDetails = walletSyncDetails{
		name:               theme.Caption("wallet-1"),
		status:             theme.Caption("Syncing..."),
		blockHeaderFetched: theme.Caption("100 of 164864"),
		syncingProgress:    theme.Caption("320 days behind"),
	}
}

// Draw adds all the widgets to the stored layout context.
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

// content lays out the entire content for overview page.
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

// recentTransactionsColumn lays out the list of recent transactions.
func (page *Overview) recentTransactionsColumn(gtx *layout.Context) {
	var transactionRows []func()
	for i := 0; i < 5; i++ {
		transactionRows = append(transactionRows, func() {
			page.recentTransactionRow(gtx)
		})
	}

	page.columnMargin.Layout(gtx, func() {
		page.column.Layout(gtx,
			layout.Rigid(func() {
				page.transactionColumnTitle.Layout(gtx)
			}),
			layout.Rigid(func() {
				page.list.Layout(gtx, len(transactionRows), func(i int) {
					layout.Inset{Top: units.Padding}.Layout(gtx, transactionRows[i])
				})
			}),
		)
	})
}

// recentTransactionRow lays out a single row of a recent transaction.
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

// syncStatusColumn lays out content for displaying sync status.
func (page *Overview) syncStatusColumn(gtx *layout.Context) {
	uniform := layout.UniformInset(units.Padding)
	syncStatusWidgets := []func(){
		func() {
			page.syncBoxTitleRow(gtx, uniform)
		},
		func() {
			page.syncStatusTextRow(gtx, uniform)
		},
		func() {
			page.progressBarRow(gtx, uniform)
		},
		func() {
			page.progressStatusRow(gtx, uniform)
		},
		func() {
			page.walletSyncRow(gtx, uniform)
		},
	}
	page.columnMargin.Layout(gtx, func() {
		page.syncStatusList.Layout(gtx, len(syncStatusWidgets), func(i int) {
			layout.UniformInset(units.NoPadding).Layout(gtx, syncStatusWidgets[i])
		})
	})
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
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

// syncBoxTitleRow lays out widgets in the title row inside the sync status box.
func (page *Overview) syncBoxTitleRow(gtx *layout.Context, inset layout.Inset) {
	page.endToEndRow(gtx, inset, page.statusTitle, page.onlineStatus)
}

// syncBoxTitleRow lays out sync status text and sync button.
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

// syncBoxTitleRow lays out the progress bar.
func (page *Overview) progressBarRow(gtx *layout.Context, inset layout.Inset) {
	inset.Layout(gtx, func() {
		page.progressBar.Layout(gtx, 25)
	})
}

// syncBoxTitleRow lays out the progress status when the wallet is syncing.
func (page *Overview) progressStatusRow(gtx *layout.Context, inset layout.Inset) {
	page.endToEndRow(gtx, inset, page.progressPercentage, page.timeLeft)
}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func (page *Overview) walletSyncRow(gtx *layout.Context, inset layout.Inset) {
	syncBoxes := []func(){
		func() {
			page.walletSyncBox(gtx, inset, page.walletSyncDetails)
		},
		func() {
			page.walletSyncBox(gtx, inset, page.walletSyncDetails)
		},
		func() {
			page.walletSyncBox(gtx, inset, page.walletSyncDetails)
		},
	}
	page.columnMargin.Layout(gtx, func() {
		page.column.Layout(gtx,
			layout.Rigid(func() {
				page.endToEndRow(gtx, inset, page.syncSteps, page.headersFetched)
			}),
			layout.Rigid(func() {
				page.endToEndRow(gtx, inset, page.connectedPeersTitle, page.connectedPeers)
			}),
			layout.Rigid(func() {
				page.walletSyncList.Layout(gtx, len(syncBoxes), func(i int) {
					layout.Inset{Left: units.ColumnMargin}.Layout(gtx, syncBoxes[i])
				})
			}),
		)
	})
}

// walletSyncBox lays out the wallet syncing details of a single wallet.
func (page *Overview) walletSyncBox(gtx *layout.Context, inset layout.Inset, details walletSyncDetails) {
	page.columnMargin.Layout(gtx, func() {
		layout.Stack{}.Layout(gtx,
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = gtx.Px(units.WalletSyncBoxWidthMin)
				gtx.Constraints.Height.Min = gtx.Px(units.WalletSyncBoxHeightMin)
				fillWalletSyncBox(gtx, values.WalletSyncBoxGray)
			}),
			layout.Stacked(func() {
				uniform := layout.UniformInset(units.SyncBoxPadding)
				uniform.Layout(gtx, func() {
					gtx.Constraints.Width.Min =  gtx.Px(units.WalletSyncBoxContentWidth)
					gtx.Constraints.Width.Max = gtx.Constraints.Width.Min
					page.column.Layout(gtx,
						layout.Rigid(func() {
							page.endToEndRow(gtx, inset, details.name, details.status)
						}),
						layout.Rigid(func() {
							page.endToEndRow(gtx, inset, page.walletHeaderFetchedTitle, details.blockHeaderFetched)
						}),
						layout.Rigid(func() {
							page.endToEndRow(gtx, inset, page.walletSyncingProgressTitle, details.syncingProgress)
						}),
					)
				})
			}),
		)
	})
}

// fillWalletSyncBox adds a coloured rectangle as a background to walletSyncBoxes
func fillWalletSyncBox(gtx *layout.Context, color color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
