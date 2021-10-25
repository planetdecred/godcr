package decredmaterial

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"image"
)

type Scrollbar struct {
	track, indicator gesture.Click
	drag             gesture.Drag
	delta            float32
}

// Layout updates the internal state of the scrollbar based on events
// since the previous call to Layout. The provided axis will be used to
// normalize input event coordinates and constraints into an axis-
// independent format. viewportStart is the position of the beginning
// of the scrollable viewport relative to the underlying content expressed
// as a value in the range [0,1]. viewportEnd is the position of the end
// of the viewport relative to the underlying content, also expressed
// as a value in the range [0,1]. For example, if viewportStart is 0.25
// and viewportEnd is .5, the viewport described by the scrollbar is
// currently showing the second quarter of the underlying content.
func (s *Scrollbar) Layout(gtx layout.Context, axis layout.Axis, viewportStart float32) layout.Dimensions {
	// Calculate the length of the major axis of the scrollbar. This is
	// the length of the track within which pointer events occur, and is
	// used to scale those interactions.
	trackHeight := float32(axis.Convert(gtx.Constraints.Max).X)
	s.delta = 0

	// Jump to a click in the track.
	for _, event := range s.track.Events(gtx) {
		if event.Type != gesture.TypeClick ||
			event.Modifiers != key.Modifiers(0) ||
			event.NumClicks > 1 {
			continue
		}
		pos := axis.Convert(image.Point{
			X: int(event.Position.X),
			Y: int(event.Position.Y),
		})
		normalizedPos := float32(pos.X) / trackHeight
		s.delta += normalizedPos - viewportStart
	}

	// Offset to account for any drags.
	for _, event := range s.drag.Events(gtx.Metric, gtx, gesture.Axis(axis)) {
		if event.Type != pointer.Drag {
			continue
		}
		dragOffset := fConvert(axis, event.Position).X
		normalizedDragOffset := dragOffset / trackHeight
		s.delta += normalizedDragOffset - viewportStart
	}

	// Process events from the indicator so that hover is
	// detected properly.
	_ = s.indicator.Events(gtx)

	return layout.Dimensions{}
}

// AddTrack configures the track click listener for the scrollbar to use
// the most recently added pointer.AreaOp.
func (s *Scrollbar) AddTrack(ops *op.Ops) {
	s.track.Add(ops)
}

// AddIndicator configures the indicator click listener for the scrollbar to use
// the most recently added pointer.AreaOp.
func (s *Scrollbar) AddIndicator(ops *op.Ops) {
	s.indicator.Add(ops)
}

// AddDrag configures the drag listener for the scrollbar to use
// the most recently added pointer.AreaOp.
func (s *Scrollbar) AddDrag(ops *op.Ops) {
	s.drag.Add(ops)
}

// IndicatorHovered returns whether the scroll indicator is currently being
// hovered by the pointer.
func (s *Scrollbar) IndicatorHovered() bool {
	return s.indicator.Hovered()
}

// ScrollDistance returns the normalized distance that the scrollbar
// moved during the last call to Layout as a value in the range [-1,1].
func (s *Scrollbar) ScrollDistance() float32 {
	return s.delta
}

func fConvert(a layout.Axis, pt f32.Point) f32.Point {
	if a == layout.Horizontal {
		return pt
	}
	return f32.Pt(pt.Y, pt.X)
}

// ListWithScroll holds the persistent state for a layout.List that has a
// scrollbar attached.
type ListWithScroll struct {
	Scrollbar
	LayoutList
}
