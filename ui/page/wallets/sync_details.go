package wallets

import (
	"fmt"
	"time"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

func (pg *AppOverviewPage) initSyncDetailsWidgets() {
	pg.walletSyncList = &layout.List{Axis: layout.Vertical}
	pg.syncClickable = pg.Theme.NewClickable(true)
	pg.cachedIcon = decredmaterial.NewIcon(pg.Theme.Icons.Cached)
}

// syncDetail returns a walletSyncDetails object containing data of a single wallet sync box
func (pg *AppOverviewPage) syncDetail(status, headersFetched, progress string) walletSyncDetails {
	return walletSyncDetails{
		status:             pg.Theme.Body2(status),
		blockHeaderFetched: pg.Theme.Body1(headersFetched),
		syncingProgress:    pg.Theme.Body1(values.StringF(values.StrSyncingProgressStat, progress)),
	}
}

func (pg *AppOverviewPage) connectionPeer(gtx C) D {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			connectedPeersInfoLabel := pg.Theme.Body1(values.String(values.StrConnectedTo))
			connectedPeersInfoLabel.Color = pg.Theme.Color.GrayText2
			return connectedPeersInfoLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			connectedPeers := pg.WL.MultiWallet.ConnectedPeers()
			return layout.Inset{Left: values.MarginPadding5, Right: values.MarginPadding5}.Layout(gtx, pg.Theme.Body1(fmt.Sprintf("%d", connectedPeers)).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			peersLabel := pg.Theme.Body1(values.String(values.StrPeers))
			peersLabel.Color = pg.Theme.Color.GrayText2
			return peersLabel.Layout(gtx)
		}),
	)
}

// syncStatusTextRow lays out sync status text and sync button.
func (pg *AppOverviewPage) syncStatusTextRow(gtx C, inset layout.Inset) D {
	syncing, rescanning := pg.WL.MultiWallet.IsSyncing(), pg.WL.MultiWallet.IsRescanning()

	syncStatusLabel := pg.Theme.Label(values.TextSize16, values.String(values.StrWalletNotSynced))
	if pg.WL.MultiWallet.IsSyncing() {
		syncStatusLabel.Text = values.String(values.StrSyncingState)
	} else if pg.WL.MultiWallet.IsRescanning() {
		syncStatusLabel.Text = values.String(values.StrRescanningBlocks)
	} else if pg.WL.MultiWallet.IsSynced() {
		syncStatusLabel.Text = values.String(values.StrSynced)
	}

	var children []layout.FlexChild
	children = append(children, layout.Rigid(syncStatusLabel.Layout))

	if syncing || rescanning {
		children = append(children, layout.Flexed(1, pg.progressBarRow))
	}

	children = append(children, layout.Flexed(1, func(gtx C) D {
		return layout.E.Layout(gtx, func(gtx C) D {
			// Set gxt to Disabled (Sets Queue to nil) if syncClickable state is disabled, prevents double click.
			if !pg.syncClickable.Enabled() {
				gtx = pg.syncClickable.SetEnabled(false, &gtx)
			}
			return decredmaterial.LinearLayout{
				Width:     decredmaterial.WrapContent,
				Height:    decredmaterial.WrapContent,
				Clickable: pg.syncClickable,
				Direction: layout.Center,
				Alignment: layout.Middle,
				Border:    decredmaterial.Border{Color: pg.Theme.Color.Gray2, Width: values.MarginPadding1, Radius: decredmaterial.Radius(10)},
				Padding:   layout.Inset{Top: values.MarginPadding3, Bottom: values.MarginPadding3, Left: values.MarginPadding8, Right: values.MarginPadding8},
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if pg.WL.MultiWallet.IsConnectedToDecredNetwork() {
						return D{}
					}

					return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
						pg.cachedIcon.Color = pg.Theme.Color.Gray1
						return pg.cachedIcon.Layout(gtx, values.MarginPadding16)
					})
				}),
				layout.Rigid(func(gtx C) D {
					sync := pg.Theme.Label(values.TextSize14, values.String(values.StrReconnect))
					sync.TextSize = values.TextSize14
					sync.Color = pg.Theme.Color.Text
					if pg.WL.MultiWallet.IsRescanning() {
						sync.Text = values.String(values.StrCancel)
					} else if pg.WL.MultiWallet.IsConnectedToDecredNetwork() {
						sync.Text = values.String(values.StrDisconnect)
					} else {
						sync.Text = values.String(values.StrReconnect)
					}

					return sync.Layout(gtx)
				}),
			)
		})
	}))

	return layout.Flex{
		Axis:      layout.Horizontal,
		Spacing:   layout.SpaceBetween,
		Alignment: layout.Middle,
	}.Layout(gtx, children...)
}

// progressBarRow lays out the progress bar.
func (pg *AppOverviewPage) progressBarRow(gtx C) D {
	return layout.Inset{Left: values.MarginPadding40, Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
		progress, _ := pg.progressStatusDetails()

		p := pg.Theme.ProgressBar(progress)
		p.Height = values.MarginPadding16
		p.Radius = decredmaterial.Radius(4)
		p.Color = pg.Theme.Color.Success
		p.TrackColor = pg.Theme.Color.Gray2

		progressTitleLabel := pg.Theme.Label(values.TextSize14, fmt.Sprintf("%v%%", progress))
		progressTitleLabel.Color = pg.Theme.Color.InvText
		return p.TextLayout(gtx, progressTitleLabel.Layout)
	})
}

// progressStatusRow lays out the progress status when the wallet is syncing.
func (pg *AppOverviewPage) progressStatusDetails() (int, string) {
	timeLeft := pg.remainingSyncTime
	progress := pg.syncProgress
	rescanUpdate := pg.rescanUpdate
	if rescanUpdate != nil && rescanUpdate.ProgressReport != nil {
		progress = int(rescanUpdate.ProgressReport.RescanProgress)
		timeLeft = components.TimeFormat(int(rescanUpdate.ProgressReport.RescanTimeRemaining), true)
	}

	timeLeftLabel := values.StringF(values.StrTimeLeft, timeLeft)
	if progress == 0 {
		timeLeftLabel = values.String(values.StrLoading)
	}

	return progress, timeLeftLabel
}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func (pg *AppOverviewPage) walletSyncRow(gtx C, inset layout.Inset) D {
	col := pg.Theme.Color.GrayText2
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				completedSteps := pg.Theme.Body2(values.StringF(values.StrSyncSteps, pg.syncStep))
				completedSteps.Color = col

				headersFetched := pg.Theme.Body1(values.StringF(values.StrFetchingBlockHeaders, pg.stepFetchProgress))
				if pg.syncStep == wallet.AddressDiscoveryStep {
					headersFetched.Text = values.StringF(values.StrDiscoveringWalletAddress, pg.stepFetchProgress)
				} else if pg.syncStep == wallet.RescanHeadersStep {
					headersFetched.Text = values.StringF(values.StrRescanningHeaders, pg.stepFetchProgress)
				}

				return inset.Layout(gtx, func(gtx C) D {
					return components.EndToEndRow(gtx, completedSteps.Layout, headersFetched.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				connectedPeersTitleLabel := pg.Theme.Body2(values.String(values.StrConnectedPeersCount))
				connectedPeersTitleLabel.Color = col
				connectedPeers := pg.WL.MultiWallet.ConnectedPeers()
				connectedPeersLabel := pg.Theme.Body1(fmt.Sprintf("%d", connectedPeers))
				return inset.Layout(gtx, func(gtx C) D {
					return components.EndToEndRow(gtx, connectedPeersTitleLabel.Layout, connectedPeersLabel.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				_, timeLeft := pg.progressStatusDetails()

				remainingSyncTime := pg.Theme.Body2("Sync completion time")
				remainingSyncTime.Color = col

				timeLeftLabel := pg.Theme.Body2(timeLeft)

				return inset.Layout(gtx, func(gtx C) D {
					return components.EndToEndRow(gtx, remainingSyncTime.Layout, timeLeftLabel.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				currentSeconds := time.Now().UnixNano() / int64(time.Second)
				w := pg.WL.SelectedWallet.Wallet

				status := values.String(values.StrSyncingState)
				if pg.WL.SelectedWallet.Wallet.IsWaiting() {
					status = values.String(values.StrWaitingState)
				}

				blockHeightProgress := values.StringF(values.StrBlockHeaderFetchedCount, w.GetBestBlock(), pg.headersToFetchOrScan)
				daysBehind := components.TimeFormat(int(currentSeconds-w.GetBestBlockTimeStamp()), true)
				details := pg.syncDetail(status, blockHeightProgress, daysBehind)
				uniform := layout.UniformInset(values.MarginPadding5)

				return pg.walletSyncBox(gtx, uniform, details)
			}),
		)
	})
}

// walletSyncBox lays out the wallet syncing details of a single wallet.
func (pg *AppOverviewPage) walletSyncBox(gtx C, inset layout.Inset, details walletSyncDetails) D {
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		card := pg.Theme.Card()
		card.Color = pg.Theme.Color.Gray4

		col := pg.Theme.Color.GrayText2
		return card.Layout(gtx, func(gtx C) D {
			return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						status := pg.Theme.Body2(values.String(values.StrStatus))
						status.Color = col
						return inset.Layout(gtx, func(gtx C) D {
							return components.EndToEndRow(gtx, status.Layout, details.status.Layout)
						})
					}),
					layout.Rigid(func(gtx C) D {
						headersFetchedTitleLabel := pg.Theme.Body2(values.String(values.StrBlockHeaderFetched))
						headersFetchedTitleLabel.Color = col
						return inset.Layout(gtx, func(gtx C) D {
							return components.EndToEndRow(gtx, headersFetchedTitleLabel.Layout, details.blockHeaderFetched.Layout)
						})
					}),
					layout.Rigid(func(gtx C) D {
						progressTitleLabel := pg.Theme.Body2(values.String(values.StrSyncingProgress))
						progressTitleLabel.Color = col
						return inset.Layout(gtx, func(gtx C) D {
							return components.EndToEndRow(gtx, progressTitleLabel.Layout, details.syncingProgress.Layout)
						})
					}),
				)
			})
		})
	})
}

func (pg *AppOverviewPage) rescanDetailsLayout(gtx C, inset layout.Inset) D {
	rescanUpdate := pg.rescanUpdate
	if rescanUpdate == nil {
		return D{}
	}
	wal := pg.WL.MultiWallet.WalletWithID(rescanUpdate.WalletID)
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		card := pg.Theme.Card()
		card.Color = pg.Theme.Color.Gray4
		return card.Layout(gtx, func(gtx C) D {
			return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return inset.Layout(gtx, func(gtx C) D {
							return pg.Theme.Body1(wal.Name).Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						headersFetchedTitleLabel := pg.Theme.Body2(values.String(values.StrBlocksScanned))
						headersFetchedTitleLabel.Color = pg.Theme.Color.GrayText2

						blocksScannedLabel := pg.Theme.Body1(fmt.Sprint(rescanUpdate.ProgressReport.CurrentRescanHeight))
						return inset.Layout(gtx, func(gtx C) D {
							return components.EndToEndRow(gtx, headersFetchedTitleLabel.Layout, blocksScannedLabel.Layout)
						})
					}),
					layout.Rigid(func(gtx C) D {
						progressTitleLabel := pg.Theme.Body2(values.String(values.StrSyncingProgress))
						progressTitleLabel.Color = pg.Theme.Color.GrayText2

						rescanProgress := values.StringF(values.StrBlocksLeft, rescanUpdate.ProgressReport.TotalHeadersToScan-rescanUpdate.ProgressReport.CurrentRescanHeight)
						blocksScannedLabel := pg.Theme.Body1(rescanProgress)
						return inset.Layout(gtx, func(gtx C) D {
							return components.EndToEndRow(gtx, progressTitleLabel.Layout, blocksScannedLabel.Layout)
						})
					}),
				)
			})
		})
	})
}
