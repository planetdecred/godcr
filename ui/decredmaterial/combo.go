package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Combo struct {
	items          []ComboItem
	isOpen         bool
	selectedIndex  int
	color          color.RGBA
	background     color.RGBA
	chevronIcon    *widget.Icon
	navigationIcon *widget.Icon
	backdrop       *widget.Clickable
}

type ComboItem struct {
	Text   string
	Icon   image.Image
	button Button
	label  Label
}

func (t *Theme) Combo(items []ComboItem) *Combo {
	c := &Combo{
		isOpen:         false,
		items:          make([]ComboItem, len(items)+1),
		color:          t.Color.Background,
		background:     t.Color.Surface,
		chevronIcon:    t.chevronDownIcon,
		navigationIcon: t.NavigationCheckIcon,
		backdrop:       new(widget.Clickable),
	}

	for i := range items {
		items[i].button = t.Button(new(widget.Clickable), items[i].Text)
		items[i].label = t.Body1(items[i].Text)
		c.items[i+1] = items[i]
	}

	if len(c.items) > 0 {
		c.items[0] = ComboItem{
			Text:   items[0].Text,
			Icon:   items[0].Icon,
			label:  t.Body1(items[0].Text),
			button: t.Button(new(widget.Clickable), items[0].Text),
		}
		c.selectedIndex = 1
	}

	return c
}

func (c *Combo) Selected() string {
	return c.items[c.SelectedIndex()].Text
}

func (c *Combo) SelectedIndex() int {
	return c.selectedIndex - 1
}

func (c *Combo) handleEvents() {
	for c.items[0].button.Button.Clicked() {
		c.isOpen = !c.isOpen
	}

	for i := range c.items {
		index := i
		if index != 0 {
			for c.items[index].button.Button.Clicked() {
				c.selectedIndex = index
				c.items[0].label.Text = c.items[index].Text
				c.isOpen = false
			}
		}
	}

	for c.backdrop.Clicked() {

	}
}

func (c *Combo) Changed() bool {
	for i := range c.items {
		index := i
		if index != 0 {
			for c.items[index].button.Button.Clicked() {
				if c.items[0].label.Text != c.items[index].Text {
					return true
				}
			}
		}
	}

	return false
}

func (c *Combo) layoutIcon(itemIndex int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		if c.items[itemIndex].Icon == nil {
			return layout.Dimensions{}
		}

		img := widget.Image{Src: paint.NewImageOp(c.items[itemIndex].Icon)}
		img.Scale = 0.045

		return img.Layout(gtx)
	})
}

func (c *Combo) layoutText(index int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		gtx.Constraints.Min.X = 80
		return layout.Inset{
			Right: unit.Dp(15),
			Left:  unit.Dp(5),
		}.Layout(gtx, func(gtx C) D {
			return c.items[index].label.Layout(gtx)
		})
	})
}

func (c *Combo) layoutActiveIcon(index int, isFirstOption bool) layout.FlexChild {
	var icon *widget.Icon
	if isFirstOption {
		icon = c.chevronIcon
	} else if index == c.selectedIndex {
		icon = c.navigationIcon
	}

	return layout.Rigid(func(gtx C) D {
		return layout.E.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				if icon != nil {
					return icon.Layout(gtx, unit.Dp(20))
				}
				return layout.Dimensions{}
			})
		})
	})
}

func (c *Combo) layoutOption(gtx layout.Context, itemIndex int, isFirstOption bool) layout.Dimensions {
	btn := c.items[itemIndex].button

	min := gtx.Constraints.Min
	min.X = 100

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min = min
				iconLayout := c.layoutIcon(itemIndex)
				textLayout := c.layoutText(itemIndex)
				activeIconLayout := c.layoutActiveIcon(itemIndex, isFirstOption)

				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, iconLayout, textLayout, activeIconLayout)
			})
		}),
		layout.Expanded(btn.Button.Layout),
	)
}

func (c *Combo) Layout(gtx layout.Context) layout.Dimensions {
	c.handleEvents()

	children := []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return c.layoutOption(gtx, 0, true)
		}),
	}

	if c.isOpen {
		return c.comboItemMenu(gtx)
	}
	return c.drawlayout(gtx, false, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}

func (c *Combo) comboItemMenu(gtx layout.Context) layout.Dimensions {
	items := c.items[1:]
	var comboItemRows []func(gtx C) D
	for i := range items {
		index := i
		comboItemRows = append(comboItemRows, func(gtx C) D {
			return c.layoutOption(gtx, index+1, false)
		})
	}

	border := widget.Border{Color: c.color, CornerRadius: unit.Dp(10), Width: unit.Dp(2)}
	return border.Layout(gtx, func(gtx C) D {
		return c.drawlayout(gtx, true, func(gtx C) D {
			list := &layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(comboItemRows), func(gtx C, i int) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.UniformInset(unit.Dp(0)).Layout(gtx, comboItemRows[i])
					}),
					layout.Rigid(func(gtx C) D {
						// if i < len(comboItemRows)-1 {
						// 	return layout.Inset{
						// 		Top:    unit.Dp(10),
						// 		Bottom: unit.Dp(10),
						// 	}.Layout(gtx, func(gtx C) D {
						// 		return c.line.Layout(gtx)
						// 	})
						// }

						return layout.Dimensions{}
					}),
				)
			})
		})
	})
}

// drawlayout wraps the page tx and sync section in a card layout
func (c *Combo) drawlayout(gtx layout.Context, isPopUp bool, body layout.Widget) layout.Dimensions {
	color := c.color
	m := unit.Dp(5)
	if isPopUp {
		color = c.background
		m = unit.Dp(15)
	}
	return Card{Color: color, Rounded: true}.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(m).Layout(gtx, body)
	})
}
