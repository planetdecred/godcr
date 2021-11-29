package overview

import (
	"fmt"
	"image/color"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *AppOverviewPage) initWalletStatusWidgets() {
	pg.toggleSyncDetails = pg.Theme.Button(values.String(values.StrShowDetails))
	pg.toggleSyncDetails.TextSize = values.TextSize16
	pg.toggleSyncDetails.Background = color.NRGBA{}
	pg.toggleSyncDetails.Color = pg.Theme.Color.Primary
	pg.toggleSyncDetails.Inset = layout.Inset{}

	pg.syncedIcon = decredmaterial.NewIcon(pg.Icons.ActionCheckCircle)
	pg.syncedIcon.Color = pg.Theme.Color.Success

	pg.syncingIcon = pg.Icons.SyncingIcon

	pg.notSyncedIcon = decredmaterial.NewIcon(pg.Icons.NavigationCancel)
	pg.notSyncedIcon.Color = pg.Theme.Color.Danger

	pg.walletStatusIcon = decredmaterial.NewIcon(pg.Icons.ImageBrightness1)
}

// syncStatusSection lays out content for displaying sync status.
func (pg *AppOverviewPage) syncStatusSection(gtx layout.Context) layout.Dimensions {
	uniform := layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPadding5}
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		return components.Container{Padding: layout.Inset{
			Top:    values.MarginPadding15,
			Bottom: values.MarginPadding16,
		}}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.syncBoxTitleRow),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(pg.syncStatusIcon),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return pg.syncStatusTextRow(gtx, uniform)
									}),
									layout.Rigid(func(gtx C) D {
										if !pg.walletSyncing && !pg.rescanningBlocks {
											return pg.syncDormantContent(gtx, uniform)
										}
										return layout.Dimensions{}
									}),
									layout.Rigid(func(gtx C) D {
										if pg.walletSyncing || pg.rescanningBlocks {
											return pg.progressBarRow(gtx, uniform)
										}
										return layout.Dimensions{}
									}),
									layout.Rigid(func(gtx C) D {
										if pg.walletSyncing || pg.rescanningBlocks {
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
					if pg.walletSyncing || pg.rescanningBlocks {
						return pg.Theme.Separator().Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.syncDetailsVisibility {
						if pg.walletSyncing {
							return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
								return pg.walletSyncRow(gtx, uniform)
							})
						} else if pg.rescanningBlocks {
							return components.Container{Padding: layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
								return pg.rescanDetailsLayout(gtx, uniform)
							})
						}
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					if (pg.walletSyncing || pg.rescanningBlocks) && pg.syncDetailsVisibility {
						return pg.Theme.Separator().Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.walletSyncing || pg.rescanningBlocks {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Inset{Top: values.MarginPadding14}.Layout(gtx, pg.toggleSyncDetails.Layout)
					}
					return layout.Dimensions{}
				}),
			)
		})
	})
}

// syncBoxTitleRow lays out widgets in the title row inside the sync status box.
func (pg *AppOverviewPage) syncBoxTitleRow(gtx layout.Context) layout.Dimensions {
	title := pg.Theme.Body2(values.String(values.StrWalletStatus))
	title.Color = pg.Theme.Color.GrayText1
	statusLabel := pg.Theme.Body1(values.String(values.StrOffline))
	pg.walletStatusIcon.Color = pg.Theme.Color.Danger
	if pg.isConnnected {
		statusLabel.Text = values.String(values.StrOnline)
		pg.walletStatusIcon.Color = pg.Theme.Color.Success
	}

	syncStatus := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding4, Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
					return pg.walletStatusIcon.Layout(gtx, values.MarginPadding14)
				})
			}),
			layout.Rigid(statusLabel.Layout),
		)
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	titlePadding := values.MarginPadding15
	return components.Container{Padding: layout.Inset{
		Left:   titlePadding,
		Right:  titlePadding,
		Bottom: titlePadding,
	}}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(title.Layout),
			layout.Rigid(func(gtx C) D {
				if len(pg.transactions) > 0 {
					return syncStatus(gtx)
				}
				return D{}
			}),
		)
	})
}

func (pg *AppOverviewPage) syncStatusIcon(gtx layout.Context) layout.Dimensions {
	syncStatusIcon := pg.notSyncedIcon
	syncStatusIcon.Color = pg.Theme.Color.Danger
	if pg.walletSynced {
		syncStatusIcon = pg.syncedIcon
		syncStatusIcon.Color = pg.Theme.Color.Success
	}
	i := layout.Inset{Right: values.MarginPadding16, Top: values.MarginPadding9}
	if pg.walletSyncing {
		return i.Layout(gtx, pg.syncingIcon.Layout24dp)
	}
	return i.Layout(gtx, func(gtx C) D {
		return syncStatusIcon.Layout(gtx, values.MarginPadding24)
	})
}

// syncDormantContent lays out sync status content when the wallet is synced or not connected
func (pg *AppOverviewPage) syncDormantContent(gtx layout.Context, uniform layout.Inset) layout.Dimensions {
	return uniform.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding12}.Layout(gtx, pg.blockInfoRow)
			}),
			layout.Rigid(func(gtx C) D {
				if pg.walletSynced {
					return pg.connectionPeer(gtx)
				}
				latestBlockTitleLabel := pg.Theme.Body1(values.String(values.StrNoConnectedPeer))
				latestBlockTitleLabel.Color = pg.Theme.Color.GrayText2
				return latestBlockTitleLabel.Layout(gtx)
			}),
		)
	})
}

func (pg *AppOverviewPage) blockInfoRow(gtx layout.Context) layout.Dimensions {
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
			}.Layout(gtx, pg.Theme.Body1(fmt.Sprintf("%d", pg.bestBlock.Height)).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			pg.walletStatusIcon.Color = pg.Theme.Color.Gray1
			return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
				return pg.walletStatusIcon.Layout(gtx, values.MarginPadding5)
			})
		}),
		layout.Rigid(pg.Theme.Body1(components.TimeAgo(pg.bestBlock.Timestamp)).Layout),
	)
}
