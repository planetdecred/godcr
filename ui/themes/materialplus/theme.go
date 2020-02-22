package materialplus

import (
	"image/color"

	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/ui"
)

// Theme is an extenstion of gio's material theme
type Theme struct {
	*material.Theme

	Danger, Disabled, Success, Secondary, Background color.RGBA
}

// NewTheme returns a new materialplus theme
func NewTheme() *Theme {
	t := material.NewTheme()
	t.Color.Primary = ui.LightBlueColor
	return &Theme{
		Theme:    t,
		Danger:   ui.DangerColor,
		Disabled: ui.GrayColor,
	}
}
