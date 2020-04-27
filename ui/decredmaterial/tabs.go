package decredmaterial

import (
	"gioui.org/font"
	"gioui.org/text"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/draw"
)

// todo: add side line for vertical tabs
// todo: add physical scroll button for scrolling horizontally
// todo: add a script for generating Decred icons as bytes
// todo: remove radius from button animation

const (
	Top Position = iota
	Right
	Bottom
	Left
)

var adaptiveTabWidth int

type Position int

type TabItem struct {
	Label
	Icon   image.Image
	iconOp paint.ImageOp
	index  int
}

func tabIndicatorDimensions(gtx *layout.Context, tabPosition Position) (width, height int) {
	switch tabPosition {
	case Top, Bottom:
		width, height = gtx.Dimensions.Size.X, 4
	default:
		width, height = 5, gtx.Dimensions.Size.Y
	}
	return
}

// tabAlignment determines the alignment of the active tab indicator relative to the tab item
// content. It is determined by the position of the tab.
func indicatorDirection(tabPosition Position) layout.Direction {
	switch tabPosition {
	case Top:
		return layout.S
	case Left:
		return layout.W
	case Bottom:
		return layout.N
	case Right:
		return layout.E
	default:
		return layout.W
	}
}

// indicator defines how the active tab indicator is drawn
func activeIndicator(gtx *layout.Context, width, height int) layout.Widget {
	return func() {
		paint.ColorOp{Color: keyblue}.Add(gtx.Ops)
		paint.PaintOp{Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(width),
				Y: float32(height),
			},
		}}.Add(gtx.Ops)
		gtx.Dimensions = layout.Dimensions{
			Size: image.Point{X: width, Y: height},
		}
	}
}

func (t *TabItem) LayoutIcon(gtx *layout.Context) {
	sz := gtx.Constraints.Width.Min
	if t.iconOp.Size().X != sz {
		img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: sz}})
		draw.ApproxBiLinear.Scale(img, img.Bounds(), t.Icon, t.Icon.Bounds(), draw.Src, nil)
		t.iconOp = paint.NewImageOp(img)
	}

	img := material.Image{Src: t.iconOp}
	img.Scale = float32(sz) / float32(gtx.Px(unit.Dp(float32(sz))))
	img.Layout(gtx)
}

func (t *TabItem) iconText(gtx *layout.Context, tabPosition Position) layout.Widget {
	widgetAxis := layout.Vertical
	if tabPosition == Left || tabPosition == Right {
		widgetAxis = layout.Horizontal
	}

	return func() {
		layout.Flex{}.Layout(gtx, layout.Rigid(func() {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
				layout.Flex{Axis: widgetAxis, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func() {
						if t.Icon != nil {
							layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
								dim := gtx.Px(unit.Dp(20))
								sz := image.Point{X: dim, Y: dim}
								gtx.Constraints = layout.RigidConstraints(gtx.Constraints.Constrain(sz))
								t.LayoutIcon(gtx)
							})
						}
					}),
					layout.Rigid(func() {
						if t.Label.shaper != nil {
							t.Label.Alignment = text.Middle
							t.Label.Layout(gtx)
						}
					}),
				)
			})
		}))
	}
}

func (t *TabItem) Layout(gtx *layout.Context, selected int, btn *widget.Button, tabPosition Position) {
	var tabWidth, tabHeight int

	layout.Stack{}.Layout(gtx,
		layout.Stacked(func() {
			t.iconText(gtx, tabPosition)()
			if gtx.Dimensions.Size.X > adaptiveTabWidth {
				adaptiveTabWidth = gtx.Dimensions.Size.X
			}
		}),
		layout.Expanded(func() {
			if tabPosition == Left || tabPosition == Right {
				gtx.Constraints.Width.Min = adaptiveTabWidth
			}
			Button{Background: color.RGBA{}, shaper: font.Default(),}.Layout(gtx, btn)
			tabWidth, tabHeight = tabIndicatorDimensions(gtx, tabPosition)
		}),
		layout.Expanded(func() {
			if selected != t.index {
				return
			}
			indicatorDirection(tabPosition).Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(activeIndicator(gtx, tabWidth, tabHeight)))
			})
		}),
	)
}

// Tabs lays out a Flexed(Size) List with Selected as the first element and Item as the rest.
type Tabs struct {
	Flex     layout.Flex
	Size     float32
	items    []TabItem
	Selected int
	changed  bool
	btns     []*widget.Button
	list     *layout.List
	Position Position
}

func NewTabs() *Tabs {
	return &Tabs{
		list:     &layout.List{},
		Position: Left,
	}
}

// SetTabs creates a button widget for each tab item
func (t *Tabs) SetTabs(tabs []TabItem) {
	t.items = tabs
	if len(t.items) != len(t.btns) {
		t.btns = make([]*widget.Button, len(t.items))
		for i := range t.btns {
			t.btns[i] = new(widget.Button)
		}
	}
}

func (t *Tabs) Changed() bool {
	return t.changed
}

// contentTabPosition depending on the specified tab position determines the order of the tab and
// the page content.
func (t *Tabs) contentTabPosition(gtx *layout.Context, body layout.Widget) (widgets []layout.FlexChild) {
	var content, tab layout.FlexChild

	widgets = make([]layout.FlexChild, 2)
	content = layout.Flexed(1, func() {
		layout.Inset{Left: unit.Dp(5)}.Layout(gtx, body)
	})
	tab = layout.Rigid(func() {
		t.list.Layout(gtx, len(t.btns), func(i int) {
			t.items[i].index = i
			t.items[i].Layout(gtx, t.Selected, t.btns[i], t.Position)
		})
	})

	switch t.Position {
	case Bottom, Right:
		widgets[0], widgets[1] = content, tab
	default:
		widgets[0], widgets[1] = tab, content
	}
	return widgets
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

	widgets := t.contentTabPosition(gtx, body)
	t.Flex.Spacing = layout.SpaceBetween
	t.Flex.Layout(gtx, widgets...)
	for i := range t.btns {
		t.changed = false
		if t.btns[i].Clicked(gtx) {
			t.changed = true
			t.Selected = i
			return
		}
	}
}
