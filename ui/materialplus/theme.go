package materialplus

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Theme is an extenstion of gio's material theme
type Theme struct {
	*material.Theme

	Palette
	Icon struct {
		Cancel, Logo, Check, Add *material.Icon
	}
}

// NewTheme returns a new materialplus theme
func NewTheme(colors Palette) *Theme {
	t := material.NewTheme()
	if t == nil {
		return nil
	}
	t.Color.Primary = colors.Primary
	return &Theme{
		Theme:   t,
		Palette: colors,
	}
}

func (t *Theme) Background(gtx *layout.Context, w layout.Widget) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillWithColor(gtx, ARGB(0x22444444))
		}),
		layout.Stacked(w),
	)
}

func (t *Theme) Surface(gtx *layout.Context, w layout.Widget) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillWithColor(gtx, RGB(0xffffff))
		}),
		layout.Stacked(w),
	)
}
