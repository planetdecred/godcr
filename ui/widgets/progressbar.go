package widgets

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"image/color"

	"github.com/raedahgroup/godcr-gio/ui/values"
)

type (
	ProgressBar struct {
		height          int
		backgroundColor color.RGBA
		progressColor   color.RGBA
	}
)

// paintArea creates an overlay of two rectangles to make a progress bar.
// The first indicates the progress to be completed, the second indicates
// the level of progress
func paintArea(ctx *layout.Context, color color.RGBA, x int, y int) {
	borderRadius := float32(6)
	borderWidth := 1
	if y < 21 {
		borderRadius = float32(4)
		borderWidth = 0
	}

	clip.Rect{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(x),
				Y: float32(y),
			},
		},
		NE: borderRadius,
		NW: borderRadius,
		SE: borderRadius,
		SW: borderRadius,
	}.Op(ctx.Ops).Add(ctx.Ops)
	fill(ctx, values.Gray, x, y)

	innerWidth := x - borderWidth
	innerHeight := y - borderWidth

	clip.Rect{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(innerWidth),
				Y: float32(innerHeight),
			},
			Min: f32.Point{
				X: float32(borderWidth),
				Y: float32(borderWidth),
			},
		},
		NE: borderRadius,
		NW: borderRadius,
		SE: borderRadius,
		SW: borderRadius,
	}.Op(ctx.Ops).Add(ctx.Ops)
	fill(ctx, color, innerWidth, innerHeight)
}

func (p *ProgressBar) SetHeight(height int) *ProgressBar {
	p.height = height
	return p
}

func (p *ProgressBar) SetBackgroundColor(col color.RGBA) *ProgressBar {
	p.backgroundColor = col
	return p
}

func (p *ProgressBar) SetProgressColor(col color.RGBA) *ProgressBar {
	p.progressColor = col
	return p
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		height:          values.DefaultProgressBarHeight,
		backgroundColor: values.Gray,
		progressColor:   values.Green,
	}
}

func (p *ProgressBar) Layout(ctx *layout.Context, progress float64) {
	layout.Stack{}.Layout(ctx,
		layout.Stacked(func() {
			paintArea(ctx, p.backgroundColor, ctx.Constraints.Width.Max, p.height)
			// calculate width of indicator with respects to progress bar width
			indicatorWidth := progress / float64(100) * float64(ctx.Constraints.Width.Max)

			if indicatorWidth > float64(ctx.Constraints.Width.Max) {
				indicatorWidth = float64(ctx.Constraints.Width.Max)
			}

			paintArea(ctx, p.progressColor, int(indicatorWidth), p.height)
		}),
	)
}
