package decredmaterial

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/values"
)

type Clickable struct {
	button    *widget.Clickable
	style     *values.ClickableStyle
	Hoverable bool
	Radius    CornerRadius
	isEnabled bool
}

func (t *Theme) NewClickable(hoverable bool) *Clickable {
	return &Clickable{
		button:    &widget.Clickable{},
		style:     t.Styles.ClickableStyle,
		Hoverable: hoverable,
		isEnabled: true,
	}
}

func (cl *Clickable) Style() values.ClickableStyle {
	return *cl.style
}

func (cl *Clickable) ChangeStyle(style *values.ClickableStyle) {
	cl.style = style
}

func (cl *Clickable) Clicked() bool {
	return cl.button.Clicked()
}

// SetEnabled enables/disables the clickable.
func (cl *Clickable) SetEnabled(enable bool, gtx *layout.Context) layout.Context {
	var mGtx layout.Context
	if gtx != nil && !enable {
		mGtx = gtx.Disabled()
	}

	cl.isEnabled = enable
	return mGtx
}

// Return clickable enabled/disabled state.
func (cl *Clickable) Enabled() bool {
	return cl.isEnabled
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
				paint.Fill(gtx.Ops, cl.style.HoverColor)
			}

			for _, c := range cl.button.History() {
				drawInk(gtx, c, cl.style.Color)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}
