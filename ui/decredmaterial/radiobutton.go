// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"gioui.org/widget"
	"gioui.org/widget/material"
)

type RadioButton struct {
	material.RadioButtonStyle
}

// RadioButton returns a RadioButton with a label. The key specifies
// the value for the Enum.
func (t *Theme) RadioButton(group *widget.Enum, key, label string, color color.NRGBA) RadioButton {
	rb := RadioButton{material.RadioButton(t.Base, group, key, label)}
	rb.Color = color
	return rb
}
