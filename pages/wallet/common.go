package wallet

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr-gio/helper"
	"github.com/raedahgroup/godcr-gio/widgets"
)

const (
	bodyHeight = 300
)

func drawHeader(ctx *layout.Context, backFunc func(), titleFunc func()) {
	inset := layout.Inset{
		Top:   unit.Dp(0),
		Left:  unit.Dp(helper.StandaloneScreenPadding),
		Right: unit.Dp(helper.StandaloneScreenPadding),
	}
	inset.Layout(ctx, func() {
		layout.Stack{}.Layout(ctx,
			layout.Expanded(func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
					layout.Rigid(func() {
						ctx.Constraints.Width.Min = 50
						backFunc()
					}),
					layout.Rigid(func() {
						inset := layout.Inset{
							Top: unit.Dp(10),
						}
						inset.Layout(ctx, func() {
							titleFunc()
						})
					}),
				)
			}),
		)
	})
}

func drawBody(ctx *layout.Context, title *widgets.Label, bodyFunc func()) {
	topInset := float32(68)
	if title != nil {
		inset := layout.Inset{
			Top:   unit.Dp(topInset),
			Left:  unit.Dp(helper.StandaloneScreenPadding),
			Right: unit.Dp(helper.StandaloneScreenPadding),
		}
		inset.Layout(ctx, func() {
			title.Draw(ctx)
		})
		topInset += 20
	}

	inset := layout.Inset{
		Top:   unit.Dp(topInset),
		Left:  unit.Dp(helper.StandaloneScreenPadding),
		Right: unit.Dp(helper.StandaloneScreenPadding),
	}
	inset.Layout(ctx, func() {
		bodyFunc()
	})
}

func drawCardBody(ctx *layout.Context, title *widgets.Label, bodyFunc func()) {
	drawBody(ctx, title, func() {
		inset := layout.Inset{
			Top: unit.Dp(0),
		}
		inset.Layout(ctx, func() {
			ctx.Constraints.Height.Min = 400
			helper.PaintArea(ctx, helper.WhiteColor, ctx.Constraints.Width.Max, bodyHeight+20)
			bodyFunc()
		})
	})

}

func drawFooter(ctx *layout.Context, footerFunc func()) {
	inset := layout.Inset{
		Top: unit.Dp(float32(bodyHeight + 120)),
	}
	inset.Layout(ctx, func() {
		helper.PaintFooter(ctx, helper.WhiteColor, ctx.Constraints.Width.Max, 200)
		inset := layout.UniformInset(unit.Dp(20))
		inset.Layout(ctx, func() {
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			footerFunc()
		})
	})
}
