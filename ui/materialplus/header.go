package materialplus

import (
	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
)

const (
	HeaderSize = float32(.15)
)

func (t *Theme) WithHeader(gtx *layout.Context, header layout.Widget, body layout.Widget) {
	layouts.FlexWithTwoCildren{
		First:   header,
		Second:  body,
		Weight:  HeaderSize,
		Flex:    layout.Flex{Axis: layout.Vertical},
		Flexing: layouts.FirstFlexed,
	}.Layout(gtx)
}
