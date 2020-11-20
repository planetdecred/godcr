package decredmaterial

import (
	//"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

const scrollbarThicknessFactor = 0.015

type Scroller struct {
	clickable widget.Clickable
	drag      gesture.Drag
	scrolled  bool
	length    int
	axis      layout.Axis
	// progress is how far from the top our scrollbar is expressed as a fraction between 0 and 1
	progress float32
}

// Update the internal state of the bar.
func (s *Scroller) Update(gtx C) {
	s.scrolled = false

	defer func() {
		if s.progress > 1 {
			s.progress = 1
		} else if s.progress < 0 {
			s.progress = 0
		}
	}()

	gestureAxis := gesture.Vertical
	if s.axis == layout.Horizontal {
		gestureAxis = gesture.Horizontal
	}

	if s.clickable.Clicked() {
		if presses := s.clickable.History(); len(presses) > 0 {
			press := presses[len(presses)-1]
			pressPosition := press.Position.Y
			if gestureAxis == gesture.Horizontal {
				pressPosition = press.Position.X
			}

			s.progress = pressPosition / float32(s.length)
			s.scrolled = true
		}
	}

	if drags := s.drag.Events(gtx.Metric, gtx, gestureAxis); len(drags) > 0 {
		drag := drags[len(drags)-1]

		delta := drag.Position.Y
		if gestureAxis == gesture.Horizontal {
			delta = drag.Position.X
		}

		s.progress = (s.progress*float32(s.length) + (delta / 2)) / float32(s.length)
		s.scrolled = true
	}
}

// Scrolled returns true if the scroll position changed within the last frame.
func (s Scroller) Scrolled() (didScroll bool, progress float32) {
	return s.scrolled, s.progress
}

type Scrollbar struct {
	*Scroller
	Progress  float32
	Scale     float32
	MinLength unit.Value

	trackColor color.RGBA
	thumbColor color.RGBA
}

func (t *Theme) Scrollbar(axis layout.Axis) *Scrollbar {
	return &Scrollbar{
		trackColor: t.Color.Gray,
		thumbColor: t.Color.Primary,
		Scroller: &Scroller{
			axis: axis,
		},
		MinLength: unit.Dp(15),
	}
}

func (s *Scrollbar) getMargins(gtx layout.Context, axis layout.Axis, progress, scale float32) (f32.Point, unit.Value, unit.Value) {
	s.Progress = progress
	s.Scale = scale

	s.Scroller.progress = s.Progress
	s.Update(gtx)
	if scrolled, _ := s.Scrolled(); scrolled {
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	maxScrollLength := gtx.Constraints.Max.Y
	if axis == layout.Horizontal {
		maxScrollLength = gtx.Constraints.Max.X
	}

	scaledLength := (s.Scale * float32(maxScrollLength))
	s.MinLength = unit.Dp(scaledLength / gtx.Metric.PxPerDp)

	s.length = maxScrollLength

	var size f32.Point
	if axis == layout.Vertical {
		size = f32.Point{
			X: float32(gtx.Px(unit.Dp(float32(gtx.Constraints.Max.X)))),
			Y: float32(gtx.Px(s.MinLength)),
		}
	} else {
		size = f32.Point{
			X: float32(gtx.Px(s.MinLength)),
			Y: scrollbarThicknessFactor * float32(gtx.Constraints.Max.Y),
		}
	}

	total := float32(gtx.Constraints.Max.X) / gtx.Metric.PxPerDp

	var top, left unit.Value
	if axis == layout.Vertical {
		left = unit.Dp(2)
		top = unit.Dp(total * s.Progress)
		if top.V+s.MinLength.V > total {
			top = unit.Dp(total - s.MinLength.V)
		}
	} else {
		top = unit.Dp(2)
		left = unit.Dp(total * s.Progress)
		if left.V+s.MinLength.V > total {
			left = unit.Dp(total - s.MinLength.V)
		}
	}

	return size, top, left
}

func (s *Scrollbar) Layout(gtx layout.Context, axis layout.Axis, progress, scale float32) layout.Dimensions {
	size, top, left := s.getMargins(gtx, axis, progress, scale)

	// don't display scrollbar if content height or width is equal to or less than viewport size
	if s.MinLength.V >= float32(s.length) {
		return layout.Dimensions{}
	}

	clickable := &s.clickable
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(clickable.Layout),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
			}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx C) D {
			dims := layout.Inset{
				Top:    top,
				Right:  unit.Dp(2),
				Left:   left,
				Bottom: unit.Dp(2),
			}.Layout(gtx, func(gtx C) D {
				pointer.Rect(image.Rectangle{
					Max: image.Point{
						X: int(size.X),
						Y: int(size.Y),
					},
				}).Add(gtx.Ops)
				s.drag.Add(gtx.Ops)
				return drawRect(gtx, s.thumbColor, size, float32(gtx.Px(unit.Dp(4))))
			})

			dims.Size.X = gtx.Constraints.Max.X
			return dims
		}),
	)
}

// drawRect creates a rectangle of the provided background color with
// Dimensions specified by size and a corner radius (on all corners)
// specified by radii.
func drawRect(gtx C, background color.RGBA, size f32.Point, radii float32) D {
	stack := op.Push(gtx.Ops)
	paintOp := paint.ColorOp{Color: background}
	paintOp.Add(gtx.Ops)
	bounds := f32.Rectangle{
		Max: size,
	}
	clip.RRect{
		Rect: bounds,
		NW:   radii,
		NE:   radii,
		SE:   radii,
		SW:   radii,
	}.Add(gtx.Ops)
	paint.PaintOp{
		Rect: bounds,
	}.Add(gtx.Ops)
	stack.Pop()
	return layout.Dimensions{Size: image.Pt(int(size.X), int(size.Y))}
}
