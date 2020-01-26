package materialplus

import (
	"gioui.org/widget/material"
)

// Theme is an extenstion of gio's material theme
type Theme struct {
	*material.Theme
}

// NewTheme returns a new materialplus theme
func NewTheme() *Theme {
	return &Theme{
		Theme: material.NewTheme(),
	}
}
