package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
)

func (win *Window) TabbedWallets(selected, body layout.Widget, item layout.ListElement) {
	layouts.Tabs{
		Selected: selected,
		Item:     item,
		Body:     body,
		List:     win.tabsList,
		Flex: layout.Flex{
			Axis: layout.Horizontal,
		},
		Size:       .3,
		ButtonSize: .2,
	}.Layout(win.gtx, &win.selected, win.buttons.tabs)
}

func (win *Window) Header() {
	win.theme.Label(unit.Dp(50), "GoDcr").Layout(win.gtx)
}
