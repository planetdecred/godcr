package materialplus

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/ui/layouts"
)

func (t *Theme) Tabbed(gtx *layout.Context, selected *int, tabs []*widget.Button, selectedItem layout.Widget, item layout.ListElement, body layout.Widget) layout.Widget {
	tabbed := layouts.Tabbed{
		Selected: selectedItem,
		Item:     item,
		Body:     body,
		List:     &layout.List{Axis: layout.Vertical},
		Flex: layout.Flex{
			Axis: layout.Horizontal,
		},
		TabSize:    .3,
		ButtonSize: .2,
	}
	return func() {
		tabbed.Layout(gtx, selected, tabs)
	}
}
