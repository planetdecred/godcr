package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// DefaultTabSize is the default flexed size of the tab section in a Tabs
const DefaultTabSize = .15

type TabItem struct {
	Button
}

func (t *TabItem) Layout(gtx *layout.Context, btn *widget.Button, selected bool) {
	gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
	t.Button.Layout(gtx, btn)
}

// Tabs laysout a Flexed(Size) List with Selected as the first element and Item as the rest.
type Tabs struct {
	List     *layout.List
	Flex     layout.Flex
	Size     float32
	items    []TabItem
	Selected int
	changed  bool
	btns     []*widget.Button
}

func NewTabs() *Tabs {
	return &Tabs{
		List: &layout.List{
			Axis: layout.Vertical,
		},
		Size: DefaultTabSize,
	}
}

func (tab *Tabs) SetTabs(tabs []TabItem) {
	tab.items = tabs
	if len(tab.items) != len(tab.btns) {
		tab.btns = make([]*widget.Button, len(tab.items))
		for i := range tab.btns {
			tab.btns[i] = new(widget.Button)
		}
	}
}

func (tab *Tabs) Changed() bool {
	return tab.changed
}

// Layout the tabs
func (tab *Tabs) Layout(gtx *layout.Context, body layout.Widget) {
	tab.Flex.Layout(gtx,
		layout.Flexed(tab.Size, func() {
			tab.List.Layout(gtx, len(tab.btns), func(i int) {
				tab.items[i].Layout(gtx, tab.btns[i], i == tab.Selected)
			})
		}),
		layout.Flexed(1-tab.Size, body),
	)
	for i := range tab.btns {
		if tab.btns[i].Clicked(gtx) {
			tab.changed = true
			tab.Selected = i
			return
		}
	}
}
