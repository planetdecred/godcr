package materialplus

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// DangerButton a button with the background set to theme.Danger
func (t *Theme) DangerButton(text string) material.Button {
	btn := t.Button(text)
	btn.Background = t.Danger
	return btn
}

// IconLabel layouts flex with Rigid ic and Rigid lbl
func IconLabel(gtx *layout.Context, flex layout.Flex, ic *material.Icon, lbl material.Label) {
	flex.Layout(gtx,
		layout.Rigid(func() { ic.Layout(gtx, unit.Dp(20)) }),
		layout.Rigid(func() { lbl.Layout(gtx) }),
	)
}

type LabelButton struct {
}

// func (b LabelButton) Layout(gtx *layout.Context, button *widget.Button) {
// 	col := b.Color
// 	bgcol := b.Background
// 	hmin := gtx.Constraints.Width.Min
// 	vmin := gtx.Constraints.Height.Min
// 	layout.Stack{Alignment: layout.Center}.Layout(gtx,
// 		layout.Expanded(func() {
// 			rr := float32(gtx.Px(unit.Dp(4)))
// 			clip.Rect{
// 				Rect: f32.Rectangle{Max: f32.Point{
// 					X: float32(gtx.Constraints.Width.Min),
// 					Y: float32(gtx.Constraints.Height.Min),
// 				}},
// 				NE: rr, NW: rr, SE: rr, SW: rr,
// 			}.Op(gtx.Ops).Add(gtx.Ops)
// 			fill(gtx, bgcol)
// 			for _, c := range button.History() {
// 				drawInk(gtx, c)
// 			}
// 		}),
// 		layout.Stacked(func() {
// 			gtx.Constraints.Width.Min = hmin
// 			gtx.Constraints.Height.Min = vmin
// 			layout.Center.Layout(gtx, func() {
// 				b.Inset.Layout(gtx, func() {
// 					paint.ColorOp{Color: col}.Add(gtx.Ops)
// 					widget.Label{}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Text)
// 				})
// 			})
// 			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
// 			button.Layout(gtx)
// 		}),
// 	)
// }
