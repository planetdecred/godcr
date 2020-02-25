package materialplus

import "image/color"

// ProgressBar fills the context with Background then fills it with Foreground.
type ProgressBar struct {
	Foreground, Background color.RGBA
}
