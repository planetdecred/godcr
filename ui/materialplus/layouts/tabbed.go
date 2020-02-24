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

	Size float32
}

// Layout a widget.
func (tab Tabs) Layout(gtx *layout.Context, selected *int, tabs []*widget.Button) {
	tab.Flex.Layout(gtx,
		layout.Flexed(tab.Size, func() {
			tab.List.Layout(gtx, len(tabs), func(i int) {
				Clicker(func() {
					if i == *selected {
						tab.Selected()
					} else {
						tab.Item(i)
					}
				}).Layout(gtx, tabs[i])
				pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
				tabs[i].Layout(gtx)
			})
		}),
		layout.Rigid(tab.Body),
	)
}

func (tab Tabs) Layedout(gtx *layout.Context, selected *int, tabs []*widget.Button) layout.Widget {
	return func() {
		tab.Layout(gtx, selected, tabs)
	}

}
