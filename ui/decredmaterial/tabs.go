package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const (
	Top Position = iota
	Right
	Bottom
	Left
)

const (
	navDrawerWidth          = 200
	navDrawerMinimizedWidth = 120
)

var adaptiveTabWidth int

// Position determines what side of the page the tab would be laid out
type Position int

// TabItem displays a single child of a tab. Label and Icon in TabItem are optional.
type TabItem struct {
	ID             string
	Title          string
	label          Label
	button         Button
	Icon           image.Image
	iconOp         paint.ImageOp
	inactiveIconOp paint.ImageOp
	index          int
}

// tabIndicatorDimensions defines the width and height of the active tab item indicator depending
// on the tab Position.
func tabIndicatorDimensions(dims layout.Dimensions, tabPosition Position) (width, height int) {
	switch tabPosition {
	case Top, Bottom:
		width, height = dims.Size.X, 4
	default:
		width, height = 0, dims.Size.Y
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
func line(gtx layout.Context, width, height int, col color.RGBA) layout.Dimensions {
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{
		Max: f32.Point{
			X: float32(width),
			Y: float32(height),
		},
	}}.Add(gtx.Ops)
	return layout.Dimensions{
		Size: image.Point{X: width, Y: height},
	}
}

// layoutIcon lays out the icon of a tab item
func (t *TabItem) layoutIcon(gtx layout.Context) layout.Dimensions {
	//iconOp := t.iconOp
	//if t.sele

	img := widget.Image{Src: t.iconOp}
	img.Scale = 0.05
	return img.Layout(gtx)
}

// iconText lays out the text of a tab item and its icon if it has one. It aligns the text and the icon
// based on the position of the tab.
func (t *TabItem) iconText(gtx layout.Context, tabPosition Position, isMinimized bool) layout.Dimensions {
	widgetAxis := layout.Vertical
	if tabPosition == Left || tabPosition == Right {
		widgetAxis = layout.Horizontal
	}

	axis := widgetAxis
	leftInset := float32(15)
	if isMinimized {
		axis = layout.Vertical
		leftInset = 0
	}

	dims := layout.Flex{}.Layout(gtx, layout.Rigid(func(gtx C) D {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {

			//gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: axis, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx C) D {
							dim := gtx.Px(unit.Dp(20))
							gtx.Constraints.Max = image.Point{X: dim, Y: dim}
							return t.layoutIcon(gtx)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Inset{
							Left: unit.Dp(leftInset),
						}.Layout(gtx, func(gtx C) D {
							lbl := t.label
							if isMinimized {
								lbl.TextSize = values.TextSize10
							}
							return lbl.Layout(gtx)
						})
					})
				}),
			)
		})
	}))
	return dims
}

func NewTabItem(title string, icon, inactiveIcon *image.Image) TabItem {
	tabItem := TabItem{
		Title: title,
		ID:    title,
	}

	if icon != nil {
		tabItem.iconOp = paint.NewImageOp(*icon)
		tabItem.inactiveIconOp = paint.NewImageOp(*inactiveIcon)
	}

	return tabItem
}

func (t *TabItem) Layout(gtx layout.Context, selected int, tabPosition Position, width int, isMinimized bool) layout.Dimensions {
	var tabWidth, tabHeight int
	var iconTextDims layout.Dimensions

	gtx.Constraints.Min.X = width
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			iconTextDims = t.iconText(gtx, tabPosition, isMinimized)
			if iconTextDims.Size.X > adaptiveTabWidth {
				adaptiveTabWidth = iconTextDims.Size.X
			}
			return iconTextDims
		}),
		layout.Expanded(func(gtx C) D {
			if tabPosition == Left || tabPosition == Right {
				gtx.Constraints.Min.X = adaptiveTabWidth
			}
			tabWidth, tabHeight = tabIndicatorDimensions(iconTextDims, tabPosition)
			b := t.button.Layout(gtx)
			return b
		}),
		layout.Expanded(func(gtx C) D {
			if selected != t.index {
				return layout.Dimensions{}
			}
			return indicatorDirection(tabPosition).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return line(gtx, tabWidth, tabHeight, keyblue)
					},
					))
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
	prevEvents       int
	events           []widget.ChangeEvent
	// btns             []*widget.Clickable
	title                   Label
	list                    *layout.List
	Position                Position
	Separator               bool
	iconButton              IconButton
	scrollLeft              *widget.Clickable
	scrollRight             *widget.Clickable
	theme                   *Theme
	Background              color.RGBA
	isNavTabs               bool
	isNavDrawerMinimized    bool
	minimizeNavDrawerButton IconButton
	maximizeNavDrawerButton IconButton
}

func NewTabs(th *Theme, isNavTabs bool) *Tabs {
	t := &Tabs{
		theme:                   th,
		list:                    &layout.List{},
		Position:                Left,
		scrollLeft:              new(widget.Clickable),
		scrollRight:             new(widget.Clickable),
		iconButton:              th.IconButton(new(widget.Clickable), new(widget.Icon)),
		flex:                    layout.Flex{},
		Background:              th.Color.Surface,
		isNavTabs:               isNavTabs,
		isNavDrawerMinimized:    false,
		minimizeNavDrawerButton: th.PlainIconButton(new(widget.Clickable), th.navigationArrowBack),
		maximizeNavDrawerButton: th.PlainIconButton(new(widget.Clickable), th.navigationArrowForward),
	}
	t.minimizeNavDrawerButton.Color = th.Color.Gray
	t.maximizeNavDrawerButton.Color = th.Color.Gray

	return t
}

// SetTabs creates a button widget for each tab item.
func (t *Tabs) SetTabs(tabs []TabItem) {
	t.items = tabs

	for i := range tabs {
		l := t.theme.Body1(t.items[i].Title)
		t.items[i].label = l
		b := t.theme.Button(new(widget.Clickable), "")
		b.Background = color.RGBA{}
		tabs[i].button = b
	}
}

// SelectedID returns the ID of the currently selected tab item
func (t *Tabs) SelectedID() string {
	if len(t.items) == 0 {
		return ""
	}
	return t.items[t.Selected].ID
}

// scrollButton lays out the right and left scroll buttons of the tab when Position is Horizontal.
func (t *Tabs) scrollButton(right bool, button *widget.Clickable) layout.FlexChild {
	show := false
	icon := mustIcon(widget.NewIcon(icons.NavigationChevronLeft))
	if right && t.list.Position.BeforeEnd {
		show = true
		icon = mustIcon(widget.NewIcon(icons.NavigationChevronRight))
	}

	if !right && t.list.Position.Offset > 0 {
		show = true
	}
	return layout.Rigid(func(gtx C) D {
		if (t.Position == Top || t.Position == Bottom) && show {
			t.iconButton.Icon = icon
			t.iconButton.Size = unit.Dp(20)
			t.iconButton.Color = rgb(0xbbbbbb)
			t.iconButton.Background = color.RGBA{}
			t.iconButton.Button = button
			return t.iconButton.Layout(gtx)
		}
		return layout.Dimensions{}
	})
}

// SetTitle setting the title of the tabs
func (t *Tabs) SetTitle(title Label) {
	t.title = title
}

// tabsTitle lays out the title of the tab when Position is Horizontal.
func (t *Tabs) tabsTitle() layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		if (t.Position == Top || t.Position == Bottom) && t.title.Text != "" {
			return layout.Inset{Top: values.MarginPadding10, Right: values.MarginPadding5, Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
				return t.title.Layout(gtx)
			})
		}
		return layout.Dimensions{}
	})
}

// contentTabPosition determines the order of the tab and page content depending on the tab Position.
func (t *Tabs) contentTabPosition(body layout.Widget) (widgets []layout.FlexChild) {
	var content, tab layout.FlexChild

	widgets = make([]layout.FlexChild, 2)
	content = layout.Flexed(1, body)
	tab = layout.Rigid(func(gtx C) D {
		dims := layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return Card{Color: t.Background}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: t.list.Axis, Spacing: layout.SpaceBetween}.Layout(gtx,
						t.tabsTitle(),
						t.scrollButton(false, t.scrollLeft),
						layout.Flexed(1, func(gtx C) D {
							mt := values.MarginPaddingMinus10
							ml := values.MarginPadding10
							if t.Position == Right || t.Position == Left {
								mt = values.MarginPadding0
								ml = values.MarginPadding0
							}

							return layout.Inset{Left: ml, Top: mt}.Layout(gtx, func(gtx C) D {
								return t.list.Layout(gtx, len(t.items), func(gtx C, i int) D {
									t.items[i].index = i
									col := t.Background
									if t.items[i].index == t.Selected {
										col = t.theme.Color.Background
									}

									width := gtx.Constraints.Min.X
									isMinimized := false
									if t.isNavTabs {
										if t.isNavDrawerMinimized {
											width = navDrawerMinimizedWidth
											isMinimized = true
										} else {
											width = navDrawerWidth
											isMinimized = false
										}
									}

									return Card{Color: col}.Layout(gtx, func(gtx C) D {
										return t.items[i].Layout(gtx, t.Selected, t.Position, width, isMinimized)
									})
								})
							})
						}),
						t.scrollButton(true, t.scrollRight),
					)
				})
			}),
			layout.Expanded(func(gtx C) D {
				direction := layout.E
				if t.Position == Right {
					direction = layout.W
				}
				// display separator only if Separator is true and tab is vertical
				if t.Separator && (t.Position == Right || t.Position == Left) {
					return direction.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return line(gtx, 2, gtx.Constraints.Max.Y, rgb(0xefefef))
							}),
						)
					})
				}
				return layout.Dimensions{}
			}),
			layout.Expanded(func(gtx C) D {
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
				return layout.SE.Layout(gtx, func(gtx C) D {
					btn := t.minimizeNavDrawerButton
					if t.isNavDrawerMinimized {
						btn = t.maximizeNavDrawerButton
					}
					return btn.Layout(gtx)
				})
			}),
		)
		return dims
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
	n := copy(t.events, t.events[t.prevEvents:])
	t.events = t.events[:n]
	t.prevEvents = n
	t.Selected = index
	if t.previousSelected != t.Selected {
		t.events = append(t.events, widget.ChangeEvent{})
		t.previousSelected = index
	}
}

// ChangeEvent returns the last change event
func (t *Tabs) ChangeEvent() bool {
	if len(t.events) == 0 {
		return false
	}
	n := copy(t.events, t.events[1:])
	t.events = t.events[:n]
	if t.prevEvents > 0 {
		t.prevEvents--
	}
	return true
}

func (t *Tabs) processChangeEvent() {
	for i := range t.items {
		if t.items[i].button.Button.Clicked() {
			t.ChangeTab(i)
			return
		}
	}

	for t.minimizeNavDrawerButton.Button.Clicked() {
		t.isNavDrawerMinimized = true
	}

	for t.maximizeNavDrawerButton.Button.Clicked() {
		t.isNavDrawerMinimized = false
	}

}

func (t *Tabs) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	t.processChangeEvent()

	for t.scrollRight.Clicked() {
		t.list.Position.Offset += 60
	}

	for t.scrollLeft.Clicked() {
		t.list.Position.Offset -= 60
	}

	switch t.Position {
	case Top, Bottom:
		t.list.Axis = layout.Horizontal
		t.flex.Axis = layout.Vertical
	default:
		t.list.Axis = layout.Vertical
		t.flex.Axis = layout.Horizontal
	}

	widgets := t.contentTabPosition(body)
	t.flex.Spacing = layout.SpaceBetween
	return t.flex.Layout(gtx, widgets...)
}
