package materialplus

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// DefaultTabSize is the default flexed size of the tab section in a Tabs
const DefaultTabSize = .3

type TabItem struct {
	material.Button
}

func (t *TabItem) Layout(gtx *layout.Context, btn *widget.Button, selected bool) {
	gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
	t.Button.Layout(gtx, btn)
}

// Tabs laysout a Flexed(Size) List with Selected as the first element and Item as the rest.
type Tabs struct {
	List *layout.List
	Flex layout.Flex
	Size float32
}

func NewTabs() *Tabs {
	return &Tabs{
		List: &layout.List{
			Axis: layout.Vertical,
		},
		Size: DefaultTabSize,
	}
}

// Layout the tabs
func (tab *Tabs) Layout(gtx *layout.Context, selected *int, tabBtns []*widget.Button, tabItems []TabItem, body layout.Widget) {
	tab.Flex.Layout(gtx,
		layout.Flexed(tab.Size, func() {
			tab.List.Layout(gtx, len(tabBtns), func(i int) {
				tabItems[i].Layout(gtx, tabBtns[i], i == *selected)
			})
		}),
		layout.Flexed(1-tab.Size, body),
	)
}
