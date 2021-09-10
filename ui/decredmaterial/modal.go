package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Modal struct {
	overlayColor color.NRGBA
	co           color.NRGBA
	list         *layout.List
	button       *widget.Clickable
	card         Card
	isFloatTitle bool
}

func (t *Theme) ModalFloatTitle() *Modal {
	mod := t.Modal()
	mod.isFloatTitle = true
	return mod
}

func (t *Theme) Modal() *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200
	co := t.Color.Surface

	return &Modal{
		overlayColor: overlayColor,
		co:           co,
		list:         &layout.List{Axis: layout.Vertical, Alignment: layout.Middle},
		button:       new(widget.Clickable),
		card:         t.Card(),
	}
}

// Layout renders the modal widget to screen. The modal assumes the size of
// its content plus padding.
func (m *Modal) Layout(gtx layout.Context, widgets []layout.Widget, margin int) layout.Dimensions {
	dims := layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			fillMax(gtx, m.overlayColor, CornerRadius{})
			return m.button.Layout(gtx)
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			var widgetFuncs []layout.Widget
			var title layout.Widget

			if m.isFloatTitle && len(widgets) > 0 {
				title = widgets[0]
				widgetFuncs = append(widgetFuncs, widgets[1:]...)
			} else {
				widgetFuncs = append(widgetFuncs, widgets...)
			}

			gtx.Constraints.Max.X = gtx.Px(unit.Dp(450))
			return LinearLayout{
				Orientation: layout.Vertical,
				Width:       WrapContent,
				Height:      WrapContent,
				Padding:     layout.UniformInset(unit.Dp(15)),
				Alignment:   layout.Middle,
				Border: Border{
					Radius: Radius(10),
				},
				Direction:  layout.Center,
				Background: m.co,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if m.isFloatTitle && len(widgets) > 0 {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, title)
					}
					return D{}
				}),
				layout.Rigid(func(gtx C) D {
					return m.list.Layout(gtx, len(widgetFuncs), func(gtx C, i int) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
					})
				}),
			)
		}),
	)

	return dims
}

func (m *Modal) BackdropClicked() bool {
	return m.button.Clicked()
}
