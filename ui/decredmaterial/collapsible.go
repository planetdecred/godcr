package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Direction int

type Collapsible struct {
	isOpen                bool
	axis                  Direction
	buttonWidget          *widget.Button
	textLabel             Label
	outline               Outline
	openIcon              *Icon
	closeIcon             *Icon
	headerBackgroundColor color.RGBA
}

const (
	LayoutVertical Direction = iota
	LayoutHorizontal
)

func (t *Theme) Collapsible(text string, axis Direction) *Collapsible {
	c := &Collapsible{
		axis:                  axis,
		textLabel:             t.Body2(text),
		buttonWidget:          new(widget.Button),
		headerBackgroundColor: t.Color.Hint,
		outline:               t.Outline(),
		openIcon:              t.chevronUpIcon,
		closeIcon:             t.chevronDownIcon,
	}

	return c
}

func (c *Collapsible) layoutHeader(gtx *layout.Context) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(4)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, c.headerBackgroundColor)
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			c.outline.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
					layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func() {
							c.textLabel.Layout(gtx)
						}),
						layout.Rigid(func() {
							gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
							layout.NE.Layout(gtx, func() {
								icon := c.closeIcon
								if c.isOpen {
									icon = c.openIcon
								}
								inset := layout.Inset{
									Right: unit.Dp(10),
								}
								inset.Layout(gtx, func() {
									icon.Layout(gtx, unit.Dp(20))
								})
							})
						}),
					)
				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			c.buttonWidget.Layout(gtx)
		}),
	)
}

func (c *Collapsible) Layout(gtx *layout.Context, contentFunc func(*layout.Context)) {
	for c.buttonWidget.Clicked(gtx) {
		c.isOpen = !c.isOpen
	}

	var axis layout.Axis
	if c.axis == LayoutHorizontal {
		axis = layout.Horizontal
	} else {
		axis = layout.Vertical
	}

	layout.Flex{Axis: axis}.Layout(gtx,
		layout.Rigid(func() {
			c.layoutHeader(gtx)
		}),
		layout.Rigid(func() {
			if c.isOpen {
				contentFunc(gtx)
			}
		}),
	)
}
