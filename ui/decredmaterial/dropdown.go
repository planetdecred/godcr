package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

var MaxWidth = unit.Dp(800)

type DropDown struct {
	items          []DropDownItem
	isOpen         bool
	selectedIndex  int
	color          color.NRGBA
	background     color.NRGBA
	dropdownIcon   *widget.Icon
	navigationIcon *widget.Icon
	backdrop       *widget.Clickable

	group            uint
	closeAllDropdown func(group uint)
	card             Card
}

type DropDownItem struct {
	Text   string
	Icon   *widget.Image
	button Button
	label  Label
}

func (t *Theme) DropDown(items []DropDownItem, group uint) *DropDown {
	c := &DropDown{
		isOpen:         false,
		items:          make([]DropDownItem, len(items)+1),
		color:          t.Color.Gray1,
		background:     t.Color.Surface,
		dropdownIcon:   t.dropDownIcon,
		navigationIcon: t.navigationCheckIcon,
		backdrop:       new(widget.Clickable),

		group:            group,
		closeAllDropdown: t.closeAllDropdownMenus,
		card:             t.Card(),
	}

	for i := range items {
		items[i].button = t.Button(new(widget.Clickable), items[i].Text)
		items[i].label = t.Body1(items[i].Text)
		c.items[i+1] = items[i]
	}

	if len(c.items) > 0 {
		c.items[0] = DropDownItem{
			Text:   items[0].Text,
			Icon:   items[0].Icon,
			label:  t.Body1(items[0].Text),
			button: t.Button(new(widget.Clickable), items[0].Text),
		}
		c.selectedIndex = 1
	}

	t.dropDownMenus = append(t.dropDownMenus, c)
	return c
}

func (c *DropDown) Selected() string {
	return c.items[c.SelectedIndex()+1].Text
}

func (c *DropDown) SelectedIndex() int {
	return c.selectedIndex - 1
}

func (c *DropDown) Len() int {
	return len(c.items) - 1
}

func (c *DropDown) handleEvents() {
	for c.items[0].button.Button.Clicked() {
		c.closeAllDropdown(c.group)
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
		c.closeAllDropdown(c.group)
	}
}

func (c *DropDown) Changed() bool {
	for i := range c.items {
		index := i
		if index != 0 {
			for c.items[index].button.Button.Clicked() {
				if c.items[0].label.Text != c.items[index].Text {
					c.selectedIndex = index
					c.items[0].label.Text = c.items[index].Text
					c.isOpen = false
					return true
				}
			}
		}
	}

	return false
}

func (c *DropDown) layoutIcon(itemIndex int) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		if c.items[itemIndex].Icon == nil {
			return D{}
		}

		img := c.items[itemIndex].Icon
		img.Scale = 0.045

		return img.Layout(gtx)
	})
}

func (c *DropDown) layoutText(index int) layout.FlexChild {
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

func (c *DropDown) layoutActiveIcon(index int, isFirstOption bool) layout.FlexChild {
	var icon *widget.Icon
	if isFirstOption {
		icon = c.dropdownIcon
	} else if index == c.selectedIndex {
		icon = c.navigationIcon
	}

	return layout.Rigid(func(gtx C) D {
		return layout.E.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				if icon != nil {
					return icon.Layout(gtx, unit.Dp(20))
				}
				return D{}
			})
		})
	})
}

func (c *DropDown) layoutOption(gtx C, itemIndex int, isFirstOption bool) D {
	btn := c.items[itemIndex].button
	min := gtx.Constraints.Min
	min.X = 100

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
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

func (c *DropDown) Layout(gtx C, dropPos int, reversePos bool) D {
	c.handleEvents()

	children := []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return c.layoutOption(gtx, 0, true)
		}),
	}
	position := dropPos
	if reversePos {
		width := gtx.Constraints.Max.X
		nw := (width * 800) / gtx.Px(MaxWidth)
		position = nw - dropPos
	}

	if c.isOpen {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				gtx.Constraints.Min = gtx.Constraints.Max
				return c.backdrop.Layout(gtx)
			}),
			layout.Stacked(func(gtx C) D {
				return layout.Inset{
					Left: unit.Dp(float32(position)),
				}.Layout(gtx, func(gtx C) D {
					return c.dropDownItemMenu(gtx)
				})
			}),
		)
	}
	return layout.Inset{
		Left: unit.Dp(float32(position)),
	}.Layout(gtx, func(gtx C) D {
		return c.drawLayout(gtx, false, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})
	})
}

func (c *DropDown) dropDownItemMenu(gtx C) D {
	items := c.items[1:]
	var dropDownItemRows []layout.Widget
	for i := range items {
		index := i
		dropDownItemRows = append(dropDownItemRows, func(gtx C) D {
			return c.layoutOption(gtx, index+1, false)
		})
	}

	border := widget.Border{Color: c.color, CornerRadius: unit.Dp(10), Width: unit.Dp(2)}
	return border.Layout(gtx, func(gtx C) D {
		return c.drawLayout(gtx, true, func(gtx C) D {
			list := &layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(dropDownItemRows), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, dropDownItemRows[i])
			})
		})
	})
}

// drawLayout wraps the page tx and sync section in a card layout
func (c *DropDown) drawLayout(gtx C, isPopUp bool, body layout.Widget) D {
	color := c.color
	m := unit.Dp(5)
	if isPopUp {
		color = c.background
		m = unit.Dp(15)
	}
	c.card.Color = color
	return c.card.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(m).Layout(gtx, body)
	})
}
