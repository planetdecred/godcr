package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/raedahgroup/godcr/wallet"

	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const (
	balanceTitle         = "Current Total Balance"
	statusTitle          = "Wallet Status"
	stepsTitle           = "Step"
	transactionsTitle    = "Recent Transactions"
	connectedPeersTitle  = "Connected peers count"
	headersFetchedTitle  = "Block header fetched"
	syncingProgressTitle = "Syncing progress"
	latestBlockTitle     = "Last Block Height"
	lastSyncedTitle      = "Last Block Mined"
	noTransaction        = "no transactions"
	offlineStatus        = "Offline"
	onlineStatus         = "Online"
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

type overviewPage struct {
	gtx                *layout.Context
	theme              *decredmaterial.Theme
	walletInfo         *wallet.MultiWalletInfo
	walletSyncStatus   *wallet.SyncStatus
	walletTransactions *wallet.Transactions
	*inputs
	*outputs
}

func (win *Window) OverviewPage() {
	if win.walletInfo.LoadedWallets == 0 {
		win.Page(func() {
			win.outputs.noWallet.Layout(win.gtx)
		})
		return
	}
	body := func() {
		page := overviewPage{
			gtx:                win.gtx,
			theme:              win.theme,
			walletInfo:         win.walletInfo,
			walletSyncStatus:   win.walletSyncStatus,
			walletTransactions: win.walletTransactions,
			inputs:             &win.inputs,
			outputs:            &win.outputs,
		}
		container := layout.Inset{Left: containerPadding, Right: containerPadding}
		container.Layout(win.gtx, func() {
			page.layoutPage()
		})
	}
	win.Page(body)
}

// syncDetail returns a walletSyncDetails object containing data of a single wallet sync box
func (page overviewPage) syncDetail(name, status, headersFetched, progress string) walletSyncDetails {
	return walletSyncDetails{
		name:               page.theme.Caption(name),
		status:             page.theme.Caption(status),
		blockHeaderFetched: page.theme.Caption(headersFetched),
		syncingProgress:    page.theme.Caption(progress),
	}
}

// layout lays out the entire content for overview page.
func (page overviewPage) layoutPage() {
	gtx := page.gtx
	walletInfo := page.walletInfo
	theme := page.theme

	pageContent := []func(){
		func() {
			layout.Inset{Top: pageMarginTop}.Layout(gtx, func() {
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func() {
						mainBalance := theme.H4("")
						subBalance := theme.H6("")
						page.layoutBalance(walletInfo.TotalBalance, mainBalance, subBalance)
					}),
					layout.Rigid(func() {
						theme.Caption(balanceTitle).Layout(gtx)
					}),
				)

			})
		},
		func() {
			page.recentTransactionsColumn()
		},
		func() {
			layout.Inset{Bottom: containerPadding}.Layout(gtx, func() {
				page.syncStatusColumn()
			})
		},
	}

	listContainer.Layout(gtx, len(pageContent), func(i int) {
		layout.UniformInset(noPadding).Layout(gtx, pageContent[i])
	})
}

// recentTransactionsColumn lays out the list of recent transactions.
func (page overviewPage) recentTransactionsColumn() {
	theme := page.theme
	gtx := page.gtx
	var transactionRows []func()
	if len(page.walletTransactions.Txs) > 0 {
		for _, txn := range page.walletTransactions.Recent {
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
				page.recentTransactionRow(txnWidgets)
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
				theme.Caption(transactionsTitle).Layout(page.gtx)
			}),
			layout.Rigid(func() {
				list := &layout.List{Axis: layout.Vertical}
				list.Layout(page.gtx, len(transactionRows), func(i int) {
					layout.Inset{Top: padding}.Layout(page.gtx, transactionRows[i])
				})
			}),
			layout.Rigid(func() {
				if len(transactionRows) > 5 {
					layout.Center.Layout(page.gtx, func() {
						layout.Inset{Top: padding}.Layout(page.gtx, func() {
							layout.Stack{}.Layout(page.gtx,
								layout.Expanded(func() {
									layout.Center.Layout(page.gtx, func() {
										gtx.Constraints.Width.Min = moreButtonWidth
										gtx.Constraints.Height.Max = moreButtonHeight
										moreButton := page.outputs.moreDiag
										moreButton.TextSize = syncButtonTextSize
										moreButton.Layout(gtx, &page.inputs.toTransactions)
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
func (page overviewPage) recentTransactionRow(txn transactionWidgets) {
	gtx := page.gtx
	margin := layout.UniformInset(transactionsRowMargin)
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			margin.Layout(gtx, func() {
				page.layoutBalance(txn.balance, txn.mainBalance, txn.subBalance)
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
func (page overviewPage) syncStatusColumn() {
	gtx := page.gtx
	uniform := layout.UniformInset(padding)
	layout.Inset{Top: columnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				page.syncBoxTitleRow(uniform)
			}),
			layout.Rigid(func() {
				page.syncStatusTextRow(uniform)
			}),
			layout.Rigid(func() {
				if page.walletInfo.Syncing {
					page.syncActiveContent(uniform)
				} else {
					page.syncDormantContent(uniform)
				}
			}),
		)
	})
}

// syncingContent lays out sync status content when the wallet is syncing
func (page overviewPage) syncActiveContent(uniform layout.Inset) {
	layout.Flex{Axis: layout.Vertical}.Layout(page.gtx,
		layout.Rigid(func() {
			page.progressBarRow(uniform)
		}),
		layout.Rigid(func() {
			page.progressStatusRow(uniform)
		}),
		layout.Rigid(func() {
			page.walletSyncRow(uniform)
		}),
	)
}

// syncDormantContent lays out sync status content when the wallet is synced or not connected
func (page overviewPage) syncDormantContent(uniform layout.Inset) {
	layout.Flex{Axis: layout.Vertical}.Layout(page.gtx,
		layout.Rigid(func() {
			latestBlockTitleLabel := page.theme.Body1(latestBlockTitle)
			blockLabel := page.theme.Body1(fmt.Sprintf("%v", page.walletInfo.BestBlockHeight))
			page.endToEndRow(uniform, latestBlockTitleLabel, blockLabel)
		}),
		layout.Rigid(func() {
			lastSyncedLabel := page.theme.Body1(lastSyncedTitle)
			blockLabel := page.theme.Body1(fmt.Sprintf("%v ago", page.walletInfo.LastSyncTime))
			page.endToEndRow(uniform, lastSyncedLabel, blockLabel)
		}),
	)
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func (page overviewPage) endToEndRow(inset layout.Inset, leftLabel, rightLabel decredmaterial.Label) {
	gtx := page.gtx
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
func (page overviewPage) syncBoxTitleRow(inset layout.Inset) {
	statusTitleLabel := page.theme.Caption(statusTitle)
	statusTitleLabel.Color = gray
	statusLabel := page.theme.Body1(offlineStatus)
	if page.walletInfo.Synced || page.walletInfo.Syncing {
		statusLabel.Text = onlineStatus
	}
	page.endToEndRow(inset, statusTitleLabel, statusLabel)
}

// syncBoxTitleRow lays out sync status text and sync button.
func (page overviewPage) syncStatusTextRow(inset layout.Inset) {
	gtx := page.gtx
	inset.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.5, func() {
				layout.W.Layout(gtx, func() {
					syncStatusLabel := page.theme.H6(notSyncedStatus)
					if page.walletInfo.Syncing {
						syncStatusLabel.Text = syncingStatus
					} else if page.walletInfo.Synced {
						syncStatusLabel.Text = syncedStatus
					}
					syncStatusLabel.Layout(page.gtx)
				})
			}),
			layout.Flexed(1, func() {
				// stack a button on a card widget to produce a transparent button.
				layout.E.Layout(gtx, func() {
					gtx.Constraints.Width.Min = syncButtonWidth
					gtx.Constraints.Height.Max = syncButtonHeight
					syncButton := page.outputs.sync
					syncButton.TextSize = syncButtonTextSize
					if page.walletInfo.Synced {
						syncButton.Text = "Disconnect"
					}
					syncButton.Layout(gtx, &page.inputs.sync)
				})
			}),
		)
	})
}

// syncBoxTitleRow lays out the progress bar.
func (page overviewPage) progressBarRow(inset layout.Inset) {
	inset.Layout(page.gtx, func() {
		progress := page.walletSyncStatus.Progress
		page.gtx.Constraints.Height.Max = 20
		page.theme.ProgressBar(float64(progress)).Layout(page.gtx)
	})
}

// syncBoxTitleRow lays out the progress status when the wallet is syncing.
func (page overviewPage) progressStatusRow(inset layout.Inset) {
	timeLeft := page.walletSyncStatus.RemainingTime
	if timeLeft == "" {
		timeLeft = "0s"
	}

	percentageLabel := page.theme.Body1(fmt.Sprintf("%v%%", page.walletSyncStatus.Progress))
	timeLeftLabel := page.theme.Body1(fmt.Sprintf("%v left", timeLeft))
	page.endToEndRow(inset, percentageLabel, timeLeftLabel)
}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func (page overviewPage) walletSyncRow(inset layout.Inset) {
	gtx := page.gtx
	layout.Inset{Top: columnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				completedSteps := page.walletSyncStatus.Steps
				totalSteps := page.walletSyncStatus.TotalSteps
				completedStepsLabel := page.theme.Caption(fmt.Sprintf("%s %d/%d", stepsTitle, completedSteps, totalSteps))
				completedStepsLabel.Color = gray
				headersFetchedLabel := page.theme.Body1(fmt.Sprintf("%s. %v%%", fetchingBlockHeaders,
					page.walletSyncStatus.HeadersFetchProgress))
				headersFetchedLabel.Color = gray
				page.endToEndRow(inset, completedStepsLabel, headersFetchedLabel)
			}),
			layout.Rigid(func() {
				connectedPeersTitleLabel := page.theme.Caption(connectedPeersTitle)
				connectedPeersTitleLabel.Color = gray
				connectedPeersLabel := page.theme.Body1(fmt.Sprintf("%d", page.walletSyncStatus.ConnectedPeers))
				page.endToEndRow(inset, connectedPeersTitleLabel, connectedPeersLabel)
			}),
			layout.Rigid(func() {
				var overallBlockHeight int32
				var walletSyncBoxes []func()

				if page.walletSyncStatus != nil {
					overallBlockHeight = page.walletSyncStatus.HeadersToFetch
				}

				for i := 0; i < len(page.walletInfo.Wallets); i++ {
					w := page.walletInfo.Wallets[i]
					if w.BestBlockHeight > overallBlockHeight {
						overallBlockHeight = w.BestBlockHeight
					}
					blockHeightProgress := fmt.Sprintf("%v of %v", w.BestBlockHeight, overallBlockHeight)
					details := page.syncDetail(w.Name, w.Status, blockHeightProgress, w.DaysBehind)
					uniform := layout.UniformInset(padding)
					walletSyncBoxes = append(walletSyncBoxes,
						func() {
							page.walletSyncBox(uniform, details)
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
func (page overviewPage) walletSyncBox(inset layout.Inset, details walletSyncDetails) {
	gtx := page.gtx
	layout.Inset{Top: columnMargin}.Layout(gtx, func() {
		gtx.Constraints.Width.Min = gtx.Px(walletSyncBoxContentWidth)
		gtx.Constraints.Width.Max = gtx.Constraints.Width.Min
		decredmaterial.Card{Inset: layout.UniformInset(unit.Dp(0))}.Layout(gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					page.endToEndRow(inset, details.name, details.status)
				}),
				layout.Rigid(func() {
					headersFetchedTitleLabel := page.theme.Caption(headersFetchedTitle)
					headersFetchedTitleLabel.Color = gray
					page.endToEndRow(inset, headersFetchedTitleLabel, details.blockHeaderFetched)
				}),
				layout.Rigid(func() {
					progressTitleLabel := page.theme.Caption(syncingProgressTitle)
					progressTitleLabel.Color = gray
					page.endToEndRow(inset, progressTitleLabel, details.syncingProgress)
				}),
			)
		})
	})
}

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page overviewPage) layoutBalance(amount string, main, sub decredmaterial.Label) {
	mainText, subText := page.breakBalance(amount)
	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(page.gtx,
		layout.Rigid(func() {
			main.Text = mainText
			main.Layout(page.gtx)
		}),
		layout.Rigid(func() {
			sub.Text = subText
			sub.Layout(page.gtx)
		}),
	)
}

// breakBalance takes the balance string and returns it in two slices
func (page overviewPage) breakBalance(balance string) (b1, b2 string) {
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
