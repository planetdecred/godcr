package widgets

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr-gio/helper"
)

type (
	Checkbox struct {
		isChecked   bool
		icon        *Icon
		padding     unit.Value
		size        unit.Value
		button      *widget.Button
		isClickable bool
	}
)

func NewCheckbox() *Checkbox {
	return &Checkbox{
		isChecked:   false,
		icon:        NavigationCheckIcon,
		button:      new(widget.Button),
		padding:     unit.Dp(5),
		size:        unit.Dp(26),
		isClickable: true,
	}
}

func (c *Checkbox) IsChecked() bool {
	return c.isChecked
}

func (c *Checkbox) MakeAsIcon() *Checkbox {
	c.isClickable = false
	return c
}

func (c *Checkbox) toggleCheckState() {
	if c.isChecked {
		c.isChecked = false
		return
	}
	c.isChecked = true
}

func (c *Checkbox) SetSize(size float32) *Checkbox {
	c.size = unit.Dp(size)
	return c
}

func (c *Checkbox) SetPadding(padding float32) *Checkbox {
	c.padding = unit.Dp(padding)
	return c
}

func (c *Checkbox) Draw(ctx *layout.Context) {
	for c.button.Clicked(ctx) {
		c.toggleCheckState()
	}

	bgcol := helper.DecredGreenColor
	if !c.isClickable {
		c.isChecked = true
	}

	if !c.isChecked {
		bgcol = helper.WhiteColor
	}

	hmin := ctx.Constraints.Width.Min
	vmin := ctx.Constraints.Height.Min
	layout.Stack{Alignment: layout.Center}.Layout(ctx,
		layout.Expanded(func() {
			ctx.Constraints.Height.Min = ctx.Constraints.Width.Min
			size := float32(ctx.Constraints.Width.Min)
			rr := float32(size) * .5
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(ctx.Constraints.Width.Min),
					Y: float32(ctx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(ctx.Ops).Add(ctx.Ops)
			Fill(ctx, helper.DecredGreenColor)

			layout.Align(layout.Center).Layout(ctx, func() {
				layout.UniformInset(unit.Dp(1)).Layout(ctx, func() {
					ctx.Constraints.Width.Min = 34
					ctx.Constraints.Height.Min = ctx.Constraints.Width.Min

					mainSize := float32(ctx.Constraints.Width.Min)
					mainRadius := float32(mainSize) * .5
					clip.Rect{
						Rect: f32.Rectangle{Max: f32.Point{
							X: mainSize,
							Y: mainSize,
						}},
						NE: mainRadius, NW: mainRadius, SE: mainRadius, SW: mainRadius,
					}.Op(ctx.Ops).Add(ctx.Ops)
					Fill(ctx, bgcol)
					for _, c := range c.button.History() {
						drawInk(ctx, c)
					}
				})
			})
		}),
		layout.Stacked(func() {
			ctx.Constraints.Width.Min = hmin
			ctx.Constraints.Height.Min = vmin
			layout.Align(layout.Center).Layout(ctx, func() {
				layout.UniformInset(c.padding).Layout(ctx, func() {
					size := ctx.Px(c.size) - 2*ctx.Px(c.padding)
					if c.isChecked {
						ico := c.icon.image(size)
						ico.Add(ctx.Ops)
						paint.PaintOp{
							Rect: f32.Rectangle{Max: toPointF(ico.Size())},
						}.Add(ctx.Ops)
					}
					ctx.Dimensions = layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				})
			})
			if c.isClickable {
				pointer.Ellipse(image.Rectangle{Max: ctx.Dimensions.Size}).Add(ctx.Ops)
				c.button.Layout(ctx)
			}
		}),
	)
}
