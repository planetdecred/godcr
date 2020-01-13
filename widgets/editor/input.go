package editor

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr-gio/helper"
	"github.com/raedahgroup/godcr-gio/widgets"
)

type (
	Input struct {
		*Editor
		hint              string
		focusBorderColor  color.RGBA
		normalBorderColor color.RGBA
	}
)

func NewInput(hint string) *Input {
	return &Input{
		Editor:            new(Editor),
		hint:              hint,
		focusBorderColor:  helper.DecredLightBlueColor,
		normalBorderColor: helper.GrayColor,
	}
}

func (i *Input) SetMask(char string) *Input {
	i.setMask(char)
	return i
}

func (i *Input) Numeric() *Input {
	i.numeric()
	return i
}

func (i *Input) OnType(onType func()) {

}

func (i *Input) SetFocusedBorderColor(col color.RGBA) *Input {
	i.focusBorderColor = col
	return i
}

func (i *Input) SetBorderColor(col color.RGBA) *Input {
	i.normalBorderColor = col
	return i
}

func (i *Input) Draw(ctx *layout.Context) {
	borderColor := i.normalBorderColor
	if i.focused {
		borderColor = i.focusBorderColor
	}

	layout.Stack{}.Layout(ctx,
		layout.Expanded(func() {
			borderRadius := float32(ctx.Px(unit.Dp(4)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(ctx.Constraints.Width.Min),
					Y: float32(ctx.Constraints.Height.Min),
				}},
				NE: borderRadius, NW: borderRadius, SE: borderRadius, SW: borderRadius,
			}.Op(ctx.Ops).Add(ctx.Ops)
			widgets.Fill(ctx, borderColor)

			layout.Align(layout.Center).Layout(ctx, func() {
				layout.UniformInset(unit.Dp(1)).Layout(ctx, func() {
					ctx.Constraints.Height.Min = 48
					ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
					clip.Rect{
						Rect: f32.Rectangle{Max: f32.Point{
							X: float32(ctx.Constraints.Width.Min),
							Y: float32(ctx.Constraints.Height.Min),
						}},
						NE: borderRadius, NW: borderRadius, SE: borderRadius, SW: borderRadius,
					}.Op(ctx.Ops).Add(ctx.Ops)
					widgets.Fill(ctx, helper.WhiteColor)
				})
			})
		}),
		layout.Stacked(func() {
			ctx.Constraints.Height.Min = 50
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			layout.Align(layout.Center).Layout(ctx, func() {
				layout.UniformInset(unit.Dp(8)).Layout(ctx, func() {
					ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
					i.draw(ctx)
				})
			})
		}),
	)

	/**
	stack := layout.Stack{}
	input := stack.Rigid(ctx, func(){
		ctx.Constraints.Height.Min = 50
		ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
		layout.Align(layout.Center).Layout(ctx, func(){
			layout.UniformInset(unit.Dp(8)).Layout(ctx, func() {
				ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
				i.draw(ctx)
			})
		})
	})

	bg := stack.Expand(ctx, func(){
		rr := float32(ctx.Px(unit.Dp(4)))
		ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
		widgets.Rrect(ctx.Ops,
			float32(ctx.Constraints.Width.Min),
			float32(ctx.Constraints.Height.Min),
			rr, rr, rr, rr,
		)
		widgets.Fill(ctx, borderColor)

		layout.Align(layout.Center).Layout(ctx, func(){
			layout.UniformInset(unit.Dp(1)).Layout(ctx, func(){
				ctx.Constraints.Height.Min = 48
				ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
				widgets.Rrect(ctx.Ops,
					float32(ctx.Constraints.Width.Min),
					float32(ctx.Constraints.Height.Min),
					rr, rr, rr, rr,
				)
				widgets.Fill(ctx, helper.WhiteColor)
			})
		})
	})
	stack.Layout(ctx, bg, input)**/
}

func (i *Input) draw(ctx *layout.Context) {
	theme := helper.GetTheme()

	var stack op.StackOp
	stack.Push(ctx.Ops)
	var macro op.MacroOp
	macro.Record(ctx.Ops)
	paint.ColorOp{
		Color: helper.GrayColor,
	}.Add(ctx.Ops)
	widgets.NewLabel(i.hint, 3).SetColor(helper.GrayColor).Draw(ctx)
	macro.Stop()
	if w := ctx.Dimensions.Size.X; ctx.Constraints.Width.Min < w {
		ctx.Constraints.Width.Min = w
	}
	if h := ctx.Dimensions.Size.Y; ctx.Constraints.Height.Min < h {
		ctx.Constraints.Height.Min = h
	}
	i.Layout(ctx, theme.Shaper, theme.Fonts.Regular)

	if i.Len() > 0 {
		paint.ColorOp{
			Color: helper.BlackColor,
		}.Add(ctx.Ops)
		i.PaintText(ctx)
	} else {
		macro.Add()
	}
	paint.ColorOp{
		Color: helper.BlackColor,
	}.Add(ctx.Ops)
	i.PaintCaret(ctx)
	stack.Pop()
}
