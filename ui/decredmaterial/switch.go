// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Switch struct {
	material.SwitchStyle
}

// RadioButton returns a RadioButton with a label. The key specifies
// the value for the Enum.
func (t *Theme) Switch(swtch *widget.Bool) Switch {
	return Switch{material.Switch(t.Base, swtch)}
}
