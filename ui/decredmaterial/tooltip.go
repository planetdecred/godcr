package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/planetdecred/godcr/ui/values"
)

type Tooltip struct {
	hoverable  *Hoverable
	background color.NRGBA
	shadow     *Shadow
}

func (t *Theme) Tooltip() *Tooltip {
	return &Tooltip{
		hoverable:  t.Hoverable(),
		background: t.Color.Surface,
		shadow:     t.Shadow(),
	}
}

func (t *Tooltip) layout(gtx C, pos layout.Inset, wdgt layout.Widget) D {
	border := Border{
		Radius: Radius(7),
	}

	return pos.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return LinearLayout{
					Width:      WrapContent,
					Height:     WrapContent,
					Padding:    layout.UniformInset(values.MarginPadding12),
					Background: t.background,
					Border:     border,
					Shadow:     t.shadow,
				}.Layout2(gtx, wdgt)
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
