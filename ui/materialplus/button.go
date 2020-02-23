package materialplus

import (
	"gioui.org/widget/material"
)

func (t *Theme) DangerButton(text string) material.Button {
	btn := t.Button(text)
	btn.Background = t.Danger
	return btn
}
