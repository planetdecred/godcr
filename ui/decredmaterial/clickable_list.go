package decredmaterial

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type ClickableList struct {
	layout.List
	clickables   []*widget.Clickable
	selectedItem int
}

func (t *Theme) NewClickableList(axis layout.Axis) *ClickableList {
	return &ClickableList{
		List:         layout.List{Axis: layout.Vertical},
		selectedItem: -1,
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
		return Clickable(gtx, cl.clickables[i], func(gtx layout.Context) layout.Dimensions {
			return w(gtx, i)
		})

	})
}
