package materialplus

import (
	"image/color"
)

type Palette struct {
	Primary, Secondary, Tertiary, Accent color.RGBA
	Success, Danger, Warn                color.RGBA
	Background, Text, Disabled           color.RGBA
}
