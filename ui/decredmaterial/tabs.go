package decredmaterial

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/image/draw"

	"gioui.org/font"
	"gioui.org/text"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr/ui/values"
)

const (
	Top Position = iota
	Right
	Bottom
	Left
)

var adaptiveTabWidth int

// Position determines what side of the page the tab would be laid out
type Position int

// TabItem displays a single child of a tab. Label and Icon in TabItem are optional.
type TabItem struct {
	Label
	Icon   image.Image
	iconOp paint.ImageOp
	index  int
}

// tabIndicatorDimensions defines the width and height of the active tab item indicator depending
// on the tab Position.
func tabIndicatorDimensions(gtx *layout.Context, tabPosition Position) (width, height int) {
	switch tabPosition {
	case Top, Bottom:
		width, height = gtx.Dimensions.Size.X, 4
	default:
		width, height = 5, gtx.Dimensions.Size.Y
	}
	return
}

// indicatorDirection determines the alignment of the active tab indicator relative to the tab item
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

// line returns a rectangle using a defined width, height and color.
func line(gtx *layout.Context, width, height int, col color.RGBA) layout.Widget {
	return func() {
		paint.ColorOp{Color: col}.Add(gtx.Ops)
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

// layoutIcon lays out the icon of a tab item
func (t *TabItem) layoutIcon(gtx *layout.Context) {
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

// iconText lays out the text of a tab item and its icon if it has one. It aligns the text and the icon
// based on the position of the tab.
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
								t.layoutIcon(gtx)
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
			Button{Background: color.RGBA{}, Radius: 0, shaper: font.Default()}.Layout(gtx, btn)
			tabWidth, tabHeight = tabIndicatorDimensions(gtx, tabPosition)
		}),
		layout.Expanded(func() {
			if selected != t.index {
				return
			}
			indicatorDirection(tabPosition).Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(line(gtx, tabWidth, tabHeight, keyblue)))
			})
		}),
	)
}

// Tabs displays succession of TabItems. Using the Position option, Tabs can be displayed on any four sides
// of a rendered page.
type Tabs struct {
	flex             layout.Flex
	Size             float32
	items            []TabItem
	Selected         int
	previousSelected int
	events           []widget.ChangeEvent
	btns             []*widget.Button
	title            Label
	list             *layout.List
	Position         Position
	Separator        bool
	scrollLeft       *widget.Button
	scrollRight      *widget.Button
}

func NewTabs() *Tabs {
	return &Tabs{
		list:        &layout.List{},
		Position:    Left,
		scrollLeft:  new(widget.Button),
		scrollRight: new(widget.Button),
	}
}

// SetTabs creates a button widget for each tab item.
func (t *Tabs) SetTabs(tabs []TabItem) {
	t.items = tabs
	if len(t.items) != len(t.btns) {
		t.btns = make([]*widget.Button, len(t.items))
		for i := range t.btns {
			t.btns[i] = new(widget.Button)
		}
	}
}

// scrollButton lays out the right and left scroll buttons of the tab when Position is Horizontal.
func (t *Tabs) scrollButton(gtx *layout.Context, right bool, button *widget.Button) layout.FlexChild {
	show := false
	icon := mustIcon(NewIcon(icons.NavigationChevronLeft))
	if right && t.list.Position.BeforeEnd {
		show = true
		icon = mustIcon(NewIcon(icons.NavigationChevronRight))
	}

	if !right && t.list.Position.Offset > 0 {
		show = true
	}
	return layout.Rigid(func() {
		if (t.Position == Top || t.Position == Bottom) && show {
			IconButton{
				Color: rgb(0xbbbbbb),
				Icon:  icon,
				Size:  unit.Dp(35),
			}.Layout(gtx, button)
		}
	})
}

// SetTitle setting the title of the tabs
func (t *Tabs) SetTitle(title Label) {
	t.title = title
}

// tabsTitle lays out the title of the tab when Position is Horizontal.
func (t *Tabs) tabsTitle(gtx *layout.Context) layout.FlexChild {
	return layout.Rigid(func() {
		if (t.Position == Top || t.Position == Bottom) && t.title.shaper != nil {
			layout.Inset{Top: values.MarginPadding10, Right: values.MarginPadding5}.Layout(gtx, func() {
				t.title.Layout(gtx)
			})
		}
	})
}

// contentTabPosition determines the order of the tab and page content depending on the tab Position.
func (t *Tabs) contentTabPosition(gtx *layout.Context, body layout.Widget) (widgets []layout.FlexChild) {
	var content, tab layout.FlexChild

	widgets = make([]layout.FlexChild, 2)
	content = layout.Flexed(1, body)
	tab = layout.Rigid(func() {
		layout.Stack{}.Layout(gtx,
			layout.Stacked(func() {
				layout.Flex{Axis: t.list.Axis, Spacing: layout.SpaceBetween}.Layout(gtx,
					t.tabsTitle(gtx),
					t.scrollButton(gtx, false, t.scrollLeft),
					layout.Flexed(1, func() {
						t.list.Layout(gtx, len(t.btns), func(i int) {
							t.items[i].index = i
							t.items[i].Layout(gtx, t.Selected, t.btns[i], t.Position)
						})
					}),
					t.scrollButton(gtx, true, t.scrollRight),
				)
			}),
			layout.Expanded(func() {
				direction := layout.E
				if t.Position == Right {
					direction = layout.W
				}
				// display separator only if Separator is true and tab is vertical
				if t.Separator && (t.Position == Right || t.Position == Left) {
					direction.Layout(gtx, func() {
						layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func() {
								line(gtx, 2, gtx.Constraints.Height.Max, rgb(0xcccccc))()
							}),
						)
					})
				}
			}),
		)
	})

	switch t.Position {
	case Bottom, Right:
		widgets[0], widgets[1] = content, tab
	default:
		widgets[0], widgets[1] = tab, content
	}
	return widgets
}

// ChangeTab changes the position of the selected tab
func (t *Tabs) ChangeTab(index int) {
	t.Selected = index
	if t.previousSelected != t.Selected {
		t.events = append(t.events, widget.ChangeEvent{})
		t.previousSelected = index
	}
}

// ChangeEvent returns the last change event
func (t *Tabs) ChangeEvent(gtx *layout.Context) []widget.ChangeEvent {
	t.processChangeEvent(gtx)
	events := t.events
	t.events = nil
	return events
}

func (t *Tabs) processChangeEvent(gtx *layout.Context) {
	for i := range t.btns {
		if t.btns[i].Clicked(gtx) {
			t.ChangeTab(i)
			return
		}
	}
}

func (t *Tabs) Layout(gtx *layout.Context, body layout.Widget) {
	switch t.Position {
	case Top, Bottom:
		t.list.Axis = layout.Horizontal
		t.flex.Axis = layout.Vertical
	default:
		t.list.Axis = layout.Vertical
		t.flex.Axis = layout.Horizontal
	}

	widgets := t.contentTabPosition(gtx, body)
	t.flex.Spacing = layout.SpaceBetween
	t.flex.Layout(gtx, widgets...)

	for t.scrollRight.Clicked(gtx) {
		t.list.Position.Offset += 60
	}

	for t.scrollLeft.Clicked(gtx) {
		t.list.Position.Offset -= 60
	}

	t.processChangeEvent(gtx)
}
