package decredmaterial

import (
	"fmt"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// DefaultTabSizeVertical is the default flexed size of the tab section in a Tabs when vertically aligned
const DefaultTabSizeVertical = .15

// DefaultTabSizeHorizontal is the default flexed size of the tab section in a Tabs when horizontally aligned
const DefaultTabSizeHorizontal = .10

const (
	Top Position = iota
	Right
	Bottom
	Left
)

type Position int

type TabItem struct {
	Button
	Label
	index     int
	icon      paint.ImageOp
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
		return layout.E
	case Bottom:
		return layout.N
	case Right:
		return layout.W
	default:
		return layout.E
	}
}

// indicator defines how the active tab indicator is drawn
func indicator(gtx *layout.Context, width, height int) layout.Widget {
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

func (t *TabItem) LayoutIcon(gtx *layout.Context, icon image.Image) {
	sz := gtx.Constraints.Width.Min
	if t.icon.Size().X != sz {
		img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: sz}})
		draw.ApproxBiLinear.Scale(img, img.Bounds(), icon, icon.Bounds(), draw.Src, nil)
		t.icon = paint.NewImageOp(img)
	}

	img := material.Image{Src: t.icon}
	img.Scale = float32(sz) / float32(gtx.Px(unit.Dp(float32(sz))))
	img.Layout(gtx)
}

func readImage() (image.Image, error) {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	file, err := os.Open(filepath.Dir(d) + "/assets/decredicons/overview.png")
	if err != nil {
		return nil, fmt.Errorf("img.jpg file not found!")
	}

	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func (t *TabItem) iconText(gtx *layout.Context, tabPosition Position) layout.Widget {
	widgetAxis := layout.Vertical
	if tabPosition == Left || tabPosition == Right {
		gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
		widgetAxis = layout.Horizontal
	}

	return func() {
		layout.Flex{}.Layout(gtx, layout.Rigid(func() {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
				layout.Flex{Axis: widgetAxis, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func() {
						if t.icon.Size().X != 0 {
							layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
								dim := gtx.Px(unit.Dp(20))
								sz := image.Point{X: dim, Y: dim}
								gtx.Constraints = layout.RigidConstraints(gtx.Constraints.Constrain(sz))
								img, _ := readImage()
								t.LayoutIcon(gtx, img)
							})
						}
					}),
					layout.Rigid(func() {
						t.Label.Alignment = text.Middle
						t.Label.Layout(gtx)
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
		}),
		layout.Expanded(func() {
			if tabPosition == Left || tabPosition == Right {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			}
			t.Button.Color = darkblue
			t.Button.Background = color.RGBA{}
			t.Button.Layout(gtx, btn)
			tabWidth, tabHeight = tabIndicatorDimensions(gtx, tabPosition)
		}),
		layout.Expanded(func() {
			if selected != t.index {
				return
			}
			indicatorDirection(tabPosition).Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(indicator(gtx, tabWidth, tabHeight)))
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
		Size:     DefaultTabSizeVertical,
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

// contentTabPosition depending on the specified tab position determines the order of the tab and
// the page content.
func (t *Tabs) contentTabPosition(gtx *layout.Context, body layout.Widget) (widgets []layout.FlexChild) {
	var content, tab layout.FlexChild

	widgets = make([]layout.FlexChild, 2)
	content = layout.Flexed(1-t.Size, func() {
		layout.Inset{Left: unit.Dp(5)}.Layout(gtx, body)
	})
	tab = layout.Flexed(t.Size, func() {
		t.list.Layout(gtx, len(t.btns), func(i int) {
			t.items[i].index = i
			t.items[i].Layout(gtx, t.Selected, t.btns[i], t.Position)
			if t.btns[i].Clicked(gtx) {
				t.Selected = i
			}
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
		if t.Size < DefaultTabSizeHorizontal {
			t.Size = DefaultTabSizeHorizontal
		}
		t.list.Axis = layout.Horizontal
		t.Flex.Axis = layout.Vertical
	default:
		t.list.Axis = layout.Vertical
		t.Flex.Axis = layout.Horizontal
	}

	widgets := t.contentTabPosition(gtx, body)
	t.Flex.Layout(gtx, widgets...)
}
