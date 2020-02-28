package values

import "image/color"

var (
	// ProgressBarGray indicates the level of sync progress that is yet to be completed.
	ProgressBarGray = color.RGBA{230, 234, 237, 255}

	// ProgressBarGreen indicates the level of sync progress that has been completed.
	ProgressBarGreen = color.RGBA{65, 190, 83, 255}

	// White is the RGBA value for white color
	White = color.RGBA{255, 255, 255, 255}

	// DefaultCardGray is the background color of the card widget.
	DefaultCardGray = color.RGBA{243, 245, 246, 255}

	// ButtonGray is the RGBA value for gray color used on buttons
	ButtonGray = color.RGBA{196, 203, 210, 255}

	// ButtonRed is the RGBA value for red color used on buttons
	ButtonRed = color.RGBA{237, 109, 71, 255}

	// TextGray is the RGBA value for light texts on the app
	TextGray = color.RGBA{137, 151, 165, 255}
)
