package materialplus

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui"
)

type (
	tabItem struct {
		label      string
		btnWidget  *widget.Button
		renderFunc func(*layout.Context)
	}

	// Tab represents a tab item
	Tab struct {
		Label      string
		RenderFunc func(*layout.Context)

		btnWidget *widget.Button
		trigger   material.Button
	}

	// TabContainer represents a tab container
	TabContainer struct {
		items           []Tab
		currentTabIndex int
	}
)

const (
	tabNavHeight = 60
)

// TabContainer returns a tabcontainer instance
func (t *Theme) TabContainer(items []Tab) *TabContainer {
	for i := range items {
		btn := t.Button(items[i].Label)
		btn.Background = ui.WhiteColor

		items[i].btnWidget = new(widget.Button)
		items[i].trigger = btn
	}

	return &TabContainer{
		items:           items,
		currentTabIndex: 0,
	}
}

// Draw renders the tabcontainer to screen
func (t *TabContainer) Draw(gtx *layout.Context) {
	w := []func(){
		func() {
			t.drawNavSection(gtx)
		},
		func() {
			t.drawContentSection(gtx)
		},
	}

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, 2, func(i int) {
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

// GetCurrentTabLabel returns the label of the open tab
func (t *TabContainer) GetCurrentTabLabel() string {
	return t.items[t.currentTabIndex].Label
}
