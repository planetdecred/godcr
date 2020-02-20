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
// WindowHeight determines how high in the vertical axis the content area of the modal will be
// WindowWidth determines how wide in the horizonatl axix the content area of the modal will be
func (t *Theme) ModalPopUp(gtx *layout.Context, WindowHeight, WindowWidth int, renderFunc func()) {
	overlayColor := ui.BlackColor
	overlayColor.A = 200

	helper.PaintArea(gtx, overlayColor, gtx.Constraints.Width.Max, gtx.Constraints.Height.Max)

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			gtx.Constraints.Height.Min = (heightPercentage / 100) * gtx.Constraints.Height.Max
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

func (t *Theme) Modal(gtx *layout.Context, renderFunc func()) {
	overlayColor := ui.BlackColor
	overlayColor.A = 200

	helper.PaintArea(gtx, overlayColor, gtx.Constraints.Width.Max, helper.WindowHeight)

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			gtx.Constraints.Height.Min = 170
		}),
		layout.Stacked(func() {
			inset := layout.Inset{
				Top: unit.Dp(50),
			}
			inset.Layout(gtx, func() {
				helper.PaintArea(gtx, ui.WhiteColor, gtx.Constraints.Width.Max, gtx.Constraints.Height.Max)
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
