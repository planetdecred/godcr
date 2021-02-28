package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Collapsible struct {
	isExpanded      bool
	button          *widget.Clickable
	BackgroundColor color.NRGBA
	card            Card
	expandedIcon    *widget.Icon
	collapsedIcon   *widget.Icon
}

type CollapsibleWithOption struct {
	isExpanded      bool
	button          *widget.Clickable
	BackgroundColor color.NRGBA
	card            Card
	expandedIcon    *widget.Image
	collapsedIcon   *widget.Image
	moreIconButton  IconButton
}

func (t *Theme) Collapsible() *Collapsible {
	c := &Collapsible{
		BackgroundColor: t.Color.Surface,
		button:          new(widget.Clickable),
		card:            t.Card(),
		expandedIcon:    t.chevronUpIcon,
		collapsedIcon:   t.chevronDownIcon,
	}
	c.card.Color = c.BackgroundColor
	return c
}

func (t *Theme) CollapsibleWithOption() *CollapsibleWithOption {
	c := &CollapsibleWithOption{
		BackgroundColor: t.Color.Surface,
		button:          new(widget.Clickable),
		card:            t.Card(),
		expandedIcon:    t.expandIcon,
		collapsedIcon:   t.collapseIcon,
		moreIconButton: IconButton{
			IconButtonStyle: material.IconButtonStyle{
				Button:     new(widget.Clickable),
				Icon:       t.NavMoreIcon,
				Size:       unit.Dp(25),
				Background: color.NRGBA{},
				Color:      t.Color.Text,
				Inset:      layout.UniformInset(unit.Dp(0)),
			},
		},
	}
	c.card.Color = c.BackgroundColor
	return c
}

func (c *Collapsible) Layout(gtx layout.Context, header, body func(C) D) layout.Dimensions {
	for c.button.Clicked() {
		c.isExpanded = !c.isExpanded
	}

	icon := c.collapsedIcon
	if c.isExpanded {
		icon = c.expandedIcon
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
								return layout.Inset{
									Right: unit.Dp(10),
								}.Layout(gtx, func(C) D {
									return icon.Layout(gtx, unit.Dp(20))
								})
							}),
						)
					}),
					layout.Expanded(c.button.Layout),
				)
			}),
			layout.Rigid(func(gtx C) D {
				if c.isExpanded {
					return body(gtx)
				}
				return D{}
			}),
		)
	})
}

func (c *CollapsibleWithOption) Layout(gtx layout.Context, header, body func(C) D, more func(C)) layout.Dimensions {
	for c.button.Clicked() {
		c.isExpanded = !c.isExpanded
	}

	icon := c.collapsedIcon
	if c.isExpanded {
		icon = c.expandedIcon
	}

	return c.card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return layout.Stack{}.Layout(gtx,
								layout.Stacked(func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											icon.Scale = 1
											return icon.Layout(gtx)
										}),
										layout.Rigid(header),
									)
								}),
								layout.Expanded(c.button.Layout),
							)
						}),
						layout.Rigid(func(gtx C) D {
							more(gtx)
							return c.moreIconButton.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if c.isExpanded {
					return body(gtx)
				}
				return D{}
			}),
		)
	})
}

func (c *CollapsibleWithOption) MoreTriggered() bool {
	return c.moreIconButton.Button.Clicked()
}
