package widgets

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr-gio/helper"
)

type (
	tabItem struct {
		label      *ClickableLabel
		renderFunc func()
	}

	TabContainer struct {
		items           []tabItem
		currentTabIndex int
	}
)

func NewTabContainer() *TabContainer {
	return &TabContainer{
		items:           []tabItem{},
		currentTabIndex: 0,
	}
}

func (t *TabContainer) AddTab(label string) *TabContainer {
	lbl := NewClickableLabel(label).
		SetSize(5).
		SetAlignment(AlignMiddle).
		SetWeight(text.Bold)

	item := tabItem{
		label:      lbl,
		renderFunc: func() {},
	}
	t.items = append(t.items, item)
	return t
}

func (t *TabContainer) Draw(ctx *layout.Context, renderFuncs ...func(*layout.Context)) {
	t.drawNavSection(ctx)

	// draw current tab content
	inset := layout.Inset{
		Top: unit.Dp(25),
	}
	inset.Layout(ctx, func() {
		// TODO make sure number of render funcs match number of tabs
		renderFuncs[t.currentTabIndex](ctx)
	})
}

func (t *TabContainer) drawNavSection(ctx *layout.Context) {
	navTabWidth := ctx.Constraints.Width.Max / len(t.items)
	columns := make([]layout.FlexChild, len(t.items))

	for index := range t.items {
		tindex := index

		color := helper.BlackColor
		if t.currentTabIndex == tindex {
			color = helper.DecredLightBlueColor
		}

		columns[tindex] = layout.Rigid(func() {
			t.items[tindex].label.SetWidth(navTabWidth).
				SetColor(color).
				Draw(ctx, func() {
					t.currentTabIndex = tindex
				})
		})
	}

	layout.Stack{}.Layout(ctx,
		layout.Expanded(func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(ctx, columns...)
		}),
	)
}

// func (t *TabContainer) drawContentSection(ctx *layout.Context, renderFuncs []func(*layout.Context)) {
// 	inset := layout.Inset{
// 		Top: unit.Dp(25),
// 	}
// 	inset.Layout(ctx, func() {

// 	})
// }
