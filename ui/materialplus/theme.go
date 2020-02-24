package materialplus

import (
	"gioui.org/widget/material"
)

// Theme is an extenstion of gio's material theme
type Theme struct {
	*material.Theme

	Palette
	Icon struct {
		Cancel, Logo, Check *material.Icon
	}
}

// NewTheme returns a new materialplus theme
func NewTheme(colors Palette) *Theme {
	t := material.NewTheme()
	if t == nil {
		return nil
	}
	t.Color.Primary = colors.Primary
	t.Color.Text = colors.Text
	return &Theme{
		Theme:   t,
		Palette: colors,
	}
}
