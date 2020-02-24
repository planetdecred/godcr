package layouts

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"
)

type Tabbed struct {
	Item     layout.ListElement
	Selected layout.Widget
	Body     layout.Widget

	List       *layout.List
	Flex       layout.Flex
	TabSize    float32
	ButtonSize float32
}

func (tab Tabbed) Layout(gtx *layout.Context, selected *int, tabs []*widget.Button) {
	FlexWithTwoCildren{
		Flex: tab.Flex,
		First: func() {
			tab.List.Layout(gtx, len(tabs), func(i int) {
				if i == *selected {
					tab.Selected()
				} else {
					Clicker{
						Widget: func() { tab.Item(i) },
					}.Layout(gtx, tabs[i])
				}
				pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
				tabs[i].Layout(gtx)
			})
		},
		Second: tab.Body,
		Weight: tab.TabSize,
	}.Layout(gtx)
}
