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
	theme      *decredmaterial.Theme
	debugItems []debugItem
	common     pageCommon
}

func (win *Window) DebugPage(common pageCommon) Page {
	debugItems := []debugItem{
		{
			clickable: new(widget.Clickable),
			text:      "Check wallet logs",
			page:      PageLog,
		},
	}

	pg := &debugPage{
		theme:      common.theme,
		debugItems: debugItems,
		common:     common,
	}

	return pg
}

func (pg *debugPage) handle() {
	for i := range pg.debugItems {
		for pg.debugItems[i].clickable.Clicked() {
			pg.common.changePage(pg.debugItems[i].page)
		}
	}
}

func (pg *debugPage) onClose() {}

func (pg *debugPage) debugItem(gtx C, i int, common pageCommon) D {
	return decredmaterial.Clickable(gtx, pg.debugItems[i].clickable, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return common.theme.Body1(pg.debugItems[i].text).Layout(gtx)
				})
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						return common.icons.chevronRight.Layout(gtx, values.MarginPadding22)
					})
				})
			}),
		)
	})
}

func (pg *debugPage) layoutDebugItems(gtx C, common pageCommon) {
	background := common.theme.Color.Surface
	card := common.theme.Card()
	card.Color = background
	card.Layout(gtx, func(gtx C) D {
		list := layout.List{Axis: layout.Vertical}
		return list.Layout(gtx, len(pg.debugItems), func(gtx C, i int) D {
			return pg.debugItem(gtx, i, common)
		})
	})
}

func (pg *debugPage) Layout(gtx C) D {
	container := func(gtx C) D {
		page := SubPage{
			title: "Debug",
			back: func() {
				pg.common.changePage(PageMore)
			},
			body: func(gtx C) D {
				pg.layoutDebugItems(gtx, pg.common)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			},
		}
		return pg.common.SubPageLayout(gtx, page)

	}
	return pg.common.Layout(gtx, func(gtx C) D {
		return pg.common.UniformPadding(gtx, container)
	})
}
