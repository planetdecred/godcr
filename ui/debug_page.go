package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageDebug = "Debug"

type debugItem struct {
	clickable *widget.Clickable
	text      string
	page      string
}

type debugPage struct {
	theme *decredmaterial.Theme

	pageTitle decredmaterial.Label

	debugItems []debugItem
}

func (win *Window) DebugPage(common pageCommon) layout.Widget {
	debugItems := []debugItem{
		{
			clickable: new(widget.Clickable),
			text:      "Check wallet logs",
			page:      PageLog,
		},
	}

	pg := &debugPage{
		theme:      common.theme,
		pageTitle:  common.theme.H5("Session log entries"),
		debugItems: debugItems,
	}

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *debugPage) handle(common pageCommon) {
	for i := range pg.debugItems {
		for pg.debugItems[i].clickable.Clicked() {
			*common.page = pg.debugItems[i].page
		}
	}
}

func (pg *debugPage) debugItem(gtx C, i int, common pageCommon) D {
	return decredmaterial.Clickable(gtx, pg.debugItems[i].clickable, func(gtx C) D {
		background := common.theme.Color.Surface
		card := common.theme.Card()
		card.Color = background
		return card.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return common.theme.Body1(pg.debugItems[i].text).Layout(gtx)
					})
				}))
		})
	})
}

func (pg *debugPage) layoutDebugItems(gtx C, common pageCommon) {
	layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(pg.debugItems), func(gtx C, i int) D {
				return pg.debugItem(gtx, i, common)
			})
		}))
}

// main settings layout
func (pg *debugPage) Layout(gtx C, common pageCommon) D {
	container := func(gtx C) D {
		page := SubPage{
			title: "Debug",
			back: func() {
				*common.page = PageMore
			},
			body: func(gtx C) D {
				pg.layoutDebugItems(gtx, common)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			},
		}
		return common.SubPageLayout(gtx, page)

	}
	return common.Layout(gtx, container)
}
