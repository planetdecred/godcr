package decredmaterial

import (
	// "fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Collapsible struct {
	IsExpanded            bool
	Button                *widget.Clickable
	expandIcon            *widget.Icon
	headerBackgroundColor color.RGBA
}

func (t *Theme) Collapsible(button *widget.Clickable) *Collapsible {
	c := &Collapsible{
		headerBackgroundColor: t.Color.Hint,
		expandIcon:            t.navMoreIcon,
		Button:                button,
	}

	return c
}

func (c *Collapsible) layoutHeader(gtx layout.Context, header func(C) D) layout.Dimensions {
	dims := layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return header(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: unit.Dp(20)}.Layout(gtx, func(C) D {
				return c.expandIcon.Layout(gtx, unit.Dp(20))
			})
		}),
	)

	return dims
}

func (c *Collapsible) Layout(gtx layout.Context, header func(C) D, content func(C) D) layout.Dimensions {
	for c.Button.Clicked() {
		c.IsExpanded = !c.IsExpanded
	}

	dims := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return c.layoutHeader(gtx, header)
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
	return dims
}
