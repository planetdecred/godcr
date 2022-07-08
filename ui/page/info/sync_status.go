package info

import (
	"fmt"
	"time"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

func (pg *InfoPage) initWalletStatusWidgets() {
	pg.syncedIcon = decredmaterial.NewIcon(pg.Theme.Icons.ActionCheckCircle)
	pg.syncedIcon.Color = pg.Theme.Color.Success

	pg.syncingIcon = pg.Theme.Icons.SyncingIcon

	pg.notSyncedIcon = decredmaterial.NewIcon(pg.Theme.Icons.NavigationCancel)
	pg.notSyncedIcon.Color = pg.Theme.Color.Danger

	pg.walletStatusIcon = decredmaterial.NewIcon(pg.Theme.Icons.ImageBrightness1)
	pg.syncSwitch = pg.Theme.Switch()
}

// syncStatusSection lays out content for displaying sync status.
func (pg *InfoPage) syncStatusSection(gtx C) D {
	syncing, rescanning := pg.WL.MultiWallet.IsSyncing(), pg.WL.MultiWallet.IsRescanning()
	uniform := layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPadding5}
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		return components.Container{Padding: layout.Inset{
			Top:    values.MarginPadding15,
			Bottom: values.MarginPadding16,
		}}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.syncBoxTitleRow),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding20, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(pg.syncStatusIcon),
							layout.Rigid(pg.syncStatusTextRow),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					switch {
					case syncing:
						return pg.walletSyncRow(gtx, uniform)
					case rescanning:
						return pg.rescanDetailsLayout(gtx, uniform)
					default:
						return pg.syncDormantContent(gtx, uniform)
					}
				}),
			)
		})
	})
}

// syncBoxTitleRow lays out widgets in the title row inside the sync status box.
func (pg *InfoPage) syncBoxTitleRow(gtx C) D {
	statusLabel := pg.Theme.Label(values.TextSize14, values.String(values.StrOffline))
	pg.walletStatusIcon.Color = pg.Theme.Color.Danger
	if pg.WL.MultiWallet.IsConnectedToDecredNetwork() {
		statusLabel.Text = values.String(values.StrOnline)
		pg.walletStatusIcon.Color = pg.Theme.Color.Success
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(pg.Theme.Label(values.TextSize14, values.String(values.StrWalletStatus)).Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding4,
				Left:  values.MarginPadding4,
			}.Layout(gtx, func(gtx C) D {
				return pg.walletStatusIcon.Layout(gtx, values.MarginPadding10)
			})
		}),
		layout.Rigid(statusLabel.Layout),
		layout.Rigid(func(gtx C) D {
			if pg.WL.MultiWallet.IsSyncing() || pg.WL.MultiWallet.IsSynced() {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						connectedPeers := fmt.Sprintf("%d", pg.WL.MultiWallet.ConnectedPeers())
						return pg.Theme.Label(values.TextSize14, values.StringF(values.StrConnectedTo, connectedPeers)).Layout(gtx)
					}),
				)
			}

			return pg.Theme.Label(values.TextSize14, values.String(values.StrNoConnectedPeer)).Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, pg.layoutAutoSyncSection)
		}),
	)
}

func (pg *InfoPage) syncStatusIcon(gtx C) D {
	syncStatusIcon := pg.notSyncedIcon
	syncStatusIcon.Color = pg.Theme.Color.Danger
	if pg.WL.MultiWallet.IsSynced() {
		syncStatusIcon = pg.syncedIcon
		syncStatusIcon.Color = pg.Theme.Color.Success
	}
	i := layout.Inset{Right: values.MarginPadding16}
	if pg.WL.MultiWallet.IsSyncing() {
		return i.Layout(gtx, func(gtx C) D {
			return pg.syncingIcon.LayoutSize(gtx, values.MarginPadding20)
		})
	}

	return i.Layout(gtx, func(gtx C) D {
		return syncStatusIcon.Layout(gtx, values.MarginPadding20)
	})
}

// syncDormantContent lays out sync status content when the wallet is synced or not connected
func (pg *InfoPage) syncDormantContent(gtx C, uniform layout.Inset) D {
	return uniform.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding12}.Layout(gtx, pg.blockInfoRow)
			}),
		)
	})
}

func (pg *InfoPage) blockInfoRow(gtx C) D {
	bestBlock := pg.WL.MultiWallet.GetBestBlock()
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			latestBlockTitleLabel := pg.Theme.Body1(values.String(values.StrLastBlockHeight))
			latestBlockTitleLabel.Color = pg.Theme.Color.GrayText2
			return latestBlockTitleLabel.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Left:  values.MarginPadding5,
				Right: values.MarginPadding5,
			}.Layout(gtx, pg.Theme.Body1(fmt.Sprintf("%d", bestBlock.Height)).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			pg.walletStatusIcon.Color = pg.Theme.Color.Gray1
			return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
				return pg.walletStatusIcon.Layout(gtx, values.MarginPadding5)
			})
		}),
		layout.Rigid(pg.Theme.Body1(components.TimeAgo(bestBlock.Timestamp)).Layout),
	)
}

func (pg *InfoPage) layoutAutoSyncSection(gtx C) D {
	return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, pg.syncSwitch.Layout)
			}),
			layout.Rigid(pg.Theme.Body2(values.String(values.StrSync)).Layout),
		)
	})
}

// syncDetail returns a walletSyncDetails object containing data of a single wallet sync box
func (pg *InfoPage) syncDetail(status, headersFetched, progress string) walletSyncDetails {
	return walletSyncDetails{
		status:             pg.Theme.Body2(status),
		blockHeaderFetched: pg.Theme.Body2(headersFetched),
		syncingProgress:    pg.Theme.Body2(values.StringF(values.StrSyncingProgressStat, progress)),
	}
}

// syncStatusTextRow lays out sync status text and sync button.
func (pg *InfoPage) syncStatusTextRow(gtx C) D {
	syncing, rescanning := pg.WL.MultiWallet.IsSyncing(), pg.WL.MultiWallet.IsRescanning()

	syncStatusLabel := pg.Theme.Label(values.TextSize14, values.String(values.StrWalletNotSynced))
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
		children = append(children, layout.Rigid(pg.progressBarRow))
	}

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
	}.Layout(gtx, children...)
}

// progressBarRow lays out the progress bar.
func (pg *InfoPage) progressBarRow(gtx C) D {
	return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
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
func (pg *InfoPage) progressStatusDetails() (int, string) {
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
func (pg *InfoPage) walletSyncRow(gtx C, inset layout.Inset) D {
	currentSeconds := time.Now().UnixNano() / int64(time.Second)
	w := pg.WL.SelectedWallet.Wallet

	status := values.String(values.StrSyncingState)
	if pg.WL.SelectedWallet.Wallet.IsWaiting() {
		status = values.String(values.StrWaitingState)
	}

	blockHeightProgress := values.StringF(values.StrBlockHeaderFetchedCount, w.GetBestBlock(), pg.headersToFetchOrScan)
	daysBehind := components.TimeFormat(int(currentSeconds-w.GetBestBlockTimeStamp()), true)
	details := pg.syncDetail(status, blockHeightProgress, daysBehind)
	col := pg.Theme.Color.GrayText2

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			completedSteps := pg.Theme.Body2(values.StringF(values.StrSyncSteps, pg.syncStep))
			completedSteps.Color = col

			headersFetched := pg.Theme.Body2(values.StringF(values.StrFetchingBlockHeaders, pg.stepFetchProgress))
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
		layout.Rigid(func(gtx C) D {
			_, timeLeft := pg.progressStatusDetails()
			timeLeftLabel := pg.Theme.Body2(timeLeft)

			remainingSyncTime := pg.Theme.Body2(values.String(values.StrSyncCompTime))
			remainingSyncTime.Color = col

			return inset.Layout(gtx, func(gtx C) D {
				return components.EndToEndRow(gtx, remainingSyncTime.Layout, timeLeftLabel.Layout)
			})
		}),
	)
}

func (pg *InfoPage) rescanDetailsLayout(gtx C, inset layout.Inset) D {
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
