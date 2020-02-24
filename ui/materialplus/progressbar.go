package materialplus

import (
	"image/color"

	"gioui.org/layout"
	"github.com/raedahgroup/godcr-gio/ui/materialplus/layouts"
)

// ProgressBar indicates the progress of a process. Height defines the thickness of the progressbar,
// BackgroundColor defines the color of the track, ProgressColor defines the color of the moving progress.
type ProgressBar struct {
	Progress float64
	Color    color.RGBA
}

// // track lays out a rectangle to represent the level of progress yet to be completed.
// func (p *ProgressBar) track(gtx *layout.Context) {
// 	borderedRectangle(gtx, values.ProgressBarGray, gtx.Constraints.Width.Max, p.Height)
// }

// // value lays out a rectangle to represent the level of progress that has been completed.
// func (p *ProgressBar) value(gtx *layout.Context, progress float64) {
// 	width := progress / 100 * float64(gtx.Constraints.Width.Max)
// 	if width > float64(gtx.Constraints.Width.Max) {
// 		width = float64(gtx.Constraints.Width.Max)
// 	}
// 	borderedRectangle(gtx, p.ProgressColor, int(width), p.Height)
// }

// Layout lays out the track and level of progress on each other.

// ProgressBar returns a new ProgressBar instance.
func (t *Theme) ProgressBar(progress float64) ProgressBar {
	return ProgressBar{
		Color:    t.Primary,
		Progress: progress,
	}
}

func (p ProgressBar) Layout(gtx *layout.Context) {
	layouts.FillWithColor(gtx, p.Color)
}
