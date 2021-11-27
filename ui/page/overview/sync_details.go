package overview

import (
	"fmt"
	"time"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *AppOverviewPage) initSyncDetailsWidgets() {
	pg.walletSyncList = &layout.List{Axis: layout.Vertical}
	pg.syncClickable = pg.Theme.NewClickable(true)
	pg.cachedIcon = decredmaterial.NewIcon(pg.Icons.Cached)

	pg.sync = pg.Theme.Label(values.MarginPadding14, values.String(values.StrReconnect))
	pg.sync.TextSize = values.TextSize14
	pg.sync.Color = pg.Theme.Color.Text
}

// syncDetail returns a walletSyncDetails object containing data of a single wallet sync box
func (pg *AppOverviewPage) syncDetail(name, status, headersFetched, progress string) walletSyncDetails {
	return walletSyncDetails{
		name:               pg.Theme.Body1(name),
		status:             pg.Theme.Body2(status),
		blockHeaderFetched: pg.Theme.Body1(headersFetched),
		syncingProgress:    pg.Theme.Body1(progress),
	}
}

func (pg *AppOverviewPage) connectionPeer(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			connectedPeersInfoLabel := pg.Theme.Body1(values.String(values.StrConnectedTo))
			connectedPeersInfoLabel.Color = pg.Theme.Color.GrayText2
			return connectedPeersInfoLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding5, Right: values.MarginPadding5}.Layout(gtx, pg.Theme.Body1(fmt.Sprintf("%d", pg.connectedPeers)).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			peersLabel := pg.Theme.Body1("peers")
			peersLabel.Color = pg.Theme.Color.GrayText2
			return peersLabel.Layout(gtx)
		}),
	)
}

// syncStatusTextRow lays out sync status text and sync button.
func (pg *AppOverviewPage) syncStatusTextRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	syncStatusLabel := pg.Theme.H6(values.String(values.StrWalletNotSynced))
	if pg.walletSyncing {
		syncStatusLabel.Text = values.String(values.StrSyncingState)
	} else if pg.rescanningBlocks {
		syncStatusLabel.Text = "Rescanning blocks"
	} else if pg.walletSynced {
		syncStatusLabel.Text = values.String(values.StrSynced)
	}

	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Flexed(1, syncStatusLabel.Layout),
			layout.Rigid(func(gtx C) D {
				return decredmaterial.LinearLayout{
					Width:     decredmaterial.WrapContent,
					Height:    decredmaterial.WrapContent,
					Clickable: pg.syncClickable,
					Direction: layout.Center,
					Alignment: layout.Middle,
					Border:    decredmaterial.Border{Color: pg.Theme.Color.Gray1, Width: values.MarginPadding1, Radius: decredmaterial.Radius(10)},
					Padding:   layout.Inset{Top: values.MarginPadding3, Bottom: values.MarginPadding3, Left: values.MarginPadding8, Right: values.MarginPadding8},
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if pg.isConnnected {
							return D{}
						}

						return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
							pg.cachedIcon.Color = pg.Theme.Color.Gray1
							return pg.cachedIcon.Layout(gtx, values.MarginPadding16)
						})
					}),
					layout.Rigid(func(gtx C) D {
						if pg.rescanningBlocks {
							pg.sync.Text = values.String(values.StrCancel)
						} else if pg.isConnnected {
							pg.sync.Text = values.String(values.StrDisconnect)
						} else {
							pg.sync.Text = values.String(values.StrReconnect)
						}

						return pg.sync.Layout(gtx)
					}),
				)
			}),
		)
	})
}

// progressBarRow lays out the progress bar.
func (pg *AppOverviewPage) progressBarRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	return inset.Layout(gtx, func(gtx C) D {
		progress := pg.syncProgress
		rescanUpdate := pg.rescanUpdate
		if rescanUpdate != nil && rescanUpdate.ProgressReport != nil {
			progress = int(rescanUpdate.ProgressReport.RescanProgress)
		}
		p := pg.Theme.ProgressBar(progress)
		p.Height = values.MarginPadding8
		p.Radius = decredmaterial.Radius(values.MarginPadding4.V)
		p.Color = pg.Theme.Color.Success
		p.TrackColor = pg.Theme.Color.Gray2
		return p.Layout(gtx)
	})
}

// progressStatusRow lays out the progress status when the wallet is syncing.
func (pg *AppOverviewPage) progressStatusRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	timeLeft := pg.remainingSyncTime
	progress := pg.syncProgress
	rescanUpdate := pg.rescanUpdate
	if rescanUpdate != nil && rescanUpdate.ProgressReport != nil {
		progress = int(rescanUpdate.ProgressReport.RescanProgress)
		timeLeft = components.TimeFormat(int(rescanUpdate.ProgressReport.RescanTimeRemaining), true)
	}

	percentageLabel := pg.Theme.Body1(fmt.Sprintf("%v%%", progress))
	timeLeftLabel := pg.Theme.Body1(fmt.Sprintf("%v left", timeLeft))
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return components.EndToEndRow(gtx, percentageLabel.Layout, timeLeftLabel.Layout)
	})

}

//	walletSyncRow layouts a list of wallet sync boxes horizontally.
func (pg *AppOverviewPage) walletSyncRow(gtx layout.Context, inset layout.Inset) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				completedSteps := pg.Theme.Body2(values.StringF(values.StrSyncSteps, pg.syncStep))
				completedSteps.Color = pg.Theme.Color.GrayText2
				headersFetched := pg.Theme.Body1(values.StringF(values.StrFetchingBlockHeaders, pg.headerFetchProgress))
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return components.EndToEndRow(gtx, completedSteps.Layout, headersFetched.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				connectedPeersTitleLabel := pg.Theme.Body2(values.String(values.StrConnectedPeersCount))
				connectedPeersTitleLabel.Color = pg.Theme.Color.GrayText2
				connectedPeersLabel := pg.Theme.Body1(fmt.Sprintf("%d", pg.connectedPeers))
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return components.EndToEndRow(gtx, connectedPeersTitleLabel.Layout, connectedPeersLabel.Layout)
				})
			}),
			layout.Rigid(func(gtx C) D {
				var walletSyncBoxes []layout.Widget

				currentSeconds := time.Now().UnixNano() / int64(time.Second)
				for i := 0; i < len(pg.allWallets); i++ {
					w := pg.allWallets[i]

					status := "syncing..."
					if w.IsWaiting() {
						status = "waiting..."
					}

					blockHeightProgress := values.StringF(values.StrBlockHeaderFetchedCount, w.GetBestBlock(), pg.headersToFetchOrScan)
					daysBehind := components.TimeFormat(int(currentSeconds-w.GetBestBlockTimeStamp()), true)
					details := pg.syncDetail(w.Name, status, blockHeightProgress, daysBehind)
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
func (pg *AppOverviewPage) walletSyncBox(gtx layout.Context, inset layout.Inset, details walletSyncDetails) layout.Dimensions {
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		card := pg.Theme.Card()
		card.Color = pg.Theme.Color.Gray4
		return card.Layout(gtx, func(gtx C) D {
			return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return components.EndToEndRow(gtx, details.name.Layout, details.status.Layout)
						})
					}),
					layout.Rigid(func(gtx C) D {
						headersFetchedTitleLabel := pg.Theme.Body2(values.String(values.StrBlockHeaderFetched))
						headersFetchedTitleLabel.Color = pg.Theme.Color.GrayText2
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return components.EndToEndRow(gtx, headersFetchedTitleLabel.Layout, details.blockHeaderFetched.Layout)
						})
					}),
					layout.Rigid(func(gtx C) D {
						progressTitleLabel := pg.Theme.Body2(values.String(values.StrSyncingProgress))
						progressTitleLabel.Color = pg.Theme.Color.GrayText2
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return components.EndToEndRow(gtx, progressTitleLabel.Layout, details.syncingProgress.Layout)
						})
					}),
				)
			})
		})
	})
}

func (pg *AppOverviewPage) rescanDetailsLayout(gtx layout.Context, inset layout.Inset) layout.Dimensions {
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
			return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return pg.Theme.Body1(wal.Name).Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						headersFetchedTitleLabel := pg.Theme.Body2("Blocks scanned")
						headersFetchedTitleLabel.Color = pg.Theme.Color.GrayText2

						blocksScannedLabel := pg.Theme.Body1(fmt.Sprint(rescanUpdate.ProgressReport.CurrentRescanHeight))
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return components.EndToEndRow(gtx, headersFetchedTitleLabel.Layout, blocksScannedLabel.Layout)
						})
					}),
					layout.Rigid(func(gtx C) D {
						progressTitleLabel := pg.Theme.Body2(values.String(values.StrSyncingProgress))
						progressTitleLabel.Color = pg.Theme.Color.GrayText2

						rescanProgress := fmt.Sprintf("%d blocks left", rescanUpdate.ProgressReport.TotalHeadersToScan-rescanUpdate.ProgressReport.CurrentRescanHeight)
						blocksScannedLabel := pg.Theme.Body1(rescanProgress)
						return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return components.EndToEndRow(gtx, progressTitleLabel.Layout, blocksScannedLabel.Layout)
						})
					}),
				)
			})
		})
	})
}
