// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type RadioButton struct {
	material.RadioButtonStyle
}

// RadioButton returns a RadioButton with a label. The key specifies
// the value for the Enum.
func (t *Theme) RadioButton(group *widget.Enum, key, label string) RadioButton {
	return RadioButton{material.RadioButton(t.Base, group, key, label)}
}
