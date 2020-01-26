package materialplus

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr-gio/ui"
)

const (
	heightPercentage = 80 // percentage of window height the modal content takes
)

// Modal renders a modal instamce to screen
func (t *Theme) Modal(gtx *layout.Context, renderFunc func()) {
	overlayColor := ui.BlackColor
	overlayColor.A = 200

	Fill(gtx, overlayColor, gtx.Constraints.Width.Max, gtx.Constraints.Height.Max)

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
		}),
		layout.Stacked(func() {
			inset := layout.Inset{
				Top: unit.Dp(50),
			}
			inset.Layout(gtx, func() {
				Fill(gtx, ui.WhiteColor, gtx.Constraints.Width.Max, gtx.Constraints.Height.Max)
				inset := layout.Inset{
					Top:   unit.Dp(7),
					Left:  unit.Dp(25),
					Right: unit.Dp(25),
				}
				inset.Layout(gtx, renderFunc)
			})
		}),
	)
}
