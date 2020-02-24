package materialplus

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
)

// TabbedLayout lays out 
func (t *Theme) TabbedLayout(gtx *layout.Context, selected *int, tabs []*widget.Button, selectedItem layout.Widget, item layout.ListElement, body layout.Widget) {
	.Layout(gtx, selected, tabs)
}
