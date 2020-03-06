package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
)

const (
	headerHeight = .15
	navSize      = .1
)

func (win *Window) Page(body layout.Widget) {
	bd := func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Flexed(headerHeight, func() {
				decredmaterial.Card{
					Inset: layout.Inset{
						Bottom: unit.Dp(1),
					},
				}.Layout(win.gtx, win.Header)

			}),
			layout.Flexed(1-headerHeight, func() {
				toMax(win.gtx)
				body()
			}),
		)
	}
	layout.Flex{Axis: layout.Horizontal}.Layout(win.gtx,
		layout.Flexed(navSize, func() {
			decredmaterial.Card{
				Inset: layout.Inset{
					Right: unit.Dp(1),
				},
			}.Layout(win.gtx, win.NavBar)
		}),
		layout.Flexed(1-navSize, bd),
	)
}

// TabbedPage layouts a layout.Tabs
func (win *Window) TabbedPage(body layout.Widget) {
	items := make([]decredmaterial.TabItem, win.walletInfo.LoadedWallets)
	for i := 0; i < win.walletInfo.LoadedWallets; i++ {
		items[i] = decredmaterial.TabItem{
			Button: win.theme.Button(win.walletInfo.Wallets[i].Name),
		}
	}
	bd := func() {
		toMax(win.gtx)
		win.tabs.Layout(win.gtx, body)
	}
	win.Page(bd)
}

// Header lays out the window header
func (win *Window) Header() {
	toMax(win.gtx)
	layout.Flex{
		Alignment: layout.Middle,
	}.Layout(win.gtx,
		layout.Flexed(0.4, func() {
			win.theme.H3("GoDcr").Layout(win.gtx)
		}),
		layout.Flexed(0.2, func() {
			layout.Center.Layout(win.gtx, func() {
				win.outputs.createDiag.Layout(win.gtx, &win.inputs.createDiag)
			})
		}),
		layout.Flexed(0.4, func() {
			layout.Center.Layout(win.gtx, func() {
				win.outputs.sync.Layout(win.gtx, &win.inputs.sync)
			})
		}),
	)
}

func (win *Window) NavBar() {
	toMax(win.gtx)
	layout.Flex{Spacing: layout.SpaceEvenly, Alignment: layout.Middle, Axis: layout.Vertical}.Layout(win.gtx,
		layout.Rigid(func() {
			layout.Center.Layout(win.gtx, func() {
				win.outputs.toOverview.Layout(win.gtx, &win.inputs.toOverview)
			})
		}),
		layout.Rigid(func() {
			layout.Center.Layout(win.gtx, func() {
				win.outputs.toWallets.Layout(win.gtx, &win.inputs.toWallets)
			})
		}),
		layout.Rigid(func() {
			layout.Center.Layout(win.gtx, func() {
				win.outputs.toTransactions.Layout(win.gtx, &win.inputs.toTransactions)
			})
		}),
		layout.Rigid(func() {
			layout.Center.Layout(win.gtx, func() {
				win.outputs.toSettings.Layout(win.gtx, &win.inputs.toSettings)
			})
		}),
	)
}

func toMax(gtx *layout.Context) {
	gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
	gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
}

func (win *Window) Err() {
	if win.err != "" {
		win.outputs.err.Text = win.err
		win.outputs.err.Layout(win.gtx)
	}
}
