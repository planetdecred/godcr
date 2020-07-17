package ui

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/ui/values"
	"github.com/raedahgroup/godcr/wallet"
)

const PageOverview = "overview"

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
	cancel,
	viewAllTx,
	connectedPeersInfo,
	noConnectedPeers string
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
	direction   *widget.Icon
	mainBalance decredmaterial.Label
	subBalance  decredmaterial.Label
	date        decredmaterial.Label
	status      decredmaterial.Label
}

type overviewPage struct {
	listContainer, walletSyncList *layout.List
	theme                         *decredmaterial.Theme
	tab                           *decredmaterial.Tabs

	walletInfo           *wallet.MultiWalletInfo
	walletSyncStatus     *wallet.SyncStatus
	walletTransactions   **wallet.Transactions
	walletTransaction    **wallet.Transaction
	toTransactions, sync decredmaterial.Button
	syncedIcon, notSyncedIcon,
	walletStatusIcon *widget.Icon
	syncingIcon          image.Image
	toTransactionDetails []*gesture.Click
	line                 *decredmaterial.Line

	text             overviewPageText
	syncButtonHeight int
	syncButtonWidth  int
	moreButtonWidth  int
	moreButtonHeight int
	gray             color.RGBA
}

func (win *Window) OverviewPage(c pageCommon) layout.Widget {
	pg := overviewPage{
		theme: c.theme,
		tab:   c.navTab,

		walletInfo:         win.walletInfo,
		walletSyncStatus:   win.walletSyncStatus,
		walletTransactions: &win.walletTransactions,
		walletTransaction:  &win.walletTransaction,
		listContainer:      &layout.List{Axis: layout.Vertical},
		walletSyncList:     &layout.List{Axis: layout.Horizontal},
		toTransactions:     c.theme.Button(new(widget.Clickable), "See all"),
		line:               c.theme.Line(),

		syncButtonHeight: 70,
		syncButtonWidth:  145,
		moreButtonWidth:  115,
		moreButtonHeight: 70,

		gray: color.RGBA{137, 151, 165, 255},
	}
	pg.text = overviewPageText{
		balanceTitle:         "Current Total Balance",
		statusTitle:          "Wallet Status",
		stepsTitle:           "Step",
		transactionsTitle:    "Recent Transactions",
		connectedPeersTitle:  "Connected peers count",
		connectedPeersInfo:   "Connected to",
		noConnectedPeers:     "No connected peers",
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
		reconnect:            "Connect",
		disconnect:           "Disconnect",
		cancel:               "Cancel",
		viewAllTx:            "See all",
	}

	pg.toTransactions = c.theme.Button(new(widget.Clickable), pg.text.viewAllTx)
	pg.toTransactions.TextSize = values.TextSize14
	pg.toTransactions.Background = color.RGBA{}
	pg.toTransactions.Color = c.theme.Color.Primary
	pg.sync = c.theme.Button(new(widget.Clickable), pg.text.reconnect)
	pg.toTransactions.Inset = layout.Inset{
		Top: values.MarginPadding10, Bottom: values.MarginPadding0,
		Left: values.MarginPadding0, Right: values.MarginPadding0,
	}

	pg.sync = c.theme.Button(new(widget.Clickable), pg.text.reconnect)
	pg.sync.TextSize = values.TextSize10
	pg.sync.Background = c.theme.Color.Background
	pg.sync.Color = c.theme.Color.Text

	pg.syncedIcon = c.icons.actionCheckCircle
	pg.syncedIcon.Color = c.theme.Color.Success

	pg.syncingIcon = c.icons.syncingIcon

	pg.notSyncedIcon = c.icons.navigationCancel
	pg.notSyncedIcon.Color = c.theme.Color.Danger

	pg.walletStatusIcon = c.icons.imageBrightness1

	pg.line.Color = c.theme.Color.Gray

	return func(gtx C) D {
		pg.Handler(gtx, c)
		return pg.Layout(gtx, c)
	}
}

// Layout lays out the entire content for overview pg.
func (pg *overviewPage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	if c.info.LoadedWallets == 0 {
		return c.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return c.theme.H3(pg.text.noWallet).Layout(gtx)
			})
		})
	}

	walletInfo := pg.walletInfo
	theme := pg.theme

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					mainBalance := theme.H4("")
					subBalance := theme.H6("")
					return pg.layoutBalance(gtx, walletInfo.TotalBalance, mainBalance, subBalance)
				}),
				layout.Rigid(func(gtx C) D {
					return theme.Caption(pg.text.balanceTitle).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return pg.recentTransactionsColumn(gtx, c)
		},
		func(gtx C) D {
			return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
				return pg.syncStatusColumn(gtx)
			})
		},
	}

	return c.Layout(gtx, func(gtx C) D {
		return pg.listContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, pageContent[i])
		})
	})
}

// syncDetail returns a walletSyncDetails object containing data of a single wallet sync box
func (pg *overviewPage) syncDetail(name, status, headersFetched, progress string) walletSyncDetails {
	return walletSyncDetails{
		name:               pg.theme.Caption(name),
		status:             pg.theme.Caption(status),
		blockHeaderFetched: pg.theme.Caption(headersFetched),
		syncingProgress:    pg.theme.Caption(progress),
	}
}

// recentTransactionsColumn lays out the list of recent transactions.
func (pg *overviewPage) recentTransactionsColumn(gtx layout.Context, c pageCommon) layout.Dimensions {
	theme := pg.theme
	var transactionRows []func(gtx C) D

	if len((*pg.walletTransactions).Txs) > 0 {
		pg.updateToTransactionDetailsButtons()

		for index, txn := range (*pg.walletTransactions).Recent {
			txnWidgets := transactionWidgets{
				wallet:      theme.Body1(txn.WalletName),
				balance:     txn.Balance,
				mainBalance: theme.Body1(""),
				subBalance:  theme.Caption(""),
				date:        theme.Body1(txn.DateTime),
				status:      theme.Body1(txn.Status),
			}
			if txn.Txn.Direction == dcrlibwallet.TxDirectionSent {
				txnWidgets.direction = c.icons.contentRemove
				txnWidgets.direction.Color = c.theme.Color.Danger
			} else {
				txnWidgets.direction = c.icons.contentAdd
				txnWidgets.direction.Color = c.theme.Color.Success
			}

			click := pg.toTransactionDetails[index]

			transactionRows = append(transactionRows, func(gtx C) D {
				pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
				click.Add(gtx.Ops)
				return pg.recentTransactionRow(gtx, txnWidgets)
			})
		}
	} else {
		transactionRows = append(transactionRows, func(gtx C) D {
			return pg.centralize(gtx, func(gtx C) D {
				label := theme.Caption(pg.text.noTransaction)
				label.Color = pg.gray
				return label.Layout(gtx)
			})
		})
	}

	return pg.drawlayout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return theme.Caption(pg.text.transactionsTitle).Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				list := &layout.List{Axis: layout.Vertical}
				return pg.centralize(gtx, func(gtx C) D {
					return list.Layout(gtx, len(transactionRows), func(gtx C, i int) D {
						return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, transactionRows[i])
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				pg.line.Width = gtx.Constraints.Max.X
				return pg.line.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.centralize(gtx, func(gtx C) D {
					return pg.toTransactions.Layout(gtx)
				})
			}),
		)
	})
}

func (pg *overviewPage) centralize(gtx layout.Context, content layout.Widget) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, content)
		}),
	)
}

// recentTransactionRow lays out a single row of a recent transaction.
func (pg *overviewPage) recentTransactionRow(gtx layout.Context, txn transactionWidgets) layout.Dimensions {
	margin := layout.UniformInset(values.MarginPadding10)

	dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding50)
			return layout.Inset{Top: values.TextSize12}.Layout(gtx, func(gtx C) D {
				return txn.direction.Layout(gtx, values.TextSize16)
			})
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
			return margin.Layout(gtx, func(gtx C) D {
				return pg.layoutBalance(gtx, txn.balance, txn.mainBalance, txn.subBalance)
			})
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
			return margin.Layout(gtx, func(gtx C) D {
				return txn.wallet.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding100)
			return margin.Layout(gtx, func(gtx C) D {
				return txn.date.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return margin.Layout(gtx, func(gtx C) D {
				return txn.status.Layout(gtx)
			})
		}),
	)
	return dims
}

// syncStatusColumn lays out content for displaying sync status.
func (pg *overviewPage) syncStatusColumn(gtx layout.Context) layout.Dimensions {
	uniform := layout.UniformInset(values.MarginPadding5)
	return pg.drawlayout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.syncBoxTitleRow(gtx, uniform)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.syncStatusTextRow(gtx, uniform)
			}),
			layout.Rigid(func(gtx C) D {
				if pg.walletInfo.Syncing {
					return pg.syncActiveContent(gtx, uniform)
				}
				return pg.syncDormantContent(gtx)
			}),
		)
	})
}

// drawlayout wraps the page tx and sync section in a card layout
func (pg *overviewPage) drawlayout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return decredmaterial.Card{Color: pg.theme.Color.Surface}.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding20).Layout(gtx, body)
	})
}

// syncingContent lays out sync status content when the wallet is syncing
func (pg *overviewPage) syncActiveContent(gtx layout.Context, uniform layout.Inset) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.progressBarRow(gtx, uniform)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.progressStatusRow(gtx, uniform)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.walletSyncRow(gtx, uniform)
		}),
	)
}

// syncDormantContent lays out sync status content when the wallet is synced or not connected
func (pg *overviewPage) syncDormantContent(gtx layout.Context) layout.Dimensions {
	return layout.Inset{Left: values.MarginPadding45}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.TextSize12}.Layout(gtx, func(gtx C) D {
					return pg.blockInfoRow(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.walletInfo.Synced {
					return pg.connectionPeer(gtx)
				}
				latestBlockTitleLabel := pg.theme.Body1(pg.text.noConnectedPeers)
				latestBlockTitleLabel.Color = pg.gray
				return latestBlockTitleLabel.Layout(gtx)
			}),
		)
	})
}

func (pg *overviewPage) blockInfoRow(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			latestBlockTitleLabel := pg.theme.Body1(pg.text.latestBlockTitle)
			latestBlockTitleLabel.Color = pg.gray
			return latestBlockTitleLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding5, Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return pg.theme.Body1(fmt.Sprintf("%v", pg.walletInfo.BestBlockHeight)).Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			pg.walletStatusIcon.Color = pg.gray
			return layout.Inset{Right: values.MarginPadding10, Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
				return pg.walletStatusIcon.Layout(gtx, values.MarginPadding5)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return pg.theme.Body1(fmt.Sprintf("%v", pg.walletInfo.LastSyncTime)).Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			lastSyncedLabel := pg.theme.Body1("ago")
			lastSyncedLabel.Color = pg.gray
			return lastSyncedLabel.Layout(gtx)
		}),
	)
}

func (pg *overviewPage) connectionPeer(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			connectedPeersInfoLabel := pg.theme.Body1(pg.text.connectedPeersInfo)
			connectedPeersInfoLabel.Color = pg.gray
			return connectedPeersInfoLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding5, Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return pg.theme.Body1(fmt.Sprintf("%d", pg.walletSyncStatus.ConnectedPeers)).Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			peersLabel := pg.theme.Body1("peers")
			peersLabel.Color = pg.gray
			return peersLabel.Layout(gtx)
		}),
	)
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func (pg *overviewPage) endToEndRow(gtx layout.Context, inset layout.Inset, leftLabel, rightLabel decredmaterial.Label) layout.Dimensions {
	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return leftLabel.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return rightLabel.Layout(gtx)
				})
			}),
		)
	})
}

// syncBoxTitleRow lays out widgets in the title row inside the sync status box.
func (pg *overviewPage) syncBoxTitleRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	statusTitleLabel := pg.theme.Caption(pg.text.statusTitle)
	statusTitleLabel.Color = pg.theme.Color.Text
	statusLabel := pg.theme.Body1(pg.text.offlineStatus)
	pg.walletStatusIcon.Color = pg.theme.Color.Danger
	if pg.walletInfo.Synced || pg.walletInfo.Syncing {
		statusLabel.Text = pg.text.onlineStatus
		pg.walletStatusIcon.Color = pg.theme.Color.Success
	}

	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return statusTitleLabel.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Top: values.MarginPadding5, Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return pg.walletStatusIcon.Layout(gtx, values.MarginPadding10)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return statusLabel.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}

// syncStatusTextRow lays out sync status text and sync button.
func (pg *overviewPage) syncStatusTextRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	syncStatusLabel := pg.theme.H6(pg.text.notSyncedStatus)
	syncStatusIcon := pg.notSyncedIcon
	if pg.walletInfo.Syncing {
		syncStatusLabel.Text = pg.text.syncingStatus
	} else if pg.walletInfo.Synced {
		syncStatusLabel.Text = pg.text.syncedStatus
		syncStatusIcon = pg.syncedIcon
	}

	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if pg.walletInfo.Syncing {
					return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return pg.theme.ImageIcon(gtx, pg.syncingIcon, 50)
					})
				}
				return layout.Inset{Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
					return syncStatusIcon.Layout(gtx, values.MarginPadding25)
				})
			}),
			layout.Flexed(0.5, func(gtx C) D {
				return layout.W.Layout(gtx, func(gtx C) D {
					return syncStatusLabel.Layout(gtx)
				})
			}),
			layout.Flexed(1, func(gtx C) D {
				// stack a button on a card widget to produce a transparent button.
				return layout.E.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = pg.syncButtonWidth
					gtx.Constraints.Min.X = pg.syncButtonWidth
					gtx.Constraints.Max.Y = pg.syncButtonHeight
					if pg.walletInfo.Synced {
						pg.sync.Text = pg.text.disconnect
					}
					return pg.sync.Layout(gtx)
				})
			}),
		)
	})
}

// syncBoxTitleRow lays out the progress bar.
func (pg *overviewPage) progressBarRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	return inset.Layout(gtx, func(gtx C) D {
		progress := pg.walletSyncStatus.Progress
		p := pg.theme.ProgressBar(int(progress))
		p.Color = pg.theme.Color.Success
		return p.Layout(gtx)
	})
}

// syncBoxTitleRow lays out the progress status when the wallet is syncing.
func (pg *overviewPage) progressStatusRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	timeLeft := pg.walletSyncStatus.RemainingTime
	if timeLeft == "" {
		timeLeft = "0s"
	}

	percentageLabel := pg.theme.Body1(fmt.Sprintf("%v%%", pg.walletSyncStatus.Progress))
	timeLeftLabel := pg.theme.Body1(fmt.Sprintf("%v Left", timeLeft))
	return pg.endToEndRow(gtx, inset, percentageLabel, timeLeftLabel)
}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func (pg *overviewPage) walletSyncRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				completedSteps := pg.walletSyncStatus.Steps
				totalSteps := pg.walletSyncStatus.TotalSteps
				completedStepsLabel := pg.theme.Caption(fmt.Sprintf("%s %d/%d", pg.text.stepsTitle, completedSteps, totalSteps))
				completedStepsLabel.Color = pg.gray
				headersFetchedLabel := pg.theme.Body1(fmt.Sprintf("%s. %v%%", pg.text.fetchingBlockHeaders,
					pg.walletSyncStatus.HeadersFetchProgress))
				headersFetchedLabel.Color = pg.gray
				return pg.endToEndRow(gtx, inset, completedStepsLabel, headersFetchedLabel)
			}),
			layout.Rigid(func(gtx C) D {
				connectedPeersTitleLabel := pg.theme.Caption(pg.text.connectedPeersTitle)
				connectedPeersTitleLabel.Color = pg.gray
				connectedPeersLabel := pg.theme.Body1(fmt.Sprintf("%d", pg.walletSyncStatus.ConnectedPeers))
				return pg.endToEndRow(gtx, inset, connectedPeersTitleLabel, connectedPeersLabel)
			}),
			layout.Rigid(func(gtx C) D {
				var overallBlockHeight int32
				var walletSyncBoxes []func(gtx C) D

				if pg.walletSyncStatus != nil {
					overallBlockHeight = pg.walletSyncStatus.HeadersToFetch
				}

				for i := 0; i < len(pg.walletInfo.Wallets); i++ {
					w := pg.walletInfo.Wallets[i]
					if w.BestBlockHeight > overallBlockHeight {
						overallBlockHeight = w.BestBlockHeight
					}
					blockHeightProgress := fmt.Sprintf("%v of %v", w.BestBlockHeight, overallBlockHeight)
					details := pg.syncDetail(w.Name, w.Status, blockHeightProgress, w.DaysBehind)
					uniform := layout.UniformInset(values.MarginPadding5)
					walletSyncBoxes = append(walletSyncBoxes,
						func(gtx C) D {
							return pg.walletSyncBox(gtx, uniform, details)
						})
				}

				return pg.walletSyncList.Layout(gtx, len(walletSyncBoxes), func(gtx C, i int) D {
					if i == 0 {
						return walletSyncBoxes[i](gtx)
					}
					return layout.Inset{Left: values.MarginPadding30}.Layout(gtx, walletSyncBoxes[i])
				})
			}),
		)
	})
}

// walletSyncBox lays out the wallet syncing details of a single wallet.
func (pg *overviewPage) walletSyncBox(gtx layout.Context, inset layout.Inset, details walletSyncDetails) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Px(values.WalletSyncBoxContentWidth)
		gtx.Constraints.Max.X = gtx.Constraints.Min.X
		return decredmaterial.Card{Color: pg.theme.Color.Background}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.endToEndRow(gtx, inset, details.name, details.status)
				}),
				layout.Rigid(func(gtx C) D {
					headersFetchedTitleLabel := pg.theme.Caption(pg.text.headersFetchedTitle)
					headersFetchedTitleLabel.Color = pg.gray
					return pg.endToEndRow(gtx, inset, headersFetchedTitleLabel, details.blockHeaderFetched)
				}),
				layout.Rigid(func(gtx C) D {
					progressTitleLabel := pg.theme.Caption(pg.text.syncingProgressTitle)
					progressTitleLabel.Color = pg.gray
					return pg.endToEndRow(gtx, inset, progressTitleLabel, details.syncingProgress)
				}),
			)
		})
	})
}

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (pg *overviewPage) layoutBalance(gtx layout.Context, amount string, main, sub decredmaterial.Label) layout.Dimensions {
	mainText, subText := pg.breakBalance(amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			main.Text = mainText
			return main.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			sub.Text = subText
			return sub.Layout(gtx)
		}),
	)
}

// breakBalance takes the balance string and returns it in two slices
func (pg *overviewPage) breakBalance(balance string) (b1, b2 string) {
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

func (pg *overviewPage) updateToTransactionDetailsButtons() {
	recentTxs := (*pg.walletTransactions).Recent
	if len(recentTxs) != len(pg.toTransactionDetails) {
		pg.toTransactionDetails = make([]*gesture.Click, len(recentTxs))
		for i := range recentTxs {
			pg.toTransactionDetails[i] = &gesture.Click{}
		}
	}
}

func (pg *overviewPage) Handler(gtx layout.Context, c pageCommon) {
	if pg.sync.Button.Clicked() {
		if pg.walletInfo.Synced || pg.walletInfo.Syncing {
			c.wallet.CancelSync()
			pg.sync.Text = pg.text.reconnect
		} else {
			c.wallet.StartSync()
			pg.sync.Text = pg.text.cancel
		}
	}
	if pg.toTransactions.Button.Clicked() {
		pg.tab.ChangeTab(4)
	}

	for index, click := range pg.toTransactionDetails {
		for _, e := range click.Events(gtx) {
			if e.Type == gesture.TypeClick {
				txn := (*pg.walletTransactions).Recent[index]
				*pg.walletTransaction = &txn
				*c.page = PageTransactionDetails
				return
			}
		}
	}
}
