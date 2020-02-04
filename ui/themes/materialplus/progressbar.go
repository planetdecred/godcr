package materialplus

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"

	"github.com/raedahgroup/godcr-gio/ui/values"
)

type (
	ProgressBar struct {
		height          int
		backgroundColor color.RGBA
		progressColor   color.RGBA
	}
)

// tracks lays out a rectangle to represent the level of progress yet to be completed.
func (p *ProgressBar) track(gtx *layout.Context) {
	borderedRectangle(gtx, values.ProgressBarGray, gtx.Constraints.Width.Max, p.height)
}

// values lays out a rectangle to represent the level of progress that has been completed.
func (p *ProgressBar) value(gtx *layout.Context, progress float64) {
	width := progress / 100 * float64(gtx.Constraints.Width.Max)
	if width > float64(gtx.Constraints.Width.Max) {
		width = float64(gtx.Constraints.Width.Max)
	}
	borderedRectangle(gtx, p.progressColor, int(width), p.height)
}

// borderedRectangle defines the dimensions of the rectangle.
//func borderedRectangle(gtx *layout.Context, color color.RGBA, x, y int) {
func borderedRectangle(gtx *layout.Context, color color.RGBA, x, y int) {
	br := float32(y / 5)
	rect := f32.Rectangle{
		Max: f32.Point{
			X: float32(x),
			Y: float32(y),
		},
	}
	clip.Rect{
		Rect: rect,
		NE:   br, NW: br, SE: br, SW: br,
	}.Op(gtx.Ops).Add(gtx.Ops)

	Fill(gtx, color, x, y)
}

// SetHeight sets the height of the progress bar
func (p *ProgressBar) SetHeight(height int) *ProgressBar {
	p.height = height
	return p
}

// SetBackgroundColor sets the color of track of the progress bar
func (p *ProgressBar) SetBackgroundColor(col color.RGBA) *ProgressBar {
	p.backgroundColor = col
	return p
}

// SetProgressColor sets the color of the level of progress that has been completed.
func (p *ProgressBar) SetProgressColor(col color.RGBA) *ProgressBar {
	p.progressColor = col
	return p
}

// Layout lays out the track and level of progress on each other.
func (p *ProgressBar) Layout(gtx *layout.Context, progress float64) {
	layout.Stack{}.Layout(gtx,
		layout.Stacked(func() {
			p.track(gtx)
			p.value(gtx, progress)
		}),
	)
}

// NewProgressBar creates a new ProgressBar object.
func (t *Theme) ProgressBar() *ProgressBar {
	return &ProgressBar{
		height:          values.DefaultProgressBarHeight,
		backgroundColor: values.ProgressBarGray,
		progressColor:   values.ProgressBarGreen,
	}
}
