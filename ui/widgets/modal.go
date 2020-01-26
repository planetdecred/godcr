package widgets

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/helper"

)

// Modal represents a popup widget with a background overlay
type Modal struct {
	isOpen       bool
	overlayColor color.RGBA
}

// NewModal returns a new Modal instance
func NewModal() *Modal {
	overlayColor := helper.BlackColor
	overlayColor.A = 200

	return &Modal{
		overlayColor: overlayColor,
		isOpen:       true,
	}
}

// Draw renders the modal instamce to screen
func (m *Modal) Draw(gtx *layout.Context, theme *material.Theme, renderFunc func()) {
	helper.PaintArea(gtx, m.overlayColor, helper.WindowWidth, helper.WindowHeight)

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			gtx.Constraints.Height.Min = 100
		}),
		layout.Stacked(func() {
			inset := layout.Inset{
				Top: unit.Dp(50),
			}
			inset.Layout(gtx, func() {
				helper.PaintArea(gtx, helper.WhiteColor, gtx.Constraints.Width.Max, gtx.Constraints.Height.Max)
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
