package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

var MaxWidth = unit.Dp(800)

type DropDown struct {
	theme          *Theme
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
	Width            int
}

type DropDownItem struct {
	Text   string
	Icon   *Image
	button *widget.Clickable
	label  Label
}

func (t *Theme) DropDown(items []DropDownItem, group uint) *DropDown {
	c := &DropDown{
		theme:          t,
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
		items[i].button = new(widget.Clickable)
		items[i].label = t.Body1(items[i].Text)
		c.items[i+1] = items[i]
	}

	if len(c.items) > 0 {
		c.items[0] = DropDownItem{
			Text:   items[0].Text,
			Icon:   items[0].Icon,
			label:  t.Body1(items[0].Text),
			button: new(widget.Clickable),
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
	for c.items[0].button.Clicked() {
		c.closeAllDropdown(c.group)
		c.isOpen = !c.isOpen
	}

	for i := range c.items {
		index := i
		if index != 0 {
			for c.items[index].button.Clicked() {
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
			for c.items[index].button.Clicked() {
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

func (c *DropDown) layoutActiveIcon(gtx layout.Context, index int, isFirstOption bool) D {
	var icon *widget.Icon
	if isFirstOption {
		icon = c.dropdownIcon
	} else if index == c.selectedIndex {
		icon = c.navigationIcon
	}

	return layout.E.Layout(gtx, func(gtx C) D {
		if icon != nil {
			return icon.Layout(gtx, unit.Dp(20))
		}
		return layout.Dimensions{}
	})
}

func (c *DropDown) layoutOption(gtx layout.Context, itemIndex int, isFirstOption bool) D {
	btn := c.items[itemIndex].button
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Px(unit.Dp(120))
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if c.items[itemIndex].Icon == nil {
							return layout.Dimensions{}
						}

						img := c.items[itemIndex].Icon
						return img.Layout24dp(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Px(unit.Dp(75))
						return layout.Inset{
							Right: unit.Dp(15),
							Left:  unit.Dp(5),
						}.Layout(gtx, func(gtx C) D {
							return c.items[itemIndex].label.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return c.layoutActiveIcon(gtx, itemIndex, isFirstOption)
					}),
				)
			}),
			layout.Expanded(btn.Layout),
		)
	})
}

func (c *DropDown) Layout(gtx C, dropPos int, reversePos bool) D {
	c.handleEvents()

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return c.drawLayout(gtx, false, func(gtx C) D {
				return c.layoutOption(gtx, 0, true)
			})
		}),
	}

	iLeft := dropPos
	iRight := 0
	alig := layout.NW
	if reversePos {
		alig = layout.NE
		iLeft = 0
		iRight = dropPos
	}

	if c.isOpen {
		return layout.Stack{Alignment: alig}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				gtx.Constraints.Min = gtx.Constraints.Max
				return c.backdrop.Layout(gtx)
			}),
			layout.Stacked(func(gtx C) D {
				return layout.Inset{
					Left:  unit.Dp(float32(iLeft)),
					Right: unit.Dp(float32(iRight)),
				}.Layout(gtx, func(gtx C) D {
					lay := c.dropDownItemMenu(gtx)
					w := (lay.Size.X * 800) / gtx.Px(MaxWidth)
					c.Width = w
					return lay
				})
			}),
		)
	}
	return layout.Stack{Alignment: alig}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Inset{
				Left:  unit.Dp(float32(iLeft)),
				Right: unit.Dp(float32(iRight)),
			}.Layout(gtx, func(gtx C) D {
				return c.drawLayout(gtx, false, func(gtx C) D {
					lay := layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
					w := (lay.Size.X * 800) / gtx.Px(MaxWidth)
					c.Width = w + 10
					return lay
				})
			})
		}),
	)
}

func (c *DropDown) dropDownItemMenu(gtx C) D {
	border := widget.Border{Color: c.color, CornerRadius: unit.Dp(10), Width: unit.Dp(2)}
	return border.Layout(gtx, func(gtx C) D {
		return c.drawLayout(gtx, true, func(gtx C) D {
			list := &layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(c.items[1:]), func(gtx C, index int) D {
				i := index + 1
				card := c.theme.Card()
				card.Color = color.NRGBA{}
				card.Radius = Radius(0)
				return card.HovarableLayout(gtx, c.items[i].button, func(gtx C) D {
					return c.layoutOption(gtx, i, false)
				})
			})
		})
	})
}

// drawLayout wraps the page tx and sync section in a card layout
func (c *DropDown) drawLayout(gtx C, isPopUp bool, body layout.Widget) D {
	color := c.color
	if isPopUp {
		color = c.background
	}
	c.card.Color = color
	return c.card.Layout(gtx, body)
}
