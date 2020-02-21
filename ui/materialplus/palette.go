package materialplus

import (
	"image/color"
)

// Palette represensts a set of colors used by a Theme
type Palette struct {
	Primary, Secondary, Tertiary, Accent color.RGBA
	Success, Danger, Warn                color.RGBA
	Disabled                             color.RGBA
}
