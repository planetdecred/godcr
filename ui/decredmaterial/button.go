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

func (t *Theme) IconButton(icon *widget.Icon, button *widget.Clickable) IconButton {
	return IconButton{material.IconButton(t.Base, button, icon)}
}

func (t *Theme) PlainIconButton(icon *widget.Icon, button *widget.Clickable) IconButton {
	btn := IconButton{material.IconButton(t.Base, button, icon)}
	btn.Background = color.RGBA{}
	return btn
}

func (b Button) Layout(gtx layout.Context) layout.Dimensions {
	return b.ButtonStyle.Layout(gtx)
}

func (b IconButton) Layout(gtx layout.Context) layout.Dimensions {
	return b.IconButtonStyle.Layout(gtx)
}
