package ui

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/raedahgroup/godcr/wallet"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageOverview = "overview"

var syncButtonTextSize = unit.Dp(10)

type overviewPageText struct {
	balanceTitle,
	statusTitle,
	stepsTitle,
	transactionsTitle,
	connectedPeersTitle,
	headersFetchedTitle,
	syncingProgressTitle,
	latestBlockTitle,
	lastSyncedTitle,
	noTransaction,
	offlineStatus,
	onlineStatus,
	syncingStatus,
	notSyncedStatus,
	syncedStatus,
	fetchingBlockHeaders,
	reconnect,
	disconnect,
	noWallet,
	cancel string
}

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
	direction   *decredmaterial.Icon
	mainBalance decredmaterial.Label
	subBalance  decredmaterial.Label
	date        decredmaterial.Label
	status      decredmaterial.Label
}

type overviewPage struct {
	listContainer, walletSyncList *layout.List
	gtx                           *layout.Context
	theme                         *decredmaterial.Theme
	walletInfo                    *wallet.MultiWalletInfo
	walletSyncStatus              *wallet.SyncStatus
	walletTransactions            **wallet.Transactions
	walletTransaction             **wallet.Transaction
	toTransactions, sync          decredmaterial.Button
	toTransactionsW, syncW        widget.Button
	toTransactionDetails          map[string]*gesture.Click

	text                      overviewPageText
	syncButtonHeight          int
	syncButtonWidth           int
	moreButtonWidth           int
	moreButtonHeight          int
	padding                   unit.Value
	containerPadding          unit.Value
	pageMarginTop             unit.Value
	columnMargin              unit.Value
	transactionsRowMargin     unit.Value
	noPadding                 unit.Value
	walletSyncBoxContentWidth unit.Value
	gray                      color.RGBA
}

func (win *Window) OverviewPage(c pageCommon) layout.Widget {
	page := overviewPage{
		gtx:                c.gtx,
		theme:              c.theme,
		walletInfo:         win.walletInfo,
		walletSyncStatus:   win.walletSyncStatus,
		walletTransactions: &win.walletTransactions,
		walletTransaction:  &win.walletTransaction,
		listContainer:      &layout.List{Axis: layout.Vertical},
		walletSyncList:     &layout.List{Axis: layout.Horizontal},
		toTransactions:     c.theme.Button("more"),

		syncButtonHeight: 70,
		syncButtonWidth:  145,
		moreButtonWidth:  115,
		moreButtonHeight: 70,

		padding:                   unit.Dp(5),
		containerPadding:          unit.Dp(20),
		pageMarginTop:             unit.Dp(50),
		columnMargin:              unit.Dp(30),
		transactionsRowMargin:     unit.Dp(10),
		noPadding:                 unit.Dp(0),
		walletSyncBoxContentWidth: unit.Dp(280),
		gray:                      color.RGBA{137, 151, 165, 255},
	}
	page.text = overviewPageText{
		balanceTitle:         "Current Total Balance",
		statusTitle:          "Wallet Status",
		stepsTitle:           "Step",
		transactionsTitle:    "Recent Transactions",
		connectedPeersTitle:  "Connected peers count",
		headersFetchedTitle:  "Block header fetched",
		syncingProgressTitle: "Syncing progress",
		latestBlockTitle:     "Last Block Height",
		lastSyncedTitle:      "Last Block Mined",
		noTransaction:        "no transactions",
		noWallet:             "No wallet loaded",
		offlineStatus:        "Offline",
		onlineStatus:         "Online",
		syncingStatus:        "Syncing...",
		notSyncedStatus:      "Not Synced",
		syncedStatus:         "Synced",
		fetchingBlockHeaders: "Fetching block headers",
		reconnect:            "Reconnect",
		disconnect:           "Disconnect",
		cancel:               "Cancel",
	}
	page.toTransactions.TextSize = unit.Dp(10)
	page.sync = c.theme.Button(page.text.reconnect)
	page.sync.TextSize = unit.Dp(10)

	return func() {
		page.Layout(c)
		page.Handler(c)
		page.updateToTransactionDetailsButtons()
	}
}

// Layout lays out the entire content for overview page.
func (page *overviewPage) Layout(c pageCommon) {
	if c.info.LoadedWallets == 0 {
		c.Layout(c.gtx, func() {
			layout.Center.Layout(c.gtx, func() {
				c.theme.H3(page.text.noWallet).Layout(c.gtx)
			})
		})
		return
	}

	gtx := page.gtx
	walletInfo := page.walletInfo
	theme := page.theme

	pageContent := []func(){
		func() {
			layout.Inset{Top: page.pageMarginTop}.Layout(gtx, func() {
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func() {
						mainBalance := theme.H4("")
						subBalance := theme.H6("")
						page.layoutBalance(walletInfo.TotalBalance, mainBalance, subBalance)
					}),
					layout.Rigid(func() {
						theme.Caption(page.text.balanceTitle).Layout(gtx)
					}),
				)

			})
		},
		func() {
			page.recentTransactionsColumn(c)
		},
		func() {
			layout.Inset{Bottom: page.containerPadding}.Layout(gtx, func() {
				page.syncStatusColumn()
			})
		},
	}

	c.Layout(c.gtx, func() {
		page.listContainer.Layout(gtx, len(pageContent), func(i int) {
			layout.UniformInset(page.noPadding).Layout(gtx, pageContent[i])
		})
	})
}

// syncDetail returns a walletSyncDetails object containing data of a single wallet sync box
func (page *overviewPage) syncDetail(name, status, headersFetched, progress string) walletSyncDetails {
	return walletSyncDetails{
		name:               page.theme.Caption(name),
		status:             page.theme.Caption(status),
		blockHeaderFetched: page.theme.Caption(headersFetched),
		syncingProgress:    page.theme.Caption(progress),
	}
}

// recentTransactionsColumn lays out the list of recent transactions.
func (page *overviewPage) recentTransactionsColumn(c pageCommon) {
	theme := page.theme
	gtx := page.gtx
	var transactionRows []func()

	if len((*page.walletTransactions).Txs) > 0 {
		for _, txn := range (*page.walletTransactions).Recent {
			txnWidgets := transactionWidgets{
				wallet:      theme.Body1(txn.WalletName),
				balance:     txn.Balance,
				mainBalance: theme.Body1(""),
				subBalance:  theme.Caption(""),
				date: theme.Body1(fmt.Sprintf("%v",
					dcrlibwallet.ExtractDateOrTime(txn.Txn.Timestamp))),
				status: theme.Body1(txn.Status),
			}
			if txn.Txn.Direction == dcrlibwallet.TxDirectionSent {
				txnWidgets.direction = c.icons.contentRemove
				txnWidgets.direction.Color = c.theme.Color.Danger
			} else {
				txnWidgets.direction = c.icons.contentAdd
				txnWidgets.direction.Color = c.theme.Color.Success
			}

			click := page.toTransactionDetails[txn.Txn.Hash]

			transactionRows = append(transactionRows, func() {
				page.recentTransactionRow(txnWidgets)
				pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
				click.Add(gtx.Ops)
			})
		}
	} else {
		transactionRows = append(transactionRows, func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layout.Flexed(1, func() {
				layout.Center.Layout(gtx, func() {
					label := theme.Caption(page.text.noTransaction)
					label.Color = page.gray
					label.Layout(gtx)
				})
			}))
		})
	}

	layout.Inset{Top: page.columnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				theme.Caption(page.text.transactionsTitle).Layout(page.gtx)
			}),
			layout.Rigid(func() {
				list := &layout.List{Axis: layout.Vertical}
				list.Layout(page.gtx, len(transactionRows), func(i int) {
					layout.Inset{Top: page.padding}.Layout(page.gtx, transactionRows[i])
				})
			}),
			layout.Rigid(func() {
				if len(transactionRows) > 5 {
					layout.Center.Layout(page.gtx, func() {
						layout.Inset{Top: page.padding}.Layout(page.gtx, func() {
							layout.Stack{}.Layout(page.gtx,
								layout.Expanded(func() {
									layout.Center.Layout(page.gtx, func() {
										gtx.Constraints.Width.Min = page.moreButtonWidth
										gtx.Constraints.Height.Max = page.moreButtonHeight
										page.toTransactions.Layout(gtx, &page.toTransactionsW)
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
func (page *overviewPage) recentTransactionRow(txn transactionWidgets) {
	gtx := page.gtx
	margin := layout.UniformInset(page.transactionsRowMargin)
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(12), Right: unit.Dp(25)}.Layout(gtx, func() {
				txn.direction.Layout(gtx, unit.Dp(16))
			})
		}),
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
func (page *overviewPage) syncStatusColumn() {
	gtx := page.gtx
	uniform := layout.UniformInset(page.padding)
	layout.Inset{Top: page.columnMargin}.Layout(gtx, func() {
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
func (page *overviewPage) syncActiveContent(uniform layout.Inset) {
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
func (page *overviewPage) syncDormantContent(uniform layout.Inset) {
	layout.Flex{Axis: layout.Vertical}.Layout(page.gtx,
		layout.Rigid(func() {
			latestBlockTitleLabel := page.theme.Body1(page.text.latestBlockTitle)
			blockLabel := page.theme.Body1(fmt.Sprintf("%v", page.walletInfo.BestBlockHeight))
			page.endToEndRow(uniform, latestBlockTitleLabel, blockLabel)
		}),
		layout.Rigid(func() {
			lastSyncedLabel := page.theme.Body1(page.text.lastSyncedTitle)
			blockLabel := page.theme.Body1(fmt.Sprintf("%v ago", page.walletInfo.LastSyncTime))
			page.endToEndRow(uniform, lastSyncedLabel, blockLabel)
		}),
	)
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func (page *overviewPage) endToEndRow(inset layout.Inset, leftLabel, rightLabel decredmaterial.Label) {
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
func (page *overviewPage) syncBoxTitleRow(inset layout.Inset) {
	statusTitleLabel := page.theme.Caption(page.text.statusTitle)
	statusTitleLabel.Color = page.gray
	statusLabel := page.theme.Body1(page.text.offlineStatus)
	if page.walletInfo.Synced || page.walletInfo.Syncing {
		statusLabel.Text = page.text.onlineStatus
	}
	page.endToEndRow(inset, statusTitleLabel, statusLabel)
}

// syncStatusTextRow lays out sync status text and sync button.
func (page *overviewPage) syncStatusTextRow(inset layout.Inset) {
	gtx := page.gtx
	inset.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.5, func() {
				layout.W.Layout(gtx, func() {
					syncStatusLabel := page.theme.H6(page.text.notSyncedStatus)
					if page.walletInfo.Syncing {
						syncStatusLabel.Text = page.text.syncingStatus
					} else if page.walletInfo.Synced {
						syncStatusLabel.Text = page.text.syncedStatus
					}
					syncStatusLabel.Layout(page.gtx)
				})
			}),
			layout.Flexed(1, func() {
				// stack a button on a card widget to produce a transparent button.
				layout.E.Layout(gtx, func() {
					gtx.Constraints.Width.Min = page.syncButtonWidth
					gtx.Constraints.Height.Max = page.syncButtonHeight
					if page.walletInfo.Synced {
						page.sync.Text = page.text.disconnect
					}
					page.sync.Layout(gtx, &page.syncW)
				})
			}),
		)
	})
}

// syncBoxTitleRow lays out the progress bar.
func (page *overviewPage) progressBarRow(inset layout.Inset) {
	inset.Layout(page.gtx, func() {
		progress := page.walletSyncStatus.Progress
		page.gtx.Constraints.Height.Max = 20
		page.theme.ProgressBar(float64(progress)).Layout(page.gtx)
	})
}

// syncBoxTitleRow lays out the progress status when the wallet is syncing.
func (page *overviewPage) progressStatusRow(inset layout.Inset) {
	timeLeft := page.walletSyncStatus.RemainingTime
	if timeLeft == "" {
		timeLeft = "0s"
	}

	percentageLabel := page.theme.Body1(fmt.Sprintf("%v%%", page.walletSyncStatus.Progress))
	timeLeftLabel := page.theme.Body1(fmt.Sprintf("%v left", timeLeft))
	page.endToEndRow(inset, percentageLabel, timeLeftLabel)
}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func (page *overviewPage) walletSyncRow(inset layout.Inset) {
	gtx := page.gtx
	layout.Inset{Top: page.columnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				completedSteps := page.walletSyncStatus.Steps
				totalSteps := page.walletSyncStatus.TotalSteps
				completedStepsLabel := page.theme.Caption(fmt.Sprintf("%s %d/%d", page.text.stepsTitle, completedSteps, totalSteps))
				completedStepsLabel.Color = page.gray
				headersFetchedLabel := page.theme.Body1(fmt.Sprintf("%s. %v%%", page.text.fetchingBlockHeaders,
					page.walletSyncStatus.HeadersFetchProgress))
				headersFetchedLabel.Color = page.gray
				page.endToEndRow(inset, completedStepsLabel, headersFetchedLabel)
			}),
			layout.Rigid(func() {
				connectedPeersTitleLabel := page.theme.Caption(page.text.connectedPeersTitle)
				connectedPeersTitleLabel.Color = page.gray
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
					uniform := layout.UniformInset(page.padding)
					walletSyncBoxes = append(walletSyncBoxes,
						func() {
							page.walletSyncBox(uniform, details)
						})
				}

				page.walletSyncList.Layout(gtx, len(walletSyncBoxes), func(i int) {
					if i == 0 {
						layout.UniformInset(page.noPadding).Layout(gtx, walletSyncBoxes[i])
					} else {
						layout.Inset{Left: page.columnMargin}.Layout(gtx, walletSyncBoxes[i])
					}
				})
			}),
		)
	})
}

// walletSyncBox lays out the wallet syncing details of a single wallet.
func (page *overviewPage) walletSyncBox(inset layout.Inset, details walletSyncDetails) {
	gtx := page.gtx
	layout.Inset{Top: page.columnMargin}.Layout(gtx, func() {
		gtx.Constraints.Width.Min = gtx.Px(page.walletSyncBoxContentWidth)
		gtx.Constraints.Width.Max = gtx.Constraints.Width.Min
		decredmaterial.Card{Inset: layout.UniformInset(unit.Dp(0))}.Layout(gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					page.endToEndRow(inset, details.name, details.status)
				}),
				layout.Rigid(func() {
					headersFetchedTitleLabel := page.theme.Caption(page.text.headersFetchedTitle)
					headersFetchedTitleLabel.Color = page.gray
					page.endToEndRow(inset, headersFetchedTitleLabel, details.blockHeaderFetched)
				}),
				layout.Rigid(func() {
					progressTitleLabel := page.theme.Caption(page.text.syncingProgressTitle)
					progressTitleLabel.Color = page.gray
					page.endToEndRow(inset, progressTitleLabel, details.syncingProgress)
				}),
			)
		})
	})
}

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page *overviewPage) layoutBalance(amount string, main, sub decredmaterial.Label) {
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
func (page *overviewPage) breakBalance(balance string) (b1, b2 string) {
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

func (page *overviewPage) updateToTransactionDetailsButtons() {
	recentTxs := (*page.walletTransactions).Recent
	shouldInit := false

	if len(recentTxs) != len(page.toTransactionDetails) {
		shouldInit = true
	}

	// When new block is added, check first block in recentTxs against map of toTransactionDetails
	if len(page.toTransactionDetails) > 0 {
		if _, found := page.toTransactionDetails[recentTxs[0].Txn.Hash]; !found {
			shouldInit = true
		}
	}

	if shouldInit {
		page.toTransactionDetails = make(map[string]*gesture.Click, len(recentTxs))
		for _, txn := range recentTxs {
			page.toTransactionDetails[txn.Txn.Hash] = &gesture.Click{}
		}
	}
}

func (page *overviewPage) Handler(c pageCommon) {
	if page.syncW.Clicked(page.gtx) {
		if page.walletInfo.Synced || page.walletInfo.Syncing {
			c.wallet.CancelSync()
			page.sync.Background = c.theme.Color.Primary
			page.sync.Text = page.text.reconnect
		} else {
			c.wallet.StartSync()
			page.sync.Background = c.theme.Color.Danger
			page.sync.Text = page.text.cancel
		}
	}
	if page.toTransactionsW.Clicked(page.gtx) {
		*c.page = PageTransactions
	}

	for has, click := range page.toTransactionDetails {
		for _, e := range click.Events(page.gtx) {
			if e.Type == gesture.TypeClick {
				for _, txn := range (*page.walletTransactions).Recent {
					if has == txn.Txn.Hash {
						*page.walletTransaction = &txn
						*c.page = PageTransactionDetails
						return
					}
				}
			}
		}
	}
}
