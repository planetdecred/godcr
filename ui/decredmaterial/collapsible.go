package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/values"
)

type IconStyle uint8

const (
	// Chevron sets the icon design chevron icon.
	Chevron IconStyle = iota
	// Caret sets the  the icon design to caret icon.
	Caret
)

type IconPosition uint8

const (
	// After the chevron icon on the left of the header.
	After IconPosition = iota
	// Before sets the chevron icon on the right of the header.
	Before
)

type Collapsible struct {
	th              *Theme
	style           *values.ColorStyle
	iconColor       color.NRGBA
	isExpanded      bool
	button          *widget.Clickable
	BackgroundColor color.NRGBA
	IconStyle       IconStyle
	IconPosition    IconPosition
	card            Card
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
		th:        t,
		button:    new(widget.Clickable),
		card:      t.Card(),
		iconColor: t.Color.Gray1,
		style:     t.Styles.CollapsibleStyle,
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

func (c *Collapsible) Layout(gtx C, header, body func(C) D) D {
	for c.button.Clicked() {
		c.isExpanded = !c.isExpanded
	}

	var icon *Image
	if c.IconStyle == Caret {
		icon = c.th.expandIcon
		if c.isExpanded {
			icon = c.th.collapseIcon
		}
	} else if c.IconStyle == Chevron {
		icon = c.th.Icons.ChevronExpand
		if c.isExpanded {
			icon = c.th.Icons.ChevronCollapse
		}
	}

	c.card.Color = c.style.Background
	return c.card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return c.button.Layout(gtx, func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X

							var children []layout.FlexChild
							var spacing layout.Spacing

							if c.IconPosition == Before {
								spacing = layout.SpaceEnd
								children = append(children, layout.Rigid(func(gtx C) D {
									return layout.Inset{Right: values.MarginPadding9}.Layout(gtx, icon.Layout24dp)
								}))
							}

							children = append(children, layout.Rigid(header))

							if c.IconPosition == After {
								spacing = layout.SpaceBetween
								children = append(children, layout.Rigid(icon.Layout24dp))
							}

							return layout.Flex{Spacing: spacing}.Layout(gtx, children...)
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

func (c *Collapsible) IsExpanded() bool {
	return c.isExpanded
}

var rememberExpand map[int]bool

func (c *CollapsibleWithOption) Layout(gtx C, header, body func(C) D, more func(C), rowID int) D {
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
	c.card.Color = c.BackgroundColor
	return c.card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return c.button.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										// TODO needs to be centered vertically
										return icon.Layout24dp(gtx)
									}),
									layout.Rigid(header),
								)
							})
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
