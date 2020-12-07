package decredmaterial

import (
	// "fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Collapsible struct {
	items          []MoreItem
	IsExpanded            bool
	Button                *widget.Clickable
	MoreIcon              IconButton
	isOpened              bool
	BackgroundColor color.RGBA
	color                 color.RGBA
		theme       *Theme
}

type MoreItem struct {
	Text   string
	button Button
	label  Label
}

func (t *Theme) Collapsible(button *widget.Clickable) *Collapsible {
	c := &Collapsible{
		BackgroundColor: t.Color.Surface,
		MoreIcon: IconButton{
			material.IconButtonStyle{
				Icon:       t.navMoreIcon,
				Size:       unit.Dp(25),
				Background: color.RGBA{},
				Color:      t.Color.Text,
				Inset:      layout.UniformInset(unit.Dp(0)),
				Button:     new(widget.Clickable),
			},
		},
		Button: button,
		color: t.Color.Background,
		theme:      t,
	}
	
	return c
}

func (c *Collapsible) layoutHeader(gtx layout.Context, header func(C) D) layout.Dimensions {
	dims := layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return header(gtx)
		}),
	)

	return dims
}

func (c *Collapsible) Layout(gtx layout.Context, header func(C) D, content func(C) D) layout.Dimensions {
	c.handleEvents()

	dims := layout.Inset{Top: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
		return Card{Color: c.BackgroundColor, Rounded: true}.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Stack{}.Layout(gtx,
									layout.Stacked(func(gtx C) D {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X - 30
										return c.layoutHeader(gtx, header)
									}),
									layout.Expanded(c.Button.Layout),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return c.MoreIcon.Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						if c.IsExpanded {
							return content(gtx)
						}
						return layout.Dimensions{}
					}),
				)
			})
		})
	})

	return layout.Stack{Alignment: layout.NE}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return dims
		}),
		layout.Stacked(func(gtx C) D {
			if c.isOpened {
				return c.moreOption(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}

func (c *Collapsible) moreItemMenu(gtx layout.Context, body layout.Widget) layout.Dimensions {
	border := widget.Border{Color: c.color, CornerRadius: unit.Dp(10), Width: unit.Dp(2)}
	return layout.Inset{Top: unit.Dp(50)}.Layout(gtx, func(gtx C) D {
		return border.Layout(gtx, func(gtx C) D {
			return Card{Color: c.BackgroundColor, Rounded: true}.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(5)).Layout(gtx, body)
			})
		})
	})
}

func (c *Collapsible) moreOption(gtx layout.Context) layout.Dimensions {
	items := c.items[1:]
	var moreItemRows []func(gtx C) D
	for i := range items {
		index := i+1
		moreItemRows = append(moreItemRows, func(gtx C) D {
			btn := c.items[index].button
			min := gtx.Constraints.Min
			min.X = 100

			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min = min
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								gtx.Constraints.Min.X = 80
								return layout.Inset{
									Right: unit.Dp(15),
									Left:  unit.Dp(5),
								}.Layout(gtx, func(gtx C) D {
									return c.items[index].label.Layout(gtx)
								})
							}),
						)
					})
				}),
				layout.Expanded(btn.Button.Layout),
			)
		})
	}

	border := widget.Border{Color: c.color, CornerRadius: unit.Dp(10), Width: unit.Dp(2)}
	return border.Layout(gtx, func(gtx C) D {
		return c.moreItemMenu(gtx, func(gtx C) D {
			list := &layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(moreItemRows), func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, moreItemRows[i])
			})
		})
	})
}

// SetTabs creates a button widget for each tab item.
func (c *Collapsible) AddItems(items []MoreItem) {
	c.items = items

	for i := range items {
		items[i].button = c.theme.Button(new(widget.Clickable), items[i].Text)
		items[i].label = c.theme.Body1(items[i].Text)
		c.items[i+1] = items[i]
	}
}

func (c *Collapsible) handleEvents() {
	for c.Button.Clicked() {
		c.IsExpanded = !c.IsExpanded
	}

	if len(c.items) > 0 {
		if c.MoreIcon.Button.Clicked() {
			c.isOpened = !c.isOpened
		}
	}

	for i := range c.items {
		index := i
		if index != 0 {
			for c.items[index].button.Button.Clicked() {
				c.isOpened = false
			}
		}
	}
}
