package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Modal struct {
	overlayColor color.NRGBA
	background   color.NRGBA
	list         *widget.List
	button       *widget.Clickable
	card         Card
	scroll       ListStyle
	isFloatTitle bool
	padding      unit.Value
}

func (t *Theme) ModalFloatTitle() *Modal {
	mod := t.Modal()
	mod.isFloatTitle = true
	return mod
}

func (t *Theme) Modal() *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	m := &Modal{
		overlayColor: overlayColor,
		background:   t.Color.Surface,
		list: &widget.List{
			List: layout.List{Axis: layout.Vertical, Alignment: layout.Middle},
		},
		button:  new(widget.Clickable),
		card:    t.Card(),
		padding: unit.Dp(16),
	}

	m.scroll = t.List(m.list)

	return m
}

// Layout renders the modal widget to screen. The modal assumes the size of
// its content plus padding.
func (m *Modal) Layout(gtx layout.Context, widgets []layout.Widget) layout.Dimensions {
	dims := layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return m.button.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return fillMax(gtx, m.overlayColor, CornerRadius{})
			})
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

			gtx.Constraints.Max.X = gtx.Px(unit.Dp(380))
			inset := layout.Inset{
				Top:    unit.Dp(50),
				Bottom: unit.Dp(50),
			}
			return inset.Layout(gtx, func(gtx C) D {
				return LinearLayout{
					Orientation: layout.Vertical,
					Width:       WrapContent,
					Height:      WrapContent,
					Padding:     layout.UniformInset(m.padding),
					Alignment:   layout.Middle,
					Border: Border{
						Radius: Radius(14),
					},
					Direction:  layout.Center,
					Background: m.background,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if m.isFloatTitle && len(widgets) > 0 {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							if m.padding == unit.Dp(0) {
								return layout.UniformInset(m.padding).Layout(gtx, title)
							}
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, title)
						}
						return D{}
					}),
					layout.Rigid(func(gtx C) D {
						return m.scroll.Layout(gtx, len(widgetFuncs), func(gtx C, i int) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
						})
					}),
				)
			})
		}),
	)

	return dims
}

func (m *Modal) BackdropClicked(minimizable bool) bool {
	if minimizable {
		return m.button.Clicked()
	}

	return false
}

func (m *Modal) SetPadding(padding unit.Value) {
	m.padding = padding
}
