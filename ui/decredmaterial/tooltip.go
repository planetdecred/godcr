package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Tooltip struct {
	hoverable   *Hoverable
	card        Card
	borderColor color.NRGBA
}

func (t *Theme) Tooltip() *Tooltip {
	return &Tooltip{
		hoverable:   t.Hoverable(),
		card:        t.Card(),
		borderColor: t.Color.Gray1,
	}
}

func (t *Tooltip) layout(gtx C, pos layout.Inset, wdgt layout.Widget) D {
	border := widget.Border{
		Color:        t.borderColor,
		CornerRadius: unit.Dp(5),
		Width:        unit.Dp(1),
	}

	return pos.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return border.Layout(gtx, func(gtx C) D {
					return t.card.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, wdgt)
					})
				})
			}),
		)
	})
}

func (t *Tooltip) Layout(gtx C, rect image.Rectangle, pos layout.Inset, wdgt layout.Widget) D {
	if t.hoverable.Hovered() {
		m := op.Record(gtx.Ops)
		t.layout(gtx, pos, wdgt)
		op.Defer(gtx.Ops, m.Stop())
	}

	t.hoverable.Layout(gtx, rect)
	return D{
		Size: rect.Min,
	}
}
