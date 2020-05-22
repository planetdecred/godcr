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
	buttonWidget          *widget.Button
	line                  *Line
	expandedIcon          *Icon
	collapsedIcon         *Icon
	headerBackgroundColor color.RGBA
}

func (t *Theme) Collapsible() *Collapsible {
	c := &Collapsible{
		isExpanded:            false,
		headerBackgroundColor: t.Color.Hint,
		expandedIcon:          t.chevronUpIcon,
		collapsedIcon:         t.chevronDownIcon,
		line:                  t.Line(),
		buttonWidget:          new(widget.Button),
	}

	c.line.Color = t.Color.Gray
	c.line.Color.A = 140

	return c
}

func (c *Collapsible) layoutHeader(gtx *layout.Context, headerFunc func(*layout.Context)) {
	icon := c.collapsedIcon
	if c.isExpanded {
		icon = c.expandedIcon
	}

	layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func() {
			headerFunc(gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Right: unit.Dp(20)}.Layout(gtx, func() {
				icon.Layout(gtx, unit.Dp(20))
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
	c.buttonWidget.Layout(gtx)
}

func (c *Collapsible) Layout(gtx *layout.Context, headerFunc, contentFunc func(*layout.Context)) {
	for c.buttonWidget.Clicked(gtx) {
		c.isExpanded = !c.isExpanded
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func() {
				c.layoutHeader(gtx, headerFunc)
			})
		}),
		layout.Rigid(func() {
			if c.isExpanded {
				contentFunc(gtx)
			}
		}),
		layout.Rigid(func() {
			c.line.Width = gtx.Constraints.Width.Max
			c.line.Layout(gtx)
		}),
	)
}
