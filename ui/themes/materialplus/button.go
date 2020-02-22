package materialplus

import (
	"image/color"

	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui"

	"gioui.org/widget/material"
)

func (t *Theme) ButtonWithColor(text string, color color.RGBA) material.Button {
	button := t.Button(text)
	button.Background = color
	return button
}

//
func (t *Theme) LayoutWithBackGround(gtx *layout.Context, block bool, widget func()) {
	ui.LayoutWithBackGround(gtx, t.Color.Primary, block, widget)
}
