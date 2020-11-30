package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

const (
	Vertical Axis = iota
	Horizontal
)

type Axis int

type Container struct {
	container *layout.List
	scrollbar *Scrollbar

	contentHeight float32
	contentWidth  float32
	axis          layout.Axis
}

// Container returns a list layout  with a visible scrollbar
func (t *Theme) Container(scrollAxis Axis) *Container {
	axis := layout.Vertical
	if scrollAxis == 1 {
		axis = layout.Horizontal
	}

	return &Container{
		container: &layout.List{
			Axis: axis,
		},
		scrollbar: t.Scrollbar(axis),
		axis:      axis,
	}
}

// calculateContentDims calculates the total height and width of the content to be displayed.
func (c *Container) calculateContentDims(gtx layout.Context, w []func(gtx C) D) {
	height := float32(0)
	width := float32(0)
	for i := range w {
		index := i
		(&layout.List{Axis: layout.Vertical}).Layout(gtx, 1, func(gtx C, i int) D {
			dim := layout.UniformInset(unit.Dp(5)).Layout(gtx, w[index])
			height += float32(dim.Size.Y)
			width += float32(dim.Size.X)
			return layout.Dimensions{}
		})
	}
	c.contentHeight = height
	c.contentWidth = width
}

func (c *Container) Layout(gtx layout.Context, w []func(gtx C) D) layout.Dimensions {
	return layout.Flex{Axis: c.axis}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			c.calculateContentDims(gtx, w)
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			return c.layout(gtx, w)
		}),
	)
}

func (c *Container) layout(gtx layout.Context, w []func(gtx C) D) layout.Dimensions {
	if scrolled, progress := c.scrollbar.Scrolled(); scrolled {
		c.container.Position.First = int(float32(len(w)) * progress)
	}

	var visibleFraction, scrollDepth float32
	scrollbarThickness := gtx.Px(unit.Dp(12))

	inset := layout.Inset{}
	if c.axis == layout.Vertical {
		inset.Right = unit.Dp(float32(scrollbarThickness + 10))
	}

	contentFunc := func(gtx C) D {
		return inset.Layout(gtx, func(gtx C) D {
			var totalVisibleHeight float32

			dims := c.container.Layout(gtx, len(w), func(gtx C, i int) D {
				dim := layout.UniformInset(unit.Dp(5)).Layout(gtx, w[i])
				maxLength := dim.Size.Y
				contentLength := c.contentHeight
				if c.axis == layout.Horizontal {
					maxLength = dim.Size.X
					contentLength = c.contentWidth
				}
				totalVisibleHeight += float32(maxLength)
				visibleFraction = totalVisibleHeight / contentLength
				scrollDepth = float32(c.container.Position.First) / float32(len(w))

				return dim
			})
			return dims
		})
	}

	scrollbarFunc := func(gtx C) D {
		if c.axis == layout.Vertical {
			gtx.Constraints.Max.X = scrollbarThickness
		} else {
			gtx.Constraints.Max.Y = scrollbarThickness * 40
		}
		return c.scrollbar.Layout(gtx, c.axis, scrollDepth, visibleFraction)
	}

	if c.axis == layout.Vertical {
		return layout.Stack{Alignment: layout.E}.Layout(gtx, layout.Stacked(contentFunc), layout.Stacked(scrollbarFunc))
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, layout.Rigid(contentFunc), layout.Rigid(scrollbarFunc))
}
