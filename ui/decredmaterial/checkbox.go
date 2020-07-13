// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CheckBoxStyle struct {
	material.CheckBoxStyle
}

func (t *Theme) CheckBox(checkBox *widget.Bool, label string) CheckBoxStyle {
	return CheckBoxStyle{material.CheckBox(t.Base, checkBox, label)}
}

// Layout updates the checkBox and displays it.
func (c CheckBoxStyle) Layout(gtx layout.Context) layout.Dimensions {
	return c.CheckBoxStyle.Layout(gtx)
}
