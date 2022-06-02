// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"
	_ "image/png" //makes png images a decodable format

	"gioui.org/unit"

	"gioui.org/widget"
)

type Icon struct {
	*widget.Icon
	Color color.NRGBA
}

// NewIcon returns a new Icon from IconVG data.
func NewIcon(icon *widget.Icon) *Icon {
	return &Icon{
		Icon: icon,
	}
}

func (icon *Icon) Layout(gtx C, iconSize unit.Dp) D {
	cl := color.NRGBA{A: 0xff}
	if icon.Color != (color.NRGBA{}) {
		cl = icon.Color
	}
	gtx.Constraints.Min.X = gtx.Dp(iconSize)
	return icon.Icon.Layout(gtx, cl)
}
