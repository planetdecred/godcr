package layouts

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

// Tabs laysout a Flexed(Size) List with Selected as the first element and Item as the rest.
type Tabs struct {
	Item     layout.ListElement
	Selected layout.Widget
	Body     layout.Widget

	List *layout.List
	Flex layout.Flex

	Size float32
}

// Layout the tabs
func (tab Tabs) Layout(gtx *layout.Context, selected *int, tabs []*widget.Button) {
	tab.Flex.Layout(gtx,
		layout.Flexed(tab.Size, func() {
			tab.List.Layout(gtx, len(tabs), func(i int) {
				Clickable(func() {
					if i != 0 {
						tab.Item(i)
					} else {
						tab.Selected()
					}
				}).Layout(gtx, tabs[i])
				tabs[i].Layout(gtx)
			})
		}),
		layout.Rigid(tab.Body),
	)
}
