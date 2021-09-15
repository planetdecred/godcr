package decredmaterial

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func Clickable(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	return material.Clickable(gtx, button, w)
}

type Cllickable struct {
	button     *widget.Clickable
	color      color.NRGBA
	HoverColor color.NRGBA
	Hoverable  bool
	Radius     CornerRadius
}

func (t *Theme) NewClickable(hoverable bool) *Cllickable {
	return &Cllickable{
		button:     &widget.Clickable{},
		color:      t.Color.SurfaceHighlight,
		HoverColor: t.Color.ActiveGray,
		Hoverable:  hoverable,
	}
}

func (cl *Cllickable) Clicked() bool {
	return cl.button.Clicked()
}

func (cl *Cllickable) Layout(gtx C, w layout.Widget) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(cl.button.Layout),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			tr := float32(gtx.Px(unit.Dp(cl.Radius.TopRight)))
			tl := float32(gtx.Px(unit.Dp(cl.Radius.TopLeft)))
			br := float32(gtx.Px(unit.Dp(cl.Radius.BottomRight)))
			bl := float32(gtx.Px(unit.Dp(cl.Radius.BottomLeft)))
			clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
				NW: tl, NE: tr, SE: br, SW: bl,
			}.Add(gtx.Ops)
			clip.Rect{Max: gtx.Constraints.Min}.Add(gtx.Ops)

			if cl.Hoverable && cl.button.Hovered() {
				paint.Fill(gtx.Ops, cl.HoverColor)
			}

			for _, c := range cl.button.History() {
				drawInk(gtx, c, cl.color)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}
