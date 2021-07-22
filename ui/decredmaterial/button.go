// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"

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
