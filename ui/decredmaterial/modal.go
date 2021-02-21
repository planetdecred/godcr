package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Modal struct {
	overlayColor color.NRGBA
	list         *layout.List
	button       *widget.Clickable
	card         Card
}

func (t *Theme) Modal() *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	return &Modal{
		overlayColor: overlayColor,
		list:         &layout.List{Axis: layout.Vertical, Alignment: layout.Middle},
		button:       new(widget.Clickable),
		card:         t.Card(),
	}
}

// Layout renders the modal widget to screen. The modal assumes the size of
// its content plus padding.
func (m *Modal) Layout(gtx layout.Context, widgets []func(gtx C) D, margin int) layout.Dimensions {
	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			fillMax(gtx, m.overlayColor)
			return m.button.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			var widgetFuncs []func(gtx C) D

			widgetFuncs = append(widgetFuncs, widgets...)
			scaled := 3840 / float32(gtx.Constraints.Max.X)
			mg := unit.Px(float32(margin) / scaled)

			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Left:  mg,
					Right: mg,
				}.Layout(gtx, func(gtx C) D {
					return m.card.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
							return m.list.Layout(gtx, len(widgetFuncs), func(gtx C, i int) D {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
							})
						})
					})
				})
			})
		}),
	)

	return dims
}
