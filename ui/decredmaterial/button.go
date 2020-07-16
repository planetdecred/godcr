// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget/material"

	"gioui.org/widget"
)

type Button struct {
	material.ButtonStyle
}

type IconButton struct {
	material.IconButtonStyle
}

func (t *Theme) Button(button *widget.Clickable, txt string) Button {
	return Button{material.Button(t.Base, button, txt)}
}

func (t *Theme) IconButton(button *widget.Clickable, icon *widget.Icon) IconButton {
	return IconButton{material.IconButton(t.Base, button, icon)}
}

func (t *Theme) PlainIconButton(button *widget.Clickable, icon *widget.Icon) IconButton {
	btn := IconButton{material.IconButton(t.Base, button, icon)}
	btn.Background = color.RGBA{}
	return btn
}

func Clickable(gtx layout.Context, button *widget.Clickable, w layout.Widget) layout.Dimensions {
	return material.Clickable(gtx, button, w)
}

func (b Button) Layout(gtx layout.Context) layout.Dimensions {
	return b.ButtonStyle.Layout(gtx)
}

func (b IconButton) Layout(gtx layout.Context) layout.Dimensions {
	return b.IconButtonStyle.Layout(gtx)
}
