// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type CheckBoxStyle struct {
	theme *Theme
	material.CheckBoxStyle
}

func (t *Theme) CheckBox(checkBox *widget.Bool, label string) CheckBoxStyle {
	cb := CheckBoxStyle{theme: t, CheckBoxStyle: material.CheckBox(t.Base, checkBox, label)}
	cb.theme = t
	return cb
}

// Layout updates the checkBox and displays it.
func (c CheckBoxStyle) Layout(gtx layout.Context) layout.Dimensions {
	c.Color = c.theme.Color.Text
	c.IconColor = c.theme.Color.Primary
	return c.CheckBoxStyle.Layout(gtx)
}
