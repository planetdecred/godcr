// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Label struct {
	material.LabelStyle
}

func (t *Theme) H1(txt string) Label {
	return t.labelWithDefaultColor(Label{material.H1(t.Base, txt)})
}

func (t *Theme) H2(txt string) Label {
	return t.labelWithDefaultColor(Label{material.H2(t.Base, txt)})
}

func (t *Theme) H3(txt string) Label {
	return t.labelWithDefaultColor(Label{material.H2(t.Base, txt)})
}

func (t *Theme) H4(txt string) Label {
	return t.labelWithDefaultColor(Label{material.H4(t.Base, txt)})
}

func (t *Theme) H5(txt string) Label {
	return t.labelWithDefaultColor(Label{material.H5(t.Base, txt)})
}

func (t *Theme) H6(txt string) Label {
	return t.labelWithDefaultColor(Label{material.H6(t.Base, txt)})
}

func (t *Theme) Body1(txt string) Label {
	return t.labelWithDefaultColor(Label{material.Body1(t.Base, txt)})
}

func (t *Theme) Body2(txt string) Label {
	return t.labelWithDefaultColor(Label{material.Body2(t.Base, txt)})
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

func (t *Theme) LabelColor(size unit.Value, txt string, color color.NRGBA) Label {
	return t.labelWithColor(Label{material.Label(t.Base, size, txt)}, color)
}

func (t *Theme) labelWithColor(l Label, color color.NRGBA) Label {
	l.Color = color
	return l
}

func (t *Theme) labelWithDefaultColor(l Label) Label {
	l.Color = t.Color.DeepBlue
	return l
}
