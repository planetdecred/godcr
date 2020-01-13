package widgets

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/helper"
)

type (
	Selectable struct {
		selected string
		items    []*Button
	}
)

func NewSelectable(items []string) *Selectable {
	btns := make([]*Button, len(items))
	for i := range items {
		btns[i] = NewButton(items[i], nil).
			SetColor(helper.GrayColor).
			SetBackgroundColor(helper.WhiteColor).
			SetBorderColor(helper.GrayColor)
	}

	return &Selectable{
		items: btns,
	}
}

func (s *Selectable) Select(index int) {
	for i := range s.items {
		if index == i {
			s.selected = s.items[i].text
			break
		}
	}
}

func (s *Selectable) SelectText(txt string) {
	s.selected = txt
}

func (s *Selectable) Selected() string {
	return s.selected
}

func (s *Selectable) Draw(ctx *layout.Context) {
	numItems := len(s.items)
	spacing := 3
	children := make([]layout.FlexChild, numItems)
	tSpacing := spacing * numItems
	width := (ctx.Constraints.Width.Max / numItems) - tSpacing

	for i := range s.items {
		index := i
		children[i] = layout.Rigid(func() {
			ctx.Constraints.Width.Min = width
			sideInset := float32(tSpacing / 2)

			inset := layout.Inset{
				Left:  unit.Dp(sideInset),
				Right: unit.Dp(sideInset),
			}
			inset.Layout(ctx, func() {
				s.items[index].Draw(ctx, func() {
					s.items[index].SetColor(helper.DecredLightBlueColor).SetBorderColor(helper.DecredLightBlueColor)
					s.selected = s.items[index].text

					for i := range s.items {
						tindex := i
						if tindex != index {
							s.items[tindex].SetColor(helper.GrayColor).SetBorderColor(helper.GrayColor)
						}
					}
				})
			})
		})
	}

	layout.Stack{}.Layout(ctx,
		layout.Expanded(func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(ctx, children...)
		}),
	)
}
