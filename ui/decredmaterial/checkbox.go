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
	cb := CheckBoxStyle{material.CheckBox(t.Base, checkBox, label)}
	cb.Color = t.Color.Text
	cb.IconColor = t.Color.Primary
	return cb
}

// Layout updates the checkBox and displays it.
func (c CheckBoxStyle) Layout(gtx layout.Context) layout.Dimensions {
	return c.CheckBoxStyle.Layout(gtx)
}
