package common

import (
	"fmt"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr-gio/helper"
	"github.com/raedahgroup/godcr-gio/widgets"
)

type report struct {
	percentageProgress float64
	timeRemaining      string
	daysBehind         string
}

type widgetItems struct {
	reconnectButton  *widgets.Button
	cancelButton     *widgets.Button
	progressBar      *widgets.ProgressBar
	disconnectButton *widgets.Button
}

type Syncer struct {
	showDetails    bool
	wallet         *helper.MultiWallet
	refreshDisplay func()
	syncError      error

	report  *report
	widgets *widgetItems
}

func NewSyncer(wallet *helper.MultiWallet, refreshDisplay func()) *Syncer {
	s := &Syncer{
		wallet:         wallet,
		refreshDisplay: refreshDisplay,
		syncError:      nil,
		showDetails:    false,
		report:         &report{},
	}

	s.widgets = &widgetItems{
		progressBar:      widgets.NewProgressBar(),
		reconnectButton:  widgets.NewButton("Reconnect", nil).SetBorderColor(helper.GrayColor).SetBackgroundColor(helper.WhiteColor).SetColor(helper.BlackColor),
		cancelButton:     widgets.NewButton("Cancel", nil).SetBorderColor(helper.GrayColor).SetBackgroundColor(helper.WhiteColor).SetColor(helper.BlackColor),
		disconnectButton: widgets.NewButton("Disconnect", nil).SetBorderColor(helper.GrayColor).SetBackgroundColor(helper.WhiteColor).SetColor(helper.BlackColor),
	}

	return s
}

func (s *Syncer) OnSyncStarted() {}

func (s *Syncer) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	s.refreshDisplay()
}

func (s *Syncer) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	s.report.percentageProgress = float64(headersFetchProgress.TotalSyncProgress)
	s.report.timeRemaining = dcrlibwallet.CalculateTotalTimeRemaining(headersFetchProgress.TotalTimeRemainingSeconds)
	s.report.daysBehind = dcrlibwallet.CalculateDaysBehind(headersFetchProgress.CurrentHeaderTimestamp)
	s.refreshDisplay()
}

func (s *Syncer) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	s.report.percentageProgress = float64(addressDiscoveryProgress.TotalSyncProgress)
	s.refreshDisplay()
}

func (s *Syncer) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	s.report.percentageProgress = float64(headersRescanProgress.TotalSyncProgress)
	s.refreshDisplay()
}

func (s *Syncer) OnSyncCompleted() {
	s.report.percentageProgress = 100
	s.refreshDisplay()
}

func (s *Syncer) OnSyncCanceled(willRestart bool) {}

func (s *Syncer) OnSyncEndedWithError(err error) {
	s.syncError = err
	s.refreshDisplay()
}

func (s *Syncer) Debug(debugInfo *dcrlibwallet.DebugInfo) {}

func (s *Syncer) OnTransaction(transaction string) {}

func (s *Syncer) OnBlockAttached(walletID int, blockHeight int32) {}

func (s *Syncer) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {}

func (s *Syncer) Render(ctx *layout.Context) {
	helper.PaintArea(ctx, helper.WhiteColor, ctx.Constraints.Width.Max, 140)

	inset := layout.UniformInset(unit.Dp(15))
	inset.Layout(ctx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
			layout.Rigid(func() {
				widgets.NewLabel("Wallet Status").
					SetColor(helper.GrayColor).
					SetSize(4).
					SetWeight(text.Bold).
					Draw(ctx)
			}),
			layout.Flexed(1, func() {
				layout.Align(layout.NE).Layout(ctx, func() {
					s.drawWalletStatus(ctx)
				})
			}),
		)

		inset := layout.Inset{
			Top: unit.Dp(25),
		}

		inset.Layout(ctx, func() {
			if s.wallet.IsSynced() {
				s.drawIsSyncedCard(ctx)
			} else if !s.wallet.IsSyncing() {
				s.drawIsSyncingCard(ctx)
			} else {
				s.drawNotSyncedCard(ctx)
			}
		})
	})
}

func (s *Syncer) drawWalletStatus(ctx *layout.Context) {
	var indicatorColor color.RGBA
	var statusText string

	if s.IsOnline() {
		indicatorColor = helper.DecredGreenColor
		statusText = "Online"
	} else {
		indicatorColor = helper.DecredOrangeColor
		statusText = "Offline"
	}

	indicatorSize := float32(8)
	radius := indicatorSize * .5

	layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
		layout.Rigid(func() {
			inset := layout.Inset{
				Top:   unit.Dp(3.5),
				Right: unit.Dp(5),
			}
			inset.Layout(ctx, func() {
				clip.Rect{
					Rect: f32.Rectangle{
						Max: f32.Point{
							X: indicatorSize,
							Y: indicatorSize,
						},
					},
					NE: radius,
					NW: radius,
					SE: radius,
					SW: radius,
				}.Op(ctx.Ops).Add(ctx.Ops)
				helper.Fill(ctx, indicatorColor, int(indicatorSize), int(indicatorSize))
			})
		}),
		layout.Rigid(func() {
			widgets.NewLabel(statusText).
				SetSize(4).
				SetColor(helper.GrayColor).
				SetWeight(text.Bold).
				Draw(ctx)
		}),
	)
}

func (s *Syncer) IsOnline() bool {
	if !s.wallet.IsSynced() || !s.wallet.IsSyncing() {
		return false
	}
	return true
}

func (s *Syncer) drawNotSyncedCard(ctx *layout.Context) {
	layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
		layout.Rigid(func() {
			widgets.CancelIcon.SetColor(helper.DangerColor).Draw(ctx, 25)
		}),
		layout.Rigid(func() {
			inset := layout.Inset{
				Left: unit.Dp(40),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("Not Synced").
					SetSize(5).
					SetColor(helper.BlackColor).
					SetWeight(text.Bold).
					Draw(ctx)
			})
		}),
		layout.Flexed(1, func() {
			layout.Align(layout.NE).Layout(ctx, func() {
				ctx.Constraints.Height.Max = 35
				ctx.Constraints.Width.Max = 130

				s.widgets.reconnectButton.Draw(ctx, func() {

				})
			})
		}),
	)

	inset := layout.Inset{
		Top:  unit.Dp(35),
		Left: unit.Dp(40),
	}
	inset.Layout(ctx, func() {
		lowestBlock := s.wallet.GetLowestBlock()
		lowestBlockHeight := int32(-1)

		if lowestBlock != nil {
			lowestBlockHeight = lowestBlock.Height
		}

		widgets.NewLabel(fmt.Sprintf("Synced to block %d - %s", lowestBlockHeight, s.report.daysBehind)).
			SetSize(4).
			SetColor(helper.GrayColor).
			Draw(ctx)
	})

	inset = layout.Inset{
		Top:  unit.Dp(65),
		Left: unit.Dp(40),
	}
	inset.Layout(ctx, func() {
		widgets.NewLabel("No connected peers").
			SetSize(4).
			SetColor(helper.GrayColor).
			Draw(ctx)
	})
}

func (s *Syncer) drawIsSyncedCard(ctx *layout.Context) {
	layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
		layout.Rigid(func() {
			inset := layout.Inset{
				Top: unit.Dp(3),
			}
			inset.Layout(ctx, func() {
				widgets.NewCheckbox().MakeAsIcon().Draw(ctx)
			})
		}),
		layout.Rigid(func() {
			inset := layout.Inset{
				Left: unit.Dp(30),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("Synced").
					SetSize(5).
					SetColor(helper.BlackColor).
					SetWeight(text.Bold).
					Draw(ctx)
			})
		}),
		layout.Flexed(1, func() {
			layout.Align(layout.NE).Layout(ctx, func() {
				ctx.Constraints.Height.Max = 35
				ctx.Constraints.Width.Max = 130

				s.widgets.disconnectButton.Draw(ctx, func() {

				})
			})
		}),
	)

	inset := layout.Inset{
		Top:  unit.Dp(35),
		Left: unit.Dp(40),
	}
	inset.Layout(ctx, func() {
		lowestBlock := s.wallet.GetLowestBlock()
		lowestBlockHeight := int32(-1)

		if lowestBlock != nil {
			lowestBlockHeight = lowestBlock.Height
		}

		widgets.NewLabel(fmt.Sprintf("Synced to block %d - %s", lowestBlockHeight, s.report.daysBehind)).
			SetSize(4).
			SetColor(helper.GrayColor).
			Draw(ctx)
	})

	inset = layout.Inset{
		Top:  unit.Dp(65),
		Left: unit.Dp(40),
	}
	inset.Layout(ctx, func() {
		widgets.NewLabel(fmt.Sprintf("Connected to %d peers", s.wallet.ConnectedPeers())).
			SetSize(4).
			SetColor(helper.GrayColor).
			Draw(ctx)
	})

}

func (s *Syncer) drawIsSyncingCard(ctx *layout.Context) {
	layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
		layout.Rigid(func() {
			inset := layout.Inset{
				Top: unit.Dp(3),
			}
			inset.Layout(ctx, func() {
				helper.PaintCircle(ctx, helper.DecredGreenColor, 15)
			})
		}),
		layout.Rigid(func() {
			inset := layout.Inset{
				Left: unit.Dp(30),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("Syncing...").
					SetSize(5).
					SetColor(helper.BlackColor).
					SetWeight(text.Bold).
					Draw(ctx)
			})
		}),
		layout.Flexed(1, func() {
			layout.Align(layout.NE).Layout(ctx, func() {
				ctx.Constraints.Height.Max = 35
				ctx.Constraints.Width.Max = 130

				s.widgets.cancelButton.Draw(ctx, func() {

				})
			})
		}),
	)

	inset := layout.Inset{
		Top: unit.Dp(38),
	}
	inset.Layout(ctx, func() {
		s.widgets.progressBar.
			SetHeight(10).
			SetBackgroundColor(helper.GrayColor).
			Draw(ctx, &s.report.percentageProgress)
	})

	inset = layout.Inset{
		Top: unit.Dp(50),
	}
	inset.Layout(ctx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
			layout.Rigid(func() {
				widgets.NewLabel(fmt.Sprintf("%d%%", int(s.report.percentageProgress))).
					SetColor(helper.BlackColor).
					SetSize(4).
					SetWeight(text.Bold).
					Draw(ctx)
			}),
			layout.Flexed(1, func() {
				layout.Align(layout.NE).Layout(ctx, func() {
					widgets.NewLabel(s.report.timeRemaining).
						SetColor(helper.BlackColor).
						SetSize(4).
						SetWeight(text.Bold).
						Draw(ctx)
				})
			}),
		)
	})
}
