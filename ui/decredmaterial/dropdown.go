package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/values"
)

var MaxWidth = unit.Dp(800)

type DropDown struct {
	theme          *Theme
	items          []DropDownItem
	isOpen         bool
	revs           bool
	selectedIndex  int
	color          color.NRGBA
	background     color.NRGBA
	dropdownIcon   *widget.Icon
	navigationIcon *widget.Icon
	backdrop       *widget.Clickable

	group            uint
	closeAllDropdown func(group uint)
	Width            int
	linearLayout     *LinearLayout
	shadow           *Shadow
}

type DropDownItem struct {
	Text      string
	Icon      *Image
	clickable *Clickable
	label     Label
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
		linearLayout: &LinearLayout{
			Width:  WrapContent,
			Height: WrapContent,
			Border: Border{Radius: Radius(8)},
		},
		shadow: t.Shadow(),
	}

	for i := range items {
		items[i].clickable = t.NewClickable(true)
		items[i].label = t.Body1(items[i].Text)
		c.items[i+1] = items[i]
	}

	if len(c.items) > 0 {
		txt := items[0].Text
		if len(items[0].Text) > 12 {
			txt = items[0].Text[:12] + "..."
		}

		c.items[0] = DropDownItem{
			Text:      items[0].Text,
			Icon:      items[0].Icon,
			label:     t.Body1(txt),
			clickable: t.NewClickable(true),
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
	for c.items[0].clickable.Clicked() {
		c.closeAllDropdown(c.group)
		c.isOpen = !c.isOpen
	}

	for i := range c.items {
		index := i
		if index != 0 {
			for c.items[index].clickable.Clicked() {
				c.selectedIndex = index
				txt := c.items[index].Text
				if len(c.items[index].Text) > 12 {
					txt = c.items[index].Text[:12] + "..."
				}
				c.items[0].label.Text = txt
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
			for c.items[index].clickable.Clicked() {
				if c.items[0].label.Text != c.items[index].Text {
					c.selectedIndex = index
					txt := c.items[index].Text
					if len(c.items[index].Text) > 12 {
						txt = c.items[index].Text[:12] + "..."
					}
					c.items[0].label.Text = txt
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
	clickable := c.items[itemIndex].clickable
	width := gtx.Px(values.MarginPadding174)
	if c.revs {
		width = gtx.Px(values.MarginPadding140)
	}

	radius := Radius(8)
	if itemIndex == 1 {
		radius.BottomRight = 0
		radius.BottomLeft = 0
	} else if itemIndex == len(c.items[1:]) {
		radius.TopLeft = 0
		radius.TopRight = 0
	} else {
		radius = Radius(0)
	}

	return LinearLayout{
		Width:     width,
		Height:    WrapContent,
		Clickable: clickable,
		Padding:   layout.UniformInset(values.MarginPadding10),
		Border:    Border{Radius: radius},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if c.items[itemIndex].Icon == nil {
				return layout.Dimensions{}
			}

			img := c.items[itemIndex].Icon
			return img.Layout24dp(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X = gtx.Px(unit.Dp(110))
			if c.revs {
				gtx.Constraints.Max.X = gtx.Px(unit.Dp(100))
			}
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{
				Right: unit.Dp(5),
				Left:  unit.Dp(5),
			}.Layout(gtx, func(gtx C) D {
				return c.items[itemIndex].label.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return c.layoutActiveIcon(gtx, itemIndex, isFirstOption)
		}),
	)
}

func (c *DropDown) Layout(gtx C, dropPos int, reversePos bool) D {
	c.handleEvents()

	iLeft := dropPos
	iRight := 0
	alig := layout.NW
	c.revs = reversePos
	if reversePos {
		alig = layout.NE
		iLeft = 10
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
					return c.dropDownItemMenu(gtx)
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
					lay := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return c.layoutOption(gtx, 0, true)
						}))
					w := (lay.Size.X * 800) / gtx.Px(MaxWidth)
					c.Width = w + 10
					return lay
				})
			})
		}),
	)
}

func (c *DropDown) dropDownItemMenu(gtx C) D {
	return c.drawLayout(gtx, true, func(gtx C) D {
		list := &layout.List{Axis: layout.Vertical}
		return list.Layout(gtx, len(c.items[1:]), func(gtx C, index int) D {
			i := index + 1
			return c.layoutOption(gtx, i, false)
		})
	})
}

// drawLayout wraps the page tx and sync section in a card layout
func (d *DropDown) drawLayout(gtx C, isPopUp bool, body layout.Widget) D {
	if isPopUp {
		d.linearLayout.Background = d.background
		d.linearLayout.Shadow = nil
	} else {
		d.linearLayout.Background = d.color
		d.linearLayout.Shadow = nil
	}

	return d.linearLayout.Layout2(gtx, body)
}
