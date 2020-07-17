// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Label struct {
	material.LabelStyle
}

func (t *Theme) H1(txt string) Label {
	return Label{material.H1(t.Base, txt)}
}

func (t *Theme) H2(txt string) Label {
	return Label{material.H2(t.Base, txt)}
}

func (t *Theme) H3(txt string) Label {
	return Label{material.H3(t.Base, txt)}
}

func (t *Theme) H4(txt string) Label {
	return Label{material.H4(t.Base, txt)}
}

func (t *Theme) H5(txt string) Label {
	return Label{material.H5(t.Base, txt)}
}

func (t *Theme) H6(txt string) Label {
	return Label{material.H6(t.Base, txt)}
}

func (t *Theme) Body1(txt string) Label {
	return Label{material.Body1(t.Base, txt)}
}

func (t *Theme) Body2(txt string) Label {
	return Label{material.Body2(t.Base, txt)}
}

func (t *Theme) Caption(txt string) Label {
	return Label{material.Caption(t.Base, txt)}
}

func (t *Theme) ErrorLabel(txt string) Label {
	label := t.Caption(txt)
	label.Color = t.Color.Danger
	return label
}

func (t *Theme) Label(size unit.Value, txt string) Label {
	return Label{material.Label(t.Base, size, txt)}
}
