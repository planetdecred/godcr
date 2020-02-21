package ui

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus"
)

const (
	headerHeight = .15
	navSize      = .15
)

func (win *Window) Page(body layout.Widget) {
	layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
		layout.Flexed(headerHeight, win.Header),
		layout.Flexed(1-headerHeight-navSize, func() {
			toMax(win.gtx)
			win.theme.Foreground(win.gtx, body)
		}),
		layout.Flexed(navSize, win.NavBar),
	)
}

// TabbedPage layouts a layout.Tabs
func (win *Window) TabbedPage(body layout.Widget) {
	items := make([]materialplus.TabItem, win.walletInfo.LoadedWallets)
	for i := 0; i < win.walletInfo.LoadedWallets; i++ {
		items[i] = materialplus.TabItem{
			Button: win.theme.Button(win.walletInfo.Wallets[i].Name),
		}
	}
	bd := func() {
		toMax(win.gtx)
		win.tabs.Layout(win.gtx, &win.selected, win.inputs.tabs, items, body)
	}
	win.Page(bd)
}

// Header lays out the window header
func (win *Window) Header() {
	toMax(win.gtx)
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(0.4, func() {
			win.theme.H3("GoDcr").Layout(win.gtx)
		}),
		layout.Rigid(func() {
			win.outputs.createDiag.Layout(win.gtx, &win.inputs.createDiag)
		}),
	)
}

func (win *Window) NavBar() {
	toMax(win.gtx)
	layout.Flex{Spacing: layout.SpaceEvenly, Alignment: layout.Middle}.Layout(win.gtx,
		layout.Rigid(func() {
			win.outputs.toWallets.Layout(win.gtx, &win.inputs.toWallets)
		}),
		layout.Rigid(func() {
			win.outputs.toOverview.Layout(win.gtx, &win.inputs.toOverview)
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
