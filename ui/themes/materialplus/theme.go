package materialplus

import (
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui"
)

// Theme is an extenstion of gio's material theme
type Theme struct {
	*material.Theme
}

// NewTheme returns a new materialplus theme
func NewTheme() *Theme {
	theme := &Theme{
		material.NewTheme(),
	}
	theme.Color.Primary = ui.LightBlueColor

	return theme
}
