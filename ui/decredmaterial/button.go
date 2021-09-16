// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/values"
)

type Button struct {
	material.ButtonStyle
	label              Label
	clickable          *widget.Clickable
	isEnabled          bool
	disabledBackground color.NRGBA
	disabledTextColor  color.NRGBA
	HighlightColor     color.NRGBA

	Margin layout.Inset
}

type ButtonLayout struct {
	material.ButtonLayoutStyle
}

type IconButton struct {
	material.IconButtonStyle
}

func (t *Theme) Button(button *widget.Clickable, txt string) Button {
	buttonStyle := material.Button(t.Base, button, txt)
	buttonStyle.TextSize = values.TextSize16
	buttonStyle.Background = t.Color.Primary
	buttonStyle.CornerRadius = values.MarginPadding8
	buttonStyle.Inset = layout.Inset{
		Top:    values.MarginPadding10,
		Right:  values.MarginPadding16,
		Bottom: values.MarginPadding10,
		Left:   values.MarginPadding16,
	}
	return Button{
		ButtonStyle:        buttonStyle,
		label:              t.Label(values.TextSize16, txt),
		clickable:          button,
		disabledBackground: t.Color.Gray,
		disabledTextColor:  t.Color.Surface,
		HighlightColor:     t.Color.PrimaryHighlight,
		isEnabled:          true,
	}
}

func (t *Theme) OutlineButton(button *widget.Clickable, txt string) Button {
	btn := t.Button(button, txt)
	btn.Background = color.NRGBA{}
	btn.HighlightColor = t.Color.SurfaceHighlight
	btn.Color = t.Color.Primary
	btn.disabledBackground = color.NRGBA{}
	btn.disabledTextColor = t.Color.InactiveGray

	return btn
}

func (t *Theme) ButtonLayout(button *widget.Clickable) ButtonLayout {
	return ButtonLayout{material.ButtonLayout(t.Base, button)}
}

func (t *Theme) IconButton(button *widget.Clickable, icon *widget.Icon) IconButton {
	return IconButton{material.IconButton(t.Base, button, icon)}
}

func (t *Theme) PlainIconButton(button *widget.Clickable, icon *widget.Icon) IconButton {
	btn := IconButton{material.IconButton(t.Base, button, icon)}
	btn.Background = color.NRGBA{}
	return btn
}

func (b *Button) SetEnabled(enabled bool) {
	b.isEnabled = enabled
}

func (b *Button) Enabled() bool {
	return b.isEnabled
}

func (b Button) Clicked() bool {
	return b.clickable.Clicked()
}

func (b Button) Hovered() bool {
	return b.clickable.Hovered()
}

func (b Button) Click() {
	b.clickable.Click()
}

func (b *Button) Layout(gtx layout.Context) layout.Dimensions {
	wdg := func(gtx layout.Context) layout.Dimensions {
		return b.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			textColor := b.Color
			if !b.Enabled() {
				textColor = b.disabledTextColor
			}

			b.label.Text = b.Text
			b.label.Font = b.Font
			b.label.Alignment = text.Middle
			b.label.TextSize = b.TextSize
			b.label.Color = textColor
			return b.label.Layout(gtx)
		})
	}

	return b.Margin.Layout(gtx, func(gtx C) D {
		return b.buttonStyleLayout(gtx, wdg)
	})
}

func (b Button) buttonStyleLayout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	min := gtx.Constraints.Min
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rr := float32(gtx.Px(b.CornerRadius))
			clip.UniformRRect(f32.Rectangle{Max: f32.Point{
				X: float32(gtx.Constraints.Min.X),
				Y: float32(gtx.Constraints.Min.Y),
			}}, rr).Add(gtx.Ops)

			background := b.Background
			if !b.Enabled() {
				background = b.disabledBackground
			} else if b.clickable.Hovered() {
				background = Hovered(b.Background)
			}

			paint.Fill(gtx.Ops, background)
			for _, c := range b.clickable.History() {
				drawInk(gtx, c, b.HighlightColor)
			}

			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = min
			return layout.Center.Layout(gtx, w)
		}),
		layout.Expanded(func(gtx C) D {
			if !b.Enabled() {
				return D{}
			}

			return b.clickable.Layout(gtx)
		}),
	)
}

func (bl ButtonLayout) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return bl.ButtonLayoutStyle.Layout(gtx, w)
}

func (ib IconButton) Layout(gtx layout.Context) layout.Dimensions {
	return ib.IconButtonStyle.Layout(gtx)
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
