package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Collapsible struct {
	isExpanded            bool
	buttonWidget          *widget.Clickable
	line                  *Line
	expandedIcon          *widget.Icon
	collapsedIcon         *widget.Icon
	headerBackgroundColor color.RGBA
}

func (t *Theme) Collapsible() *Collapsible {
	c := &Collapsible{
		isExpanded:            false,
		headerBackgroundColor: t.Color.Hint,
		expandedIcon:          t.chevronUpIcon,
		collapsedIcon:         t.chevronDownIcon,
		line:                  t.Line(),
		buttonWidget:          new(widget.Clickable),
	}

	c.line.Color = t.Color.Gray
	c.line.Color.A = 140

	return c
}

func (c *Collapsible) layoutHeader(gtx layout.Context, header func(C) D) layout.Dimensions {
	icon := c.collapsedIcon
	if c.isExpanded {
		icon = c.expandedIcon
	}

	dims := layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return header(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: unit.Dp(20)}.Layout(gtx, func(C) D {
				return icon.Layout(gtx, unit.Dp(20))
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return c.buttonWidget.Layout(gtx)
}

func (c *Collapsible) Layout(gtx layout.Context, header func(C) D, content func(C) D) layout.Dimensions {
	for c.buttonWidget.Clicked() {
		c.isExpanded = !c.isExpanded
	}

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			c.line.Width = gtx.Constraints.Max.X
			return c.line.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
				return c.layoutHeader(gtx, header)
			})
		}),
		layout.Rigid(func(gtx C) D {
			if c.isExpanded {
				return content(gtx)
			}
			return layout.Dimensions{}
		}),
	)
	return dims
}
