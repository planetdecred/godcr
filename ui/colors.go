package ui

import (
	"image/color"

	"github.com/raedahgroup/godcr-gio/ui/materialplus"
)

var (
	white     = color.RGBA{255, 255, 255, 255}
	black     = color.RGBA{0, 0, 0, 255}
	gray      = color.RGBA{128, 128, 128, 255}
	lightGray = color.RGBA{200, 200, 200, 205}

	dangerRed    = color.RGBA{215, 58, 73, 255}
	successGreen = color.RGBA{9, 227, 98, 255}

	lightBlue = color.RGBA{41, 112, 255, 255}

	orange = color.RGBA{237, 109, 71, 255}
	green  = color.RGBA{46, 214, 161, 255}

	keyBlue   = color.RGBA{0x29, 0x70, 0xFF, 255}
	turquiose = color.RGBA{0x2E, 0xD6, 0xA1, 255}
	darkBlue  = color.RGBA{0x09, 0x14, 0x40, 255}
)

var decredPalette = materialplus.Palette{
	Primary:   keyBlue,
	Secondary: turquiose,
	Success:   successGreen,
	Danger:    dangerRed,
	Tertiary:  darkBlue,
}
