package ui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
)

const (
	balanceTitle         = "Current Total Balance"
	statusTitle          = "Wallet Status"
	stepsTitle           = "Step"
	transactionsTitle    = "Recent Transactions"
	connectedPeersTitle  = "Connected peers count"
	headersFetchedTitle  = "Block header fetched"
	syncingProgressTitle = "Syncing progress"
	latestBlockTitle     = "Latest Block"
	noTransaction        = "no transactions"
	onlineStatus         = "Offline"
	syncingStatus        = "Syncing..."
	notSyncedStatus      = "Not Synced"
	syncedStatus         = "Synced"
	fetchingBlockHeaders = "Fetching block headers"
)

var (
	listContainer  = &layout.List{Axis: layout.Vertical}
	walletSyncList = &layout.List{Axis: layout.Horizontal}

	syncButtonHeight = 70
	syncButtonWidth  = 145
	moreButtonWidth  = 115
	moreButtonHeight = 70

	padding                   = unit.Dp(5)
	containerPadding          = unit.Dp(20)
	pageMarginTop             = unit.Dp(50)
	columnMargin              = unit.Dp(30)
	transactionsRowMargin     = unit.Dp(10)
	noPadding                 = unit.Dp(0)
	walletSyncBoxContentWidth = unit.Dp(280)
	syncButtonTextSize        = unit.Dp(10)

	gray = color.RGBA{137, 151, 165, 255}
)

// walletSyncDetails contains sync data for each wallet when a sync
// is in progress.
type walletSyncDetails struct {
	name               decredmaterial.Label
	status             decredmaterial.Label
	blockHeaderFetched decredmaterial.Label
	syncingProgress    decredmaterial.Label
}

type transactionWidgets struct {
	wallet      decredmaterial.Label
	balance     string
	mainBalance decredmaterial.Label
	subBalance  decredmaterial.Label
	date        decredmaterial.Label
	status      decredmaterial.Label
}

func (win *Window) OverviewPage() {
	if win.walletInfo.LoadedWallets == 0 {
		win.Page(func() {
			win.outputs.noWallet.Layout(win.gtx)
		})
		return
	}
	body := func() {
		container := layout.Inset{Left: containerPadding, Right: containerPadding}
		container.Layout(win.gtx, func() {
			layoutPage(win)
		})
	}
	win.Page(body)
}

// syncDetail returns a walletSyncDetails object containing data of a single wallet sync box
func syncDetail(theme *decredmaterial.Theme, name, status, headersFetched, progress string) walletSyncDetails {
	return walletSyncDetails{
		name:               theme.Caption(name),
		status:             theme.Caption(status),
		blockHeaderFetched: theme.Caption(headersFetched),
		syncingProgress:    theme.Caption(progress),
	}
}

// layout lays out the entire content for overview page.
func layoutPage(win *Window) {
	gtx := win.gtx
	walletInfo := win.walletInfo
	theme := win.theme

	pageContent := []func(){
		func() {
			layout.Inset{Top: pageMarginTop}.Layout(gtx, func() {
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func() {
						mainBalance := theme.H4("")
						subBalance := theme.H6("")
						layoutBalance(gtx, walletInfo.TotalBalance, mainBalance, subBalance)
					}),
					layout.Rigid(func() {
						theme.Caption(balanceTitle).Layout(gtx)
					}),
				)

			})
		},
		func() {
			recentTransactionsColumn(win)
		},
		func() {
			layout.Inset{Bottom: containerPadding}.Layout(gtx, func() {
				syncStatusColumn(win)
			})
		},
	}

	listContainer.Layout(gtx, len(pageContent), func(i int) {
		layout.UniformInset(noPadding).Layout(gtx, pageContent[i])
	})
}

// recentTransactionsColumn lays out the list of recent transactions.
func recentTransactionsColumn(win *Window) {
	theme := win.theme
	gtx := win.gtx
	recentTransactions := win.walletTransactions.Recent

	var transactionRows []func()
	if len(win.walletTransactions.Txs) > 0 {
		for _, txn := range recentTransactions {
			txnWidgets := transactionWidgets{
				wallet:      theme.Body1(txn.WalletName),
				balance:     txn.Balance,
				mainBalance: theme.Body1(""),
				subBalance:  theme.Caption(""),
				date: theme.Body1(fmt.Sprintf("%v",
					dcrlibwallet.ExtractDateOrTime(txn.Txn.Timestamp))),
				status: theme.Body1(txn.Status),
			}

			transactionRows = append(transactionRows, func() {
				recentTransactionRow(gtx, txnWidgets)
			})
		}
	} else {
		transactionRows = append(transactionRows, func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layout.Flexed(1, func() {
				layout.Center.Layout(gtx, func() {
					label := theme.Caption(noTransaction)
					label.Color = gray
					label.Layout(gtx)
				})
			}))
		})
	}

	layout.Inset{Top: columnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				theme.Caption(transactionsTitle).Layout(gtx)
			}),
			layout.Rigid(func() {
				list := &layout.List{Axis: layout.Vertical}
				list.Layout(gtx, len(transactionRows), func(i int) {
					layout.Inset{Top: padding}.Layout(gtx, transactionRows[i])
				})
			}),
			layout.Rigid(func() {
				if len(transactionRows) > 5 {
					layout.Center.Layout(gtx, func() {
						layout.Inset{Top: padding}.Layout(gtx, func() {
							layout.Stack{}.Layout(gtx,
								layout.Expanded(func() {
									layout.Center.Layout(gtx, func() {
										gtx.Constraints.Width.Min = moreButtonWidth
										gtx.Constraints.Height.Max = moreButtonHeight
										moreButton := win.outputs.more
										moreButton.TextSize = syncButtonTextSize
										moreButton.Layout(gtx, &win.inputs.toTransactions)
									})
								}),
							)
						})
					})
				}
			}),
		)
	})
}

// recentTransactionRow lays out a single row of a recent transaction.
func recentTransactionRow(gtx *layout.Context, txn transactionWidgets) {
	margin := layout.UniformInset(transactionsRowMargin)
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			margin.Layout(gtx, func() {
				layoutBalance(gtx, txn.balance, txn.mainBalance, txn.subBalance)
			})
		}),
		layout.Flexed(1, func() {
			layout.E.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						margin.Layout(gtx, func() {
							txn.wallet.Layout(gtx)
						})
					}),
					layout.Rigid(func() {
						margin.Layout(gtx, func() {
							txn.date.Layout(gtx)
						})
					}),
					layout.Rigid(func() {
						margin.Layout(gtx, func() {
							txn.status.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

// syncStatusColumn lays out content for displaying sync status.
func syncStatusColumn(win *Window) {
	gtx := win.gtx
	uniform := layout.UniformInset(padding)
	layout.Inset{Top: columnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				syncBoxTitleRow(win, uniform)
			}),
			layout.Rigid(func() {
				syncStatusTextRow(win, uniform)
			}),
			layout.Rigid(func() {
				if win.walletInfo.Syncing {
					syncActiveContent(win, uniform)
				} else {
					syncDormantContent(win, uniform)
				}
			}),
		)
	})
}

// syncingContent lays out sync status content when the wallet is syncing
func syncActiveContent(win *Window, uniform layout.Inset) {
	gtx := win.gtx
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			progressBarRow(win, uniform)
		}),
		layout.Rigid(func() {
			progressStatusRow(win, uniform)
		}),
		layout.Rigid(func() {
			walletSyncRow(win, uniform)
		}),
	)
}

// syncDormantContent lays out sync status content when the wallet is synced or not connected
func syncDormantContent(win *Window, uniform layout.Inset) {
	layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
		layout.Rigid(func() {
			latestBlockTitleLabel := win.theme.Body1(latestBlockTitle)
			blockLabel := win.theme.Body1(fmt.Sprintf("%v . %v ago", win.walletInfo.BestBlockHeight,
				win.walletInfo.LastSyncTime))
			endToEndRow(win.gtx, uniform, latestBlockTitleLabel, blockLabel)
		}),
	)
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func endToEndRow(gtx *layout.Context, inset layout.Inset, leftLabel, rightLabel decredmaterial.Label) {
	inset.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func() {
				leftLabel.Layout(gtx)
			}),
			layout.Flexed(1, func() {
				layout.E.Layout(gtx, func() {
					rightLabel.Layout(gtx)
				})
			}),
		)
	})
}

// syncBoxTitleRow lays out widgets in the title row inside the sync status box.
func syncBoxTitleRow(win *Window, inset layout.Inset) {
	statusTitleLabel := win.theme.Caption(statusTitle)
	statusTitleLabel.Color = gray
	statusLabel := win.theme.Body1(onlineStatus)
	endToEndRow(win.gtx, inset, statusTitleLabel, statusLabel)
}

// syncBoxTitleRow lays out sync status text and sync button.
func syncStatusTextRow(win *Window, inset layout.Inset) {
	gtx := win.gtx
	inset.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.5, func() {
				layout.W.Layout(gtx, func() {
					syncText := ""
					if win.walletInfo.Syncing {
						syncText = syncingStatus
					} else if win.walletInfo.Synced {
						syncText = syncedStatus
					} else {
						syncText = notSyncedStatus
					}
					win.theme.H6(syncText).Layout(gtx)
				})
			}),
			layout.Flexed(1, func() {
				// stack a button on a card widget to produce a transparent button.
				layout.E.Layout(gtx, func() {
					gtx.Constraints.Width.Min = syncButtonWidth
					gtx.Constraints.Height.Max = syncButtonHeight
					syncButton := win.outputs.sync
					syncButton.TextSize = syncButtonTextSize
					if win.walletInfo.Synced {
						syncButton.Text = "Disconnect"
					}
					syncButton.Layout(gtx, &win.inputs.sync)
				})
			}),
		)
	})
}

// syncBoxTitleRow lays out the progress bar.
func progressBarRow(win *Window, inset layout.Inset) {
	inset.Layout(win.gtx, func() {
		progress := win.walletSyncStatus.Progress
		win.theme.ProgressBar().Layout(win.gtx, float64(progress))
	})
}

// syncBoxTitleRow lays out the progress status when the wallet is syncing.
func progressStatusRow(win *Window, inset layout.Inset) {
	gtx := win.gtx
	percentageLabel := win.theme.Body1(fmt.Sprintf("%v%%", win.walletSyncStatus.Progress))
	timeLeftLabel := win.theme.Body1(fmt.Sprintf("%v left", win.walletSyncStatus.RemainingTime))
	endToEndRow(gtx, inset, percentageLabel, timeLeftLabel)
}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func walletSyncRow(win *Window, inset layout.Inset) {
	gtx := win.gtx
	layout.Inset{Top: columnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				completedSteps := win.walletSyncStatus.Steps
				totalSteps := win.walletSyncStatus.TotalSteps
				completedStepsLabel := win.theme.Caption(fmt.Sprintf("%s %d/%d", stepsTitle, completedSteps, totalSteps))
				completedStepsLabel.Color = gray
				headersFetchedLabel := win.theme.Body1(fmt.Sprintf("%s. %v%%", fetchingBlockHeaders,
					win.walletSyncStatus.HeadersFetchProgress))
				headersFetchedLabel.Color = gray
				endToEndRow(win.gtx, inset, completedStepsLabel, headersFetchedLabel)
			}),
			layout.Rigid(func() {
				connectedPeersTitleLabel := win.theme.Caption(connectedPeersTitle)
				connectedPeersTitleLabel.Color = gray
				connectedPeersLabel := win.theme.Body1(fmt.Sprintf("%d", win.walletSyncStatus.ConnectedPeers))
				endToEndRow(gtx, inset, connectedPeersTitleLabel, connectedPeersLabel)
			}),
			layout.Rigid(func() {
				var overallBlockHeight int32
				var walletSyncBoxes []func()

				if win.walletSyncStatus != nil {
					overallBlockHeight = win.walletSyncStatus.HeadersToFetch
				}

				for i := 0; i < len(win.walletInfo.Wallets); i++ {
					w := win.walletInfo.Wallets[i]
					if w.BestBlockHeight > overallBlockHeight {
						overallBlockHeight = w.BestBlockHeight
					}
					blockHeightProgress := fmt.Sprintf("%v of %v", w.BestBlockHeight, overallBlockHeight)
					details := syncDetail(win.theme, w.Name, w.Status, blockHeightProgress, w.DaysBehind)
					uniform := layout.UniformInset(padding)
					walletSyncBoxes = append(walletSyncBoxes,
						func() {
							walletSyncBox(win, uniform, details)
						})
				}

				walletSyncList.Layout(gtx, len(walletSyncBoxes), func(i int) {
					if i == 0 {
						layout.UniformInset(noPadding).Layout(gtx, walletSyncBoxes[i])
					} else {
						layout.Inset{Left: columnMargin}.Layout(gtx, walletSyncBoxes[i])
					}
				})
			}),
		)
	})
}

// walletSyncBox lays out the wallet syncing details of a single wallet.
func walletSyncBox(win *Window, inset layout.Inset, details walletSyncDetails) {
	gtx := win.gtx
	layout.Inset{Top: columnMargin}.Layout(gtx, func() {
		gtx.Constraints.Width.Min = gtx.Px(walletSyncBoxContentWidth)
		gtx.Constraints.Width.Max = gtx.Constraints.Width.Min
		decredmaterial.Card{Inset: layout.UniformInset(unit.Dp(0))}.Layout(gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					endToEndRow(gtx, inset, details.name, details.status)
				}),
				layout.Rigid(func() {
					headersFetchedTitleLabel := win.theme.Caption(headersFetchedTitle)
					headersFetchedTitleLabel.Color = gray
					endToEndRow(gtx, inset, headersFetchedTitleLabel, details.blockHeaderFetched)
				}),
				layout.Rigid(func() {
					progressTitleLabel := win.theme.Caption(syncingProgressTitle)
					progressTitleLabel.Color = gray
					endToEndRow(gtx, inset, progressTitleLabel, details.syncingProgress)
				}),
			)
		})
	})
}

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func layoutBalance(gtx *layout.Context, amount string, main, sub decredmaterial.Label) {
	mainText, subText := breakBalance(amount)
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

// breakBalance takes the balance string and returns it in two slices
func breakBalance(balance string) (b1, b2 string) {
	balanceParts := strings.Split(balance, ".")
	if len(balanceParts) == 1 {
		return balanceParts[0], ""
	}
	b1 = balanceParts[0]
	b2 = balanceParts[1]
	b1 = b1 + "." + b2[:2]
	b2 = b2[2:]
	return
}
