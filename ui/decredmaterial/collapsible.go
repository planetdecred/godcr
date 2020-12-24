package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type MoreItem struct {
	Text   string
	Button *widget.Clickable
	label  Label
}

type Collapsible struct {
	IsExpanded      bool
	Button          *widget.Clickable
	expandIcon      *widget.Icon
	BackgroundColor color.RGBA
	expandedIcon    *widget.Icon
	collapsedIcon   *widget.Icon
	card            Card
}

type CollapsibleWithOption struct {
	Items             []MoreItem
	Collapsible       *Collapsible
	MoreIcon          IconButton
	isOpened          bool
	color             color.RGBA
	selectedMoreIndex int
	theme             *Theme
}

func (t *Theme) Collapsible() *Collapsible {
	c := &Collapsible{
		BackgroundColor: t.Color.Surface,
		Button:          new(widget.Clickable),
		expandedIcon:    t.chevronUpIcon,
		collapsedIcon:   t.chevronDownIcon,
		card:            t.Card(),
	}
	c.card.Color = c.BackgroundColor

	return c
}

func (t *Theme) CollapsibleWithOption(Items []MoreItem) *CollapsibleWithOption {
	c := &CollapsibleWithOption{
		Items: make([]MoreItem, len(Items)+1),
		color: t.Color.Background,
		theme: t,
		Collapsible: &Collapsible{
			BackgroundColor: t.Color.Surface,
			Button:          new(widget.Clickable),
			expandedIcon:    t.chevronUpIcon,
			collapsedIcon:   t.chevronDownIcon,
		},
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
	}

	for i := range Items {
		Items[i].Button = new(widget.Clickable)
		Items[i].label = c.theme.Body1(Items[i].Text)
		c.Items[i+1] = Items[i]
	}

	return c
}

func (c *Collapsible) collapseIconLayout(gtx layout.Context) layout.Dimensions {
	icon := c.collapsedIcon
	if c.IsExpanded {
		icon = c.expandedIcon
	}

	return layout.Inset{
		Right: unit.Dp(10),
	}.Layout(gtx, func(C) D {
		return icon.Layout(gtx, unit.Dp(20))
	})
}

func (c *Collapsible) Layout(gtx layout.Context, header func(C) D, content func(C) D) layout.Dimensions {
	for c.Button.Clicked() {
		c.IsExpanded = !c.IsExpanded
	}

	return c.card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return header(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return c.collapseIconLayout(gtx)
							}),
						)
					}),
					layout.Expanded(c.Button.Layout),
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
}

func (c *CollapsibleWithOption) Layout(gtx layout.Context, header func(C) D, content func(C) D, footer func(C) D) layout.Dimensions {
	c.handleEvents()

	// dims := layout.Inset{Top: unit.Dp(15)}.Layout(gtx, func(gtx C) D {
	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			card := c.theme.Card()
			card.Color = c.Collapsible.BackgroundColor
			return card.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Flexed(0.93, func(gtx C) D {
									return layout.Stack{}.Layout(gtx,
										layout.Stacked(func(gtx C) D {
											return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													return c.Collapsible.collapseIconLayout(gtx)
												}),
												layout.Rigid(func(gtx C) D {
													return header(gtx)
												}),
											)
										}),
										layout.Expanded(c.Collapsible.Button.Layout),
									)
								}),
								layout.Flexed(0.07, func(gtx C) D {
									return layout.E.Layout(gtx, func(gtx C) D {
										return c.MoreIcon.Layout(gtx)
									})
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							if c.Collapsible.IsExpanded {
								return content(gtx)
							}
							return layout.Dimensions{}
						}),
					)
				})
			})
		}),
		layout.Rigid(func(gtx C) D {
			if footer != nil {
				return layout.Inset{Top: unit.Dp(-10)}.Layout(gtx, func(gtx C) D {
					card := c.theme.Card()
					card.Radius = CornerRadius{
						NE: 0,
						NW: 0,
						SW: 10,
						SE: 10,
					}
					card.Color = c.theme.Color.Orange
					return card.Layout(gtx, func(gtx C) D {
						return footer(gtx)
					})
				})
			}
			return layout.Dimensions{}
		}),
	)
	// })

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

func (c *CollapsibleWithOption) moreItemMenu(gtx layout.Context, body layout.Widget) layout.Dimensions {
	border := widget.Border{Color: c.color, CornerRadius: unit.Dp(10), Width: unit.Dp(2)}
	return layout.Inset{Top: unit.Dp(50)}.Layout(gtx, func(gtx C) D {
		return border.Layout(gtx, func(gtx C) D {
			card := c.theme.Card()
			card.Color = c.Collapsible.BackgroundColor
			return card.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(5)).Layout(gtx, body)
			})
		})
	})
}

func (c *CollapsibleWithOption) moreOption(gtx layout.Context) layout.Dimensions {
	return c.moreItemMenu(gtx, func(gtx C) D {
		list := &layout.List{Axis: layout.Vertical}
		Items := c.Items[1:]
		return list.Layout(gtx, len(Items), func(gtx C, i int) D {
			return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx C) D {
				index := i + 1
				btn := c.Items[index].Button
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
										return c.Items[index].label.Layout(gtx)
									})
								}),
							)
						})
					}),
					layout.Expanded(btn.Layout),
				)
			})
		})
	})
}

func (c *CollapsibleWithOption) SelectedIndex() int {
	return c.selectedMoreIndex - 1
}

func (c *CollapsibleWithOption) Selected() string {
	return c.Items[c.SelectedIndex()].Text
}

func (c *CollapsibleWithOption) Hide() {
	c.isOpened = false
}

func (c *CollapsibleWithOption) handleEvents() {
	if len(c.Items) > 0 {
		if c.MoreIcon.Button.Clicked() {
			c.isOpened = !c.isOpened
		}
	}

	for i := range c.Items {
		index := i
		if index != 0 {
			for c.Items[index].Button.Clicked() {
				c.selectedMoreIndex = index
				c.isOpened = false
			}
		}
	}
}
