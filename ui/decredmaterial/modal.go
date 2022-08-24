package decredmaterial

import (
	"fmt"
	"image/color"

	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/app"
)

type Modal struct {
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Modal satisfy the app.Modal interface. It also defines
	// helper methods for accessing the WindowNavigator that displayed this
	// modal.
	*app.GenericPageModal

	overlayColor color.NRGBA
	background   color.NRGBA
	list         *widget.List
	button       *widget.Clickable
	card         Card
	scroll       ListStyle
	padding      unit.Dp

	isFloatTitle  bool
	isDisabled    bool
	showScrollBar bool
}

func (t *Theme) ModalFloatTitle(id string) *Modal {
	mod := t.Modal(id)
	mod.isFloatTitle = true
	return mod
}

func (t *Theme) Modal(id string) *Modal {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	uniqueID := fmt.Sprintf("%s-%d", id, GenerateRandomNumber())
	m := &Modal{
		GenericPageModal: app.NewGenericPageModal(uniqueID),
		overlayColor:     overlayColor,
		background:       t.Color.Surface,
		list: &widget.List{
			List: layout.List{Axis: layout.Vertical, Alignment: layout.Middle},
		},
		button:  new(widget.Clickable),
		card:    t.Card(),
		padding: unit.Dp(24),
	}

	m.scroll = t.List(m.list)

	return m
}

// Dismiss removes the modal from the window. Does nothing if the modal was
// not previously pushed into a window.
func (m *Modal) Dismiss() {
	// ParentWindow will only be accessible if this modal has been
	// pushed into display by a WindowNavigator.
	if parentWindow := m.ParentWindow(); parentWindow != nil {
		parentWindow.DismissModal(m.ID())
	} else {
		panic("can't dismiss a modal that has not been displayed")
	}
}

// IsShown is true if this modal has been pushed into a window and is currently
// the top modal in the window.
func (m *Modal) IsShown() bool {
	if parentWindow := m.ParentWindow(); parentWindow != nil {
		topModal := parentWindow.TopModal()
		return topModal != nil && topModal.ID() == m.ID()
	}
	return false
}

// Layout renders the modal widget to screen. The modal assumes the size of
// its content plus padding.
func (m *Modal) Layout(gtx layout.Context, widgets []layout.Widget) layout.Dimensions {
	mGtx := gtx
	if m.isDisabled {
		mGtx = gtx.Disabled()
	}
	dims := layout.Stack{Alignment: layout.Center}.Layout(mGtx,
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			fillMax(gtx, m.overlayColor, CornerRadius{})

			return m.button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				semantic.Button.Add(gtx.Ops)
				return layout.Dimensions{Size: gtx.Constraints.Min}
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

			gtx.Constraints.Max.X = gtx.Dp(unit.Dp(360))
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

							inset := layout.Inset{
								Top:    unit.Dp(10),
								Bottom: unit.Dp(10),
							}
							return inset.Layout(gtx, title)
						}
						return D{}
					}),
					layout.Rigid(func(gtx C) D {
						mTB := unit.Dp(10)
						mLR := unit.Dp(0)
						if m.padding == unit.Dp(0) {
							mLR = mTB
						}
						inset := layout.Inset{
							Top:    mTB,
							Bottom: mTB,
							Left:   mLR,
							Right:  mLR,
						}
						if m.showScrollBar {
							return m.scroll.Layout(gtx, len(widgetFuncs), func(gtx C, i int) D {
								gtx.Constraints.Min.X = gtx.Constraints.Max.X
								return inset.Layout(gtx, widgetFuncs[i])
							})
						}
						list := &layout.List{Axis: layout.Vertical}
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return list.Layout(gtx, len(widgetFuncs), func(gtx C, i int) D {
							return inset.Layout(gtx, widgetFuncs[i])
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

func (m *Modal) SetPadding(padding unit.Dp) {
	m.padding = padding
}

func (m *Modal) ShowScrollbar(showScrollBar bool) {
	m.showScrollBar = showScrollBar
}

func (m *Modal) SetDisabled(disabled bool) {
	m.isDisabled = disabled
}
