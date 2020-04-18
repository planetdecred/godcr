package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/widget/material"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

const sizeOfIcon float32 = 24

type TabContainer struct {
	Tabs            []Tab
	buttonContainer *layout.List

	Font               text.Font
	TextSize           unit.Value
	shaper             text.Shaper
	ActiveButtonStripe color.RGBA
	// buttonContentDivider divides button tabs and button view.
	buttonContentDivider color.RGBA
	color                color.RGBA

	tabLocation  *TabLocation
	currenTab    *int
	maxDimension *layout.Dimensions
}

type TabLocation int

const (
	TabLocationLeading TabLocation = iota
	TabLocationTrailing
	TabLocationBottom
	TabLocationTop
)

func NewTabContainer(t *material.Theme, tabItems ...Tab) TabContainer {
	return TabContainer{
		Tabs:            tabItems,
		buttonContainer: &layout.List{Axis: layout.Vertical},

		ActiveButtonStripe:   t.Color.Primary,
		buttonContentDivider: color.RGBA{200, 255, 144, 255},
		color:                t.Color.Text,
		TextSize:             t.TextSize.Scale(1),
		shaper:               t.Shaper,

		currenTab:    new(int),
		tabLocation:  new(TabLocation),
		maxDimension: new(layout.Dimensions),
	}
}

func (b TabContainer) Layout(gtx *layout.Context) {
	for i, Tab := range b.Tabs {
		Tab := Tab
		i := i
		for Tab.Button.Clicked(gtx) {
			*b.currenTab = i
		}
	}

	var tabContainerAxis layout.Axis

	switch *b.tabLocation {
	case TabLocationLeading, TabLocationTrailing:
		b.buttonContainer.Axis = layout.Vertical
		tabContainerAxis = layout.Horizontal
	default:
		b.buttonContainer.Axis = layout.Horizontal
		tabContainerAxis = layout.Vertical
	}

	buttonLayout := material.ButtonLayout{
		Inset:        layout.UniformInset(unit.Dp(0)),
		Background:   color.RGBA{},
		CornerRadius: unit.Dp(0),
	}

	buttonWidget := make([]func(), len(b.Tabs)+1)
	// Set a spacer so that users can identify
	// as a scrollable tabcontainer.
	//
	// https://material.io/components/tabs/#scrollable-tabs
	buttonWidget[0] = func() {
		layout.UniformInset(unit.Dp(30)).Layout(gtx, func() {})
	}

	for i, Tab := range b.Tabs {
		Tab := Tab
		i := i

		buttonWidget[i+1] = func() {
			buttonLayout.Layout(gtx, Tab.Button, func() {
				if *b.tabLocation == TabLocationLeading || *b.tabLocation == TabLocationTrailing {
					b.tabLayout(gtx, i, layout.Horizontal)
				} else {
					b.tabLayout(gtx, i, layout.Vertical)
				}
			})
		}
	}

	buttonContainer := layout.Rigid(func() {
		b.buttonContainer.Layout(gtx, len(buttonWidget), func(i int) {
			buttonWidget[i]()
		})
	})

	contentContainer := layout.Flexed(1, b.Tabs[*b.currenTab].Content())

	// Divides tab buttons and tab button.
	buttonContentDivider := layout.Rigid(func() {
		var x, y = gtx.Constraints.Width.Max, gtx.Px(unit.Dp(2))
		if *b.tabLocation == TabLocationLeading || *b.tabLocation == TabLocationTrailing {
			x, y = gtx.Px(unit.Dp(2)), gtx.Constraints.Height.Max
		}

		d := image.Point{X: x, Y: y}
		dr := f32.Rectangle{
			Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
		}

		paint.ColorOp{Color: b.buttonContentDivider}.Add(gtx.Ops)
		paint.PaintOp{Rect: dr}.Add(gtx.Ops)

		gtx.Dimensions = layout.Dimensions{Size: d}
	})

	tabFlexContainer := layout.Flex{Axis: tabContainerAxis}
	switch *b.tabLocation {
	case TabLocationLeading, TabLocationTop:
		tabFlexContainer.Layout(gtx, buttonContainer, buttonContentDivider, contentContainer)
	case TabLocationTrailing, TabLocationBottom:
		tabFlexContainer.Layout(gtx, contentContainer, buttonContentDivider, buttonContainer)
	}
}

func (b *TabContainer) Append(gtx *layout.Context, item Tab) {
	b.Tabs = append(b.Tabs, item)
	op.InvalidateOp{}.Add(gtx.Ops)
}

func (b TabContainer) ChangeTabIndex(gtx *layout.Context, index int) {
	if index >= len(b.Tabs) || index < 0 {
		return
	}

	*b.currenTab = index
	op.InvalidateOp{}.Add(gtx.Ops)
}

func (b TabContainer) ChangeTabLocation(gtx *layout.Context, tabLocation TabLocation) {
	*b.tabLocation = tabLocation
	op.InvalidateOp{}.Add(gtx.Ops)
}

func (b TabContainer) CurrentTabIndex() int {
	return *b.currenTab
}

func (b TabContainer) CurrentTabName() string {
	return b.Tabs[*b.currenTab].TabName
}

func (b *TabContainer) Pop(gtx *layout.Context, index int) {
	if len(b.Tabs)-1 == 0 {
		return
	}
	if index >= len(b.Tabs) || index < 0 {
		return
	}

	changeLocation := index - 1
	if changeLocation < 0 {
		changeLocation = 0
	}
	b.Tabs = append(b.Tabs[:index], b.Tabs[index+1:]...)
	*b.currenTab = changeLocation
	b.maxDimension.Size.X, b.maxDimension.Size.X = 0, 0
	op.InvalidateOp{}.Add(gtx.Ops)
}

func (b *TabContainer) Prepend(gtx *layout.Context, item Tab) {
	b.Tabs = append([]Tab{item}, b.Tabs...)
	*b.currenTab = *b.currenTab + 1

	op.InvalidateOp{}.Add(gtx.Ops)
}

func (b TabContainer) tabLayout(gtx *layout.Context, index int, axis layout.Axis) {
	buttonWithActiveIndicator := layout.Flex{Axis: axis, Alignment: layout.Middle}
	var maxDim layout.Dimensions

	// Add up size of icon and text.
	//
	// For the max X axis, we need only need to
	// get the max of max of either the text or icon
	getMax := func() {
		maxDim.Size.Y += gtx.Dimensions.Size.Y
		maxDim.Size.X += gtx.Dimensions.Size.X

		if b.maxDimension.Size.X < gtx.Dimensions.Size.X {
			b.maxDimension.Size.X = gtx.Dimensions.Size.X
		}
		if b.maxDimension.Size.Y < maxDim.Size.Y {
			b.maxDimension.Size.Y = maxDim.Size.Y
		}
	}

	buttonContent := layout.Rigid(func() {
		iconAndLabelContainer := layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}

		icon := layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(10), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, func() {
				// Set icon maximum size to 24 dp.
				pxSizeOfIcon := gtx.Px(unit.Dp(sizeOfIcon))
				if b.Tabs[index].Icon != nil {
					imgOp := *b.Tabs[index].Icon
					var max = imgOp.Size().Y
					if max < imgOp.Size().X {
						max = imgOp.Size().X
					}

					val := float32(pxSizeOfIcon) / float32(max)
					material.Image{Src: imgOp, Scale: val}.Layout(gtx)
					if gtx.Dimensions.Size.X < 96 {
						gtx.Dimensions.Size.X = 96
					}
					if gtx.Dimensions.Size.Y < 96 {
						gtx.Dimensions.Size.Y = 96
					}

					getMax()
				}
			})
		})

		label := layout.Rigid(func() {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
				gtx.Constraints.Width.Min = b.maxDimension.Size.X
				paint.ColorOp{Color: b.color}.Add(gtx.Ops)
				widget.Label{Alignment: text.Middle}.Layout(gtx, b.shaper, text.Font{}, b.TextSize, b.Tabs[index].TabName)
				getMax()
			})
		})

		if b.Tabs[index].Icon == nil {
			iconAndLabelContainer.Layout(gtx, label)
			return
		}

		iconAndLabelContainer.Layout(gtx, icon, label)
	})

	activeTabIndicator := layout.Rigid(func() {
		// Set stripline indicator to max length of icon, text and spacer.
		var stripLine = f32.Point{
			X: float32(gtx.Px(unit.Dp(4))),
			Y: float32(b.maxDimension.Size.Y + gtx.Px(unit.Dp(30))),
		}

		if axis == layout.Vertical {
			stripLine = f32.Point{
				X: float32(b.maxDimension.Size.X + gtx.Px(unit.Dp(20))),
				Y: float32(gtx.Px(unit.Dp(4))),
			}
		}

		if *b.currenTab == index {
			paint.ColorOp{Color: b.ActiveButtonStripe}.Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rectangle{Max: stripLine}}.Add(gtx.Ops)
		}

		gtx.Dimensions.Size.Y, gtx.Dimensions.Size.X = int(stripLine.Y), int(stripLine.X)
	})

	switch *b.tabLocation {
	case TabLocationLeading, TabLocationTop:
		buttonWithActiveIndicator.Layout(gtx, buttonContent, activeTabIndicator)
	case TabLocationTrailing, TabLocationBottom:
		buttonWithActiveIndicator.Layout(gtx, activeTabIndicator, buttonContent)
	}
}
