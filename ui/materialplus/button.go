package materialplus

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (t *Theme) DangerButton(text string) material.Button {
	btn := t.Button(text)
	btn.Background = t.Danger
	return btn
}

func IconLabel(gtx *layout.Context, flex layout.Flex, ic *material.Icon, lbl material.Label) {
	flex.Layout(gtx,
		layout.Rigid(func() { ic.Layout(gtx, unit.Dp(20)) }),
		layout.Rigid(func() { lbl.Layout(gtx) }),
	)
}
