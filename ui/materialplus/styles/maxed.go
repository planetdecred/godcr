package styles

import (
	"gioui.org/layout"
)

type maxed struct{}

var Maxed maxed

func (maxed) Styled(gtx *layout.Context, w layout.Widget) layout.Widget {
	return func() {
		gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
		gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
		w()
	}
}
