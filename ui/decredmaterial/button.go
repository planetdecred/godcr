// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"
	"math"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/op/clip"
	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"

	"gioui.org/widget"
)

type Button struct {
	material.ButtonStyle
	isEnabled          bool
	disabledBackground color.NRGBA
	surfaceColor       color.NRGBA
}

type IconButton struct {
	material.IconButtonStyle
}

func (t *Theme) Button(button *widget.Clickable, txt string) Button {
	return Button{
		ButtonStyle:        material.Button(t.Base, button, txt),
		disabledBackground: t.Color.Gray,
		surfaceColor:       t.Color.Surface,
		isEnabled:          true,
	}
}

func (t *Theme) IconButton(button *widget.Clickable, icon *widget.Icon) IconButton {
	return IconButton{material.IconButton(t.Base, button, icon)}
}

func (t *Theme) PlainIconButton(button *widget.Clickable, icon *widget.Icon) IconButton {
	btn := IconButton{material.IconButton(t.Base, button, icon)}
	btn.Background = color.NRGBA{}
	return btn
}

func Clickable(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	return material.Clickable(gtx, button, w)
}

func Clickable2(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(button.Layout),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			// clip.Rect{Max: gtx.Constraints.Min}.Add(gtx.Ops)
			rr := float32(gtx.Px(unit.Dp(4)))
			clip.UniformRRect(f32.Rectangle{Max: f32.Point{
				X: float32(gtx.Constraints.Min.X),
				Y: float32(gtx.Constraints.Min.Y),
			}}, rr).Add(gtx.Ops)
			// background := b.Background
			// switch {
			// case gtx.Queue == nil:
			// 	background = Disabled(b.Background)
			// case b.Button.Hovered():
			// 	background = Hovered(b.Background)
			// }
			// paint.Fill(gtx.Ops, background)
			for _, c := range button.History() {
				layoutInfoTooltip(gtx, c)
			}
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(w),
	)
}

func layoutInfoTooltip(gtx layout.Context, c widget.Press) {
	layout.Inset{}.Layout(gtx, func(gtx C) D {
		return drawInk(gtx, c)
	})
}

func (b *Button) SetEnabled(enabled bool) {
	b.isEnabled = enabled
}

func (b *Button) Enabled() bool {
	return b.isEnabled
}

func (b Button) Clicked() bool {
	return b.Button.Clicked()
}

func (b *Button) Layout(gtx layout.Context) layout.Dimensions {
	if !b.Enabled() {
		gtx.Queue = nil
		b.Background, b.Color = b.disabledBackground, b.surfaceColor
	}

	return b.ButtonStyle.Layout(gtx)
}

func (b IconButton) Layout(gtx layout.Context) layout.Dimensions {
	return b.IconButtonStyle.Layout(gtx)
}

type TextAndIconButton struct {
	theme           *Theme
	Button          *widget.Clickable
	icon            *widget.Icon
	text            string
	Color           color.NRGBA
	BackgroundColor color.NRGBA
}

func (t *Theme) TextAndIconButton(btn *widget.Clickable, text string, icon *widget.Icon) TextAndIconButton {
	return TextAndIconButton{
		theme:           t,
		Button:          btn,
		icon:            icon,
		text:            text,
		Color:           t.Color.Surface,
		BackgroundColor: t.Color.Primary,
	}
}

func (b TextAndIconButton) Layout(gtx layout.Context) layout.Dimensions {
	btnLayout := material.ButtonLayout(b.theme.Base, b.Button)
	btnLayout.Background = b.BackgroundColor
	b.icon.Color = b.Color

	return btnLayout.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					var d D
					size := gtx.Px(unit.Dp(46)) - 2*gtx.Px(unit.Dp(16))
					b.icon.Layout(gtx, unit.Px(float32(size)))
					d = layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					l := material.Body1(b.theme.Base, b.text)
					l.Color = b.Color
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layLabel, layIcon)
		})
	})
}

func drawInk(gtx layout.Context, c widget.Press) {
	// duration is the number of seconds for the
	// completed animation: expand while fading in, then
	// out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := gtx.Now

	t := float32(now.Sub(c.Start).Seconds())

	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}

	endt := float32(end.Sub(c.Start).Seconds())

	// Compute the fade-in/out position in [0;1].
	var alphat float32
	{
		var haste float32
		if c.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}

		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := float32(gtx.Constraints.Min.X)
	if h := float32(gtx.Constraints.Min.Y); h > size {
		size = h
	}
	// Cover the entire constraints min rectangle.
	size *= 2 * float32(math.Sqrt(2))
	// Apply curve values to size and color.
	size *= sizeBezier
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(col*0xff)
	defer op.Save(gtx.Ops).Load()
	rgba := mulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	// rr := size * .5
	// op.Offset(c.Position.Add(f32.Point{
	// 	X: -rr,
	// 	Y: -rr,
	// })).Add(gtx.Ops)
	// tr := float32(gtx.Px(unit.Dp(4)))
	// tl := float32(gtx.Px(unit.Dp(4)))
	// br := float32(gtx.Px(unit.Dp(4)))
	// bl := float32(gtx.Px(unit.Dp(4)))
	// clip.RRect{
	// 	Rect: f32.Rectangle{Max: f32.Point{
	// 		X: float32(gtx.Constraints.Min.X),
	// 		Y: float32(gtx.Constraints.Min.Y),
	// 	}},
	// 	NW: tl, NE: tr, SE: br, SW: bl,
	// }.Add(gtx.Ops)
	// clip.UniformRRect(f32.Rectangle{Max: f32.Pt(size, size)}, rr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}