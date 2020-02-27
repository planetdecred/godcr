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
	fillWithColor(gtx, t.Tertiary)
	w()
}

func (t *Theme) Foreground(gtx *layout.Context, w layout.Widget) {
	fillWithColor(gtx, t.Secondary)
	w()
}
