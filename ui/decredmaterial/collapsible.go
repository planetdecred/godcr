package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Collapsible struct {
	IsExpanded      bool
	Button          *widget.Clickable
	BackgroundColor color.NRGBA
	expandedIcon    *widget.Image
	collapsedIcon   *widget.Image
	card            Card
}

func (t *Theme) Collapsible() *Collapsible {
	c := &Collapsible{
		BackgroundColor: t.Color.Surface,
		Button:          new(widget.Clickable),
		card:            t.Card(),
		collapsedIcon:   t.collapseIcon,
		expandedIcon:    t.expandIcon,
	}
	c.card.Color = c.BackgroundColor
	c.collapsedIcon.Scale, c.expandedIcon.Scale = 1, 1

	return c
}

func (c *Collapsible) getHeaderParts(header, option func(C) D) []layout.FlexChild {
	icon := c.collapsedIcon
	if c.IsExpanded {
		icon = c.expandedIcon
	}

	children := []layout.FlexChild{
		layout.Flexed(1, func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return icon.Layout(gtx)
						}),
						layout.Rigid(header),
					)
				}),
				layout.Expanded(c.Button.Layout),
			)
		}),
	}

	if option != nil {
		children = append(children, layout.Rigid(option))
	}

	return children
}

func (c *Collapsible) Layout(gtx C, header, content, option func(C) D) layout.Dimensions {
	for c.Button.Clicked() {
		c.IsExpanded = !c.IsExpanded
	}

	children := []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, c.getHeaderParts(header, option)...)
			})
		}),
	}

	if c.IsExpanded {
		children = append(children, layout.Rigid(content))
	}

	return c.card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}
