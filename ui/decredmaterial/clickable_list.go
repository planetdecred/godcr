package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type ClickableList struct {
	layout.List
	theme              *Theme
	clickables         []*widget.Clickable
	selectedItem       int
	ClickableHighlight bool
	DividerHeight      unit.Value
}

func (t *Theme) NewClickableList(axis layout.Axis) *ClickableList {
	return &ClickableList{
		theme:              t,
		List:               layout.List{Axis: axis},
		ClickableHighlight: true,
		selectedItem:       -1,
	}
}

func (cl *ClickableList) ItemClicked() (bool, int) {
	defer func() {
		cl.selectedItem = -1
	}()
	return cl.selectedItem != -1, cl.selectedItem
}

func (cl *ClickableList) handleClickables(count int) {
	if len(cl.clickables) != count {

		cl.clickables = make([]*widget.Clickable, count)
		for i := 0; i < count; i++ {
			cl.clickables[i] = new(widget.Clickable)
		}
	}

	for index, clickable := range cl.clickables {
		for clickable.Clicked() {
			cl.selectedItem = index
		}
	}
}

func (cl *ClickableList) Layout(gtx layout.Context, count int, w layout.ListElement) layout.Dimensions {
	cl.handleClickables(count)
	return cl.List.Layout(gtx, count, func(gtx layout.Context, i int) layout.Dimensions {
		var row layout.Dimensions = w(gtx, i)
		if cl.ClickableHighlight {
			row = Clickable(gtx, cl.clickables[i], func(gtx layout.Context) layout.Dimensions {
				return row
			})
		}

		// add divider to all rows except last
		if i < (count-1) && cl.DividerHeight.V > 0 {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return row
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y += gtx.Px(cl.DividerHeight)

					return layout.Dimensions{Size: gtx.Constraints.Min}
				}),
			)
		}
		return row
	})
}
