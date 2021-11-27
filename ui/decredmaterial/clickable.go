package decredmaterial

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Clickable struct {
	button     *widget.Clickable
	Color      color.NRGBA
	HoverColor color.NRGBA
	Hoverable  bool
	Radius     CornerRadius
}

func (t *Theme) NewClickable(hoverable bool) *Clickable {
	// TODO: Temp fix until dark mode colors are sorted out.
	color := t.Color.Gray4
	if !t.DarkMode {
		color = Hovered(color)
	}

	return &Clickable{
		button:     &widget.Clickable{},
		Color:      t.Color.SurfaceHighlight,
		HoverColor: color,
		Hoverable:  hoverable,
	}
}

func (cl *Clickable) Clicked() bool {
	return cl.button.Clicked()
}

func (cl *Clickable) Layout(gtx C, w layout.Widget) D {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(cl.button.Layout),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			tr := float32(gtx.Px(unit.Dp(cl.Radius.TopRight)))
			tl := float32(gtx.Px(unit.Dp(cl.Radius.TopLeft)))
			br := float32(gtx.Px(unit.Dp(cl.Radius.BottomRight)))
			bl := float32(gtx.Px(unit.Dp(cl.Radius.BottomLeft)))
			defer clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
				NW: tl, NE: tr, SE: br, SW: bl,
			}.Push(gtx.Ops).Pop()
			clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()

			if cl.Hoverable && cl.button.Hovered() {
				paint.Fill(gtx.Ops, cl.HoverColor)
			}

			for _, c := range cl.button.History() {
				drawInk(gtx, c, cl.Color)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}
