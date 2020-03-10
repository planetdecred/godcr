package ui

import (
	"fmt"
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/helper"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/ui/values"
)

const (
	balanceTitle         = "Current Total Balance"
	statusTitle          = "Wallet Status"
	stepsTitle           = "Step"
	transactionsTitle    = "Recent Transactions"
	connectedPeersTitle  = "Connected peers count"
	headersFetchedTitle  = "Block header fetched"
	syncingProgressTitle = "Syncing progress"
	latestBlockTitle 	 = "Latest Block"
	noTransaction        = "no transactions"
	onlineStatus         = "Offline"
	syncingStatus        = "Syncing..."
	notSyncedStatus      = "Not Synced"
	syncedStatus         = "Synced"
	fetchingBlockHeaders = "Fetching block headers"
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
	balance     int64
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
		container := layout.Inset{Left: units.ContainerPadding, Right: units.ContainerPadding}
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
			layout.Inset{Top: units.PageMarginTop}.Layout(gtx, func() {
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
			layout.Inset{Bottom: units.ContainerPadding}.Layout(gtx, func() {
				syncStatusColumn(win)
			})
		},
	}

	listContainer := layout.List{Axis: layout.Vertical}
	listContainer.Layout(gtx, len(pageContent), func(i int) {
		layout.UniformInset(units.NoPadding).Layout(gtx, pageContent[i])
	})
}

// recentTransactionsColumn lays out the list of recent transactions.
func recentTransactionsColumn(win *Window) {
	theme := win.theme
	gtx := win.gtx
	walletInfo := win.walletInfo
	recentTransactions := win.walletTransactions.Recent

	var transactionRows []func()
	if len(recentTransactions) > 0 {
		for _, txn := range recentTransactions {
			txnWidgets := transactionWidgets{
				wallet:      theme.Body1(""),
				balance:     txn.Amount,
				mainBalance: theme.Body1(""),
				subBalance:  theme.Caption(""),
				date: theme.Body1(fmt.Sprintf("%v",
					dcrlibwallet.ExtractDateOrTime(txn.Timestamp))),
				status: theme.Body1(helper.TransactionStatus(walletInfo.BestBlockHeight,
					txn.BlockHeight)),
			}
			walletName, err := helper.WalletNameFromID(txn.WalletID, walletInfo.Wallets)
			if err != nil {
				fmt.Printf("%v \n", err.Error())
			} else {
				txnWidgets.wallet.Text = walletName
			}

			transactionRows = append(transactionRows, func() {
				recentTransactionRow(gtx, txnWidgets)
			})
		}
	} else {
		transactionRows = append(transactionRows, func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layout.Flexed(values.EntireSpace, func() {
				layout.Center.Layout(gtx, func() {
					label := theme.Caption(noTransaction)
					label.Color = values.TextGray
					label.Layout(gtx)
				})
			}))
		})
	}

	layout.Inset{Top: units.ColumnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				theme.Caption(transactionsTitle).Layout(gtx)
			}),
			layout.Rigid(func() {
				list := &layout.List{Axis: layout.Vertical}
				list.Layout(gtx, len(transactionRows), func(i int) {
					layout.Inset{Top: units.Padding}.Layout(gtx, transactionRows[i])
				})
			}),
			layout.Rigid(func() {
				if len(transactionRows) > 5 {
					layout.Center.Layout(gtx, func() {
						layout.Inset{Top: units.Padding}.Layout(gtx, func() {
							layout.Stack{}.Layout(gtx,
								layout.Stacked(func() {
									//page.moreButtonCard.Color = values.ButtonGray
									//page.moreButtonCard.Width = values.MoreButtonWidth
									//page.moreButtonCard.Height = values.MoreButtonHeight
									//page.moreButtonCard.Layout(gtx, float32(gtx.Px(units.DefaultButtonRadius)))
								}),
								layout.Expanded(func() {
									layout.Center.Layout(gtx, func() {
										//gtx.Constraints.Width.Min = values.MoreButtonWidth - values.ButtonBorder
										//gtx.Constraints.Height.Max = values.MoreButtonHeight - values.ButtonBorder
										//page.moreButton.Color = values.ButtonGray
										//page.moreButton.Background = values.White
										//page.moreButton.TextSize = units.SyncButtonTextSize
										//page.moreButton.Layout(gtx, page.moreButtonWidget)
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
	margin := layout.UniformInset(units.TransactionsRowMargin)
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func() {
			margin.Layout(gtx, func() {
				layoutBalance(gtx, helper.Balance(txn.balance), txn.mainBalance, txn.subBalance)
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
	uniform := layout.UniformInset(units.Padding)
	layout.Inset{Top: units.ColumnMargin}.Layout(gtx, func() {
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
	layout.Flex{Axis:layout.Vertical}.Layout(win.gtx,
		layout.Rigid(func() {
			latestBlockTitleLabel := win.theme.Body1(latestBlockTitle)
			blockLabel := win.theme.Body1(fmt.Sprintf("%v . %v ago", win.walletInfo.BestBlockHeight,
				helper.LastBlockSync(win.walletInfo.BestBlockTime)))
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
			layout.Flexed(values.EntireSpace, func() {
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
	statusTitleLabel.Color = values.TextGray
	statusLabel := win.theme.Body1(onlineStatus)
	endToEndRow(win.gtx, inset, statusTitleLabel, statusLabel)
}

// syncBoxTitleRow lays out sync status text and sync button.
func syncStatusTextRow(win *Window, inset layout.Inset) {
	gtx := win.gtx
	inset.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func() {
				syncText := ""
				if win.walletInfo.Syncing {
					syncText = syncingStatus
				} else if win.walletInfo.Synced {
					syncText = syncedStatus
				} else {
					syncText = notSyncedStatus
				}
				win.theme.H6(syncText).Layout(gtx)
			}),
			layout.Flexed(values.EntireSpace, func() {
				// stack a button on a card widget to produce a transparent button.
				layout.E.Layout(gtx, func() {
					layout.Stack{}.Layout(gtx,
						layout.Stacked(func() {
							//page.syncButtonCard.Width = values.SyncButtonWidth
							//page.syncButtonCard.Height = values.SyncButtonHeight
							//page.syncButtonCard.Layout(gtx, float32(gtx.Px(units.DefaultButtonRadius)))
						}),
						layout.Expanded(func() {
							layout.Center.Layout(gtx, func() {
								gtx.Constraints.Width.Min = values.SyncButtonWidth - values.ButtonBorder
								gtx.Constraints.Height.Max = values.SyncButtonHeight - values.ButtonBorder
								syncButton := win.outputs.sync
								syncButton.Background = values.White
								syncButton.TextSize = units.SyncButtonTextSize
								syncButton.Layout(gtx, &win.inputs.sync)
							})
						}),
					)
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
	timeLeftLabel := win.theme.Body1(fmt.Sprintf("%v left",
		helper.RemainingSyncTime(win.walletSyncStatus.RemainingTime)))
	endToEndRow(gtx, inset, percentageLabel, timeLeftLabel)
}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func walletSyncRow(win *Window, inset layout.Inset) {
	gtx := win.gtx
	layout.Inset{Top: units.ColumnMargin}.Layout(gtx, func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func() {
				completedSteps := win.walletSyncStatus.Steps
				totalSteps := win.walletSyncStatus.TotalSteps
				completedStepsLabel := win.theme.Caption(fmt.Sprintf("%s %d/%d", stepsTitle, completedSteps, totalSteps))
				completedStepsLabel.Color = values.TextGray
				headersFetchedLabel := win.theme.Body1(fmt.Sprintf("%s. %v%%", fetchingBlockHeaders,
					win.walletSyncStatus.HeadersFetchProgress))
				headersFetchedLabel.Color = values.TextGray
				endToEndRow(win.gtx, inset, completedStepsLabel, headersFetchedLabel)
			}),
			layout.Rigid(func() {
				connectedPeersTitleLabel := win.theme.Caption(connectedPeersTitle)
				connectedPeersTitleLabel.Color = values.TextGray
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
					status := helper.WalletSyncStatus(w, overallBlockHeight)
					progress := helper.WalletSyncProgressTime(w.BlockTimestamp)
					details := syncDetail(win.theme, w.Name, status, blockHeightProgress, progress)
					uniform := layout.UniformInset(units.Padding)
					walletSyncBoxes = append(walletSyncBoxes,
						func() {
							walletSyncBox(win, uniform, details)
						})
				}

				walletSyncList := &layout.List{Axis: layout.Horizontal}
				walletSyncList.Layout(gtx, len(walletSyncBoxes), func(i int) {
					if i == 0 {
						layout.UniformInset(units.NoPadding).Layout(gtx, walletSyncBoxes[i])
					} else {
						layout.Inset{Left: units.ColumnMargin}.Layout(gtx, walletSyncBoxes[i])
					}
				})
			}),
		)
	})
}

// walletSyncBox lays out the wallet syncing details of a single wallet.
func walletSyncBox(win *Window, inset layout.Inset, details walletSyncDetails) {
	gtx := win.gtx
	layout.Inset{Top: units.ColumnMargin}.Layout(gtx, func() {
		layout.Stack{}.Layout(gtx,
			layout.Stacked(func() {
				//page.walletSyncCard.Width = gtx.Px(units.WalletSyncBoxWidthMin)
				//page.walletSyncCard.Height = gtx.Px(units.WalletSyncBoxHeightMin)
				//page.walletSyncCard.Layout(gtx, 0)
			}),
			layout.Stacked(func() {
				uniform := layout.UniformInset(units.SyncBoxPadding)
				uniform.Layout(gtx, func() {
					gtx.Constraints.Width.Min = gtx.Px(units.WalletSyncBoxContentWidth)
					gtx.Constraints.Width.Max = gtx.Constraints.Width.Min
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							endToEndRow(gtx, inset, details.name, details.status)
						}),
						layout.Rigid(func() {
							headersFetchedTitleLabel := win.theme.Caption(headersFetchedTitle)
							headersFetchedTitleLabel.Color = values.TextGray
							endToEndRow(gtx, inset, headersFetchedTitleLabel, details.blockHeaderFetched)
						}),
						layout.Rigid(func() {
							progressTitleLabel := win.theme.Caption(syncingProgressTitle)
							progressTitleLabel.Color = values.TextGray
							endToEndRow(gtx, inset, progressTitleLabel, details.syncingProgress)
						}),
					)
				})
			}),
		)
	})
}

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func layoutBalance(gtx *layout.Context, amount string, main, sub decredmaterial.Label) {
	mainText, subText := helper.BreakBalance(amount)
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
