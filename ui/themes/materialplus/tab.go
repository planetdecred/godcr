package materialplus

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui"
)

type (
	// Tab represents a tab item
	Tab struct {
		ID         int32
		Label      string
		RenderFunc func(*layout.Context)

		btnWidget *widget.Button
		trigger   material.Button
	}

	// TabContainer represents a tab container
	TabContainer struct {
		title           material.Label
		items           []Tab
		currentTabIndex int
	}
)

const (
	tabNavHeight = 60
)

// TabContainer initializes an instance of TabContainer
func (t *Theme) TabContainer(items []Tab) *TabContainer {
	for i := range items {
		btn := t.Button(items[i].Label)
		btn.Background = ui.WhiteColor

		items[i].btnWidget = new(widget.Button)
		items[i].trigger = btn
	}

	return &TabContainer{
		title:           t.H6(""),
		items:           items,
		currentTabIndex: 0,
	}
}

// Layout renders the tabcontainer to screen
func (t *TabContainer) Layout(gtx *layout.Context, title string) {
	w := []func(){
		func() {
			t.title.Text = title
			t.title.Layout(gtx)
		},
		func() {
			t.drawNavSection(gtx)
		},
		func() {
			t.drawContentSection(gtx)
		},
	}

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(w), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, w[i])
	})
}

func (t *TabContainer) drawNavSection(gtx *layout.Context) {
	navWidth := gtx.Constraints.Width.Max / len(t.items)
	columns := make([]layout.FlexChild, len(t.items))

	for index := range t.items {
		i := index

		if t.currentTabIndex == i {
			t.items[i].trigger.Color = ui.LightBlueColor
		} else {
			t.items[i].trigger.Color = ui.BlackColor
		}

		columns[i] = layout.Rigid(func() {
			gtx.Constraints.Width.Min = navWidth
			for t.items[i].btnWidget.Clicked(gtx) {
				t.currentTabIndex = i
			}

			t.items[i].trigger.Layout(gtx, t.items[i].btnWidget)
		})
	}

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx, columns...)
		}),
	)
}

func (t *TabContainer) drawContentSection(gtx *layout.Context) {
	t.items[t.currentTabIndex].RenderFunc(gtx)
}

// GetCurrentTabID returns the id of the open tab
func (t *TabContainer) GetCurrentTabID() int32 {
	return t.items[t.currentTabIndex].ID
}
