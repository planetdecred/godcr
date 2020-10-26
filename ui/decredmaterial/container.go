package decredmaterial

import (
	//"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
)

type Container struct {
	container *layout.List
	scrollbar *Scrollbar

	contentHeight float32
}

// Container returns a list layout  with a visible scrollbar
func (t *Theme) Container() *Container {
	return &Container{
		container: &layout.List{
			Axis: layout.Vertical,
		},
		scrollbar: t.Scrollbar(),
	}
}

// calculateContentHeight calculates the total height of the content to be displayed.
// if the total content height is longer than the window viewport size, a scrollbar is displayed.
func (c *Container) calculateContentHeight(gtx layout.Context, w []func(gtx C) D) {
	height := float32(0)
	for i := range w {
		index := i
		(&layout.List{Axis: layout.Vertical}).Layout(gtx, 1, func(gtx C, i int) D {
			dim := layout.UniformInset(unit.Dp(5)).Layout(gtx, w[index])
			height += float32(dim.Size.Y)
			return layout.Dimensions{}
		})
	}
	c.contentHeight = height
}

func (c *Container) Layout(gtx layout.Context, w []func(gtx C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			c.calculateContentHeight(gtx, w)
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

	scrollbarWidth := 0.015 * float32(gtx.Constraints.Max.X)

	var visibleFraction, scrollDepth float32
	return layout.Stack{Alignment: layout.E}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Inset{
				Right: unit.Dp(scrollbarWidth + 10),
			}.Layout(gtx, func(gtx C) D {
				var totalVisibleHeight float32

				dims := c.container.Layout(gtx, len(w), func(gtx C, i int) D {
					dim := layout.UniformInset(unit.Dp(5)).Layout(gtx, w[i])
					totalVisibleHeight += float32(dim.Size.Y)
					return dim
				})

				visibleFraction = totalVisibleHeight / c.contentHeight
				scrollDepth = float32(c.container.Position.First) / float32(len(w))

				return dims
			})
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Max.X = int(scrollbarWidth)
			return c.scrollbar.Layout(gtx, scrollDepth, visibleFraction)
		}),
	)
}
