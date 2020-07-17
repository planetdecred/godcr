// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"gioui.org/widget/material"
)

type ProgressBarStyle struct {
	material.ProgressBarStyle
}

func (t *Theme) ProgressBar(progress int) ProgressBarStyle {
	return ProgressBarStyle{material.ProgressBar(t.Base, progress)}
}
