package decredmaterial

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image"
	"image/color"
)

// DefaultTabSize is the default flexed size of the tab section in a Tabs
const DefaultTabSize = .15

const (
	Top Position = iota
	Right
	Bottom
	Left
)

type Position int

type TabItem struct {
	Button
	Icon
}

func (t *TabItem) Layout(gtx *layout.Context, tabIndex, selected int, btn *widget.Button) {
	var tabWidth, tabHeight int

	layout.Stack{Alignment: layout.E}.Layout(gtx,
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			t.Button.Color = darkblue
			t.Button.Background = color.RGBA{}
			t.Button.Layout(gtx, btn)
			tabHeight = gtx.Dimensions.Size.Y
		}),
		layout.Stacked(func() {
			layout.Flex{Axis:layout.Horizontal}.Layout(gtx, layout.Flexed(0, func() {
				if selected != tabIndex {
					return
				}
				paint.ColorOp{Color: keyblue}.Add(gtx.Ops)
				paint.PaintOp{Rect: f32.Rectangle{
					Max: f32.Point{
						X: float32(5),
						Y: float32(tabHeight),
					},
				}}.Add(gtx.Ops)
				gtx.Dimensions = layout.Dimensions{
					Size: image.Point{X: tabWidth, Y: tabHeight},
				}
			}))
		}),
	)
}

// Tabs laysout a Flexed(Size) List with Selected as the first element and Item as the rest.
type Tabs struct {
	Flex        layout.Flex
	Size        float32
	items       []TabItem
	Selected    int
	changed     bool
	btns        []*widget.Button
	list        *layout.List
	Position 	Position
}

func NewTabs() *Tabs {
	return &Tabs{
		list: &layout.List{},
		Position: Left,
		Size: DefaultTabSize,
	}
}

func (t *Tabs) SetTabs(tabs []TabItem) {
	t.items = tabs
	if len(t.items) != len(t.btns) {
		t.btns = make([]*widget.Button, len(t.items))
		for i := range t.btns {
			t.btns[i] = new(widget.Button)
		}
	}
}

// Layout the tabs
func (t *Tabs) Layout(gtx *layout.Context, body layout.Widget) {
	switch t.Position {
	case Top, Bottom:
		t.list.Axis = layout.Horizontal
		t.Flex.Axis = layout.Vertical
	default:
		t.list.Axis = layout.Vertical
		t.Flex.Axis = layout.Horizontal
	}

	t.Flex.Layout(gtx,
		layout.Flexed(t.Size, func() {
			t.list.Layout(gtx, len(t.btns), func(i int) {
				t.items[i].Layout(gtx, i, t.Selected, t.btns[i])
				if t.btns[i].Clicked(gtx) {
					t.Selected = i
				}
			})
		}),
		layout.Flexed(1-t.Size, body),
	)
}
