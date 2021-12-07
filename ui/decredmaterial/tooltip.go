package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/planetdecred/godcr/ui/values"
)

type Tooltip struct {
	theme      *Theme
	hoverable  *Hoverable
	background color.NRGBA
	shadow     *Shadow
}

func (t *Theme) Tooltip() *Tooltip {
	return &Tooltip{
		theme:     t,
		hoverable: t.Hoverable(),
		shadow:    t.Shadow(),
	}
}

func (t *Tooltip) layout(gtx C, pos layout.Inset, wdgt layout.Widget) D {

	border := Border{
		Width:  values.MarginPadding1,
		Radius: Radius(8),
	}

	return pos.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return LinearLayout{
					Width:      WrapContent,
					Height:     WrapContent,
					Padding:    layout.UniformInset(values.MarginPadding12),
					Background: t.theme.Color.Surface,
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
