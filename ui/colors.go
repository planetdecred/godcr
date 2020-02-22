package ui

import "image/color"

var (
	WhiteColor     = color.RGBA{255, 255, 255, 255}
	BlackColor     = color.RGBA{0, 0, 0, 255}
	GrayColor      = color.RGBA{128, 128, 128, 255}
	LightGrayColor = color.RGBA{200, 200, 200, 205}

	DangerColor  = color.RGBA{215, 58, 73, 255}
	SuccessColor = color.RGBA{227, 98, 9, 255}

	DarkBlueColor  = color.RGBA{9, 20, 64, 255}
	LightBlueColor = color.RGBA{41, 112, 255, 255}

	OrangeColor = color.RGBA{237, 109, 71, 255}
	GreenColor  = color.RGBA{46, 214, 161, 255} //color.RGBA{65, 191, 83, 255}

	BackgroundColor = color.RGBA{248, 249, 250, 255}

	FadedColor = color.RGBA{0, 0, 0, 128}
)

// Faded redues the color Aplha by half
func Faded(color color.RGBA) color.RGBA {
	color.A = color.A / 2
	return color
}
