package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type collapsibleCommon struct {
	isExpanded      bool
	button          *widget.Clickable
	BackgroundColor color.NRGBA
	card            Card
	expandedIcon    interface{}
	collapsedIcon   interface{}
	currentIcon     interface{}
}

type Collapsible struct {
	common         collapsibleCommon
	moreIconButton IconButton
	withOptions    bool
}

func (t *Theme) Collapsible() *Collapsible {
	c := &Collapsible{
		common: collapsibleCommon{
			BackgroundColor: t.Color.Surface,
			button:          new(widget.Clickable),
			card:            t.Card(),
			expandedIcon:    t.chevronUpIcon,
			collapsedIcon:   t.chevronDownIcon,
			currentIcon:     t.chevronDownIcon,
		},
	}
	c.common.card.Color = c.common.BackgroundColor
	return c
}

func (t *Theme) CollapsibleWithOption() *Collapsible {
	collapsible := t.Collapsible()
	collapsible.withOptions = true
	collapsible.moreIconButton = IconButton{
		IconButtonStyle: material.IconButtonStyle{
			Button:     new(widget.Clickable),
			Icon:       t.navMoreIcon,
			Size:       unit.Dp(25),
			Background: color.NRGBA{},
			Color:      t.Color.Text,
			Inset:      layout.UniformInset(unit.Dp(0)),
		},
	}
	collapsible.common.expandedIcon = t.expandIcon
	collapsible.common.collapsedIcon = t.collapseIcon
	collapsible.common.currentIcon = t.collapseIcon
	return collapsible
}

func (c *Collapsible) layoutContent(content func(C) D) layout.FlexChild {
	return layout.Rigid(func(gtx C) D {
		if c.common.isExpanded {
			return content(gtx)
		}
		return D{}
	})
}

func (c *Collapsible) layoutCollapsible(contents []func(C) D) []layout.FlexChild {
	if len(contents) != 2 {
		return []layout.FlexChild{}
	}

	return []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return contents[0](gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Right: unit.Dp(10),
							}.Layout(gtx, func(C) D {
								icon := c.common.currentIcon.(*widget.Icon)
								return icon.Layout(gtx, unit.Dp(20))
							})
						}),
					)
				}),
				layout.Expanded(c.common.button.Layout),
			)
		}),
		c.layoutContent(contents[1]),
	}
}

func (c *Collapsible) layoutCollapsibleWithOption(contents []func(C) D) []layout.FlexChild {
	if len(contents) != 3 {
		return []layout.FlexChild{}
	}

	return []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.Stack{}.Layout(gtx,
						layout.Stacked(func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									icon := c.common.currentIcon.(*widget.Image)
									icon.Scale = 1
									return icon.Layout(gtx)
								}),
								layout.Rigid(contents[0]))
						}),
						layout.Expanded(c.common.button.Layout),
					)
				}),
				layout.Rigid(func(gtx C) D {
					contents[2](gtx)
					return c.moreIconButton.Layout(gtx)
				}),
			)
		}),
		c.layoutContent(contents[1]),
	}
}

func (c *Collapsible) Layout(gtx layout.Context, contents ...func(C) D) D {
	c.handleEvents()

	var children []layout.FlexChild
	var padding float32 = 0
	if c.withOptions {
		children = c.layoutCollapsibleWithOption(contents)
		padding = 10
	} else {
		children = c.layoutCollapsible(contents)
	}

	return c.common.card.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(padding)).Layout(gtx, func(C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})
	})
}

func (c *Collapsible) MoreTriggered() bool {
	return c.withOptions && c.moreIconButton.Button.Clicked()
}

func (c *Collapsible) handleEvents() {
	if c.common.isExpanded {
		c.common.currentIcon = c.common.expandedIcon
	} else {
		c.common.currentIcon = c.common.collapsedIcon
	}

	for c.common.button.Clicked() {
		c.common.isExpanded = !c.common.isExpanded
	}
}
