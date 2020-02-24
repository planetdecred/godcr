package layouts

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
)

// Tabs displays a tabbed document
// TODO: Doc
type Tabs struct {
	Item     layout.ListElement
	Selected layout.Widget
	Body     layout.Widget

	List *layout.List
	Flex layout.Flex

	Size       float32
	ButtonSize float32
}

// Layout a widget.
func (tab Tabs) Layout(gtx *layout.Context, selected *int, tabs []*widget.Button) {
	tab.Flex.Layout(gtx,
		layout.Flexed(tab.Size, func() {
			tab.List.Layout(gtx, len(tabs), func(i int) {
				if i == *selected {
					tab.Selected()
				} else {
					Clicker(func() { tab.Item(i) }).Layout(gtx, tabs[i])
				}
				pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
				tabs[i].Layout(gtx)
			})
		}),
		layout.Flexed(1-tab.Size, tab.Body),
	)
}
