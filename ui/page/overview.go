package page

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/decred/dcrd/dcrutil"

	"github.com/raedahgroup/godcr-gio/ui/helper"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/ui/values"
	"github.com/raedahgroup/godcr-gio/ui/widgets"
	"github.com/raedahgroup/godcr-gio/wallet"
)

const OverviewID = "overview"

// Overview represents the overview page of the app.
// It is the first page the user sees on launch when a wallet exists.
type Overview struct {
	theme *materialplus.Theme

	syncButtonWidget    *widget.Button
	progressBar         *materialplus.ProgressBar
	balanceTitle        material.Label
	mainBalance         material.Label
	subBalance          material.Label
	statusTitle         material.Label
	latestBlockTitle    material.Label
	latestBlock         material.Label
	latestBlockTime 	material.Label
	syncStatus          material.Label
	onlineStatus        material.Label
	syncButton          material.Button
	syncButtonCard      widgets.Card
	progressPercentage  material.Label
	timeLeft            material.Label
	syncSteps           material.Label
	headersFetched      material.Label
	connectedPeersTitle material.Label
	connectedPeers      material.Label
	test                material.Label

	walletHeaderFetchedTitle   material.Label
	walletSyncingProgressTitle material.Label
	walletSyncDetails          walletSyncDetails
	walletSyncCard             widgets.Card
	walletSyncBoxes            []func()

	transactionColumnTitle material.Label
	transactionIcon        material.Label
	transactionMainAmount  material.Label
	transactionSubAmount   material.Label
	transactionWallet      material.Label
	transactionStatus      material.Label
	transactionDate        material.Label

	column         layout.Flex
	columnMargin   layout.Inset
	row            layout.Flex
	list           layout.List
	listContainer  layout.List
	walletSyncList layout.List

	transactionAmount string
	balance           string
	progress          float64
	states            map[string]interface{}
	wallet            *wallet.Wallet
	walletInfo        *wallet.MultiWalletInfo
	walletSyncStatus  *wallet.SyncStatus
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
func (page *Overview) Init(theme *materialplus.Theme, w *wallet.Wallet, states map[string]interface{}) {
	page.theme = theme
	page.states = states
	page.wallet = w
	page.row = layout.Flex{Axis: layout.Horizontal}
	page.column = layout.Flex{Axis: layout.Vertical}
	page.columnMargin = layout.Inset{Top: units.ColumnMargin}
	page.list = layout.List{Axis: layout.Vertical}
	page.walletSyncList = layout.List{Axis: layout.Horizontal}
	page.listContainer = layout.List{Axis: layout.Vertical}

	page.balanceTitle = theme.Caption("Current Total Balance")
	page.balance = "315.08193725 DCR"
	page.mainBalance = theme.H4("")
	page.subBalance = theme.H6("")
	page.statusTitle = theme.Caption("Wallet Status")
	page.syncStatus = theme.H6("Syncing...")
	page.onlineStatus = theme.Caption("Online")
	page.syncButtonWidget = new(widget.Button)
	page.syncButton = theme.Button("Cancel")
	page.progressBar = theme.ProgressBar()
	page.progressPercentage = theme.Caption("0%")
	page.timeLeft = theme.Caption("0s left")
	page.syncStatus = theme.H5("Syncing...")
	page.syncSteps = theme.Caption("Step 0/0")
	page.headersFetched = theme.Caption("Fetching block headers. 0%")
	page.connectedPeersTitle = theme.Caption("Connected peers count")
	page.connectedPeers = theme.Caption("0")
	page.walletHeaderFetchedTitle = theme.Caption("Block header fetched")
	page.walletSyncingProgressTitle = theme.Caption("SyncingProgress")
	page.walletSyncCard = widgets.NewCard()
	page.transactionColumnTitle = theme.Caption("Recent Transactions")
	page.transactionIcon = theme.Caption("icon")
	page.transactionAmount = "34.17458878 DCR"
	page.transactionMainAmount = theme.Label(units.TransactionBalanceMain, "")
	page.transactionSubAmount = theme.Label(units.TransactionBalanceSub, "")
	page.transactionWallet = theme.Caption("Default")
	page.transactionDate = theme.Caption("11 Jan 2020, 13:24")
	page.transactionStatus = theme.Caption("Pending")
	page.test = theme.Caption("t")
	page.syncButtonCard = widgets.NewCard()
	page.latestBlockTitle = theme.Caption("Latest Block")
	page.latestBlock = theme.Caption("")

	page.walletSyncDetails = walletSyncDetails{
		name:               theme.Caption("wallet-1"),
		status:             theme.Caption("Syncing..."),
		blockHeaderFetched: theme.Caption("100 of 164864"),
		syncingProgress:    theme.Caption("320 days behind"),
	}
}

func (page *Overview) syncDetail(name, status, headersFetched, progress string) walletSyncDetails {
	theme := page.theme
	return walletSyncDetails{
		name:               theme.Caption(name),
		status:             theme.Caption(status),
		blockHeaderFetched: theme.Caption(headersFetched),
		syncingProgress:    theme.Caption(progress),
	}
}

// Draw adds all the widgets to the stored layout context.
func (page *Overview) Draw(gtx *layout.Context) interface{} {
	page.walletInfo = page.states[StateWalletInfo].(*wallet.MultiWalletInfo)
	page.update(gtx)

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			container := layout.Inset{Left: units.ContainerPadding, Right: units.ContainerPadding}
			container.Layout(gtx, func() {
				page.content(gtx)
			})
		}),
	)
	return nil
}

func (page *Overview) checkState(state string) bool {
	if _, ok := page.states[state]; ok {
		return true
	}
	return false
}

func (page *Overview) update(gtx *layout.Context) {
	page.updateBalance()
	page.updateSyncData()
	page.updateSyncProgressData()
	page.updateWalletSyncBox(gtx)
}

// updatePage updates the state of the overview page
func (page *Overview) updateBalance() {
	page.balance = dcrutil.Amount(page.walletInfo.TotalBalance).String()
}

func (page *Overview) updateSyncData() {
	if page.walletInfo.Synced {
		page.syncButton.Text = "Disconnect"
		page.syncStatus.Text = "Synced"
		page.onlineStatus.Text = "Online"
	} else if page.walletInfo.Syncing {
		page.syncButton.Text = "Cancel"
		page.syncStatus.Text = "Syncing..."
		page.onlineStatus.Text = "Online"
	} else {
		page.syncStatus.Text = "Not synced"
		page.syncButton.Text = "Reconnect"
		page.onlineStatus.Text = "Offline"
	}
}

func (page *Overview) updateSyncProgressData() {
	if page.checkState(StateSyncStatus) {
		page.walletSyncStatus = page.states[StateSyncStatus].(*wallet.SyncStatus)
		page.progress = float64(page.walletSyncStatus.Progress)
		page.progressPercentage.Text = fmt.Sprintf("%v%%", page.progress)
		page.timeLeft.Text = fmt.Sprintf("%v left", helper.RemainingSyncTime(page.walletSyncStatus.RemainingTime))
		page.headersFetched.Text = fmt.Sprintf("Fetching block headers. %v%%", page.walletSyncStatus.HeadersFetchProgress)
		page.connectedPeers.Text = fmt.Sprintf("%d", page.walletSyncStatus.ConnectedPeers)
		page.syncSteps.Text = fmt.Sprintf("Step %d/%d", page.walletSyncStatus.Steps, page.walletSyncStatus.TotalSteps)
	}
	page.latestBlock.Text = fmt.Sprintf("%v . %v ago", page.walletInfo.BestBlockHeight,
		helper.LastBlockSync(page.walletInfo.BestBlockTime))
}

func (page *Overview) updateWalletSyncBox(gtx *layout.Context) {
	var overallBlockHeight int32

	page.walletSyncBoxes = []func(){}
	if page.checkState(StateSyncStatus) {
		overallBlockHeight = page.walletSyncStatus.HeadersToFetch
	}

	for i := 0; i < len(page.walletInfo.Wallets); i++ {
		w := page.walletInfo.Wallets[i]
		if w.BestBlockHeight > overallBlockHeight {
			overallBlockHeight = w.BestBlockHeight
		}
		blockHeightProgress := fmt.Sprintf("%v of %v", w.BestBlockHeight, overallBlockHeight)
		status := helper.WalletSyncStatus(w, overallBlockHeight)
		progress := helper.WalletSyncProgressTime(w.BlockTimestamp)
		details := page.syncDetail(w.Name, status, blockHeightProgress, progress)
		uniform := layout.UniformInset(units.Padding)
		page.walletSyncBoxes = append(page.walletSyncBoxes,
			func() {
				page.walletSyncBox(gtx, uniform, details)
			})
	}
}

func (page *Overview) updateRecentTransactions(gtx *layout.Context) {
	// todo
}

// content lays out the entire content for overview page.
func (page *Overview) content(gtx *layout.Context) {
	pageContent := []func(){
		func() {
			layout.Inset{Top: units.PageMarginTop}.Layout(gtx, func() {
				page.column.Layout(gtx,
					layout.Rigid(func() {
						layoutBalance(gtx, page.balance, page.mainBalance, page.subBalance)
					}),
					layout.Rigid(func() {
						page.balanceTitle.Layout(gtx)
					}),
				)

			})
		},
		func() {
			page.recentTransactionsColumn(gtx)
		},
		func() {
			layout.Inset{Bottom: units.ContainerPadding}.Layout(gtx, func() {
				page.syncStatusColumn(gtx)
			})
		},
	}
	page.listContainer.Layout(gtx, len(pageContent), func(i int) {
		layout.UniformInset(units.NoPadding).Layout(gtx, pageContent[i])
	})
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
				layoutBalance(gtx, page.transactionAmount, page.transactionMainAmount, page.transactionSubAmount)
			})
		}),
		layout.Flexed(1, func() {
			layout.Align(layout.E).Layout(gtx, func() {
				page.row.Layout(gtx,
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
					layout.Rigid(func() {
						margin.Layout(gtx, func() {
							page.transactionStatus.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// syncStatusColumn lays out content for displaying sync status.
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
				if page.walletInfo.Syncing {
					page.syncActiveContent(gtx, uniform)
				} else {
					page.syncDormantContent(gtx, uniform)
				}
			}),
		)
	})
}

// syncingContent lays out sync status content when the wallet is syncing
func (page *Overview) syncActiveContent(gtx *layout.Context, uniform layout.Inset) {
	page.column.Layout(gtx,
		layout.Rigid(func() {
			page.progressBarRow(gtx, uniform)
		}),
		layout.Rigid(func() {
			page.progressStatusRow(gtx, uniform)
		}),
		layout.Rigid(func() {
			page.walletSyncRow(gtx, uniform)
		}),
	)
}

// syncDormantContent lays out sync status content when the wallet is synced or not connected
func (page *Overview) syncDormantContent(gtx *layout.Context, uniform layout.Inset) {
	page.column.Layout(gtx,
		layout.Rigid(func() {
			page.endToEndRow(gtx, uniform, page.latestBlockTitle, page.latestBlock)
		}),
	)
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
				// stack a button on a card widget to produce a transparent button.
				layout.Align(layout.E).Layout(gtx, func() {
					layout.Stack{}.Layout(gtx,
						layout.Stacked(func() {
							page.syncButtonCard.SetColor(values.ButtonGray)
							page.syncButtonCard.SetWidth(values.SyncButtonWidth)
							page.syncButtonCard.SetHeight(values.SyncButtonHeight)
							page.syncButtonCard.Layout(gtx, float32(gtx.Px(units.DefaultButtonRadius)))
						}),
						layout.Expanded(func() {
							gtx.Constraints.Width.Max = values.SyncButtonWidth
							gtx.Constraints.Height.Max = values.SyncButtonHeight
							layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Flexed(1, func() {
									layout.Align(layout.Center).Layout(gtx, func() {
										layout.Flex{Axis: layout.Vertical}.Layout(gtx,
											layout.Flexed(1, func() {
												layout.Align(layout.Center).Layout(gtx, func() {
													page.syncButton.Font.Size = units.SyncButtonTextSize
													page.syncButton.Color = values.ButtonGray
													page.syncButton.Background = values.White
													gtx.Constraints.Width.Min = values.SyncButtonWidth - values.SyncButtonBorder
													page.syncButton.Layout(gtx, page.syncButtonWidget)
													for page.syncButtonWidget.Clicked(gtx) {
														page.triggerSync()
													}
												})
											}),
										)
									})
								}),
							)
						}),
					)
				})
			}),
		)
	})
}

func (page *Overview) triggerSync() {
	if page.walletInfo.Syncing || page.walletInfo.Synced {
		page.wallet.CancelSync()
	} else {
		page.wallet.StartSync()
	}
}

// syncBoxTitleRow lays out the progress bar.
func (page *Overview) progressBarRow(gtx *layout.Context, inset layout.Inset) {
	inset.Layout(gtx, func() {
		page.progressBar.Layout(gtx, page.progress)
	})
}

// syncBoxTitleRow lays out the progress status when the wallet is syncing.
func (page *Overview) progressStatusRow(gtx *layout.Context, inset layout.Inset) {
	page.endToEndRow(gtx, inset, page.progressPercentage, page.timeLeft)
}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func (page *Overview) walletSyncRow(gtx *layout.Context, inset layout.Inset) {
	page.columnMargin.Layout(gtx, func() {
		page.column.Layout(gtx,
			layout.Rigid(func() {
				page.endToEndRow(gtx, inset, page.syncSteps, page.headersFetched)
			}),
			layout.Rigid(func() {
				page.endToEndRow(gtx, inset, page.connectedPeersTitle, page.connectedPeers)
			}),
			layout.Rigid(func() {
				page.walletSyncList.Layout(gtx, len(page.walletSyncBoxes), func(i int) {
					if i == 0 {
						layout.UniformInset(units.NoPadding).Layout(gtx, page.walletSyncBoxes[i])
					} else {
						layout.Inset{Left: units.ColumnMargin}.Layout(gtx, page.walletSyncBoxes[i])
					}
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
				page.walletSyncCard.SetWidth(gtx.Px(units.WalletSyncBoxWidthMin))
				page.walletSyncCard.SetHeight(gtx.Px(units.WalletSyncBoxHeightMin))
				page.walletSyncCard.Layout(gtx, 0)
			}),
			layout.Stacked(func() {
				uniform := layout.UniformInset(units.SyncBoxPadding)
				uniform.Layout(gtx, func() {
					gtx.Constraints.Width.Min = gtx.Px(units.WalletSyncBoxContentWidth)
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

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func layoutBalance(gtx *layout.Context, balance string, main, sub material.Label) {
	mainText, subText := helper.BreakBalance(balance)
	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func() {
			main.Text = mainText
			main.Layout(gtx)
		}),
		layout.Rigid(func() {
			sub.Text = subText
			sub.Layout(gtx)
		}),
	)
}
