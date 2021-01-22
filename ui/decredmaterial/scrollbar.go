// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

type Scrollbar struct {
	color    color.RGBA
	float    *Float
	position float32
}

func (t *Theme) Scrollbar(float *Float, min, max float32) Scrollbar {
	return Scrollbar{
		color: t.Color.Primary,
		float: float,
	}
}

func (s *Scrollbar) scrolled() bool {
	return s.float.Scrolled()
}

func (s *Scrollbar) Layout(gtx layout.Context, contentLength, initialPosition float32) layout.Dimensions {
	maxSize := gtx.Constraints.Max

	windowLength := float32(maxSize.Y)
	visiblePortion := windowLength / contentLength
	scrollbarLength := windowLength * visiblePortion

	s.position = initialPosition * float32(maxSize.Y)

	st := op.Push(gtx.Ops)
	op.Offset(f32.Pt(0, 0)).Add(gtx.Ops)
	s.float.Layout(gtx, int(scrollbarLength), int(contentLength))
	st.Pop()
	if s.float.Scrolled() {
		s.position = s.float.Pos()
	}

	color := s.color
	if gtx.Queue == nil {
		color = mulAlpha(color, 150)
	}

	track := f32.Rectangle{
		Min: f32.Point{
			Y: 0,
		},
		Max: f32.Point{
			X: float32(maxSize.X),
			Y: float32(maxSize.Y),
		},
	}
	clip.RRect{Rect: track}.Add(gtx.Ops)
	fill(gtx, mulAlpha(color, 96))

	contentOffset := (s.position / float32(maxSize.Y)) * contentLength
	lengthOffset := int((windowLength - scrollbarLength) * (contentOffset / (contentLength - windowLength)))

	minY := float32(lengthOffset)
	maxY := minY + scrollbarLength
	if maxY > float32(maxSize.Y) {
		maxY = float32(maxSize.Y)
		minY = maxY - scrollbarLength
	}

	thumb := f32.Rectangle{
		Min: f32.Point{
			Y: minY,
		},
		Max: f32.Point{
			Y: maxY,
			X: float32(maxSize.X),
		},
	}
	rr := 0.5 * float32(maxSize.X)
	clip.RRect{
		Rect: thumb,
		NE:   rr, NW: rr, SE: rr, SW: rr,
	}.Add(gtx.Ops)
	fill(gtx, color)

	return layout.Dimensions{Size: maxSize}
}

func (s *Scrollbar) Position() float32 {
	return s.position
}

type ScrollContainer struct {
	container                       *layout.List
	scrollbar                       Scrollbar
	totalContentHeight              float32
	hasCalculatedTotalContentHeight bool
	scrollbarThicknessPercentage    float32
}

func (t *Theme) ScrollContainer() *ScrollContainer {
	return &ScrollContainer{
		container:                       &layout.List{Axis: layout.Vertical},
		scrollbar:                       t.Scrollbar(new(Float), 0, 100),
		totalContentHeight:              0,
		hasCalculatedTotalContentHeight: false,
		scrollbarThicknessPercentage:    1.3,
	}
}

func (s *ScrollContainer) calculateTotalContentHeight(gtx layout.Context, widgets []func(gtx C) D) {
	for i := range widgets {
		index := i
		(&layout.List{Axis: layout.Vertical}).Layout(gtx, 1, func(gtx C, i int) D {
			dim := layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[index])
			s.totalContentHeight += float32(dim.Size.Y)
			return layout.Dimensions{}
		})
	}
	s.hasCalculatedTotalContentHeight = true
}

func (s *ScrollContainer) Layout(gtx layout.Context, widgets []func(gtx C) D) layout.Dimensions {
	maxSize := gtx.Constraints.Max

	if s.scrollbar.scrolled() {
		s.container.Position.First = int(float32(len(widgets)) * (s.scrollbar.Position() / float32(maxSize.Y)))
	}

	totalVisibleHeight := float32(maxSize.Y)
	scrollbarThickness := (s.scrollbarThicknessPercentage / 100) * float32(1)

	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1-scrollbarThickness, func(gtx C) D {
			if !s.hasCalculatedTotalContentHeight {
				s.calculateTotalContentHeight(gtx, widgets)
				return layout.Dimensions{}
			}

			return layout.Inset{
				Right: unit.Dp(15),
			}.Layout(gtx, func(gtx C) D {
				return s.container.Layout(gtx, len(widgets), func(gtx C, i int) D {
					return layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
				})
			})
		}),
		layout.Flexed(scrollbarThickness, func(gtx C) D {
			// don't display scrollbar if total content height is less than or equal to container height
			if s.totalContentHeight <= totalVisibleHeight {
				return layout.Dimensions{}
			}

			initialPosition := float32(s.container.Position.First) / float32(len(widgets))
			return s.scrollbar.Layout(gtx, s.totalContentHeight, initialPosition)
		}),
	)
}
