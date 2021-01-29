package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Collapsible struct {
	IsExpanded      bool
	Button          *widget.Clickable
	BackgroundColor color.NRGBA
	expandedIcon    *widget.Icon
	collapsedIcon   *widget.Icon
	card            Card
}

type CollapsibleWithOption struct {
	moreIcon        IconButton
	backgroundColor color.NRGBA
	button          *widget.Clickable
	isExpanded      bool
	expandedIcon    *widget.Image
	collapsedIcon   *widget.Image
	card            Card
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

func (t *Theme) CollapsibleWithOption() *CollapsibleWithOption {
	expandedIcon := t.expandIcon
	collapsedIcon := t.collapseIcon

	expandedIcon.Scale = 1
	collapsedIcon.Scale = 1

	return &CollapsibleWithOption{
		backgroundColor: t.Color.Surface,
		collapsedIcon:   collapsedIcon,
		expandedIcon:    expandedIcon,
		card:            t.Card(),
		button:          new(widget.Clickable),
		moreIcon: IconButton{
			IconButtonStyle: material.IconButtonStyle{
				Button:     new(widget.Clickable),
				Icon:       t.navMoreIcon,
				Size:       unit.Dp(25),
				Background: color.NRGBA{},
				Color:      t.Color.Text,
				Inset:      layout.UniformInset(unit.Dp(0)),
			},
		},
	}
}

func (c *Collapsible) Layout(gtx layout.Context, header func(C) D, content func(C) D) layout.Dimensions {
	icon := c.collapsedIcon
	if c.IsExpanded {
		icon = c.expandedIcon
	}

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
								return layout.Inset{
									Right: unit.Dp(10),
								}.Layout(gtx, func(C) D {
									return icon.Layout(gtx, unit.Dp(20))
								})
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

func (c *CollapsibleWithOption) Layout(gtx layout.Context, header, content func(C) D, more func(C)) layout.Dimensions {
	for c.button.Clicked() {
		c.isExpanded = !c.isExpanded
	}

	icon := c.collapsedIcon
	if c.isExpanded {
		icon = c.expandedIcon
	}

	headerFlex := layout.Rigid(func(gtx C) D {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							return layout.Flex{}.Layout(gtx, layout.Rigid(icon.Layout), layout.Rigid(header))
						}),
						layout.Expanded(c.button.Layout),
					)
				}),
				layout.Rigid(func(gtx C) D {
					more(gtx)
					return c.moreIcon.Layout(gtx)
				}),
			)
		})
	})

	children := []layout.FlexChild{headerFlex}
	if c.isExpanded {
		children = append(children, layout.Rigid(content))
	}

	return c.card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	})
}

func (c *CollapsibleWithOption) MoreTriggered() bool {
	return c.moreIcon.Button.Clicked()
}
