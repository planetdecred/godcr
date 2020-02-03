package widgets

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr-gio/ui/helpers"
)

type selectItem struct {
	key string
	val string

	clicker helpers.Clicker
}

// Select represents a combo widget
type Select struct {
	items         []selectItem
	selectedIndex int
	isOpen        bool
	textSize      float32
	shaper        text.Shaper
}

// Select returns an instance of the select widget
func (t *Theme) Select(items map[string]string) *Select {
	s := &Select{
		isOpen:   false,
		items:    make([]selectItem, len(items)+1),
		textSize: t.TextSize.V,
		shaper:   t.Shaper,
	}

	if len(items) > 0 {
		counter := 0
		for key, val := range items {
			// the item at the zeroeth index is the trigger
			if counter == 0 {
				s.items[0] = selectItem{
					key:     key,
					val:     val,
					clicker: helpers.NewClicker(),
				}
			}
			s.items[counter+1] = selectItem{
				key:     key,
				val:     val,
				clicker: helpers.NewClicker(),
			}
			counter++
		}
	}

	return s
}

// Draw renders the select instance on screen
func (s *Select) Draw(gtx *layout.Context) {
	gtx.Constraints.Width.Min = 100

	container := layout.List{Axis: layout.Vertical}
	container.Layout(gtx, len(s.items), func(i int) {
		if s.isOpen || i == 0 {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				for s.items[i].clicker.Clicked(gtx) {
					if i != 0 {
						s.setSelected(i)
					}
					s.isOpen = !s.isOpen
				}

				s.drawItem(gtx, &s.items[i])
			})
		}
	})
}

func (s *Select) setSelected(itemIndex int) {
	s.selectedIndex = itemIndex
	s.items[0].key = s.items[itemIndex].key
}

func (s *Select) drawItem(gtx *layout.Context, item *selectItem) {
	col := helpers.BlackColor
	bgcol := helpers.LightGrayColor
	vmin := gtx.Constraints.Height.Min

	font := text.Font{
		Size: unit.Dp(s.textSize),
	}

	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rounding := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{
					Max: f32.Point{
						X: float32(gtx.Constraints.Width.Min),
						Y: float32(gtx.Constraints.Height.Min),
					},
				},
				NE: rounding,
				NW: rounding,
				SE: rounding,
				SW: rounding,
			}.Op(gtx.Ops).Add(gtx.Ops)
			helpers.Fill(gtx, bgcol)
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = 120
			gtx.Constraints.Height.Min = vmin

			layout.Align(layout.Start).Layout(gtx, func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					paint.ColorOp{Color: col}.Add(gtx.Ops)
					widget.Label{}.Layout(gtx, s.shaper, font, item.key)
				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			item.clicker.Register(gtx)
		}),
	)
}
