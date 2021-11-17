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
	clickable      *Clickable

	group            uint
	closeAllDropdown func(group uint)
	Width            int
	linearLayout     *LinearLayout
	padding          layout.Inset
	shadow           *Shadow
}

type DropDownItem struct {
	Text      string
	Icon      *Image
	clickable *Clickable
}

func (t *Theme) DropDown(items []DropDownItem, group uint) *DropDown {
	d := &DropDown{
		theme:          t,
		isOpen:         false,
		selectedIndex:  0,
		items:          make([]DropDownItem, 0),
		color:          t.Color.Gray1,
		background:     t.Color.Surface,
		dropdownIcon:   t.dropDownIcon,
		navigationIcon: t.navigationCheckIcon,
		backdrop:       new(widget.Clickable),
		clickable:      t.NewClickable(true),

		group:            group,
		closeAllDropdown: t.closeAllDropdownMenus,
		linearLayout: &LinearLayout{
			Width:  WrapContent,
			Height: WrapContent,
			Border: Border{Radius: Radius(8)},
		},
		padding: layout.Inset{Top: values.MarginPadding8, Bottom: values.MarginPadding8},
		shadow:  t.Shadow(),
	}

	d.clickable.HoverColor = Hovered(t.Color.InactiveGray)
	d.clickable.Radius = Radius(8)

	for i := range items {
		items[i].clickable = t.NewClickable(true)
		d.items = append(d.items, items[i])
	}

	t.dropDownMenus = append(t.dropDownMenus, d)
	return d
}

func (d *DropDown) Selected() string {
	return d.items[d.SelectedIndex()].Text
}

func (d *DropDown) SelectedIndex() int {
	return d.selectedIndex
}

func (d *DropDown) Len() int {
	return len(d.items)
}

func (d *DropDown) handleEvents() {
	if d.isOpen {
		for i := range d.items {
			index := i
			for d.items[index].clickable.Clicked() {
				d.selectedIndex = index
				d.isOpen = false
				break
			}
		}
	} else {
		for d.clickable.Clicked() {
			d.isOpen = true
		}
	}

	for d.backdrop.Clicked() {
		d.closeAllDropdown(d.group)
	}
}

func (d *DropDown) Changed() bool {
	if d.isOpen {
		for i := range d.items {
			index := i
			for d.items[index].clickable.Clicked() {
				d.selectedIndex = index
				d.isOpen = false
				return true
			}
		}
	}

	return false
}

func (d *DropDown) layoutActiveIcon(gtx layout.Context, index int) D {
	var icon *Icon
	if !d.isOpen {
		icon = NewIcon(d.dropdownIcon)
	} else if index == d.selectedIndex {
		icon = NewIcon(d.navigationIcon)
	}

	return layout.E.Layout(gtx, func(gtx C) D {
		if icon != nil {
			return icon.Layout(gtx, values.MarginPadding20)
		}
		return layout.Dimensions{}
	})
}

func (d *DropDown) layoutOption(gtx layout.Context, itemIndex int) D {
	item := d.items[itemIndex]

	width := gtx.Px(values.MarginPadding174)
	if d.revs {
		width = gtx.Px(values.MarginPadding140)
	}

	radius := Radius(0)
	clickable := item.clickable
	if !d.isOpen {
		radius = Radius(8)
		clickable = d.clickable
	}

	return LinearLayout{
		Width:     width,
		Height:    WrapContent,
		Clickable: clickable,
		Padding:   layout.UniformInset(values.MarginPadding10),
		Border:    Border{Radius: radius},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if item.Icon == nil {
				return layout.Dimensions{}
			}

			return item.Icon.Layout24dp(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Max.X = gtx.Px(unit.Dp(110))
			if d.revs {
				gtx.Constraints.Max.X = gtx.Px(unit.Dp(100))
			}
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{
				Right: unit.Dp(5),
				Left:  unit.Dp(5),
			}.Layout(gtx, func(gtx C) D {
				lbl := d.theme.Body2(item.Text)
				if !d.isOpen && len(item.Text) > 9 {
					lbl.Text = item.Text[:9] + "..."
				}
				return lbl.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return d.layoutActiveIcon(gtx, itemIndex)
		}),
	)
}

func (d *DropDown) Layout(gtx C, dropPos int, reversePos bool) D {
	d.handleEvents()

	iLeft := dropPos
	iRight := 0
	alig := layout.NW
	d.revs = reversePos
	if reversePos {
		alig = layout.NE
		iLeft = 10
		iRight = dropPos
	}

	if d.isOpen {
		return layout.Stack{Alignment: alig}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				gtx.Constraints.Min = gtx.Constraints.Max
				return d.backdrop.Layout(gtx)
			}),
			layout.Stacked(func(gtx C) D {
				return layout.Inset{
					Left:  unit.Dp(float32(iLeft)),
					Right: unit.Dp(float32(iRight)),
				}.Layout(gtx, func(gtx C) D {
					return d.dropDownItemMenu(gtx)
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
				return d.drawLayout(gtx, func(gtx C) D {
					lay := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return d.layoutOption(gtx, d.selectedIndex)
						}))
					w := (lay.Size.X * 800) / gtx.Px(MaxWidth)
					d.Width = w + 10
					return lay
				})
			})
		}),
	)
}

func (d *DropDown) dropDownItemMenu(gtx C) D {
	return d.drawLayout(gtx, func(gtx C) D {
		list := &layout.List{Axis: layout.Vertical}
		return list.Layout(gtx, len(d.items), func(gtx C, index int) D {
			return d.layoutOption(gtx, index)
		})
	})
}

// drawLayout wraps the page tx and sync section in a card layout
func (d *DropDown) drawLayout(gtx C, body layout.Widget) D {
	if d.isOpen {
		d.linearLayout.Background = d.background
		d.linearLayout.Padding = d.padding
		d.linearLayout.Shadow = d.shadow
	} else {
		d.linearLayout.Background = d.color
		d.linearLayout.Padding = layout.Inset{}
		d.linearLayout.Shadow = nil
	}

	return d.linearLayout.Layout2(gtx, body)
}

// Reslice the dropdowns
func ResliceDropdown(dropdowns []*DropDown, indexToRemove int) []*DropDown {
	dropdowns[indexToRemove] = dropdowns[len(dropdowns)-1]
	return dropdowns[:len(dropdowns)-1]
}

// Display one dropdown at a time
func DisplayOneDropdown(dropdowns ...*DropDown) {
	var menus []*DropDown
	for i, menu := range dropdowns {
		if menu.isOpen {
			menus = ResliceDropdown(dropdowns, i)
		}
		for _, newMenus := range menus {
			newMenus.isOpen = false
		}
	}
}
