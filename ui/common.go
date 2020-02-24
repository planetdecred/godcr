package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
)

// TabbedWallets layouts a layout.Tabs
func (win *Window) TabbedWallets(selected, body layout.Widget, item layout.ListElement) {
	layouts.Tabs{
		Selected: selected,
		Item:     item,
		Body:     body,
		List:     win.tabsList,
		Flex: layout.Flex{
			Axis: layout.Horizontal,
		},
		Size: .3,
	}.Layout(win.gtx, &win.selected, win.buttons.tabs)
}

// Header lays out the window header
func (win *Window) Header() {
	win.theme.Label(unit.Dp(50), "GoDcr").Layout(win.gtx)
}

// Background fills the context with theme Background
func (win *Window) Background() {
	layouts.FillWithColor(win.gtx, win.theme.Background)
}
