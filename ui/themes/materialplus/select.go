package materialplus

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

	"github.com/raedahgroup/godcr-gio/ui"
)

// SelectItem represents a select option
// the key is the option value. the text is the text to be displayed
type SelectItem struct {
	Key  int
	Text string

	button *widget.Button
}

// Select represents a combo widget
type Select struct {
	items         []SelectItem
	selectedIndex int
	isOpen        bool
	textSize      float32
	shaper        text.Shaper
}

// Select returns an instance of the select widget
func (t *Theme) Select(items []SelectItem) *Select {
	s := &Select{
		isOpen:   false,
		textSize: t.TextSize.V,
		shaper:   t.Shaper,
	}

	s.SetOptions(items)

	return s
}

// SetOptions sets the select options
func (s *Select) SetOptions(items []SelectItem) {
	s.items = make([]SelectItem, len(items)+1)

	if len(items) > 0 {
		// init option buttons
		for i := range items {
			items[i].button = new(widget.Button)
			s.items[i+1] = items[i]
		}

		// set the current selected option to be the first option
		s.items[0] = items[0]
		s.items[0].button = new(widget.Button)
	}
}

// UpdateOptions updates already set options
func (s *Select) UpdateOptions(items []SelectItem) {
	if len(items) > 0 {
		s.items[0].Key = items[0].Key
		s.items[0].Text = items[0].Text

		for i := range items {
			s.items[i+1].Key = items[i].Key
			s.items[i+1].Text = items[i].Text
		}
	}
}

// Layout renders the select instance on screen
func (s *Select) Layout(gtx *layout.Context) {
	gtx.Constraints.Width.Min = 100

	container := layout.List{Axis: layout.Vertical}
	container.Layout(gtx, len(s.items), func(i int) {
		if s.isOpen || i == 0 && s.items[i].button != nil {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				for s.items[i].button.Clicked(gtx) {
					if i != 0 {
						s.setSelected(i)
					}

					s.isOpen = !s.isOpen
				}

				s.layoutItem(gtx, &s.items[i])
			})
		}
	})
}

func (s *Select) setSelected(itemIndex int) {
	s.selectedIndex = itemIndex
	s.items[0].Key = s.items[itemIndex].Key
	s.items[0].Text = s.items[itemIndex].Text
}

func (s *Select) layoutItem(gtx *layout.Context, item *SelectItem) {
	col := ui.BlackColor
	bgcol := ui.LightGrayColor
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
			Fill(gtx, bgcol, gtx.Constraints.Width.Min, gtx.Constraints.Height.Min)
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = 120
			gtx.Constraints.Height.Min = vmin

			layout.Align(layout.Start).Layout(gtx, func() {
				layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
					paint.ColorOp{Color: col}.Add(gtx.Ops)
					widget.Label{}.Layout(gtx, s.shaper, font, item.Text)
				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			item.button.Layout(gtx)
		}),
	)
}

// GetSelected returns the currently selected item
func (s *Select) GetSelected() *SelectItem {
	return &s.items[s.selectedIndex]
}
