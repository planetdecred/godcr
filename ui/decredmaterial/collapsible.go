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

type Collapsible struct {
	isOpen                bool
	buttonWidget          *widget.Button
	outline               Outline
	openIcon              *Icon
	closeIcon             *Icon
	headerBackgroundColor color.RGBA
}

func (t *Theme) Collapsible() *Collapsible {
	c := &Collapsible{
		isOpen:                false,
		headerBackgroundColor: t.Color.Hint,
		outline:               t.Outline(),
		openIcon:              t.chevronUpIcon,
		closeIcon:             t.chevronDownIcon,
	}

	if c.buttonWidget == nil {
		c.buttonWidget = new(widget.Button)
	}

	return c
}

func (c *Collapsible) layoutHeader(gtx *layout.Context, headerFunc func(*layout.Context)) {
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
				layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
					layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func() {
							icon := c.closeIcon
							if c.isOpen {
								icon = c.openIcon
							}
							icon.Layout(gtx, unit.Dp(20))
						}),
						layout.Rigid(func(){
							inset := layout.Inset{
								Left: unit.Dp(20),
							}
							inset.Layout(gtx, func(){
								headerFunc(gtx)
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

func (c *Collapsible) Layout(gtx *layout.Context, headerFunc, contentFunc func(*layout.Context)) {
	for c.buttonWidget.Clicked(gtx) {
		c.isOpen = !c.isOpen
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			c.layoutHeader(gtx, headerFunc)
		}),
		layout.Rigid(func() {
			if c.isOpen {
				contentFunc(gtx)
			}
		}),
	)
}
