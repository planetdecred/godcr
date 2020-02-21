package materialplus

import (
	"image/color"

	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/ui"
)

// Theme is an extenstion of gio's material theme
type Theme struct {
	*material.Theme

	Danger   color.RGBA
	Disabled color.RGBA
	Primary  color.RGBA
	White    color.RGBA
}

// NewTheme returns a new materialplus theme
func NewTheme() *Theme {
	theme := &Theme{
		Theme: material.NewTheme(),
	}
	theme.setColors()

	return theme
}

func (t *Theme) setColors() {
	t.White = ui.WhiteColor
	t.Danger = ui.DangerColor
	t.Disabled = ui.GrayColor

	t.Color.Primary = ui.LightBlueColor
}
