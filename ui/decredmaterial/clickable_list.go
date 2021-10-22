package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type ClickableList struct {
	// layout.List
	theme          *Theme
	clickables     []*Clickable
	ClickHighlight color.NRGBA
	Radius         CornerRadius // this radius is used by the clickable
	selectedItem   int
	DividerHeight  unit.Value
	list           *widget.List
}

func (t *Theme) NewClickableList(axis layout.Axis) *ClickableList {
	click := &ClickableList{
		theme: t,
		// List:              layout.List{Axis: axis},
		ClickHighlight: t.Color.SurfaceHighlight,
		selectedItem:   -1,
	}

	click.list = &widget.List{
		List: layout.List{
			Axis: axis,
		},
	}
	return click
}

func (cl *ClickableList) ItemClicked() (bool, int) {
	defer func() {
		cl.selectedItem = -1
	}()
	return cl.selectedItem != -1, cl.selectedItem
}

func (cl *ClickableList) handleClickables(count int) {
	if len(cl.clickables) != count {

		cl.clickables = make([]*Clickable, count)
		for i := 0; i < count; i++ {
			clickable := cl.theme.NewClickable(true)
			clickable.Color = cl.ClickHighlight
			cl.clickables[i] = clickable
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
	return cl.theme.List(cl.list).Layout(gtx, count, func(gtx C, i int) D {
		return cl.row(gtx, count, i, w)
	})
}

func (cl *ClickableList) row(gtx layout.Context, count int, i int, w layout.ListElement) layout.Dimensions {
	if i == 0 { // first item
		cl.clickables[i].Radius.TopLeft = cl.Radius.TopLeft
		cl.clickables[i].Radius.TopRight = cl.Radius.TopRight
	}
	if i == count-1 { // last item
		cl.clickables[i].Radius.BottomLeft = cl.Radius.BottomLeft
		cl.clickables[i].Radius.BottomRight = cl.Radius.BottomRight
	}
	row := cl.clickables[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return w(gtx, i)
	})

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
}
