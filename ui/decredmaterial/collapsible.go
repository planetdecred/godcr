package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/values"
)

type Collapsible struct {
	style           *values.ColorStyle
	iconColor       color.NRGBA
	isExpanded      bool
	button          *widget.Clickable
	BackgroundColor color.NRGBA
	card            Card
	expandedIcon    *Icon
	collapsedIcon   *Icon
}

type CollapsibleWithOption struct {
	isExpanded      bool
	button          *widget.Clickable
	BackgroundColor color.NRGBA
	card            Card
	expandedIcon    *Image
	collapsedIcon   *Image
	moreIconButton  IconButton
}

func (t *Theme) Collapsible() *Collapsible {
	c := &Collapsible{
		button:        new(widget.Clickable),
		card:          t.Card(),
		expandedIcon:  NewIcon(t.chevronUpIcon),
		collapsedIcon: NewIcon(t.chevronDownIcon),
		iconColor:     t.Color.Gray1,
		style:         t.Styles.CollapsibleStyle,
	}
	return c
}

func (t *Theme) CollapsibleWithOption() *CollapsibleWithOption {
	c := &CollapsibleWithOption{
		BackgroundColor: t.Color.Surface,
		button:          new(widget.Clickable),
		card:            t.Card(),
		expandedIcon:    t.collapseIcon,
		collapsedIcon:   t.expandIcon,
		moreIconButton: IconButton{
			IconButtonStyle{
				Button: new(widget.Clickable),
				Icon:   t.navMoreIcon,
				Size:   unit.Dp(25),
				Inset:  layout.UniformInset(unit.Dp(0)),
			},
			&values.ColorStyle{
				Background: color.NRGBA{},
				Foreground: t.Color.Text,
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

	icon.Color = c.iconColor
	c.card.Color = c.style.Background
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
								return icon.Layout(gtx, values.MarginPadding20)
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

func (c *Collapsible) IsExpanded() bool {
	return c.isExpanded
}

var rememberExpand map[int]bool

func (c *CollapsibleWithOption) Layout(gtx layout.Context, header, body func(C) D, more func(C), rowID int) layout.Dimensions {
	if rememberExpand == nil {
		rememberExpand = make(map[int]bool)
	}

	if c.button.Clicked() {
		rememberExpand[rowID] = !rememberExpand[rowID]
	}

	icon := c.collapsedIcon
	if rememberExpand[rowID] {
		icon = c.expandedIcon
	}

	return c.card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return layout.Stack{}.Layout(gtx,
								layout.Stacked(func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											// TODO needs to be centered vertically
											return icon.Layout24dp(gtx)
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
				if rememberExpand[rowID] {
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
