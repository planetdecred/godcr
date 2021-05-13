package ui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageOverview = "Overview"

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
	noConnectedPeers,
	showSyncDetails,
	hideSyncDetails string
}

// walletSyncDetails contains sync data for each wallet when a sync
// is in progress.
type walletSyncDetails struct {
	name               decredmaterial.Label
	status             decredmaterial.Label
	blockHeaderFetched decredmaterial.Label
	syncingProgress    decredmaterial.Label
}

type overviewPage struct {
	listContainer, walletSyncList,
	transactionsList *layout.List
	theme *decredmaterial.Theme
	tab   *decredmaterial.Tabs

	walletInfo         *wallet.MultiWalletInfo
	walletSyncStatus   *wallet.SyncStatus
	walletTransactions **wallet.Transactions
	walletTransaction  **wallet.Transaction
	toTransactions     decredmaterial.TextAndIconButton
	sync               decredmaterial.Button
	toggleSyncDetails  decredmaterial.Button
	syncedIcon, notSyncedIcon,
	walletStatusIcon, cachedIcon *widget.Icon
	syncingIcon          *widget.Image
	toTransactionDetails []*gesture.Click

	autoSyncWallet bool

	text                  overviewPageText
	syncButtonHeight      int
	moreButtonWidth       int
	moreButtonHeight      int
	isCheckingLockWL      bool
	syncDetailsVisibility bool
	txnRowHeight          int
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
		walletSyncList:     &layout.List{Axis: layout.Vertical},
		transactionsList:   &layout.List{Axis: layout.Vertical},

		syncButtonHeight: 50,
		moreButtonWidth:  115,
		moreButtonHeight: 70,
		txnRowHeight:     56,

		isCheckingLockWL: false,
		autoSyncWallet:   true,
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
		noTransaction:        "No transactions yet",
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
		viewAllTx:            "See all",
		showSyncDetails:      "Show details",
		hideSyncDetails:      "Hide details",
	}

	pg.toTransactions = c.theme.TextAndIconButton(new(widget.Clickable), pg.text.viewAllTx, c.icons.navigationArrowForward)
	pg.toTransactions.Color = c.theme.Color.Primary
	pg.toTransactions.BackgroundColor = c.theme.Color.Surface

	pg.sync = c.theme.Button(new(widget.Clickable), pg.text.reconnect)
	pg.sync.TextSize = values.TextSize10
	pg.sync.Background = color.NRGBA{}
	pg.sync.Color = c.theme.Color.Text

	pg.toggleSyncDetails = c.theme.Button(new(widget.Clickable), pg.text.showSyncDetails)
	pg.toggleSyncDetails.TextSize = values.TextSize16
	pg.toggleSyncDetails.Background = color.NRGBA{}
	pg.toggleSyncDetails.Color = c.theme.Color.Primary
	pg.toggleSyncDetails.Inset = layout.Inset{}

	pg.syncedIcon = c.icons.actionCheckCircle
	pg.syncedIcon.Color = c.theme.Color.Success

	pg.syncingIcon = c.icons.syncingIcon
	pg.syncingIcon.Scale = 1

	pg.notSyncedIcon = c.icons.navigationCancel
	pg.notSyncedIcon.Color = c.theme.Color.Danger

	pg.walletStatusIcon = c.icons.imageBrightness1
	pg.cachedIcon = c.icons.cached

	return func(gtx C) D {
		pg.Handler(gtx, c, win)
		return pg.Layout(gtx, c)
	}
}

// Layout lays out the entire content for overview pg.
func (pg *overviewPage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	if c.info.LoadedWallets == 0 {
		return c.Layout(gtx, func(gtx C) D {
			return c.UniformPadding(gtx, func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					return c.theme.H3(pg.text.noWallet).Layout(gtx)
				})
			})
		})
	}

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.recentTransactionsSection(gtx, c)
		},
		func(gtx C) D {
			return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
				return pg.syncStatusSection(gtx)
			})
		},
	}

	return c.Layout(gtx, func(gtx C) D {
		return c.UniformPadding(gtx, func(gtx C) D {
			return pg.listContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, pageContent[i])
			})
		})
	})
}

// syncDetail returns a walletSyncDetails object containing data of a single wallet sync box
func (pg *overviewPage) syncDetail(name, status, headersFetched, progress string) walletSyncDetails {
	return walletSyncDetails{
		name:               pg.theme.Body1(name),
		status:             pg.theme.Body2(status),
		blockHeaderFetched: pg.theme.Body1(headersFetched),
		syncingProgress:    pg.theme.Body1(progress),
	}
}

// recentTransactionsSection lays out the list of recent transactions.
func (pg *overviewPage) recentTransactionsSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	var recentTransactions []wallet.Transaction
	if len((*pg.walletTransactions).Txs) > 0 {
		recentTransactions = (*pg.walletTransactions).Recent
		if len(recentTransactions) != len(pg.toTransactionDetails) {
			pg.toTransactionDetails = createClickGestures(len(recentTransactions))
		}
	}

	return pg.theme.Card().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		padding := values.MarginPadding15
		return Container{layout.Inset{Top: padding, Bottom: padding}}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					title := pg.theme.Body2(pg.text.transactionsTitle)
					title.Color = pg.theme.Color.Gray3
					return pg.titleRow(gtx, title.Layout, pg.toTransactions.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return pg.theme.Separator().Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if len((*pg.walletTransactions).Txs) == 0 {
						message := pg.theme.Body1(pg.text.noTransaction)
						message.Color = pg.theme.Color.Gray2
						return Container{layout.Inset{
							Left:   values.MarginPadding16,
							Bottom: values.MarginPadding3,
							Top:    values.MarginPadding18,
						}}.Layout(gtx, func(gtx C) D {
							return message.Layout(gtx)
						})
					}

					return pg.transactionsList.Layout(gtx, len(recentTransactions), func(gtx C, i int) D {
						click := pg.toTransactionDetails[i]
						pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
						click.Add(gtx.Ops)
						var row = TransactionRow{
							transaction: recentTransactions[i],
							index:       i,
							showBadge:   showLabel(recentTransactions),
						}
						return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
							return transactionRow(gtx, common, row)
						})
					})
				}),
			)
		})
	})
}

// syncStatusSection lays out content for displaying sync status.
func (pg *overviewPage) syncStatusSection(gtx layout.Context) layout.Dimensions {
	uniform := layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPadding5}
	return pg.theme.Card().Layout(gtx, func(gtx C) D {
		return Container{layout.Inset{
			Top:    values.MarginPadding15,
			Bottom: values.MarginPadding16,
		}}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.syncBoxTitleRow(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, pg.theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return Container{layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return pg.syncStatusIcon(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return pg.syncStatusTextRow(gtx, uniform)
									}),
									layout.Rigid(func(gtx C) D {
										if !pg.walletInfo.Syncing {
											return pg.syncDormantContent(gtx, uniform)
										}
										return layout.Dimensions{}
									}),
									layout.Rigid(func(gtx C) D {
										if pg.walletInfo.Syncing {
											return pg.progressBarRow(gtx, uniform)
										}
										return layout.Dimensions{}
									}),
									layout.Rigid(func(gtx C) D {
										if pg.walletInfo.Syncing {
											return pg.progressStatusRow(gtx, uniform)
										}
										return layout.Dimensions{}
									}),
								)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					if pg.walletInfo.Syncing {
						return pg.theme.Separator().Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.walletInfo.Syncing && pg.syncDetailsVisibility {
						return Container{layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
							return pg.walletSyncRow(gtx, uniform)
						})
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.walletInfo.Syncing && pg.syncDetailsVisibility {
						return pg.theme.Separator().Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.walletInfo.Syncing {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Inset{Top: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
							return pg.toggleSyncDetails.Layout(gtx)
						})
					}
					return layout.Dimensions{}
				}),
			)
		})
	})
}

func (pg overviewPage) titleRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	titlePadding := values.MarginPadding15
	return Container{layout.Inset{
		Left:   titlePadding,
		Right:  titlePadding,
		Bottom: titlePadding,
	}}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return leftWidget(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return rightWidget(gtx)
			}),
		)
	})
}

// syncDormantContent lays out sync status content when the wallet is synced or not connected
func (pg *overviewPage) syncDormantContent(gtx layout.Context, uniform layout.Inset) layout.Dimensions {
	return uniform.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
					return pg.blockInfoRow(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.walletInfo.Synced {
					return pg.connectionPeer(gtx)
				}
				latestBlockTitleLabel := pg.theme.Body1(pg.text.noConnectedPeers)
				latestBlockTitleLabel.Color = pg.theme.Color.Gray
				return latestBlockTitleLabel.Layout(gtx)
			}),
		)
	})
}

func (pg *overviewPage) blockInfoRow(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			latestBlockTitleLabel := pg.theme.Body1(pg.text.latestBlockTitle)
			latestBlockTitleLabel.Color = pg.theme.Color.Gray
			return latestBlockTitleLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding5, Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return pg.theme.Body1(fmt.Sprintf("%v", pg.walletInfo.BestBlockHeight)).Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			pg.walletStatusIcon.Color = pg.theme.Color.Gray
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
			lastSyncedLabel.Color = pg.theme.Color.Gray
			return lastSyncedLabel.Layout(gtx)
		}),
	)
}

func (pg *overviewPage) connectionPeer(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			connectedPeersInfoLabel := pg.theme.Body1(pg.text.connectedPeersInfo)
			connectedPeersInfoLabel.Color = pg.theme.Color.Gray
			return connectedPeersInfoLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding5, Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return pg.theme.Body1(fmt.Sprintf("%d", pg.walletSyncStatus.ConnectedPeers)).Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			peersLabel := pg.theme.Body1("peers")
			peersLabel.Color = pg.theme.Color.Gray
			return peersLabel.Layout(gtx)
		}),
	)
}

// syncBoxTitleRow lays out widgets in the title row inside the sync status box.
func (pg *overviewPage) syncBoxTitleRow(gtx layout.Context) layout.Dimensions {
	title := pg.theme.Body2(pg.text.statusTitle)
	title.Color = pg.theme.Color.Gray3
	statusLabel := pg.theme.Body1(pg.text.offlineStatus)
	pg.walletStatusIcon.Color = pg.theme.Color.Danger
	if pg.walletInfo.Synced || pg.walletInfo.Syncing {
		statusLabel.Text = pg.text.onlineStatus
		pg.walletStatusIcon.Color = pg.theme.Color.Success
	}

	syncStatus := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding4, Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
					return pg.walletStatusIcon.Layout(gtx, values.MarginPadding14)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return statusLabel.Layout(gtx)
			}),
		)
	}
	return pg.titleRow(gtx, title.Layout, syncStatus)
}

// syncStatusTextRow lays out sync status text and sync button.
func (pg *overviewPage) syncStatusTextRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	syncStatusLabel := pg.theme.H6(pg.text.notSyncedStatus)
	if pg.walletInfo.Syncing {
		syncStatusLabel.Text = pg.text.syncingStatus
	} else if pg.walletInfo.Synced {
		syncStatusLabel.Text = pg.text.syncedStatus
	}

	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Flexed(1, syncStatusLabel.Layout),
			layout.Rigid(func(gtx C) D {
				// stack a button on a card widget to produce a transparent button.
				return layout.E.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Max.Y = pg.syncButtonHeight
					border := widget.Border{Color: pg.theme.Color.Hint, CornerRadius: values.MarginPadding10, Width: values.MarginPadding1}
					return border.Layout(gtx, func(gtx C) D {
						pg.sync.Inset = layout.Inset{
							Top:    values.MarginPadding5,
							Bottom: values.MarginPadding5,
							Left:   values.MarginPadding10,
							Right:  values.MarginPadding10,
						}
						pg.sync.CornerRadius = values.MarginPadding10

						if pg.sync.Text == pg.text.reconnect {
							pg.sync.Inset.Left = values.MarginPadding25
							layout.Inset{Top: values.MarginPadding4, Left: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
								pg.cachedIcon.Color = pg.theme.Color.Gray
								return pg.cachedIcon.Layout(gtx, values.TextSize14)
							})
						}

						return pg.sync.Layout(gtx)
					})
				})
			}),
		)
	})
}

func (pg *overviewPage) syncStatusIcon(gtx layout.Context) layout.Dimensions {
	syncStatusIcon := pg.notSyncedIcon
	if pg.walletInfo.Synced {
		syncStatusIcon = pg.syncedIcon
	}
	i := layout.Inset{Right: values.MarginPadding16, Top: values.MarginPadding9}
	if pg.walletInfo.Syncing {
		return i.Layout(gtx, func(gtx C) D {
			return pg.syncingIcon.Layout(gtx)
		})
	}
	return i.Layout(gtx, func(gtx C) D {
		return syncStatusIcon.Layout(gtx, values.MarginPadding24)
	})
}

// progressBarRow lays out the progress bar.
func (pg *overviewPage) progressBarRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	return inset.Layout(gtx, func(gtx C) D {
		progress := pg.walletSyncStatus.Progress
		p := pg.theme.ProgressBar(int(progress))
		p.Height = values.MarginPadding8
		p.Radius = values.MarginPadding4
		p.Color = pg.theme.Color.Success
		return p.Layout(gtx)
	})
}

// progressStatusRow lays out the progress status when the wallet is syncing.
func (pg *overviewPage) progressStatusRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	timeLeft := pg.walletSyncStatus.RemainingTime
	if timeLeft == "" {
		timeLeft = "0s"
	}

	percentageLabel := pg.theme.Body1(fmt.Sprintf("%v%%", pg.walletSyncStatus.Progress))
	timeLeftLabel := pg.theme.Body1(fmt.Sprintf("%v Left", timeLeft))
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return endToEndRow(gtx, percentageLabel.Layout, timeLeftLabel.Layout)
	})

}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func (pg *overviewPage) walletSyncRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				totalSteps := pg.walletSyncStatus.TotalSteps
				completedSteps := pg.theme.Body2(fmt.Sprintf("%s %d/%d", pg.text.stepsTitle,
					pg.walletSyncStatus.Steps, totalSteps))
				completedSteps.Color = pg.theme.Color.Gray
				headersFetched := pg.theme.Body1(fmt.Sprintf("%s  ·  %v%%", pg.text.fetchingBlockHeaders,
					pg.walletSyncStatus.HeadersFetchProgress))
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return endToEndRow(gtx, completedSteps.Layout, headersFetched.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				connectedPeersTitleLabel := pg.theme.Body2(pg.text.connectedPeersTitle)
				connectedPeersTitleLabel.Color = pg.theme.Color.Gray
				connectedPeersLabel := pg.theme.Body1(fmt.Sprintf("%d", pg.walletSyncStatus.ConnectedPeers))
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return endToEndRow(gtx, connectedPeersTitleLabel.Layout, connectedPeersLabel.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				var overallBlockHeight int32
				var walletSyncBoxes []layout.Widget

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
					return walletSyncBoxes[i](gtx)
				})
			}),
		)
	})
}

// walletSyncBox lays out the wallet syncing details of a single wallet.
func (pg *overviewPage) walletSyncBox(gtx layout.Context, inset layout.Inset, details walletSyncDetails) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		card := pg.theme.Card()
		card.Color = pg.theme.Color.LightGray
		return card.Layout(gtx, func(gtx C) D {
			return Container{
				layout.UniformInset(values.MarginPadding16),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return endToEndRow(gtx, details.name.Layout, details.status.Layout)
						})
					}),
					layout.Rigid(func(gtx C) D {
						headersFetchedTitleLabel := pg.theme.Body2(pg.text.headersFetchedTitle)
						headersFetchedTitleLabel.Color = pg.theme.Color.Gray
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return endToEndRow(gtx, headersFetchedTitleLabel.Layout, details.blockHeaderFetched.Layout)
						})
					}),
					layout.Rigid(func(gtx C) D {
						progressTitleLabel := pg.theme.Body2(pg.text.syncingProgressTitle)
						progressTitleLabel.Color = pg.theme.Color.Gray
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return endToEndRow(gtx, progressTitleLabel.Layout, details.syncingProgress.Layout)
						})
					}),
				)
			})
		})
	})
}

func (pg *overviewPage) Handler(eq event.Queue, c pageCommon, win *Window) {

	if win.wallet != nil {
		isDarkModeOn := win.wallet.ReadBoolConfigValueForKey("isDarkModeOn")
		if isDarkModeOn != win.theme.DarkMode {
			win.theme.SwitchDarkMode(isDarkModeOn)
			win.reloadPage(c)
		}
	}

	if pg.walletInfo.Synced {
		pg.sync.Text = pg.text.disconnect
	}

	if pg.autoSyncWallet && !pg.walletInfo.Synced {
		walletsLocked := getLockedWallets(c.wallet.AllWallets())
		if len(walletsLocked) == 0 {
			c.wallet.StartSync()
			pg.sync.Text = pg.text.cancel
			pg.autoSyncWallet = false
		}
	}

	if !pg.isCheckingLockWL {
		if lockedWallets := getLockedWallets(c.wallet.AllWallets()); len(lockedWallets) > 0 {
			showWalletUnlockModal(c, lockedWallets)
		}
		pg.isCheckingLockWL = true
	}

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
		c.changePage(PageTransactions)
	}

	for index, click := range pg.toTransactionDetails {
		for _, e := range click.Events(eq) {
			if e.Type == gesture.TypeClick {
				txn := (*pg.walletTransactions).Recent[index]
				*pg.walletTransaction = &txn

				c.setReturnPage(PageOverview)
				c.changePage(PageTransactionDetails)
				return
			}
		}
	}

	if pg.toggleSyncDetails.Button.Clicked() {
		pg.syncDetailsVisibility = !pg.syncDetailsVisibility
		if pg.syncDetailsVisibility {
			pg.toggleSyncDetails.Text = pg.text.hideSyncDetails
		} else {
			pg.toggleSyncDetails.Text = pg.text.showSyncDetails
		}
	}
}

func showWalletUnlockModal(c pageCommon, lockedWallets []*dcrlibwallet.Wallet) {
	go func() {
		c.modalReceiver <- &modalLoad{
			template: UnlockWalletRestoreTemplate,
			title:    "Unlock to resume restoration",
			confirm: func(pass string) {
				err := c.wallet.UnlockWallet(lockedWallets[0].ID, []byte(pass))
				if err != nil {
					errText := err.Error()
					if err.Error() == "invalid_passphrase" {
						errText = "Invalid passphrase"
					}
					c.notify(errText, false)
				} else {
					c.closeModal()
				}
			},
			confirmText: "Unlock",
		}
	}()
}
