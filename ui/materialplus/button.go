package materialplus

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// DangerButton a button with the background set to theme.Danger
func (t *Theme) DangerButton(text string) material.Button {
	btn := t.Button(text)
	btn.Background = t.Danger
	return btn
}

// IconLabel layouts flex with Rigid ic and Rigid lbl
func IconLabel(gtx *layout.Context, flex layout.Flex, ic *material.Icon, lbl material.Label) {
	flex.Layout(gtx,
		layout.Rigid(func() { ic.Layout(gtx, unit.Dp(20)) }),
		layout.Rigid(func() { lbl.Layout(gtx) }),
	)
}
